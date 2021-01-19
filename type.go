package ghost

import (
	"io"
	"net/http"
	"os"
)

// Shell is a lifeless object, until developer gives an interesting ghost to it.
type Shell interface {

	// Startup starts up the shell manually, use this when you want to
	// control the shell lifecycle by yourself.
	// Otherwise, use Run instead.
	Startup() error

	// Shutdown shuts down the shell manually, use this when you want to
	// control the shell lifecycle by yourself.
	// Otherwise, use Run instead.
	Shutdown()

	// Done returns a read-only error channel, you will get notification
	// from it when the shell completely shutdown, use this when you
	// control the shell lifecycle by yourself.
	// Otherwise, use Run instead.
	Done() <-chan error

	// Run automatically runs the shell, and shutdown it when receive specific
	// OS signals, Run will exit after the shell completely shutdown.
	// If no signal specified, handles SIGINT and SIGTERM as default.
	Run(sig ...os.Signal) error
}

// Context describes the request context.
type Context interface {

	// Request returns the original HTTP request object.
	Request() *http.Request

	// Scheme returns request scheme(http/https).
	Scheme() string

	// Host returns request server host.
	Host() string

	// Method returns request method in upper-case.
	Method() string

	// Path returns request path.
	Path() string

	// RemoteIp returns the client IP.
	RemoteIp() string

	// PathVar return the variable value in request path.
	PathVar(name string) string

	// Header returns a value in request header with given name.
	Header(name string) string

	// HeaderArray returns all values in request header who has the given name.
	HeaderArray(name string) []string

	// Query returns the parameter value in query string who has the given name.
	Query(name string) string

	// QueryArray returns all parameter values in query string who has the given name.
	QueryArray(name string) []string

	// Cookie returns the cookie value who has the given name.
	Cookie(name string) string

	// CookieArray returns all cookie values who has the given name.
	CookieArray(name string) []string

	// Body return the request body, do not use it in individual goroutine,
	// because it will be closed after controller return.
	Body() io.Reader
}

// View describes the response.
type View interface {

	// HTTP status code
	Status() int

	// Response body
	Body() io.Reader
}

// Controller is a function to handle request.
type Controller func(ctx Context) (v View, err error)
