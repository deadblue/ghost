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

	// PathVar return the variable value in request path.
	PathVar(name string) string

	// Query returns the parameter value in query string.
	Query(name string) string
}

// View describes the response.
type View interface {

	// HTTP status code
	Status() int

	// Response headers
	Header() http.Header

	// Response body
	Body() io.Reader
}

type Controller func(ctx Context) (v View, err error)
