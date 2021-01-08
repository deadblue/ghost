package ghost

// AwareStartup is an optional interface, when developer need do some initialization
// on his ghost, implements this, and do initialization in OnStartup().
type AwareStartup interface {
	OnStartup() error
}

// AwareShutdown is an optional interface, when developer need do some finalization
// on his ghost, implements this, and do finalization in OnShutdown().
type AwareShutdown interface {
	OnShutdown() error
}

/*
Binder is an optional interface, if developer implements it on his ghost, he will
get 50x speedup on controller invoking.

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
