package codemetagenerator

import (
	"errors"
	"fmt"

	"github.com/cacoco/codemetagenerator/internal"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove 'key'",
	Args:  cobra.ExactArgs(1),
	Short: "Remove a property by key, e.g., 'foo' or 'foo.bar' or 'foo[1]', or 'foo[1].bar' from the in-progress codemeta.json file",
	Long: `Removes a given property and its value from the in-progress codemeta.json file. 
For example, to remove the 'name' property, run:

codemetagenerator remove 'name'
	
This will remove the 'name' property from the JSON. To remove a nested 
property, use dot notation. For example, to remove the 'email' property of the 
'maitainer' object, run:
	
codemetagenerator remove 'maintainer.email'
	
This will remove the 'email' property of the 'maintainer' object. To remove a 
nested array property, use a bracket index with dot notation. For example, 
remove the 'givenName' property of the first 'author' object in the array, run:
	
codemetagenerator remove 'author[0].givenName'
	
This will remove the 'givenName' property of the first 'author' object in the 
array. Note that the index starts at 0, not 1.

When complete, run "generate" to generate the final codemeta.json file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		codemeta, err := internal.LoadInProgressCodeMetaFile()
		if err != nil {
			msg := "unable to load the in-progress codemeta.json file for editing. Have you run \"codemetagenerator new\" yet?"
			return errors.New(msg)
		}

		key := args[0]
		mutateMap := *codemeta
		updateErr := internal.RemoveMapValue(mutateMap, key)
		if updateErr != nil {
			return fmt.Errorf("unable to remove the `%s` key from the in-progress codemeta.json file. Does it exist?", key)
		}
		saveErr := internal.SaveInProgressCodeMetaFile(&mutateMap)
		if saveErr != nil {
			return fmt.Errorf("unable to save in-progress codemeta.json file after editing: %s", saveErr.Error())
		}

		fmt.Println("✅ Successfully removed the `" + key + "` key from the in-progress codemeta.json file.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
