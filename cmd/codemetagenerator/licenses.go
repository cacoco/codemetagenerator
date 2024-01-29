package codemetagenerator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cacoco/codemetagenerator/internal"
	"github.com/spf13/cobra"
)

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
		foundLicenses, err := internal.GetSupportedLicenses()
		if err != nil {
			return err
		}
		SupportedLicenses = *foundLicenses
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// list licenses
			for _, license := range SupportedLicenses {
				fmt.Println(license)
			}
			return nil
		} else if len(args) == 1 && args[0] == "refresh" {
			// update licenses file
			err := internal.GetAndCacheLicenseFile(true)
			if err != nil {
				return err
			}
			fmt.Println("✅ Successfully updated SPDX licenses file.")
			return nil
		} else {
			msg := fmt.Sprintf("unrecognized argument(s): '%s'", strings.Join(args, " "))
			return errors.New(msg)
		}
	},
}

func init() {
	rootCmd.AddCommand(licensesCmd)
}
