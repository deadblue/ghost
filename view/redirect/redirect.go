package redirect

import (
	"io"
	"net/http"
)

const (
	headerLocation = "Location"
)

type MovedPermanently string

func (v MovedPermanently) Status() int {
	return http.StatusMovedPermanently
}

func (v MovedPermanently) Body() io.Reader {
	return nil
}

func (v MovedPermanently) BeforeSendHeader(h http.Header) {
	h.Set(headerLocation, string(v))
}

type Found string

func (v Found) Status() int {
	return http.StatusFound
}

func (v Found) Body() io.Reader {
	return nil
}

func (v Found) BeforeSendHeader(h http.Header) {
	h.Set(headerLocation, string(v))
}

type TemporaryRedirect string

func (v TemporaryRedirect) Status() int {
	return http.StatusTemporaryRedirect
}

func (v TemporaryRedirect) Body() io.Reader {
	return nil
}

func (v TemporaryRedirect) BeforeSendHeader(h http.Header) {
	h.Set(headerLocation, string(v))
}

type PermanentRedirect string

func (v PermanentRedirect) Status() int {
	return http.StatusPermanentRedirect
}

func (v PermanentRedirect) Body() io.Reader {
	return nil
}

func (v PermanentRedirect) BeforeSendHeader(h http.Header) {
	h.Set(headerLocation, string(v))
}
