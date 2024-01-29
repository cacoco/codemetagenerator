package codemetagenerator

import (
	"fmt"

	"github.com/spf13/cobra"
)

func checkArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		// must specify a sub-command
		return fmt.Errorf("this command must be run with a specific resource sub-command: author, contributor, or keyword")
	}
	if len(args) == 1 {
		if args[0] == "keyword" {
			// keywords MUST have at least one argument
			return fmt.Errorf("this command must be run with at least one keyword argument")
		}

		// ensure the sub-command is either authors and contributors
		if (args[0] != "author") && (args[0] != "contributor") {
			return fmt.Errorf("unrecognized resource sub-command: %s", args[0])
		}
	}
	if len(args) > 1 {
		if args[0] != "keyword" {
			// only keywords can have at least one argument
			if (args[0] != "author") && (args[0] != "contributor") {
				return fmt.Errorf("unrecognized resource sub-command: %s", args[0])
			}
			return fmt.Errorf("no args expected for the %s sub-command", args[0])
		}
	}

	return nil
}

func AddCmdRunE(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("this command must be run with a resource sub-command like author, contributor or keyword")
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:       "add [command]",
	ValidArgs: []string{"author", "contributor", "keyword"},
	Args:      checkArgs,
	Short:     "Adds resources [authors, contributors, keywords] to the in-progress codemeta.json file",
	Long: `Use this command to add authors, contributors, keywords, to the in-progress 
codemeta.json file. You can choose to clear all of the data in a field by 
running the "remove" command. When you are done adding resources, run 
"generate" to generate the final codemeta.json file. 

Note that this command must be run with a resource sub-command like author, contributor or keyword.`,
	RunE: AddCmdRunE,
}

func init() {
	rootCmd.AddCommand(addCmd)
}
