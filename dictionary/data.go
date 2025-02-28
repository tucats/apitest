package dictionary

import (
	"sync/atomic"
)

var Dictionary = make(map[string]string)

var sequence = atomic.Int32{}
