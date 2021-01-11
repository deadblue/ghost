package view

import (
	"fmt"
	"net/http"
	"strings"
)

type Status404 struct {
	impl
}

func (v Status404) Status() int {
	return http.StatusNotFound
}

func NotFound(method, path string) Status404 {
	message := fmt.Sprintf("No Route To: %s %s",
		strings.ToUpper(method), path)
	return Status404{impl(message)}
}
