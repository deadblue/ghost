package ghost

import (
	"fmt"
	"github.com/deadblue/ghost/internal/context"
	"github.com/deadblue/ghost/internal/route"
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

type _Engine struct {
	// Route registry
	rt *route.Registry[_Handler]

	// Startup observers
	startObs []StartupObserver
	// Shutdown observers
	shutObs []ShutdownObserver

	// HTTP status observer
	sh StatusHandler
}

func (e *_Engine) BeforeStartup() (err error) {
	if e.startObs == nil {
		return
	}
	for _, ob := range e.startObs {
		if err = ob.BeforeStartup(); err != nil {
			return
		}
	}
	return
}

func (e *_Engine) AfterShutdown() (err error) {
	if e.shutObs == nil {
		return
	}
	for _, ob := range e.shutObs {
		_ = ob.AfterShutdown()
	}
	return
}

func (e *_Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Make context
	ctx := (&context.Impl{}).FromRequest(r)
	// TODO: Automatically handle OPTION request
	// Resolve handler
	h, err := e.rt.Resolve(ctx.Method(), ctx.Path(), ctx)
	v := View(nil)
	// Invoke handler
	if err != nil {
		v = e.sh.OnStatus(http.StatusNotFound, ctx, nil)
	} else {
		v = e.invoke(h, ctx)
	}
	// Render view
	e.render(ctx, v, w)
}

func (e *_Engine) invoke(h _Handler, ctx Context) (v View) {
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
	v, err := h.Handle(ctx)
	if err != nil {
		v = e.sh.OnStatus(http.StatusInternalServerError, ctx, err)
	}
	return
}

func (e *_Engine) render(ctx Context, v View, w http.ResponseWriter) {
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
			if impl1, ok1 := v.(ViewSizeAdviser); ok1 {
				// Get size from view
				size = impl1.ContentLength()
			} else if impl2, ok2 := body.(bodyHasLength); ok2 {
				// Auto-detect body size
				size = int64(impl2.Len())
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

func (e *_Engine) addStartupObserver(ob StartupObserver) {
	e.startObs = append(e.startObs, ob)
}

func (e *_Engine) addShutdownObserver(ob ShutdownObserver) {
	e.shutObs = append(e.shutObs, ob)
}

// "install" installs ghost on engine
func install[G any](engine *_Engine, ghost G) {
	// Scan implemented interfaces
	var v any = ghost
	if sh, isImpl := v.(StatusHandler); isImpl {
		engine.sh = sh
	} else {
		engine.sh = defaultStatusHandler
	}
	if ob, isImpl := v.(StartupObserver); isImpl {
		engine.addStartupObserver(ob)
	}
	if ob, isImpl := v.(ShutdownObserver); isImpl {
		engine.addShutdownObserver(ob)
	}

	// Scan handle functions
	if engine.rt == nil {
		engine.rt = (&route.Registry[_Handler]{}).Init()
	}
	rt := reflect.TypeOf(ghost)
	for i := 0; i < rt.NumMethod(); i++ {
		m := rt.Method(i)
		mn, mf := m.Name, m.Func.Interface()
		// Check method function signature
		if h, ok := asHandler(mf, ghost); ok {
			if err := engine.rt.Mount(mn, h); err == nil {
				log.Printf("Mount method [%s] => %p", mn, h)
			} else {
				log.Printf("Mount method [%s] failed: %s", mn, err.Error())
			}
		}
	}
}
