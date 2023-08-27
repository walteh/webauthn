//	rhttp.go
//	golang/rhttp
//
//	Created by walter on 2023-01-08.
//	Copyright © 2023, nugg.xyz LLC. All rights reserved.
//	---------------------------------------------------------------------
//	adapted from hashicorp/go-retryablehttp
//	Copyright © 2015, HashiCorp, Inc. MPL 2.0
//	---------------------------------------------------------------------
//

package rhttp

// Package retryablehttp provides a familiar HTTP client interface with
// automatic retries and exponential backoff. It is a thin wrapper over the
// standard net/http client library and exposes nearly the same public API.
// This makes retryablehttp very easy to drop into existing programs.
//
// retryablehttp performs automatic retries under certain conditions. Mainly, if
// an error is returned by the client (connection errors etc), or if a 500-range
// response is received, then a retry is invoked. Otherwise, the response is
// returned and left to the caller to interpret.
//
// Requests which take a request body should provide a non-nil function
// parameter. The best choice is to provide either a function satisfying
// ReaderFunc which provides multiple io.Readers in an efficient manner, a
// *bytes.Buffer (the underlying raw byte slice will be used) or a raw byte
// slice. As it is a reference type, and we will wrap it as needed by readers,
// we can efficiently re-use the request body without needing to copy it. If an
// io.Reader (such as a *bytes.Reader) is provided, the full body will be read
// prior to the first request, and will be efficiently re-used for any retries.
// ReadSeeker can be used, but some users have observed occasional data races
// between the net/http library and the Seek functionality of some
// implementations of ReadSeeker, so should be avoided if possible.

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/smithy-go/encoding/xml"
	"github.com/rs/zerolog"
)

var (
	// Default retry configuration
	defaultRetryWaitMin = 1 * time.Second
	defaultRetryWaitMax = 30 * time.Second
	defaultRetryMax     = 4

	// defaultLogger is the logger provided with defaultClient
	defaultLogger = zerolog.New(os.Stderr).With().Timestamp().Logger()

	// defaultClient is used for performing requests without explicitly making
	// a new client. It is purposely private to avoid modifications.
	defaultClient = NewClient()

	// We need to consume response bodies to maintain http connections, but
	// limit the size we consume to respReadLimit.
	respReadLimit = int64(4096)

	// A regular expression to match the error returned by net/http when the
	// configured number of redirects is exhausted. This error isn't typed
	// specifically so we resort to matching on the error string.
	redirectsErrorRe = regexp.MustCompile(`stopped after \d+ redirects\z`)

	// A regular expression to match the error returned by net/http when the
	// scheme specified in the URL is invalid. This error isn't typed
	// specifically so we resort to matching on the error string.
	schemeErrorRe = regexp.MustCompile(`unsupported protocol scheme`)

	// A regular expression to match the error returned by net/http when the
	// TLS certificate is not trusted. This error isn't typed
	// specifically so we resort to matching on the error string.
	notTrustedErrorRe = regexp.MustCompile(`certificate is not trusted`)
)

// ResponseHandlerFunc is a type of function that takes in a Response, and does something with it.
// The ResponseHandlerFunc is called when the HTTP client successfully receives a response and the
// CheckRetry function indicates that a retry of the base request is not necessary.
// If an error is returned from this function, the CheckRetry policy will be used to determine
// whether to retry the whole request (including this handler).
//
// Make sure to check status codes! Even if the request was completed it may have a non-2xx status code.
//
// The response body is not automatically closed. It must be closed either by the ResponseHandlerFunc or
// by the caller out-of-band. Failure to do so will result in a memory leak.
type ResponseHandlerFunc func(*Response) error

// LenReader is an interface implemented by many in-memory io.Reader's. Used
// for automatically sending the right Content-Length header when possible.
type LenReader interface {
	Len() int
}

type ctxKey string

const defaultContextKey ctxKey = "rhttp"

func Ctx(ctx context.Context) *Client {
	if ctx == nil {
		return defaultClient
	}

	if c, ok := ctx.Value(defaultContextKey).(*Client); ok {
		return c
	}

	return defaultClient
}

func (me *Client) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, defaultContextKey, me)
}

func (me *Client) WithAwsConfig(cfg aws.Config) aws.Config {
	cfg.HTTPClient = &http.Client{
		Transport: &RoundTripper{Client: me},
	}
	return cfg
}

// // Logger interface allows to use other loggers than
// // standard log.Logger.
// type Logger interface {
// 	Printf(string, ...interface{})
// }

// LeveledLogger is an interface that can be implemented by any logger or a
// logger wrapper to provide leveled logging. The methods accept a message
// string and a variadic number of key-value pairs. For log.Printf style
// formatting where message string contains a format specifier, use Logger
// interface.
// type LeveledLogger interface {
// 	Error(msg string, keysAndValues ...interface{})
// 	Info(msg string, keysAndValues ...interface{})
// 	Debug(msg string, keysAndValues ...interface{})
// 	Warn(msg string, keysAndValues ...interface{})
// }

// // hookLogger adapts an LeveledLogger to Logger for use by the existing hook functions
// // without changing the API.
// type hookLogger struct {
// 	LeveledLogger
// }

// func (h hookLogger) Printf(s string, args ...interface{}) {
// 	h.Info(fmt.Sprintf(s, args...))
// }

// RequestLogHook allows a function to run before each retry. The HTTP
// request which will be made, and the retry number (0 for the initial
// request) are available to users. The internal logger is exposed to
// consumers.
type RequestLogHook func(*zerolog.Logger, *http.Request, int)

// ResponseLogHook is like RequestLogHook, but allows running a function
// on each HTTP response. This function will be invoked at the end of
// every HTTP request executed, regardless of whether a subsequent retry
// needs to be performed or not. If the response body is read or closed
// from this method, this will affect the response returned from Do().
type ResponseLogHook func(*zerolog.Logger, *Response)

// CheckRetry specifies a policy for handling retries. It is called
// following each request with the response and error values returned by
// the http.Client. If CheckRetry returns false, the Client stops retrying
// and returns the response to the caller. If CheckRetry returns an error,
// that error value is returned in lieu of the error from the request. The
// Client will close any response body when retrying, but if the retry is
// aborted it is up to the CheckRetry callback to properly close any
// response body before returning.
type CheckRetry func(ctx context.Context, resp *Response, err error) (bool, error)

// Backoff specifies a policy for how long to wait between retries.
// It is called after a failing request to determine the amount of time
// that should pass before trying again.
type Backoff func(min, max time.Duration, attemptNum int, resp *Response) time.Duration

// ErrorHandler is called if retries are expired, containing the last status
// from the http library. If not specified, default behavior for the library is
// to close the body and return an error indicating how many tries were
// attempted. If overriding this, be sure to close the body if needed.
type ErrorHandler func(resp *Response, err error, numTries int) (*Response, error)

// Client is used to make HTTP requests. It adds additional functionality
// like automatic retries to tolerate minor outages.
type Client struct {
	HTTPClient *http.Client   // Internal HTTP client.
	Logger     zerolog.Logger // Customer logger instance. Can be either Logger or LeveledLogger

	RetryWaitMin time.Duration // Minimum time to wait
	RetryWaitMax time.Duration // Maximum time to wait
	RetryMax     int           // Maximum number of retries

	// RequestLogHook allows a user-supplied function to be called
	// before each retry.
	RequestLogHook RequestLogHook

	// ResponseLogHook allows a user-supplied function to be called
	// with the response from each HTTP request executed.
	ResponseLogHook ResponseLogHook

	// CheckRetry specifies the policy for handling retries, and is called
	// after each request. The default policy is DefaultRetryPolicy.
	CheckRetry CheckRetry

	// Backoff specifies the policy for how long to wait between retries
	Backoff Backoff

	// ErrorHandler specifies the custom error handler to use, if any
	ErrorHandler ErrorHandler

	loggerInit sync.Once
	clientInit sync.Once
}

// NewClient creates a new Client with default settings.
func NewClient() *Client {
	return &Client{
		HTTPClient:   DefaultPooledClient(),
		Logger:       defaultLogger,
		RetryWaitMin: defaultRetryWaitMin,
		RetryWaitMax: defaultRetryWaitMax,
		RetryMax:     defaultRetryMax,
		CheckRetry:   DefaultRetryPolicy,
		Backoff:      DefaultBackoff,
	}
}

// NewClient creates a new Client with default settings.
func NewClientContext(ctx context.Context) *Client {
	return &Client{
		HTTPClient:   DefaultPooledClient(),
		Logger:       *zerolog.Ctx(ctx),
		RetryWaitMin: defaultRetryWaitMin,
		RetryWaitMax: defaultRetryWaitMax,
		RetryMax:     defaultRetryMax,
		CheckRetry:   DefaultRetryPolicy,
		Backoff:      DefaultBackoff,
	}
}

func (c *Client) logger() *zerolog.Logger {
	// c.loggerInit.Do(func() {
	// 	if c.Logger == nil {
	// 		return
	// 	}

	// 	switch c.Logger.(type) {
	// 	case Logger, LeveledLogger:
	// 		// ok
	// 	default:
	// 		// This should happen in dev when they are setting Logger and work on code, not in prod.
	// 		panic(fmt.Sprintf("invalid logger type passed, must be Logger or LeveledLogger, was %T", c.Logger))
	// 	}
	// })

	return &c.Logger
}

func (c *Request) logger() *zerolog.Logger {
	// c.loggerInit.Do(func() {
	// 	if c.Logger == nil {
	// 		return
	// 	}

	// 	switch c.Logger.(type) {
	// 	case Logger, LeveledLogger:
	// 		// ok
	// 	default:
	// 		// This should happen in dev when they are setting Logger and work on code, not in prod.
	// 		panic(fmt.Sprintf("invalid logger type passed, must be Logger or LeveledLogger, was %T", c.Logger))
	// 	}
	// })

	return zerolog.Ctx(c.Request.Context())
}

// DefaultRetryPolicy provides a default callback for Client.CheckRetry, which
// will retry on connection errors and server errors.
func DefaultRetryPolicy(ctx context.Context, resp *Response, err error) (bool, error) {
	// do not retry on context.Canceled or context.DeadlineExceeded
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	// don't propagate other errors
	shouldRetry, _ := baseRetryPolicy(resp, err)
	return shouldRetry, nil
}

// ErrorPropagatedRetryPolicy is the same as DefaultRetryPolicy, except it
// propagates errors back instead of returning nil. This allows you to inspect
// why it decided to retry or not.
func ErrorPropagatedRetryPolicy(ctx context.Context, resp *Response, err error) (bool, error) {
	// do not retry on context.Canceled or context.DeadlineExceeded
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	return baseRetryPolicy(resp, err)
}

func baseRetryPolicy(resp *Response, err error) (bool, error) {
	if err != nil {
		if v, ok := err.(*url.Error); ok {
			// Don't retry if the error was due to too many redirects.
			if redirectsErrorRe.MatchString(v.Error()) {
				return false, v
			}

			// Don't retry if the error was due to an invalid protocol scheme.
			if schemeErrorRe.MatchString(v.Error()) {
				return false, v
			}

			// Don't retry if the error was due to TLS cert verification failure.
			if notTrustedErrorRe.MatchString(v.Error()) {
				return false, v
			}
			if _, ok := v.Err.(x509.UnknownAuthorityError); ok {
				return false, v
			}
		}

		// The error is likely recoverable so retry.
		return true, nil
	}

	// 429 Too Many Requests is recoverable. Sometimes the server puts
	// a Retry-After response header to indicate when the server is
	// available to start processing request from client.
	if resp.StatusCode == http.StatusTooManyRequests {
		return true, nil
	}

	if resp.StatusCode != http.StatusOK {
		// If the response was bad, we don't want to retry unless it was a
		// 400 Bad Request response with a Retry-After header.
		// See https://tools.ietf.org/html/rfc7231#section-6.5.1
		if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
			return true, nil
		}
	}

	if resp.StatusCode == http.StatusBadRequest {
		b, err := io.ReadAll(resp.Body)

		resp.Body = io.NopCloser(bytes.NewBuffer(b))
		if err == nil {
			var inter map[string]interface{}
			if err = json.Unmarshal(b, &inter); err == nil {
				if _type, ok := inter["__type"].(string); ok {
					if strings.Contains(_type, "ThrottlingException") {
						return true, nil
					}

					if strings.Contains(_type, "ProvisionedThroughputExceeded") {
						return true, nil
					}
				}
			}
		}
	}

	// Check the response code. We retry on 500-range responses to allow
	// the server time to recover, as 500's are typically not permanent
	// errors and may relate to outages on the server side. This will catch
	// invalid response codes as well, like 0 and 999.
	if resp.StatusCode == 0 || (resp.StatusCode >= 500 && resp.StatusCode != http.StatusNotImplemented) {
		return true, fmt.Errorf("unexpected HTTP status %s", resp.Status)
	}

	return false, nil
}

// DefaultBackoff provides a default callback for Client.Backoff which
// will perform exponential backoff based on the attempt number and limited
// by the provided minimum and maximum durations.
//
// It also tries to parse Retry-After response header when a http.StatusTooManyRequests
// (HTTP Code 429) is found in the resp parameter. Hence it will return the number of
// seconds the server states it may be ready to process more requests from this client.
func DefaultBackoff(min, max time.Duration, attemptNum int, resp *Response) time.Duration {
	if resp != nil {
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
			if s, ok := resp.Header["Retry-After"]; ok {
				if sleep, err := strconv.ParseInt(s[0], 10, 64); err == nil {
					return time.Second * time.Duration(sleep)
				}
			}
		}
	}

	mult := math.Pow(2, float64(attemptNum)) * float64(min)
	sleep := time.Duration(mult)
	if float64(sleep) != mult || sleep > max {
		sleep = max
	}
	return sleep
}

// LinearJitterBackoff provides a callback for Client.Backoff which will
// perform linear backoff based on the attempt number and with jitter to
// prevent a thundering herd.
//
// min and max here are *not* absolute values. The number to be multiplied by
// the attempt number will be chosen at random from between them, thus they are
// bounding the jitter.
//
// For instance:
// * To get strictly linear backoff of one second increasing each retry, set
// both to one second (1s, 2s, 3s, 4s, ...)
// * To get a small amount of jitter centered around one second increasing each
// retry, set to around one second, such as a min of 800ms and max of 1200ms
// (892ms, 2102ms, 2945ms, 4312ms, ...)
// * To get extreme jitter, set to a very wide spread, such as a min of 100ms
// and a max of 20s (15382ms, 292ms, 51321ms, 35234ms, ...)
func LinearJitterBackoff(min, max time.Duration, attemptNum int, resp *Response) time.Duration {
	// attemptNum always starts at zero but we want to start at 1 for multiplication
	attemptNum++

	if max <= min {
		// Unclear what to do here, or they are the same, so return min *
		// attemptNum
		return min * time.Duration(attemptNum)
	}

	// Seed rand; doing this every time is fine
	rand := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	// Pick a random number that lies somewhere between the min and max and
	// multiply by the attemptNum. attemptNum starts at zero so we always
	// increment here. We first get a random percentage, then apply that to the
	// difference between min and max, and add to min.
	jitter := rand.Float64() * float64(max-min)
	jitterMin := int64(jitter) + int64(min)
	return time.Duration(jitterMin * int64(attemptNum))
}

// PassthroughErrorHandler is an ErrorHandler that directly passes through the
// values from the net/http library for the final request. The body is not
// closed.
func PassthroughErrorHandler(resp *Response, err error, _ int) (*Response, error) {
	return resp, err
}

type urlValues url.Values

var _ zerolog.LogObjectMarshaler = (*urlValues)(nil)

func (v urlValues) MarshalZerologObject(e *zerolog.Event) {
	for k, vv := range v {
		for _, vvv := range vv {
			e.Str(k, vvv)
		}
	}
}

// Do wraps calling an HTTP method with retries.
func (c *Client) Do(req *Request) (*Response, error) {
	c.clientInit.Do(func() {
		if c.HTTPClient == nil {
			c.HTTPClient = DefaultPooledClient()
		}
	})

	start := time.Now()

	// req.logger().Debug().Str("method", req.Method).Str("url", req.URL.String()).Msg("performing request")

	var resp *Response
	var attempt int
	var shouldRetry bool
	var doErr, respErr, checkErr error

	defer func() {
		even := req.logger().Trace().Dur("duration", time.Since(start)).Int("attempts", attempt).Str("method", req.Method).Str("url", req.URL.String())

		if req.Header.Get("X-Amz-Target") != "" {
			even.Str("x-amz-target", req.Header.Get("X-Amz-Target"))
		}

		stat := 0

		backup := func(name string, r io.Reader) {
			body2, err := ioutil.ReadAll(r)
			if err != nil {
				even.Str(fmt.Sprintf("%s:body", name), "failed to read response body")
			} else {

				even.Str(fmt.Sprintf("%s:body:raw", name), string(body2))
			}
		}

		if resp != nil {
			even.Int("status", resp.StatusCode)
			stat = resp.StatusCode

			contentType := resp.Header.Get("Content-Type")
			even.Str("response:header[content-type]", contentType)

			if stat != 200 {
				reader, err := resp.BodyReader()
				if err != nil {
					even.Str("response:body", "failed to read response body")
				} else {

					if contentType == "application/json" || contentType == "application/x-amz-json-1.0" {
						var res map[string]interface{}
						if err = json.NewDecoder(reader).Decode(&res); err != nil {
							backup("response", reader)
						} else {
							even.Interface("response:body", res)
						}
					} else if contentType == "application/xml" || contentType == "text/xml" {
						j, err := xml.GetErrorResponseComponents(reader, false)
						if err != nil {
							backup("response", reader)
						} else {
							even.Interface("response:body", j)
						}
					} else {
						backup("response", reader)
					}
				}

			}
		}

		if stat != 200 {
			if req.body != nil {
				reader, err := req.body()
				if err != nil {
					req.logger().Error().Err(err).Msg("failed to read request body")
				}

				var res map[string]interface{}
				if err = json.NewDecoder(reader).Decode(&res); err != nil {
					reader, err := req.body()
					if err != nil {
						backup("request", reader)
					}
					bod, err := io.ReadAll(reader)
					if err != nil {
						backup("request", reader)
					}
					res, err := url.ParseQuery(string(bod))
					if err != nil {
						backup("request", reader)
					}

					even.Object("request:body", urlValues(res))
				} else {
					even.Interface("request:body", res)
				}
			}
		}

		even.Msg("request completed")
	}()

	for i := 0; ; i++ {
		doErr, respErr = nil, nil
		attempt++

		// Always rewind the request body when non-nil.
		if req.body != nil {
			body, err := req.body()
			if err != nil {
				c.HTTPClient.CloseIdleConnections()
				return resp, err
			}
			if c, ok := body.(io.ReadCloser); ok {
				req.Request.Body = c
			} else {
				req.Body = ioutil.NopCloser(body)
			}
		}

		if c.RequestLogHook != nil {
			c.RequestLogHook(req.logger(), req.Request, attempt)
		}

		// Attempt the request
		httpresp, ierr := c.HTTPClient.Do(req.Request)
		if ierr != nil {
			doErr = ierr
		} else {
			resp, doErr = NewResponse(httpresp)
		}
		// Check if we should continue with retries.
		shouldRetry, checkErr = c.CheckRetry(req.Context(), resp, doErr)
		if !shouldRetry && doErr == nil && req.responseHandler != nil {
			respErr = req.responseHandler(resp)
			shouldRetry, checkErr = c.CheckRetry(req.Context(), resp, respErr)
		}

		err := doErr
		if respErr != nil {
			err = respErr
		}
		if err != nil {
			req.logger().Error().Err(err).Str("method", req.Method).Str("url", req.URL.String()).Dur("duration", time.Since(start)).Msg("request failed")
		} else {
			// Call this here to maintain the behavior of logging all requests,
			// even if CheckRetry signals to stop.
			if c.ResponseLogHook != nil {
				c.ResponseLogHook(req.logger(), resp)
			}
		}

		if !shouldRetry {
			break
		}

		// We do this before drainBody because there's no need for the I/O if
		// we're breaking out
		remain := c.RetryMax - i
		if remain <= 0 {
			break
		}

		// We're going to retry, consume any response to reuse the connection.
		if doErr == nil {
			c.drainBody(resp.Body)
		}

		wait := c.Backoff(c.RetryWaitMin, c.RetryWaitMax, i, resp)
		status := ""
		if resp != nil {
			status = resp.Status
		}
		req.logger().Debug().Str("method", req.Method).Str("url", req.URL.String()).Str("status", status).Msg("retrying request")
		timer := time.NewTimer(wait)
		select {
		case <-req.Context().Done():
			timer.Stop()
			c.HTTPClient.CloseIdleConnections()
			return nil, req.Context().Err()
		case <-timer.C:
		}

		// Make shallow copy of http Request so that we can modify its body
		// without racing against the closeBody call in persistConn.writeLoop.
		httpreq := *req.Request
		req.Request = &httpreq
	}

	// this is the closest we have to success criteria
	if doErr == nil && respErr == nil && checkErr == nil && !shouldRetry {
		return resp, nil
	}

	defer c.HTTPClient.CloseIdleConnections()

	var err error
	if checkErr != nil {
		err = checkErr
	} else if respErr != nil {
		err = respErr
	} else {
		err = doErr
	}

	if c.ErrorHandler != nil {
		return c.ErrorHandler(resp, err, attempt)
	}

	// By default, we close the response body and return an error without
	// returning the response
	if resp != nil {
		c.drainBody(resp.Body)
	}

	// this means CheckRetry thought the request was a failure, but didn't
	// communicate why
	if err == nil {
		return nil, fmt.Errorf("%s %s giving up after %d attempt(s)",
			req.Method, req.URL, attempt)
	}

	return nil, fmt.Errorf("%s %s giving up after %d attempt(s): %w",
		req.Method, req.URL, attempt, err)
}

// Try to read the response body so we can reuse this connection.
func (c *Client) drainBody(body io.ReadCloser) {
	defer body.Close()
	_, err := io.Copy(ioutil.Discard, io.LimitReader(body, respReadLimit))
	if err != nil {
		c.logger().Error().Err(err).Msg("error reading response body")
	}
}

// Get is a shortcut for doing a GET request without making a new client.
func Get(ctx context.Context, url string) (*Response, error) {
	return defaultClient.Get(ctx, url)
}

// Get is a convenience helper for doing simple GET requests.
func (c *Client) Get(ctx context.Context, url string) (*Response, error) {
	req, err := NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Head is a shortcut for doing a HEAD request without making a new client.
func Head(ctx context.Context, url string) (*Response, error) {
	return defaultClient.Head(ctx, url)
}

// Head is a convenience method for doing simple HEAD requests.
func (c *Client) Head(ctx context.Context, url string) (*Response, error) {
	req, err := NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Post is a shortcut for doing a POST request without making a new client.
func Post(ctx context.Context, url, bodyType string, body interface{}) (*Response, error) {
	return defaultClient.Post(ctx, url, bodyType, body)
}

// Post is a convenience method for doing simple POST requests.
func (c *Client) Post(ctx context.Context, url, bodyType string, body interface{}) (*Response, error) {
	req, err := NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", bodyType)
	return c.Do(req)
}

// PostForm is a shortcut to perform a POST with form data without creating
// a new client.
func PostForm(ctx context.Context, url string, data url.Values) (*Response, error) {
	return defaultClient.PostForm(ctx, url, data)
}

// PostForm is a convenience method for doing simple POST operations using
// pre-filled url.Values form data.
func (c *Client) PostForm(ctx context.Context, url string, data url.Values) (*Response, error) {
	return c.Post(ctx, url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

// StandardClient returns a stdlib *http.Client with a custom Transport, which
// shims in a *retryablehttp.Client for added retries.
func (c *Client) StandardClient() *http.Client {
	return &http.Client{
		Transport: &RoundTripper{Client: c},
	}
}
