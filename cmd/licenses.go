package cmd

import (
	"fmt"
	"strings"

	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/ohler55/ojg/oj"
	"github.com/spf13/cobra"
)

func validateLicenseId(basedir string) func(string) error {
	return func(id string) error {
		supportedLicenses := SupportedLicenses.getSupportedLicenses()
		if supportedLicenses == nil {
			return fmt.Errorf("SPDX licenses have not be downloaded, please run `codemeta licenses refresh` to download the SPDX licenses")
		}
		for _, license := range supportedLicenses {
			if license == id {
				return nil
			}
		}
		return fmt.Errorf("invalid SPDX license ID: " + id + ", see: https://spdx.org/licenses/ for a list of valid values")
	}
}

func getLicenseReference(basedir string, id string) (*string, error) {
	bytes, err := utils.LoadFile(utils.GetLicensesFilePath(basedir))
	if err != nil {
		return nil, err
	}
	var licenses map[string]string
	err = oj.Unmarshal(bytes, &licenses)

	if err != nil {
		return nil, err
	}

	reference, ok := licenses[id]
	if !ok {
		return nil, fmt.Errorf("Invalid license ID: " + id)
	}
	return &reference, nil
}

// licensesCmd represents the licenses command
var licensesCmd = &cobra.Command{
	Use:   "licenses [refresh]",
	Args:  cobra.RangeArgs(0, 1),
	Short: "List (or refresh cached) SPDX license IDs",
	Long: `Use this command to list license SPDX ids from the https://spdx.org/licenses/. 
This is a long list and as such you may want to pipe the output into "more" or 
"less" to view it. Pass the "refresh" argument to update the cached list of licenses.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			supportedLicenses := SupportedLicenses.getSupportedLicenses()
			// list licenses
			for _, license := range supportedLicenses {
				fmt.Println(license)
			}
			return nil
		} else if len(args) == 1 && args[0] == "refresh" {
			// update licenses file
			err := downloadSPDXLicenses(utils.UserHomeDir, utils.MkHttpClient(), true)
			if err != nil {
				return fmt.Errorf("unable to update SPDX licenses file: %s", err.Error())
			}
			fmt.Println("✅ Successfully updated SPDX licenses file.")
			return nil
		} else {
			return fmt.Errorf("unrecognized argument(s): '%s'", strings.Join(args, " "))
		}
	},
}

func init() {
	rootCmd.AddCommand(licensesCmd)
}
