package context

import (
	"io"
	"net/http"
	"net/url"
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
	// Cookie values
	ck map[string][]string
}

func (i *Impl) Request() *http.Request {
	return i.r
}

func (i *Impl) MethodAndPath() (method string, path string) {
	return i.r.Method, i.r.URL.Path
}

func (i *Impl) PathVar(name string) string {
	return i.pv[name]
}

func (i *Impl) HeaderArray(name string) []string {
	return i.r.Header[name]
}

func (i *Impl) Header(name string) string {
	return i.r.Header.Get(name)
}

func (i *Impl) QueryArray(name string) []string {
	if i.qs == nil {
		i.qs = i.r.URL.Query()
	}
	return i.qs[name]
}

func (i *Impl) Query(name string) string {
	a := i.QueryArray(name)
	if a == nil || len(a) == 0 {
		return ""
	} else {
		return a[0]
	}
}

func (i *Impl) CookieArray(name string) []string {
	if i.ck == nil {
		i.ck = make(map[string][]string)
		// TODO: Parse cookies by myself.
		for _, c := range i.r.Cookies() {
			k, v := c.Name, c.Value
			i.ck[k] = append(i.ck[k], v)
		}
	}
	return i.ck[name]
}

func (i *Impl) Cookie(name string) string {
	a := i.CookieArray(name)
	if a == nil || len(a) == 0 {
		return ""
	} else {
		return a[0]
	}
}

func (i *Impl) Body() io.Reader {
	return i.r.Body
}

func FromRequest(r *http.Request) *Impl {
	return &Impl{
		r: r,
	}
}
