package ghost

import (
	"errors"
	"github.com/deadblue/ghost/internal/route"
	"log"
	"reflect"
)

var (
	errInvalidRootPath = errors.New("root path should start with slash")
	errInvalidShell    = errors.New("do not make shell by yourself")
)

// Implant implants other ghosts into shell, and use "root" as prefix. "root"
// should not be empty and should start with "/".
//
// NOTE: This function is not goroutine-safety, DO NOT call it in multiple goroutines.
//
// Example:
//
//     shell := Born(&MasterGhost{})
//     Implant(shell, &FooGhost{}, "/foo")
//     Implant(shell, &BarGhost{}, "/bar")
//     shell.Run()
//
// Q: Why is Implant not a method on Shell?
// A: Because Go DOES NOT allow using type parameter in method :(
func Implant[G any](shell Shell, ghost G, root string) (err error) {
	if root[0] != '/' {
		return errInvalidRootPath
	}
	if shImpl, ok := shell.(*_ShellImpl); !ok {
		return errInvalidShell
	} else {
		internalImplant(shImpl.kn, ghost, root, false)
	}
	return
}

func internalImplant[G any](kernel *_Kernel, ghost G, root string, isMaster bool) {
	// Scan implemented interfaces
	var v any = ghost
	if ob, isImpl := v.(StartupObserver); isImpl {
		kernel.startupObs.Append(ob)
	}
	if ob, isImpl := v.(ShutdownObserver); isImpl {
		kernel.shutdownObs.Append(ob)
	}
	// Only master ghost can set up ErrorHandler
	if isMaster {
		if eh, isImpl := v.(ErrorHandler); isImpl {
			kernel.eh = eh
		} else {
			kernel.eh = _ErrorHandlerImpl{}
		}
	}

	rt := reflect.TypeOf(ghost)
	gn := getGhostName(rt)
	log.Printf("Implanting ghost: %s", gn)

	// Scan handle functions
	if kernel.rr == nil {
		kernel.rr = (&route.Registry[_Handler]{}).Init()
	}
	for i := 0; i < rt.NumMethod(); i++ {
		m := rt.Method(i)
		mn, mf := m.Name, m.Func.Interface()
		// Check method function signature
		if h, ok := asHandler(mf, ghost); ok {
			if rule, err := kernel.rr.Mount(mn, root, h); err == nil {
				log.Printf("Mount method [%s.%s] => [%s]", gn, mn, rule)
			} else {
				log.Printf("Mount method [%s.%s] failed: %s", gn, mn, err.Error())
			}
		}
	}
}

func getGhostName(rt reflect.Type) string {
	for rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	return rt.Name()
}
