package rhttp

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)

// ReaderFunc is the type of function that can be given natively to NewRequest
type ReaderFunc func() (io.Reader, error)

type Bodyable interface {
	SetBodyReader(rawBody ReaderFunc)
	SetContentLength(int64)
	BodyReader() (io.Reader, error)
}

// BodyBytes allows accessing the request body. It is an analogue to
// http.Request's Body variable, but it returns a copy of the underlying data
// rather than consuming it.
//
// This function is not thread-safe; do not call it at the same time as another
// call, or at the same time this request is being used with Client.Do.
func BodyBytes(r Bodyable) ([]byte, error) {
	body, err := r.BodyReader()
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(body)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// SetBody allows setting the request body.
//
// It is useful if a new body needs to be set without constructing a new Request.
func SetBody(bdy Bodyable, rawBody interface{}) error {
	bodyReader, contentLength, err := getBodyReaderAndContentLength(rawBody)
	if err != nil {
		return err
	}
	bdy.SetBodyReader(bodyReader)
	bdy.SetContentLength(contentLength)
	return nil
}

// WriteTo allows copying the request body into a writer.
//
// It writes data to w until there's no more data to write or
// when an error occurs. The return int64 value is the number of bytes
// written. Any error encountered during the write is also returned.
// The signature matches io.WriterTo interface.
func WriteTo(r Bodyable, w io.Writer) (int64, error) {
	body, err := r.BodyReader()
	if err != nil {
		return 0, err
	}
	if c, ok := body.(io.Closer); ok {
		defer c.Close()
	}
	return io.Copy(w, body)
}

func getBodyReaderAndContentLength(rawBody interface{}) (ReaderFunc, int64, error) {
	var bodyReader ReaderFunc
	var contentLength int64

	switch body := rawBody.(type) {
	// If they gave us a function already, great! Use it.
	case ReaderFunc:
		bodyReader = body
		tmp, err := body()
		if err != nil {
			return nil, 0, err
		}
		if lr, ok := tmp.(LenReader); ok {
			contentLength = int64(lr.Len())
		}
		if c, ok := tmp.(io.Closer); ok {
			c.Close()
		}

	case func() (io.Reader, error):
		bodyReader = body
		tmp, err := body()
		if err != nil {
			return nil, 0, err
		}
		if lr, ok := tmp.(LenReader); ok {
			contentLength = int64(lr.Len())
		}
		if c, ok := tmp.(io.Closer); ok {
			c.Close()
		}

	// If a regular byte slice, we can read it over and over via new
	// readers
	case []byte:
		buf := body
		bodyReader = func() (io.Reader, error) {
			return bytes.NewReader(buf), nil
		}
		contentLength = int64(len(buf))

	// If a bytes.Buffer we can read the underlying byte slice over and
	// over
	case *bytes.Buffer:
		buf := body
		bodyReader = func() (io.Reader, error) {
			return bytes.NewReader(buf.Bytes()), nil
		}
		contentLength = int64(buf.Len())

	// We prioritize *bytes.Reader here because we don't really want to
	// deal with it seeking so want it to match here instead of the
	// io.ReadSeeker case.
	case *bytes.Reader:
		buf, err := ioutil.ReadAll(body)
		if err != nil {
			return nil, 0, err
		}
		bodyReader = func() (io.Reader, error) {
			return bytes.NewReader(buf), nil
		}
		contentLength = int64(len(buf))

	// Compat case
	case io.ReadSeeker:
		raw := body
		bodyReader = func() (io.Reader, error) {
			_, err := raw.Seek(0, 0)
			return ioutil.NopCloser(raw), err
		}
		if lr, ok := raw.(LenReader); ok {
			contentLength = int64(lr.Len())
		}

	// Read all in so we can reset
	case io.Reader:
		buf, err := ioutil.ReadAll(body)
		if err != nil {
			return nil, 0, err
		}
		bodyReader = func() (io.Reader, error) {
			return bytes.NewReader(buf), nil
		}
		contentLength = int64(len(buf))

	// No body provided, nothing to do
	case nil:

	// Unrecognized type
	default:
		return nil, 0, fmt.Errorf("cannot handle type %T", rawBody)
	}
	return bodyReader, contentLength, nil
}
