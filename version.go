package ghost

import (
	"fmt"
	"runtime"
)

const (
	Version = "0.0.3"
)

var (
	_Server = fmt.Sprintf("Ghost/%s (%s/%s %s)", Version,
		runtime.GOOS, runtime.GOARCH, runtime.Version())
)
