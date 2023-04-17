package providers

import (
	"reflect"
	"testing"

	"git.nugg.xyz/webauthn/pkg/webauthn/types"
)

func Test_verifyPackedFormat(t *testing.T) {
	type args struct {
		att            types.AttestationObject
		clientDataHash []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []interface{}
		wantErr bool
	}{
		// {
		// 	name: "Successful Self Attestation",
		// 	args: args{
		// 		att: AttestationObject
		// 	}
		// }
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got, got1, err := NewPackedAttestationProvider().Attest(tt.args.att, tt.args.clientDataHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("verifyPackedFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("verifyPackedFormat() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("verifyPackedFormat() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

var testPackedAttestationOptions = []string{}

var testPackedAttestationResponses = []string{}
