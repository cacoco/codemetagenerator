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

func Test_ExecuteNewCmd1(t *testing.T) {
	g := gomega.NewWithT(t)

	temp := t.TempDir()
	// setup
	os.Mkdir(utils.GetHomeDir(temp), 0755)
	file, err := os.ReadFile("../testdata/spdx-licenses.json")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = os.WriteFile(utils.GetLicensesFilePath(temp), file, 0644)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	supportedLicenses, err := utils.GetSupportedLicenses(temp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	SupportedLicenses.setSupportedLicenses(*supportedLicenses)

	givenName := "givenName"
	familyName := "familyName"
	email := "person@email.org"
	id := "id"
	programmingLanguageURL := "https://programmingLanguageURL.com"
	programmingLanguageName := "programmingLanguageName"

	var stack utils.Stack[string]
	stack.Push(id + "\n")
	stack.Push(email + "\n")
	stack.Push(familyName + "\n")
	stack.Push(givenName + "\n")
	stack.Push("\n") // select first option for maintainer person or org prompt -- person
	stack.Push("https://readme.com\n")
	stack.Push("Apache-2.0\n")
	stack.Push("version\n")
	stack.Push("runtimePlatform\n")
	stack.Push(programmingLanguageURL + "\n")
	stack.Push(programmingLanguageName + "\n")
	stack.Push("https://codeRepository.org\n")
	stack.Push("\n") // select first developmentStatus, Abandoned
	stack.Push("description\n")
	stack.Push("name\n")
	stack.Push("identifier\n")
	reader := utils.TestReader{In: utils.TestStdin{Data: stack}}

	writer := utils.TestWriter{}

	new := &cobra.Command{Use: "new", RunE: func(cmd *cobra.Command, args []string) error {
		return new(temp, &reader, &writer, "")
	},
	}
	buf := bytes.NewBufferString("")
	new.SetOut(buf)
	new.SetErr(buf)
	new.SetArgs([]string{})

	err = new.Execute()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// check file
	fileBytes, le := utils.LoadFile(utils.GetInProgressFilePath(temp))
	if le != nil {
		t.Errorf("Unexpected error: %v", le)
	}
	var m = make(map[string]any)
	oj.Unmarshal(fileBytes, &m)

	maintainer := model.NewPerson(&givenName, &familyName, &email, &id)
	programmingLang := model.NewProgrammingLanguage(&programmingLanguageName, &programmingLanguageURL)

	expected := map[string]any{
		model.Context:             model.DefaultContext,
		model.Type:                model.SoftwareSourceCodeType,
		model.Identifier:          "identifier",
		model.Name:                "name",
		model.Description:         "description",
		model.Version:             "version",
		model.Maintainer:          *maintainer,
		model.ProgrammingLanguage: *programmingLang,
		model.DevelopmentStatus:   "Abandoned",
		model.License:             "https://spdx.org/licenses/Apache-2.0.html",
		model.RuntimePlatform:     "runtimePlatform",
		model.CodeRepository:      "https://codeRepository.org",
		model.Readme:              "https://readme.com",
	}

	g.Ω(m).Should(gomega.Equal(expected))
}

func Test_ExecuteNewCmd2(t *testing.T) {
	g := gomega.NewWithT(t)

	temp := t.TempDir()
	// setup
	os.Mkdir(utils.GetHomeDir(temp), 0755)
	file, err := os.ReadFile("../testdata/spdx-licenses.json")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = os.WriteFile(utils.GetLicensesFilePath(temp), file, 0644)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	supportedLicenses, err := utils.GetSupportedLicenses(temp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	SupportedLicenses.setSupportedLicenses(*supportedLicenses)

	reader := utils.TestReader{In: utils.TestStdin{Data: *utils.NilStack}}
	writer := utils.TestWriter{}

	new := &cobra.Command{Use: "new", RunE: func(cmd *cobra.Command, args []string) error {
		return new(temp, &reader, &writer, "../testdata/testmeta.json")
	},
	}
	buf := bytes.NewBufferString("")
	new.SetOut(buf)
	new.SetErr(buf)
	new.SetArgs([]string{})

	err = new.Execute()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// check file
	fileBytes, le := utils.LoadFile(utils.GetInProgressFilePath(temp))
	if le != nil {
		t.Errorf("Unexpected error: %v", le)
	}
	var m = make(map[string]any)
	oj.Unmarshal(fileBytes, &m)

	expected := map[string]any{
		model.Context:     model.DefaultContext,
		model.Type:        model.SoftwareSourceCodeType,
		model.Identifier:  "testmeta",
		model.Name:        "TestMeta",
		model.Description: "A test codemeta.json file.",
	}

	g.Ω(m).Should(gomega.Equal(expected))
}
