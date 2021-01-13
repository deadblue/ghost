package view

import (
	"io"
	"net/http"
	"strings"
)

type impl string

func (i impl) Status() int {
	return http.StatusOK
}

func (i impl) Body() io.Reader {
	return strings.NewReader(string(i))
}

func (i impl) ContentType() string {
	return "text/plain;charset=utf-8"
}
