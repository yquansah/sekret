package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/yquansah/sekret/pkg/k8s"

	"github.com/spf13/cobra"
)

var (
	listNamespace string
)

type SecretItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SecretData struct {
	Data []SecretItem `json:"data"`
}

var listCmd = &cobra.Command{
	Use:   "list [secret-name]",
	Short: "List the contents of a Kubernetes secret in JSON format",
	Long: `List retrieves and displays the contents of a Kubernetes secret
in JSON format with the secret data as an array of key-value pairs.

Example:
  sekret list my-secret --namespace default`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		secretName := args[0]

		k8sClient, err := k8s.NewClient(kubeconfig)
		if err != nil {
			return fmt.Errorf("failed to create Kubernetes client: %w", err)
		}

		secretData, err := k8sClient.GetSecretData(ctx, secretName, listNamespace)
		if err != nil {
			return fmt.Errorf("failed to get secret data: %w", err)
		}

		var items []SecretItem
		for key, value := range secretData {
			items = append(items, SecretItem{
				Key:   key,
				Value: value,
			})
		}

		output := SecretData{Data: items}
		jsonOutput, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}

		fmt.Println(string(jsonOutput))
		return nil
	},
}

func init() {
	listCmd.Flags().StringVarP(&listNamespace, "namespace", "n", "default", "Kubernetes namespace")
}