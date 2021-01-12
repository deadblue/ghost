package view

import (
	"fmt"
	"net/http"
	"strings"
)

type Http404 struct {
	impl
}

func (v Http404) Status() int {
	return http.StatusNotFound
}

func NotFound(method, path string) Http404 {
	message := fmt.Sprintf("No Route To: %s %s",
		strings.ToUpper(method), path)
	return Http404{impl(message)}
}
