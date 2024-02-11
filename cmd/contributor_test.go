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

func TestAddContributor(t *testing.T) {
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

	var stack utils.Stack[string]
	stack.Push("id\n")
	stack.Push("https://url.com\n")
	stack.Push("name\n")
	stack.Push("j\n") // down arrow to second option
	reader := utils.TestReader{In: utils.TestStdin{Data: stack}}

	writer := utils.TestWriter{}

	contributor, err := contributor(&reader, &writer, temp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := map[string]any{
		model.Type: model.OrganizationType,
		model.Name: "name",
		model.URL:  "https://url.com",
		model.Id:   "id",
	}
	g.Ω(*contributor).Should(gomega.Equal(expected))
}

func Test_ExecuteContributorCmd(t *testing.T) {
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

	var stack utils.Stack[string]
	stack.Push("id\n")
	stack.Push("https://url.com\n")
	stack.Push("name\n")
	stack.Push("j\n") // down arrow to second option
	reader := utils.TestReader{In: utils.TestStdin{Data: stack}}
	writer := utils.TestWriter{}

	contributorCmd := &cobra.Command{Use: "contributor", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, args []string) error {
		_, err := contributor(&reader, &writer, temp)
		return err
	},
	}
	buf := bytes.NewBufferString("")
	contributorCmd.SetOut(buf)
	contributorCmd.SetErr(buf)
	contributorCmd.SetArgs([]string{})

	err = contributorCmd.Execute()
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
		model.Contributor: []any{map[string]any{
			model.Type: model.OrganizationType,
			model.Name: "name",
			model.URL:  "https://url.com",
			model.Id:   "id",
		}},
	}

	g.Ω(m).Should(gomega.Equal(expected))
}
