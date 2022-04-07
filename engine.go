package ghost

import (
	"fmt"
	"github.com/deadblue/ghost/internal/context"
	"github.com/deadblue/ghost/internal/view"
	"io"
	"log"
	"mime"
	"net/http"
	"path"
	"reflect"
	"runtime"
	"strconv"
)

const (
	// Special response headers
	_HeaderServer        = "Server"
	_HeaderContentType   = "Content-Type"
	_HeaderContentLength = "Content-Length"
	// Default header value
	_DefaultContentType = "application/octet-stream"
)

var (
	// Server token in response header
	_ServerToken = fmt.Sprintf("Ghost/%s (%s/%s %s)", Version,
		runtime.GOOS, runtime.GOARCH, runtime.Version())
)

type _Engine interface {
	http.Handler
	StartupObserver
	ShutdownObserver
}

type _EngineImpl[Ghost any] struct {
	// The ghost which implemented by user
	g Ghost
	// Route tree
	rt *_RouteTable
	// HTTP status observer
	sh StatusHandler
}

func (e *_EngineImpl[Ghost]) BeforeStartup() (err error) {
	var ghost any = e.g
	if ob, isImpl := ghost.(StartupObserver); isImpl {
		err = ob.BeforeStartup()
	}
	return
}

func (e *_EngineImpl[Ghost]) AfterShutdown() (err error) {
	var ghost any = e.g
	if ob, isImpl := ghost.(ShutdownObserver); isImpl {
		err = ob.AfterShutdown()
	}
	return
}

func (e *_EngineImpl[Ghost]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Make context
	ctx := (&context.Impl{}).FromRequest(r)
	// Resolve controller
	h, v := e.rt.Resolve(ctx), View(nil)
	// Invoke controller
	if h == nil {
		v = e.sh.OnStatus(http.StatusNotFound, ctx, nil)
	} else {
		v = e.invoke(h, ctx)
	}
	// Render view
	e.render(ctx, v, w)
}

func (e *_EngineImpl[Ghost]) invoke(h _Handler[Ghost], ctx Context) (v View) {
	defer func() {
		// Catch panic here
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				v = e.sh.OnStatus(http.StatusInternalServerError, ctx, err)
			} else {
				v = e.sh.OnStatus(http.StatusInternalServerError, ctx, fmt.Errorf("panic: %v", r))
			}
		}
	}()
	v, err := h(ctx)
	if err != nil {
		v = e.sh.OnStatus(http.StatusInternalServerError, ctx, err)
	}
	return
}

func (e *_EngineImpl[Ghost]) render(ctx Context, v View, w http.ResponseWriter) {
	// Ensure the view is not nil
	if v == nil {
		v = e.defaultView()
	}
	headers, body := w.Header(), v.Body()
	// Prepare response header
	headers.Set(_HeaderServer, _ServerToken)
	// Allow view manipulates the response header
	if hi, ok := v.(ViewHeaderInterceptor); ok {
		hi.BeforeSendHeader(headers)
	}
	// Set Content-Type and Content-Length at last
	if body != nil {
		// Try to set "Content-Type", when view does not set.
		if headers.Get(_HeaderContentType) == "" {
			ct := ""
			if a, ok := v.(ViewTypeAdviser); ok {
				ct = a.ContentType()
			} else {
				ct = mime.TypeByExtension(path.Ext(ctx.Path()))
			}
			if ct == "" {
				ct = _DefaultContentType
			}
			headers.Set(_HeaderContentType, ct)
		}
		// Try to set "Content-Length" when view does not set.
		if headers.Get(_HeaderContentLength) == "" {
			size := int64(0)
			if a, ok := v.(ViewSizeAdviser); ok {
				// Get size from view
				size = a.ContentLength()
			} else if l, ok := body.(bodyHasLength); ok {
				// Auto-detect body size
				size = int64(l.Len())
			}
			if size > 0 {
				headers.Set(_HeaderContentLength, strconv.FormatInt(size, 10))
			}
		}
	}
	// Send response header
	w.WriteHeader(v.Status())
	// Send response body
	if body != nil {
		// Auto close the closable body
		if c, ok := body.(io.Closer); ok {
			defer func() {
				_ = c.Close()
			}()
		}
		_, _ = io.Copy(w, body)
	}
}

func (e *_EngineImpl[Ghost]) install(ghost Ghost) {
	// Initialize engine
	e.g, e.rt = ghost, &_RouteTable{
		st: make(map[string]_Handler),
	}
	// Scan implemented interfaces on ghost
	var v any = ghost
	if sh, isImpl := v.(StatusHandler); isImpl {
		e.sh = sh
	} else {
		e.sh = defaultStatusHandler
	}
	// Scan controller
	rt := reflect.TypeOf(ghost)
	for i := 0; i < rt.NumMethod(); i++ {
		m := rt.Method(i)
		mn, mf := m.Name, m.Func.Interface()
		// Check method function signature
		if fn, ok := mf.(func(Ghost, Context) (View, error)); ok {
			hdl := func(ctx Context) (View, error) {
				return fn(ghost, ctx)
			}
			if err := e.rt.Mount(mn, hdl); err == nil {
				log.Printf("Mount controller: %s", m.Name)
			}
		}
	}
}

// defaultView returns HTTP 200 empty view.
func (e *_EngineImpl[Ghost]) defaultView() View {
	return view.Http200
}
