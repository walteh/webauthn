package lambda

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/walteh/webauthn/pkg/awsepoch"
)

func TestDynamoDBStreamRecord_UnmarshalJSON(t *testing.T) {

	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  DynamoDBStreamRecord
		args    args
		wantErr bool
	}{
		{
			name: "TestDynamoDBStreamRecord_UnmarshalJSON",
			args: args{
				b: []byte("{\"eventID\":\"10f8c554826ee6973d3541738134e0a0\",\"eventName\":\"INSERT\",\"eventVersion\":\"1.1\",\"eventSource\":\"aws:dynamodb\",\"awsRegion\":\"us-east-1\",\"dynamodb\":{\"ApproximateCreationDateTime\":1.674921798E9,\"Keys\":{\"group_id\":{\"S\":\"nuggft-20221105\"},\"deduplication_id\":{\"S\":\"0xca4882e8edab41322411e8a138b00932d8cf708aedbda7afcf109909d429fbcc_3\"}},\"NewImage\":{\"group_id\":{\"S\":\"nuggft-20221105\"},\"payload\":{\"S\":\"abc\"},\"deduplication_id\":{\"S\":\"0xca4882e8edab41322411e8a138b00932d8cf708aedbda7afcf109909d429fbcc_3\"},\"ttl\":{\"N\":\"1674922098\"}},\"SequenceNumber\":\"3977200000000040943643552\",\"SizeBytes\":1119,\"StreamViewType\":\"NEW_AND_OLD_IMAGES\"},\"eventSourceARN\":\"arn:aws:dynamodb:us-east-1:876249556319:table/dev6-sns-buffer/stream/2023-01-27T18:59:24.761\"}"),
			},
			fields: DynamoDBStreamRecord{
				ApproximateCreationDateTime: awsepoch.AWSEpochTime{Time: time.Unix(1674921798, 0)},
				Keys: map[string]types.AttributeValue{
					"group_id":         &types.AttributeValueMemberS{Value: "nuggft-20221105"},
					"deduplication_id": &types.AttributeValueMemberS{Value: "0xca4882e8edab41322411e8a138b00932d8cf708aedbda7afcf109909d429fbcc_3"},
				},
				NewImage: map[string]types.AttributeValue{
					"group_id":         &types.AttributeValueMemberS{Value: "nuggft-20221105"},
					"payload":          &types.AttributeValueMemberS{Value: `abc`},
					"deduplication_id": &types.AttributeValueMemberS{Value: "0xca4882e8edab41322411e8a138b00932d8cf708aedbda7afcf109909d429fbcc_3"},
					"ttl":              &types.AttributeValueMemberN{Value: "1674922098"},
				},
				OldImage:       nil,
				SequenceNumber: "3977200000000040943643552",
				SizeBytes:      1119,
				StreamViewType: "NEW_AND_OLD_IMAGES",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			me := &DynamoDBStreamEventRecord{}

			if err := json.Unmarshal(tt.args.b, &me); (err != nil) != tt.wantErr {
				t.Errorf("DynamoDBStreamRecord.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, me.Change, tt.fields)

		})
	}
}
