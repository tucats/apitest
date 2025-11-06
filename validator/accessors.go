package validator

import (
	"strconv"
	"time"

	"github.com/araddon/dateparse"
	"github.com/google/uuid"
)

func getTimeValue(v any) (time.Time, error) {
	switch value := v.(type) {
	case time.Time:
		return value, nil

	case string:
		t, err := dateparse.ParseAny(value)
		if err != nil {
			return time.Time{}, ErrInvalidData.Value(value)
		}

		return t, nil

	default:
		return time.Time{}, ErrInvalidData.Value(value)
	}
}

func getUUIDValue(v any) (uuid.UUID, error) {
	switch value := v.(type) {
	case uuid.UUID:
		return value, nil

	case string:
		if value == "" {
			return uuid.Nil, nil
		}

		id, err := uuid.Parse(value)
		if err != nil {
			return uuid.Nil, ErrInvalidData.Value(value)
		}

		return id, nil

	case []byte:
		id, err := uuid.FromBytes(value)
		if err != nil {
			return uuid.Nil, ErrInvalidData.Value(value)
		}

		return id, nil

	default:
		return uuid.Nil, ErrInvalidData.Value(value)
	}
}

func getIntValue(v any) (int, error) {
	switch value := v.(type) {
	case int:
		return value, nil
	case float64:
		return int(value), nil
	case string:
		n, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return 0, ErrInvalidData.Value(value)
		}

		return int(n), nil

	default:
		return 0, ErrInvalidData.Value(value)
	}
}

func getFloatValue(v any) (float64, error) {
	switch value := v.(type) {
	case float64:
		return value, nil
	case int:
		return float64(value), nil

	case string:
		n, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0, ErrInvalidData.Value(value)
		}

		return n, nil

	default:
		return 0, ErrInvalidData.Value(value)
	}
}

func getStringValue(v any) (string, error) {
	switch value := v.(type) {
	case string:
		return value, nil
	default:
		return "", ErrInvalidData.Value(value)
	}
}

func getBoolValue(v any) (bool, error) {
	switch value := v.(type) {
	case bool:
		return value, nil

	case string:
		switch value {
		case "true":
			return true, nil
		case "false":
			return false, nil
		}

		return false, ErrInvalidData.Value(value)

	default:
		return false, ErrUnsupportedType
	}
}
