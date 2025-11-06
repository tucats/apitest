package validator

import (
	"encoding/json"
	"reflect"
	"strings"
)

func ValidateByName(name string, text string) error {
	dictionaryLock.Lock()
	item, exists := dictionary[name]
	dictionaryLock.Unlock()

	if !exists {
		return ErrUndefinedStructure.Context(name)
	}

	return item.Validate(text)
}

func (i *Item) Validate(text string) error {
	var (
		err error
		v   any
	)

	// Parse the JSON into an abstract object.
	err = json.Unmarshal([]byte(text), &v)
	if err != nil {
		return err
	}

	return i.validateValue(v, 0)
}

func (i *Item) validateValue(v any, depth int) error {
	if depth > maxValidationDepth {
		return ErrMaxDepthExceeded.Value(depth)
	}

	// If this is an array, validate each element.
	if i.IsArray || i.ValueType == TypeArray {
		array, ok := v.([]any)
		if !ok {
			return ErrInvalidData.Context(i.Name).Value(v)
		}

		if i.HasMinLength && len(array) < i.MinLength {
			return ErrArrayLengthOutOfRange.Context(i.Name).Value(len(array)).Expected(i.MinLength)
		}

		if i.HasMaxLength && len(array) > i.MaxLength {
			return ErrArrayLengthOutOfRange.Context(i.Name).Value(len(array)).Expected(i.MaxLength)
		}

		for _, element := range array {
			err := i.BaseType.validateValue(element, depth+1)
			if err != nil {
				return err
			}
		}

		return nil
	}

	switch i.ValueType {
	case TypeAny:
		return nil // Accept anything.

	case TypeMap:
		// For maps, we can only validate keys, not values. IF there is an enum list, let's check it out.
		if len(i.Enums) > 0 {
			actual := reflect.ValueOf(v)
			keys := actual.MapKeys()

			for _, key := range keys {
				found := false
				keyString := key.String()

				for _, enum := range i.Enums {
					if i.CaseSensitive {
						if keyString == enum {
							found = true

							break
						}
					} else {
						if strings.EqualFold(keyString, enum) {
							found = true

							break
						}
					}
				}

				if !found {
					return ErrInvalidEnumeratedValue.Context(i.Name).Value(keyString).Expected(i.Enums)
				}
			}
		}

	case TypeTime:
		_, err := getTimeValue(v)
		if err != nil {
			return ErrInvalidData.Context(i.Name).Value(v)
		}

		if i.HasMinValue {
			t, _ := getTimeValue(i.MinValue)
			vv, _ := getTimeValue(v)

			if vv.Before(t) {
				return ErrValueOutOfRange.Context(i.Name).Value(v)
			}
		}

		if i.HasMaxValue {
			t, _ := getTimeValue(i.MaxValue)
			vv, _ := getTimeValue(v)

			if vv.After(t) {
				return ErrValueOutOfRange.Context(i.Name).Value(v)
			}
		}

	case TypeUUID:
		_, err := getUUIDValue(v)
		if err != nil {
			return ErrInvalidData.Context(i.Name).Value(v)
		}

	case TypeBool:
		_, err := getBoolValue(v)
		if err != nil {
			return ErrInvalidData.Context(i.Name).Value(v)
		}

	case TypeInt:
		value, err := getIntValue(v)
		if err != nil {
			return ErrInvalidData.Context(i.Name).Value(v)
		}

		if i.HasMinValue {
			t, _ := getIntValue(i.MinValue)
			if value < t {
				return ErrValueOutOfRange.Context(i.Name).Value(value)
			}
		}

		if i.HasMaxValue {
			t, _ := getIntValue(i.MaxValue)
			if value > t {
				return ErrValueOutOfRange.Context(i.Name).Value(value)
			}
		}

		found := false

		for _, enum := range i.Enums {
			enumValue, _ := getIntValue(enum)

			if value == enumValue {
				found = true

				break
			}
		}

		if len(i.Enums) > 0 && !found {
			return ErrInvalidEnumeratedValue.Context(i.Name).Value(value).Expected(i.Enums)
		}

	case TypeFloat:
		value, err := getFloatValue(v)
		if err != nil {
			return ErrInvalidData.Context(i.Name).Value(v)
		}

		if i.HasMinValue {
			t, _ := getFloatValue(i.MinValue)
			if value < t {
				return ErrValueOutOfRange.Context(i.Name).Value(value)
			}
		}

		if i.HasMaxValue {
			t, _ := getFloatValue(i.MinValue)
			if value > t {
				return ErrValueOutOfRange.Context(i.Name).Value(value)
			}
		}

	case TypeString:
		value, err := getStringValue(v)
		if err != nil {
			return ErrInvalidData.Context(i.Name).Value(v)
		}

		if i.HasMinLength {
			if len(value) < i.MinLength {
				return ErrValueLengthOutOfRange.Context(i.Name).Value(value)
			}
		}

		if i.HasMaxLength {
			if len(value) > i.MaxLength {
				return ErrValueLengthOutOfRange.Context(i.Name).Value(value)
			}
		}

		found := false

		for _, enum := range i.Enums {
			if i.CaseSensitive {
				if value == enum {
					found = true

					break
				}
			} else {
				if strings.EqualFold(value, enum) {
					found = true

					break
				}
			}
		}

		if len(i.Enums) > 0 && !found {
			return ErrInvalidEnumeratedValue.Context(i.Name).Value(value).Expected(i.Enums)
		}

	case TypeStruct:
		m, ok := v.(map[string]any)
		if !ok {
			return ErrInvalidData.Context(i.Name).Value(v)
		}

		for _, field := range i.Fields {
			fieldValue, exists := m[field.Name]
			if !exists {
				if field.Required {
					return ErrRequired.Value(field.Name)
				}

				continue
			}

			err := field.validateValue(fieldValue, depth+1)
			if err != nil {
				return err
			}
		}

	default:
		return ErrUnimplemented.Context(i.Name).Value(i.ValueType.String())
	}

	return nil
}
