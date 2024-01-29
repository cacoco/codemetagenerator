package codemetagenerator

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Args:  cobra.NoArgs,
	Short: "Clean the $HOME/.codemetagenerator directory",
	Long: `Removes the $HOME/.codemetagenerator directory used to store the in-progress 
codemeta.json file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return errors.New("unable to find $HOME directory")
		}
		codemetageneratorHomeDir := homeDir + "/.codemetagenerator"
		if _, err := os.Stat(codemetageneratorHomeDir); err == nil {
			err := os.RemoveAll(codemetageneratorHomeDir)
			if err != nil {
				return fmt.Errorf("unable to remove codemetagenerator home directory: %s", err.Error())
			}
			fmt.Println("✅ Successfully cleaned the $HOME/.codemetagenerator directory.")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
