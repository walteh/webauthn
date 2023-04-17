package signinwithapple

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// ValidationURL is the endpoint for verifying tokens
	ValidationURL string = "https://appleid.apple.com/auth/token"
	// ContentType is the one expected by Apple
	ContentType string = "application/x-www-form-urlencoded"
	// UserAgent is required by Apple or the request will fail
	UserAgent string = "nugg.xyz/aws"
	// AcceptHeader is the content that we are willing to accept
	AcceptHeader string = "application/json"
)

// ValidationClient is an interface to call the validation API
type ValidationClient interface {
	VerifyWebToken(ctx context.Context, reqBody WebValidationTokenRequest, result interface{}) error
	VerifyAppToken(ctx context.Context, reqBody AppValidationTokenRequest, result interface{}) error
	VerifyRefreshToken(ctx context.Context, reqBody ValidationRefreshRequest, result interface{}) error
}

// Client implements ValidationClient
type Client struct {
	validationURL string
	httpClient    *http.Client

	config *ClientConfig
}

// New creates a Client object
func NewClient(url *url.URL, teamId string, serviceId string, keyId string) *Client {
	return newClientWithUrlString(url.String(), teamId, serviceId, keyId)
}

func newClientWithUrlString(url string, teamId string, serviceId string, keyId string) *Client {

	client := &Client{
		validationURL: url,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		config: &ClientConfig{
			TeamID:   teamId,
			ClientID: serviceId,
			KeyID:    keyId,
		},
	}

	return client
}

// NewWithURL creates a Client object with a custom URL provided

// VerifyWebToken sends the WebValidationTokenRequest and gets validation result
func (c *Client) VerifyWebToken(ctx context.Context, reqBody WebValidationTokenRequest) (res *ValidationResponse, err error) {
	data := url.Values{}
	data.Set("client_id", reqBody.ClientID)
	data.Set("client_secret", reqBody.ClientSecret)
	data.Set("code", reqBody.Code)
	data.Set("redirect_uri", reqBody.RedirectURI)
	data.Set("grant_type", "authorization_code")

	err = doRequest(ctx, c.httpClient, &res, c.validationURL, data)

	return
}

// VerifyAppToken sends the AppValidationTokenRequest and gets validation result
func (c *Client) VerifyAppToken(ctx context.Context, reqBody AppValidationTokenRequest) (res *ValidationResponse, err error) {
	data := url.Values{}
	data.Set("client_id", reqBody.ClientID)
	data.Set("client_secret", reqBody.ClientSecret)
	data.Set("code", reqBody.Code)
	data.Set("grant_type", "authorization_code")

	err = doRequest(ctx, c.httpClient, &res, c.validationURL, data)

	return
}

// VerifyRefreshToken sends the WebValidationTokenRequest and gets validation result
func (c *Client) VerifyRefreshToken(ctx context.Context, reqBody ValidationRefreshRequest) (res *RefreshResponse, err error) {
	data := url.Values{}
	data.Set("client_id", reqBody.ClientID)
	data.Set("client_secret", reqBody.ClientSecret)
	data.Set("refresh_token", reqBody.RefreshToken)
	data.Set("grant_type", "refresh_token")

	err = doRequest(ctx, c.httpClient, &res, c.validationURL, data)

	return
}

func doRequest(ctx context.Context, client *http.Client, result interface{}, url string, data url.Values) error {
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("content-type", ContentType)
	req.Header.Add("accept", AcceptHeader)
	req.Header.Add("user-agent", UserAgent) // apple requires a user agent

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(result)
}
