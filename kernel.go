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

type _Kernel struct {
	// The ghost given by developer
	g interface{}
	// Route tree
	rt *_RouteTable

	// HTTP status observer
	sh StatusHandler
}

func (k *_Kernel) BeforeStartup() (err error) {
	if k.g != nil {
		if o, isImpl := k.g.(StartupObserver); isImpl {
			err = o.BeforeStartup()
		}
	}
	return
}

func (k *_Kernel) AfterShutdown() (err error) {
	if k.g != nil {
		if o, isImpl := k.g.(ShutdownObserver); isImpl {
			err = o.AfterShutdown()
		}
	}
	return
}

func (k *_Kernel) Install(ghost interface{}) *_Kernel {
	// Initial kernel
	k.g, k.rt = ghost, &_RouteTable{
		mapping:  make(map[_RouteKey]Controller),
		branches: make(map[string][]*_RoutePath),
	}
	// Setup HTTP status handler
	if sh, ok := ghost.(StatusHandler); ok {
		k.sh = sh
	} else {
		k.sh = defaultStatusHandler
	}
	// Scan controller
	binder, hasBinder := ghost.(Binder)
	rt, rv := reflect.TypeOf(ghost), reflect.ValueOf(ghost)
	for i := 0; i < rt.NumMethod(); i++ {
		mt, mv := rt.Method(i), rv.Method(i)
		// Try to convert the method function to controller.
		ctrl, isCtrl := mv.Interface().(func(Context) (View, error))
		if !isCtrl {
			continue
		}
		// Use binder if developer implements it.
		if hasBinder {
			// Double check, avoid nil controller.
			if fastCtrl := binder.Bind(mt.Func.Interface()); fastCtrl != nil {
				ctrl = fastCtrl
			}
		}
		// Mount controller
		if err := k.rt.Mount(mt.Name, ctrl); err == nil {
			log.Printf("Mount controller: %s", mt.Name)
		}
	}
	return k
}

func (k *_Kernel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (k *_Kernel) invoke(ctrl Controller, ctx Context) (v View) {
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

func (k *_Kernel) render(ctx Context, v View, w http.ResponseWriter) {
	// Ensure the view is not nil
	if v == nil {
		v = k.defaultView()
	}
	headers, body := w.Header(), v.Body()
	// Prepare response header
	headers.Set(_HeaderServer, _Server)
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
				// Auto detect body size
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
		// Copy data to response writer
		if _, err := io.Copy(w, body); err != nil {
			log.Printf("Render view error: %s", err)
		}
	}
}

// defaultView returns HTTP 200 empty view.
func (k *_Kernel) defaultView() View {
	return view.Http200
}
