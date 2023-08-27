package lambda

import (
	"encoding/json"
	"regexp"

	"github.com/walteh/webauthn/pkg/awsepoch"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

// DynamoDBEventResponse is the outer structure to report batch item failures for DynamoDBEvent.
type DynamoDBStreamEventResponse struct {
	BatchItemFailures []BatchItemFailure `json:"batchItemFailures"`
}

// The DynamoDBStreamEvent stream event handled to Lambda
// http://docs.aws.amazon.com/lambda/latest/dg/eventsources.html#eventsources-ddb-update
type DynamoDBStreamEvent struct {
	Records []DynamoDBStreamEventRecord `json:"Records"`
}

// DynamoDBStreamEventRecord stores information about each record of a DynamoDB stream event
type DynamoDBStreamEventRecord struct {
	// The region in which the GetRecords request was received.
	AWSRegion string `json:"awsRegion"`

	// The main body of the stream record, containing all of the DynamoDB-specific
	// fields.
	Change DynamoDBStreamRecord `json:"dynamodb"`

	// A globally unique identifier for the event that was recorded in this stream
	// record.
	EventID string `json:"eventID"`

	// The type of data modification that was performed on the DynamoDB table:
	//
	//    * INSERT - a new item was added to the table.
	//
	//    * MODIFY - one or more of an existing item's attributes were modified.
	//
	//    * REMOVE - the item was deleted from the table
	EventName string `json:"eventName"`

	// The AWS service from which the stream record originated. For DynamoDB Streams,
	// this is aws:dynamodb.
	EventSource string `json:"eventSource"`

	// The version number of the stream record format. This number is updated whenever
	// the structure of Record is modified.
	//
	// Client applications must not assume that eventVersion will remain at a particular
	// value, as this number is subject to change at any time. In general, eventVersion
	// will only increase as the low-level DynamoDB Streams API evolves.
	EventVersion string `json:"eventVersion"`

	// The event source ARN of DynamoDB
	EventSourceArn string `json:"eventSourceARN"` //nolint: stylecheck

	// Items that are deleted by the Time to Live process after expiration have
	// the following fields:
	//
	//    * Records[].userIdentity.type
	//
	// "Service"
	//
	//    * Records[].userIdentity.principalId
	//
	// "dynamodb.amazonaws.com"
	UserIdentity *DynamoDBUserIdentity `json:"userIdentity,omitempty"`
}

func (me *DynamoDBStreamEventRecord) TableName() string {
	// get the table name from the dyanmo table arn
	// arn:aws:dynamodb:us-east-1:123456789012:table/ExampleTableWithStream/stream/2015-06-27T00:48:05.899
	// ExampleTableWithStream
	reg := regexp.MustCompile(`table/([^/]+)/stream`)
	matches := reg.FindStringSubmatch(me.EventSourceArn)
	if len(matches) == 2 {
		return matches[1]
	}
	return me.EventSourceArn
}

type DynamoDBUserIdentity struct {
	Type        string `json:"type"`
	PrincipalID string `json:"principalId"`
}

// DynamoDBStreamRecord represents a description of a single data modification that was performed on an item
// in a DynamoDB table.
type DynamoDBStreamRecord struct {

	// The approximate date and time when the stream record was created, in UNIX
	// epoch time (http://www.epochconverter.com/) format.
	ApproximateCreationDateTime awsepoch.AWSEpochTime `json:"ApproximateCreationDateTime,omitempty"`

	// The primary key attribute(s) for the DynamoDB item that was modified.
	Keys DynamoDBImage `json:"Keys,omitempty"`

	// The item in the types. table as it appeared after it was modified.
	NewImage DynamoDBImage `json:"NewImage,omitempty"`

	// The item in the types. table as it appeared before it was modified.
	OldImage DynamoDBImage `json:"OldImage,omitempty"`

	// The sequence number of the stream record.
	SequenceNumber string `json:"SequenceNumber"`

	// The size of the stream record, in bytes.
	SizeBytes int64 `json:"SizeBytes"`

	// The type of data from the modified DynamoDB item that was captured in this
	// stream record.
	StreamViewType string `json:"StreamViewType"`
}

type DynamoDBImage map[string]types.AttributeValue

func (me *DynamoDBImage) UnmarshalJSON(b []byte) error {

	m := &types.AttributeValueMemberM{
		Value: make(map[string]types.AttributeValue),
	}

	var c map[string]map[string]interface{}

	err := json.Unmarshal(b, &c)
	if err != nil {
		return err
	}
	for k, v := range c {
		xx, err := unmarshalDynamoDBAttributeValueMap(v)
		if err != nil {
			return err
		}
		m.Value[k] = xx
	}

	// m.Value =

	*me = m.Value

	return nil

}

func (me DynamoDBImage) AsAttributeValues() map[string]types.AttributeValue {
	return me
}

type DynamoDBKeyType string

const (
	DynamoDBKeyTypeHash  DynamoDBKeyType = "HASH"
	DynamoDBKeyTypeRange DynamoDBKeyType = "RANGE"
)

type DynamoDBOperationType string

const (
	DynamoDBOperationTypeInsert DynamoDBOperationType = "INSERT"
	DynamoDBOperationTypeModify DynamoDBOperationType = "MODIFY"
	DynamoDBOperationTypeRemove DynamoDBOperationType = "REMOVE"
)

type DynamoDBSharedIteratorType string

const (
	DynamoDBShardIteratorTypeTrimHorizon         DynamoDBSharedIteratorType = "TRIM_HORIZON"
	DynamoDBShardIteratorTypeLatest              DynamoDBSharedIteratorType = "LATEST"
	DynamoDBShardIteratorTypeAtSequenceNumber    DynamoDBSharedIteratorType = "AT_SEQUENCE_NUMBER"
	DynamoDBShardIteratorTypeAfterSequenceNumber DynamoDBSharedIteratorType = "AFTER_SEQUENCE_NUMBER"
)

type DynamoDBStreamStatus string

const (
	DynamoDBStreamStatusEnabling  DynamoDBStreamStatus = "ENABLING"
	DynamoDBStreamStatusEnabled   DynamoDBStreamStatus = "ENABLED"
	DynamoDBStreamStatusDisabling DynamoDBStreamStatus = "DISABLING"
	DynamoDBStreamStatusDisabled  DynamoDBStreamStatus = "DISABLED"
)

type DynamoDBStreamViewType string

const (
	DynamoDBStreamViewTypeNewImage        DynamoDBStreamViewType = "NEW_IMAGE"          // the entire item, as it appeared after it was modified.
	DynamoDBStreamViewTypeOldImage        DynamoDBStreamViewType = "OLD_IMAGE"          // the entire item, as it appeared before it was modified.
	DynamoDBStreamViewTypeNewAndOldImages DynamoDBStreamViewType = "NEW_AND_OLD_IMAGES" // both the new and the old item images of the item.
	DynamoDBStreamViewTypeKeysOnly        DynamoDBStreamViewType = "KEYS_ONLY"          // only the key attributes of the modified item.
)

// "{\"Message\":\"{\\\"default\\\":\\\"{\\\\\\\"address\\\\\\\":\\\\\\\"0x69420000bad0605988626169e32aa82fb3981add\\\\\\\",\\\\\\\"data\\\\\\\":\\\\\\\"0xc0000000059c0000000001f49b0e2b16f57648c7baf28edd7772a815af266e770000000000000000000000000446177900000000000000000bec0831040e000600000000000007d40000000097b2c73628cc125400000000027bdeecd5caedac\\\\\\\",\\\\\\\"blockNumber\\\\\\\":7900761,\\\\\\\"transactionHash\\\\\\\":\\\\\\\"0x488f4549c9286cc808d366ad65dd6d23db071f4f696270c8950c03de7155bafa\\\\\\\",\\\\\\\"transactionIndex\\\\\\\":3,\\\\\\\"blockHash\\\\\\\":\\\\\\\"0xca4882e8edab41322411e8a138b00932d8cf708aedbda7afcf109909d429fbcc\\\\\\\",\\\\\\\"logIndex\\\\\\\":3,\\\\\\\"removed\\\\\\\":false,\\\\\\\"chain\\\\\\\":1,\\\\\\\"origin\\\\\\\":\\\\\\\"\\\\\\\"}\\\"}\",\"MessageAttributes\":null,\"MessageDeduplicationId\":\"0xca4882e8edab41322411e8a138b00932d8cf708aedbda7afcf109909d429fbcc_3\",\"MessageGroupId\":\"nuggft-20221105\",\"MessageStructure\":null,\"PhoneNumber\":null,\"Subject\":null,\"TargetArn\":null,\"TopicArn\":\"\"}",

// "{\"Message\":\"{\\\"default\\\":\\\"{\\\"address\\\":\\\"0x69420000bad0605988626169e32aa82fb3981add\\\",\\\"data\\\":\\\"0xc0000000059c0000000001f49b0e2b16f57648c7baf28edd7772a815af266e770000000000000000000000000446177900000000000000000bec0831040e000600000000000007d40000000097b2c73628cc125400000000027bdeecd5caedac\\\",\\\"blockNumber\\\":7900761,\\\"transactionHash\\\":\\\"0x488f4549c9286cc808d366ad65dd6d23db071f4f696270c8950c03de7155bafa\\\",\\\"transactionIndex\\\":3,\\\"blockHash\\\":\\\"0xca4882e8edab41322411e8a138b00932d8cf708aedbda7afcf109909d429fbcc\\\",\\\"logIndex\\\":3,\\\"removed\\\":false,\\\"chain\\\":1,\\\"origin\\\":\\\"\\\"}\\\"}\",\"MessageAttributes\":null,\"MessageDeduplicationId\":\"0xca4882e8edab41322411e8a138b00932d8cf708aedbda7afcf109909d429fbcc_3\",\"MessageGroupId\":\"nuggft-20221105\",\"MessageStructure\":null,\"PhoneNumber\":null,\"Subject\":null,\"TargetArn\":null,\"TopicArn\":\"\"}"
