package cmd

import (
	"fmt"
	"os"

	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/spf13/cobra"
)

func clean(basedir string, writer utils.Writer) error {
	codemetageneratorHomeDir := utils.GetHomeDir(basedir)
	if _, err := os.Stat(codemetageneratorHomeDir); err == nil {
		err := os.RemoveAll(codemetageneratorHomeDir)
		if err != nil {
			handleErr(writer, err)
			return fmt.Errorf("unable to remove codemetagenerator home directory")
		}
		writer.Println("✅ Successfully cleaned the $HOME/.codemetagenerator directory.")
	}
	return nil
}

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Args:  cobra.NoArgs,
	Short: "Clean the $HOME/.codemetagenerator directory",
	Long: `
Removes the $HOME/.codemetagenerator directory used to store the in-progress codemeta.json file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return clean(utils.UserHomeDir, &utils.StdoutWriter{})
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
