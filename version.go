package ghost

import (
	"fmt"
	"runtime"
)

const (
	Version = "0.0.1"
)

var (
	_HeaderServer = fmt.Sprintf("Ghost/%s (%s/%s %s)", Version,
		runtime.GOOS, runtime.GOARCH, runtime.Version())
)
