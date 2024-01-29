package internal

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/nexidian/gocliselect"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

const (
	SPDXLicensesURL = "https://raw.githubusercontent.com/spdx/license-list-data/master/json/licenses.json"

	CodemetaGeneratorDirectoryName = ".codemetagenerator"
	InProgressFilePath             = "/" + CodemetaGeneratorDirectoryName + "/codemeta.inprogress.json"
	SPDXLicensesFilePath           = "/" + CodemetaGeneratorDirectoryName + "/spdx-licenses.json"
)

var UserHomeDir string = getUserHomeDir()
var SPDXLicensesFile string = UserHomeDir + SPDXLicensesFilePath

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func MkPrompt(text string) (*string, error) {
	fmt.Println(text)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		return nil, errors.New("Error: unable to read input from the console for prompt: " + text)
	}
	selection := scanner.Text()
	return &selection, nil
}

func LoadInProgressCodeMetaFile() (*map[string]any, error) {
	homeDir := UserHomeDir
	_, error := os.Stat(homeDir + InProgressFilePath)
	if error != nil {
		return nil, error
	} else {
		bytes, err := os.ReadFile(homeDir + InProgressFilePath)
		if err != nil {
			return nil, err
		}
		var codemeta map[string]any
		json.Unmarshal(bytes, &codemeta)
		return &codemeta, nil
	}
}

func SaveInProgressCodeMetaFile(codemeta *map[string]interface{}) error {
	homeDir := UserHomeDir
	file, err := json.MarshalIndent(codemeta, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(homeDir+InProgressFilePath, file, 0644)
	return err
}

func DeleteInProgressCodeMetaFile() error {
	homeDir := UserHomeDir
	err := os.Remove(homeDir + InProgressFilePath)
	return err
}

func NewPersonOrOrganizationPrompt(key string) (*map[string]string, error) {
	keyTypeMenu := gocliselect.NewMenu("Please enter a " + key + " type")
	keyTypeMenu.AddItem("Person", "person")
	keyTypeMenu.AddItem("Organization", "organization")
	keyType := keyTypeMenu.Display()

	var thing *map[string]string // Person or Organization map
	switch keyType {
	case "person":
		givenName, err := MkPrompt("Enter the given (first) name of the person:")
		if err != nil {
			return nil, err
		}
		familyName, err := MkPrompt("Enter the family (last) name of the person:")
		if err != nil {
			return nil, err
		}
		email, err := MkPrompt("Enter the email address of the person:")
		if err != nil {
			return nil, err
		}
		id, err := MkPrompt("Enter the identifier of the person (see: https://orcid.org):")
		if err != nil {
			return nil, err
		}
		thing = NewPerson(givenName, familyName, email, id)
	case "organization":
		name, err := MkPrompt("Enter the name of the organization:")
		if err != nil {
			return nil, err
		}
		url, err := MkPrompt("Enter the URL of the organization:")
		if err != nil {
			return nil, err
		}
		id, err := MkPrompt("Enter the identifier of the organization:")
		if err != nil {
			return nil, err
		}
		thing = NewOrganization(name, url, id)
	}
	return thing, nil
}

// cannot 'overwrite' in the insert case, has to strictly be a new key
func InsertMapValue(m map[string]any, key string, value any) error {
	if value == nil {
		return errors.New("cannot insert a nil value")
	}
	// cannot overwrite an existing key -- must not already exist
	keys := strings.Split(key, ".")
	if len(keys) == 1 {
		indexedKey, index, _ := getIndexedKey(key)
		// cannot be a indexed key, e.g, foo[1] because can't insert into an array, only an object witin an array
		if index != nil {
			return errors.New("cannot insert a new key at an array index: " + key)
		}
		_, ok := m[*indexedKey]
		// single non-indexed key: foo, nil, nil --> m["foo"]
		if ok {
			return errors.New("key already exists: " + key)
		} else {
			// insert value at the key
			m[key] = value
		}
		return nil
	} else {
		// traverse the map using keys
		for i := 0; i < len(keys)-1; i++ {
			currentKey := keys[i]                             // foo[1].bar -> foo[1]
			indexedKey, index, _ := getIndexedKey(currentKey) // foo[1] -> foo, 1
			nextValue := m[*indexedKey]                       // m["foo"] -> []map[string]interface{}{...} or map[string]interface{}{...}
			if index != nil {
				// there must be at least one more element in the keys array which represents the key of the map to insert
				if i == len(keys)-1 {
					return errors.New("cannot insert a new key at an array index: " + currentKey)
				}
				// array -- check that index is within bounds of the array value at the indexed
				if *index >= len(nextValue.([]any)) {
					return errors.New("index of key is out of bounds: " + *indexedKey)
				}
				// indexed key --> m["foo"][index], check if value exists at index (TODO: is this necessary?)
				if nextValue.([]any)[*index] == nil {
					return errors.New("key does not exist: " + *indexedKey)
				} else {
					// move to the map at the index
					m = nextValue.([]any)[*index].(map[string]any)
				}
			} else {
				// map -- move to the map at the key
				// map at the key
				m = nextValue.(map[string]any)
			}
		}

		// insert the value at the final key, value cannot already exist and cannot be an array index
		finalKey := keys[len(keys)-1]
		indexedKey, index, _ := getIndexedKey(finalKey)
		_, ok := m[*indexedKey]
		if ok {
			return errors.New("key already exist: " + finalKey)
		}
		if index != nil {
			return errors.New("cannot insert a new key at an array index: " + finalKey)
		}

		m[*indexedKey] = value

		return nil
	}
}

func UpdateMapValue(m map[string]any, key string, newValue any) error {
	if newValue == nil {
		return errors.New("cannot update a key to a nil value")
	}
	// cannot insert new key -- must already exist
	keys := strings.Split(key, ".")
	if len(keys) == 1 {
		indexedKey, index, _ := getIndexedKey(key) // foo[1] -> foo, 1, nil
		currentValue, ok := m[*indexedKey]         // m["foo"] -> []map[string]interface{}{...} or map[string]interface{}{...}
		if !ok {
			return errors.New("key does not exist: " + *indexedKey)
		}
		if index != nil {
			// array m["foo"][index] -- check that index is within bounds of the array value at the indexed
			if *index >= len(currentValue.([]any)) {
				msg := fmt.Sprintf("index: %d of key is out of bounds: %s", index, *indexedKey)
				return errors.New(msg)
			}
			// indexed key --> m["foo"][index], check if value exists at index (TODO: is this necessary?)
			if currentValue.([]any)[*index] == nil {
				return errors.New("key does not exist: " + *indexedKey)
			} else {
				// update the value at the array index to the new value -- m["foo"][1] = newValue
				currentValue.([]any)[*index] = newValue
			}
			return nil
		} else {
			// non-indexed key: foo, nil, nil --> m["foo"]
			// update the map to the new value
			m[key] = newValue
			return nil
		}
	} else {
		// traverse the map using keys
		for i := 0; i < len(keys)-1; i++ {
			currentKey := keys[i]                             // foo[1].bar -> foo[1]
			indexedKey, index, _ := getIndexedKey(currentKey) // foo[1] -> foo, 1
			nextValue, ok := m[*indexedKey]                   // m["foo"] -> []map[string]interface{}{...} or map[string]interface{}{...}
			if !ok {
				return errors.New("key does not exist: " + *indexedKey)
			}
			if index != nil {
				// array -- check that index is within bounds of the array value at the indexed
				if *index >= len(nextValue.([]any)) {
					msg := fmt.Sprintf("index: %d of key is out of bounds: %s", index, *indexedKey)
					return errors.New(msg)
				}
				// indexed key --> m["foo"][index], check if value exists at index (TODO: is this necessary?)
				if nextValue.([]any)[*index] == nil {
					return errors.New("key does not exist: " + *indexedKey)
				} else {
					// move to the map at the index
					m = nextValue.([]any)[*index].(map[string]any)
				}
			} else {
				// non-indexed key: foo, nil, nil --> m["foo"]
				// map -- move to the map at the key
				m = nextValue.(map[string]any)
			}
		}

		// update the value at the final key, value must already exist and can be an array index
		finalKey := keys[len(keys)-1]
		indexedKey, index, _ := getIndexedKey(finalKey)
		_, ok := m[*indexedKey]
		if !ok {
			return errors.New("key does not exist: " + finalKey)
		}
		if index != nil {
			m[*indexedKey].([]any)[*index] = newValue
		} else {
			m[*indexedKey] = newValue
		}

		return nil
	}
}

func RemoveMapValue(m map[string]any, key string) error {
	keys := strings.Split(key, ".")
	if len(keys) == 1 {
		indexedKey, index, _ := getIndexedKey(key) // foo[1] -> foo, 1, nil
		currentValue, ok := m[*indexedKey]         // m["foo"] -> []map[string]interface{}{...} or map[string]interface{}{...}
		if !ok {
			return errors.New("key does not exist: " + *indexedKey)
		}
		if index != nil {
			// array m["foo"][index] -- check that index is within bounds of the array value at the indexed
			if *index >= len(currentValue.([]any)) {
				msg := fmt.Sprintf("index: %d of key is out of bounds: %s", index, *indexedKey)
				return errors.New(msg)
			}
			// indexed key --> m["foo"][index], check if value exists at index (TODO: is this necessary?)
			if currentValue.([]any)[*index] == nil {
				return errors.New("key does not exist: " + *indexedKey)
			} else {
				// delete value from slice at the array index -- m["foo"][1] and set it to the index in the map
				m[*indexedKey] = slices.Delete(currentValue.([]any), *index, *index+1)
			}
			return nil
		} else {
			// non-indexed key: foo, nil, nil --> m["foo"]
			// delete key from map
			delete(m, key)
			return nil
		}
	} else {
		// traverse the map using keys
		for i := 0; i < len(keys)-1; i++ {
			currentKey := keys[i]                             // foo[1].bar -> foo[1]
			indexedKey, index, _ := getIndexedKey(currentKey) // foo[1] -> foo, 1
			nextValue, ok := m[*indexedKey]                   // m["foo"] -> []map[string]interface{}{...} or map[string]interface{}{...}
			if !ok {
				return errors.New("key does not exist: " + *indexedKey)
			}
			if index != nil {
				// array -- check that index is within bounds of the array value at the indexed
				if *index >= len(nextValue.([]any)) {
					msg := fmt.Sprintf("index: %d of key is out of bounds: %s", index, *indexedKey)
					return errors.New(msg)
				}
				// indexed key --> m["foo"][index], check if value exists at index (TODO: is this necessary?)
				if nextValue.([]any)[*index] == nil {
					return errors.New("key does not exist: " + *indexedKey)
				} else {
					// move to the map at the index
					m = nextValue.([]any)[*index].(map[string]any)
				}
			} else {
				// non-indexed key: foo, nil, nil --> m["foo"]
				// map -- move to the map at the key
				m = nextValue.(map[string]any)
			}
		}

		// delete the value at the final key, value must already exist
		finalKey := keys[len(keys)-1]
		indexedKey, index, _ := getIndexedKey(finalKey)
		nextValue, ok := m[*indexedKey]
		if !ok {
			return errors.New("key does not exist: " + finalKey)
		}
		if index != nil {
			// array -- check that index is within bounds of the array value at the indexed
			if *index >= len(nextValue.([]any)) {
				msg := fmt.Sprintf("index: %d of key is out of bounds: %s", index, *indexedKey)
				return errors.New(msg)
			}
			m[*indexedKey] = slices.Delete(nextValue.([]any), *index, *index+1)
		} else {
			delete(m, *indexedKey)
		}

		return nil
	}
}

func GetAndCacheLicenseFile(overwrite bool) error {
	_, error := os.Stat(SPDXLicensesFile)
	if error != nil || overwrite {
		// file does not exist - download and store it
		spdxClient := http.Client{
			Timeout: time.Second * 2, // Timeout after 2 seconds
		}
		request, err := http.NewRequest(http.MethodGet, SPDXLicensesURL, nil)
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

		// covnert into reference keyed by licenseId
		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		var licensesList LicensesList
		json.Unmarshal(bytes, &licensesList)

		var licensesMap map[string]string = make(map[string]string)
		lo.ForEach(licensesList.Licenses, func(license LicenseStruct, _ int) {
			licensesMap[license.LicenseId] = license.Reference
		})

		json, err := json.MarshalIndent(licensesMap, "", " ")
		if err != nil {
			return err
		}
		// Write new to file
		writeErr := os.WriteFile(SPDXLicensesFile, json, 0644)
		if writeErr != nil {
			return writeErr
		}
	}
	return nil
}

func GetSupportedLicenses() (*[]string, error) {
	licenses, err := loadLicenseFile()
	if err != nil {
		return nil, err
	}

	keys := maps.Keys(*licenses)
	return &keys, nil
}

func CheckAndConvertLicenseId(s *string) (*string, error) {
	reference, err := getLicenseReferenceById(s)
	if err != nil {
		return nil, err
	}
	if reference == nil {
		return nil, errors.New("Invalid license ID: " + *s)
	}
	return reference, nil
}

// Private

func getUserHomeDir() string {
	homeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)
	return homeDir
}

func getIndexedKey(key string) (*string, *int, error) {
	if strings.Contains(key, "[") && strings.Contains(key, "]") {
		keyAndIndex := strings.Split(key, "[")
		indexStr := strings.TrimRight(keyAndIndex[1], "]")
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			return nil, nil, errors.New("Invalid index: " + indexStr)
		}
		return &keyAndIndex[0], &index, nil
	} else {
		return &key, nil, nil
	}
}

func loadLicenseFile() (*map[string]string, error) {
	err := GetAndCacheLicenseFile(false)
	if err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(SPDXLicensesFile)
	if err != nil {
		return nil, err
	}
	var licenses map[string]string
	json.Unmarshal(bytes, &licenses)

	return &licenses, nil
}

func getLicenseReferenceById(id *string) (*string, error) {
	licenses, err := loadLicenseFile()
	if err != nil {
		return nil, err
	}

	reference := (*licenses)[*id]
	return &reference, nil
}
