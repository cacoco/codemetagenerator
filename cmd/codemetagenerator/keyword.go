package codemetagenerator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cacoco/codemetagenerator/internal"
	"github.com/spf13/cobra"
)

// keywordCmd represents the keyword command
var keywordCmd = &cobra.Command{
	Use:   "keyword",
	Args:  cobra.MinimumNArgs(1),
	Short: "Adds a keyword to the in-progress codemeta.json file",
	Long: `Add a single keyword to the in-progres codemeta.json file. A keyword can be a 
person or an organization. Prompts for the information needed to add a keyword 
and then adds it to the 'keywords' array in the in-progress codemeta.json file. 
You can add multiple keywords by running this command multiple times. If you 
need to remove a keyword, run the "remove" command to remove keywords. Run the 
"edit" command to edit properties of a keyword.

When complete, run "generate" to generate the final codemeta.json file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		codemeta, err := internal.LoadInProgressCodeMetaFile()
		if err != nil {
			msg := "unable to load the in-progress codemeta.json file for editing. Have you run \"codemetagenerator new\" yet?"
			return errors.New(msg)
		}
		mutateMap := *codemeta
		currentValue := mutateMap[internal.Keywords]
		if currentValue == nil {
			mutateMap[internal.Keywords] = args
		} else {
			for _, arg := range args {
				mutateMap[internal.Keywords] = append(mutateMap[internal.Keywords].([]string), arg)
			}
		}
		fmt.Println("Added keyword(s): " + strings.Join(args, ", "))

		saveErr := internal.SaveInProgressCodeMetaFile(&mutateMap)
		if saveErr != nil {
			return errors.New("unable to save in-progress codemeta.json file after editing")
		} else {
			fmt.Println("⭐ Successfully updated in-progress codemeta.json file.")
		}
		return nil
	},
}

func init() {
	addCmd.AddCommand(keywordCmd)
}
