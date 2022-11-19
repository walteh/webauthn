package webauthncbor

import (
	"encoding/base64"
	"testing"
)

type AttestationObject struct {
	// The authenticator data, including the newly created public key. See AuthenticatorData for more info
	AuthData AuthenticatorData
	// The byteform version of the authenticator data, used in part for signature validation
	RawAuthData []byte `json:"authData" cbor:"authData"`
	// The format of the Attestation data.
	Format string `json:"fmt" cbor:"fmt"`
	// The attestation statement data sent back if attestation is requested.
	AttStatement map[string]interface{} `json:"attStmt,omitempty"`
}

type AuthenticatorFlags uint8

type AuthenticatorData struct {
	RPIDHash []byte                 `json:"rpid"`
	Flags    AuthenticatorFlags     `json:"flags"`
	Counter  uint32                 `json:"sign_count"`
	AttData  AttestedCredentialData `json:"att_data"`
	ExtData  []byte                 `json:"ext_data"`
}

type AttestedCredentialData struct {
	AAGUID       []byte `json:"aaguid"`
	CredentialID []byte `json:"credential_id"`
	// The raw credential public key bytes received from the attestation data
	CredentialPublicKey []byte `json:"public_key"`
}

func TestUnmarshal(t *testing.T) {
	type args struct {
		data []byte
		v    interface{}
	}

	r, err := base64.RawStdEncoding.DecodeString("o2NmbXRkbm9uZWdhdHRTdG10oGhhdXRoRGF0YViYqbmr9/xGsTVktJ1c+FvL83H5y2MODWs1S8YLUeBl2khdAAAAAAAAAAAAAAAAAAAAAAAAAAAAFAxyGEJnT2kSs7Q+SXoyzZZ5+2fFpQECAyYgASFYIBOLuWNlKh1gf4tejUV5dvIiEjfuQ2Ci5tIpuPkbNdVEIlggBt5P704/l0OzpaQ4/Oc+P7pE3Kvzd3L77L+T9Z970Yc")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestUnmarshal",
			args: args{
				data: r,
				v: &AttestationObject{
					AuthData: AuthenticatorData{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Unmarshal(tt.args.data, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
