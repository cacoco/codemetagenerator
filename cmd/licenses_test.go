package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

func Test_ExecuteLicensesCmd1(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	temp := t.TempDir()
	// setup
	os.Mkdir(utils.GetHomeDir(temp), 0755)
	file, le := os.ReadFile("../testdata/spdx-licenses.json")
	if le != nil {
		t.Errorf("Unexpected error: %v", le)
	}
	we := os.WriteFile(utils.GetLicensesFilePath(temp), file, 0644)
	if we != nil {
		t.Errorf("Unexpected error: %v", we)
	}

	writer := &utils.TestWriter{}

	var supported *[]string
	new := &cobra.Command{Use: "licenses", PreRunE: func(cmd *cobra.Command, args []string) error {
		res, err := getSupportedLicenses(temp)
		if err != nil {
			return err
		}
		supported = res
		return nil
	}, RunE: func(cmd *cobra.Command, args []string) error {
		return listLicenses(temp, *supported, writer)
	},
	}
	buf := bytes.NewBufferString("")
	new.SetOut(buf)
	new.SetErr(buf)
	new.SetArgs([]string{})

	err := new.Execute()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	g.Ω(len(*supported)).Should(gomega.Equal(624))
}
