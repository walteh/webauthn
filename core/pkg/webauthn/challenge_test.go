package protocol

import (
	"nugg-webauthn/core/pkg/hex"
	"reflect"
	"testing"
)

func TestCreateChallenge(t *testing.T) {
	tests := []struct {
		name    string
		want    hex.Hash
		wantErr bool
	}{
		{
			"Successfull Challenge Create",
			hex.Hash{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateChallenge()
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateChallenge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.want = got
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateChallenge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChallenge_String(t *testing.T) {
	newChallenge, err := CreateChallenge()
	if err != nil {
		t.Errorf("CreateChallenge() error = %v", err)
		return
	}
	wantChallenge := newChallenge.Hex()
	tests := []struct {
		name string
		c    hex.Hash
		want string
	}{
		{
			"Successful String",
			newChallenge,
			wantChallenge,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("Challenge.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// "{\"id\":\"0xeb1aed0612505ba414e2d39a078aeb1da4b62acd320c32ffabc8c30b5e8312367e54b55fe2409192f976c69d83ef5b133d89ecf9ec6c1d01d3cb71becc944af2\",\"rawId\":\"0xeb1aed0612505ba414e2d39a078aeb1da4b62acd320c32ffabc8c30b5e8312367e54b55fe2409192f976c69d83ef5b133d89ecf9ec6c1d01d3cb71becc944af2\",\"type\":\"public-key\",\"clientExtensionResults\":{\"appid\":true},\"response\":{\"attestationObject\":\"0xa363666d74646e6f6e656761747453746d74a068617574684461746158c474a6ea9213c99c2f74b22492b320cf40262a94c1a950a0397f29250b60841ef041000000000000000000000
// 00000000000000000000040eb1aed0612505ba414e2d39a078aeb1da4b62acd320c32ffabc8c30b5e8312367e54b55fe2409192f976c69d83ef5b133d89ecf9ec6c1d01d3cb71becc944af2a50102032620012158202ca179c52d3f067b6db5422b9b6676b60678b900f09656dd21faeb04c00108d7225820c7c229fc65d92be7a892221d0cac1e6a66213360b42910e44484c52463aab857\",\"clientDataJSON\":\"{\"challenge\":\"0x5bc1b3154f291a386845b5ab2c395a9807eaff2e12d42646d55ba87912c046b1\",\"origin\":\"https://webauthn.io\",\"type\":\"webauthn.create\"}\"}}}"
