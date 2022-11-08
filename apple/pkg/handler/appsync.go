package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"nugg-auth/apple/pkg/constants"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type AppSyncHandler struct{}

func (me AppSyncHandler) ParseRequest(handler LambdaHander, _event map[string]interface{}) (Request, error) {

	encoded, err := json.Marshal(_event)
	if err != nil {
		return Request{}, err
	}

	var event events.AppSyncLambdaAuthorizerRequest
	err = json.Unmarshal(encoded, &event)
	if err != nil {
		return Request{}, err
	}

	// unbase 64 the token

	fixed := strings.HasPrefix(event.AuthorizationToken, constants.AppsyncAuthHeaderPrefix)
	if !fixed {
		return Request{}, fmt.Errorf("invalid authorization header")
	}

	str, err := base64.StdEncoding.DecodeString((strings.Replace(event.AuthorizationToken, constants.AppsyncAuthHeaderPrefix, "", 1)))
	if err != nil {
		return Request{}, err
	}

	var appsyncAuthToken map[string]string
	err = json.Unmarshal(str, &appsyncAuthToken)
	if err != nil {
		return Request{}, err
	}

	return Request{
		AppleJwtToken: appsyncAuthToken[constants.AppleJwtHeader],
		Operation:     event.RequestContext.APIID,
	}, nil
}

func (me AppSyncHandler) FormatResponse(handler LambdaHander, isAuthorized bool, context map[string]interface{}, _err error) (interface{}, error) {
	if _err != nil {
		return &events.AppSyncLambdaAuthorizerResponse{
			IsAuthorized: false,
			ResolverContext: map[string]interface{}{
				"message": _err.Error(),
			},
		}, nil
	}

	if context == nil {
		context = map[string]interface{}{}
	}

	return &events.AppSyncLambdaAuthorizerResponse{
		IsAuthorized:    isAuthorized,
		ResolverContext: context,
	}, nil

}
