package signinwithapple

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang-jwt/jwt/v4"
)

type MockClient struct {
	keyfunc jwt.Keyfunc
}

func (me *MockClient) BuildKeyFunc() (jwt.Keyfunc, error) {
	return me.keyfunc, nil
}

var DummyApplePublicKeyResponse = &PublicKeyResponse{
	Keys: []ApplePublicKey{
		{KTY: "RSA",
			KID: "fh6Bs8C",
			Use: "sig",
			Alg: "RS256",
			N:   "u704gotMSZc6CSSVNCZ1d0S9dZKwO2BVzfdTKYz8wSNm7R_KIufOQf3ru7Pph1FjW6gQ8zgvhnv4IebkGWsZJlodduTC7c0sRb5PZpEyM6PtO8FPHowaracJJsK1f6_rSLstLdWbSDXeSq7vBvDu3Q31RaoV_0YlEzQwPsbCvD45oVy5Vo5oBePUm4cqi6T3cZ-10gr9QJCVwvx7KiQsttp0kUkHM94PlxbG_HAWlEZjvAlxfEDc-_xZQwC6fVjfazs3j1b2DZWsGmBRdx1snO75nM7hpyRRQB4jVejW9TuZDtPtsNadXTr9I5NjxPdIYMORj9XKEh44Z73yfv0gtw",
			E:   "AQAB"},
	}}

func TestJwtApple(t *testing.T) {
	type args struct {
		ctx   context.Context
		token string
	}

	tests := []struct {
		name    string
		args    args
		want    *jwt.Token
		wantErr bool
	}{
		{
			name: "A",
			args: args{
				ctx:   context.Background(),
				token: "eyJraWQiOiJmaDZCczhDIiwiYWxnIjoiUlMyNTYifQ.eyJpc3MiOiJodHRwczovL2FwcGxlaWQuYXBwbGUuY29tIiwiYXVkIjoieHl6Lm51Z2cuYXBwIiwiZXhwIjoxNjY3OTI0NzEwLCJpYXQiOjE2Njc4MzgzMTAsInN1YiI6IjAwMTQzNy5kZWY1MzVkZGQ5ZTI0YzRmYTQzNjdkY2E1MGZkZmVkYi4xOTUxIiwiY19oYXNoIjoidTFPdHRlUDh6Ym5vM2lnNVZpd2owUSIsImF1dGhfdGltZSI6MTY2NzgzODMxMCwibm9uY2Vfc3VwcG9ydGVkIjp0cnVlfQ.H1vlDXhXXXj_OQFwyzklrMu3r1qQupK_M5Ot6lbIpHNK9eB8WQHCjvD2dND-ov_eNFWae8gIKRxH4ev_lMfgZFc-IomcB4aqld5uIbEy53VInfCJ9sGNyAjt31kZidLdasqTSSzXxtKbBf89gtiZqHuqHmWF0rSPqs7xM1x2CHfvQYu5642COq8vfwcajbXdYsBLEBnfsSyKVikD7_6Ggl0fwUpixOgGg0i2syNEm8uJ5eCTi0K5k3fEhNESxV7A-voKw-wERNQnxsh1Isr-6s2YHXMhWR2iHLSDn1-H-k2sKgkqTE7P20Y2BAlfAvXccctUtfHAplS7CbFJ0Zs3pw",
			},
			want: &jwt.Token{
				Raw: "eyJraWQiOiJmaDZCczhDIiwiYWxnIjoiUlMyNTYifQ.eyJpc3MiOiJodHRwczovL2FwcGxlaWQuYXBwbGUuY29tIiwiYXVkIjoieHl6Lm51Z2cuYXBwIiwiZXhwIjoxNjY3OTI0NzEwLCJpYXQiOjE2Njc4MzgzMTAsInN1YiI6IjAwMTQzNy5kZWY1MzVkZGQ5ZTI0YzRmYTQzNjdkY2E1MGZkZmVkYi4xOTUxIiwiY19oYXNoIjoidTFPdHRlUDh6Ym5vM2lnNVZpd2owUSIsImF1dGhfdGltZSI6MTY2NzgzODMxMCwibm9uY2Vfc3VwcG9ydGVkIjp0cnVlfQ.H1vlDXhXXXj_OQFwyzklrMu3r1qQupK_M5Ot6lbIpHNK9eB8WQHCjvD2dND-ov_eNFWae8gIKRxH4ev_lMfgZFc-IomcB4aqld5uIbEy53VInfCJ9sGNyAjt31kZidLdasqTSSzXxtKbBf89gtiZqHuqHmWF0rSPqs7xM1x2CHfvQYu5642COq8vfwcajbXdYsBLEBnfsSyKVikD7_6Ggl0fwUpixOgGg0i2syNEm8uJ5eCTi0K5k3fEhNESxV7A-voKw-wERNQnxsh1Isr-6s2YHXMhWR2iHLSDn1-H-k2sKgkqTE7P20Y2BAlfAvXccctUtfHAplS7CbFJ0Zs3pw",
				Method: &jwt.SigningMethodRSA{
					Name: "RS256",
					Hash: 0x5,
				},
				Header: map[string]interface{}{
					"kid": "fh6Bs8C",
					"alg": "RS256",
				},
				Claims: jwt.MapClaims{
					"nonce_supported": true,
					"iss":             "https://appleid.apple.com",
					"aud":             "xyz.nugg.app",
					"exp":             1667924710.000000,
					"iat":             1667838310.000000,
					"sub":             "001437.def535ddd9e24c4fa4367dca50fdfedb.1951",
					"c_hash":          "u1OtteP8zbno3ig5Viwj0Q",
					"auth_time":       1667838310.000000,
				},
				Signature: "H1vlDXhXXXj_OQFwyzklrMu3r1qQupK_M5Ot6lbIpHNK9eB8WQHCjvD2dND-ov_eNFWae8gIKRxH4ev_lMfgZFc-IomcB4aqld5uIbEy53VInfCJ9sGNyAjt31kZidLdasqTSSzXxtKbBf89gtiZqHuqHmWF0rSPqs7xM1x2CHfvQYu5642COq8vfwcajbXdYsBLEBnfsSyKVikD7_6Ggl0fwUpixOgGg0i2syNEm8uJ5eCTi0K5k3fEhNESxV7A-voKw-wERNQnxsh1Isr-6s2YHXMhWR2iHLSDn1-H-k2sKgkqTE7P20Y2BAlfAvXccctUtfHAplS7CbFJ0Zs3pw",
				Valid:     true,
			},
			wantErr: false,
		},
		{
			name: "A",
			args: args{
				ctx:   context.Background(),
				token: "eyJraWQiOiJmaDZCczhDIiwiYWxnIjoiUlMyNTYifQ.eyJpc3MiOiJodHRwczovL2FwcGxlaWQuYXBwbGUuY29tIiwiYXVkIjoieHl6Lm51Z2cuYXBwIiwiZXhwIjoxNjY3OTI0NzEwLCJpYXQiOjE2Njc4MzgzMTAsInN1YiI6IjAwMTQzNy5kZWY1MzVkZGQ5ZTI0YzRmYTQzNjdkY2E1MGZkZmVkYi4xOTUxIiwiY19oYXNoIjoidTFPdHRlUDh6Ym5vM2lnNVZpd2owUSIsImF1dGhfdGltZSI6MTY2NzgzODMxMCwibm9uY2Vfc3VwcG9ydGVkIjp0cnVlfQ.H1vlDXhXXXj_OQFwyzklrMu3r1qQupK_M5Ot6lbIpHNK9eB8WQHCjvD2dND-ov_eNFWae8gIKRxH4ev_lMfgZFc-IomcB4aqld5uIbEy53VInfCJ9sGNyAjt31kZidLdasqTSSzXxtKbBf89gtiZqHuqHmWF0rSPqs7xM1x2CHfvQYu5642COq8vfwcajbXdYsBLEBnfsSyKVikD7_6Ggl0fwUpixOgGg0i2syNEm8uJ5eCTi0K5k3fEhNESxV7A-voKw-wERNQnxsh1Isr-6s2YHXMhWR2iHLSDn1-H-k2sKgkqTE7P20Y2BAlfAvXccctUtfHAplS7CbFJ0Zs3pw",
			},
			want: &jwt.Token{
				Raw: "eyJraWQiOiJmaDZCczhDIiwiYWxnIjoiUlMyNTYifQ.eyJpc3MiOiJodHRwczovL2FwcGxlaWQuYXBwbGUuY29tIiwiYXVkIjoieHl6Lm51Z2cuYXBwIiwiZXhwIjoxNjY3OTI0NzEwLCJpYXQiOjE2Njc4MzgzMTAsInN1YiI6IjAwMTQzNy5kZWY1MzVkZGQ5ZTI0YzRmYTQzNjdkY2E1MGZkZmVkYi4xOTUxIiwiY19oYXNoIjoidTFPdHRlUDh6Ym5vM2lnNVZpd2owUSIsImF1dGhfdGltZSI6MTY2NzgzODMxMCwibm9uY2Vfc3VwcG9ydGVkIjp0cnVlfQ.H1vlDXhXXXj_OQFwyzklrMu3r1qQupK_M5Ot6lbIpHNK9eB8WQHCjvD2dND-ov_eNFWae8gIKRxH4ev_lMfgZFc-IomcB4aqld5uIbEy53VInfCJ9sGNyAjt31kZidLdasqTSSzXxtKbBf89gtiZqHuqHmWF0rSPqs7xM1x2CHfvQYu5642COq8vfwcajbXdYsBLEBnfsSyKVikD7_6Ggl0fwUpixOgGg0i2syNEm8uJ5eCTi0K5k3fEhNESxV7A-voKw-wERNQnxsh1Isr-6s2YHXMhWR2iHLSDn1-H-k2sKgkqTE7P20Y2BAlfAvXccctUtfHAplS7CbFJ0Zs3pw",
				Method: &jwt.SigningMethodRSA{
					Name: "RS256",
					Hash: 0x5,
				},
				Header: map[string]interface{}{
					"kid": "fh6Bs8C",
					"alg": "RS256",
				},
				Claims: jwt.MapClaims{
					"nonce_supported": true,
					"iss":             "https://appleid.apple.com",
					"aud":             "xyz.nugg.app",
					"exp":             1667924710.000000,
					"iat":             1667838310.000000,
					"sub":             "001437.def535ddd9e24c4fa4367dca50fdfedb.1951",
					"c_hash":          "u1OtteP8zbno3ig5Viwj0Q",
					"auth_time":       1667838310.000000,
				},
				Signature: "H1vlDXhXXXj_OQFwyzklrMu3r1qQupK_M5Ot6lbIpHNK9eB8WQHCjvD2dND-ov_eNFWae8gIKRxH4ev_lMfgZFc-IomcB4aqld5uIbEy53VInfCJ9sGNyAjt31kZidLdasqTSSzXxtKbBf89gtiZqHuqHmWF0rSPqs7xM1x2CHfvQYu5642COq8vfwcajbXdYsBLEBnfsSyKVikD7_6Ggl0fwUpixOgGg0i2syNEm8uJ5eCTi0K5k3fEhNESxV7A-voKw-wERNQnxsh1Isr-6s2YHXMhWR2iHLSDn1-H-k2sKgkqTE7P20Y2BAlfAvXccctUtfHAplS7CbFJ0Zs3pw",
				Valid:     true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DummyApplePublicKeyResponse.ParseToken(tt.args.token)

			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}
