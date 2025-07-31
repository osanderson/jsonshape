package jsonshape

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
)

func TestSanitize_PreservesShape(t *testing.T) {
	input := "{" +
		"\"id\": \"01985dcc-875c-7178-b54e-9a52ca230d18\", " +
		"\"name\": \"Alice\", " +
		"\"email\": \"alice@example.com\"," +
		"\"age\": 35.4," +
		"\"ageAdjustment\": -2," +
		"\"isActive\": true," +
		"\"phones\": [\"+6598765432\", \"+123456789\"]," +
		"\"profile\": {" +
		"\"id\": \"ABC:123_XY$Z\"," +
		"\"bio\": \"This is a long text\"}," +
		"\"updated\": \"2025-07-31T08:11:04.123Z\"" +
		"}"

	log.Printf("%s", input)

	var data interface{}
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		t.Fatal(err)
	}

	log.Printf("%x", data)

	opts := Options{
		Matchers: DefaultMatchers,
	}

	sanitized, err := Sanitize(data, opts)
	if err != nil {
		t.Fatalf("Sanitize error: %s", err)
	}

	jsonOut, err := json.Marshal(sanitized)
	if err != nil {
		t.Fatalf("Failed to marshal sanitized: %v", err)
	}

	var reParsed interface{}
	if err := json.Unmarshal(jsonOut, &reParsed); err != nil {
		t.Fatalf("Sanitized JSON is invalid: %v", err)
	}

	if !shapesMatch(data, reParsed) {
		t.Error("Sanitized JSON does not match original shape")
	}
}

func shapesMatch(a, b interface{}) bool {
	switch aVal := a.(type) {
	case map[string]interface{}:
		bVal, ok := b.(map[string]interface{})
		if !ok || len(aVal) != len(bVal) {
			return false
		}
		for k := range aVal {
			if !shapesMatch(aVal[k], bVal[k]) {
				return false
			}
		}
		return true
	case []interface{}:
		bVal, ok := b.([]interface{})
		if !ok || len(aVal) != len(bVal) {
			return false
		}
		for i := range aVal {
			if !shapesMatch(aVal[i], bVal[i]) {
				return false
			}
		}
		return true
	default:
		return reflect.TypeOf(a) == reflect.TypeOf(b)
	}
}
