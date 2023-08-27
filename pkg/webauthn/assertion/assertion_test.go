package assertion_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rs/zerolog"
	"github.com/walteh/webauthn/pkg/hex"
	"github.com/walteh/webauthn/pkg/webauthn/assertion"
	"github.com/walteh/webauthn/pkg/webauthn/types"
)

// func TestParseCredentialRequestResponse(t *testing.T) {
// 	reqBody := ioutil.NopCloser(bytes.NewReader([]byte(testAssertionResponses["success"])))
// 	httpReq := &http.Request{Body: reqBody}

// 	byteID := hex.MustBase64ToHash("AI7D5q2P0LS-Fal9ZT7CHM2N5BLbUunF92T8b6iYC199bO2kagSuU05-5dZGqb1SP0A0lyTWng")
// 	byteAAGUID := hex.MustBase64ToHash("rc4AAjW8xgpkiwsl8fBVAw")
// 	byteRPIDHash := hex.MustBase64ToHash("dKbqkhPJnC90siSSsyDPQCYqlMGpUKA5fyklC2CEHvA")
// 	byteAuthData := hex.MustBase64ToHash("dKbqkhPJnC90siSSsyDPQCYqlMGpUKA5fyklC2CEHvBFXJJiGa3OAAI1vMYKZIsLJfHwVQMANwCOw-atj9C0vhWpfWU-whzNjeQS21Lpxfdk_G-omAtffWztpGoErlNOfuXWRqm9Uj9ANJck1p6lAQIDJiABIVggKAhfsdHcBIc0KPgAcRyAIK_-Vi-nCXHkRHPNaCMBZ-4iWCBxB8fGYQSBONi9uvq0gv95dGWlhJrBwCsj_a4LJQKVHQ")
// 	byteSignature := hex.MustBase64ToHash("MEUCIBtIVOQxzFYdyWQyxaLR0tik1TnuPhGVhXVSNgFwLmN5AiEAnxXdCq0UeAVGWxOaFcjBZ_mEZoXqNboY5IkQDdlWZYc")
// 	byteUserHandle := hex.MustBase64ToHash("0ToAAAAAAAAAAA")
// 	byteCredentialPubKey := hex.MustBase64ToHash("pQMmIAEhWCAoCF-x0dwEhzQo-ABxHIAgr_5WL6cJceREc81oIwFn7iJYIHEHx8ZhBIE42L26-rSC_3l0ZaWEmsHAKyP9rgslApUdAQI")
// 	byteClientDataJSON := "{\"challenge\":\"E4PTcIH_HfX1pC6Sigk1SC9NAlgeztN0439vi8z_c9k\",\"new_keys_may_be_added_here\":\"do not compare clientDataJSON against a template. See https://goo.gl/yabPex\",\"origin\":\"https://webauthn.io\",\"type\":\"webauthn.get\"}"

// 	// type args struct {
// 	// 	response *http.Request
// 	// }
// 	tests := []struct {
// 		name    string
// 		args    types.CredentialAssertion
// 		want    types.CredentialAssertion
// 		wantErr bool
// 	}{
// 		{
// 			name: "Successfully Parse Credential Assertion",
// 			args: args{
// 				httpReq,
// 			},
// 			want: &ParsedCredentialAssertionData{
// 				ParsedPublicKeyCredential: ParsedPublicKeyCredential{
// 					ParsedCredential: ParsedCredential{
// 						ID:   byteID,
// 						Type: "public-key",
// 					},
// 					RawID: byteID,
// 				},
// 				Response: ParsedAssertionResponse{
// 					CollectedClientData: clientdata.CollectedClientData{
// 						Type:      clientdata.CeremonyType("webauthn.get"),
// 						Challenge: hex.MustBase64ToHash("E4PTcIH_HfX1pC6Sigk1SC9NAlgeztN0439vi8z_c9k"),
// 						Origin:    "https://webauthn.io",
// 						Hint:      "do not compare clientDataJSON against a template. See https://goo.gl/yabPex",
// 					},
// 					AuthenticatorData: AuthenticatorData{
// 						RPIDHash: byteRPIDHash,
// 						Counter:  1553097241,
// 						Flags:    0x045,
// 						AttData: AttestedCredentialData{
// 							AAGUID:              byteAAGUID,
// 							CredentialID:        byteID,
// 							CredentialPublicKey: byteCredentialPubKey,
// 						},
// 					},
// 					Signature: byteSignature,
// 					SessionID: byteUserHandle,
// 				},
// 				Raw: CredentialAssertionResponse{
// 					PublicKeyCredential: PublicKeyCredential{
// 						Credential: Credential{
// 							Type: "public-key",
// 							ID:   byteID,
// 						},
// 						RawID: byteID,
// 					},
// 					AssertionResponse: AuthenticatorAssertionResponse{
// 						AuthenticatorResponse: AuthenticatorResponse{
// 							UTF8ClientDataJSON: byteClientDataJSON,
// 						},
// 						AuthenticatorData: byteAuthData,
// 						Signature:         byteSignature,
// 						SessionID:         byteUserHandle,
// 					},
// 				},
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := ParseCredentialRequestResponse(tt.args.response)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ParseCredentialRequestResponse() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got.ClientExtensionResults, tt.want.ClientExtensionResults) {
// 				t.Errorf("ClientExtensionResults = %v \n want: %v", got, tt.want)
// 			}
// 			if !reflect.DeepEqual(got.ID, tt.want.ID) {
// 				t.Errorf("ID = %v \n want: %v", got.ID, tt.want.ID)
// 			}
// 			if !reflect.DeepEqual(got.ParsedCredential, tt.want.ParsedCredential) {
// 				t.Errorf("ParsedCredential = %v \n want: %v", got, tt.want)
// 			}
// 			if !reflect.DeepEqual(got.ParsedPublicKeyCredential, tt.want.ParsedPublicKeyCredential) {
// 				t.Errorf("ParsedPublicKeyCredential = %v \n want: %v", got.ParsedPublicKeyCredential.ClientExtensionResults, tt.want.ParsedPublicKeyCredential.ClientExtensionResults)
// 			}
// 			if !reflect.DeepEqual(got.Raw, tt.want.Raw) {
// 				t.Errorf("Raw = %+v \n want: %+v", got.Raw, tt.want.Raw)
// 			}
// 			if !reflect.DeepEqual(got.RawID, tt.want.RawID) {
// 				t.Errorf("RawID = %v \n want: %v", got, tt.want)
// 			}
// 			if !reflect.DeepEqual(got.Response, tt.want.Response) {
// 				var pkInterfaceMismatch bool
// 				if !reflect.DeepEqual(got.Response.CollectedClientData, tt.want.Response.CollectedClientData) {
// 					t.Errorf("Collected Client Data = %v \n want: %v", got.Response.CollectedClientData, tt.want.Response.CollectedClientData)
// 				}
// 				if !reflect.DeepEqual(got.Response.Signature, tt.want.Response.Signature) {
// 					t.Errorf("Signature = %v \n want: %v", got.Response.Signature, tt.want.Response.Signature)
// 				}
// 				if !reflect.DeepEqual(got.Response.AuthenticatorData.AttData.CredentialPublicKey, tt.want.Response.AuthenticatorData.AttData.CredentialPublicKey) {
// 					// Unmarshall CredentialPublicKey
// 					var pkWant interface{}
// 					keyBytesWant := tt.want.Response.AuthenticatorData.AttData.CredentialPublicKey
// 					webauthncbor.Unmarshal(keyBytesWant, &pkWant)
// 					var pkGot interface{}
// 					keyBytesGot := got.Response.AuthenticatorData.AttData.CredentialPublicKey
// 					webauthncbor.Unmarshal(keyBytesGot, &pkGot)
// 					if !reflect.DeepEqual(pkGot, pkWant) {
// 						t.Errorf("Response = %+v \n want: %+v", pkGot, pkWant)
// 					} else {
// 						pkInterfaceMismatch = true
// 					}
// 				}
// 				if pkInterfaceMismatch {
// 					return
// 				} else {
// 					t.Errorf("Response = %+v \n want: %+v", got.Response, tt.want.Response)
// 				}
// 			}
// 		})
// 	}
// }

func TestParsedCredentialAssertionData_Verify(t *testing.T) {

	tests := []struct {
		name    string
		args    types.VerifyAssertionInputArgs
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := zerolog.New(zerolog.NewConsoleWriter()).Level(zerolog.TraceLevel).With().Caller().Logger().WithContext(context.Background())
			if err := assertion.VerifyAssertionInput(ctx, tt.args); (err != nil) != tt.wantErr {
				t.Errorf("ParsedCredentialAssertionData.Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

var testAssertionOptions = map[string]string{
	// None Attestation - MacOS TouchID
	`success`: fmt.Sprintf(`{
		"id":"%s",
		"rawId":"%s",
		"type":"public-key",
		"response":{
			"attestationObject":"%s",
			"clientDataJSON":"{\"challenge\":\"fyeuuGP9zuej2FGjevi79bzqMKWxi4PYIa_5wj261W0\",\"origin\":\"https://webauthn.io\",\"type\":\"webauthn.create\"}"}
		}`,
		hex.MustBase64ToHash("AI7D5q2P0LS-Fal9ZT7CHM2N5BLbUunF92T8b6iYC199bO2kagSuU05-5dZGqb1SP0A0lyTWng").Hex(),
		hex.MustBase64ToHash("AI7D5q2P0LS-Fal9ZT7CHM2N5BLbUunF92T8b6iYC199bO2kagSuU05-5dZGqb1SP0A0lyTWng").Hex(),
		hex.MustBase64ToHash("o2NmbXRkbm9uZWdhdHRTdG10oGhhdXRoRGF0YVi7dKbqkhPJnC90siSSsyDPQCYqlMGpUKA5fyklC2CEHvBFXJJiFa3OAAI1vMYKZIsLJfHwVQMANwCOw-atj9C0vhWpfWU-whzNjeQS21Lpxfdk_G-omAtffWztpGoErlNOfuXWRqm9Uj9ANJck1p6lAQIDJiABIVggKAhfsdHcBIc0KPgAcRyAIK_-Vi-nCXHkRHPNaCMBZ-4iWCBxB8fGYQSBONi9uvq0gv95dGWlhJrBwCsj_a4LJQKVHQ").Hex(),
	),
}

var testAssertionResponses = map[string]string{
	// None Attestation - MacOS TouchID
	`success`: fmt.Sprintf(`{
		"id":"%s",
		"rawId":"%s",
		"type":"public-key",
		"response":{
			"authenticatorData":"%s",
			"clientDataJSON":"{\"challenge\":\"E4PTcIH_HfX1pC6Sigk1SC9NAlgeztN0439vi8z_c9k\",\"new_keys_may_be_added_here\":\"do not compare clientDataJSON against a template. See https://goo.gl/yabPex\",\"origin\":\"https://webauthn.io\",\"type\":\"webauthn.get\"}",
			"signature":"%s",
			"sessionId":"%s"
		}}`,
		hex.MustBase64ToHash("AI7D5q2P0LS-Fal9ZT7CHM2N5BLbUunF92T8b6iYC199bO2kagSuU05-5dZGqb1SP0A0lyTWng").Hex(),
		hex.MustBase64ToHash("AI7D5q2P0LS-Fal9ZT7CHM2N5BLbUunF92T8b6iYC199bO2kagSuU05-5dZGqb1SP0A0lyTWng").Hex(),
		hex.MustBase64ToHash("dKbqkhPJnC90siSSsyDPQCYqlMGpUKA5fyklC2CEHvBFXJJiGa3OAAI1vMYKZIsLJfHwVQMANwCOw-atj9C0vhWpfWU-whzNjeQS21Lpxfdk_G-omAtffWztpGoErlNOfuXWRqm9Uj9ANJck1p6lAQIDJiABIVggKAhfsdHcBIc0KPgAcRyAIK_-Vi-nCXHkRHPNaCMBZ-4iWCBxB8fGYQSBONi9uvq0gv95dGWlhJrBwCsj_a4LJQKVHQ").Hex(),
		hex.MustBase64ToHash("MEUCIBtIVOQxzFYdyWQyxaLR0tik1TnuPhGVhXVSNgFwLmN5AiEAnxXdCq0UeAVGWxOaFcjBZ_mEZoXqNboY5IkQDdlWZYc").Hex(),
		hex.MustBase64ToHash("0ToAAAAAAAAAAA"),
	),
}
