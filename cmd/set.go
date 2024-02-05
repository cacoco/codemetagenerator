package cmd

import (
	"fmt"
	"os"

	internal "github.com/cacoco/codemetagenerator/internal/json"
	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/ohler55/ojg/oj"
	"github.com/spf13/cobra"
)

func set(basedir string, args []string) error {
	path := args[0]
	value := args[1]

	bytes, err := utils.LoadFile(utils.GetInProgressFilePath(basedir))
	if err != nil {
		return fmt.Errorf("unable to load the in-progress codemeta.json file: %s", err.Error())
	}

	result, err := setValue(bytes, path, value)
	if err != nil {
		return err
	}

	return os.WriteFile(utils.GetInProgressFilePath(basedir), []byte(*result), 0644)
}

func setValue(jsonBytes []byte, path string, value string) (*string, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("path is empty")
	}
	if len(value) == 0 {
		return nil, fmt.Errorf("value is empty")
	}
	json := string(jsonBytes)
	var val any
	err := oj.Unmarshal([]byte(value), &val)
	if err != nil {
		// just treat as a string
		val = value
	}
	result, err := internal.Set(json, path, val)
	if err != nil {
		return nil, fmt.Errorf("unable to set the value of the property with path, `%s` in the in-progress codemeta.json file: %s", path, err.Error())
	}
	return &result, nil
}

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Args:  cobra.ExactArgs(2),
	Short: "Set the value of an arbitrary key in the in-progress codemeta.json file",
	Long: `
Set a property by key in the in-progress codemeta.json file. This will insert or edit.

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

"-1" is a special index that appends to an array. For example:
"key3.-1"		=> appends a new value to the end of the "key3" array.

If the key does not exist, it will be created. If the key exists, its value will be overwritten.

Examples:

codemetagenerator set 'key6.key7' 'seven'
codemetagenerator set 'key14' 'fourteen'
codemetagenerator set 'key3.1' 'two'
codemetagenerator set 'key1.1.second' '5'
codemetagenerator set 'key3.-1' 'three'
codemetagenerator set 'key15' '{"first": 15, "second": 16}'
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return set(utils.UserHomeDir, args)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
