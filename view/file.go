package view

import (
	"github.com/deadblue/ghost"
	"mime"
	"net/http"
	"os"
	"path"
)

func FromFile(f *os.File) (v ghost.View, err error) {
	// Get file size
	info, err := f.Stat()
	if err != nil {
		return
	}
	// Determine media type from file extension
	mediaType := mime.TypeByExtension(path.Ext(f.Name()))
	if mediaType == "" {
		mediaType = "application/octet-stream"
	}
	// Build view
	v = Generic(http.StatusOK, f).
		BodySize(info.Size()).MediaType(mediaType)
	return
}
