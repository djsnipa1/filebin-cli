package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var baseURL string

var rootCmd = &cobra.Command{
	Use:   "filebin",
	Short: "CLI client for filebin.net",
	Long:  "Upload, download, list, and manage files on filebin.net",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "https://filebin.net", "Filebin server URL")
}
