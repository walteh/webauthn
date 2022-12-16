package types

import (
	"strconv"

	"github.com/nuggxyz/webauthn/pkg/errors"
	"github.com/nuggxyz/webauthn/pkg/hex"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// type S = types.AttributeValueMemberS

// type N = types.AttributeValueMemberN

// type B = types.AttributeValueMemberB

// type SS = types.AttributeValueMemberSS

// type NS = types.AttributeValueMemberNS

// type BS = types.AttributeValueMemberBS

// type M = types.AttributeValueMemberM

// type L = types.AttributeValueMemberL

// type NULL = types.AttributeValueMemberNULL

// type BOOL = types.AttributeValueMemberBOOL

// make function wrappers for each av type

func S(value string) *types.AttributeValueMemberS {
	return &types.AttributeValueMemberS{Value: value}
}

func N(value string) *types.AttributeValueMemberN {
	return &types.AttributeValueMemberN{Value: value}
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
	return "", errors.NewError(0x11).WithMessage("could not unmarshal " + key).WithCallerStepingBack(1)
}

func GetSHash(av *types.AttributeValueMemberM, key string) (hex.Hash, error) {
	h, err := GetS(av, key)
	if err != nil {
		return hex.Hash{}, err
	}
	return hex.HexToHash(h), nil
}

func GetSHashNotZero(av *types.AttributeValueMemberM, key string) (hex.Hash, error) {
	h, err := GetSHash(av, key)
	if err != nil {
		return hex.Hash{}, err
	}
	if h.IsZero() {
		return hex.Hash{}, errors.NewError(0x11).WithMessage("could not unmarshal " + key).WithCallerStepingBack(1)
	}
	return h, nil
}

func GetNUint64(av *types.AttributeValueMemberM, key string) (uint64, error) {
	if x, ok := av.Value[key].(*types.AttributeValueMemberN); ok {
		return strconv.ParseUint(x.Value, 10, 64)
	}
	return 0, errors.NewError(0x11).WithMessage("could not unmarshal " + key).WithCallerStepingBack(1)
}

func GetBOOL(av *types.AttributeValueMemberM, key string) (bool, error) {
	if x, ok := av.Value[key].(*types.AttributeValueMemberBOOL); ok {
		return x.Value, nil
	}
	return false, errors.NewError(0x11).WithMessage("could not unmarshal " + key).WithCallerStepingBack(1)
}
