package context

import "net/url"

func (i *Impl) IsCrossOrigin() bool {
	origin := i.r.Header.Get(headerOrigin)
	if origin == "" {
		return false
	}
	ou, err := url.Parse(origin)
	if err != nil {
		return false
	}
	return ou.Scheme != i.Scheme() || ou.Host != i.Host()
}
