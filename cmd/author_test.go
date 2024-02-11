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

func TestAddAuthor(t *testing.T) {
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
	stack.Push("person@email.org\n")
	stack.Push("familyName\n")
	stack.Push("givenName\n")
	stack.Push("\n") // enter to select the first option
	reader := utils.TestReader{In: utils.TestStdin{Data: stack}}
	writer := utils.TestWriter{}

	author, err := author(&reader, &writer, temp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := map[string]any{
		model.Type:       model.PersonType,
		model.FamilyName: "familyName",
		model.GivenName:  "givenName",
		model.Email:      "person@email.org",
		model.Id:         "id",
	}
	g.Ω(*author).Should(gomega.Equal(expected))
}

func Test_ExecuteAuthorCmd(t *testing.T) {
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
	stack.Push("person@email.org\n")
	stack.Push("familyName\n")
	stack.Push("givenName\n")
	stack.Push("\n") // enter to select the first option
	reader := utils.TestReader{In: utils.TestStdin{Data: stack}}
	writer := utils.TestWriter{}

	authorCmd := &cobra.Command{Use: "author", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, args []string) error {
		_, err := author(&reader, &writer, temp)
		return err
	},
	}
	buf := bytes.NewBufferString("")
	authorCmd.SetOut(buf)
	authorCmd.SetErr(buf)
	authorCmd.SetArgs([]string{})

	err = authorCmd.Execute()
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
		model.Author: []any{map[string]any{
			model.Type:       model.PersonType,
			model.FamilyName: "familyName",
			model.GivenName:  "givenName",
			model.Email:      "person@email.org",
			model.Id:         "id",
		}},
	}

	g.Ω(m).Should(gomega.Equal(expected))
}
