package ghost

import (
	"log"
	"net/http"
)

const (
	DefaultNetwork = "tcp"
	DefaultAddress = "127.0.0.1:9057"
)

// Born creates a Shell with your ghost, which will listen at default
// network and address.
func Born[Ghost any](ghost Ghost, options ...Option) Shell {
	// Create engine
	engine := &_EngineImpl[Ghost]{}
	engine.install(ghost)
	// Create shell
	shell := &_ShellImpl{
		// Listener network and address
		ln: DefaultNetwork,
		la: DefaultAddress,
		// HTTP server
		hs: &http.Server{
			Handler: engine,
		},
		e: engine,

		cf:    0,
		errCh: make(chan error),
	}
	// Apply options
	// TODO: Support more options
	for _, opt := range options {
		switch opt.(type) {
		case optListen:
			ol := opt.(optListen)
			shell.ln, shell.la = ol.network, ol.address
		default:
			log.Printf("Unsupported option: %#v", opt)
		}
	}
	return shell
}
