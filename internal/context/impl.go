package context

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// TODO:
//   Add more useful methods on Context.
//   The goal is to make developer forget HTTP request ...

// Impl is the internal ghost.Context implementation.
// It is not goroutine-safe, so DO NOT use it in multi-goroutine.
type Impl struct {
	r *http.Request

	// Path variables
	pv map[string]string
	// Query-string values
	qs url.Values
}

func (i *Impl) Request() *http.Request {
	return i.r
}

func (i *Impl) PathVar(name string) string {
	return i.pv[name]
}

func (i *Impl) Header(name string) string {
	return i.r.Header.Get(name)
}

func (i *Impl) Query(name string) string {
	if i.qs == nil {
		i.qs = i.r.URL.Query()
	}
	return i.qs.Get(name)
}

func (i *Impl) Json(v interface{}) (err error) {
	if ct := i.r.Header.Get("content-type"); strings.HasPrefix(ct, "application/json") {
		err = json.NewDecoder(i.r.Body).Decode(v)
	} else {
		err = errors.New("malformed request")
	}
	return
}

func (i *Impl) Body() io.ReadCloser {
	return i.r.Body
}

func FromRequest(r *http.Request) *Impl {
	return &Impl{
		r: r,
	}
}
