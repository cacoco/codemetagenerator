package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/cacoco/codemetagenerator/internal/model"
	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func new(basedir string, reader utils.Reader, writer utils.Writer, inFile string) error {
	stdin := reader.Stdin()
	stdout := writer.Stdout()

	// clean up any previous file
	inProgressFilePath := utils.GetInProgressFilePath(basedir)
	utils.DeleteFile(inProgressFilePath)

	var successMsg string = "⭐ Successfully created new in-progress codemeta.json file."
	if inFile != "" {
		bytes, err := utils.LoadFile(inFile)
		if err != nil {
			handleErr(writer, err)
			return writer.Errorf("unable to read input file %s", inFile)
		}
		err = utils.WriteFile(inProgressFilePath, bytes)
		if err != nil {
			handleErr(writer, err)
			return writer.Errorf("unable to write in-progress codemeta.json file")
		}

		fileBase := filepath.Base(inFile)
		successMsg = fmt.Sprintf("⭐ Successfully loaded '%s' as new in-progress codemeta.json file.", fileBase)
	} else {
		var result = make(map[string]any)

		identifier, err := utils.MkPrompt(&stdin, &stdout, "Enter a unique identifier for your software source code", utils.Nop)
		if err != nil {
			handleErr(writer, err)
			return writer.Errorf("unable to create new identifier")
		}

		name, err := utils.MkPrompt(&stdin, &stdout, "Enter a name for your software source code", utils.Nop)
		if err != nil {
			handleErr(writer, err)
			return writer.Errorf("unable to create new name")
		}

		description, err := utils.MkPrompt(&stdin, &stdout, "Enter a description for your software source code", utils.Nop)
		if err != nil {
			handleErr(writer, err)
			return writer.Errorf("unable to create new description")
		}

		developmentStatusOptions := []model.MenuOption{
			{Name: "Abandoned", Type: "abandoned"},
			{Name: "Active", Type: "active"},
			{Name: "Concept", Type: "concept"},
			{Name: "Inactive", Type: "inactive"},
			{Name: "Moved", Type: "moved"},
			{Name: "Suspended", Type: "suspended"},
			{Name: "Unsupported", Type: "unsupported"},
			{Name: "WIP", Type: "wip"},
		}

		templates := &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "➞ {{ .Name | cyan }}",
			Inactive: "  {{ .Name | cyan }}",
			Selected: `{{ "Select a development status (see: https://www.repostatus.org/):" | faint}} {{ .Name | faint }}`,
			Details: `--------- Status ----------
{{ "Name:" | faint }}	{{ .Name }}`,
		}

		prompt := promptui.Select{
			Label:     "Select a development status from the list below (see: https://www.repostatus.org/)",
			Items:     developmentStatusOptions,
			Templates: templates,
			Size:      8,
			Searcher:  nil,
			Stdin:     reader.Stdin(),
			Stdout:    writer.Stdout(),
		}

		i, _, err := prompt.Run()
		if err != nil {
			return err
		}
		developmentStatus := developmentStatusOptions[i].Name

		codeRepository, err := utils.MkPrompt(&stdin, &stdout, "Enter the URL of the code repository for the project", utils.ValidUrl)
		if err != nil {
			return err
		}

		programmingLanguageName, err := utils.MkPrompt(&stdin, &stdout, "Enter the name of the programming language of the project", utils.Nop)
		if err != nil {
			return err
		}
		programmingLanguageURL, err := utils.MkPrompt(&stdin, &stdout, "Enter the URL of the programming language of the project", utils.ValidUrl)
		if err != nil {
			return err
		}
		programmingLanguage := model.NewProgrammingLanguage(programmingLanguageName, programmingLanguageURL)

		runtimePlatform, err := utils.MkPrompt(&stdin, &stdout, "Enter the name of the runtime platform of the project", utils.Nop)
		if err != nil {
			return err
		}

		version, err := utils.MkPrompt(&stdin, &stdout, "Enter the version of the project", utils.Nop)
		if err != nil {
			return err
		}

		validateFn := validateLicenseId(basedir)
		license, err := utils.MkPrompt(&stdin, &stdout, "Enter the SPDX license ID for the project (see: https://spdx.org/licenses/)", validateFn)
		if err != nil {
			return err
		}

		var licenseDetailsUrl string
		if (*license) != "" {
			referenceUrl, err := getLicenseReference(basedir, *license)
			if err != nil {
				handleErr(writer, err)
				return writer.Errorf("unable to create new license details URL")
			}
			licenseDetailsUrl = *referenceUrl
		}

		readme, err := utils.MkPrompt(&stdin, &stdout, "Enter the URL of the README file for the project", utils.ValidUrl)
		if err != nil {
			return err
		}

		maintainer, err := utils.NewPersonOrOrganizationPrompt(&reader, &writer, "Maintainer")
		if err != nil {
			return err
		}

		result[model.Identifier] = identifier
		result[model.Name] = name
		result[model.Description] = description
		result[model.Version] = version
		result[model.Maintainer] = maintainer
		result[model.ProgrammingLanguage] = programmingLanguage
		result[model.DevelopmentStatus] = developmentStatus
		result[model.License] = licenseDetailsUrl
		result[model.RuntimePlatform] = runtimePlatform
		result[model.CodeRepository] = codeRepository
		result[model.Readme] = readme

		codemeta := *model.NewCodemeta(&result)

		err = utils.Marshal(inProgressFilePath, codemeta)
		if err != nil {
			handleErr(writer, err)
			return writer.Errorf("unable to save in-progress codemeta.json file after editing")
		}
	}

	writer.Println(successMsg)
	writer.Println("👇 You can now add authors, contributors, keywords, and other fields to the in-progress codemeta.json file.")
	writer.Println("➡️  To add/remove authors, contributors or keywords, run the following commands:")
	writer.Println("\tcodemetagenerator add author")
	writer.Println("\tcodemetagenerator remove authors")
	writer.Println("\tcodemetagenerator add contributor")
	writer.Println("\tcodemetagenerator remove contributors")
	writer.Println("\tcodemetagenerator add keyword")
	writer.Println("\tcodemetagenerator remove keywords")
	writer.Println("↔️  To edit any key in the in-progress codemeta.json file, run the following command:")
	writer.Println("\tcodemetagenerator edit key.subkey newValue")
	writer.Println("✅ To generate the resultant 'codemeta.json' file, run the following command:")
	writer.Println("\tcodemetagenerator generate [-o|--output] <output file path>")

	return nil
}

var inputFile string

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Args:  cobra.NoArgs,
	Short: "Start a new codemeta.json file for editing. When complete, run \"codemetagenerator generate\" to generate the resultant 'codemeta.json' file",
	Long: `
This command starts a new 'codemeta.json' file with basic information about your 
software source code project. 

It is expected that you will also add authors, contributors, keywords, and other 
fields to the in-progress codemeta.json, then run:

codemetagenerator generate

to generate the resultant 'codemeta.json' file, optionally selecting the file destination.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return new(utils.UserHomeDir, &utils.StdinReader{}, &utils.StdoutWriter{}, inputFile)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	utils.MkHomeDir(utils.UserHomeDir)

	newCmd.Flags().StringVarP(&inputFile, "input", "i", "", "path to an input 'codemeta.json' file. If not specified, a new file will be started.")
}
