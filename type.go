package ghost

import (
	"io"
	"net/http"
	"os"
)

/*
Shell is the shell of the developer made ghost, it covers the basic reactions what an HTTP
server should do, and dispatches requests to the ghost.

Developer has two ways to use the shell: manually manage the lifecycle, or just run it.
Here are the examples of the two ways:

	// Create a shell from developer's ghost.
	shell := ghost.Born(&YourGhost{})

	// Way 1: Just run the shell, wait for it shut down completely.
	if err := shell.Run(); err != nil {
		panic(err)
	}

	// Way 2: Manually manage the lifecycle of the shell.
	// Start up the shell.
	if err := shell.Startup(); err != nil {
		panic(err)
	}
	for running := true; running; {
		select {
		case <- someEventArrived:
			// Shut down the shell.
			shell.Shutdown()
		case err := <- shell.Done():
			// The shell completely shut down.
			if err != nil {
				// Shut down with error
			} else {
				// Normally shut down.
			}
			// Exit the for-loop.
			running = false;
		}
	}
*/
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
