package cmd

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/spf13/cobra"
)

var Debug bool

func SetVersionInfo(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s (Built on %s from Git SHA %s)", version, date, commit)
}

func loadSupportedLicenses(basedir string, httpClient *http.Client) (*[]string, error) {
	// if the licenses file doesn't exist, download a new SPDX file and cache it
	if _, err := os.Stat(utils.GetLicensesFilePath(basedir)); os.IsNotExist(err) {
		err := downloadSPDXLicenses(basedir, httpClient, false)
		if err != nil {
			return nil, err
		}
	}
	return utils.GetSupportedLicenses(basedir)
}

// downloads and caches a translation of SPDX licenses JSON
func downloadSPDXLicenses(basedir string, httpClient *http.Client, overwrite bool) error {
	// download and cache the licenses file
	request, err := utils.MkJSONRequest(http.MethodGet, utils.SPDXLicensesURL)
	if err != nil {
		return err
	}
	bytes, err := utils.DoRequest(httpClient, request)
	if err != nil {
		return err
	}
	return utils.CacheLicensesFile(basedir, bytes, true)
}

type Licenses struct {
	m                 sync.Mutex
	supportedLicenses []string
}

var SupportedLicenses Licenses = Licenses{}

func (l *Licenses) getSupportedLicenses() []string {
	l.m.Lock()
	defer l.m.Unlock()
	return l.supportedLicenses
}

func (l *Licenses) setSupportedLicenses(supportedLicenses []string) {
	l.m.Lock()
	defer l.m.Unlock()
	l.supportedLicenses = supportedLicenses
}

func handleErr(writer utils.Writer, err error) {
	if err != nil {
		if Debug {
			fmt.Fprintln(writer.StdErr(), "Error:", err.Error())
		}
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "codemetagenerator",
	Short: "An interactive codemeta.json file generator for projects",
	Long: `
CodeMeta (https://codemeta.github.io) is a JSON-LD file format used to describe software projects. 
'codemetagenerator' is an interactive tool that helps you generate a valid 'codemeta.json' file.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		foundLicenses, err := loadSupportedLicenses(utils.UserHomeDir, utils.MkHttpClient())
		if err != nil {
			return fmt.Errorf("unable to retrieve supported licenses: %s", err.Error())
		}
		SupportedLicenses.setSupportedLicenses(*foundLicenses)
		return nil
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

	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "enable debug mode")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
