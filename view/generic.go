package view

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type GenericView struct {
	status int
	header http.Header
	body   io.Reader
}

func (gv *GenericView) Status() int {
	return gv.status
}

func (gv *GenericView) Body() io.Reader {
	return gv.body
}

func (gv *GenericView) BeforeSendHeader(h http.Header) {
	// Pouring stored headers
	if gv.header != nil && len(gv.header) > 0 {
		for k, vs := range gv.header {
			for _, v := range vs {
				h.Add(k, v)
			}
		}
	}
}

func (gv *GenericView) BodySize(size int64) *GenericView {
	gv.header.Set("Content-Length", strconv.FormatInt(size, 10))
	return gv
}

func (gv *GenericView) MediaType(t string) *GenericView {
	gv.header.Set("Content-Type", t)
	return gv
}

func (gv *GenericView) Download(mode string, filename string) *GenericView {
	value := "inline"
	if mode == "attachment" {
		value = "attachment"
	}
	if filename != "" {
		value += "; filename=\"" + filename + "\""
	}
	gv.header.Set("Content-Disposition", value)
	return gv
}

func (gv *GenericView) CachePrivate(age time.Duration) *GenericView {
	gv.header.Set("Cache-Control", fmt.Sprintf(
		"private, max-age=%d", int64(age.Seconds())))
	return gv
}

func (gv *GenericView) CachePublic(age time.Duration) *GenericView {
	gv.header.Set("Cache-Control", fmt.Sprintf(
		"public, max-age=%d", int64(age.Seconds())))
	return gv
}

func (gv *GenericView) CacheDisable() *GenericView {
	gv.header.Set("Cache-Control", "no-store")
	return gv
}

func (gv *GenericView) AddHeader(name, value string) *GenericView {
	gv.header.Add(name, value)
	return gv
}

// Generic returns a GenericView, which has a lot of methods to setup the response headers.
func Generic(status int, body io.Reader) *GenericView {
	return &GenericView{
		status: status,
		body:   body,
		header: http.Header{},
	}
}
