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

type _Kernel[Ghost any] struct {
	// The ghost which implemented by user
	g Ghost
	// Route tree
	rt *_RouteTable
	// HTTP status observer
	sh StatusHandler
}

func (k *_Kernel[Ghost]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Make context
	ctx := context.FromRequest(r)
	// Resolve controller
	ctrl, v := k.rt.Resolve(ctx), View(nil)
	// Invoke controller
	if ctrl == nil {
		v = k.sh.OnStatus(http.StatusNotFound, ctx, nil)
	} else {
		v = k.invoke(ctrl, ctx)
	}
	// Render view
	k.render(ctx, v, w)
}

func (k *_Kernel[Ghost]) BeforeStartup() (err error) {
	var i any = k.g
	if o, isImpl := i.(StartupObserver); isImpl {
		err = o.BeforeStartup()
	}
	return
}

func (k *_Kernel[Ghost]) AfterShutdown() (err error) {
	var i interface{} = k.g
	if o, isImpl := i.(ShutdownObserver); isImpl {
		err = o.AfterShutdown()
	}
	return
}

func (k *_Kernel[Ghost]) install(ghost Ghost) {
	// Initial kernel
	k.g, k.rt = ghost, &_RouteTable{
		mapping:  make(map[_RouteKey]_Controller),
		branches: make(map[string][]*_RoutePath),
	}
	// Setup HTTP status handler
	var i interface{} = ghost
	if sh, ok := i.(StatusHandler); ok {
		k.sh = sh
	} else {
		k.sh = defaultStatusHandler
	}
	// Scan controller
	rt := reflect.TypeOf(ghost)
	for i := 0; i < rt.NumMethod(); i++ {
		m := rt.Method(i)
		mf := m.Func.Interface()
		if f, ok := mf.(func(Ghost, Context) (View, error)); ok {
			ctrl := func(ctx Context) (View, error) {
				return f(ghost, ctx)
			}
			if err := k.rt.Mount(m.Name, ctrl); err == nil {
				log.Printf("Mount controller: %s", m.Name)
			}
		}
	}
}

func (k *_Kernel[Ghost]) invoke(ctrl _Controller, ctx Context) (v View) {
	defer func() {
		// Catch panic here
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				v = k.sh.OnStatus(http.StatusInternalServerError, ctx, err)
			} else {
				v = k.sh.OnStatus(http.StatusInternalServerError, ctx, fmt.Errorf("panic: %v", r))
			}
		}
	}()
	v, err := ctrl(ctx)
	if err != nil {
		v = k.sh.OnStatus(http.StatusInternalServerError, ctx, err)
	}
	return
}

func (k *_Kernel[Ghost]) render(ctx Context, v View, w http.ResponseWriter) {
	// Ensure the view is not nil
	if v == nil {
		v = k.defaultView()
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

// defaultView returns HTTP 200 empty view.
func (k *_Kernel[Ghost]) defaultView() View {
	return view.Http200
}

func createKernel[Ghost any](ghost Ghost) *_Kernel[Ghost] {
	k := &_Kernel[Ghost]{}
	k.install(ghost)
	return k
}
