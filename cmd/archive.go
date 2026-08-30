package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var archiveCmd = &cobra.Command{
	Use:   "archive <bin>",
	Short: "Download bin as tar or zip archive",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bin := args[0]
		format, _ := cmd.Flags().GetString("format")
		outPath, _ := cmd.Flags().GetString("output")

		if format == "" {
			format = "tar"
		}

		var url string
		var defaultOut string
		switch format {
		case "tar":
			url = fmt.Sprintf("%s/archive/%s/tar", baseURL, bin)
			defaultOut = bin + ".tar"
		case "zip":
			url = fmt.Sprintf("%s/archive/%s/zip", baseURL, bin)
			defaultOut = bin + ".zip"
		default:
			return fmt.Errorf("unsupported format %q (use tar or zip)", format)
		}

		if outPath == "" {
			outPath = defaultOut
		}

		jar, _ := cookiejar.New(nil)
		client := &http.Client{Jar: jar}

		resp, err := client.Get(url)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}
		resp.Body.Close()

		contentType := resp.Header.Get("Content-Type")
		if strings.Contains(contentType, "text/html") {
			resp, err = client.Get(url)
			if err != nil {
				return fmt.Errorf("archive failed: %w", err)
			}
		}

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return fmt.Errorf("archive failed (%d): %s", resp.StatusCode, string(respBody))
		}

		f, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer f.Close()

		if _, err := io.Copy(f, resp.Body); err != nil {
			return fmt.Errorf("failed to write archive: %w", err)
		}
		resp.Body.Close()

		fmt.Printf("Saved %s archive to %s\n", format, outPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(archiveCmd)
	archiveCmd.Flags().StringP("format", "f", "", "Archive format: tar or zip (default: tar)")
	archiveCmd.Flags().StringP("output", "o", "", "Output file path")
}
