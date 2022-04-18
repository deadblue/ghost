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

	// Body returns a reader for request body.
	Body() io.Reader

	// Header returns request header with specific name.
	Header(name string) (string, bool)

	// Scheme returns request scheme(http/https).
	Scheme() string

	// Host returns request server host.
	Host() string

	// RemoteIp returns the client IP.
	RemoteIp() string

	// PathVar return the variable value in request path.
	PathVar(name string) string
}
