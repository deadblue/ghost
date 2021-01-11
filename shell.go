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

type _ShellImpl struct {
	// Listener network and address
	ln, la string
	// HTTP server
	hs *http.Server
	// Request dispatch kernel
	kn *_Kernel
	// Closed flag
	cf int32
	// Error channel
	ec chan error
}

func (s *_ShellImpl) die(err error) {
	if atomic.CompareAndSwapInt32(&s.cf, 0, 1) {
		if err != nil {
			s.ec <- err
		}
		close(s.ec)
	}
}

func (s *_ShellImpl) Startup() error {
	// Start network listener
	l, err := net.Listen(s.ln, s.la)
	if err != nil {
		return err
	}
	log.Printf("Shell working at: %s://%s", s.ln, s.la)
	// For unix listener, delete it after close.
	if ul, ok := l.(*net.UnixListener); ok {
		ul.SetUnlinkOnClose(true)
	}
	go func(nl net.Listener) {
		err := s.kn.BeforeStartup()
		if err == nil {
			err = s.hs.Serve(nl)
		}
		if err != nil && err != http.ErrServerClosed {
			_ = s.kn.AfterShutdown()
			s.die(err)
		}
	}(l)
	return nil
}

func (s *_ShellImpl) Shutdown() {
	go func() {
		if atomic.LoadInt32(&s.cf) == 0 {
			err := s.hs.Shutdown(context.Background())
			_ = s.kn.AfterShutdown()
			s.die(err)
		}
	}()
}

func (s *_ShellImpl) Done() <-chan error {
	return s.ec
}

func (s *_ShellImpl) Run(sig ...os.Signal) (err error) {
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

// Born creates a Shell with your ghost, and will listen at default
// network and address: "http://127.0.0.1:8066".
func Born(ghost interface{}) Shell {
	return BornAt(ghost, "tcp", "127.0.0.1:8066")
}

// BornAt creates a Shell with your ghost, and will listen at the
// network and address where you give.
func BornAt(ghost interface{}, network, address string) Shell {
	// Create kernel
	k := (&_Kernel{}).Install(ghost)
	// Make shell
	s := &_ShellImpl{
		// Listener network and address
		ln: network,
		la: address,
		// HTTP server and handler
		hs: &http.Server{
			Handler: k,
		},
		kn: k,

		cf: 0,
		ec: make(chan error),
	}
	return s
}
