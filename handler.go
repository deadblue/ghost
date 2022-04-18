package ghost

import "fmt"

type _Handler interface {
	Handle(Context) (View, error)
}

type _HandlerImpl[R any] struct {
	receiver R
	methodFn func(R, Context) (View, error)
}

func (h *_HandlerImpl[R]) Handle(ctx Context) (v View, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(error); !ok {
				err = fmt.Errorf("panic: %v", r)
			}
		}
	}()
	return h.methodFn(h.receiver, ctx)
}

// asHandler checks the signature of method function, and converts it to |_Handler| if support.
func asHandler[R any](fnVal any, receiver R) (_Handler, bool) {
	if fn, ok := fnVal.(func(R, Context) (View, error)); ok {
		return &_HandlerImpl[R]{
			receiver: receiver,
			methodFn: fn,
		}, true
	}
	return nil, false
}
