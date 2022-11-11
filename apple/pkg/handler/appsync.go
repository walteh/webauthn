package handler

type AppSyncHandler struct{}

// func (me AppSyncHandler) ParseRequest(handler LambdaHander, _event map[string]interface{}) (Request, error) {

// 	encoded, err := json.Marshal(_event)
// 	if err != nil {
// 		log.Println("error encoding event: ", err.Error())
// 		return Request{}, err
// 	}

// 	var event events.AppSyncLambdaAuthorizerRequest
// 	err = json.Unmarshal(encoded, &event)
// 	if err != nil {
// 		log.Println("error unmarshalling event: ", err.Error())
// 		return Request{}, err
// 	}

// 	// unbase 64 the token

// 	fixed := strings.HasPrefix(event.AuthorizationToken, constants.AppsyncAuthHeaderPrefix)
// 	if !fixed {
// 		return Request{}, fmt.Errorf("invalid authorization header")
// 	}

// 	str, err := base64.StdEncoding.DecodeString((strings.Replace(event.AuthorizationToken, constants.AppsyncAuthHeaderPrefix, "", 1)))
// 	if err != nil {
// 		log.Println("error decoding token: ", err.Error())
// 		return Request{}, fmt.Errorf("invalid authorization header")
// 	}

// 	var appsyncAuthToken map[string]string
// 	err = json.Unmarshal(str, &appsyncAuthToken)
// 	if err != nil {
// 		log.Println("error unmarshalling token: ", err.Error())
// 		return Request{}, err
// 	}

// 	return Request{
// 		AppleJwtToken: appsyncAuthToken[constants.AppleJwtHeader],
// 		Operation:     event.RequestContext.APIID,
// 	}, nil
// }

// func (me AppSyncHandler) FormatResponse(handler LambdaHander, isAuthorized bool, context map[string]interface{}, _err error) (interface{}, error) {
// 	if _err != nil {
// 		log.Println("error: ", _err.Error())
// 		return events.AppSyncLambdaAuthorizerResponse{
// 			IsAuthorized: false,
// 			ResolverContext: map[string]interface{}{
// 				"message": _err.Error(),
// 			},
// 		}, nil
// 	}

// 	if context == nil {
// 		context = map[string]interface{}{}
// 	}

// 	return events.AppSyncLambdaAuthorizerResponse{
// 		IsAuthorized:    isAuthorized,
// 		ResolverContext: context,
// 	}, nil

// }
