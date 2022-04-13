package view

import (
	"github.com/deadblue/ghost"
	"io"
	"net/http"
)

type statusView int

func (v statusView) Status() int {
	return int(v)
}

func (v statusView) Body() io.Reader {
	return nil
}

var (
	NoContent   ghost.View = statusView(http.StatusNoContent)
	NotModified ghost.View = statusView(http.StatusNotModified)
)
