package cmd

import (
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

var lockCmd = &cobra.Command{
	Use:   "lock <bin>",
	Short: "Lock a bin (read-only)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bin := args[0]
		url := fmt.Sprintf("%s/%s", baseURL, bin)

		req, err := http.NewRequest("PUT", url, nil)
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
			return fmt.Errorf("lock failed (%d): %s", resp.StatusCode, string(body))
		}

		fmt.Printf("Locked bin %s\n", bin)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lockCmd)
}
