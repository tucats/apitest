package parser

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func parse(body interface{}, item string) ([]string, error) {
	var (
		index   int
		isIndex bool
		name    string
	)

	// If the item is just a "dot" it means the entire body is the result
	if item == "." {
		return []string{fmt.Sprintf("%v", body)}, nil
	}

	item = dotQuote(item)

	// Split out the item we seek plus whatever might be after it
	parts := strings.SplitN(item, ".", 2)
	if len(parts) == 1 {
		parts = append(parts, ".")
	}

	for i, part := range parts {
		if i == 0 {
			parts[i] = dotUnquote(part, ".")
		} else {
			parts[i] = dotUnquote(part, "\\.")
		}
	}

	// Determine if this is a numeric index or a name
	if i, err := strconv.Atoi(parts[0]); err == nil {
		index = i
		isIndex = true
	} else {
		name = parts[0]
	}

	// Is the name the wildcard array index?
	if name == "*" {
		return anyArrayElement(body, parts, item)
	}

	// If it's an index, the current item must be an array
	if isIndex {
		return arrayElement(body, index, parts, item)
	}

	// If it's a name, the current item must be a map of some type.
	val := reflect.ValueOf(body)

	if val.Kind() == reflect.Map {
		for _, e := range val.MapKeys() {
			if e.String() == name {
				v := val.MapIndex(e)

				return parse(v.Interface(), parts[1])
			}
		}

		return nil, fmt.Errorf("Map element not found: %s", name)
	}

	return nil, fmt.Errorf("Item is not a map, item, or array: %T", item)
}

func dotQuote(s string) string {
	return strings.ReplaceAll(s, "\\.", "$$DOT$$")
}

func dotUnquote(s string, target string) string {
	return strings.ReplaceAll(s, "$$DOT$$", target)
}
