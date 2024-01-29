package internal

import (
	"testing"

	"encoding/json"
)

func TestUpdateMapValue(t *testing.T) {
	// Initialize a testMap
	var testMap map[string]any
	text := `{
		"key1": [
			{"key2": 2},
			{"key4": 4}
		],
		"key3": ["one", "two"],
		"key5": "hello",
		"key6": {
			"key7": "seven",
			"key8": 8,
			"key9": {
				"key10": [1, 2, 3],
				"key11": {
					"key12": "twelve"
				},
				"key13": "world"
			}
		}
	}`
	jsonErr := json.Unmarshal([]byte(text), &testMap)
	if jsonErr != nil {
		t.Errorf("Marshal of JSON text to object returned an error: %v", jsonErr)
	}

	// update nested key
	err1 := UpdateMapValue(testMap, "key1[0].key2", 20)
	if err1 != nil {
		t.Errorf("UpdateMapValue returned an error: %v", err1)
	}
	// update nested key again
	err2 := UpdateMapValue(testMap, "key1[0].key2", 30)
	if err2 != nil {
		t.Errorf("UpdateMapValue returned an error: %v", err2)
	}
	// update array key
	err3 := UpdateMapValue(testMap, "key3[0]", "four")
	if err3 != nil {
		t.Errorf("UpdateMapValue returned an error: %v", err3)
	}
	err4 := UpdateMapValue(testMap, "key6.key9.key11.key12", "twenty-twelve")
	if err4 != nil {
		t.Errorf("UpdateMapValue returned an error: %v", err4)
	}
	err5 := UpdateMapValue(testMap, "key5", "goodbye")
	if err5 != nil {
		t.Errorf("UpdateMapValue returned an error: %v", err5)
	}
	// update leaf array index
	err6 := UpdateMapValue(testMap, "key6.key9.key10[2]", 30)
	if err6 != nil {
		t.Errorf("UpdateMapValue returned an error: %v", err6)
	}
	// try to "update" a non-existing key
	err7 := UpdateMapValue(testMap, "key4", 4)
	if err7 == nil {
		t.Errorf("UpdateMapValue did not return an error when it should have")
	}
	// try to update to a nil value
	err8 := UpdateMapValue(testMap, "key5", nil)
	if err8 == nil {
		t.Errorf("UpdateMapValue did not return an error when it should have")
	}

	// Check the state of the testMap
	expected := `{"key1":[{"key2":30},{"key4":4}],"key3":["four","two"],"key5":"goodbye","key6":{"key7":"seven","key8":8,"key9":{"key10":[1,2,30],"key11":{"key12":"twenty-twelve"},"key13":"world"}}}`

	// Convert the testMap to a string
	b, err := json.Marshal(testMap)
	if err != nil {
		t.Errorf("Marshal returned an error: %v", err)
	}
	actual := string(b)
	if actual != expected {
		t.Errorf("Final map state comparison failed, got: %v, want: %v.", actual, expected)
	}
}

func TestInsertMapValue(t *testing.T) {
	// Initialize a testMap
	var testMap map[string]any
	text := `{
		"key1": [
			{"key2": 2},
			{"key4": 4}
		],
		"key3": ["one", "two"],
		"key5": "hello",
		"key6": {
			"key7": "seven",
			"key8": 8,
			"key9": {
				"key10": [1, 2, 3],
				"key11": {
					"key12": "twelve"
				},
				"key13": "world"
			}
		}
	}`
	jsonErr := json.Unmarshal([]byte(text), &testMap)
	if jsonErr != nil {
		t.Errorf("Marshal of JSON text to object returned an error: %v", jsonErr)
	}

	err1 := InsertMapValue(testMap, "key1[0].key3", 3)
	if err1 != nil {
		t.Errorf("InsertMapValue returned an error: %v", err1)
	}
	err2 := InsertMapValue(testMap, "key1[1].key14", 14)
	if err2 != nil {
		t.Errorf("InsertMapValue returned an error: %v", err2)
	}
	// try to insert a nil value
	err3 := InsertMapValue(testMap, "key1[1].key14", nil)
	if err3 == nil {
		t.Errorf("InsertMapValue did not return an error when it should have")
	}
	err4 := InsertMapValue(testMap, "key6.key9.key11.key23", "twenty-three")
	if err4 != nil {
		t.Errorf("InsertMapValue returned an error: %v", err4)
	}
	// try to insert into an array -- out of bounds
	err5 := InsertMapValue(testMap, "key3[2]", "three")
	if err5 == nil {
		t.Errorf("InsertMapValue did not return an error when it should have")
	}
	// try to insert into an array -- existing index
	err6 := InsertMapValue(testMap, "key3[0]", "three")
	if err6 == nil {
		t.Errorf("InsertMapValue did not return an error when it should have")
	}
	// try to insert over an existing key
	err7 := InsertMapValue(testMap, "key5", "goodbye")
	if err7 == nil {
		t.Errorf("InsertMapValue did not return an error when it should have")
	}
	// insert to leaf array index
	err8 := InsertMapValue(testMap, "key6.key9.key10[2]", 30)
	if err8 == nil {
		t.Errorf("InsertMapValue did not return an error when it should have")
	}

	// Check the state of the testMap
	expected := `{"key1":[{"key2":2,"key3":3},{"key14":14,"key4":4}],"key3":["one","two"],"key5":"hello","key6":{"key7":"seven","key8":8,"key9":{"key10":[1,2,3],"key11":{"key12":"twelve","key23":"twenty-three"},"key13":"world"}}}`
	// Convert the testMap to a string
	b, err := json.Marshal(testMap)
	if err != nil {
		t.Errorf("Marshal returned an error: %v", err)
	}
	actual := string(b)
	if actual != expected {
		t.Errorf("Final map state comparison failed, got: %v, want: %v.", actual, expected)
	}
}

func TestRemoveMapValue(t *testing.T) {
	// Initialize a testMap
	var testMap map[string]any
	text := `{
		"key1": [
			{"key2": 2},
			{"key4": 4}
		],
		"key3": ["one", "two"],
		"key5": "hello",
		"key6": {
			"key7": "seven",
			"key8": 8,
			"key9": {
				"key10": [1, 2, 3],
				"key11": {
					"key12": "twelve"
				},
				"key13": "world"
			}
		}
	}`
	jsonErr := json.Unmarshal([]byte(text), &testMap)
	if jsonErr != nil {
		t.Errorf("Marshal of JSON text to object returned an error: %v", jsonErr)
	}

	err1 := RemoveMapValue(testMap, "key1[0].key2") // removes "key2": 2 -> [{}, {"key4": 4}]
	if err1 != nil {
		t.Errorf("RemoveMapValue returned an error: %v", err1)
	}
	err2 := RemoveMapValue(testMap, "key1[0]") // removes {} -> [{"key4": 4}]
	if err2 != nil {
		t.Errorf("RemoveMapValue returned an error: %v", err2)
	}
	err3 := RemoveMapValue(testMap, "key1[0]") // removes "key4": 4 -> []
	if err3 != nil {
		t.Errorf("RemoveMapValue returned an error: %v", err3)
	}
	err4 := RemoveMapValue(testMap, "key1[0]") // index no longer exists -- should be index out of bounds
	if err4 == nil {
		t.Errorf("RemoveMapValue did not return an error when it should have")
	}
	err5 := RemoveMapValue(testMap, "key6.key9.key10[3]") // nested index does not exists -- should be index out of bounds
	if err5 == nil {
		t.Errorf("RemoveMapValue did not return an error when it should have")
	}
	err6 := RemoveMapValue(testMap, "key16") // key doesn't exist
	if err6 == nil {
		t.Errorf("RemoveMapValue did not return an error when it should have")
	}
	err7 := RemoveMapValue(testMap, "key6.key9.key13") // remove leaf
	if err7 != nil {
		t.Errorf("RemoveMapValue returned an error: %v", err7)
	}
	err8 := RemoveMapValue(testMap, "key6.key9.key10")
	if err8 != nil {
		t.Errorf("RemoveMapValue returned an error: %v", err8)
	}

	// Check the state of the testMap
	expected := `{"key1":[],"key3":["one","two"],"key5":"hello","key6":{"key7":"seven","key8":8,"key9":{"key11":{"key12":"twelve"}}}}`
	// Convert the testMap to a string
	b, err := json.Marshal(testMap)
	if err != nil {
		t.Errorf("Marshal returned an error: %v", err)
	}
	actual := string(b)
	if actual != expected {
		t.Errorf("Final map state comparison failed, got: %v, want: %v.", actual, expected)
	}
}
