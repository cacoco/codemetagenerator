package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/cacoco/codemetagenerator/internal/model"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/oj"
	"github.com/samber/lo"
	"golang.org/x/exp/maps"
)

const (
	sPDXLicensesURL = "https://raw.githubusercontent.com/spdx/license-list-data/master/json/licenses.json"

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
	_, err := os.Stat(homedir)
	if err != nil {
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
	return os.WriteFile(path, []byte(json), 0644)
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

func Marshal(path string, m *map[string]any, args ...any) error {
	var p gen.Parser
	bytes, err := oj.Marshal(*m, args...)
	if err != nil {
		return err
	}

	node, err := p.Parse(bytes)
	if err != nil {
		return err
	}
	json := oj.JSON(node, &oj.Options{Sort: true, Indent: 2, OmitNil: true})
	return WriteJSON(path, json)
}

func LoadFile(filePath string) ([]byte, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	} else {
		bytes, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		return bytes, nil
	}
}

func DeleteFile(path string) error {
	return os.Remove(path)
}

func GetAndCacheLicenseFile(basedir string, overwrite bool) error {
	// ensure we have a home directory
	err := MkHomeDir(basedir)
	if err != nil {
		return err
	}

	licensesFilePath := basedir + sPDXLicensesFilePath

	_, err = os.Stat(licensesFilePath)
	if err != nil || overwrite {
		// file does not exist - download and store it
		spdxClient := http.Client{
			Timeout: time.Second * 2, // Timeout after 2 seconds
		}
		request, err := http.NewRequest(http.MethodGet, sPDXLicensesURL, nil)
		if err != nil {
			return err
		}
		request.Header.Set("User-Agent", "codemetagenerator")
		request.Header.Set("Accept", "application/json")

		response, getErr := spdxClient.Do(request)
		if getErr != nil {
			return getErr
		}
		defer response.Body.Close()

		// convert into reference keyed by licenseId => [licenseId] -> reference (url)
		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("unable to read response body when downloading SPDX license file: %s", err.Error())
		}
		var licensesList model.LicensesList
		oj.Unmarshal(bytes, &licensesList)

		var licensesMap map[string]any = make(map[string]any)
		lo.ForEach(licensesList.Licenses, func(license model.LicenseStruct, _ int) {
			licensesMap[license.LicenseId] = license.Reference
		})

		// marshal to file
		err = Marshal(licensesFilePath, &licensesMap)
		if err != nil {
			return fmt.Errorf("unable to save translated SPDX licenses file: %s", err.Error())
		}
	}
	return nil
}

func GetSupportedLicenses(basedir string) (*[]string, error) {
	err := GetAndCacheLicenseFile(basedir, false)
	if err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(basedir + sPDXLicensesFilePath)
	if err != nil {
		return nil, err
	}
	var licenses map[string]string
	oj.Unmarshal(bytes, &licenses)

	keys := maps.Keys(licenses)
	return &keys, nil
}
