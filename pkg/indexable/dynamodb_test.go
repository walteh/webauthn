package indexable

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type fakeUnmarshaler struct {
	Id float64 `dynamodbav:"ID"`
}

func (f *fakeUnmarshaler) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) error {
	type G *fakeUnmarshaler

	if val, ok := av.(*types.AttributeValueMemberM); ok {
		return attributevalue.UnmarshalMap(val.Value, G(f))
	}
	//
	return errors.New("invalid type")

}

func TestNewDynamoDBAVMap(t *testing.T) {
	type args struct {
		out     []map[string]types.AttributeValue
		builder DynamoDBAVMapIDBuilder[*fakeUnmarshaler]
	}
	tests := []struct {
		name    string
		args    args
		want    *DynamoDBAVMap[*fakeUnmarshaler]
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				out: []map[string]types.AttributeValue{},
				builder: func(val *fakeUnmarshaler) string {
					return fmt.Sprintf("%v", val.Id)
				},
			},
			want: &DynamoDBAVMap[*fakeUnmarshaler]{
				b: map[string]bool{},
				m: map[string]*fakeUnmarshaler{},
			},
		},
		{
			name: "one",
			args: args{
				out: []map[string]types.AttributeValue{
					{
						"Id": &types.AttributeValueMemberN{Value: "1.0"},
					},
				},
				builder: func(val *fakeUnmarshaler) string {
					return fmt.Sprintf("%v", val.Id)
				},
			},
			want: &DynamoDBAVMap[*fakeUnmarshaler]{
				b: map[string]bool{
					"1": true,
				},
				m: map[string]*fakeUnmarshaler{
					"1": {
						Id: 1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDynamoDBAVMap(tt.args.out, tt.args.builder)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDynamoDBAVMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDynamoDBAVMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
