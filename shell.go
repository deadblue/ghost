package ghost

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

/*
The implementation of Shell.
*/
type _ShellImpl[Ghost any] struct {
	// Listener network and address
	ln, la string
	// HTTP server
	hs *http.Server
	// Request dispatch kernel
	kn *_Kernel[Ghost]
	// Closed flag
	cf int32
	// Error channel
	ec chan error
}

func (s *_ShellImpl[Ghost]) die(err error) {
	if atomic.CompareAndSwapInt32(&s.cf, 0, 1) {
		_ = s.kn.AfterShutdown()
		if err != nil {
			s.ec <- err
		}
		close(s.ec)
	}
}

func (s *_ShellImpl[Ghost]) Startup() error {
	if err := s.kn.BeforeStartup(); err != nil {
		return err
	}
	// Start network listener
	l, err := net.Listen(s.ln, s.la)
	if err != nil {
		_ = s.kn.AfterShutdown()
		return err
	}
	log.Printf("Shell working at: %s://%s", s.ln, s.la)
	// For unix listener, delete it after close.
	if ul, ok := l.(*net.UnixListener); ok {
		ul.SetUnlinkOnClose(true)
	}
	go func(nl net.Listener) {
		err := s.hs.Serve(nl)
		if err != nil && err != http.ErrServerClosed {
			s.die(err)
		}
	}(l)
	return nil
}

func (s *_ShellImpl[Ghost]) Shutdown() {
	go func() {
		if atomic.LoadInt32(&s.cf) == 0 {
			err := s.hs.Shutdown(context.Background())
			s.die(err)
		}
	}()
}

func (s *_ShellImpl[Ghost]) Done() <-chan error {
	return s.ec
}

func (s *_ShellImpl[Ghost]) Run(sig ...os.Signal) (err error) {
	sc := make(chan os.Signal)
	if len(sig) > 0 {
		signal.Notify(sc, sig...)
	} else {
		// Handle SIGINT and SIGTERM by default.
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	}
	defer func() {
		signal.Stop(sc)
		close(sc)
	}()
	// Startup server
	if err = s.Startup(); err != nil {
		return err
	}
	// Loop
	for running := true; running; {
		select {
		case <-sc:
			log.Println("killing the shell ...")
			// Shutdown when receive OS signal
			s.Shutdown()
		case err = <-s.Done():
			log.Println("Shell is dead!")
			running = false
		}
	}
	return
}

func createShell[Ghost any](network, address string, kernel *_Kernel[Ghost]) Shell {
	return &_ShellImpl[Ghost]{
		// Listener network and address
		ln: network,
		la: address,
		// HTTP server and handler
		hs: &http.Server{
			Handler: kernel,
		},
		kn: kernel,

		cf: 0,
		ec: make(chan error),
	}
}
