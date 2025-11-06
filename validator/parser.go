package validator

import (
	"strconv"
	"strings"
)

func (item *Item) ParseTag(tag string) error {
	var err error

	// Split the tag string into parts by commas
	parts := strings.SplitSeq(tag, ",")

	// Scan over each part, and apply it to the item as appropriate
	for part := range parts {
		var (
			key   string
			value string
		)

		part = strings.TrimSpace(part)
		if part == "" {
			return ErrEmptyTag
		}

		elements := strings.SplitN(part, "=", 2)
		key = strings.ToLower(strings.TrimSpace(elements[0]))

		if len(elements) > 1 {
			value = strings.TrimSpace(elements[1])
			if len(value) == 0 {
				return ErrEmptyTagValue.Context(key)
			}
		}

		switch key {
		case "name":
			item.Name = value

		case "required":
			item.Required = true

		case "minlength", "minlen":
			item.HasMinLength = true

			n, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return ErrInvalidInteger.Context(key).Value(value)
			}

			item.MinLength = int(n)

		case "maxlength", "maxlen":
			item.HasMaxLength = true

			n, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return ErrInvalidInteger.Context(key).Value(value)
			}

			item.MaxLength = int(n)

		case "maxvalue", "max":
			item.HasMaxValue = true

			item.MaxValue = value

		case "minvalue", "min":
			item.HasMinValue = true

			item.MinValue = value

		case "enum":
			if item.ValueType != TypeString && item.ValueType != TypeInt && item.ValueType != TypeMap {
				return ErrInvalidEnumType.Context(key).Value(item.ValueType.String())
			}

			enums := strings.Split(value, "|")
			if len(enums) == 0 {
				return ErrMissingEnumValue.Context(key)
			}

			item.Enums = make([]string, len(enums))

			for i, enum := range enums {
				item.Enums[i] = strings.TrimSpace(enum)
			}

		case "matchcase", "casesensitive":
			item.CaseSensitive = true

		default:
			return ErrInvalidKeyword.Value(key)
		}
	}

	return err
}
