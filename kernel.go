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
	// Route tree
	rt *_RouteTree
	// The ghost
	g interface{}
}

func (k *_Kernel) beforeStartup() (err error) {
	if k.g != nil {
		if a, ok := k.g.(AwareStartup); ok {
			err = a.OnStartup()
		}
	}
	return
}

func (k *_Kernel) afterShutdown() (err error) {
	if k.g != nil {
		if a, ok := k.g.(AwareShutdown); ok {
			err = a.OnShutdown()
		}
	}
	return
}

func (k *_Kernel) handleFallback(method, path string) View {
	// TODO: Allow developer to handle this
	return view.NotFound(method, path)
}

func (k *_Kernel) handleError(err error) View {
	// TODO: Allow developer to handle this
	return view.InternalError(err)
}

func (k *_Kernel) Install(ghost interface{}) *_Kernel {
	// Initial kernel
	k.g, k.rt = ghost, &_RouteTree{
		bs: make(map[string]*_RouteBranch),
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
		v = k.handleFallback(r.Method, r.URL.Path)
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
				v = k.handleError(err)
			} else {
				v = k.handleError(fmt.Errorf("panic: %v", r))
			}
		}
	}()
	v, err := ctrl(ctx)
	if err != nil {
		v = k.handleError(err)
	}
	// TODO: Double check the v, to make sure it is NOT nil.
	return
}
