package jwt

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"

	"github.com/golang-jwt/jwt/v4"
)

// RSA returns a corresponding *rsa.PublicKey
func (k appleKey) RSA() *rsa.PublicKey {
	return &rsa.PublicKey{
		N: k.N,
		E: k.E,
	}
}

// decodeBase64BigInt decodes a base64-encoded larger integer from Apple's key format.
func decodeBase64BigInt(s string) *big.Int {
	buffer, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(s)
	if err != nil {
		log.Fatalf("failed to decode base64: %v", err)
	}

	return big.NewInt(0).SetBytes(buffer)
}

// appleKey is a type of public key.
type appleKey struct {
	KTY string
	KID string
	Use string
	Alg string
	N   *big.Int
	E   int
}

type stringfiedAppleKey struct {
	KTY string `json:"kty"`
	KID string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func (k *stringfiedAppleKey) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &k); err != nil {
		return err
	}
	return nil
}

func (k *stringfiedAppleKey) toAppleKey() appleKey {
	return appleKey{
		KTY: k.KTY,
		KID: k.KID,
		Use: k.Use,
		Alg: k.Alg,
		N:   decodeBase64BigInt(k.N),
		E:   int(decodeBase64BigInt(k.E).Int64()),
	}
}

// UnmarshalJSON parses a JSON-encoded value and stores the result in the object.
func (k *appleKey) UnmarshalJSON(b []byte) error {
	var tmp stringfiedAppleKey

	err := tmp.UnmarshalJSON(b)
	if err != nil {
		return err
	}

	*k = tmp.toAppleKey()

	return nil
}

func NewAppleClient(ctx context.Context, endpoint *url.URL) (*Client, error) {

	res, err := http.Get(endpoint.String())
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	var appleJWKs struct {
		Keys []appleKey `json:"keys"`
	}

	err = json.Unmarshal(body, &appleJWKs)
	if err != nil {
		log.Fatal(err)
	}

	keyfunc := func(t *jwt.Token) (interface{}, error) {

		// check the signing method
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		// check the kid
		kid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid not found in token header")
		}

		var match appleKey

		// find the key
		for _, key := range appleJWKs.Keys {
			if key.KID == kid {
				match = key
				break
			}
		}

		if match.KID == "" {
			return nil, fmt.Errorf("unable to find key for kid: %s", kid)
		}

		return match.RSA(), nil
	}

	return &Client{keyfunc}, nil
}
