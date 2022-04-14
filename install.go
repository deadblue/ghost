package ghost

import (
	"github.com/deadblue/ghost/internal/route"
	"log"
	"reflect"
)

// "install" installs ghost on engine
func install[G any](engine *_Engine, ghost G) {
	// Scan implemented interfaces
	var v any = ghost
	if eh, isImpl := v.(ErrorHandler); isImpl {
		engine.eh = eh
	} else {
		engine.eh = _ErrorHandlerImpl{}
	}
	if ob, isImpl := v.(StartupObserver); isImpl {
		engine.startupObs.Append(ob)
	}
	if ob, isImpl := v.(ShutdownObserver); isImpl {
		engine.shutdownObs.Append(ob)
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

/*
W.I.P: Support multiple ghosts in one shell.

var (
	errInvalidRootPath = errors.New("root path should start with slash")
	errInvalidShell    = errors.New("do not make shell by yourself")
)

func Install[G any](shell Shell, root string, ghost G) (err error) {
	if root[0] != '/' {
		return errInvalidRootPath
	}
	if _, ok := shell.(*_ShellImpl); !ok {
		return errInvalidShell
	} else {
		// TODO
	}
	return
}
*/
