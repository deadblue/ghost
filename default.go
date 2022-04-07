package ghost

import (
	"github.com/deadblue/ghost/internal/view"
	"net/http"
)

type _StatusHandlerImpl struct{}

func (h _StatusHandlerImpl) OnStatus(status int, ctx Context, err error) View {
	m, p := ctx.Method(), ctx.Path()
	switch status {
	case http.StatusNotFound:
		return view.NotFound(m, p)
	case http.StatusInternalServerError:
		return view.InternalError(err)
	default:
		return nil
	}
}

var defaultStatusHandler = _StatusHandlerImpl{}
