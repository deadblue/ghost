package view

import (
	"github.com/deadblue/ghost"
	"net/http"
)

func redirectView(status int, url string) ghost.View {
	return Generic(status, nil).
		AddHeader("Location", url)
}

func MovedPermanently(url string) ghost.View {
	return redirectView(http.StatusMovedPermanently, url)
}

func Found(url string) ghost.View {
	return redirectView(http.StatusFound, url)
}

func TemporaryRedirect(url string) ghost.View {
	return redirectView(http.StatusTemporaryRedirect, url)
}

func PermanentRedirect(url string) ghost.View {
	return redirectView(http.StatusPermanentRedirect, url)
}
