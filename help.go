package main

import (
	"fmt"
	"os"

	"github.com/tucats/apitest/dictionary"
)

var helpText = `
apitest {{VERSION}} - A simple JSON-driven rest API testing tool (C) 2025 Tom Cole
		  
usage: apitest [options] test-path

options:

  -d, --define <key=value>  Define a value for a variable in the test dictionary (can be repeated)
  -f, --filter <string>     Only run tests that contain the given string in their names
  -h, --help                Show this help message and exit
  -r, --rest                Enable REST mode, which makes the server return a JSON response
  -v, --Verbose             Enable Verbose output
  
  See the project README.md file for information on the format of test files that are located
  in the test path directory tree.
  `

func help() {
	fmt.Println(dictionary.Apply(helpText))

	os.Exit(0)
}
