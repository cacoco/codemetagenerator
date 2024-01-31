package codemetagenerator

import (
	"errors"
	"fmt"

	"github.com/cacoco/codemetagenerator/internal"
	"github.com/spf13/cobra"
)

func EditCmdRunE(cmd *cobra.Command, args []string) error {
	key := args[0]
	newValue := args[1]

	codemeta, err := internal.LoadInProgressCodeMetaFile()
	if err != nil {
		return errors.New("unable to load the in-progress codemeta.json file for editing. Have you run \"codemetagenerator new\" yet?")
	} else {
		mutateMap := *codemeta
		updateErr := internal.UpdateMapValue(mutateMap, key, newValue)
		if updateErr != nil {
			return fmt.Errorf("unable to update the `%s` key for the in-progress codemeta.json file. Does it exist?", key)
		} else {
			err := internal.SaveInProgressCodeMetaFile(&mutateMap)
			if err != nil {
				return errors.New("unable to save in-progress codemeta.json file after editing")
			} else {
				fmt.Println("✅ Successfully edited the `" + key + "` key for the in-progress codemeta.json file.")
			}
		}
	}
	return nil
}

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit 'key' 'newValue'",
	Args:  cobra.ExactArgs(2),
	Short: "Edit an existing property value by key, e.g., 'foo' or 'foo.bar' or 'foo[1]', or 'foo[1].bar' in the in-progress codemeta.json file",
	Long: `Edits a property value for the given key in the in-progress codemeta.json file. For 
example, to edit the 'name' property, run:

codemetagenerator edit 'name' 'My New Name'

This will update the 'name' property to 'My New Name'. To edit a nested property, use dot 
notation. For example, to edit the 'email' property of the 'maitainer' object, run:

codemetagenerator edit 'maintainer.email' 'newemail@updated.org'

This will update the 'email' property of the 'maintainer' object to 'My New Email'. To edit 
a nested array property, use a bracket index with dot notation. For example, edit the 
'givenName' property of the first 'author' object in the array, run:

codemetagenerator edit 'author[0].givenName' 'My New Given Name'

This will update the 'givenName' property of the first 'author' object in the array to 
'My New Given Name'. Note that the index starts at 0, not 1. Also note that this command 
only replaces leaf/scalar values. It does not replace whole arrays or objects. For example, 
if you have two authors, there is no way to replace a full 'author' object at this time, 
nor similarily to replace the full 'maintainer' object. It is only possible to edit 
property values WITHIN an object. Keys MUST also already exist as this DOES NOT INSERT any 
new keys but only updates existing keys. To insert a new key, run the "insert" command. To 
fully remove a key and its value, run the "remove" command.

When complete, run "generate" to generate the final codemeta.json file.`,
	RunE: EditCmdRunE,
}

func init() {
	rootCmd.AddCommand(editCmd)
}
