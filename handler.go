package ghost

import (
	"github.com/deadblue/ghost/internal/view"
	"net/http"
)

type implStatusHandler struct{}

func (h implStatusHandler) OnStatus(status int, ctx Context, err error) View {
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

var defaultStatusHandler = implStatusHandler{}
