package tests

import (
	"testing"

	"github.com/tucats/validator"
)

// Structure to validate.
type MapObject struct {
	Items map[string]string `json:"items" validate:"required,enum=key1|key2"`
}

type MapObject2 struct {
	Items map[string]int `json:"items" validate:"required,enum=key1|key2"`
}

type MapObject3 struct {
	Items map[string][]string `json:"items" validate:"required,enum=key1|key2"`
}

func Test_Maps(t *testing.T) {
	type TestITem struct {
		name     string
		object   any
		jsonText string
		expected error
	}

	tests := []TestITem{
		{
			"Valid map[string]string",
			&MapObject{},
			`{
			    "items": {
				    "key1": "value1",
                    "key2": "value2"
				}
			}`,
			nil,
		},
		{
			"Valid map[string]int",
			&MapObject2{},
			`{
			    "items": {
				    "key1": 55,
                    "key2": 67
				}
			}`,
			nil,
		},
		{
			"Valid map[string][]string]",
			&MapObject2{},
			`{
			    "items": {
				    "key1": ["value1","value2"],
                    "key2": ["value3","value4"]
				}
			}`,
			nil,
		},
		{
			"invalid map key value",
			&MapObject2{},
			`{
			    "items": {
				    "key3": 55,
                    "key2": 67
				}
			}`,
			validator.ErrInvalidEnumeratedValue.Context("items").Value("key3").Expected([]string{"key1", "key2"}),
		},
		{
			"invalid map, missing required field",
			&MapObject{},
			`{
			}`,
			validator.ErrRequired.Value("items"),
		},
	}

	for _, test := range tests {
		var (
			msg1 string
			msg2 string
		)

		// Define the structures we want to validate.
		item, err := validator.New(test.object)
		if err != nil {
			t.Fatal("Failed to define structure:", err)
		}

		err = item.Validate(test.jsonText)
		if err != nil {
			msg1 = err.Error()
		}

		if test.expected != nil {
			msg2 = test.expected.Error()
		}

		expected := (msg1 == msg2)

		if !expected {
			t.Fatalf("In \"%s\", unexpected result: %v\n", test.name, err)
		}
	}
}
