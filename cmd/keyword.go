package cmd

import (
	"fmt"
	"strings"

	"github.com/cacoco/codemetagenerator/internal/model"
	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/spf13/cobra"
)

func addKeywords(writer utils.Writer, codemeta map[string]any, args []string) []string {
	currentValue := codemeta[model.Keywords]
	if currentValue == nil {
		codemeta[model.Keywords] = args
	} else {
		for _, arg := range args {
			codemeta[model.Keywords] = append(codemeta[model.Keywords].([]string), arg)
		}
	}
	writer.Println("Added keyword(s): " + strings.Join(args, ", "))
	keywords := codemeta[model.Keywords].([]string)
	return keywords
}

// keywordCmd represents the keyword command
var keywordCmd = &cobra.Command{
	Use:   "keyword",
	Args:  cobra.MinimumNArgs(1),
	Short: "Adds a keyword to the in-progress codemeta.json file",
	Long: `
Add a single keyword to the in-progress codemeta.json file. 

A keyword can be any string value. Prompts for the information needed to add 
a keyword and then appends it to the 'keywords' array in the in-progress codemeta.json 
file. 

You can add multiple keywords by either running this command multiple times or passing
multiple values to the command, e.g.

codemetagenerator add keyword keyword1 keyword2 keyword3

will add three keywords to the in-progress codemeta.json file. If you need to remove a 
keyword, run the "delete" command to remove keywords. Run the "set" command to edit 
properties of a keyword.

When complete, run "generate" to generate the resultant 'codemeta.json' file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		stdout := &utils.StdoutWriter{}

		inProgressFilePath := utils.GetInProgressFilePath(utils.UserHomeDir)

		codemeta, err := utils.Unmarshal(inProgressFilePath)
		if err != nil {
			handleErr(stdout, err)
			return fmt.Errorf("unable to load the in-progress codemeta.json file for editing. Have you run \"codemetagenerator new\" yet?")
		}
		mutateMap := *codemeta
		addKeywords(stdout, mutateMap, args)

		err = utils.Marshal(inProgressFilePath, mutateMap)
		if err != nil {
			return fmt.Errorf("unable to save in-progress codemeta.json file after editing: %s", err.Error())
		} else {
			fmt.Println("⭐ Successfully updated in-progress codemeta.json file.")
		}
		return nil
	},
}

func init() {
	addCmd.AddCommand(keywordCmd)
}
