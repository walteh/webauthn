package indexable

// func TestNoWay(t *testing.T) {
// 	type args struct {
// 		a interface{}
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want reflect.StructField
// 	}{
// 		{
// 			name: "test",
// 			args: args{
// 				a: &struct {
// 					ABC *string
// 				}{
// 					ABC: "test",
// 				}.ABC,
// 			},
// 			want: reflect.StructField{
// 				Name: "A",
// 				Type: reflect.TypeOf(""),
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := NoWay(tt.args.a); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("NoWay() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"google.golang.org/protobuf/proto"
)

var _ Indexable = (*SampleIndexable)(nil)
var _ Identifiable = (*SampleIndexable)(nil)

type SampleIndexable struct {
	Id string `dynamodbav:"ID"`
}

func (me *SampleIndexable) Combine([]Indexable) error {
	return nil
}

/*///////////////////////////////////////////////////////////////////
 ///                     STRUCT FIELDS 							 ///
///////////////////////////////////////////////////////////////////*/

/*///////////////////////////////////////////////////////////////////
 ///                     	DYNAMODB 							 ///
///////////////////////////////////////////////////////////////////*/

func (me SampleIndexable) MarshalDynamoDBAttributeValue() (av types.AttributeValue, err error) {
	type G SampleIndexable
	res, err := attributevalue.MarshalMap(G(me))
	if err != nil {
		return
	}

	IndexableMarshalIndexes(&me, res)

	return M(res), nil

}

func (me *SampleIndexable) ID() string {
	return me.Id + ":" + me.Id
}

func (me *SampleIndexable) ToProtobuf() proto.Message {
	return nil
}

func (me *SampleIndexable) IsWorthy(string, string) bool {
	return true
}

/*///////////////////////////////////////////////////////////////////
 ///                     	INDEXABLE 							 ///
///////////////////////////////////////////////////////////////////*/

func (me *SampleIndexable) PrimaryIndex() *DynamoDBIndex {
	return &DynamoDBIndex{
		HashKey:    "pk",
		HashValue:  S(me.ID()),
		RangeKey:   "sk",
		RangeValue: S(me.ID()),
	}
}

func (me *SampleIndexable) SecondaryIndexes() map[string]*DynamoDBIndex {
	return map[string]*DynamoDBIndex{}
}

func (me *SampleIndexable) ResolvableTableName() string {
	return "event"
}

func NewSampleIndexable(id string) *SampleIndexable {
	return &SampleIndexable{
		Id: id,
	}
}
