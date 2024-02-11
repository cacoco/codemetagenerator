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

func TestAddKeywords(t *testing.T) {
	g := gomega.NewWithT(t)

	temp := t.TempDir()
	// setup
	os.Mkdir(utils.GetHomeDir(temp), 0755)

	testMap := map[string]any{
		model.Context:               model.DefaultContext,
		model.Type:                  model.SoftwareSourceCodeType,
		model.Description:           "description",
		model.ContinuousIntegration: "https://url.org",
		model.Keywords:              []string{"one", "two"},
	}

	inProgressFilePath := utils.GetInProgressFilePath(temp)
	// need an in-progress code meta file
	err := utils.Marshal(inProgressFilePath, testMap)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	args := []string{"three", "four", "five"}

	writer := &utils.TestWriter{}

	keywords, err := keyword(writer, temp, args)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := []any{"one", "two", "three", "four", "five"}
	g.Ω(keywords).Should(gomega.Equal(expected))
}

func Test_ExecuteKeywordCmd(t *testing.T) {
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
	writer := utils.TestWriter{}

	cmdArgs := []string{"three", "four", "five"}
	keywordCmd := &cobra.Command{Use: "contributor", Args: cobra.MinimumNArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		_, err := keyword(&writer, temp, args)
		return err
	},
	}
	buf := bytes.NewBufferString("")
	keywordCmd.SetOut(buf)
	keywordCmd.SetErr(buf)
	keywordCmd.SetArgs(cmdArgs)

	err = keywordCmd.Execute()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// check file
	fileBytes, err := utils.LoadFile(utils.GetInProgressFilePath(temp))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	var m = make(map[string]any)
	oj.Unmarshal(fileBytes, &m)

	expected := map[string]any{
		model.Context:               model.DefaultContext,
		model.Type:                  model.SoftwareSourceCodeType,
		model.Description:           "description",
		model.ContinuousIntegration: "https://url.org",
		model.Keywords:              []any{"three", "four", "five"},
	}

	g.Ω(m).Should(gomega.Equal(expected))
}
