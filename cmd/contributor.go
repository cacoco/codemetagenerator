package cmd

import (
	"github.com/cacoco/codemetagenerator/internal/model"
	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/spf13/cobra"
)

func contributor(reader utils.Reader, writer utils.Writer, basedir string) (*map[string]any, error) {
	inProgressFilePath := utils.GetInProgressFilePath(basedir)

	codemeta, err := utils.Unmarshal(inProgressFilePath)
	if err != nil {
		handleErr(writer, err)
		return nil, writer.Errorf("unable to load the in-progress codemeta.json file for editing. Have you run \"codemetagenerator new\" yet?")
	}
	mutateMap := *codemeta
	currrentValue := mutateMap[model.Contributor]
	contributor, err := utils.NewPersonOrOrganizationPrompt(&reader, &writer, "Contributor")
	if err != nil {
		return nil, err
	}
	if currrentValue == nil {
		// set the new contributor as the value
		mutateMap[model.Contributor] = []any{contributor}
	} else {
		// append new contributor to existing contributors
		mutateMap[model.Contributor] = append(currrentValue.([]any), contributor)
	}

	err = utils.Marshal(inProgressFilePath, mutateMap)
	if err != nil {
		handleErr(writer, err)
		return nil, writer.Errorf("unable to save in-progress codemeta.json file after editing")
	} else {
		writer.Println("⭐ Successfully updated in-progress codemeta.json file.")
	}
	return contributor, nil
}

// contributorCmd represents the contributor command
var contributorCmd = &cobra.Command{
	Use:   "contributor",
	Args:  cobra.NoArgs,
	Short: "Adds a contributor to the in-progress codemeta.json file",
	Long: `
Add a single contributor to the in-progres codemeta.json file. 

A contributor can be a person or an organization. Prompts for the information 
needed to add a contributor and then add it to the in-progress codemeta.json 
file. You can add multiple contributors by running this command multiple times.
If you need to remove a contributor, run the "delete" command to remove contributors. 
Run the "set" command to edit properties of a contributor. 

When complete, run "generate" to generate the resultant 'codemeta.json' file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := contributor(&utils.StdinReader{}, &utils.StdoutWriter{}, utils.UserHomeDir)
		return err
	},
}

func init() {
	addCmd.AddCommand(contributorCmd)
}
