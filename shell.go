package ghost

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

/*
Shell is the shell of your ghost, it covers the basic reactions what an HTTP
server should do, and dispatches requests to your ghost.

You can use Shell in two ways:

1. Simply run it:

	// Create a shell from your ghost.
	shell := ghost.Born(&YourGhost{})

	// Way 1: Just run the shell, wait for it shut down completely.
	if err := shell.Run(); err != nil {
		panic(err)
	}

2. Manage its lifecycle by yourself:

	// Create a shell from your ghost.
	shell := ghost.Born(&YourGhost{})

	// Start up the shell.
	if err := shell.Startup(); err != nil {
		panic(err)
	}
	// Waiting for shell shut down.
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

	// Startup starts up the shell manually, use this when you want to manage
	// shell's lifecycle by yourself. Otherwise, just use `Run`.
	Startup() error

	// Shutdown shuts down the shell manually, use this when you want to manage
	// shell's lifecycle by yourself. Otherwise, just use Run.
	Shutdown()

	// Done returns a read-only error channel, you will get error events from it
	// when the shell shut down, use this when you manage shell's lifecycle by
	// yourself. Otherwise, just use Run.
	Done() <-chan error

	// Run automatically runs the shell, and shutdown it when receive specific
	// OS signals, Run will exit after the shell completely shutdown.
	// If no signal set, shell will handle SIGINT and SIGTERM as default.
	Run(sig ...os.Signal) error
}

// _ShellImpl is implementation of Shell interface.
type _ShellImpl struct {
	// Listener network and address
	ln, la string
	// TLS config
	tc *tls.Config
	// HTTP server
	hs *http.Server
	// kernel
	kn *_Kernel
	// Closed flag
	cf int32
	// Error channel
	errCh chan error
}

func (s *_ShellImpl) die(err error) {
	if atomic.CompareAndSwapInt32(&s.cf, 0, 1) {
		_ = s.kn.AfterShutdown()
		if err != nil {
			s.errCh <- err
		}
		close(s.errCh)
	}
}

func (s *_ShellImpl) serve(l net.Listener) {
	err := s.hs.Serve(l)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.die(err)
	}
}

func (s *_ShellImpl) Startup() (err error) {
	// Call
	if err = s.kn.BeforeStartup(); err != nil {
		return
	}
	// Start network listener
	var nl net.Listener
	if nl, err = net.Listen(s.ln, s.la); err != nil {
		_ = s.kn.AfterShutdown()
		return
	}
	// For unix listener, delete it after close.
	if ul, ok := nl.(*net.UnixListener); ok {
		ul.SetUnlinkOnClose(true)
	}
	// TLS listener
	if s.tc != nil {
		nl = tls.NewListener(nl, s.tc)
		log.Printf("Shell working at: %s+tls://%s", s.ln, s.la)
	} else {
		log.Printf("Shell working at: %s://%s", s.ln, s.la)
	}
	// Start serve
	go s.serve(nl)
	return
}

func (s *_ShellImpl) Shutdown() {
	go func() {
		if atomic.LoadInt32(&s.cf) == 0 {
			err := s.hs.Shutdown(context.Background())
			s.die(err)
		}
	}()
}

func (s *_ShellImpl) Done() <-chan error {
	return s.errCh
}

func (s *_ShellImpl) Run(sig ...os.Signal) (err error) {
	// Setup signal handler
	sigCh := make(chan os.Signal)
	if len(sig) > 0 {
		signal.Notify(sigCh, sig...)
	} else {
		// Handle SIGINT and SIGTERM by default.
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	}
	// Release signal handler after exit
	defer func() {
		signal.Stop(sigCh)
		close(sigCh)
	}()

	// Startup server
	if err = s.Startup(); err != nil {
		return err
	}
	// Loop
	for running := true; running; {
		select {
		case <-sigCh:
			log.Println("Killing the shell ...")
			// Shutdown when receive OS signal
			s.Shutdown()
		case err = <-s.Done():
			log.Println("Shell is gone!")
			running = false
		}
	}
	return
}
