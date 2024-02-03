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

func compare(g *gomega.WithT, actual string, expected string) error {
	var actualMap map[string]any
	oj.Unmarshal([]byte(actual), &actualMap)

	var expectedMap map[string]any
	oj.Unmarshal([]byte(expected), &expectedMap)

	g.Ω(actualMap).Should(gomega.Equal(expectedMap))
	return nil
}

func TestSetStringValue(t *testing.T) {
	// initialize gomega
	g := gomega.NewWithT(t)

	json := `{"key": "value"}`
	path := "key2"
	result, err := setValue([]byte(json), path, "\"value2\"")
	if err != nil {
		t.Errorf("Setting property with path `%s` returned unexpected error: %v", path, err)
	}
	expected := `{"key":"value","key2":"value2"}`
	compare(g, *result, expected)
}

func TestSetIntValue(t *testing.T) {
	// initialize gomega
	g := gomega.NewWithT(t)

	json := `{"key": "value"}`
	path := "key2"
	result, err := setValue([]byte(json), path, "2")
	if err != nil {
		t.Errorf("Setting property with path `%s` returned unexpected error: %v", path, err)
	}
	expected := `{"key":"value","key2":2}`
	compare(g, *result, expected)
}

func TestSetStructValue(t *testing.T) {
	// initialize gomega
	g := gomega.NewWithT(t)

	json := `{"key": "value"}`
	path := "key2"
	result, err := setValue([]byte(json), path, `{"first": "value3", "second": "value4"}`)
	if err != nil {
		t.Errorf("Setting property with path `%s` returned unexpected error: %v", path, err)
	}
	expected := `{"key":"value","key2":{"first":"value3","second":"value4"}}`
	compare(g, *result, expected)
}

func TestSetArrayValue(t *testing.T) {
	// initialize gomega
	g := gomega.NewWithT(t)

	json := `{"key": "value"}`
	path := "key2"
	result, err := setValue([]byte(json), path, `[1, 2, 3]`)
	if err != nil {
		t.Errorf("Setting property with path `%s` returned unexpected error: %v", path, err)
	}
	expected := `{"key":"value","key2":[1,2,3]}`
	compare(g, *result, expected)
}

func TestSetItemInArrayValue1(t *testing.T) {
	// initialize gomega
	g := gomega.NewWithT(t)

	json := `{"key": "value", "key2":[1,2,3]}`
	path := "key2.-1" // -1 represents append as the last element of an array
	result, err := setValue([]byte(json), path, `4`)
	if err != nil {
		t.Errorf("Setting property with path `%s` returned unexpected error: %v", path, err)
	}
	expected := `{"key":"value","key2":[1,2,3,4]}`
	compare(g, *result, expected)
}

func TestSetItemInArrayValue2(t *testing.T) {
	// initialize gomega
	g := gomega.NewWithT(t)

	json := `{"key": "value", "key2":[1,2,3]}`
	path := "key2.0"
	result, err := setValue([]byte(json), path, `0`)
	if err != nil {
		t.Errorf("Setting property with path `%s` returned unexpected error: %v", path, err)
	}
	expected := `{"key":"value","key2":[0,2,3]}`
	compare(g, *result, expected)
}

func TestSetItemInArrayValue3(t *testing.T) {
	// initialize gomega
	g := gomega.NewWithT(t)

	json := `{"key": "value", "key2":[1,2,3]}`
	path := "key2.2"
	result, err := setValue([]byte(json), path, `0`)
	if err != nil {
		t.Errorf("Setting property with path `%s` returned unexpected error: %v", path, err)
	}
	expected := `{"key":"value","key2":[1,2,0]}`
	compare(g, *result, expected)
}

func Test_ExecuteSetCmd(t *testing.T) {
	g := gomega.NewWithT(t)

	// setup
	temp := t.TempDir()
	os.Mkdir(utils.GetHomeDir(temp), 0755)

	inProgressFilePath := utils.GetInProgressFilePath(temp)

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

	set := &cobra.Command{Use: "set", RunE: func(cmd *cobra.Command, args []string) error {
		return set(temp, args)
	},
	}
	buf := bytes.NewBufferString("")
	set.SetOut(buf)
	set.SetErr(buf)
	set.SetArgs([]string{"key3", "\"topic2\""})

	err = set.Execute()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	// check file
	fileBytes, le := utils.LoadFile(inProgressFilePath)
	if le != nil {
		t.Errorf("Unexpected error: %v", le)
	}
	var m = make(map[string]any)
	oj.Unmarshal(fileBytes, &m)
	g.Ω(m["key3"]).ShouldNot(gomega.BeNil())
	g.Ω(m["key3"]).ShouldNot(gomega.BeEmpty())
	g.Ω(m["key3"]).Should(gomega.Equal("topic2"))
}
