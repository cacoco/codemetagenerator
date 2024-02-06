package cmd

import (
	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/spf13/cobra"
)

func generate(basedir string, writer utils.Writer, outFile string) error {
	inProgressFilePath := utils.GetInProgressFilePath(basedir)

	json, err := utils.ReadJSON(inProgressFilePath)
	if err != nil {
		handleErr(writer, err)
		return writer.Errorf("unable to read codemeta.inprogress.json file, ensure you have run `codemetagenerator new` at least once")
	}

	if outFile != "" {
		err = utils.WriteJSON(outFile, *json)
		if err != nil {
			handleErr(writer, err)
			return writer.Errorf("unable to write codemeta.json file to output file %s", outFile)
		}
	} else {
		writer.Println(*json)

	}
	return nil
}

var outputFile string

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate [-o | --output <path/to/codemeta.json>]",
	Args:  cobra.NoArgs,
	Short: "Generate the resultant 'codemeta.json' file to the optional output file or to the console",
	Long: `
Generate the resultant 'codemeta.json' file from the in-progress file. 

Output can be written to a file [-o | --output  <path/to/codemeta.json>] or 
printed to the console.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return generate(utils.UserHomeDir, &utils.StdoutWriter{}, outputFile)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&outputFile, "output", "o", "", "The path to the output 'codemeta.json' file. If not specified, the output will be printed to the console.")
}
