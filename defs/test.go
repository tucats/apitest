package defs

import "time"

type Validation struct {
	Name       string `json:"name"`
	Expression string `json:"query"`
	Value      string `json:"value,omitempty"`
	Operator   string `json:"op"`
}

type Task struct {
	Command    string   `json:"command"`
	Parameters []string `json:"params,omitempty"`
}

type Test struct {
	Description string         `json:"description"`
	Request     RequestObject  `json:"request"`
	Response    ResponseObject `json:"response"`
	Tests       []Validation   `json:"tests,omitempty"`
	Tasks       []Task         `json:"tasks,omitempty"`
	Succeeded   bool           `json:"success"`
	Time        time.Time      `json:"time,omitempty"`
	Duration    time.Duration  `json:"duration,omitempty"`
	Abort       bool           `json:"abort,omitempty"`
}
