package protocol

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"nugg-auth/core/pkg/hex"
	"nugg-auth/core/pkg/webauthn/protocol/webauthncbor"

	"github.com/k0kubun/pp"
)

var byteID = hex.MustBase64ToHash("6xrtBhJQW6QU4tOaB4rrHaS2Ks0yDDL_q8jDC16DEjZ-VLVf4kCRkvl2xp2D71sTPYns-exsHQHTy3G-zJRK8g")
var byteClientDataJSON = `{\"challenge\":\"0x5bc1b3154f291a386845b5ab2c395a9807eaff2e12d42646d55ba87912c046b1\",\"origin\":\"https://webauthn.io\",\"type\":\"webauthn.create\"}`
var byteAttObject = hex.MustBase64ToHash("o2NmbXRkbm9uZWdhdHRTdG10oGhhdXRoRGF0YVjEdKbqkhPJnC90siSSsyDPQCYqlMGpUKA5fyklC2CEHvBBAAAAAAAAAAAAAAAAAAAAAAAAAAAAQOsa7QYSUFukFOLTmgeK6x2ktirNMgwy_6vIwwtegxI2flS1X-JAkZL5dsadg-9bEz2J7PnsbB0B08txvsyUSvKlAQIDJiABIVggLKF5xS0_BntttUIrm2Z2tgZ4uQDwllbdIfrrBMABCNciWCDHwin8Zdkr56iSIh0MrB5qZiEzYLQpEOREhMUkY6q4Vw")
var abc = "{\"challenge\":\"0x5bc1b3154f291a386845b5ab2c395a9807eaff2e12d42646d55ba87912c046b1\",\"origin\":\"https://webauthn.io\",\"type\":\"webauthn.create\"}"

func TestParseCredentialCreationResponse(t *testing.T) {

	reqBody := ioutil.NopCloser(bytes.NewReader([]byte(testCredentialRequestBody)))
	httpReq := &http.Request{Body: reqBody}
	type args struct {
		response *http.Request
	}

	byteAuthData := hex.MustBase64ToHash("dKbqkhPJnC90siSSsyDPQCYqlMGpUKA5fyklC2CEHvBBAAAAAAAAAAAAAAAAAAAAAAAAAAAAQOsa7QYSUFukFOLTmgeK6x2ktirNMgwy_6vIwwtegxI2flS1X-JAkZL5dsadg-9bEz2J7PnsbB0B08txvsyUSvKlAQIDJiABIVggLKF5xS0_BntttUIrm2Z2tgZ4uQDwllbdIfrrBMABCNciWCDHwin8Zdkr56iSIh0MrB5qZiEzYLQpEOREhMUkY6q4Vw")
	byteRPIDHash := hex.MustBase64ToHash("dKbqkhPJnC90siSSsyDPQCYqlMGpUKA5fyklC2CEHvA")
	byteCredentialPubKey := hex.MustBase64ToHash("pSJYIMfCKfxl2SvnqJIiHQysHmpmITNgtCkQ5ESExSRjqrhXAQIDJiABIVggLKF5xS0_BntttUIrm2Z2tgZ4uQDwllbdIfrrBMABCNc")

	tests := []struct {
		name    string
		args    args
		want    *ParsedCredentialCreationData
		wantErr bool
	}{
		{
			name: "Successful Credential Request Parsing",
			args: args{
				response: httpReq,
			},
			want: &ParsedCredentialCreationData{
				ParsedPublicKeyCredential: ParsedPublicKeyCredential{
					ParsedCredential: ParsedCredential{
						ID:   byteID,
						Type: "public-key",
					},
					ClientExtensionResults: AuthenticationExtensionsClientOutputs{
						"appid": true,
					},
					RawID: byteID,
				},
				Response: ParsedAttestationResponse{
					CollectedClientData: CollectedClientData{
						Type:      CeremonyType("webauthn.create"),
						Challenge: hex.MustBase64ToHash("W8GzFU8pGjhoRbWrLDlamAfq_y4S1CZG1VuoeRLARrE"),
						Origin:    "https://webauthn.io",
					},
					AttestationObject: AttestationObject{
						Format:      "none",
						RawAuthData: byteAuthData,
						AuthData: AuthenticatorData{
							RPIDHash: byteRPIDHash,
							Counter:  0,
							Flags:    0x041,
							AttData: AttestedCredentialData{
								AAGUID:              make([]byte, 16),
								CredentialID:        byteID,
								CredentialPublicKey: byteCredentialPubKey,
							},
						},
					},
				},
				Raw: CredentialCreationResponse{
					PublicKeyCredential: PublicKeyCredential{
						Credential: Credential{
							Type: "public-key",
							ID:   byteID,
						},
						ClientExtensionResults: AuthenticationExtensionsClientOutputs{
							"appid": true,
						},
						RawID: byteID,
					},
					AttestationResponse: AuthenticatorAttestationResponse{
						AuthenticatorResponse: AuthenticatorResponse{
							UTF8ClientDataJSON: abc,
						},
						AttestationObject: byteAttObject,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := ParseCredentialCreationResponse(tt.args.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCredentialCreationResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.ClientExtensionResults, tt.want.ClientExtensionResults) {
				t.Errorf("Extensions = %v \n want: %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.ID, tt.want.ID) {
				t.Errorf("ID = %v \n want: %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.ParsedCredential, tt.want.ParsedCredential) {
				t.Errorf("ParsedCredential = %v \n want: %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.ParsedPublicKeyCredential, tt.want.ParsedPublicKeyCredential) {
				t.Errorf("ParsedPublicKeyCredential = %v \n want: %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.Raw, tt.want.Raw) {
				pp.Println(got.Raw)
				pp.Println(tt.want.Raw)
				t.Errorf("Raw = %v \n want: %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.RawID, tt.want.RawID) {
				t.Errorf("RawID = %v \n want: %v", got, tt.want)
			}
			// Unmarshall CredentialPublicKey
			var pkWant interface{}
			keyBytesWant := tt.want.Response.AttestationObject.AuthData.AttData.CredentialPublicKey
			webauthncbor.Unmarshal(keyBytesWant, &pkWant)
			var pkGot interface{}
			keyBytesGot := got.Response.AttestationObject.AuthData.AttData.CredentialPublicKey
			webauthncbor.Unmarshal(keyBytesGot, &pkGot)
			if !reflect.DeepEqual(pkGot, pkWant) {
				t.Errorf("Response = %+v \n want: %+v", pkGot, pkWant)
			}
			if !reflect.DeepEqual(got.Type, tt.want.Type) {
				t.Errorf("Type = %v \n want: %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.Response.CollectedClientData, tt.want.Response.CollectedClientData) {
				t.Errorf("CollectedClientData = %v \n want: %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.Response.AttestationObject.Format, tt.want.Response.AttestationObject.Format) {
				t.Errorf("Format = %v \n want: %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.Response.AttestationObject.AuthData.AttData.CredentialID, tt.want.Response.AttestationObject.AuthData.AttData.CredentialID) {
				t.Errorf("CredentialID = %v \n want: %v", got, tt.want)
			}
		})
	}
}

func TestParsedCredentialCreationData_Verify(t *testing.T) {
	byteID := hex.MustBase64ToHash("6xrtBhJQW6QU4tOaB4rrHaS2Ks0yDDL_q8jDC16DEjZ-VLVf4kCRkvl2xp2D71sTPYns-exsHQHTy3G-zJRK8g")
	byteChallenge := hex.MustBase64ToHash("W8GzFU8pGjhoRbWrLDlamAfq_y4S1CZG1VuoeRLARrE")
	byteAuthData := hex.MustBase64ToHash("dKbqkhPJnC90siSSsyDPQCYqlMGpUKA5fyklC2CEHvBBAAAAAAAAAAAAAAAAAAAAAAAAAAAAQOsa7QYSUFukFOLTmgeK6x2ktirNMgwy_6vIwwtegxI2flS1X-JAkZL5dsadg-9bEz2J7PnsbB0B08txvsyUSvKlAQIDJiABIVggLKF5xS0_BntttUIrm2Z2tgZ4uQDwllbdIfrrBMABCNciWCDHwin8Zdkr56iSIh0MrB5qZiEzYLQpEOREhMUkY6q4Vw")
	byteRPIDHash := hex.MustBase64ToHash("dKbqkhPJnC90siSSsyDPQCYqlMGpUKA5fyklC2CEHvA")
	byteCredentialPubKey := hex.MustBase64ToHash("pSJYIMfCKfxl2SvnqJIiHQysHmpmITNgtCkQ5ESExSRjqrhXAQIDJiABIVggLKF5xS0_BntttUIrm2Z2tgZ4uQDwllbdIfrrBMABCNc")
	byteAttObject := hex.MustBase64ToHash("o2NmbXRkbm9uZWdhdHRTdG10oGhhdXRoRGF0YVjEdKbqkhPJnC90siSSsyDPQCYqlMGpUKA5fyklC2CEHvBBAAAAAAAAAAAAAAAAAAAAAAAAAAAAQOsa7QYSUFukFOLTmgeK6x2ktirNMgwy_6vIwwtegxI2flS1X-JAkZL5dsadg-9bEz2J7PnsbB0B08txvsyUSvKlAQIDJiABIVggLKF5xS0_BntttUIrm2Z2tgZ4uQDwllbdIfrrBMABCNciWCDHwin8Zdkr56iSIh0MrB5qZiEzYLQpEOREhMUkY6q4Vw")
	byteClientDataJSON := hex.MustBase64ToHash("eyJjaGFsbGVuZ2UiOiJXOEd6RlU4cEdqaG9SYldyTERsYW1BZnFfeTRTMUNaRzFWdW9lUkxBUnJFIiwib3JpZ2luIjoiaHR0cHM6Ly93ZWJhdXRobi5pbyIsInR5cGUiOiJ3ZWJhdXRobi5jcmVhdGUifQ")

	type fields struct {
		ParsedPublicKeyCredential ParsedPublicKeyCredential
		Response                  ParsedAttestationResponse
		Raw                       CredentialCreationResponse
	}
	type args struct {
		storedChallenge    hex.Hash
		verifyUser         bool
		relyingPartyID     string
		relyingPartyOrigin string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Successful Verification Test",
			fields: fields{
				ParsedPublicKeyCredential: ParsedPublicKeyCredential{
					ParsedCredential: ParsedCredential{
						ID:   byteID,
						Type: "public-key",
					},
					RawID: byteID,
				},
				Response: ParsedAttestationResponse{
					CollectedClientData: CollectedClientData{
						Type:      CeremonyType("webauthn.create"),
						Challenge: hex.MustBase64ToHash("W8GzFU8pGjhoRbWrLDlamAfq_y4S1CZG1VuoeRLARrE"),
						Origin:    "https://webauthn.io",
					},
					AttestationObject: AttestationObject{
						Format:      "none",
						RawAuthData: byteAuthData,
						AuthData: AuthenticatorData{
							RPIDHash: byteRPIDHash,
							Counter:  0,
							Flags:    0x041,
							AttData: AttestedCredentialData{
								AAGUID:              make([]byte, 16),
								CredentialID:        byteID,
								CredentialPublicKey: byteCredentialPubKey,
							},
						},
					},
				},
				Raw: CredentialCreationResponse{
					PublicKeyCredential: PublicKeyCredential{
						Credential: Credential{
							Type: "public-key",
							ID:   byteID,
						},
						RawID: byteID,
					},
					AttestationResponse: AuthenticatorAttestationResponse{
						AuthenticatorResponse: AuthenticatorResponse{
							UTF8ClientDataJSON: byteClientDataJSON.Utf8(),
						},
						AttestationObject: byteAttObject,
					},
				},
			},
			args: args{
				storedChallenge:    byteChallenge,
				verifyUser:         false,
				relyingPartyID:     `webauthn.io`,
				relyingPartyOrigin: `https://webauthn.io`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pcc := &ParsedCredentialCreationData{
				ParsedPublicKeyCredential: tt.fields.ParsedPublicKeyCredential,
				Response:                  tt.fields.Response,
				Raw:                       tt.fields.Raw,
			}
			if _, err := pcc.Verify(tt.args.storedChallenge, hex.Hash{}, tt.args.verifyUser, tt.args.relyingPartyID, tt.args.relyingPartyOrigin); (err != nil) != tt.wantErr {
				t.Errorf("ParsedCredentialCreationData.Verify() error = %+v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

var testCredentialRequestBody = fmt.Sprintf(
	`{"id":"%s","rawId":"%s","type":"public-key","clientExtensionResults":{"appid":true},"response":{"attestationObject":"%s","clientDataJSON":"%s"}}`,
	byteID.Hex(), byteID.Hex(), byteAttObject.Hex(), byteClientDataJSON)
