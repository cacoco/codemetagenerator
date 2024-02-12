package cmd

import (
	"os"
	"testing"

	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/ohler55/ojg/oj"
	"github.com/onsi/gomega"
	"golang.org/x/exp/maps"
)

func TestGetLicenseReference(t *testing.T) {
	g := gomega.NewWithT(t)

	temp := t.TempDir()
	// setup
	os.Mkdir(utils.GetHomeDir(temp), 0755)
	file, le := os.ReadFile("../testdata/spdx-licenses.json")
	if le != nil {
		t.Errorf("Unexpected error: %v", le)
	}
	we := utils.WriteFile(utils.GetLicensesFilePath(temp), file)
	if we != nil {
		t.Errorf("Unexpected error: %v", we)
	}
	writer := utils.TestWriter{}

	reference, err := getLicenseReference(writer, temp, "Apache-2.0")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	g.Ω(*reference).Should(gomega.Equal("https://spdx.org/licenses/Apache-2.0.html"))
}

func reset() {
	SupportedLicenses = Licenses{}
}

func TestValidateLicenseId1(t *testing.T) {
	g := gomega.NewWithT(t)

	temp := t.TempDir()
	// setup
	os.Mkdir(utils.GetHomeDir(temp), 0755)
	file, err := os.ReadFile("../testdata/spdx-licenses.json")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = utils.WriteFile(utils.GetLicensesFilePath(temp), file)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	supported, err := utils.GetSupportedLicenses(temp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	SupportedLicenses.setSupportedLicenses(*supported)
	defer reset() // make sure to reset the global variable
	writer := utils.TestWriter{}

	var validateFn = validateLicenseId(writer, temp)
	err = validateFn("Apache-2.0")
	g.Expect(err).To(gomega.BeNil())
}

func TestValidateLicenseId2(t *testing.T) {
	g := gomega.NewWithT(t)

	temp := t.TempDir()
	writer := utils.TestWriter{}

	// SupportedLicenses is nil, should error
	var validateFn = validateLicenseId(writer, temp)
	err := validateFn("Apache-2.0")
	g.Expect(err).ToNot(gomega.BeNil())
}

func TestLoadSupportedLicenses(t *testing.T) {
	g := gomega.NewWithT(t)

	temp := t.TempDir()
	// setup
	os.Mkdir(utils.GetHomeDir(temp), 0755)
	file, err := os.ReadFile("../testdata/spdx-full-licenses.json")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	var stack utils.Stack[string]
	stack.Push(string(file))
	httpClient := utils.NewTestHttpClient(&stack)

	actual, err := loadSupportedLicenses(temp, httpClient)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// read converted test file
	var expectedMap map[string]string = make(map[string]string)
	file, err = os.ReadFile("../testdata/spdx-licenses.json")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	oj.Unmarshal(file, &expectedMap)
	expected := maps.Keys(expectedMap)

	// check "downloaded" file against converted test file
	g.Ω(*actual).Should(gomega.ConsistOf(expected))
}
