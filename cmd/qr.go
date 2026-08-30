package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var qrCmd = &cobra.Command{
	Use:   "qr <bin>",
	Short: "Download QR code for a bin",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bin := args[0]
		outPath, _ := cmd.Flags().GetString("output")

		if outPath == "" {
			outPath = bin + "_qr.png"
		}

		url := fmt.Sprintf("%s/qr/%s", baseURL, bin)

		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("QR generation failed (%d): %s", resp.StatusCode, string(body))
		}

		f, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer f.Close()

		if _, err := io.Copy(f, resp.Body); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		fmt.Printf("QR code saved to %s\n", outPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(qrCmd)
	qrCmd.Flags().StringP("output", "o", "", "Output file path")
}
