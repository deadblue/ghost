package ghost

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/deadblue/ghost/internal/container"
	"github.com/deadblue/ghost/internal/context"
	"github.com/deadblue/ghost/internal/route"
	"github.com/deadblue/ghost/internal/view"
	"io"
	"mime"
	"net/http"
	"path"
	"runtime"
	"strconv"
)

const (
	// Server header
	_HeaderServer = "Server"
	// Content headers
	_HeaderContentType   = "Content-Type"
	_HeaderContentLength = "Content-Length"

	// Default header value
	_DefaultContentType = "application/octet-stream"
)

var (
	// Server token in response header
	_ServerToken = fmt.Sprintf("Ghost/%s (%s/%s %s)", Version,
		runtime.GOOS, runtime.GOARCH, runtime.Version())

	ErrNotFound = errors.New("not found")
)

type _Engine struct {
	// Route registry
	rt *route.Registry[_Handler]

	// Startup observers
	startupObs container.List[StartupObserver]
	// Shutdown observers
	shutdownObs container.List[ShutdownObserver]
	// Error handler
	eh ErrorHandler
}

func (e *_Engine) BeforeStartup() (err error) {
	for ok := e.startupObs.GoFirst(); ok; ok = e.startupObs.Forward() {
		_, ob := e.startupObs.Get()
		if err = ob.BeforeStartup(); err != nil {
			return
		}
	}
	return
}

func (e *_Engine) AfterShutdown() (err error) {
	for ok := e.shutdownObs.GoLast(); ok; ok = e.shutdownObs.Backward() {
		_, ob := e.shutdownObs.Get()
		if err = ob.AfterShutdown(); err != nil {
			return
		}
	}
	return
}

func (e *_Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Make context
	ctx := (&context.Impl{}).FromRequest(r)
	// Handle special request method
	method := ctx.Method()
	switch method {
	case http.MethodHead:
		// For HEAD request, change method to GET
		method = http.MethodGet
		w = &_HeadResponseWriter{rw: w}
	case http.MethodOptions:
		// TODO: Handle CORS request
	}
	// Resolve handler
	h := e.rt.Resolve(method, ctx.Path(), ctx)
	v := e.dispatch(h, ctx)
	e.renderView(ctx, v, w)
}

// dispatch calls handler with context, and catch error.
func (e *_Engine) dispatch(h _Handler, ctx Context) (v View) {
	if h == nil {
		return e.eh.OnError(ctx, ErrNotFound)
	}
	defer func() {
		// Catch panic here
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				v = e.eh.OnError(ctx, err)
			} else {
				v = e.eh.OnError(ctx, fmt.Errorf("panic: %v", r))
			}
		}
	}()
	var err error
	if v, err = h.Handle(ctx); err != nil {
		v = e.eh.OnError(ctx, err)
	}
	return
}

func (e *_Engine) renderView(ctx Context, v View, w http.ResponseWriter) {
	// Ensure the view is not nil
	if v == nil {
		v = view.Http200
	}
	headers, body := w.Header(), v.Body()
	// Prepare response header
	headers.Set(_HeaderServer, _ServerToken)
	// Allow view manipulates the response header
	if hi, ok := v.(ViewHeaderInterceptor); ok {
		hi.BeforeSendHeader(headers)
	}
	// Set Content-Type and Content-Length
	if body != nil {
		// Try to set "Content-Type", when view does not set it.
		if headers.Get(_HeaderContentType) == "" {
			headers.Set(_HeaderContentType, determineContentType(v, path.Ext(ctx.Path())))
		}
		// Try to set "Content-Length" when view does not set it.
		if headers.Get(_HeaderContentLength) == "" {
			if size := determineContentLength(v); size >= 0 {
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

func determineContentType(v View, ext string) string {
	if a, ok := v.(ViewTypeAdviser); ok {
		return a.ContentType()
	}
	if mt := mime.TypeByExtension(ext); mt != "" {
		return mt
	}
	return _DefaultContentType
}

type _SizeInterface interface {
	Size() int64
}

func determineContentLength(v View) int64 {
	if a, ok := v.(ViewSizeAdviser); ok {
		return a.ContentLength()
	}
	body := v.Body()
	if i, ok := body.(_SizeInterface); ok {
		return i.Size()
	}
	if b, ok := body.(*bytes.Buffer); ok {
		return int64(b.Len())
	}
	return -1
}
