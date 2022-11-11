package handler

// type ApiGatewayV2AuthorizerService struct{}

// func (me ApiGatewayV2AuthorizerService) ParseRequest(handler LambdaHander, _event map[string]interface{}) (Request, error) {
// 	encoded, err := json.Marshal(_event)
// 	if err != nil {
// 		log.Println("error encoding event: ", err.Error())
// 		return Request{}, err
// 	}

// 	var event events.APIGatewayV2CustomAuthorizerV2Request
// 	err = json.Unmarshal(encoded, &event)
// 	if err != nil {
// 		log.Println("error unmarshalling event: ", err.Error())
// 		return Request{}, err
// 	}
// 	return Request{
// 		AppleJwtToken: event.Headers[constants.AppleJwtHeader],
// 	}, nil
// }

// func (me ApiGatewayV2AuthorizerService) FormatResponse(handler LambdaHander, isAuthorized bool, context map[string]interface{}, _err error) (interface{}, error) {
// 	if _err != nil {
// 		log.Println("error: ", _err.Error())
// 		return &events.APIGatewayV2CustomAuthorizerSimpleResponse{
// 			IsAuthorized: false,
// 			Context: map[string]interface{}{
// 				"message": _err.Error(),
// 			},
// 		}, nil
// 	}

// 	if context == nil {
// 		context = map[string]interface{}{}
// 	}

// 	return &events.APIGatewayV2CustomAuthorizerSimpleResponse{
// 		IsAuthorized: isAuthorized,
// 		Context:      context,
// 	}, nil
// }
