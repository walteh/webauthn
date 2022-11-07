package handler

import (
	"encoding/json"
	"nugg-crypto/go/pkg/constants"

	"github.com/aws/aws-lambda-go/events"
)

type ApigwV2Handler struct{}

func (me ApigwV2Handler) ParseRequest(handler LambdaHander, _event map[string]interface{}) (Request, error) {
	encoded, err := json.Marshal(_event)
	if err != nil {
		return Request{}, err
	}

	var event events.APIGatewayV2CustomAuthorizerV2Request
	err = json.Unmarshal(encoded, &event)
	if err != nil {
		return Request{}, err
	}
	return Request{
		AppleJwtToken: event.Headers[constants.AppleJwtHeader],
	}, nil
}

func (me ApigwV2Handler) FormatResponse(handler LambdaHander, isAuthorized bool, context map[string]interface{}, _err error) (interface{}, error) {
	if _err != nil {
		return &events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
			Context: map[string]interface{}{
				"message": _err.Error(),
			},
		}, nil
	}

	if context == nil {
		context = map[string]interface{}{}
	}

	return &events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: isAuthorized,
		Context:      context,
	}, nil
}
