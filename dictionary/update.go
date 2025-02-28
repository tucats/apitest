package dictionary

import (
	"fmt"

	"github.com/tucats/apitest/logging"
	"github.com/tucats/apitest/parser"
)

func Update(text string, items map[string]string) error {
	for key, value := range items {
		item, err := parser.GetOneItem(text, value)
		if err != nil {
			return err
		}

		Dictionary[key] = item

		if logging.Verbose {
			if key == "API_TOKEN" {
				item = "***REDACTED***"
			}

			fmt.Printf("  Updating   {{%s}} = %s\n", key, item)
		}
	}

	return nil
}
