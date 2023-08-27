package lambda

import (
	"encoding/json"
	"testing"

	"github.com/k0kubun/pp/v3"
	"github.com/stretchr/testify/assert"
)

func TestSqsMessage_UnmarshalJSON(t *testing.T) {

	type anything struct {
		A string
	}

	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  SQSEntity[anything]
		args    args
		wantErr bool
	}{
		{
			name: "TestDynamoDBStreamRecord_UnmarshalJSON",
			args: args{
				b: []byte("{\"Body\": \"{\\\"A\\\": \\\"b\\\"}\"}"),
			},
			fields: SQSEntity[anything]{
				Body: anything{A: "b"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			me := SQSEntity[anything]{}

			if err := json.Unmarshal(tt.args.b, &me); (err != nil) != tt.wantErr {
				t.Errorf("DynamoDBStreamRecord.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, me, tt.fields)

		})
	}
}

func TestSqsMessage_UnmarshalJSONReal(t *testing.T) {

	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "A",
			args: args{
				b: []byte("{\"Records\":[{\"messageId\":\"6c79ced4-001c-4d93-b42f-80457e1c818c\",\"receiptHandle\":\"AQEBbDptQOMRWLLScvkJ/QFigNxkRHlzM/bryb3UAp6P2XuV8hLp1OFgY66gnZQ9YBL+mrXyLIEXhcPh0xrIO2iX+N54oRZBbOyAxozTYQo+IZ7iEy1O5jzCplA9iZl7WlkJ3TfiPlth3PDoIlpozFCu1Dup+qsArK1Uh0cxWr4BaYqu6TWbB9EXLh2BftW6YxjqXPR8olaoh6yUesjNkwmoWkm30xwKOg1eeqW4o1sr3Tpzj/Cp+6RioA6Nxe4KYENuO26cPxNzTs3qi/4IdNp0rudr3VL2KiCb5oG+wWOzuAqsefsCTadmB/YWTcwc6Imm\",\"body\":\"{\\n  \\\"Type\\\" : \\\"Notification\\\",\\n  \\\"MessageId\\\" : \\\"f08bff8b-203a-57f2-84b1-a3c4ac185312\\\",\\n  \\\"SequenceNumber\\\" : \\\"10000000000002316000\\\",\\n  \\\"TopicArn\\\" : \\\"arn:aws:sns:us-east-1:876249556319:us-dev6-bloom-sns.fifo\\\",\\n  \\\"Message\\\" : \\\"{\\\\\\\"default\\\\\\\":\\\\\\\"{\\\\\\\\\\\\\\\"address\\\\\\\\\\\\\\\":\\\\\\\\\\\\\\\"0x69420000bad0605988626169e32aa82fb3981add\\\\\\\\\\\\\\\",\\\\\\\\\\\\\\\"data\\\\\\\\\\\\\\\":\\\\\\\\\\\\\\\"0xc0000002af140000000001f49b0e2b16f57648c7baf28edd7772a815af266e7700000000000000000000000007fa1b7b00000000000000000bd807e00424001000000000000007ec00000000a11a2171dc6d55990000000009bbc65584a1aa67\\\\\\\\\\\\\\\",\\\\\\\\\\\\\\\"blockNumber\\\\\\\\\\\\\\\":8249663,\\\\\\\\\\\\\\\"transactionHash\\\\\\\\\\\\\\\":\\\\\\\\\\\\\\\"0x17c118aa01d3141ed786298a5d40d3d7845f97bb4740352df492de1bae92dbf9\\\\\\\\\\\\\\\",\\\\\\\\\\\\\\\"transactionIndex\\\\\\\\\\\\\\\":33,\\\\\\\\\\\\\\\"blockHash\\\\\\\\\\\\\\\":\\\\\\\\\\\\\\\"0xe3995dbf5cab729690795020277e61bc7a7eb37b3213d32c84f30dd20387e424\\\\\\\\\\\\\\\",\\\\\\\\\\\\\\\"logIndex\\\\\\\\\\\\\\\":20,\\\\\\\\\\\\\\\"removed\\\\\\\\\\\\\\\":false,\\\\\\\\\\\\\\\"chain\\\\\\\\\\\\\\\":1,\\\\\\\\\\\\\\\"origin\\\\\\\\\\\\\\\":\\\\\\\\\\\\\\\"\\\\\\\\\\\\\\\"}\\\\\\\"}\\\",\\n  \\\"Timestamp\\\" : \\\"2023-01-28T20:47:30.270Z\\\",\\n  \\\"UnsubscribeURL\\\" : \\\"https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:876249556319:us-dev6-bloom-sns.fifo:116c3a38-7454-4d07-9000-3cc36c4b8a0e\\\",\\n  \\\"MessageAttributes\\\" : {\\n    \\\"i_want_to_be_persisted\\\" : {\\\"Type\\\":\\\"String\\\",\\\"Value\\\":\\\"true\\\"}\\n  }\\n}\",\"attributes\":{\"ApproximateReceiveCount\":\"1\",\"SentTimestamp\":\"1674938850293\",\"SequenceNumber\":\"18875528419384560128\",\"MessageGroupId\":\"nuggft-20221105\",\"SenderId\":\"AIDAYRRVD2ENU4DSO2WBX\",\"MessageDeduplicationId\":\"0xe3995dbf5cab729690795020277e61bc7a7eb37b3213d32c84f30dd20387e424_20\",\"ApproximateFirstReceiveTimestamp\":\"1674938850293\"},\"messageAttributes\":{},\"md5OfBody\":\"06e7e5f067e92d064544c802b1d3ca19\",\"eventSource\":\"aws:sqs\",\"eventSourceARN\":\"arn:aws:sqs:us-east-1:876249556319:us-dev6-bloom-sns-buffer-persist.fifo\",\"awsRegion\":\"us-east-1\"}]}"),
			},
			wantErr: false,
		},
		{
			name: "B",
			args: args{
				b: []byte("{\"Records\":[{\"messageId\":\"af248e30-49e1-414b-8376-7395a6fa2dcf\",\"receiptHandle\":\"AQEBe+kQ1cn3aPZewsr84IzgpNbKqRHVTDhiMovdk/7pVF9TL/d84xgdYQMzZy9ZqniXZsMFQMISIrBPQ3IWZC1qnVwaAQPG7pzuo0Nh4Fg01MF4UuZsuaRMCBMCkQJaVMGnlONP6NCFGF6OkJWan4h4j3678N9sIL2tCrEttS866n8wO6K4CeLkC2UnxHGU+7R/rSZunf+gS4kEoM9v4IIpIKiMBlX1J4FfQDxWbX7L5itlXx2fx28Hnb6iFTdUm2rVtK11Yn/80xMmw8VXJeW7VcQZYwfRm1dbGLAfeTrh+1CBVS+/DKyeMUQcZFKJMVs9\",\"body\":\"{\\n  \\\"Type\\\" : \\\"Notification\\\",\\n  \\\"MessageId\\\" : \\\"0a2b2237-314c-57c0-92ca-7e78be2aab27\\\",\\n  \\\"TopicArn\\\" : \\\"arn:aws:sns:us-east-1:876249556319:us-dev6-bloom-sns.fifo\\\",\\n  \\\"Subject\\\" : \\\"ethrpc-header\\\",\\n  \\\"Message\\\" : \\\"{\\\\\\\"parentHash\\\\\\\":\\\\\\\"0xd89332844476d1b29805ca9e8fc56a49c9fe66d890df556b9a1297c52bf212be\\\\\\\",\\\\\\\"sha3Uncles\\\\\\\":\\\\\\\"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347\\\\\\\",\\\\\\\"number\\\\\\\":8393537,\\\\\\\"gasLimit\\\\\\\":30000000,\\\\\\\"gasUsed\\\\\\\":29992464,\\\\\\\"timestamp\\\\\\\":1674939684,\\\\\\\"baseFeePerGas\\\\\\\":\\\\\\\"1/62500000000000000\\\\\\\",\\\\\\\"hash\\\\\\\":\\\\\\\"0x14e8f1a5bc270bce69b25b63decaca10df523ce878b32e5863d4548605b5254b\\\\\\\",\\\\\\\"Origin\\\\\\\":\\\\\\\"nuggft-20221105\\\\\\\",\\\\\\\"chain\\\\\\\":0}\\\",\\n  \\\"Timestamp\\\" : \\\"2023-01-28T21:01:26.372Z\\\",\\n  \\\"UnsubscribeURL\\\" : \\\"https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:876249556319:us-dev6-bloom-sns.fifo:116c3a38-7454-4d07-9000-3cc36c4b8a0e\\\"\\n}\",\"attributes\":{\"ApproximateReceiveCount\":\"1\",\"SentTimestamp\":\"1674939686410\",\"SequenceNumber\":\"18875528633430512128\",\"MessageGroupId\":\"nuggft-20221105\",\"SenderId\":\"AIDAYRRVD2ENU4DSO2WBX\",\"MessageDeduplicationId\":\"0x14e8f1a5bc270bce69b25b63decaca10df523ce878b32e5863d4548605b5254b-8393537\",\"ApproximateFirstReceiveTimestamp\":\"1674939686410\"},\"messageAttributes\":{},\"md5OfBody\":\"f03b00faaecdd6d3b1a6f48dcce65adf\",\"eventSource\":\"aws:sqs\",\"eventSourceARN\":\"arn:aws:sqs:us-east-1:876249556319:us-dev6-bloom-sns-buffer-persist.fifo\",\"awsRegion\":\"us-east-1\"}]}"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			me := SQSEventViaSNS[interface{}]{}

			if err := json.Unmarshal(tt.args.b, &me); (err != nil) != tt.wantErr {
				t.Errorf("DynamoDBStreamRecord.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			pp.Println(me)

			// g := map[string]interface{}{}
			// if err := json.Unmarshal(tt.args.b, &g); (err != nil) != tt.wantErr {
			// 	t.Errorf("DynamoDBStreamRecord.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			// }
			// pp.Println(g)

		})
	}
}
