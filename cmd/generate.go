package cmd

import (
	"fmt"
	"os"

	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/ohler55/ojg/oj"
	"github.com/spf13/cobra"
)

func generate(basedir string, writer utils.Writer, outputPath string) error {
	inProgressFilePath := utils.GetInProgressFilePath(basedir)

	codemeta, err := utils.Unmarshal(inProgressFilePath)
	// only generate if there are no errors
	if err != nil {
		return writer.Errorf("unable to load the in-progress codemeta.json file for generation. Have you run \"codemetagenerator new\" yet?")
	}
	bytes, err := oj.Marshal(codemeta, 80.2)
	if err != nil {
		return writer.Errorf("unable to marshal codemeta.json file in bytes for generation: %s", err.Error())
	}
	if outputPath != "" {
		// write to file
		err = os.WriteFile(outputPath, bytes, 0644)
		if err != nil {
			return fmt.Errorf("unable to write codemeta.json file to output file %s: %s", outputPath, err.Error())
		}
	} else {
		writer.Println(string(bytes))
	}
	return nil
}

var codeMetaFilePath string

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate [-o | --output <path/to/codemeta.json>]",
	Args:  cobra.NoArgs,
	Short: "Generate the final codemeta.json file to the optional output file or to the console",
	Long: `Generates the final codemeta.json file from the in-progress codemeta.json file. 
Output can be written to a file [-o | --output  <path/to/codemeta.json>] or 
printed to the console.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return generate(utils.UserHomeDir, &utils.StdoutWriter{}, codeMetaFilePath)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&codeMetaFilePath, "output", "o", "", "The path to the output codemeta.json file. If not specified, the output will be printed to the console.")
}
