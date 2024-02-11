package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/cacoco/codemetagenerator/internal/cue"
	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/spf13/cobra"
)

func validate(basedir string, writer utils.Writer, inFile string) error {
	var path string
	if inFile == "" {
		path = utils.GetInProgressFilePath(basedir)
	} else {
		path = inFile
	}
	json, err := utils.ReadJSON(path)
	if err != nil {
		handleErr(writer, err)
		return writer.Errorf("unable to read codemeta.inprogress.json file, ensure you have run `codemetagenerator new` at least once or specify a file with the --input flag")
	}

	// validate the file
	err = cue.Validate([]byte(*json))
	if err != nil {
		handleErr(writer, err)
		return writer.Errorf("invalid codemeta.json file: %v", err)
	}

	writer.Println(fmt.Sprintf("✅ The codemeta file '%s' is valid.", filepath.Base(path)))
	return nil
}

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validates a codemeta.json file",
	Long:  `Validates a codemeta.json file. If no input file is specified, the current in progress file will be used.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return validate(utils.UserHomeDir, &utils.StdoutWriter{}, inputFile)
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringVarP(&inputFile, "input", "i", "", "path to an input 'codemeta.json' file. If not specified, the current in progress file will be used.")
}
