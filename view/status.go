package view

import (
	"io"
	"net/http"
)

type statusView int

func (v statusView) Status() int {
	return int(v)
}

func (v statusView) Header() http.Header {
	return nil
}

func (v statusView) Body() io.Reader {
	return nil
}

var (
	NoContent   = statusView(http.StatusNoContent)
	NotModified = statusView(http.StatusNotModified)
)
