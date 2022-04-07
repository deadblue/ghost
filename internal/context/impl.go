package context

import (
	"io"
	"net/http"
	"strings"
)

// Impl is the internal ghost.Context implementation.
// It is not goroutine-safe, so DO NOT use it in multi-goroutine.
type Impl struct {
	r *http.Request

	// Path variables
	pv map[string]string
}

const (
	headerForwardedFor   = "X-Forwarded-For"
	headerForwardedProto = "X-Forwarded-Proto"
	headerForwardedHost  = "X-Forwarded-Host"

	headerCfConnectingIp = "Cf-Connecting-Ip"
)

func (i *Impl) FromRequest(r *http.Request) *Impl {
	i.r = r
	return i
}

func (i *Impl) Request() *http.Request {
	return i.r
}

func (i *Impl) Method() string {
	return i.r.Method
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

func (i *Impl) Path() string {
	return i.r.URL.Path
}

func (i *Impl) BaseName() string {
	return ""
}

func (i *Impl) Body() io.Reader {
	return i.r.Body
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
