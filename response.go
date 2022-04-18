package ghost

import (
	"net/http"
)

type _HeadResponseWriter struct {
	rw http.ResponseWriter
}

func (w *_HeadResponseWriter) Header() http.Header {
	return w.rw.Header()
}

func (w *_HeadResponseWriter) WriteHeader(statusCode int) {
	w.rw.WriteHeader(statusCode)
}

func (w *_HeadResponseWriter) Write(bytes []byte) (int, error) {
	return len(bytes), nil
}
