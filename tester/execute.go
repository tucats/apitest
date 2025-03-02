package tester

import (
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tucats/apitest/defs"
	"github.com/tucats/apitest/dictionary"
	"gopkg.in/resty.v1"
)

func ExecuteTest(test *defs.Test) error {
	var (
		err  error
		kind contentType = unknownContent
	)

	// Form the URL string
	urlString := test.Request.Endpoint

	parms := make([]string, 0)
	for key, value := range test.Request.Parameters {
		parms = append(parms, fmt.Sprintf("%s=%s", key, value))
	}

	if len(parms) > 0 {
		urlString += "?" + strings.Join(parms, "&")
	}

	// Create an HTTP client
	client := resty.New()
	tlsConfiguration := &tls.Config{InsecureSkipVerify: true}
	client.SetTLSClientConfig(tlsConfiguration)

	r := client.NewRequest()

	// Update the body, headers and URLstring with the dictionary values
	for key, values := range test.Request.Headers {
		for _, value := range values {
			value = dictionary.Apply(value)

			r.Header.Add(key, value)

			if strings.EqualFold(key, "content-type") {
				v := strings.ToLower(value)
				if strings.Contains(v, "json") {
					kind = jsonContent
				} else if strings.Contains(v, "text") {
					kind = textContent
				}
			}
		}
	}

	urlString = dictionary.Apply(urlString)

	// If the request body is a file specification, substitute that now.
	if test.Request.File != "" {
		test.Request.File = dictionary.Apply(test.Request.File)

		path, err := filepath.Abs(test.Request.File)
		if err != nil {
			return err
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		test.Request.Body = string(data)
	}

	b := []byte(dictionary.Apply(test.Request.Body))

	r.Body = b

	restLog("Request body", b, kind)

	// Make the HTTP request
	now := time.Now()

	resp, err := r.Execute(test.Request.Method, urlString)
	if err != nil {
		return err
	}

	test.Duration = time.Since(now)

	// Verify that the response status code matches the expected status code
	if resp.StatusCode() != test.Response.Status {
		return fmt.Errorf("%s, expected status %d, got %d", test.Description, test.Response.Status, resp.StatusCode())
	}

	// Validate the response body if present
	b = resp.Body()
	if len(b) > 0 {
		test.Response.Body = string(b)

		kind = unknownContent

		for key, value := range resp.Header() {
			if strings.EqualFold(key, "content-type") {
				v := strings.ToLower(strings.Join(value, ","))
				if strings.Contains(v, "json") {
					kind = jsonContent
				} else if strings.Contains(v, "text") {
					kind = textContent
				}
			}
		}

		restLog("Response body", b, kind)

		err = validateTest(test)
	}

	// If there were no errors, execute any tasks in the test.
	if err == nil {
		for _, task := range test.Tasks {
			err = executeTask(task)
			if err != nil {
				return err
			}
		}
	}

	return err
}
