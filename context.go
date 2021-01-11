package ghost

import (
	"net/http"
	"net/url"
)

type _ContextImpl struct {
	// Original request
	r *http.Request

	// Path variables
	pv map[string]string
	// Query-string values
	qs url.Values
}

func (c *_ContextImpl) Request() *http.Request {
	return c.r
}

func (c *_ContextImpl) PathVar(name string) string {
	return c.pv[name]
}

func (c *_ContextImpl) Query(name string) string {
	return c.qs.Get(name)
}

// TODO:
//  Currently, I have no idea what methods should Context provide.
//  I will add more methods in future when I think it need.
//  Anyway, you can call Request() to retrieve the original request, and
//  get information from it.

// fromRequest makes a Context from http request.
func fromRequest(r *http.Request) *_ContextImpl {
	return &_ContextImpl{
		r:  r,
		pv: make(map[string]string),
		qs: r.URL.Query(),
	}
}
