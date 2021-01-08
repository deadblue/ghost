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

func (v *GenericView) Status() int {
	return v.status
}

func (v *GenericView) Header() http.Header {
	return v.header
}

func (v *GenericView) Body() io.Reader {
	return v.body
}

func (v *GenericView) ContentType(mimeType string) *GenericView {
	v.header.Set("Content-Type", mimeType)
	return v
}

func (v *GenericView) ContentLength(length int) *GenericView {
	v.header.Set("Content-Length", strconv.Itoa(length))
	return v
}

func (v *GenericView) ContentLength64(length int64) *GenericView {
	v.header.Set("Content-Length", strconv.FormatInt(length, 10))
	return v
}

func (v *GenericView) PrivateCache(age time.Duration) *GenericView {
	v.header.Set("Cache-Control", fmt.Sprintf(
		"private, max-age=%d", int64(age.Seconds())))
	return v
}

func (v *GenericView) PublicCache(age time.Duration) *GenericView {
	v.header.Set("Cache-Control", fmt.Sprintf(
		"public, max-age=%d", int64(age.Seconds())))
	return v
}

func (v *GenericView) DisableCache() *GenericView {
	v.header.Set("Cache-Control", "no-store")
	return v
}

func (v *GenericView) AddHeader(name, value string) *GenericView {
	v.header.Add(name, value)
	return v
}

// Generic returns a GenericView, which has a lot of methods to setup the response headers.
func Generic(status int, body io.Reader) *GenericView {
	return &GenericView{
		status: status,
		body:   body,
		header: http.Header{},
	}
}
