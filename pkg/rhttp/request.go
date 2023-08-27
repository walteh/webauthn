package rhttp

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

var _ Bodyable = (*Request)(nil)

// Request wraps the metadata needed to create HTTP requests.
type Request struct {
	// body is a seekable reader over the request body payload. This is
	// used to rewind the request data in between retries.
	body ReaderFunc

	responseHandler ResponseHandlerFunc

	// Embed an HTTP request directly. This makes a *Request act exactly
	// like an *http.Request so that all meta methods are supported.
	*http.Request
}

// WithContext returns wrapped Request with a shallow copy of underlying *http.Request
// with its context changed to ctx. The provided ctx must be non-nil.
func (r *Request) WithContext(ctx context.Context) *Request {
	return &Request{
		body:            r.body,
		responseHandler: r.responseHandler,
		Request:         r.Request.WithContext(ctx),
	}
}

// SetResponseHandler allows setting the response handler.
func (r *Request) SetResponseHandler(fn ResponseHandlerFunc) {
	r.responseHandler = fn
}

func (r *Request) SetBodyReader(rawBody ReaderFunc) {
	r.body = rawBody
}

func (r *Request) SetContentLength(o int64) {
	r.ContentLength = o
}

// func (r *Request) WriteTo(w io.Writer) (int64, error) {
// 	return WriteTo(r, w)

// BodyBytes allows accessing the request body. It is an analogue to
// http.Request's Body variable, but it returns a copy of the underlying data
// rather than consuming it.
//
// This function is not thread-safe; do not call it at the same time as another
// call, or at the same time this request is being used with Client.Do.
// func (r *Request) BodyBytes() ([]byte, error) {
// 	BodyBytes(r)
// }

func (r *Request) BodyReader() (io.Reader, error) {
	if r.body == nil {
		return nil, fmt.Errorf("no body set")
	}

	return r.body()
}

// FromRequest wraps an http.Request in a retryablehttp.Request
func FromRequest(r *http.Request) (*Request, error) {
	bodyReader, _, err := getBodyReaderAndContentLength(r.Body)
	if err != nil {
		return nil, err
	}
	// Could assert contentLength == r.ContentLength
	return &Request{body: bodyReader, Request: r}, nil
}

// NewRequest creates a new wrapped request.
// func NewRequest(method, url string, rawBody interface{}) (*Request, error) {
// 	return NewRequestWithContext(context.Background(), method, url, rawBody)
// }

// NewRequestWithContext creates a new wrapped request with the provided context.
//
// The context controls the entire lifetime of a request and its response:
// obtaining a connection, sending the request, and reading the response headers and body.
func NewRequestWithContext(ctx context.Context, method, url string, rawBody interface{}) (*Request, error) {
	bodyReader, contentLength, err := getBodyReaderAndContentLength(rawBody)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.ContentLength = contentLength

	return &Request{body: bodyReader, Request: httpReq}, nil
}
