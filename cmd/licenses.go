package cmd

import (
	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/spf13/cobra"
)

func licenses(writer utils.Writer) ([]string, error) {
	supportedLicenses := SupportedLicenses.getSupportedLicenses()
	// list licenses
	var list []string = make([]string, 0)
	for _, license := range supportedLicenses {
		list = append(list, license)
		writer.Println(license)
	}
	return list, nil
}

// licensesCmd represents the licenses command
var licensesCmd = &cobra.Command{
	Use:   "licenses [refresh]",
	Args:  cobra.NoArgs,
	Short: "List current SPDX IDs (https://spdx.org/licenses/)",
	Long: `
Use this command to list currently supported SPDX IDs from https://spdx.org/licenses/. 

This is a long list and as such you may want to pipe the output into "more" or 
"less" to view it. See: https://spdx.dev/learn/handling-license-info/#why`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := licenses(&utils.StdoutWriter{})
		return err
	},
}

func init() {
	rootCmd.AddCommand(licensesCmd)
}
