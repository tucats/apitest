package parser

import (
	"encoding/json"
	"fmt"
)

func GetOneItem(text string, item string) (string, error) {
	items, err := GetItem(text, item)
	if err == nil {
		if len(items) == 1 {
			return items[0], nil
		}

		if len(items) > 1 {
			return "", fmt.Errorf("Ambiguious expresssion (multiple values): %s", item)
		} else {
			return "", fmt.Errorf("No such item found: %s", item)
		}
	}

	return "", err
}

// For a given JSON payload string, extract a specific item from the payload. The item specification
// is a dot-notation string that can include integer indices and string map key values. The value is
// always returned as a string representation.
func GetItem(text string, item string) ([]string, error) {
	// Convert the body text to an arbitrary interface object using JSON
	var body interface{}

	if err := json.Unmarshal([]byte(text), &body); err != nil {
		return nil, err
	}

	return parse(body, item)
}
