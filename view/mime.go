package view

import (
	"bytes"
	"encoding/json"
	"github.com/deadblue/ghost"
	"io"
	"net/http"
	"strings"
)

type Text string

func (v Text) Status() int {
	return http.StatusOK
}

func (v Text) Body() io.Reader {
	return strings.NewReader(string(v))
}

func (v Text) ContentType() string {
	return "text/plain; charset=utf-8"
}

type Binary []byte

func (v Binary) Status() int {
	return http.StatusOK
}

func (v Binary) Body() io.Reader {
	return bytes.NewReader(v)
}

func (v Binary) ContentType() string {
	return "application/octet-stream"
}

type Json []byte

func (v Json) Status() int {
	return http.StatusOK
}

func (v Json) Body() io.Reader {
	return bytes.NewReader(v)
}

func (v Json) ContentType() string {
	return "application/json; charset=utf-8"
}

func AsJson(v interface{}) (view ghost.View, err error) {
	body, err := json.Marshal(v)
	if err == nil {
		view = Json(body)
	}
	return
}
