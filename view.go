package ghost

import "net/http"

// HeaderInterceptor is an interface that can optionally implement by View. It will be called
// after kernel sets normal response headers, developer can manipulate response header here.
// TODO: This interface need to be re-defined.
type HeaderInterceptor interface {

	// BeforeSend will be called before kernel send response headers to client.
	BeforeSend(h http.Header)
}

// ViewTypeAdviser is an optional interface for View, when a view implements it, kernel will get
// content type from it, and set to response header as "Content-Type".
type ViewTypeAdviser interface {
	ContentType() string
}

// ViewSizeAdviser is an optional interface for View, when a view implements it, kernel will get
// content size from it, and set to response header as "Content-Length".
type ViewSizeAdviser interface {
	ContentLength() int64
}
