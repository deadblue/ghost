package view

import (
	"io"
	"net/http"
)

type HeaderView struct {
	Headers map[string]string
}

func (v *HeaderView) Status() int {
	return http.StatusNoContent
}

func (v *HeaderView) Body() io.Reader {
	return nil
}

func (v *HeaderView) BeforeSendHeader(h http.Header) {
	if v.Headers == nil {
		return
	}
	for name, value := range v.Headers {
		h.Set(name, value)
	}
}
