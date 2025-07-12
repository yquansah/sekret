package cmd

import (
	"fmt"

	"sekret/pkg/dotenv"
	"sekret/pkg/k8s"

	"github.com/spf13/cobra"
)

var (
	namespace string
	envFile   string
	replace   bool
)

var upsertCmd = &cobra.Command{
	Use:   "upsert [secret-name]",
	Short: "Create or update a Kubernetes secret from a dotenv file",
	Long: `Upsert creates or updates a Kubernetes secret using key-value pairs
from a dotenv file. Values are automatically base64 encoded for storage
in the secret.

Example:
  sekret upsert my-secret --namespace default --env-file .env`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		secretName := args[0]

		envVars, err := dotenv.LoadFromFile(envFile)
		if err != nil {
			return fmt.Errorf("failed to load environment file: %w", err)
		}

		if len(envVars) == 0 {
			return fmt.Errorf("no environment variables found in %s", envFile)
		}

		k8sClient, err := k8s.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create Kubernetes client: %w", err)
		}

		keysModified, err := k8sClient.UpsertSecret(ctx, secretName, namespace, envVars, replace)
		if err != nil {
			return fmt.Errorf("failed to update secret: %w", err)
		}

		if replace {
			fmt.Printf("Successfully replaced %d key(s) for secret '%s' in namespace '%s'\n",
				keysModified, secretName, namespace)
		} else {
			fmt.Printf("Successfully upserted %d key(s) for secret '%s' in namespace '%s'\n",
				keysModified, secretName, namespace)
		}

		return nil
	},
}

func init() {
	upsertCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes namespace")
	upsertCmd.Flags().StringVarP(&envFile, "env-file", "f", ".env", "Path to dotenv file")
	upsertCmd.Flags().BoolVar(&replace, "replace", false, "Replace all existing values in secret instead of merging")
}
