package view

import (
	"io"
	"net/http"
	"strconv"
	"strings"
)

type impl string

func (i impl) Status() int {
	return http.StatusOK
}

func (i impl) Header() http.Header {
	hdr := http.Header{}
	hdr.Set("Content-Type", "text/plain;charset=utf-8")
	hdr.Set("Content-Length", strconv.Itoa(len(([]byte)(i))))
	return hdr
}

func (i impl) Body() io.Reader {
	return strings.NewReader(string(i))
}
