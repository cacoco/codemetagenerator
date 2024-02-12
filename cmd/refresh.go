package cmd

import (
	"net/http"

	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/spf13/cobra"
)

func refresh(writer utils.Writer, basedir string, httpClient *http.Client) error {
	// update licenses file
	err := downloadSPDXLicenses(basedir, httpClient, true)
	if err != nil {
		handleErr(writer, err)
		return writer.Errorf("unable to update SPDX licenses file: %s", err.Error())
	}
	writer.Println("✅ Successfully updated SPDX licenses file.")
	return nil
}

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh the list of current SPDX IDs (https://spdx.org/licenses/)",
	Long: `
Use this command to refresh the stored list of currently supported SPDX IDs from https://spdx.org/licenses/. 

	See: https://spdx.dev/learn/handling-license-info/#why`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return refresh(&utils.StdoutWriter{}, utils.UserHomeDir, utils.MkHttpClient())
	},
}

func init() {
	licensesCmd.AddCommand(refreshCmd)
}
