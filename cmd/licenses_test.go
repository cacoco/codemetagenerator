package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

func Test_ExecuteLicensesCmd(t *testing.T) {
	g := gomega.NewWithT(t)

	temp := t.TempDir()
	// setup
	os.Mkdir(utils.GetHomeDir(temp), 0755)
	file, err := os.ReadFile("../testdata/spdx-licenses.json")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = utils.WriteFile(utils.GetLicensesFilePath(temp), file)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	supported, err := utils.GetSupportedLicenses(temp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	SupportedLicenses.setSupportedLicenses(*supported)
	defer reset() // make sure to reset the global variable
	writer := utils.TestWriter{}

	var list []string
	licensesCmd := &cobra.Command{Use: "licenses", RunE: func(cmd *cobra.Command, args []string) error {
		list, err = licenses(&writer)
		return err
	},
	}
	buf := bytes.NewBufferString("")
	licensesCmd.SetOut(buf)
	licensesCmd.SetErr(buf)
	licensesCmd.SetArgs([]string{})

	err = licensesCmd.Execute()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	g.Ω(list).Should(gomega.Equal(*supported))
}
