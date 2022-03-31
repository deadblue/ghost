package ghost

// StartupObserver is an optional interface for your ghost.
type StartupObserver interface {

	// BeforeStartup will be called before Shell starts up, you can do some
	// initializing jos here. When BeforeStartup returns an error, the shell
	// won't start up, and return the error.
	BeforeStartup() error
}

// ShutdownObserver is an optional interface for your ghost.
type ShutdownObserver interface {

	// AfterShutdown will always be called after Shell completely shut down, even
	// Shell shuts down with an error, developer can do finalizing works in it.
	// Currently, the returned error will be ignored.
	AfterShutdown() error
}

// StatusHandler is an optional interface for your ghost, which allows developer to
// customize the view when HTTP 40x and 50x error occurred.
type StatusHandler interface {

	// OnStatus will be called when HTTP 40x and 50x error occurred.
	OnStatus(status int, ctx Context, err error) View
}
