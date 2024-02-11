package utils

import (
	"fmt"
	"os"

	"github.com/cacoco/codemetagenerator/internal/model"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/oj"
	"github.com/samber/lo"
	"golang.org/x/exp/maps"
)

const (
	SPDXLicensesURL = "https://raw.githubusercontent.com/spdx/license-list-data/master/json/licenses.json"

	codemetaGeneratorDirectoryName = ".codemetagenerator"
	inProgressFilePath             = "/" + codemetaGeneratorDirectoryName + "/codemeta.inprogress.json"
	sPDXLicensesFilePath           = "/" + codemetaGeneratorDirectoryName + "/spdx-licenses.json"
)

var UserHomeDir, _ = getUserHomeDir()

func getUserHomeDir() (string, error) {
	return os.UserHomeDir()
}

func MkHomeDir(basedir string) error {
	homedir := basedir + "/" + codemetaGeneratorDirectoryName
	if _, err := os.Stat(homedir); os.IsNotExist(err) {
		err := os.Mkdir(homedir, 0755)
		if err != nil {
			return fmt.Errorf("unable to create codemetagenerator directory: %s", err.Error())
		}
	}
	return nil
}

func GetHomeDir(basedir string) string {
	return basedir + "/" + codemetaGeneratorDirectoryName
}

func GetInProgressFilePath(basedir string) string {
	return basedir + inProgressFilePath
}

func GetLicensesFilePath(basedir string) string {
	return basedir + sPDXLicensesFilePath
}

func ReadJSON(path string) (*string, error) {
	var p gen.Parser
	bytes, err := LoadFile(path)
	if err != nil {
		return nil, err
	}

	node, _ := p.Parse(bytes)
	json := oj.JSON(node, &oj.Options{Sort: true, Indent: 2, OmitNil: true})
	return &json, nil
}

func WriteJSON(path string, json string) error {
	return WriteFile(path, []byte(json))
}

func Unmarshal(path string) (*map[string]any, error) {
	bytes, err := LoadFile(path)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	oj.Unmarshal(bytes, &m)
	return &m, nil
}

func MarshalBytes(path string, bytes []byte, args ...any) error {
	var p gen.Parser
	node, err := p.Parse(bytes)
	if err != nil {
		return err
	}
	json := oj.JSON(node, &oj.Options{Sort: true, Indent: 2, OmitNil: true})
	return WriteJSON(path, json)
}

func Marshal(path string, m map[string]any, args ...any) error {
	bytes, err := oj.Marshal(m, args...)
	if err != nil {
		return err
	}

	return MarshalBytes(path, bytes, args...)
}

func LoadFile(filePath string) ([]byte, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func DeleteFile(path string) error {
	return os.Remove(path)
}

func WriteFile(path string, bytes []byte) error {
	return os.WriteFile(path, bytes, 0644)
}

// converts the full SPDX JSON file into a JSON file of licenseId => reference and store it
func CacheLicensesFile(basedir string, spdxFileBytes *[]byte, overwrite bool) error {
	// ensure we have a home directory
	err := MkHomeDir(basedir)
	if err != nil {
		return err
	}

	licensesFilePath := GetLicensesFilePath(basedir)
	if _, err = os.Stat(licensesFilePath); os.IsNotExist(err) || overwrite {
		var licensesList model.LicensesList
		err := oj.Unmarshal(*spdxFileBytes, &licensesList)
		if err != nil {
			return fmt.Errorf("unable to unmarshal SPDX licenses file: %s", err.Error())
		}

		var licensesMap map[string]any = make(map[string]any)
		lo.ForEach(licensesList.Licenses, func(license model.LicenseStruct, _ int) {
			licensesMap[license.LicenseId] = license.Reference
		})

		// marshal to file
		err = Marshal(licensesFilePath, licensesMap)
		if err != nil {
			return fmt.Errorf("unable to save translated SPDX licenses file: %s", err.Error())
		}
	}

	return nil
}

func GetSupportedLicenses(basedir string) (*[]string, error) {
	bytes, err := os.ReadFile(GetLicensesFilePath(basedir))
	if err != nil {
		return nil, err
	}
	var licenses map[string]string
	oj.Unmarshal(bytes, &licenses)

	keys := maps.Keys(licenses)
	return &keys, nil
}
