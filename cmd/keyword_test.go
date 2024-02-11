package cmd

import (
	"os"
	"testing"

	"github.com/cacoco/codemetagenerator/internal/model"
	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/onsi/gomega"
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

	keywords, err := addKeywords(writer, temp, args)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := []any{"one", "two", "three", "four", "five"}
	g.Ω(keywords).Should(gomega.Equal(expected))
}
