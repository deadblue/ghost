package ghost

import (
	"fmt"
	"github.com/deadblue/ghost/internal/context"
	"github.com/deadblue/ghost/internal/view"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"
)

type _Kernel struct {
	// The ghost given by developer
	g interface{}
	// Route tree
	rt *_RouteTable

	// HTTP status handler
	hsh HttpStatusHandler
	// Ghost header interceptor
	ghi HeaderInterceptor
}

func (k *_Kernel) BeforeStartup() (err error) {
	if k.g != nil {
		if h, ok := k.g.(StartupHandler); ok {
			err = h.OnStartup()
		}
	}
	return
}

func (k *_Kernel) AfterShutdown() (err error) {
	if k.g != nil {
		if h, ok := k.g.(ShutdownHandler); ok {
			err = h.OnShutdown()
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
	if sh, ok := ghost.(HttpStatusHandler); ok {
		k.hsh = sh
	} else {
		k.hsh = defaultStatusHandler
	}
	// Setup global header interceptor
	if hi, ok := ghost.(HeaderInterceptor); ok {
		k.ghi = hi
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
			// The 50x fast controller
			ctrl = binder.Bind(mt.Func.Interface())
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
		v = k.hsh.OnStatus(http.StatusNotFound, ctx, nil)
	} else {
		v = k.invoke(ctrl, ctx)
	}
	// Render view
	k.render(v, w)
}

func (k *_Kernel) invoke(ctrl Controller, ctx Context) (v View) {
	defer func() {
		// Catch panic here
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				v = k.hsh.OnStatus(http.StatusInternalServerError, ctx, err)
			} else {
				v = k.hsh.OnStatus(http.StatusInternalServerError, ctx, fmt.Errorf("panic: %v", r))
			}
		}
	}()
	v, err := ctrl(ctx)
	if err != nil {
		v = k.hsh.OnStatus(http.StatusInternalServerError, ctx, err)
	}
	return
}

func (k *_Kernel) render(v View, w http.ResponseWriter) {
	// Ensure the view is not nil
	if v == nil {
		v = k.defaultView()
	}
	// Send response header
	headers, body := w.Header(), v.Body()
	headers.Set("Server", _HeaderServer)
	// Auto detect content length
	if body != nil {
		if l, ok := body.(hasLength); ok {
			headers.Set("Content-Length", strconv.Itoa(l.Len()))
		}
	}
	// Allow view manipulates response header
	if hi, ok := v.(HeaderInterceptor); ok {
		hi.BeforeSend(headers)
	}
	// Allow ghost manipulates response header
	if k.ghi != nil {
		k.ghi.BeforeSend(headers)
	}
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

type hasLength interface {
	Len() int
}
