package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/cacoco/codemetagenerator/internal/model"
	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/ohler55/ojg/oj"
	"github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

func Test_ExecuteGenerateCmd(t *testing.T) {
	g := gomega.NewWithT(t)

	temp := t.TempDir()
	// setup
	os.Mkdir(utils.GetHomeDir(temp), 0755)

	testMap := map[string]any{
		model.Context:               model.DefaultContext,
		model.Type:                  model.SoftwareSourceCodeType,
		model.Description:           "description",
		model.ContinuousIntegration: "https://url.org",
	}

	inProgressFilePath := utils.GetInProgressFilePath(temp)
	// need an in-progress code meta file
	err := utils.Marshal(inProgressFilePath, testMap)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	writer := &utils.TestWriter{}

	tempOutputFilePath := temp + "/codemeta.json"
	generate := &cobra.Command{Use: "generate", RunE: func(cmd *cobra.Command, args []string) error {
		return generate(temp, writer, tempOutputFilePath)
	},
	}
	buf := bytes.NewBufferString("")
	generate.SetOut(buf)
	generate.SetErr(buf)
	generate.SetArgs([]string{})

	err = generate.Execute()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// check files are the same
	fileBytes, le := utils.LoadFile(inProgressFilePath)
	if le != nil {
		t.Errorf("Unexpected error: %v", le)
	}
	var inprogress = make(map[string]any)
	oj.Unmarshal(fileBytes, &inprogress)

	outputBytes, ole := utils.LoadFile(tempOutputFilePath)
	if ole != nil {
		t.Errorf("Unexpected error: %v", ole)
	}
	var output = make(map[string]any)
	oj.Unmarshal(outputBytes, &output)

	// actual should equal expected
	g.Ω(output).Should(gomega.Equal(inprogress))
}

func Test_ExecuteGenerateCmdBadFile(t *testing.T) {
	temp := t.TempDir()
	// setup
	os.Mkdir(utils.GetHomeDir(temp), 0755)

	testMap := map[string]any{
		model.Context: model.DefaultContext,
		model.Type:    model.SoftwareSourceCodeType,
		"key":         "not a valid codemeta term",
		"key2":        []string{"also", "not", "a", "valid", "codemeta", "term"},
	}

	inProgressFilePath := utils.GetInProgressFilePath(temp)
	// need an in-progress code meta file
	err := utils.Marshal(inProgressFilePath, testMap)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	writer := &utils.TestWriter{}

	generate := &cobra.Command{Use: "generate", RunE: func(cmd *cobra.Command, args []string) error {
		return generate(temp, writer, "")
	},
	}
	buf := bytes.NewBufferString("")
	generate.SetOut(buf)
	generate.SetErr(buf)
	generate.SetArgs([]string{})

	err = generate.Execute()
	if err == nil {
		t.Errorf("Expected error")
	}
}
