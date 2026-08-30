package cmd

import (
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <bin> [file]",
	Short: "Delete a file or an entire bin",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bin := args[0]

		var url string
		var msg string
		if len(args) == 2 {
			file := args[1]
			url = fmt.Sprintf("%s/%s/%s", baseURL, bin, file)
			msg = fmt.Sprintf("Deleted file %s from bin %s", file, bin)
		} else {
			url = fmt.Sprintf("%s/%s", baseURL, bin)
			msg = fmt.Sprintf("Deleted bin %s", bin)
		}

		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("delete failed (%d): %s", resp.StatusCode, string(body))
		}

		fmt.Println(msg)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
