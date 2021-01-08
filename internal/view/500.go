package view

import (
	"fmt"
	"net/http"
)

type View500 struct {
	impl
}

func (v View500) Status() int {
	return http.StatusInternalServerError
}

func InternalError(err error) View500 {
	message := fmt.Sprintf("Internal Error: %s", err.Error())
	return View500{impl(message)}
}
