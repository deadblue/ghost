package ghost

/*
Binder is an optional but strongly recommended interface that need developer to implement
on his ghost, that will speed up controller invoking.

See the benchmark for detail:
https://gist.github.com/deadblue/b232340144acd20f48d38602fd628a1b#file-benchmark_test-go

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

// StartupObserver is an optional interface for developer's ghost.
type StartupObserver interface {

	// BeforeStartup will be called before Shell starts up, developer can do
	// initializing works in it. When BeforeStartup returns an error, the shell
	// won't start up, and return the error.
	BeforeStartup() error
}

// ShutdownObserver is an optional interface for developer's ghost.
type ShutdownObserver interface {

	// AfterShutdown will always be called after Shell completely shut down, even
	// Shell shut down with an error, developer can do finalizing works in it.
	// Currently, the returned error will be ignored.
	AfterShutdown() error
}

// StatusHandler is an optional interface for developer's ghost, it allows developer to
// customize the view when HTTP 40x and 50x error occurred.
type StatusHandler interface {

	// OnStatus will be called when HTTP 40x and 50x error occurred.
	OnStatus(status int, ctx Context, err error) View
}
