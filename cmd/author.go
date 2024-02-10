package cmd

import (
	"fmt"

	"github.com/cacoco/codemetagenerator/internal/model"
	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/spf13/cobra"
)

func newAuthor(reader utils.Reader, writer utils.Writer) (*map[string]any, error) {
	author, err := utils.NewPersonOrOrganizationPrompt(&reader, &writer, "Author")
	if err != nil {
		return nil, err
	}
	return author, nil
}

// authorCmd represents the author command
var authorCmd = &cobra.Command{
	Use:   "author",
	Args:  cobra.NoArgs,
	Short: "Adds an author to the in-progress codemeta.json file",
	Long: `
Add a single author to the in-progres codemeta.json file. 

An author can be a person or an organization. Prompts for the information 
needed to add an author and then add it to the in-progress codemeta.json file. 
You can add multiple authors by running this command multiple times. If you 
need to remove an author, run the "delete" command to remove authors. Run the 
"set" command to edit properties of an author. 

When complete, run "generate" to generate the resultant 'codemeta.json' file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := &utils.StdinReader{}
		writer := &utils.StdoutWriter{}

		inProgressFilePath := utils.GetInProgressFilePath(utils.UserHomeDir)

		codemeta, err := utils.Unmarshal(inProgressFilePath)
		if err != nil {
			handleErr(writer, err)
			return fmt.Errorf("unable to load the in-progress codemeta.json file for editing. Have you run \"codemetagenerator new\" yet?")
		}
		mutateMap := *codemeta
		currentValue := mutateMap[model.Author]
		author, err := newAuthor(reader, writer)
		if err != nil {
			return err
		}
		if currentValue == nil {
			mutateMap[model.Author] = []any{author}
		} else {
			mutateMap[model.Author] = append(currentValue.([]any), author)
		}

		err = utils.Marshal(inProgressFilePath, mutateMap)
		if err != nil {
			handleErr(writer, err)
			return fmt.Errorf("unable to save in-progress codemeta.json file after editing")
		} else {
			fmt.Println("⭐ Successfully updated in-progress codemeta.json file.")
		}
		return nil
	},
}

func init() {
	addCmd.AddCommand(authorCmd)
}
