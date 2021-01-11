package ghost

import (
	"fmt"
	"github.com/deadblue/ghost/internal/view"
	"io"
	"log"
	"net/http"
	"reflect"
)

type _Kernel struct {
	// The ghost given by developer
	g interface{}
	// Route tree
	rt *_RouteTable

	// HTTP 404 handler
	h404 Http404Handler
	// HTTP 500 handler
	h500 Http500Handler
}

func (k *_Kernel) BeforeStartup() (err error) {
	if k.g != nil {
		if a, ok := k.g.(StartupHandler); ok {
			err = a.OnStartup()
		}
	}
	return
}

func (k *_Kernel) AfterShutdown() (err error) {
	if k.g != nil {
		if a, ok := k.g.(ShutdownHandler); ok {
			err = a.OnShutdown()
		}
	}
	return
}

// defaultView returns HTTP 200 empty view.
func (k *_Kernel) defaultView() View {
	return view.Status200
}

func (k *_Kernel) Install(ghost interface{}) *_Kernel {
	// Initial kernel
	k.g, k.rt = ghost, &_RouteTable{
		mapping:  make(map[_RouteKey]Controller),
		branches: make(map[string][]*_RoutePath),
	}
	dh := defaultStatusHandler{}
	k.h404, k.h500 = dh, dh
	if h404, ok := ghost.(Http404Handler); ok {
		k.h404 = h404
	}
	if h500, ok := ghost.(Http500Handler); ok {
		k.h500 = h500
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
	ctx := fromRequest(r)
	// Resolve controller
	ctrl, v := k.rt.Resolve(r, ctx), View(nil)
	if ctrl == nil {
		v = k.h404.OnHttp404(r.Method, r.URL.Path)
	} else {
		v = k.invoke(ctrl, ctx)
	}
	// Render view to response
	hdr := w.Header()
	hdr.Set("Server", _HeaderServer)
	if vh := v.Header(); vh != nil {
		for key, vals := range v.Header() {
			for _, val := range vals {
				hdr.Add(key, val)
			}
		}
	}
	w.WriteHeader(v.Status())
	if body := v.Body(); body != nil {
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

func (k *_Kernel) invoke(ctrl Controller, ctx Context) (v View) {
	defer func() {
		// Catch panic
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				v = k.h500.OnHttp500(err)
			} else {
				v = k.h500.OnHttp500(fmt.Errorf("panic: %v", r))
			}
		}
	}()
	v, err := ctrl(ctx)
	if err != nil {
		v = k.h500.OnHttp500(err)
	}
	if v == nil {
		v = k.defaultView()
	}
	return
}

type defaultStatusHandler struct{}

func (h defaultStatusHandler) OnHttp404(method, path string) View {
	return view.NotFound(method, path)
}

func (h defaultStatusHandler) OnHttp500(err error) View {
	return view.InternalError(err)
}
