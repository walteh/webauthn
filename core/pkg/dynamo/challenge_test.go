package dynamo

// func Test_clientDataToSafeID(t *testing.T) {
// 	type args struct {
// 		clientData string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    string
// 		wantErr bool
// 	}{
// 		{
// 			name: "A",
// 			args: args{
// 				clientData: `{"challenge":"QVlSOWJxa29wSXJKVU9wVXRndzRDZw","origin":"https://localhost:8080","type":"webauthn.create"}`,
// 			},
// 			want:    "QVlSOWJxa29wSXJKVU9wVXRndzRDZw",
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := clientDataToSafeID(tt.args.clientData)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("clientDataToSafeID() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("clientDataToSafeID() = %v, want %v", got, tt.want)
// 			}
// 			log.Println(got)
// 		})
// 	}
// }
