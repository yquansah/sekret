package cmd

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/yquansah/sekret/pkg/k8s"

	"github.com/spf13/cobra"
)

var (
	deleteNamespace string
	keysToDelete    string
)

var deleteKeysCmd = &cobra.Command{
	Use:   "delete-keys [secret-name]",
	Short: "Delete specific keys from a Kubernetes secret",
	Long: `Delete specific keys from a Kubernetes secret either interactively
or by specifying keys with the --keys flag.

Without --keys flag: Interactive mode will list all keys and prompt for selection
With --keys flag: Non-interactive mode will delete the specified comma-separated keys

Examples:
  sekret delete-keys my-secret --namespace default
  sekret delete-keys my-secret --namespace default --keys key1,key2,key3`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		secretName := args[0]

		k8sClient, err := k8s.NewClient(kubeconfig)
		if err != nil {
			return fmt.Errorf("failed to create Kubernetes client: %w", err)
		}

		if keysToDelete != "" {
			// Non-interactive mode
			keys := strings.Split(keysToDelete, ",")
			for i, key := range keys {
				keys[i] = strings.TrimSpace(key)
			}

			deletedCount, err := k8sClient.DeleteKeysFromSecret(ctx, secretName, deleteNamespace, keys)
			if err != nil {
				return fmt.Errorf("failed to delete keys: %w", err)
			}

			fmt.Printf("Successfully deleted %d key(s) from secret '%s' in namespace '%s'\n",
				deletedCount, secretName, deleteNamespace)
			return nil
		}

		// Interactive mode
		secretData, err := k8sClient.GetSecretData(ctx, secretName, deleteNamespace)
		if err != nil {
			return fmt.Errorf("failed to get secret data: %w", err)
		}

		if len(secretData) == 0 {
			fmt.Printf("Secret '%s' in namespace '%s' has no keys to delete\n", secretName, deleteNamespace)
			return nil
		}

		// Create list of keys for selection
		var keys []string
		for key := range secretData {
			keys = append(keys, key)
		}

		// Prompt for key selection
		prompt := promptui.Select{
			Label: fmt.Sprintf("Select key to delete from secret '%s'", secretName),
			Items: keys,
		}

		_, selectedKey, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("selection cancelled: %w", err)
		}

		// Confirm deletion
		confirmPrompt := promptui.Prompt{
			Label:     fmt.Sprintf("Are you sure you want to delete key '%s' from secret '%s'", selectedKey, secretName),
			IsConfirm: true,
		}

		_, err = confirmPrompt.Run()
		if err != nil {
			fmt.Println("Deletion cancelled")
			return nil
		}

		// Delete the selected key
		deletedCount, err := k8sClient.DeleteKeysFromSecret(ctx, secretName, deleteNamespace, []string{selectedKey})
		if err != nil {
			return fmt.Errorf("failed to delete key: %w", err)
		}

		fmt.Printf("Successfully deleted %d key(s) from secret '%s' in namespace '%s'\n",
			deletedCount, secretName, deleteNamespace)
		return nil
	},
}

func init() {
	deleteKeysCmd.Flags().StringVarP(&deleteNamespace, "namespace", "n", "default", "Kubernetes namespace")
	deleteKeysCmd.Flags().StringVar(&keysToDelete, "keys", "", "Comma-separated list of keys to delete (non-interactive mode)")
}