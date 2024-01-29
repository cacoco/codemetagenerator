package codemetagenerator

import (
	"errors"
	"fmt"

	"github.com/cacoco/codemetagenerator/internal"
	"github.com/spf13/cobra"
)

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
		codemeta, err := internal.LoadInProgressCodeMetaFile()
		if err != nil {
			msg := "unable to load the in-progress codemeta.json file for editing. Have you run \"codemetagenerator new\" yet?"
			return errors.New(msg)
		}
		contributor, newContributorErr := internal.NewPersonOrOrganizationPrompt("Contributor")
		if newContributorErr != nil {
			return errors.New("unable to create new contributor")
		}
		mutateMap := *codemeta
		currrentValue := mutateMap[internal.Contributor]
		if currrentValue == nil {
			mutateMap[internal.Contributor] = []any{contributor}
		} else {
			mutateMap[internal.Contributor] = append(currrentValue.([]any), contributor)
		}

		saveErr := internal.SaveInProgressCodeMetaFile(codemeta)
		if saveErr != nil {
			return errors.New("unable to save in-progress codemeta.json file after editing")
		} else {
			fmt.Println("⭐ Successfully updated in-progress codemeta.json file.")
		}
		return nil
	},
}

func init() {
	addCmd.AddCommand(contributorCmd)
}
