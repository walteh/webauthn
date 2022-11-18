package dynamo

// func TestFindInOnePerTableGetResult(t *testing.T) {

// 	type run[I interface{}] struct {
// 		name    string
// 		args    []*GetOutput
// 		want    *I
// 		wantErr bool
// 	}
// 	tests := []run[Ceremony]{
// 		{
// 			name: "A",
// 			args: []*GetOutput{
// 				// {
// 				// 	TableName: "test",
// 				// 	Item: &Ceremony{
// 				// 		UserId: "test",
// 				// 	},
// 				// },
// 			},

// 			// want: &Ceremony{
// 			// 	UserId: "test",
// 			// },
// 			wantErr: false,
// 		},

// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := FindInOnePerTableGetResult(tt.args.result, tt.args.tableName)
// 			// if (err != nil) != tt.wantErr {
// 			// 	t.Errorf("FindInOnePerTableGetResult() error = %v, wantErr %v", err, tt.wantErr)
// 			// 	return
// 			// }
// 			// if !reflect.DeepEqual(got, tt.want) {
// 			// 	t.Errorf("FindInOnePerTableGetResult() = %v, want %v", got, tt.want)
// 			// }
// 		})
// 	}
// }
