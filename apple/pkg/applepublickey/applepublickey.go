package applepublickey

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/context/ctxhttp"
)

type ApplePublicKey struct {
	KTY string `json:"kty"`
	KID string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type PublicKeyResponse struct {
	Keys []ApplePublicKey `json:"keys"`
}

type Client struct {
	validationURL string
	publicKeys    *PublicKeyResponse
	ttl           time.Time
}

func NewClient(endpoint *url.URL) *Client {
	return &Client{
		validationURL: endpoint.String(),
	}
}

func (client *Client) Refresh(ctx context.Context) (*PublicKeyResponse, error) {

	if client.publicKeys != nil && client.ttl.Before(time.Now()) {
		return client.publicKeys, nil
	}

	keys, err := fetchPublicKeys(ctx, client.validationURL)
	if err != nil {
		return nil, err
	}

	client.ttl = time.Now().Add(time.Minute * 5)
	client.publicKeys = keys

	return keys, nil
}

// RSA returns a corresponding *rsa.PublicKey
func (k *ApplePublicKey) RSA() *rsa.PublicKey {
	return &rsa.PublicKey{
		N: decodeBase64BigInt(k.N),
		E: int(decodeBase64BigInt(k.E).Int64()),
	}
}

func fetchPublicKeys(ctx context.Context, url string) (*PublicKeyResponse, error) {

	response, err := ctxhttp.Get(ctx, nil, url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	appleKeys := PublicKeyResponse{}
	err = json.NewDecoder(response.Body).Decode(&appleKeys)
	if err != nil {
		return nil, err
	}

	return &appleKeys, nil
}

// decodeBase64BigInt decodes a base64-encoded larger integer from Apple's key format.
func decodeBase64BigInt(s string) *big.Int {
	buffer, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(s)
	if err != nil {
		log.Println("error decoding base64: ", err.Error())
		return big.NewInt(0)
	}

	return big.NewInt(0).SetBytes(buffer)
}

func (client *PublicKeyResponse) GetPublicKey(kid string) *ApplePublicKey {
	for _, key := range client.Keys {
		if key.KID == kid {
			return &key
		}
	}
	return nil
}
