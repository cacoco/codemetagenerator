package internal

import (
	"reflect"
	"testing"
)

func TestNewCodeMetaDefinition(t *testing.T) {
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

	actual := NewCodeMetaDefinition(&id, &name, &description, &version, maintainer, programmingLanguage, developmentStatus, &license, &runtimePlatform, &codeRepository, &readme)

	matching := reflect.DeepEqual(*actual, expected)
	if !matching {
		t.Errorf("NewCodeMetaDefinition returned an unexpected value, got: %v, want: %v.", actual, expected)
	}
}
