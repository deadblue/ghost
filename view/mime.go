package view

import (
	"bytes"
	"encoding/json"
	"github.com/deadblue/ghost"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Text string

func (v Text) Status() int {
	return http.StatusOK
}

func (v Text) Header() http.Header {
	hdr := http.Header{}
	hdr.Set("Content-Type", "text/plain;charset=utf-8")
	hdr.Set("Content-Length", strconv.Itoa(len(([]byte)(v))))
	return hdr
}

func (v Text) Body() io.Reader {
	return strings.NewReader(string(v))
}

func Json(data interface{}) (v ghost.View, err error) {
	body, err := json.Marshal(data)
	if err != nil {
		return
	}
	v = Generic(http.StatusOK, bytes.NewReader(body)).
		ContentType("application/json;charset=utf-8").
		ContentLength(len(body))
	return
}
