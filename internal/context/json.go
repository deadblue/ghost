package context

import (
	"encoding/json"
	"errors"
	"mime"
)

var (
	errNotJson = errors.New("request content is not application/json")
)

func (i *Impl) Json(v interface{}) (err error) {
	// Check Content-Type
	mt, _, err := mime.ParseMediaType(i.r.Header.Get("Content-Type"))
	if err != nil || mt != "application/json" {
		return errNotJson
	}
	// Parse json
	// TODO: Allow developer using a customized JSON parser.
	err = json.NewDecoder(i.Body()).Decode(v)
	return
}
