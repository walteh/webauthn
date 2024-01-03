package credential_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/walteh/webauthn/pkg/hex"
	"github.com/walteh/webauthn/pkg/webauthn/credential"
	"github.com/walteh/webauthn/pkg/webauthn/providers"
	"github.com/walteh/webauthn/pkg/webauthn/types"
	"github.com/walteh/webauthn/pkg/webauthn/webauthncbor"
)

var byteID = types.CredentialID(hex.MustBase64ToHash("6xrtBhJQW6QU4tOaB4rrHaS2Ks0yDDL_q8jDC16DEjZ-VLVf4kCRkvl2xp2D71sTPYns-exsHQHTy3G-zJRK8g"))
var byteClientDataJSON = `{\"challenge\":\"W8GzFU8pGjhoRbWrLDlamAfq_y4S1CZG1VuoeRLARrE\",\"origin\":\"https://webauthn.io\",\"type\":\"webauthn.create\"}`
var byteAttObject = hex.MustBase64ToHash("o2NmbXRkbm9uZWdhdHRTdG10oGhhdXRoRGF0YVjEdKbqkhPJnC90siSSsyDPQCYqlMGpUKA5fyklC2CEHvBBAAAAAAAAAAAAAAAAAAAAAAAAAAAAQOsa7QYSUFukFOLTmgeK6x2ktirNMgwy_6vIwwtegxI2flS1X-JAkZL5dsadg-9bEz2J7PnsbB0B08txvsyUSvKlAQIDJiABIVggLKF5xS0_BntttUIrm2Z2tgZ4uQDwllbdIfrrBMABCNciWCDHwin8Zdkr56iSIh0MrB5qZiEzYLQpEOREhMUkY6q4Vw")
var clientDataJson = "{\"challenge\":\"W8GzFU8pGjhoRbWrLDlamAfq_y4S1CZG1VuoeRLARrE\",\"origin\":\"https://webauthn.io\",\"type\":\"webauthn.create\"}"

func TestParseCredentialCreationResponse(t *testing.T) {

	// reqBody := ioutil.NopCloser(bytes.NewReader([]byte(testCredentialRequestBody)))
	// httpReq := &http.Request{Body: reqBody}
	// type args struct {
	// 	response *http.Request
	// }

	byteAuthData := hex.MustBase64ToHash("dKbqkhPJnC90siSSsyDPQCYqlMGpUKA5fyklC2CEHvBBAAAAAAAAAAAAAAAAAAAAAAAAAAAAQOsa7QYSUFukFOLTmgeK6x2ktirNMgwy_6vIwwtegxI2flS1X-JAkZL5dsadg-9bEz2J7PnsbB0B08txvsyUSvKlAQIDJiABIVggLKF5xS0_BntttUIrm2Z2tgZ4uQDwllbdIfrrBMABCNciWCDHwin8Zdkr56iSIh0MrB5qZiEzYLQpEOREhMUkY6q4Vw")
	byteRPIDHash := hex.MustBase64ToHash("dKbqkhPJnC90siSSsyDPQCYqlMGpUKA5fyklC2CEHvA")
	byteCredentialPubKey := hex.MustBase64ToHash("pSJYIMfCKfxl2SvnqJIiHQysHmpmITNgtCkQ5ESExSRjqrhXAQIDJiABIVggLKF5xS0_BntttUIrm2Z2tgZ4uQDwllbdIfrrBMABCNc")

	tests := []struct {
		name    string
		args    types.AttestationInput
		want    *types.AttestationObject
		wantErr bool
	}{
		{
			name: "Successful Credential Request Parsing",
			args: types.AttestationInput{
				UTF8ClientDataJSON: clientDataJson,
				AttestationObject:  byteAttObject,
				CredentialID:       byteID,
				CredentialType:     "public-key",
				ClientExtensions: map[string]interface{}{
					"appid": true,
				},
			},
			want: &types.AttestationObject{
				ClientData: types.CollectedClientData{
					Type:      types.CeremonyType("webauthn.create"),
					Challenge: types.CeremonyID(hex.MustBase64ToHash("W8GzFU8pGjhoRbWrLDlamAfq_y4S1CZG1VuoeRLARrE")),
					Origin:    "https://webauthn.io",
				},
				AuthData: types.AuthenticatorData{
					RPIDHash: byteRPIDHash,
					Counter:  0,
					Flags:    0x041,
					AttData: types.AttestedCredentialData{
						AAGUID:              make([]byte, 16),
						CredentialID:        byteID,
						CredentialPublicKey: byteCredentialPubKey,
					},
					ExtData: nil,
				},

				Format:       "none",
				RawAuthData:  byteAuthData,
				AttStatement: map[string]interface{}{},
				Extensions: map[string]interface{}{
					"appid": true,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()

			got, err := credential.ParseAttestationInput(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCredentialCreationResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got.AttStatement, tt.want.AttStatement) {
				t.Errorf("Extensions = %v \n want: %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.AuthData.AttData.AAGUID, tt.want.AuthData.AttData.AAGUID) {
				t.Errorf("Extensions = %v \n want: %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.AuthData.AttData.CredentialID, tt.want.AuthData.AttData.CredentialID) {
				t.Errorf("Extensions = %v \n want: %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.AuthData.Counter, tt.want.AuthData.Counter) {
				t.Errorf("Extensions = %v \n want: %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.AuthData.ExtData, tt.want.AuthData.ExtData) {
				t.Errorf("Extensions = %v \n want: %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.AuthData.Flags, tt.want.AuthData.Flags) {
				t.Errorf("Extensions = %v \n want: %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.AuthData.RPIDHash, tt.want.AuthData.RPIDHash) {
				t.Errorf("Extensions = %v \n want: %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.AuthData.RPIDHash, tt.want.AuthData.RPIDHash) {
				t.Errorf("Extensions = %v \n want: %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.ClientData, tt.want.ClientData) {
				t.Errorf("Extensions = %v \n want: %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.Extensions, tt.want.Extensions) {
				t.Errorf("Extensions = %v \n want: %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.Format, tt.want.Format) {
				t.Errorf("Extensions = %v \n want: %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.RawAuthData, tt.want.RawAuthData) {
				t.Errorf("Extensions = %v \n want: %v", got, tt.want)
			}

			// Unmarshall CredentialPublicKey
			var pkWant interface{}
			keyBytesWant := tt.want.AuthData.AttData.CredentialPublicKey
			webauthncbor.Unmarshal(keyBytesWant, &pkWant)
			var pkGot interface{}
			keyBytesGot := got.AuthData.AttData.CredentialPublicKey
			webauthncbor.Unmarshal(keyBytesGot, &pkGot)
			if !reflect.DeepEqual(pkGot, pkWant) {
				t.Errorf("Response = %+v \n want: %+v", pkGot, pkWant)
			}

			// if !reflect.DeepEqual(got.Type, tt.want.Type) {
			// 	t.Errorf("Type = %v \n want: %v", got, tt.want)
			// }
			// if !reflect.DeepEqual(got.Response.CollectedClientData, tt.want.Response.CollectedClientData) {
			// 	t.Errorf("CollectedClientData = %v \n want: %v", got, tt.want)
			// }
			// if !reflect.DeepEqual(got.Response.AttestationObject.Format, tt.want.Response.AttestationObject.Format) {
			// 	t.Errorf("Format = %v \n want: %v", got, tt.want)
			// }
			// if !reflect.DeepEqual(got.Response.AttestationObject.AuthData.AttData.CredentialID, tt.want.Response.AttestationObject.AuthData.AttData.CredentialID) {
			// 	t.Errorf("CredentialID = %v \n want: %v", got, tt.want)
			// }
		})
	}
}

func TestParsedCredentialCreationData_Verify(t *testing.T) {
	// byteID := hex.MustBase64ToHash("6xrtBhJQW6QU4tOaB4rrHaS2Ks0yDDL_q8jDC16DEjZ-VLVf4kCRkvl2xp2D71sTPYns-exsHQHTy3G-zJRK8g")
	byteChallenge := hex.MustBase64ToHash("W8GzFU8pGjhoRbWrLDlamAfq_y4S1CZG1VuoeRLARrE")
	// byteAuthData := hex.MustBase64ToHash("dKbqkhPJnC90siSSsyDPQCYqlMGpUKA5fyklC2CEHvBBAAAAAAAAAAAAAAAAAAAAAAAAAAAAQOsa7QYSUFukFOLTmgeK6x2ktirNMgwy_6vIwwtegxI2flS1X-JAkZL5dsadg-9bEz2J7PnsbB0B08txvsyUSvKlAQIDJiABIVggLKF5xS0_BntttUIrm2Z2tgZ4uQDwllbdIfrrBMABCNciWCDHwin8Zdkr56iSIh0MrB5qZiEzYLQpEOREhMUkY6q4Vw")
	// byteRPIDHash := hex.MustBase64ToHash("dKbqkhPJnC90siSSsyDPQCYqlMGpUKA5fyklC2CEHvA")
	// byteCredentialPubKey := hex.MustBase64ToHash("pSJYIMfCKfxl2SvnqJIiHQysHmpmITNgtCkQ5ESExSRjqrhXAQIDJiABIVggLKF5xS0_BntttUIrm2Z2tgZ4uQDwllbdIfrrBMABCNc")
	byteAttObject := hex.MustBase64ToHash("o2NmbXRkbm9uZWdhdHRTdG10oGhhdXRoRGF0YVjEdKbqkhPJnC90siSSsyDPQCYqlMGpUKA5fyklC2CEHvBBAAAAAAAAAAAAAAAAAAAAAAAAAAAAQOsa7QYSUFukFOLTmgeK6x2ktirNMgwy_6vIwwtegxI2flS1X-JAkZL5dsadg-9bEz2J7PnsbB0B08txvsyUSvKlAQIDJiABIVggLKF5xS0_BntttUIrm2Z2tgZ4uQDwllbdIfrrBMABCNciWCDHwin8Zdkr56iSIh0MrB5qZiEzYLQpEOREhMUkY6q4Vw")
	// byteClientDataJSON := hex.MustBase64ToHash("eyJjaGFsbGVuZ2UiOiJXOEd6RlU4cEdqaG9SYldyTERsYW1BZnFfeTRTMUNaRzFWdW9lUkxBUnJFIiwib3JpZ2luIjoiaHR0cHM6Ly93ZWJhdXRobi5pbyIsInR5cGUiOiJ3ZWJhdXRobi5jcmVhdGUifQ")

	tests := []struct {
		name    string
		args    types.VerifyAttestationInputArgs
		wantErr bool
	}{
		{
			name: "Successful Verification Test",
			args: types.VerifyAttestationInputArgs{
				Provider: providers.NewNoneAttestationProvider(),
				Input: types.AttestationInput{
					UTF8ClientDataJSON: clientDataJson,
					AttestationObject:  byteAttObject,
				},
				StoredChallenge:    types.CeremonyID(byteChallenge),
				VerifyUser:         false,
				RelyingPartyID:     `webauthn.io`,
				RelyingPartyOrigin: `https://webauthn.io`,
				SessionId:          hex.Hash{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()

			if _, err := credential.VerifyAttestationInput(ctx, tt.args); (err != nil) != tt.wantErr {
				t.Errorf("ParsedCredentialCreationData.Verify() error = %+v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

var testCredentialRequestBody = fmt.Sprintf(
	`{"id":"%s","rawId":"%s","type":"public-key","clientExtensionResults":{"appid":true},"response":{"attestationObject":"%s","clientDataJSON":"%s"}}`,
	byteID.Ref().Hex(), byteID.Ref().Hex(), byteAttObject.Hex(), byteClientDataJSON)

// &{{0x74a6ea9213c99c2f74b22492b320cf40262a94c1a950a0397f29250b60841ef0 65 0 {0x00000000000000000000000000000000 0xeb1aed0612505ba414e2d39a078aeb1da4b62acd320c32ffabc8c30b5e8312367e54b55fe2409192f976c69d83ef5b133d89ecf9ec6c1d01d3cb71becc944af2 0xa50102032620012158202ca179c52d3f067b6db5422b9b6676b60678b900f09656dd21faeb04c00108d7225820c7c229fc65d92be7a892221d0cac1e6a66213360b42910e44484c52463aab857} 0x} 0x74a6ea9213c99c2f74b22492b320cf40262a94c1a950a0397f29250b60841ef04100000000000000000000000000000000000000000040eb1aed0612505ba414e2d39a078aeb1da4b62acd320c32ffabc8c30b5e8312367e54b55fe2409192f976c69d83ef5b133d89ecf9ec6c1d01d3cb71becc944af2a50102032620012158202ca179c52d3f067b6db5422b9b6676b60678b900f09656dd21faeb04c00108d7225820c7c229fc65d92be7a892221d0cac1e6a66213360b42910e44484c52463aab857 none map[] {webauthn.create 0x5bc1b3154f291a386845b5ab2c395a9807eaff2e12d42646d55ba87912c046b1 https://webauthn.io <nil> } map[]}
// &{{0x74a6ea9213c99c2f74b22492b320cf40262a94c1a950a0397f29250b60841ef0 65 0 {0x00000000000000000000000000000000 0xeb1aed0612505ba414e2d39a078aeb1da4b62acd320c32ffabc8c30b5e8312367e54b55fe2409192f976c69d83ef5b133d89ecf9ec6c1d01d3cb71becc944af2 0xa5225820c7c229fc65d92be7a892221d0cac1e6a66213360b42910e44484c52463aab8570102032620012158202ca179c52d3f067b6db5422b9b6676b60678b900f09656dd21faeb04c00108d7} 0x} 0x74a6ea9213c99c2f74b22492b320cf40262a94c1a950a0397f29250b60841ef04100000000000000000000000000000000000000000040eb1aed0612505ba414e2d39a078aeb1da4b62acd320c32ffabc8c30b5e8312367e54b55fe2409192f976c69d83ef5b133d89ecf9ec6c1d01d3cb71becc944af2a50102032620012158202ca179c52d3f067b6db5422b9b6676b60678b900f09656dd21faeb04c00108d7225820c7c229fc65d92be7a892221d0cac1e6a66213360b42910e44484c52463aab857 none map[] {webauthn.create 0x5bc1b3154f291a386845b5ab2c395a9807eaff2e12d42646d55ba87912c046b1 https://webauthn.io <nil> } map[appid:true]}

// &{{0x74a6ea9213c99c2f74b22492b320cf40262a94c1a950a0397f29250b60841ef0 65 0 {0x00000000000000000000000000000000 0xeb1aed0612505ba414e2d39a078aeb1da4b62acd320c32ffabc8c30b5e8312367e54b55fe2409192f976c69d83ef5b133d89ecf9ec6c1d01d3cb71becc944af2 0xa5225820c7c229fc65d92be7a892221d0cac1e6a66213360b42910e44484c52463aab8570102032620012158202ca179c52d3f067b6db5422b9b6676b60678b900f09656dd21faeb04c00108d7} 0x} 0x74a6ea9213c99c2f74b22492b320cf40262a94c1a950a0397f29250b60841ef04100000000000000000000000000000000000000000040eb1aed0612505ba414e2d39a078aeb1da4b62acd320c32ffabc8c30b5e8312367e54b55fe2409192f976c69d83ef5b133d89ecf9ec6c1d01d3cb71becc944af2a50102032620012158202ca179c52d3f067b6db5422b9b6676b60678b900f09656dd21faeb04c00108d7225820c7c229fc65d92be7a892221d0cac1e6a66213360b42910e44484c52463aab857 none map[] {webauthn.create 0x5bc1b3154f291a386845b5ab2c395a9807eaff2e12d42646d55ba87912c046b1 https://webauthn.io <nil> } map[appid:true]}
