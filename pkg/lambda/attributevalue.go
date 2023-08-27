// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package lambda

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func unmarshalNull() (types.AttributeValue, error) {
	return &types.AttributeValueMemberNULL{Value: true}, nil
}

func unmarshalString(value interface{}) (types.AttributeValue, error) {
	var ok bool
	v, ok := value.(string)
	if !ok {
		return nil, errors.New("DynamoDBAttributeValue: S type should contain a string")
	}
	return &types.AttributeValueMemberS{Value: v}, nil

}

func unmarshalBinary(value interface{}) (types.AttributeValue, error) {
	stringValue, ok := value.(string)
	if !ok {
		return nil, errors.New("DynamoDBAttributeValue: B type should contain a base64 string")
	}

	binaryValue, err := base64.StdEncoding.DecodeString(stringValue)
	if err != nil {
		return nil, err
	}

	return &types.AttributeValueMemberB{Value: binaryValue}, nil

}

func unmarshalBoolean(value interface{}) (types.AttributeValue, error) {
	booleanValue, ok := value.(bool)
	if !ok {
		return nil, errors.New("DynamoDBAttributeValue: BOOL type should contain a boolean")
	}

	return &types.AttributeValueMemberBOOL{Value: booleanValue}, nil

}

func unmarshalBinarySet(value interface{}) (types.AttributeValue, error) {
	list, ok := value.([]interface{})
	if !ok {
		return nil, errors.New("DynamoDBAttributeValue: BS type should contain a list of base64 strings")
	}

	binarySet := make([][]byte, len(list))

	for index, element := range list {
		var err error
		elementString := element.(string)
		binarySet[index], err = base64.StdEncoding.DecodeString(elementString)
		if err != nil {
			return nil, err
		}
	}

	return &types.AttributeValueMemberBS{Value: binarySet}, nil

}

func unmarshalList(value interface{}) (types.AttributeValue, error) {
	list, ok := value.([]interface{})
	if !ok {
		return nil, errors.New("DynamoDBAttributeValue: L type should contain a list")
	}

	r := &types.AttributeValueMemberL{Value: make([]types.AttributeValue, len(list))}
	for index, element := range list {

		elementMap, ok := element.(map[string]interface{})
		if !ok {
			return nil, errors.New("DynamoDBAttributeValue: element of a list is not an DynamoDBAttributeValue")
		}

		elementDynamoDBAttributeValue, err := unmarshalDynamoDBAttributeValueMap(elementMap)
		if err != nil {
			return nil, errors.New("DynamoDBAttributeValue: unmarshal of child DynamoDBAttributeValue failed")
		}
		r.Value[index] = elementDynamoDBAttributeValue
	}
	return r, nil

}

func unmarshalMap(value interface{}) (types.AttributeValue, error) {
	m, ok := value.(map[string]interface{})
	if !ok {
		return nil, errors.New("DynamoDBAttributeValue: M type should contain a map")
	}

	r := types.AttributeValueMemberM{Value: make(map[string]types.AttributeValue)}

	for k, v := range m {

		elementMap, ok := v.(map[string]interface{})
		if !ok {
			return nil, errors.New("DynamoDBAttributeValue: element of a map is not an DynamoDBAttributeValue")
		}

		elementDynamoDBAttributeValue, err := unmarshalDynamoDBAttributeValueMap(elementMap)
		if err != nil {
			return nil, errors.New("DynamoDBAttributeValue: unmarshal of child DynamoDBAttributeValue failed")
		}
		r.Value[k] = elementDynamoDBAttributeValue
	}

	return &r, nil

}

func unmarshalNumber(value interface{}) (types.AttributeValue, error) {
	var ok bool
	v, ok := value.(string)
	if !ok {
		return nil, errors.New("DynamoDBAttributeValue: N type should contain a string")
	}
	return &types.AttributeValueMemberN{Value: v}, nil

}

func unmarshalNumberSet(value interface{}) (types.AttributeValue, error) {
	list, ok := value.([]interface{})
	if !ok {
		return nil, errors.New("DynamoDBAttributeValue: NS type should contain a list of strings")
	}

	numberSet := make([]string, len(list))

	for index, element := range list {
		numberSet[index], ok = element.(string)
		if !ok {
			return nil, errors.New("DynamoDBAttributeValue: NS type should contain a list of strings")
		}
	}

	return &types.AttributeValueMemberNS{Value: numberSet}, nil

}

func unmarshalStringSet(value interface{}) (types.AttributeValue, error) {
	list, ok := value.([]interface{})
	if !ok {
		return nil, errors.New("DynamoDBAttributeValue: SS type should contain a list of strings")
	}

	stringSet := make([]string, len(list))

	for index, element := range list {
		stringSet[index], ok = element.(string)
		if !ok {
			return nil, errors.New("DynamoDBAttributeValue: SS type should contain a list of strings")
		}
	}

	return &types.AttributeValueMemberSS{Value: stringSet}, nil

}

func unmarshalDynamoDBAttributeValue(typeLabel string, jsonValue interface{}) (types.AttributeValue, error) {

	switch typeLabel {
	case "NULL":
		return unmarshalNull()
	case "B":
		return unmarshalBinary(jsonValue)
	case "BOOL":
		return unmarshalBoolean(jsonValue)
	case "BS":
		return unmarshalBinarySet(jsonValue)
	case "L":
		return unmarshalList(jsonValue)
	case "M":
		return unmarshalMap(jsonValue)
	case "N":
		return unmarshalNumber(jsonValue)
	case "NS":
		return unmarshalNumberSet(jsonValue)
	case "S":
		return unmarshalString(jsonValue)
	case "SS":
		return unmarshalStringSet(jsonValue)
	default:
		return nil, &attributevalue.InvalidUnmarshalError{Type: reflect.TypeOf(jsonValue)}
	}
}

// // UnmarshalJSON unmarshals a JSON description of this DynamoDBAttributeValue
// func UnmarshalJSON(b []byte, mapper types.AttributeValue) error {
// 	var m interface{}

// 	err := json.Unmarshal(b, &m)
// 	if err != nil {
// 		return err
// 	}
// 	return unmarshalDynamoDBAttributeValue(mapper, "M", m)
// }

func (me DynamoDBImage) MarshalJSON() ([]byte, error) {

	mapper := make(map[string]map[string]interface{})
	for k, v := range me {
		vb, err := marshalDynamoDBAttributeValue(v)
		if err != nil {
			return nil, err
		}
		mapper[k] = vb
	}

	return json.Marshal(mapper)
}

func unmarshalDynamoDBAttributeValueMap(m map[string]interface{}) (types.AttributeValue, error) {
	if m == nil {
		return nil, errors.New("DynamoDBAttributeValue: does not contain a map")
	}

	if len(m) != 1 {
		return nil, errors.New("DynamoDBAttributeValue: map must contain a single type")
	}

	for k, v := range m {
		return unmarshalDynamoDBAttributeValue(k, v)
	}

	return nil, nil
}

func marshalDynamoDBAttributeValue(av types.AttributeValue) (map[string]interface{}, error) {

	switch x := av.(type) {
	case *types.AttributeValueMemberB:
		return map[string]interface{}{"B": x.Value}, nil
	case *types.AttributeValueMemberBOOL:
		return map[string]interface{}{"BOOL": x.Value}, nil
	case *types.AttributeValueMemberBS:
		return map[string]interface{}{"BS": x.Value}, nil
	case *types.AttributeValueMemberL:
		return marshalList(x)
	case *types.AttributeValueMemberM:
		return marshalMap(x)
	case *types.AttributeValueMemberN:
		return map[string]interface{}{"N": x.Value}, nil
	case *types.AttributeValueMemberNS:
		return map[string]interface{}{"NS": x.Value}, nil
	case *types.AttributeValueMemberNULL:
		return map[string]interface{}{"NULL": x.Value}, nil
	case *types.AttributeValueMemberS:
		return map[string]interface{}{"S": x.Value}, nil
	case *types.AttributeValueMemberSS:
		return map[string]interface{}{"SS": x.Value}, nil
	default:
		return nil, &attributevalue.InvalidUnmarshalError{Type: reflect.TypeOf(av)}
	}
}

func marshalList(l types.AttributeValue) (map[string]interface{}, error) {

	if l == nil {
		return nil, errors.New("DynamoDBAttributeValue: does not contain a list")
	}

	if wrk := l.(*types.AttributeValueMemberL); wrk == nil {
		return nil, errors.New("DynamoDBAttributeValue: list must contain a single type")
	} else {
		r := make([]interface{}, len(wrk.Value))
		for i, v := range wrk.Value {
			r[i] = v
		}
		return map[string]interface{}{"L": r}, nil
	}
}

func marshalMap(m types.AttributeValue) (map[string]interface{}, error) {

	if m == nil {
		return nil, errors.New("DynamoDBAttributeValue: does not contain a map")
	}

	if wrk := m.(*types.AttributeValueMemberM); wrk == nil {
		return nil, errors.New("DynamoDBAttributeValue: map must contain a single type")
	} else {
		r := make(map[string]interface{}, len(wrk.Value))
		for k, v := range wrk.Value {
			r[k] = v
		}
		return map[string]interface{}{"M": r}, nil
	}
}

func marshalDynamoDBAttributeValueMap(m types.AttributeValue) (map[string]interface{}, error) {
	if m == nil {
		return nil, errors.New("DynamoDBAttributeValue: does not contain a map")
	}

	if wrk := m.(*types.AttributeValueMemberM); wrk == nil {
		return nil, errors.New("DynamoDBAttributeValue: map must contain a single type")
	} else {

		return marshalDynamoDBAttributeValue(wrk)
	}

}
