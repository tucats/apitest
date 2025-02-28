package defs

type RequestObject struct {
	Endpoint   string              `json:"endpoint"`
	Parameters map[string]string   `json:"parameters,omitempty"`
	Headers    map[string][]string `json:"headers,omitempty"`
	Method     string              `json:"method"`
	Body       string              `json:"body,omitempty"`
	File       string              `json:"file,omitempty"`
}
