package tests

import (
	"testing"

	"github.com/tucats/validator"
)

// Structure to validate.
type AnyObject struct {
	ID   any    `json:"id"   validate:"required"`
	Name string `json:"name" validate:"required,minlength=1,maxlength=100"`
}

func Test_Any(t *testing.T) {
	type TestITem struct {
		name     string
		object   any
		jsonText string
		expected error
	}

	tests := []TestITem{
		{
			"Any object using string",
			&AnyObject{},
			`{
			    "id": "test value",
				"name": "First"
			}`,
			nil,
		},
		{
			"Any object using integer",
			&AnyObject{},
			`{
			    "id": 55,
				"name": "First"
			}`,
			nil,
		},
		{
			"Any object using object",
			&AnyObject{},
			`{
			    "id": {
					"sub_id": 12345
				},
				"name": "First"
			}`,
			nil,
		},
		{
			"Any object, missing field",
			&UUIDObject{},
			`{
				"name": "First"
			}`,
			validator.ErrRequired.Value("id"),
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
