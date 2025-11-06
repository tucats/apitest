package tests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/tucats/validator"
)

// Structure to validate.
type UUIDObject struct {
	ID   uuid.UUID `json:"id"   validate:"required"`
	Name string    `json:"name" validate:"required,minlength=1,maxlength=100"`
}

func Test_UUID(t *testing.T) {
	type TestITem struct {
		name     string
		object   any
		jsonText string
		expected error
	}

	tests := []TestITem{
		{
			"Valid UUID",
			&UUIDObject{},
			`{
			    "id": "29c80af7-b490-497c-85a4-a3df8233c051",
				"name": "First"
			}`,
			nil,
		},
		{
			"invalid UUID",
			&UUIDObject{},
			`{
			    "id": "29c80af7-b490-497c-85a4-a3df8c051",
				"name": "First"
			}`,
			validator.ErrInvalidData.Context("id").Value("29c80af7-b490-497c-85a4-a3df8c051"),
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
