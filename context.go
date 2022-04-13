package ghost

import (
	"io"
	"net/http"
)

// Context wraps HTTP request, and provides some helpful methods.
type Context interface {

	// Request returns the original HTTP request object.
	Request() *http.Request

	// Method returns request method.
	Method() string

	// Path returns request path.
	Path() string

	// Scheme returns request scheme(http/https).
	Scheme() string

	// Host returns request server host.
	Host() string

	// Body return a reader for request body.
	Body() io.Reader

	// RemoteIp returns the client IP.
	RemoteIp() string

	// PathVar return the variable value in request path.
	PathVar(name string) string
}
