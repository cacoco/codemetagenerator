package cmd

import (
	"fmt"
	"os"

	internal "github.com/cacoco/codemetagenerator/internal/json"
	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/ohler55/ojg/oj"
	"github.com/spf13/cobra"
)

func Delete(jsonBytes []byte, path string) (*string, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("path is empty")
	}

	json := string(jsonBytes)
	result, err := internal.Delete(json, path)
	if err != nil {
		return nil, fmt.Errorf("unable to delete the property with path, `%s` in the in-progress codemeta.json file: %s", path, err.Error())
	}
	return &result, nil
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Args:  cobra.ExactArgs(1),
	Short: "Delete an arbitrary key and its value in the in-progress codemeta.json file.",
	Long: `This expects a path syntax for the key. A path is a series of keys separated by a dot.
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
		inProgressFilePath := utils.GetInProgressFilePath(utils.UserHomeDir)

		path := args[0]
		data, err := utils.Unmarshal(inProgressFilePath)
		if err != nil {
			return err
		}

		json := oj.JSON(data)
		result, err := internal.Delete(json, path)
		if err != nil {
			return err
		}
		return os.WriteFile(inProgressFilePath, []byte(result), 0644)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
