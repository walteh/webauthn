package lambda

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/walteh/webauthn/pkg/snsable"
)

type SQSEvent[I any] struct {
	Records []SQSEntity[I] `json:"Records"`
}

type SQSEntityViaSNS[I any] SQSEntity[SNSEntity[I]]

type SQSEventViaSNS[I any] struct {
	Records []SQSEntityViaSNS[I] `json:"Records"`
}

type SQSEventResponse struct {
	BatchItemFailures []BatchItemFailure `json:"batchItemFailures"`
}

type SNSEvent[I any] struct {
	Records []SNSEntity[I] `json:"Records"`
}

type sqsMessageString SQSEntity[string]

// An Amazon SQS message.
type SQSEntity[I any] struct {
	Attributes        map[string]string            `json:"attributes"`
	Body              I                            `json:"body"`
	MD5OfBody         *string                      `json:"md5OfBody"`
	MessageAttributes map[string]SQSAttributeValue `json:"messageAttributes"`
	MessageId         *string                      `json:"messageId"`
	ReceiptHandle     *string                      `json:"receiptHandle"`
	EventSourceArn    *string                      `json:"eventSourceArn"`
	EventSource       *string                      `json:"eventSource"`
	AwsRegion         *string                      `json:"awsRegion"`
}

func (me *SQSEntity[I]) UnmarshalJSON(b []byte) error {
	var g sqsMessageString
	if err := json.Unmarshal(b, &g); err != nil {
		return err
	}

	var bdy I
	if err := json.Unmarshal([]byte(g.Body), &bdy); err != nil {
		return err
	}

	*me = SQSEntity[I]{
		Attributes:        g.Attributes,
		Body:              bdy,
		MD5OfBody:         g.MD5OfBody,
		MessageAttributes: g.MessageAttributes,
		MessageId:         g.MessageId,
		ReceiptHandle:     g.ReceiptHandle,
	}

	return nil
}

func (me SQSEntity[I]) MarshalJSON() ([]byte, error) {

	g := sqsMessageString{
		Attributes:        me.Attributes,
		MD5OfBody:         me.MD5OfBody,
		MessageAttributes: me.MessageAttributes,
		MessageId:         me.MessageId,
		ReceiptHandle:     me.ReceiptHandle,
	}

	bdy, err := json.Marshal(me.Body)
	if err != nil {
		return nil, err
	}
	g.Body = string(bdy)

	return json.Marshal(g)
}

type SNSAttributeValue struct {

	// Amazon SNS supports the following logical data types: String, String.Array,
	// Number, and Binary. For more information, see Message Attribute Data Types
	// (https://docs.aws.amazon.com/sns/latest/dg/SNSMessageAttributes.html#SNSMessageAttributes.DataTypes).
	//
	// This member is required.
	DataType *string `json:"Type"`

	// Binary type attributes can store any binary data, for example, compressed data,
	// encrypted data, or images.
	BinaryValue []byte `json:"-,omitempty"`

	// Strings are Unicode with UTF8 binary encoding. For a list of code values, see
	// ASCII Printable Characters
	// (https://en.wikipedia.org/wiki/ASCII#ASCII_printable_characters).
	StringValue *string `json:"Value"`
}

type SNSEntity[I any] struct {
	Signature         string                       `json:"Signature"`
	MessageID         string                       `json:"MessageId"`
	Type              string                       `json:"Type"`
	TopicArn          string                       `json:"TopicArn"` //nolint: stylecheck
	MessageAttributes map[string]SNSAttributeValue `json:"MessageAttributes"`
	SignatureVersion  string                       `json:"SignatureVersion"`
	Timestamp         time.Time                    `json:"Timestamp"`
	SigningCertURL    string                       `json:"SigningCertUrl"`
	Message           I                            `json:"Message"`
	UnsubscribeURL    string                       `json:"UnsubscribeUrl"`
	Subject           string                       `json:"Subject"`
}

type snsMessageString SNSEntity[string]

func (me *SNSEntity[I]) UnmarshalJSON(b []byte) error {

	var g snsMessageString
	if err := json.Unmarshal(b, &g); err != nil {
		// try to unmarshal as string
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(s), &g); err != nil {
			return err
		}
	}

	var def struct {
		Def string `json:"default"`
	}
	if err := json.Unmarshal([]byte(g.Message), &def); err == nil && def.Def != "" {
		g.Message = def.Def
	}

	var bdy I
	if err := json.Unmarshal([]byte(g.Message), &bdy); err != nil {
		return err
	}

	*me = SNSEntity[I]{
		Signature:         g.Signature,
		MessageID:         g.MessageID,
		Type:              g.Type,
		TopicArn:          g.TopicArn,
		MessageAttributes: g.MessageAttributes,
		SignatureVersion:  g.SignatureVersion,
		Timestamp:         g.Timestamp,
		SigningCertURL:    g.SigningCertURL,
		UnsubscribeURL:    g.UnsubscribeURL,
		Subject:           g.Subject,
		Message:           bdy,
	}

	return nil
}

func (me SNSEntity[I]) MarshalJSON() ([]byte, error) {

	g := snsMessageString{
		Signature:         me.Signature,
		MessageID:         me.MessageID,
		Type:              me.Type,
		TopicArn:          me.TopicArn,
		MessageAttributes: me.MessageAttributes,
		SignatureVersion:  me.SignatureVersion,
		Timestamp:         me.Timestamp,
		SigningCertURL:    me.SigningCertURL,
		UnsubscribeURL:    me.UnsubscribeURL,
		Subject:           me.Subject,
	}

	bdy, err := json.Marshal(me.Message)
	if err != nil {
		return nil, err
	}

	bdy, err = json.Marshal(map[string]string{"default": string(bdy)})
	if err != nil {
		return nil, err
	}

	g.Message = string(bdy)

	return json.Marshal(g)
}

type SNSNotificationMetadata struct {
	ApproximateFirstReceiveTimestamp string `json:"ApproximateFirstReceiveTimestamp"`
	ApproximateReceiveCount          string `json:"ApproximateReceiveCount"`
	MessageDeduplicationId           string `json:"MessageDeduplicationId"`
	MessageGroupId                   string `json:"MessageGroupId"`
	SenderId                         string `json:"SenderId"`
	SentTimestamp                    string `json:"SentTimestamp"`
	SequenceNumber                   string `json:"SequenceNumber"`
}

func (me *SQSEntityViaSNS[I]) SNSMetadata() SNSNotificationMetadata {
	return SNSNotificationMetadata{
		ApproximateFirstReceiveTimestamp: me.Attributes["ApproximateFirstReceiveTimestamp"],
		ApproximateReceiveCount:          me.Attributes["ApproximateReceiveCount"],
		MessageDeduplicationId:           me.Attributes["MessageDeduplicationId"],
		MessageGroupId:                   me.Attributes["MessageGroupId"],
		SenderId:                         me.Attributes["SenderId"],
		SentTimestamp:                    me.Attributes["SentTimestamp"],
		SequenceNumber:                   me.Attributes["SequenceNumber"],
	}
}

func (me *SQSEntityViaSNS[I]) SNSAttributes() map[string]SNSAttributeValue {
	return me.Body.MessageAttributes
}

type SQSAttributeValue = sqstypes.MessageAttributeValue

func SQSEntityFromMessage[I any](mes *sqstypes.Message, queuearn string) (*SQSEntity[I], error) {

	var bdy I
	if err := json.Unmarshal([]byte(*mes.Body), &bdy); err != nil {
		return nil, err
	}

	return &SQSEntity[I]{
		Attributes:        mes.Attributes,
		Body:              bdy,
		MD5OfBody:         mes.MD5OfBody,
		MessageId:         mes.MessageId,
		ReceiptHandle:     mes.ReceiptHandle,
		MessageAttributes: mes.MessageAttributes,
		EventSourceArn:    aws.String(queuearn),
		EventSource:       aws.String("aws:sqs"),
		AwsRegion:         aws.String("unknown"),
	}, nil
}

func SNSEntityFromSnsable[I snsable.Snsable](sns I) SQSEntityViaSNS[I] {
	_, pub := snsable.SnsableToPublishInput(sns)

	newattrs := make(map[string]SNSAttributeValue, len(pub.MessageAttributes))

	for k, v := range pub.MessageAttributes {
		newattrs[k] = SNSAttributeValue{
			DataType:    v.DataType,
			StringValue: v.StringValue,
		}
	}

	subj := ""

	if pub.Subject != nil {
		subj = *pub.Subject
	}

	snss := SNSEntity[I]{
		Signature:         "",
		MessageID:         "",
		Type:              "",
		TopicArn:          "",
		MessageAttributes: newattrs,
		SignatureVersion:  "",
		Timestamp:         *aws.Time(time.Now()),
		SigningCertURL:    "",
		UnsubscribeURL:    "",
		Subject:           subj,
		Message:           sns,
	}

	return SQSEntityViaSNS[I]{
		Attributes: map[string]string{
			"ApproximateFirstReceiveTimestamp": "",
			"ApproximateReceiveCount":          "",
			"MessageDeduplicationId":           "",
			"MessageGroupId":                   "",
			"SenderId":                         "",
			"SentTimestamp":                    "",
			"SequenceNumber":                   "",
		},
		Body:              snss,
		MD5OfBody:         aws.String(""),
		MessageId:         aws.String(""),
		ReceiptHandle:     aws.String(""),
		MessageAttributes: nil,
		EventSourceArn:    pub.TopicArn,
		EventSource:       aws.String("aws:sns"),
		AwsRegion:         aws.String("unknown"),
	}
}

type XRayable interface {
	XRayTraceID() []string
}

func (me SQSEntity[I]) XRayTraceID() []string {
	return []string{me.Attributes["AWSTraceHeader"]}
}

func (me SQSEvent[I]) XRayTraceID() []string {
	var traces []string
	for _, r := range me.Records {
		traces = append(traces, r.XRayTraceID()...)
	}
	return traces
}

func (me SQSEventViaSNS[I]) XRayTraceID() []string {
	var traces []string
	for _, r := range me.Records {
		traces = append(traces, r.XRayTraceID()...)
	}
	return traces
}

func (me SQSEntityViaSNS[I]) XRayTraceID() []string {
	return []string{me.Attributes["AWSTraceHeader"]}
}
