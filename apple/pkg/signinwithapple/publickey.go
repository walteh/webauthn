package signinwithapple

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

func fetchPublicKeys(client *http.Client, url string) (*PublicKeyResponse, error) {
	response, err := client.Get(url)
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

func (client *Client) RefreshPublicKeys() error {

	keys, err := fetchPublicKeys(client.httpClient, client.validationURL)
	if err != nil {
		return err
	}

	client.publicKeys = keys

	return nil
}

// RSA returns a corresponding *rsa.PublicKey
func (k ApplePublicKey) RSA() *rsa.PublicKey {
	return &rsa.PublicKey{
		N: decodeBase64BigInt(k.N),
		E: int(decodeBase64BigInt(k.E).Int64()),
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

func (client *PublicKeyResponse) GetPublicKey(kid string) *ApplePublicKey {
	for _, key := range client.Keys {
		if key.KID == kid {
			return &key
		}
	}
	return nil
}

func (client *PublicKeyResponse) BuildKeyFunc() (jwt.Keyfunc, error) {

	return func(t *jwt.Token) (interface{}, error) {

		// check the signing method
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		// check the kid
		kid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid not found in token header")
		}

		// get the public key
		key := client.GetPublicKey(kid)
		if key == nil {
			return nil, fmt.Errorf("key not found for kid: %s", kid)
		}

		return key.RSA(), nil

	}, nil

}
