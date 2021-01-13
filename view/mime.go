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
	return "text/plain;charset=utf-8"
}

func Json(data interface{}) (v ghost.View, err error) {
	body, err := json.Marshal(data)
	if err != nil {
		return
	}
	v = Generic(http.StatusOK, bytes.NewReader(body)).
		ContentType("application/json;charset=utf-8")
	return
}
