package view

import (
	"fmt"
	"net/http"
)

type Status500 struct {
	impl
}

func (v Status500) Status() int {
	return http.StatusInternalServerError
}

func InternalError(err error) Status500 {
	message := fmt.Sprintf("Internal Error: %s", err.Error())
	return Status500{impl(message)}
}
