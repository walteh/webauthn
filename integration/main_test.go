package main

import (
	"io"
	"net/http"
	"nugg-auth/integration/pkg/terraform"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Output struct {
	// The name of the output
	AuthExpApiInvokeUrl string `json:"auth_exp_api_invoke_url"`
	ApiGatewayHost      string `json:"apigw_host"`
	AppsyncAuthorizer   string `json:"appsync_lambda_authorizer_function_name"`
}

func TestHttp(t *testing.T) {

	r := terraform.LoadOutput[Output](t)

	type args struct {
		endpoint string
		meathod  string
		body     io.Reader
		headers  map[string]string
	}
	type want struct {
		statusCode int
		body       io.Reader
		headers    map[string]string
	}

	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name:    "pass A",
			args:    args{endpoint: r.AuthExpApiInvokeUrl + "/challenge", meathod: http.MethodPost, body: nil, headers: map[string]string{"x-nugg-challenge-state": "abc123"}},
			want:    want{statusCode: http.StatusNoContent, body: nil, headers: map[string]string{"x-nugg-challenge": ""}},
			wantErr: false,
		},
		{
			name:    "pass B",
			args:    args{endpoint: r.AuthExpApiInvokeUrl + "/challenge", meathod: http.MethodPost, body: nil, headers: map[string]string{"X-Nugg-Challenge-State": "abc123"}},
			want:    want{statusCode: http.StatusNoContent, body: nil, headers: map[string]string{"x-nugg-challenge": ""}},
			wantErr: false,
		},
		{
			name:    "fail A",
			args:    args{endpoint: r.AuthExpApiInvokeUrl + "/challenge", meathod: http.MethodPost, body: nil, headers: map[string]string{"x-nugg-challenge-stat": "abc123"}},
			want:    want{statusCode: http.StatusBadRequest, body: nil},
			wantErr: false,
		},
		{
			name:    "fail B",
			args:    args{endpoint: r.AuthExpApiInvokeUrl + "/challenge", meathod: http.MethodPost, body: nil, headers: map[string]string{}},
			want:    want{statusCode: http.StatusBadRequest, body: nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *http.Request
			var err error

			if got, err = http.NewRequest(tt.args.meathod, tt.args.endpoint, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("http.NewRequest() error = %v, wantErr %v", err, tt.wantErr)
			}

			for k, v := range tt.args.headers {
				got.Header.Set(k, v)
			}

			var resp *http.Response

			if resp, err = http.DefaultClient.Do(got); (err != nil) != tt.wantErr {
				t.Errorf("http.DefaultClient.Do() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.want.statusCode, resp.StatusCode, resp)

			for k, v := range tt.want.headers {
				if v == "" {
					assert.NotEmpty(t, resp.Header.Get(k), k)
				} else {
					assert.Equal(t, v, resp.Header.Get(k), k)
				}
			}
		})
	}

}

func TestAppsyncAuthorizer(t *testing.T) {

	// r := terraform.LoadOutput[Output](t)

	// endpoint, realtime, teardown := terraform.BuildAppsyncExample(t, r.AppsyncAuthorizer)

	// defer teardown()

	// type args struct {
	// 	endpoint string
	// 	meathod  string
	// 	body     io.Reader
	// 	headers  map[string]string
	// }

	// tests := []struct {
	// 	name    string
	// 	args    args
	// 	want    int
	// 	wantErr bool
	// }{
	// 	{
	// 		name:    "pass A",
	// 		args:    args{endpoint: endpoint, meathod: http.MethodPost, body: nil, headers: map[string]string{"x-nugg-challenge-state": "abc123"}},
	// 		want:    http.StatusNoContent,
	// 		wantErr: false,
	// 	},
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		var got *http.Request
	// 		var err error

	// 		if got, err = http.NewRequest(tt.args.meathod, tt.args.endpoint, tt.args.body); (err != nil) != tt.wantErr {
	// 			t.Errorf("http.NewRequest() error = %v, wantErr %v", err, tt.wantErr)
	// 		}

	// 		for k, v := range tt.args.headers {
	// 			got.Header.Set(k, v)
	// 		}

	// 		var resp *http.Response

	// 		if resp, err = realtime.Do(got); (err != nil) != tt.wantErr {
	// 			t.Errorf("http.DefaultClient.Do() error = %v, wantErr %v", err, tt.wantErr)
	// 		}

	// 		assert.Equal(t, tt.want, resp.StatusCode, resp)
	// 	})
	// }

}
