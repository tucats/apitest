package tester

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tucats/apitest/defs"
	"github.com/tucats/apitest/dictionary"
	"github.com/tucats/apitest/logging"
)

func executeTask(task defs.Task) error {
	var err error

	switch strings.ToLower(task.Command) {
	case "delete":
		for _, name := range task.Parameters {
			name = dictionary.Apply(name)

			name, err = filepath.Abs(filepath.Clean(name))
			if err != nil {
				return err
			}

			if logging.Verbose {
				fmt.Printf("  Task: deleting file: %s\n", name)
			}

			err = os.Remove(name)
			if err != nil {
				return err
			}
		}

	default:
		err = fmt.Errorf("Unknown task command: %s", task.Command)
	}

	return err
}
