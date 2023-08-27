package rhttp

import (
	"fmt"
	"io"
	"net/http"
)

var _ Bodyable = (*Response)(nil)

// Request wraps the metadata needed to create HTTP requests.
type Response struct {
	// body is a seekable reader over the request body payload. This is
	// used to rewind the request data in between retries.
	body ReaderFunc

	// Embed an HTTP request directly. This makes a *Response act exactly
	// like an *http.Response so that all meta methods are supported.
	*http.Response
}

// SetBody allows setting the request body.
func (r *Response) SetBodyReader(rawBody ReaderFunc) {
	r.body = rawBody
}

func (r *Response) SetContentLength(o int64) {
	r.ContentLength = o
}

func (r *Response) BodyReader() (io.Reader, error) {
	if r.body == nil {
		return nil, fmt.Errorf("no body set")
	}
	return r.body()
}

func NewResponse(resp *http.Response) (*Response, error) {
	reader, contentLength, err := getBodyReaderAndContentLength(resp.Body)
	if err != nil {
		return nil, err
	}

	a, err := reader()
	if err != nil {
		return nil, err
	}

	resp.Body = io.NopCloser(a)

	resp.ContentLength = contentLength

	return &Response{
		Response: resp,
		body:     reader,
	}, nil
}
