package dynamo

import (
	"log"
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
		want    string
		wantErr bool
	}{
		{
			name: "A",
			args: args{
				clientData:            `{"challenge":"QVlSOWJxa29wSXJKVU9wVXRndzRDZw","origin":"https://localhost:8080","type":"webauthn.create"}`,
				expectedChallengeType: "webauthn.create",
				expectedOrigin:        "https://localhost:8080",
			},
			want:    "QVlSOWJxa29wSXJKVU9wVXRndzRDZw",
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("clientDataToSafeID() = %v, want %v", got, tt.want)
			}
			log.Println(got)
		})
	}
}
