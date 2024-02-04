package cmd

import (
	"testing"

	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/onsi/gomega"
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
