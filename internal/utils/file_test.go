package utils

import (
	"os"
	"testing"

	"github.com/onsi/gomega"
)

func TestMkHomeDir(t *testing.T) {
	temp := t.TempDir()

	err := MkHomeDir(temp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// check directory  exists
	if _, err := os.Stat(GetHomeDir(temp)); err != nil {
		t.Errorf("Directory should have been created")
	}
}

func TestWriteReadJSON(t *testing.T) {
	g := gomega.NewWithT(t)

	temp := t.TempDir()
	path := temp + "/test.json"

	// write
	err := WriteJSON(path, `{"key": "value"}`)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// read
	data, err := ReadJSON(path)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	// data in file is pretty printed
	g.Ω(*data).Should(gomega.Equal("{\n  \"key\": \"value\"\n}"))
}

func TestMarshalUnMasrshal(t *testing.T) {
	g := gomega.NewWithT(t)

	temp := t.TempDir()
	path := temp + "/test.json"

	m := map[string]any{"key": "value"}

	// marshal to file
	err := Marshal(path, m)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// unmarshal from file
	data, err := Unmarshal(path)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	g.Ω(*data).Should(gomega.Equal(m))
}

func TestDeleteFile(t *testing.T) {
	temp := t.TempDir()
	path := temp + "/test.json"

	// write
	err := WriteJSON(path, `{"key": "value"}`)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// delete
	err = DeleteFile(path)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// check file does not exist
	if _, err := os.Stat(path); err == nil {
		t.Errorf("File should have been deleted")
	}
}

func TestGetSupportedLicenses(t *testing.T) {
	g := gomega.NewWithT(t)

	temp := t.TempDir()
	file, err := os.ReadFile("../../testdata/spdx-full-licenses.json")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	// get and store
	err = CacheLicensesFile(temp, &file, false)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// get supported licenses
	res, err := GetSupportedLicenses(temp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	g.Ω(len(*res)).Should(gomega.Equal(624))
}

func TestCacheLicensesFile(t *testing.T) {
	temp := t.TempDir()
	file, le := os.ReadFile("../../testdata/spdx-full-licenses.json")
	if le != nil {
		t.Errorf("Unexpected error: %v", le)
	}

	// get and store
	err := CacheLicensesFile(temp, &file, false)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// check file exists
	if _, err := os.Stat(GetLicensesFilePath(temp)); err != nil {
		t.Errorf("Licenses file should have been created")
	}
}
