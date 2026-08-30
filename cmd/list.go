package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

type BinInfo struct {
	ID                string `json:"id"`
	ReadOnly          bool   `json:"readonly"`
	Bytes             int64  `json:"bytes"`
	BytesReadable     string `json:"bytes_readable"`
	Files             int    `json:"files"`
	UpdatedAt         string `json:"updated_at"`
	UpdatedAtRelative string `json:"updated_at_relative"`
	CreatedAt         string `json:"created_at"`
	CreatedAtRelative string `json:"created_at_relative"`
	ExpiredAt         string `json:"expired_at"`
	ExpiredAtRelative string `json:"expired_at_relative"`
}

type FileInfo struct {
	Filename          string `json:"filename"`
	ContentType       string `json:"content-type"`
	Bytes             int64  `json:"bytes"`
	BytesReadable     string `json:"bytes_readable"`
	MD5               string `json:"md5"`
	SHA256            string `json:"sha256"`
	UpdatedAt         string `json:"updated_at"`
	UpdatedAtRelative string `json:"updated_at_relative"`
	CreatedAt         string `json:"created_at"`
	CreatedAtRelative string `json:"created_at_relative"`
}

type BinResponse struct {
	Bin   BinInfo    `json:"bin"`
	Files []FileInfo `json:"files"`
}

var listCmd = &cobra.Command{
	Use:   "list <bin>",
	Short: "List bin contents",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bin := args[0]
		useJSON, _ := cmd.Flags().GetBool("json")

		url := fmt.Sprintf("%s/%s", baseURL, bin)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("list failed (%d): %s", resp.StatusCode, string(body))
		}

		var binResp BinResponse
		if err := json.Unmarshal(body, &binResp); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}

		if useJSON {
			fmt.Println(string(body))
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "FILE\tSIZE\tTYPE\tCREATED\n")
		for _, f := range binResp.Files {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", f.Filename, f.BytesReadable, f.ContentType, f.CreatedAtRelative)
		}
		w.Flush()

		fmt.Fprintf(os.Stderr, "\nBin: %s | Files: %d | Size: %s | Expires: %s\n",
			binResp.Bin.ID, binResp.Bin.Files, binResp.Bin.BytesReadable, binResp.Bin.ExpiredAtRelative)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().Bool("json", false, "Output raw JSON")
}
