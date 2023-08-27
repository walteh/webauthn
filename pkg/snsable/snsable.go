package snsable

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/rs/zerolog/log"
)

type Snsable interface {
	MarshalJSON() ([]byte, error)
	SNSAttributes() map[string][]string
	AddSNSAttribute(key, value string)
}

func SnsableToPublishInput[I Snsable](event I) (ok bool, res *sns.PublishInput) {
	var err error

	bpayload, err := event.MarshalJSON()
	if err != nil {
		log.Error().Err(err).Interface("event", event).Msg("error getting sns payload, skipping")
		return false, nil
	}

	payload := string(bpayload)

	// these two fields are required
	if payload == "" {
		log.Error().Err(err).Interface("event", event).Msg("error getting sns payload, skipping")
		return false, nil
	}

	// you get Invalid parameter: Message Structure - No default entry in JSON message body
	// if you don't do this
	pl2 := make(map[string]string)
	pl2["default"] = payload

	bpayload, err = json.Marshal(pl2)
	if err != nil {
		log.Error().Err(err).Interface("event", event).Msg("error getting sns payload, skipping")
		return false, nil
	}

	attrs := make(map[string]types.MessageAttributeValue)

	for k, v := range event.SNSAttributes() {
		attrs[k] = types.MessageAttributeValue{
			DataType:    aws.String("String.Array"),
			StringValue: aws.String(EncodeSnsableArrayAttribute(v)),
		}
	}

	attrs["payload_type"] = types.MessageAttributeValue{
		DataType:    aws.String("String"),
		StringValue: aws.String(reflect.TypeOf(event).String()),
	}

	pi := &sns.PublishInput{
		Message:           aws.String(string(bpayload)),
		MessageAttributes: attrs,
		MessageStructure:  aws.String("json"),
		Subject:           aws.String(reflect.TypeOf(event).String()),
	}

	return true, pi
}

func EncodeSnsableArrayAttribute(value []string) string {
	return "[\"" + strings.Join(value, "\",\"") + "\"]"
}

func DecodeSnsableArrayAttribute(value string) []string {
	return strings.Split(strings.TrimPrefix(strings.TrimSuffix(value, "\"]"), "[\""), "\",\"")
}

func SnsableToPublishBatchInput[I Snsable](events []I) (*sns.PublishBatchInput, error) {

	batch := &sns.PublishBatchInput{
		PublishBatchRequestEntries: []types.PublishBatchRequestEntry{},
	}

	for id, event := range events {
		ok, input := SnsableToPublishInput(event)
		if !ok {
			return nil, fmt.Errorf("error getting sns payload, skipping")
		}

		batch.PublishBatchRequestEntries = append(batch.PublishBatchRequestEntries, types.PublishBatchRequestEntry{
			Message:                input.Message,
			MessageAttributes:      input.MessageAttributes,
			MessageStructure:       input.MessageStructure,
			MessageDeduplicationId: input.MessageDeduplicationId,
			MessageGroupId:         input.MessageGroupId,
			Subject:                input.Subject,
			Id:                     aws.String(fmt.Sprintf("%d", id)),
		})
	}
	return batch, nil
}

func PublishInputToSnsableTest[I Snsable](in *sns.PublishInput) I {
	var err error
	var event I

	payload := make(map[string]string)
	err = json.Unmarshal([]byte(*in.Message), &payload)
	if err != nil {
		return *new(I)
	}

	err = json.Unmarshal([]byte(payload["default"]), &event)
	if err != nil {
		return *new(I)
	}

	for k, v := range in.MessageAttributes {
		if k == "payload_type" {
			continue
		}
		for _, xx := range DecodeSnsableArrayAttribute(*v.StringValue) {
			event.AddSNSAttribute(k, xx)
		}
	}

	return event

}

var _ WithContext = (*SnsableWithContext[Snsable])(nil)

type WithContext interface {
	Context() context.Context
}

type SnsableWithContext[I Snsable] struct {
	object  I
	context context.Context
}

func WrapSnsableWithContext[I Snsable](object I, context context.Context) *SnsableWithContext[I] {
	return &SnsableWithContext[I]{
		object:  object,
		context: context,
	}
}

func (s *SnsableWithContext[I]) Context() context.Context {
	return s.context
}

func (s *SnsableWithContext[I]) MarshalJSON() ([]byte, error) {
	return s.object.MarshalJSON()
}

func (s *SnsableWithContext[I]) SNSAttributes() map[string][]string {
	return s.object.SNSAttributes()
}

func (s *SnsableWithContext[I]) AddSNSAttribute(key, value string) {
	s.object.AddSNSAttribute(key, value)
}
