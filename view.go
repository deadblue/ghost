package ghost

import "net/http"

// ViewTypeAdviser is an optional interface for View.
// Developer need not implement it when the view does not have a body.
type ViewTypeAdviser interface {

	// ContentType returns the content type of the view, it will be set in response header
	// as "Content-Type".
	ContentType() string
}

// ViewSizeAdviser is an optional interface for View.
// Developer need not implement it when the view does not have a body, or the body is one
// of "bytes.Buffer", "bytes.Reader", "strings.Reader".
type ViewSizeAdviser interface {

	// ContentLength returns the body size of the view, it will be set in response header
	// as "Content-Length", DO NOT return a incorrect value which is less or more than the
	// body size, that may cause some strange issues.
	ContentLength() int64
}

// ViewHeaderInterceptor is an optional interface for View. When a view implements it, kernel
// will pass response header to the view before send to client, view can manipulate the response
// header here.
type ViewHeaderInterceptor interface {

	// BeforeSendHeader will be called before kernel send the response headers to client.
	// View can add/update/remove any headers in it.
	BeforeSendHeader(h http.Header)
}

type bodyHasLength interface {
	Len() int
}
