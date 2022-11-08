package jwt

import (
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
		log.Println("error decoding base64: ", err.Error())
		return big.NewInt(0)
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

type AppleJwtClient struct {
	endpoint *url.URL
}

func NewAppleJwtClient(endpoint *url.URL) (*AppleJwtClient, error) {
	return &AppleJwtClient{
		endpoint: endpoint,
	}, nil

}

type AppleResponsePayload struct {
	Keys []stringfiedAppleKey `json:"keys"`
}

func (me *AppleJwtClient) BuildKeyFunc() (jwt.Keyfunc, error) {

	res, err := http.Get(me.endpoint.String())
	if err != nil {
		log.Println("error getting apple keys: ", err.Error())
		log.Println("endpoint: ", me.endpoint)
		return nil, fmt.Errorf("error fetching keys from apple: %v", err)
	}

	log.Println("response status from apple: ", res.Status)

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Println("error reading apple keys: ", err.Error())
		log.Println("response: ", res)
		log.Println("body: ", string(body))
		return nil, fmt.Errorf("error reading keys from apple: %v", err)
	}

	log.Println("body", string(body))

	var payload AppleResponsePayload

	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Println("error unmarshalling apple keys: ", err.Error())
		return nil, err
	}

	return func(t *jwt.Token) (interface{}, error) {

		// check the signing method
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			log.Println("error signing method")
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		// check the kid
		kid, ok := t.Header["kid"].(string)
		if !ok {
			log.Println("error kid")
			return nil, fmt.Errorf("kid not found in token header")
		}

		var match appleKey

		// find the key
		for _, key := range payload.Keys {
			if key.KID == kid {
				match = key.toAppleKey()
				break
			}
		}

		if match.KID == "" {
			log.Println("error match kid")
			return nil, fmt.Errorf("unable to find key for kid: %s", kid)
		}

		return match.RSA(), nil

	}, nil

}
