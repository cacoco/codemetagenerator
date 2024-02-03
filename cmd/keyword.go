package cmd

import (
	"fmt"
	"strings"

	"github.com/cacoco/codemetagenerator/internal/model"
	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/spf13/cobra"
)

func newKeyword(writer utils.Writer, codemeta map[string]any, args []string) ([]string, error) {
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
	return keywords, nil
}

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
		inProgressFilePath := utils.GetInProgressFilePath(utils.UserHomeDir)

		codemeta, err := utils.Unmarshal(inProgressFilePath)
		if err != nil {
			return fmt.Errorf("unable to load the in-progress codemeta.json file for editing. Have you run \"codemetagenerator new\" yet?")
		}
		mutateMap := *codemeta
		_, err = newKeyword(&utils.StdoutWriter{}, mutateMap, args)
		if err != nil {
			return fmt.Errorf("unable to create new keyword: %s", err.Error())
		}
		err = utils.Marshal(inProgressFilePath, &mutateMap)
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
