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

	inProgressFilePath := utils.GetInProgressFilePath(temp)
	tempOutputFilePath := temp + "/codemeta.json"

	testMap := map[string]any{
		model.Context: model.DefaultContext,
		model.Type:    model.SoftwareSourceCodeType,
		"key":         "value",
		"key2":        []string{"one", "two"},
	}

	// need an in-progress code meta file
	err := utils.Marshal(inProgressFilePath, &testMap)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	writer := &utils.TestWriter{}

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
