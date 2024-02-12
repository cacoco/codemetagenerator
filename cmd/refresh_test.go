package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/ohler55/ojg/oj"
	"github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

func Test_ExecuteRefreshCmd(t *testing.T) {
	g := gomega.NewWithT(t)

	temp := t.TempDir()
	// setup
	os.Mkdir(utils.GetHomeDir(temp), 0755)
	file, err := os.ReadFile("../testdata/spdx-full-licenses.json")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	var stack utils.Stack[string]
	stack.Push(string(file))
	httpClient := utils.NewTestHttpClient(&stack)
	writer := utils.TestWriter{}

	refreshCmd := &cobra.Command{Use: "refresh", RunE: func(cmd *cobra.Command, args []string) error {
		return refresh(writer, temp, httpClient)
	},
	}
	buf := bytes.NewBufferString("")
	refreshCmd.SetOut(buf)
	refreshCmd.SetErr(buf)
	refreshCmd.SetArgs([]string{})

	err = refreshCmd.Execute()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// read converted test file
	var expected map[string]string = make(map[string]string)
	file, err = os.ReadFile("../testdata/spdx-licenses.json")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	oj.Unmarshal(file, &expected)

	var actual map[string]string = make(map[string]string)
	fileBytes, err := utils.LoadFile(utils.GetLicensesFilePath(temp))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	oj.Unmarshal(fileBytes, &actual)

	// check "downloaded" file against converted test file
	g.Ω(actual).Should(gomega.Equal(expected))
}
