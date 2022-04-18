package ghost

import "github.com/deadblue/ghost/internal/view"

type _ErrorHandlerImpl struct{}

func (h _ErrorHandlerImpl) OnError(ctx Context, err error) View {
	m, p := ctx.Method(), ctx.Path()
	if err == ErrNotFound {
		return view.NotFound(m, p)
	} else {
		return view.InternalError(err)
	}
}
