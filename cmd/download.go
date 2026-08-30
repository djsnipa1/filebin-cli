package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download <bin> <file>",
	Short: "Download a file from a bin",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bin := args[0]
		filename := args[1]

		outPath, _ := cmd.Flags().GetString("output")
		if outPath == "" {
			outPath = filename
		}

		jar, _ := cookiejar.New(nil)
		client := &http.Client{Jar: jar}

		url := fmt.Sprintf("%s/%s/%s", baseURL, bin, filename)

		resp, err := client.Get(url)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}
		resp.Body.Close()

		if resp.StatusCode == http.StatusFound {
			location := resp.Header.Get("Location")
			resp, err = client.Get(location)
			if err != nil {
				return fmt.Errorf("download failed: %w", err)
			}
			resp.Body.Close()
		}

		contentType := resp.Header.Get("Content-Type")
		if strings.Contains(contentType, "text/html") {
			resp, err = client.Get(url)
			if err != nil {
				return fmt.Errorf("download failed: %w", err)
			}
		}

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return fmt.Errorf("download failed (%d): %s", resp.StatusCode, string(respBody))
		}

		f, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer f.Close()

		bar := progressbar.NewOptions64(
			resp.ContentLength,
			progressbar.OptionSetDescription("Downloading"),
			progressbar.OptionSetWriter(os.Stderr),
			progressbar.OptionShowBytes(true),
			progressbar.OptionSetWidth(40),
			progressbar.OptionClearOnFinish(),
			progressbar.OptionSpinnerType(14),
		)

		if _, err := io.Copy(io.MultiWriter(f, bar), resp.Body); err != nil {
			return fmt.Errorf("download failed: %w", err)
		}
		resp.Body.Close()
		bar.Finish()

		absPath, _ := filepath.Abs(outPath)
		fmt.Fprintf(os.Stderr, "\nDownloaded to %s\n", absPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringP("output", "o", "", "Output file path")
}
