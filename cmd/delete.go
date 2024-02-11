package cmd

import (
	internal "github.com/cacoco/codemetagenerator/internal/json"
	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/spf13/cobra"
)

func delete(writer utils.Writer, basedir string, propertyPath string) error {
	bytes, err := utils.LoadFile(utils.GetInProgressFilePath(basedir))
	if err != nil {
		handleErr(writer, err)
		return writer.Errorf("unable to load the in-progress codemeta.json file")
	}

	result, err := deleteValue(writer, bytes, propertyPath)
	if err != nil {
		return err
	}
	return utils.MarshalBytes(utils.GetInProgressFilePath(basedir), []byte(*result))
}

func deleteValue(writer utils.Writer, jsonBytes []byte, path string) (*string, error) {
	if len(path) == 0 {
		return nil, writer.Errorf("path is empty")
	}

	json := string(jsonBytes)
	result, err := internal.Delete(json, path)
	if err != nil {
		handleErr(writer, err)
		return nil, writer.Errorf("unable to delete the property with path, `%s` in the in-progress codemeta.json file", path)
	}
	return &result, nil
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Args:  cobra.ExactArgs(1),
	Short: "Delete an arbitrary key and its value from the in-progress codemeta.json file",
	Long: `
Delete a property by key in the in-progress codemeta.json file.

This expects a path syntax for the key. A path is a series of keys separated by a dot.
For example, using this json:
	{
		"key1": [
				{"first": 2, "second": 3},
				{"first": 4, "second": 5}
			],
		"key3": ["one", "two"],
		"key5": "hello",
		"key6": {
			"key7": "seven",
			"key8": 8,
			"key9": {
				"key10": [1, 2, 3],
				"key11": {
					"key12": "twelve"
				},
				"key13": "world"
		}
			}
		"key14": "fourteen"
	}
	
	"key6.key7"		=> "seven"
	"key14"			=> "fourteen"
	"key3.1"		=> "two"
	"key1.1.second"	=> 5

"-1" is a special index which represents the last element of an array. For example:	
"key3.-1"  => deletes the last element of the "key3" array, ("two").
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return delete(&utils.StdoutWriter{}, utils.UserHomeDir, args[0])
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
