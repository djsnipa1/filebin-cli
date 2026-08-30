package cmd

import (
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

var checksumsCmd = &cobra.Command{
	Use:   "checksums <bin>",
	Short: "List SHA256 checksums of files in a bin",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bin := args[0]
		url := fmt.Sprintf("%s/sha256/%s", baseURL, bin)

		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("checksums failed (%d): %s", resp.StatusCode, string(body))
		}

		fmt.Print(string(body))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(checksumsCmd)
}
