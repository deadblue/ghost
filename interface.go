package ghost

import "net/http"

// StartupHandler is an optional interface, when developer need do some initialization
// on his ghost, implements this, and do initialization in OnStartup().
type StartupHandler interface {
	OnStartup() error
}

// ShutdownHandler is an optional interface, when developer need do some finalization
// on his ghost, implements this, and do finalization in OnShutdown().
type ShutdownHandler interface {
	OnShutdown() error
}

// HttpStatusHandler is an optional interface, when developer wants to customize the
// error view, implement this on his ghost.
type HttpStatusHandler interface {

	// OnStatus will be called when HTTP 40x and 50x error occurred.
	OnStatus(status int, context Context, err error) View
}

// HeaderInterceptor is an interface that can optionally implement by View and developer's
// ghost. It will be called after kernel sets normal response headers, developer can manipulate
// response header here.
type HeaderInterceptor interface {

	// BeforeSend will be called before kernel send response headers to client.
	BeforeSend(h http.Header)
}

/*
Binder is an optional interface, but it is highly recommended developer to implement this, to
speed up the controller invocation.

Benchmark: https://gist.github.com/deadblue/b232340144acd20f48d38602fd628a1b#file-benchmark_test-go

A standard implementation looks like this:

	func (g *YourGhost) Bind(f interface{}) ghost.Controller {
		if ctrl, ok := f.(func(*YourGhost, ghost.Context)(ghost.View, error)); ok {
			return func(ctx ghost.Context) (ghost.View, error) {
				return ctrl(g, ctx)
			}
		} else {
			return nil
		}
	}
*/
type Binder interface {
	Bind(f interface{}) Controller
}
