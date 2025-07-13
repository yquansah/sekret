package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig string
)

var rootCmd = &cobra.Command{
	Use:   "sekret",
	Short: "A tool for managing Kubernetes secrets",
	Long: `sekret is a CLI tool that simplifies Kubernetes secret management.
It can create or update secrets from dotenv files with automatic base64 encoding,
list existing secret contents in JSON format, and delete specific keys from secrets.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", clientcmd.RecommendedHomeFile, "Path to kubeconfig file")
	rootCmd.AddCommand(upsertCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(deleteKeysCmd)
}