package ghost

import (
	"github.com/deadblue/ghost/internal/view"
	"net/http"
)

type implStatusHandler struct{}

func (h implStatusHandler) OnStatus(status int, context Context, err error) View {
	m, p := context.MethodAndPath()
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
