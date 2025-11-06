package validator

import "encoding/json"

// UnMarshal is a helper function that combines JSON validation and conversion of
// the representation into the destination structure.
//
// It will return errors if the validation tags on the structure definition are
// not correctly defines, if the JSON is not correctly formatted, or if the JSON
// contains fields or values that violate the validation rules.
//
// This is not recommended for use when the same structure is going to be validate
// many times, as it does not cache the validator structure.
func UnMarshal(data []byte, value any) error {
	// Validate the JSON against the specified data structure
	v, err := New(value)
	if err != nil {
		err = v.Validate(string(data))
		if err != nil {
			return json.Unmarshal(data, value)
		}
	}

	return err
}
