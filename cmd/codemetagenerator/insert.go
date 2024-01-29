package codemetagenerator

import (
	"errors"
	"fmt"

	"github.com/cacoco/codemetagenerator/internal"
	"github.com/spf13/cobra"
)

// insertCmd represents the insert command
var insertCmd = &cobra.Command{
	Use:   "insert 'key' 'value'",
	Args:  cobra.ExactArgs(2),
	Short: "Insert a new property with the given value by key, e.g., 'foo' or 'foo.bar' or 'foo[1].bar' into the in-progress codemeta.json file",
	Long: `Inserts a new property with the given key and value into the in-progress 
codemeta.json file. For example, to insert a 'name' key with a value of 
'My Project', run:
	
codemetagenerator insert 'name' 'My Project'

To insert a nested property, use dot notation. For example, to insert an 
'email' property of the 'maitainer' object, run:

codemetagenerator insert 'maintainer.email' 'foo@bar.com'

This will insert the 'email' property of the 'maintainer' object. To insert 
a nested array property, use a bracket index with dot notation. For example, 
insert the 'givenName' property into the first 'author' object in the array, 
run:

codemetagenerator insert 'author[0].givenName' 'Foo'

This will insert the 'givenName' property into the first 'author' object in 
the array. Note that the index starts at 0, not 1. Also note that you CANNOT 
insert over an existing key nor can you insert into an array, only into an
object value within an array.

To edit an existing key, run the "edit" command. 
To fully remove a key and its value, run the "remove" command.

When complete, run "generate" to generate the final codemeta.json file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		codemeta, err := internal.LoadInProgressCodeMetaFile()
		if err != nil {
			msg := "unable to load the in-progress codemeta.json file for editing. Have you run \"codemetagenerator new\" yet?"
			return errors.New(msg)
		}
		mutateMap := *codemeta
		key := args[0]
		value := args[1]
		updateErr := internal.InsertMapValue(mutateMap, key, value)
		if updateErr != nil {
			msg := fmt.Sprintf("unable to insert the `%s` key to the in-progress codemeta.json file. Try editing an existing key instead.", key)
			return errors.New(msg)
		} else {
			err := internal.SaveInProgressCodeMetaFile(&mutateMap)
			if err != nil {
				return errors.New("unable to save in-progress codemeta.json file after editing")
			} else {
				fmt.Println("✅ Successfully inserted the `" + key + "` key to the in-progress codemeta.json file.")
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(insertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// insertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// insertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
