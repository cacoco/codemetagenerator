package codemetagenerator

import (
	"errors"
	"fmt"

	"github.com/cacoco/codemetagenerator/internal"
	"github.com/spf13/cobra"
)

// authorCmd represents the author command
var authorCmd = &cobra.Command{
	Use:   "author",
	Args:  cobra.NoArgs,
	Short: "Adds an author to the in-progress codemeta.json file",
	Long: `Add a single author to the in-progres codemeta.json file. An author can be a 
person or an organization. Prompts for the information needed to add an 
author and then add it to the in-progress codemeta.json file. You can add 
multiple authors by running this command multiple times. If you need to remove 
an author, run the "remove" command to remove authors. Run the "edit" command 
to edit properties of an author. 

When complete, run "generate" to generate the final codemeta.json file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		codemeta, err := internal.LoadInProgressCodeMetaFile()
		if err != nil {
			msg := "unable to load the in-progress codemeta.json file for editing. Have you run \"codemetagenerator new\" yet?"
			return errors.New(msg)
		}
		author, newAuthorErr := internal.NewPersonOrOrganizationPrompt("Author")
		if newAuthorErr != nil {
			return fmt.Errorf("unable to create new author: %s", newAuthorErr.Error())
		}
		mutateMap := *codemeta
		currentValue := mutateMap[internal.Author]
		if currentValue == nil {
			mutateMap[internal.Author] = []any{author}
		} else {
			mutateMap[internal.Author] = append(currentValue.([]any), author)
		}

		saveErr := internal.SaveInProgressCodeMetaFile(&mutateMap)
		if saveErr != nil {
			return fmt.Errorf("unable to save in-progress codemeta.json file after editing: %s", saveErr.Error())
		} else {
			fmt.Println("⭐ Successfully updated in-progress codemeta.json file.")
		}
		return nil
	},
}

func init() {
	addCmd.AddCommand(authorCmd)
}
