package view

import (
	"github.com/deadblue/ghost"
	"net/http"
	"os"
)

func File(f *os.File) (v ghost.View, err error) {
	info, err := f.Stat()
	if err != nil {
		return
	}
	v = Generic(http.StatusOK, f).
		ContentLength64(info.Size()).
		ContentType("application/octet-stream")
	return
}
