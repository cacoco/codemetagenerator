package codemetagenerator

import (
	"errors"
	"fmt"
	"os"

	"github.com/nexidian/gocliselect"
	"github.com/spf13/cobra"

	"github.com/cacoco/codemetagenerator/internal"
)

type RepositoryStatus struct {
	Description string `json:"description"`
	DisplayName string `json:"display_name"`
	URL         string `json:"url"`
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Args:  cobra.NoArgs,
	Short: "Start a new codemeta.json file. When complete, run \"codemetagenerator generate\" to generate the final codemeta.json file",
	Long: `This command starts a new codemeta.json file with basic information about your 
software source code project. It is expected that you will also add authors, 
contributors, keywords, and other fields to the in-progress codemeta.json, 
then run:

codemetagenerator generate

to generate the final codemeta.json file optionally selecting the file 
destination.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// clean up any previous file
		internal.DeleteInProgressCodeMetaFile()

		identifier, idErr := internal.MkPrompt("Enter a unique identifier for your software source code:")
		if idErr != nil {
			return errors.New("unable to create new identifier")
		}

		name, nameErr := internal.MkPrompt("Enter a name for your software source code:")
		if nameErr != nil {
			return errors.New("unable to create new name")
		}

		description, descErr := internal.MkPrompt("Enter a description for your software source code:")
		if descErr != nil {
			return errors.New("unable to create new description")
		}

		developmentStatusMenu := gocliselect.NewMenu("Select a development status from the list below (see: https://www.repostatus.org/)")
		developmentStatusMenu.AddItem("Abandoned", "abandoned")
		developmentStatusMenu.AddItem("Active", "active")
		developmentStatusMenu.AddItem("Concept", "concept")
		developmentStatusMenu.AddItem("Inactive", "inactive")
		developmentStatusMenu.AddItem("Moved", "moved")
		developmentStatusMenu.AddItem("Suspended", "suspended")
		developmentStatusMenu.AddItem("Unsupported", "unsupported")
		developmentStatusMenu.AddItem("WIP", "wip")
		developmentStatus := developmentStatusMenu.Display()

		codeRepository, crErr := internal.MkPrompt("Enter the URL of the code repository for the project:")
		if crErr != nil {
			return fmt.Errorf("unable to create new code repository: %s", crErr)
		}

		programmingLanguageName, plnErr := internal.MkPrompt("Enter the name of the programming language of the project:")
		if plnErr != nil {
			return fmt.Errorf("unable to create new programming language name: %s", plnErr.Error())
		}
		programmingLanguageURL, pluErr := internal.MkPrompt("Enter the URL of the programming language of the project:")
		if pluErr != nil {
			return fmt.Errorf("unable to create new programming language URL: %s", pluErr.Error())
		}
		programmingLanguage := internal.NewProgrammingLanguage(programmingLanguageName, programmingLanguageURL)

		runtimePlatform, rpErr := internal.MkPrompt("Enter the name of the runtime platform of the project:")
		if rpErr != nil {
			return fmt.Errorf("unable to create new runtime platform: %s", rpErr.Error())
		}

		version, verErr := internal.MkPrompt("Enter the version of the project:")
		if verErr != nil {
			return fmt.Errorf("unable to create new version: %s", verErr.Error())
		}

		license, licErr := internal.MkPrompt("Enter the SPDX license ID for the project (see: https://spdx.org/licenses/):")
		if licErr != nil {
			return fmt.Errorf("unable to create new license: %s", licErr.Error())
		}
		licenseDetailsUrl, lcdErr := internal.CheckAndConvertLicenseId(license)
		if lcdErr != nil {
			return lcdErr
		}

		readme, rmErr := internal.MkPrompt("Enter the URL of the README file for the project:")
		if rmErr != nil {
			return fmt.Errorf("unable to create new README: %s", rmErr.Error())
		}

		maintainer, mErr := internal.NewPersonOrOrganizationPrompt("Maintainer")
		if mErr != nil {
			return fmt.Errorf("unable to create new maintainer: %s", mErr.Error())
		}

		codemeta := internal.NewCodeMetaDefinition(identifier, name, description, version, maintainer, programmingLanguage, developmentStatus, licenseDetailsUrl, runtimePlatform, codeRepository, readme)
		saveErr := internal.SaveInProgressCodeMetaFile(codemeta)
		if saveErr != nil {
			return fmt.Errorf("unable to save in-progress codemeta.json file after editing: %s", saveErr.Error())
		}
		fmt.Println("⭐ Successfully created new in-progress codemeta.json file.")
		fmt.Println("➡️  To add/remove authors, contributors or keywords, run the following commands:")
		fmt.Println("\tcodemetagenerator add author")
		fmt.Println("\tcodemetagenerator remove authors")
		fmt.Println("\tcodemetagenerator add contributor")
		fmt.Println("\tcodemetagenerator remove contributors")
		fmt.Println("\tcodemetagenerator add keyword")
		fmt.Println("\tcodemetagenerator remove keywords")
		fmt.Println("↔️  To edit any key in the in-progress codemeta.json file, run the following command:")
		fmt.Println("\tcodemetagenerator edit key.subkey newValue")
		fmt.Println("✅ To generate the final codemeta.json file, run the following command:")
		fmt.Println("\tcodemetagenerator generate [-o|--output] <output file path>")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error: unable to get the user's home directory.")
		return
	}
	codemetageneratorHomeDir := homeDir + "/" + internal.CodemetaGeneratorDirectoryName
	if _, err := os.Stat(codemetageneratorHomeDir); err != nil {
		// doesn't exist
		err := os.Mkdir(codemetageneratorHomeDir, os.ModePerm)
		if err != nil {
			fmt.Println("Error: unable to create " + internal.CodemetaGeneratorDirectoryName + " in $HOME directory.")
			return
		}
	}
}
