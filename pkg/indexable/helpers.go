package indexable

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ethereum/go-ethereum/common"
)

func S(value string) *types.AttributeValueMemberS {
	return &types.AttributeValueMemberS{Value: value}
}

func NbyS(value string) *types.AttributeValueMemberN {
	return &types.AttributeValueMemberN{Value: value}
}

func N[I int | uint | uint8 | uint16 | uint64 | uint32 | int32 | int8 | int64](value I) *types.AttributeValueMemberN {
	return &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", value)}
}

func F[I float32 | float64](value I) *types.AttributeValueMemberN {
	return &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", value)}
}

func BOOL(value bool) *types.AttributeValueMemberBOOL {
	return &types.AttributeValueMemberBOOL{Value: value}
}

func M(value map[string]types.AttributeValue) *types.AttributeValueMemberM {
	return &types.AttributeValueMemberM{Value: value}
}

func GetS(av *types.AttributeValueMemberM, key string) (string, error) {
	if x, ok := av.Value[key].(*types.AttributeValueMemberS); ok {
		return x.Value, nil
	}
	return "", errors.New("could not unmarshal " + key)
}

func GetNUint64(av *types.AttributeValueMemberM, key string) (uint64, error) {
	if x, ok := av.Value[key].(*types.AttributeValueMemberN); ok {
		return strconv.ParseUint(x.Value, 10, 64)
	}
	return 0, errors.New("could not unmarshal " + key)
}

func GetN[I int | uint | uint8 | uint64 | uint32 | int32 | int8](av *types.AttributeValueMemberM, key string) (I, error) {
	x, err := GetNUint64(av, key)
	return I(x), err
}

func GetBOOL(av *types.AttributeValueMemberM, key string) (bool, error) {
	if x, ok := av.Value[key].(*types.AttributeValueMemberBOOL); ok {
		return x.Value, nil
	}
	return false, errors.New("could not unmarshal " + key)
}

func SS(values []string) *types.AttributeValueMemberSS {
	return &types.AttributeValueMemberSS{Value: values}
}

func NS[I int | uint | uint8 | uint16 | uint64 | uint32 | int32 | int8 | int64](values ...I) *types.AttributeValueMemberNS {
	var ret []string
	for _, v := range values {
		ret = append(ret, fmt.Sprintf("%d", v))
	}
	return &types.AttributeValueMemberNS{Value: ret}
}

func FS[I float32 | float64](values ...I) *types.AttributeValueMemberNS {
	var ret []string
	for _, v := range values {
		ret = append(ret, fmt.Sprintf("%f", v))
	}
	return &types.AttributeValueMemberNS{Value: ret}
}

func L(values ...types.AttributeValue) *types.AttributeValueMemberL {
	return &types.AttributeValueMemberL{Value: values}
}

func GetSS(av *types.AttributeValueMemberM, key string) ([]string, error) {
	if x, ok := av.Value[key].(*types.AttributeValueMemberSS); ok {
		return x.Value, nil
	}
	return nil, errors.New("could not unmarshal " + key)
}

func GetSSHash(av *types.AttributeValueMemberM, key string) ([]common.Hash, error) {
	x, err := GetSS(av, key)
	if err != nil {
		return nil, err
	}
	var ret []common.Hash
	for _, v := range x {
		ret = append(ret, common.HexToHash(v))
	}
	return ret, nil
}

func GetSHash(av *types.AttributeValueMemberM, key string) (common.Hash, error) {
	if x, ok := av.Value[key].(*types.AttributeValueMemberS); ok {
		return common.HexToHash(x.Value), nil
	}
	return common.Hash{}, errors.New("could not unmarshal " + key)
}

func GetSAddress(av *types.AttributeValueMemberM, key string) (common.Address, error) {
	if x, ok := av.Value[key].(*types.AttributeValueMemberS); ok {
		return common.HexToAddress(x.Value), nil
	}
	return common.Address{}, errors.New("could not unmarsal " + key)
}

type SingleMarshaler struct {
	attr types.AttributeValue
}

func (me SingleMarshaler) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return me.attr, nil
}

func AsMarshaler(attr types.AttributeValue) attributevalue.Marshaler {
	return SingleMarshaler{attr}
}

func NULL(isnull bool) *types.AttributeValueMemberNULL {
	return &types.AttributeValueMemberNULL{
		Value: isnull,
	}
}
