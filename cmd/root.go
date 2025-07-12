package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sekret",
	Short: "A tool for managing Kubernetes secrets from dotenv files",
	Long: `sekret is a CLI tool that simplifies Kubernetes secret management
by reading environment variables from dotenv files and automatically
creating or updating secrets in your cluster with proper base64 encoding.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(upsertCmd)
}