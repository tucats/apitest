package defs

type ResponseObject struct {
	Headers map[string][]string `json:"headers"`
	Status  int                 `json:"status"`
	Body    string              `json:"body"`
	Save    map[string]string   `json:"save,omitempty"`
}
