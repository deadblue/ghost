package context

import (
	"io"
	"net/http"
	"strings"
)

// Impl is the internal implementation of |ghost.Context|.
// It is not goroutine-safe, so DO NOT use it in multi-goroutine.
type Impl struct {
	// Underlying HTTP request
	r *http.Request

	// Path variables
	pv map[string]string
}

func (i *Impl) FromRequest(r *http.Request) *Impl {
	i.r = r
	i.pv = make(map[string]string)
	return i
}

func (i *Impl) Request() *http.Request {
	return i.r
}

func (i *Impl) Method() string {
	return i.r.Method
}

func (i *Impl) Path() string {
	return i.r.URL.Path
}

func (i *Impl) Body() io.Reader {
	return i.r.Body
}

func (i *Impl) Header(name string) (value string, found bool) {
	value = i.r.Header.Get(name)
	found = value != ""
	return
}

func (i *Impl) Scheme() string {
	scheme := i.r.Header.Get(headerForwardedProto)
	if scheme == "" {
		scheme = i.r.URL.Scheme
	}
	return scheme
}

func (i *Impl) Host() string {
	host := i.r.Header.Get(headerForwardedHost)
	if host == "" {
		host = i.r.Host
	}
	return host
}

func (i *Impl) RemoteIp() string {
	// Try to get IP from XFF header
	xff := i.r.Header.Get(headerForwardedFor)
	if xff != "" {
		if ips := strings.Split(xff, ","); len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	addr := i.r.RemoteAddr
	for n := len(addr) - 1; n >= 0; n-- {
		if addr[n] == ':' {
			return addr[:n]
		}
	}
	return addr
}

func (i *Impl) PathVar(name string) string {
	return i.pv[name]
}

func (i *Impl) SetPathVar(name, value string) {
	i.pv[name] = value
}
