package handler

type ApiGatewayV2Service struct{}

// func (me ApiGatewayV2Service) ParseRequest(handler LambdaHander, _event map[string]interface{}) (Request, error) {
// 	encoded, err := json.Marshal(_event)
// 	if err != nil {
// 		log.Println("error encoding event: ", err.Error())
// 		return Request{}, err
// 	}

// 	var event events.APIGatewayV2HTTPRequest
// 	err = json.Unmarshal(encoded, &event)
// 	if err != nil {
// 		log.Println("error unmarshalling event: ", err.Error())
// 		return Request{}, err
// 	}
// 	return Request{
// 		AppleJwtToken: event.Headers[constants.AppleJwtHeader],
// 	}, nil
// }

// func (me ApiGatewayV2Service) FormatResponse(handler LambdaHander, isAuthorized bool, context map[string]interface{}, _err error) (interface{}, error) {
// 	if _err != nil {
// 		log.Println("error: ", _err.Error())
// 		return &events.APIGatewayV2HTTPResponse{}, nil
// 	}

// 	if context == nil {
// 		context = map[string]interface{}{}
// 	}

// 	return &events.APIGatewayV2CustomAuthorizerSimpleResponse{
// 		IsAuthorized: isAuthorized,
// 		Context:      context,
// 	}, nil
// }
