package codemetagenerator

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/cacoco/codemetagenerator/internal"
	"github.com/spf13/cobra"
)

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
		codemeta, err := internal.LoadInProgressCodeMetaFile()
		// only generate if there are no errors
		if err != nil {
			msg := "unable to load the in-progress codemeta.json file for generation. Have you run \"codemetagenerator new\" yet?"
			return errors.New(msg)
		}
		bytes, err := json.MarshalIndent(codemeta, "", " ")
		if err != nil {
			return errors.New("unable to marshal codemeta.json file in bytes for generation")
		}
		if codeMetaFilePath != "" {
			// write to file
			err = os.WriteFile(codeMetaFilePath, bytes, 0644)
			if err != nil {
				return fmt.Errorf("unable to write codemeta.json file to output file %s", codeMetaFilePath)
			}
		} else {
			fmt.Println(string(bytes))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&codeMetaFilePath, "output", "o", "", "The path to the output codemeta.json file. If not specified, the output will be printed to the console.")
}
