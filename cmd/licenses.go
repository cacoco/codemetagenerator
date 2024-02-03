package cmd

import (
	"fmt"
	"strings"

	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/spf13/cobra"
)

func validateLicenseId(basedir string) func(string) error {
	return func(id string) error {
		var supported []string
		if SupportedLicenses == nil {
			keys, err := getSupportedLicenses(basedir)
			if err != nil {
				return fmt.Errorf("unable to retrieve supported licenses: %s", err.Error())
			}
			supported = *keys
		} else {
			supported = SupportedLicenses
		}
		for _, license := range supported {
			if license == id {
				return nil
			}
		}
		return fmt.Errorf("Invalid license ID: " + id)
	}
}

func getSupportedLicenses(basedir string) (*[]string, error) {
	return utils.GetSupportedLicenses(basedir)
}

func listLicenses(basedir string, supported []string, writer utils.Writer) error {
	// list licenses
	for _, license := range supported {
		writer.Println(license)
	}
	return nil
}

func updateLicenses(basedir string, writer utils.Writer) error {
	err := utils.GetAndCacheLicenseFile(utils.UserHomeDir, true)
	if err != nil {
		return err
	}
	writer.Println("✅ Successfully updated SPDX licenses file.")
	return nil
}

var SupportedLicenses []string

// licensesCmd represents the licenses command
var licensesCmd = &cobra.Command{
	Use:   "licenses [refresh]",
	Args:  cobra.RangeArgs(0, 1),
	Short: "List (or refresh cached) SPDX license IDs",
	Long: `Use this command to list license SPDX ids from the https://spdx.org/licenses/. 
This is a long list and as such you may want to pipe the output into "more" or 
"less" to view it. Pass the "refresh" argument to update the cached list of licenses.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		foundLicenses, err := getSupportedLicenses(utils.UserHomeDir)
		if err != nil {
			return fmt.Errorf("unable to retrieve supported licenses: %s", err.Error())
		}
		SupportedLicenses = *foundLicenses
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// list licenses
			return listLicenses(utils.UserHomeDir, SupportedLicenses, &utils.StdoutWriter{})
		} else if len(args) == 1 && args[0] == "refresh" {
			// update licenses file
			return updateLicenses(utils.UserHomeDir, &utils.StdoutWriter{})
		} else {
			return fmt.Errorf("unrecognized argument(s): '%s'", strings.Join(args, " "))
		}
	},
}

func init() {
	rootCmd.AddCommand(licensesCmd)
}
