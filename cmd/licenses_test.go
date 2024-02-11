package cmd

import (
	"os"
	"testing"

	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/onsi/gomega"
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

	reference, err := getLicenseReference(temp, "Apache-2.0")
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

	var validateFn = validateLicenseId(temp)
	err = validateFn("Apache-2.0")
	g.Expect(err).To(gomega.BeNil())
}

func TestValidateLicenseId2(t *testing.T) {
	g := gomega.NewWithT(t)

	temp := t.TempDir()

	// SupportedLicenses is nil, should error
	var validateFn = validateLicenseId(temp)
	err := validateFn("Apache-2.0")
	g.Expect(err).ToNot(gomega.BeNil())
}
