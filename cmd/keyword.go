package cmd

import (
	"strings"

	"github.com/cacoco/codemetagenerator/internal/model"
	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/spf13/cobra"
)

func keyword(writer utils.Writer, basedir string, args []string) ([]any, error) {
	inProgressFilePath := utils.GetInProgressFilePath(basedir)

	codemeta, err := utils.Unmarshal(inProgressFilePath)
	if err != nil {
		handleErr(writer, err)
		return nil, writer.Errorf("unable to load the in-progress codemeta.json file for editing. Have you run \"codemetagenerator new\" yet?")
	}
	mutateMap := *codemeta
	currentValue := mutateMap[model.Keywords]
	if currentValue == nil {
		var arr []any = make([]any, len(args))
		// need to copy into an	array of any type
		for i, v := range args {
			arr[i] = v
		}
		mutateMap[model.Keywords] = arr
	} else {
		for _, keyword := range args {
			// since we're ranging over the args, we need to ensure we're appending to the latest value in the map for model.Keywords key
			mutateMap[model.Keywords] = append(mutateMap[model.Keywords].([]any), keyword)
		}
	}
	writer.Println("Added keyword(s): " + strings.Join(args, ", "))
	keywords := mutateMap[model.Keywords].([]any)

	err = utils.Marshal(inProgressFilePath, mutateMap)
	if err != nil {
		return nil, writer.Errorf("unable to save in-progress codemeta.json file after editing: %s", err.Error())
	} else {
		writer.Println("⭐ Successfully updated in-progress codemeta.json file.")
	}
	return keywords, nil

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
		_, err := keyword(&utils.StdoutWriter{}, utils.UserHomeDir, args)
		return err
	},
}

func init() {
	addCmd.AddCommand(keywordCmd)
}
