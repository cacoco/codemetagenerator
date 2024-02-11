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

var testWriter = &utils.TestWriter{}

func TestDeleteStringValue(t *testing.T) {
	// initialize gomega
	g := gomega.NewWithT(t)

	json := `{"key":"value","key2":"value2"}`
	path := "key2"

	result, err := deleteValue(testWriter, []byte(json), path)
	if err != nil {
		t.Errorf("Deleting property with path `%s` returned unexpected error: %v", path, err)
	}
	expected := `{"key":"value"}`
	compare(g, *result, expected)
}

func TestDeleteIntValue(t *testing.T) {
	// initialize gomega
	g := gomega.NewWithT(t)

	json := `{"key":"value","key2":2}`
	path := "key2"
	result, err := deleteValue(testWriter, []byte(json), path)
	if err != nil {
		t.Errorf("Setting property with path `%s` returned unexpected error: %v", path, err)
	}
	expected := `{"key": "value"}`
	compare(g, *result, expected)
}

func TestDeleteStructValue(t *testing.T) {
	// initialize gomega
	g := gomega.NewWithT(t)

	json := `{"key":"value","key2":{"first":"value3","second":"value4"}}`
	path := "key2"
	result, err := deleteValue(testWriter, []byte(json), path)
	if err != nil {
		t.Errorf("Setting property with path `%s` returned unexpected error: %v", path, err)
	}
	expected := `{"key": "value"}`
	compare(g, *result, expected)
}

func TestDeleteItemInArrayValue1(t *testing.T) {
	// initialize gomega
	g := gomega.NewWithT(t)

	json := `{"key":"value","key2":[1,2,3]}`
	path := "key2.-1"
	result, err := deleteValue(testWriter, []byte(json), path)
	if err != nil {
		t.Errorf("Setting property with path `%s` returned unexpected error: %v", path, err)
	}
	expected := `{"key":"value","key2":[1,2]}`
	compare(g, *result, expected)
}

func TestDeleteItemInArrayValue2(t *testing.T) {
	// initialize gomega
	g := gomega.NewWithT(t)

	json := `{"key":"value","key2":[1,2,3]}`
	path := "key2.0"
	result, err := deleteValue(testWriter, []byte(json), path)
	if err != nil {
		t.Errorf("Setting property with path `%s` returned unexpected error: %v", path, err)
	}
	expected := `{"key":"value","key2":[2,3]}`
	compare(g, *result, expected)
}

func TestDeleteItemInArrayValue3(t *testing.T) {
	// initialize gomega
	g := gomega.NewWithT(t)

	json := `{"key":"value","key2":[1,2,3]}`
	path := "key2.2"
	result, err := deleteValue(testWriter, []byte(json), path)
	if err != nil {
		t.Errorf("Setting property with path `%s` returned unexpected error: %v", path, err)
	}
	expected := `{"key":"value","key2":[1,2]}`
	compare(g, *result, expected)
}

func TestDeleteArrayValue(t *testing.T) {
	// initialize gomega
	g := gomega.NewWithT(t)

	json := `{"key":"value","key2":[1,2,3]}`
	path := "key2"
	result, err := deleteValue(testWriter, []byte(json), path)
	if err != nil {
		t.Errorf("Setting property with path `%s` returned unexpected error: %v", path, err)
	}
	expected := `{"key":"value"}`
	compare(g, *result, expected)
}

func Test_ExecuteDeleteCmd(t *testing.T) {
	// initialize gomega
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

	deleteCmd := &cobra.Command{Use: "delete", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		return delete(&writer, temp, args[0])
	},
	}
	buf := bytes.NewBufferString("")
	deleteCmd.SetOut(buf)
	deleteCmd.SetErr(buf)
	deleteCmd.SetArgs([]string{model.Description})

	err = deleteCmd.Execute()
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
		model.ContinuousIntegration: "https://url.org",
	}

	g.Ω(m).Should(gomega.Equal(expected))
}
