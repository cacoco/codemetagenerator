package model

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestNewCodeMetaDefinition(t *testing.T) {
	g := gomega.NewWithT(t)

	id := "id"
	description := "description"
	name := "name"
	GivenName := "givenName"
	FamilyName := "familyName"
	Email := "email"
	Id := "id"
	maintainer := NewPerson(&GivenName, &FamilyName, &Email, &Id)
	codeRepository := "codeRepository"
	version := "version"
	developmentStatus := "developmentStatus"
	ProgrammingLanguageName := "name"
	ProgrammingLanguageURL := "url"
	programmingLanguage := NewProgrammingLanguage(&ProgrammingLanguageName, &ProgrammingLanguageURL)
	runtimePlatform := "runtimePlatform"
	license := "license"
	readme := "readme"

	var testbase = make(map[string]any)
	testbase[Identifier] = id
	testbase[Name] = name
	testbase[Description] = description
	testbase[Version] = version
	testbase[Maintainer] = maintainer
	testbase[ProgrammingLanguage] = programmingLanguage
	testbase[DevelopmentStatus] = developmentStatus
	testbase[License] = license
	testbase[RuntimePlatform] = runtimePlatform
	testbase[CodeRepository] = codeRepository
	testbase[Readme] = readme

	expected := map[string]any{
		Context:             DefaultContext,
		Type:                SoftwareSourceCodeType,
		Identifier:          id,
		Description:         description,
		Name:                name,
		Maintainer:          maintainer,
		CodeRepository:      codeRepository,
		Version:             version,
		DevelopmentStatus:   developmentStatus,
		ProgrammingLanguage: programmingLanguage,
		RuntimePlatform:     runtimePlatform,
		License:             license,
		Readme:              readme,
	}

	actual := NewCodemeta(&testbase)

	g.Ω(*actual).Should(gomega.Equal(expected))
}
