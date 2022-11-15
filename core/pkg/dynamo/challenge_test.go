package dynamo

import (
	"nugg-auth/core/pkg/safeid"
	"reflect"
	"testing"
)

func Test_clientDataToSafeID(t *testing.T) {
	type args struct {
		clientData            string
		expectedChallengeType string
		expectedOrigin        string
	}
	tests := []struct {
		name    string
		args    args
		want    safeid.SafeID
		wantErr bool
	}{
		{
			name: "A",
			args: args{
				clientData:            `{"challenge":"01GHYQFHNYP06KNCQ1NCJW2TVK","origin":"https://localhost:8080","type":"webauthn.create"}`,
				expectedChallengeType: "webauthn.create",
				expectedOrigin:        "https://localhost:8080",
			},
			want:    safeid.MustParseStrict("01GHYQFHNYP06KNCQ1NCJW2TVK"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := clientDataToSafeID(tt.args.clientData, tt.args.expectedChallengeType, tt.args.expectedOrigin)
			if (err != nil) != tt.wantErr {
				t.Errorf("clientDataToSafeID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("clientDataToSafeID() = %v, want %v", *got, tt.want)
			}
		})
	}
}
