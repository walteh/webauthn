package protocol

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"nugg-auth/core/pkg/webauthn/protocol/webauthncose"
)

// The raw response returned to us from an authenticator when we request a
// credential for login/assertion.
type CredentialAssertionResponse struct {
	PublicKeyCredential
	AssertionResponse AuthenticatorAssertionResponse `json:"response"`
}

// The parsed CredentialAssertionResponse that has been marshalled into a format
// that allows us to verify the client and authenticator data inside the response
type ParsedCredentialAssertionData struct {
	ParsedPublicKeyCredential
	Response ParsedAssertionResponse
	Raw      CredentialAssertionResponse
}

// The AuthenticatorAssertionResponse contains the raw authenticator assertion data and is parsed into
// ParsedAssertionResponse
type AuthenticatorAssertionResponse struct {
	AuthenticatorResponse
	AuthenticatorData URLEncodedBase64 `json:"authenticatorData"`
	Signature         URLEncodedBase64 `json:"signature"`
	UserHandle        URLEncodedBase64 `json:"userHandle,omitempty"`
}

// Parsed form of AuthenticatorAssertionResponse
type ParsedAssertionResponse struct {
	CollectedClientData CollectedClientData
	AuthenticatorData   AuthenticatorData
	Signature           []byte
	UserHandle          []byte
}

type BetterCredentialAssertionResponse struct {
	UserID               []byte `json:"userID"`
	CredentialID         []byte `json:"credentialID"`
	RawClientDataJSON    []byte `json:"rawClientDataJSON"`
	RawAuthenticatorData []byte `json:"rawAuthenticatorData"`
	Signature            []byte `json:"signature"`
	Type                 string `json:"credentialType"`
}

// Parse the credential request response into a format that is either required by the specification
// or makes the assertion verification steps easier to complete. This takes an http.Request that contains
// the assertion response data in a raw, mostly base64 encoded format, and parses the data into
// manageable structures
func ParseCredentialRequestResponse(response *http.Request) (*ParsedCredentialAssertionData, error) {
	if response == nil || response.Body == nil {
		return nil, ErrBadRequest.WithDetails("No response given")
	}
	return ParseCredentialRequestResponseBody(response.Body)
}

// Parse the credential request response into a format that is either required by the specification
// or makes the assertion verification steps easier to complete. This takes an io.Reader that contains
// the assertion response data in a raw, mostly base64 encoded format, and parses the data into
// manageable structures
func ParseCredentialRequestResponseBody(body io.Reader) (*ParsedCredentialAssertionData, error) {
	var car CredentialAssertionResponse
	err := json.NewDecoder(body).Decode(&car)
	if err != nil {
		return nil, ErrBadRequest.WithDetails("Parse error for Assertion")
	}

	return ParseCredentialAssertionResponse(car)
}

// Parse the credential request response into a format that is either required by the specification
// or makes the assertion verification steps easier to complete. This takes an CredentialAssertionResponse that contains
// the assertion response data in a raw, mostly base64 encoded format, and parses the data into
func ParseCredentialAssertionResponsePayload(body string) (*BetterCredentialAssertionResponse, error) {

	dec, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		return nil, ErrBadRequest.WithDetails("Parse error for Assertion")
	}

	var car BetterCredentialAssertionResponse
	err = json.Unmarshal(dec, &car)
	if err != nil {
		return nil, ErrBadRequest.WithDetails("Parse error for Assertion")
	}

	return &car, nil
}

// func DecodeCredentialAssertionJSON(body *BetterCredentialAssertionResponseString) *BetterCredentialAssertionResponse {

// 	return &BetterCredentialAssertionResponse{
// 		UserID:               URLEncodedBase64(ResolveToRawURLEncoding(body.UserID)),
// 		CredentialID:         URLEncodedBase64(ResolveToRawURLEncoding(body.CredentialID)),
// 		Signature:            URLEncodedBase64(ResolveToBase64(body.Signature)),
// 		Type:                 body.Type,
// 		RawClientDataJSON:    []byte(body.RawClientDataJSON),
// 		RawAuthenticatorData: URLEncodedBase64(ResolveToRawURLEncoding(body.RawAuthenticatorData)),
// 	}
// }

// Parse the credential request response into a format that is either required by the specification
// or makes the assertion verification steps easier to complete. This takes an io.Reader that contains
// the assertion response data in a raw, mostly base64 encoded format, and parses the data into
// manageable structures
func DecodeCredentialAssertionResponse(car *BetterCredentialAssertionResponse) *CredentialAssertionResponse {

	return &CredentialAssertionResponse{
		PublicKeyCredential: PublicKeyCredential{
			Credential: Credential{
				Type: car.Type,
				ID:   base64.RawURLEncoding.EncodeToString(car.CredentialID),
			},
			RawID: car.CredentialID,
			ClientExtensionResults: map[string]interface{}{
				"appid": true,
			},
		},
		AssertionResponse: AuthenticatorAssertionResponse{
			AuthenticatorResponse: AuthenticatorResponse{
				ClientDataJSON: car.RawClientDataJSON,
			},
			AuthenticatorData: car.RawAuthenticatorData,
			Signature:         car.Signature,
			UserHandle:        car.UserID,
		},
	}

}

// Parse the credential request response into a format that is either required by the specification
// or makes the assertion verification steps easier to complete. This takes an io.Reader that contains
// the assertion response data in a raw, mostly base64 encoded format, and parses the data into
// manageable structures
func ParseCredentialAssertionResponse(car CredentialAssertionResponse) (*ParsedCredentialAssertionData, error) {

	if car.ID == "" {
		return nil, ErrBadRequest.WithDetails("CredentialAssertionResponse with ID missing")
	}

	_, err := base64.RawURLEncoding.DecodeString(car.ID)
	if err != nil {
		return nil, ErrBadRequest.WithDetails("CredentialAssertionResponse with ID not base64url encoded").WithParent(err)
	}
	if car.Type != "public-key" {
		return nil, ErrBadRequest.WithDetails("CredentialAssertionResponse with bad type")
	}
	var par ParsedCredentialAssertionData
	par.ID, par.RawID, par.Type, par.ClientExtensionResults = car.ID, car.RawID, car.Type, car.ClientExtensionResults
	par.Raw = car

	par.Response.Signature = car.AssertionResponse.Signature
	par.Response.UserHandle = car.AssertionResponse.UserHandle

	// Step 5. Let JSONtext be the result of running UTF-8 decode on the value of cData.
	// We don't call it cData but this is Step 5 in the spec.
	err = json.Unmarshal(car.AssertionResponse.ClientDataJSON, &par.Response.CollectedClientData)
	if err != nil {
		return nil, err
	}

	err = par.Response.AuthenticatorData.Unmarshal(car.AssertionResponse.AuthenticatorData)
	if err != nil {
		return nil, ErrParsingData.WithDetails("Error unmarshalling auth data").WithParent(err)
	}
	return &par, nil
}

// Follow the remaining steps outlined in §7.2 Verifying an authentication assertion
// (https://www.w3.org/TR/webauthn/#verifying-assertion) and return an error if there
// is a failure during each step.
func (p *ParsedCredentialAssertionData) Verify(storedChallenge Challenge, relyingPartyID, relyingPartyOrigin, attestationType string, verifyUser bool, credentialBytes []byte, extensions AuthenticationExtensions) error {
	// Steps 4 through 6 in verifying the assertion data (https://www.w3.org/TR/webauthn/#verifying-assertion) are
	// "assertive" steps, i.e "Let JSONtext be the result of running UTF-8 decode on the value of cData."
	// We handle these steps in part as we verify but also beforehand
	var (
		key interface{}
		err error
	)

	if extensions == nil {
		extensions = AuthenticationExtensions{}
	}

	appID, err := p.GetAppID(extensions, attestationType)
	if err != nil {
		return err
	}

	// Handle steps 7 through 10 of assertion by verifying stored data against the Collected Client Data
	// returned by the authenticator
	validError := p.Response.CollectedClientData.Verify(storedChallenge, AssertCeremony, relyingPartyOrigin)
	if validError != nil {
		return validError
	}

	// Begin Step 11. Verify that the rpIdHash in authData is the SHA-256 hash of the RP ID expected by the RP.
	rpIDHash := sha256.Sum256([]byte(relyingPartyID))

	var appIDHash [32]byte
	if appID != "" {
		appIDHash = sha256.Sum256([]byte(appID))
	}

	// Handle steps 11 through 14, verifying the authenticator data.
	validError = p.Response.AuthenticatorData.Verify(rpIDHash[:], appIDHash[:], verifyUser, true)
	if validError != nil {
		return validError
	}

	// allowedUserCredentialIDs := session.AllowedCredentialIDs

	// Step 15. Let hash be the result of computing a hash over the cData using SHA-256.
	clientDataHash := sha256.Sum256(p.Raw.AssertionResponse.ClientDataJSON)

	// Step 16. Using the credential public key looked up in step 3, verify that sig is
	// a valid signature over the binary concatenation of authData and hash.

	sigData := append(p.Raw.AssertionResponse.AuthenticatorData, clientDataHash[:]...)

	if appID == "" {
		key, err = webauthncose.ParsePublicKey(credentialBytes)
	} else {
		key, err = webauthncose.ParseFIDOPublicKey(credentialBytes)
	}

	if err != nil {
		return ErrAssertionSignature.WithDetails(fmt.Sprintf("Error parsing the assertion public key: %+v", err))
	}

	valid, err := webauthncose.VerifySignature(key, sigData, p.Response.Signature)
	if !valid || err != nil {
		log.Println("valid", valid, "err", err)
		return ErrAssertionSignature.WithDetails(fmt.Sprintf("Error validating the assertion signature: %+v\n", err))
	}

	return nil
}
