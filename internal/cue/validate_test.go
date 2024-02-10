package cue

import (
	"fmt"
	"testing"

	"github.com/cacoco/codemetagenerator/internal/utils"
)

func TestValidate(t *testing.T) {
	bytes := []byte(`
{
	"@context": "https://w3id.org/codemeta/3.0",
	"@type": "SoftwareSourceCode",
	"identifier": "testmeta",
	"description": "A test codemeta.json file.",
	"name": "TestMeta"
}`)

	err := Validate(bytes)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestValidateType1(t *testing.T) {
	bytes := []byte(`
{
	"@context": "https://w3id.org/codemeta/3.0",
	"@type":    "NOTVALID"
}`)

	err := Validate(bytes)
	if err == nil {
		t.Errorf("Expected error")
	}
}

func TestValidateType2(t *testing.T) {
	bytes := []byte(`
{
	"@context": "https://w3id.org/codemeta/3.0",
	"@type": "SoftwareSourceCode"
}`)

	err := Validate(bytes)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestValidateEmail1(t *testing.T) {
	bytes := []byte(`
{
	"@context": "https://w3id.org/codemeta/3.0",
	"@type": "SoftwareSourceCode",
	"author": [
		{
			"@type":"Person",
			"familyName":"Smith",
			"givenName":"Alice",
			"email":"NOTVALID"
			"@id":"https://orcid.org/0000-0000-0000-0000"
		}
	]
}`)

	err := Validate(bytes)
	if err == nil {
		t.Errorf("Expected error")
	}
}

func TestValidateEmail2(t *testing.T) {
	bytes := []byte(`
{
	"@context": "https://w3id.org/codemeta/3.0",
	"@type": "SoftwareSourceCode",
	"author": [
		{
			"@type":"Person",
			"familyName":"Smith",
			"givenName":"Alice",
			"email":"asmith@person.org",
			"@id":"https://orcid.org/0000-0000-0000-0000"
		}
	]
}`)

	err := Validate(bytes)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func testTime(timestamp string) error {
	tc := "\"" + timestamp + "\""
	if timestamp == "null" {
		tc = "null"
	}
	bytes := []byte(fmt.Sprintf(`
{
	"@context": "https://w3id.org/codemeta/3.0",
	"@type":    "SoftwareSourceCode",
	"dateCreated": %s       
}`, tc))

	err := Validate(bytes)
	if err != nil {
		return err
	}
	return nil
}

func TestValidateDate(t *testing.T) {
	validTimes := []string{
		// valid times
		"2019-01-02T15:04:05Z",
		"2019-01-02T15:04:05-08:00",
		"2019-01-02T15:04:05.0-08:00",
		"2019-01-02T15:04:05.01-08:00",
		"2019-01-02T15:04:05.012345678-08:00",
		"2019-02-28T15:04:59Z",
		"2019-01-02T15:04:05.01234567890-08:00",
	}

	for _, tc := range validTimes {
		err := testTime(tc)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	}

	invalidTimes := []string{
		"2019-01-02T15:04:05",        // missing time zone
		"2019-01-02T15:04:61Z",       // seconds out of range
		"2019-01-02T15:60:00Z",       // minute out of range
		"2019-01-02T24:00:00Z",       // hour out of range
		"2019-01-32T23:00:00Z",       // day out of range
		"2019-01-00T23:00:00Z",       // day out of range
		"2019-00-15T23:00:00Z",       // month out of range
		"2019-13-15T23:00:00Z",       // month out of range
		"2019-01-02T15:04:05Z+08:00", // double time zone
		"2019-01-02T15:04:05+08",     // partial time zone
	}

	for _, tc := range invalidTimes {
		err := testTime(tc)
		if err == nil {
			t.Errorf("Expected error")
		}
	}
}

func TestValidateUrl1(t *testing.T) {
	bytes := []byte(`
{
	"@context": "https://w3id.org/codemeta/3.0",
	"@type": "SoftwareSourceCode",
	"codeRepository": "https://github.com/codemeta/codemeta"
}`)

	err := Validate(bytes)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestValidateUrl2(t *testing.T) {
	bytes := []byte(`
{
	"@context": "https://w3id.org/codemeta/3.0",
	"@type": "SoftwareSourceCode",
	"codeRepository": "ht:/nopeynopenopsi?/ssjshds"
}`)

	err := Validate(bytes)
	if err == nil {
		t.Errorf("Expected error")
	}
}

func TestCodeMetaFile(t *testing.T) {
	bytes, err := utils.LoadFile("../../testdata/codemeta.json")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	err = Validate(bytes)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestCodemetaRFile(t *testing.T) {
	bytes, err := utils.LoadFile("../../testdata/codemetaR.json")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	err = Validate(bytes)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
