package context

import "errors"

var (
	errNotJson = errors.New("request content is not application/json")
)
