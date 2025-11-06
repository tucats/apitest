package validator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type ValidationError struct {
	err      error
	context  string
	value    string
	expected string
}

var ErrEmptyTag = NewError("empty tag")
var ErrEmptyTagValue = NewError("empty tag value")
var ErrInvalidKeyword = NewError("invalid keyword")
var ErrInvalidInteger = NewError("invalid integer value")
var ErrMissingEnumValue = NewError("missing enum values")
var ErrInvalidEnumType = NewError("invalid field type for enum, must be string or int")
var ErrUnsupportedType = NewError("unsupported type")
var ErrValueOutOfRange = NewError("value out of range")
var ErrValueLengthOutOfRange = NewError("value length out of range")
var ErrInvalidData = NewError("invalid data")
var ErrUnimplemented = NewError("unimplemented type")
var ErrRequired = NewError("required field missing")
var ErrMaxDepthExceeded = NewError("maximum validation depth exceeded")
var ErrArrayLengthOutOfRange = NewError("array length out of range")
var ErrInvalidEnumeratedValue = NewError("invalid enumerated value")
var ErrUndefinedStructure = NewError("undefined structure")

func NewError(msg string) *ValidationError {
	return &ValidationError{
		err: errors.New(msg),
	}
}

func (e *ValidationError) Context(context string) *ValidationError {
	if e == nil {
		return nil
	}

	e2 := e.Clone()
	e2.context = context

	return e2
}

func (e *ValidationError) Value(value any) *ValidationError {
	if e == nil {
		return nil
	}

	if value == nil {
		return e
	}

	e2 := e.Clone()
	e2.value = fmt.Sprintf("%v", value)

	return e2
}

func (e *ValidationError) Expected(expected ...any) *ValidationError {
	if e == nil {
		return nil
	}

	e2 := e.Clone()

	list := make([]string, 0, len(expected))

	for _, v := range expected {
		switch v := v.(type) {
		case []string:
			for _, s := range v {
				list = append(list, strconv.Quote(s))
			}

			continue

		case []int:
			for _, n := range v {
				list = append(list, fmt.Sprintf("%d", n))
			}

		default:
			list = append(list, fmt.Sprintf("%v", v))
		}
	}

	e2.expected = strings.Join(list, ", ")
	if len(list) > 1 {
		e2.expected = "one of " + e2.expected
	}

	return e2
}

func (e *ValidationError) Clone() *ValidationError {
	if e == nil {
		return nil
	}

	return &ValidationError{
		err:     e.err,
		context: e.context,
		value:   e.value,
	}
}

func (e *ValidationError) Error() string {
	if e == nil {
		return "Success"
	}

	result := e.err.Error()
	if e.context != "" {
		result += ", in " + e.context
	}

	if e.value != "" {
		result += fmt.Sprintf(": %s", strconv.Quote(e.value))
	}

	if e.expected != "" {
		result += fmt.Sprintf(", expected %s", e.expected)
	}

	return result
}
