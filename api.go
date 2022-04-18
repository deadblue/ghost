package ghost

import (
	"github.com/deadblue/ghost/option"
	"log"
	"net/http"
	"time"
)

const (
	DefaultNetwork = "tcp"
	DefaultAddress = "127.0.0.1:9057"
)

// Born builds a Shell for your ghost.
func Born[Ghost any](ghost Ghost, options ...option.Option) Shell {
	// Create kernel
	kernel := &_Kernel{}
	internalImplant(kernel, ghost, "/", true)
	// Create shell
	shell := &_ShellImpl{
		// Listener network and address
		ln: DefaultNetwork,
		la: DefaultAddress,
		// HTTP server
		hs: &http.Server{
			Handler: kernel,
		},
		kn: kernel,

		cf:    0,
		errCh: make(chan error),
	}
	// Apply options
	applyOptions(shell, options)
	return shell
}

func applyOptions(shell *_ShellImpl, options []option.Option) {
	for _, opt := range options {
		switch opt.(type) {
		case option.ListenOption:
			no := opt.(option.ListenOption)
			shell.ln, shell.la = no.Network, no.Address
		case *option.TlsOption:
			shell.tc = opt.(*option.TlsOption).Config
		case option.ReadTimeoutOption:
			shell.hs.ReadTimeout = time.Duration(opt.(option.ReadTimeoutOption))
		case option.ReadHeaderTimeoutOption:
			shell.hs.ReadHeaderTimeout = time.Duration(opt.(option.ReadHeaderTimeoutOption))
		case option.WriteTimeoutOption:
			shell.hs.WriteTimeout = time.Duration(opt.(option.WriteTimeoutOption))
		case option.IdleTimeoutOption:
			shell.hs.IdleTimeout = time.Duration(opt.(option.IdleTimeoutOption))
		default:
			log.Printf("Unsupported option: %#v", opt)
		}
	}
}
