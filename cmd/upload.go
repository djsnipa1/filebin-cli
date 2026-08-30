package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

type barWriter struct {
	bar *progressbar.ProgressBar
	r   io.Reader
}

func (bw *barWriter) Read(p []byte) (n int, err error) {
	n, err = bw.r.Read(p)
	bw.bar.Add64(int64(n))
	return
}

var uploadCmd = &cobra.Command{
	Use:   "upload <bin> <file>",
	Short: "Upload a file to a bin",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bin := args[0]
		filePath := args[1]

		f, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer f.Close()

		stat, err := f.Stat()
		if err != nil {
			return fmt.Errorf("failed to stat file: %w", err)
		}

		hasher := sha256.New()
		if _, err := io.Copy(hasher, f); err != nil {
			return fmt.Errorf("failed to compute checksum: %w", err)
		}
		sha256Hex := hex.EncodeToString(hasher.Sum(nil))

		if _, err := f.Seek(0, io.SeekStart); err != nil {
			return fmt.Errorf("failed to seek file: %w", err)
		}

		bar := progressbar.NewOptions64(
			stat.Size(),
			progressbar.OptionSetDescription("Uploading"),
			progressbar.OptionSetWriter(os.Stderr),
			progressbar.OptionShowBytes(true),
			progressbar.OptionSetWidth(40),
			progressbar.OptionClearOnFinish(),
			progressbar.OptionSpinnerType(14),
		)
		writer := &barWriter{bar: bar, r: f}

		filename := filepath.Base(filePath)
		url := fmt.Sprintf("%s/%s/%s", baseURL, bin, filename)

		req, err := http.NewRequest("POST", url, writer)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		req.ContentLength = stat.Size()
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("Content-SHA256", sha256Hex)

		client := &http.Client{Timeout: 0}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("upload failed: %w", err)
		}
		defer resp.Body.Close()

		bar.Finish()

		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("upload failed (%d): %s", resp.StatusCode, string(body))
		}

		fmt.Fprintf(os.Stderr, "\nUploaded %s to bin %s\n", filename, bin)
		fmt.Fprintf(os.Stderr, "URL: %s/%s/%s\n", baseURL, bin, filename)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
