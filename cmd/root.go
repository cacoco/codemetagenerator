package cmd

import (
	"fmt"
	"os"

	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/spf13/cobra"
)

func SetVersionInfo(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s (Built on %s from Git SHA %s)", version, date, commit)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "codemetagenerator",
	Short: "An interactive codemeta.json file generator for projects",
	Long: `
CODEMETA JSON FILE GENERATOR
----------------------------
CodeMeta (https://codemeta.github.io) is a JSON-LD file format used to describe software projects. 
This is an interactive tool that helps you generate a 'codemeta.json' file for your project.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		go utils.GetAndCacheLicenseFile(utils.UserHomeDir, false)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.CompletionOptions.DisableNoDescFlag = true
	rootCmd.CompletionOptions.DisableDescriptions = true
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
