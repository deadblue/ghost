package view

import (
	"fmt"
	"net/http"
	"strings"
)

type View404 struct {
	impl
}

func (v View404) Status() int {
	return http.StatusNotFound
}

func NotFound(method, path string) View404 {
	message := fmt.Sprintf("No Route To: %s %s",
		strings.ToUpper(method), path)
	return View404{impl(message)}
}
