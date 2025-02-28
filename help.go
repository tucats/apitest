package main

import (
	"fmt"
	"os"
)

var helpText = `
apitest - A simple JSON-driven rest API testing tool (c) Copyright 2025 Tom Cole
		  
usage: apitest [options]

options:

  -d, --define <key=value>  Define a value for a variable in the test dictionary (can be repeated)
  -h, --help                Show this help message and exit
  -p, --path <path>         Path to the test suite directory (required)
  -r, --rest                Enable REST mode, which makes the server return a JSON response
  -v, --Verbose             Enable Verbose output
  
  See the project README.md file for information on the format of test files that are located
  in the test path directory tree.
  `

func help() {
	fmt.Println(helpText)

	os.Exit(0)
}
