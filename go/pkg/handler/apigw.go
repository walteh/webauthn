package handler

import (
	"nugg-crypto/go/pkg/constants"

	"github.com/aws/aws-lambda-go/events"
)

type ApigwV2Handler struct{}

func (me ApigwV2Handler) ParseRequest(handler LambdaHander, _event interface{}) (Request, error) {
	event := _event.(events.APIGatewayV2CustomAuthorizerV2Request)

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
