package cmd

import (
	"fmt"

	"github.com/cacoco/codemetagenerator/internal/model"
	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/spf13/cobra"
)

func newContributor(reader utils.Reader, writer utils.Writer) (*map[string]any, error) {
	contributor, err := utils.NewPersonOrOrganizationPrompt(&reader, &writer, "Contributor")
	if err != nil {
		return nil, err
	}
	return contributor, nil
}

// contributorCmd represents the contributor command
var contributorCmd = &cobra.Command{
	Use:   "contributor",
	Args:  cobra.NoArgs,
	Short: "Adds a contributor to the in-progress codemeta.json file",
	Long: `Add a single contributor to the in-progres codemeta.json file. An contributor 
can be a person or an organization. Prompts for the information needed to add 
a contributor and then add it to the in-progress codemeta.json file. You can 
add multiple contributors by running this command multiple times. If you need 
to remove a contributor, run the "remove" command to remove contributors. Run 
the "edit" command to edit properties of a contributor. 

When complete, run "generate" to generate the final codemeta.json file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		inProgressFilePath := utils.GetInProgressFilePath(utils.UserHomeDir)

		codemeta, err := utils.Unmarshal(inProgressFilePath)
		if err != nil {
			return fmt.Errorf("unable to load the in-progress codemeta.json file for editing. Have you run \"codemetagenerator new\" yet?")
		}

		mutateMap := *codemeta
		currrentValue := mutateMap[model.Contributor]
		contributor, err := newContributor(&utils.StdinReader{}, &utils.StdoutWriter{})
		if err != nil {
			return fmt.Errorf("unable to create new contributor: %s", err.Error())
		}
		if currrentValue == nil {
			mutateMap[model.Contributor] = []any{contributor}
		} else {
			mutateMap[model.Contributor] = append(currrentValue.([]any), contributor)
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
	addCmd.AddCommand(contributorCmd)
}
