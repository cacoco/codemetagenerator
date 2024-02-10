package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/cacoco/codemetagenerator/internal/model"
	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/spf13/cobra"
)

func Test_ExecuteValidateCmd(t *testing.T) {
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

	validate := &cobra.Command{Use: "validate", RunE: func(cmd *cobra.Command, args []string) error {
		return validate(temp, writer, "")
	},
	}
	buf := bytes.NewBufferString("")
	validate.SetOut(buf)
	validate.SetErr(buf)
	validate.SetArgs([]string{})

	err = validate.Execute()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func Test_ExecuteValidateCmdBadFile(t *testing.T) {
	temp := t.TempDir()
	// setup
	os.Mkdir(utils.GetHomeDir(temp), 0755)

	testMap := map[string]any{
		model.Context:               model.DefaultContext,
		model.Type:                  model.SoftwareSourceCodeType,
		model.Description:           "description",
		model.ContinuousIntegration: "https://url.org",
		"key":                       "notvalid",
	}

	inProgressFilePath := utils.GetInProgressFilePath(temp)
	// need an in-progress code meta file
	err := utils.Marshal(inProgressFilePath, testMap)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	writer := &utils.TestWriter{}

	validate := &cobra.Command{Use: "validate", RunE: func(cmd *cobra.Command, args []string) error {
		return validate(temp, writer, "")
	},
	}
	buf := bytes.NewBufferString("")
	validate.SetOut(buf)
	validate.SetErr(buf)
	validate.SetArgs([]string{})

	err = validate.Execute()
	if err == nil {
		t.Errorf("expected error for validate command")
	}
}
