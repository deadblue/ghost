package view

import (
	"fmt"
	"net/http"
)

type Http500 struct {
	impl
}

func (v Http500) Status() int {
	return http.StatusInternalServerError
}

func InternalError(err error) Http500 {
	message := fmt.Sprintf("Internal Error: %s", err.Error())
	return Http500{impl(message)}
}
