package tests

import (
	"testing"

	"github.com/tucats/validator"
)

// Structures to validate.

type Address struct {
	Street string `json:"street" validate:"required,minlength=1,maxlength=100"`
	City   string `json:"city"   validate:"required,minlength=1,maxlength=100"`
}

type Person struct {
	Name    string  `json:"name"    validate:"required,minlength=1,maxlength=100"`
	Age     int     `json:"age"     validate:"required,min=18,max=65"`
	Address Address `json:"address" validate:"required"`
}

type Employees struct {
	Department string   `json:"department" validate:"required"`
	Division   string   `json:"division"   validate:"required,enum=HR|Finance|Marketing|Engineering"`
	Staff      []Person `json:"staff"      validate:"minlen=1"`
}

func Test_AddressStruct(t *testing.T) {
	type TestITem struct {
		name     string
		object   any
		jsonText string
		expected error
	}

	tests := []TestITem{
		{
			"valid Employees",
			&Employees{},
			`{
			    "department": "Space Research",
				"division": "Engineering",
			    "staff": [
					{
						"name": "John Doe",
						"age": 35,
						"address":{
							"street": "123 Main St",
							"city": "New York"
						}
					},
					{
						"name": "Sue Smith",
						"age": 52,
						"address":{
							"street": "155 Oak Ave",
							"city": "New York"
						}
					}
			    ]
            }`,
			nil,
		},
		{
			"invalid Employees, bad division enum value",
			&Employees{},
			`{
			    "department": "Space Research",
				"division": "Science",
			    "staff": [
					{
						"name": "John Doe",
						"age": 35,
						"address":{
							"street": "123 Main St",
							"city": "New York"
						}
					},
					{
						"name": "Sue Smith",
						"age": 52,
						"address":{
							"street": "155 Oak Ave",
							"city": "New York"
						}
					}
			    ]
            }`,
			validator.ErrInvalidEnumeratedValue.Context("division").Value("Science").Expected([]string{"HR", "Finance", "Marketing", "Engineering"}),
		},
		{
			"invalid Employees, age out of range",
			&Employees{},
			`{
			    "department": "Space Research",
				"division": "Engineering",
			    "staff": [
					{
						"name": "John Doe",
						"age": 75,
						"address":{
							"street": "123 Main St",
							"city": "New York"
						}
					},
					{
						"name": "Sue Smith",
						"age": 52,
						"address":{
							"street": "155 Oak Ave",
							"city": "New York"
						}
					}
			    ]
            }`,
			validator.ErrValueOutOfRange.Context("age").Value(75),
		},
		{
			"invalid Employees, empty staff array",
			&Employees{},
			`{
			    "department": "Space Research",
				"division": "Engineering",
			    "staff": []
            }`,
			validator.ErrArrayLengthOutOfRange.Context("staff").Value(0).Expected(1),
		},
		{
			"valid JSON for address",
			&Address{},
			`{
				"street": "123 Main St",
				"city": "New York"
			}`,
			nil,
		},
		{
			"street string is too short",
			&Address{},
			`{
				"street": "",
				"city": "New York"
			}`,
			validator.ErrValueLengthOutOfRange.Context("street"),
		},
		{
			"city field not present",
			&Address{},
			`{
				"street": "123 Main St"
			}`,
			validator.ErrRequired.Value("city"),
		},
		{
			"valid Person",
			&Person{},
			`{
				"name": "John Doe",
				"age": 35,
				"address":{
					"street": "123 Main St",
					"city": "New York"
				}
			}`,
			nil,
		},
		{
			"invalid Person, age out of range",
			&Person{},
			`{
				"name": "John Doe",
				"age": 15,
				"address":{
					"street": "123 Main St",
					"city": "New York"
				}
			}`,
			validator.ErrValueOutOfRange.Context("age").Value(15),
		},
		{
			"invalid Person, missing city field",
			&Person{},
			`{
				"name": "John Doe",
				"age": 42,
				"address":{
					"street": "123 Main St"
				}
			}`,
			validator.ErrRequired.Value("city"),
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
