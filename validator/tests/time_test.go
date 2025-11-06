package tests

import (
	"testing"
	"time"

	"github.com/tucats/validator"
)

// Structure to validate.
type TimeObject struct {
	When time.Time `json:"when" validate:"required,min=2000-01-01"`
	Name string    `json:"name" validate:"required,minlength=1,maxlength=100"`
}

func Test_Time(t *testing.T) {
	type TestITem struct {
		name     string
		object   any
		jsonText string
		expected error
	}

	tests := []TestITem{
		{
			"Valid time",
			&TimeObject{},
			`{
			    "when": "Dec 15, 2023 10:00AM",
				"name": "First"
			}`,
			nil,
		},
		{
			"invalid time format",
			&TimeObject{},
			`{
			    "when": "Yesterday",
				"name": "First"
			}`,
			validator.ErrInvalidData.Context("when").Value("Yesterday"),
		},
		{
			"date too early",
			&TimeObject{},
			`{
			    "when": "July 20, 1969 08:18AM",
				"name": "First"
			}`,
			validator.ErrValueOutOfRange.Context("when").Value("July 20, 1969 08:18AM"),
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
