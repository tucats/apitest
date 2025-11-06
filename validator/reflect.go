package validator

import (
	"reflect"
	"strings"
)

// Maximum allowed recursion depth for validation.
const maxValidationDepth = 10

func New(v any) (*Item, error) {
	return defineItem(v, 0)
}

func defineItem(v any, depth int) (*Item, error) {
	var err error

	// IF we exceed maximum recursion depth, return an error
	if depth > maxValidationDepth {
		return nil, ErrMaxDepthExceeded.Value(depth)
	}

	// Based on the reflected type of this item, form an Item{} structure that defines it.
	item := &Item{}

	valueType := reflect.TypeOf(v)
	if valueType == nil {
		item.ValueType = TypeAny

		return item, nil
	}

	kind := valueType.Kind()

	// Handle well-known external types first
	typeName := valueType.String()
	switch typeName {
	case "uuid.UUID":
		item.ValueType = TypeUUID

		return item, nil

	case "time.Time":
		item.ValueType = TypeTime

		return item, nil
	}

	// Handle based on the kind of the reflected type
	switch kind {
	case reflect.Interface:
		item.ValueType = TypeAny

	case reflect.Pointer:
		// Dereference the pointer and create an item for the base type
		valueType = valueType.Elem()
		v = reflect.Zero(valueType).Interface()

		item, err = New(v)
		if err != nil {
			return nil, err
		}

		item.IsPointer = true

	case reflect.Map:
		item.ValueType = TypeMap

	case reflect.String:
		item.ValueType = TypeString

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		item.ValueType = TypeInt

	case reflect.Float32, reflect.Float64:
		item.ValueType = TypeFloat

	case reflect.Bool:
		item.ValueType = TypeBool

	case reflect.Array, reflect.Slice:
		// Create an item for the base type of the array/slice
		item.ValueType = TypeArray

		baseItem, err := New(reflect.Zero(valueType.Elem()).Interface())
		if err != nil {
			return nil, err
		}

		item.BaseType = baseItem

	case reflect.Struct:
		// Iterate over the fields of the struct and build items for each field
		item.ValueType = TypeStruct

		for i := 0; i < valueType.NumField(); i++ {
			field := valueType.Field(i)

			fieldItem, err := defineItem(reflect.Zero(field.Type).Interface(), depth+1)
			if err != nil {
				return nil, err
			}

			fieldItem.Name = field.Name

			// See if there is a JSON tag to get the field name from
			jsonTag := field.Tag.Get("json")
			if len(jsonTag) > 0 {
				jsonParts := strings.SplitN(jsonTag, ",", 2)
				if len(jsonParts) > 0 && jsonParts[0] != "" {
					fieldItem.Name = jsonParts[0]
				}
			}

			// Parse the field's validate tag if present and build an item for it
			tagString := field.Tag.Get(validateTagName)
			if len(strings.TrimSpace(tagString)) > 0 {
				err = fieldItem.ParseTag(tagString)
				if err != nil {
					return nil, err
				}
			}

			item.Fields = append(item.Fields, *fieldItem)
		}

	default:
		err = ErrUnsupportedType.Value(valueType.String())
	}

	return item, err
}
