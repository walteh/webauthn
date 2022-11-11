package handler

type NoopHandler struct{}

// func (me NoopHandler) ParseRequest(handler LambdaHander, _event map[string]interface{}) (Request, error) {
// 	return Request{}, fmt.Errorf("unable to parse request")
// }

// func (me NoopHandler) FormatResponse(handler LambdaHander, isAuthorized bool, result map[string]interface{}, err error) (interface{}, error) {
// 	return nil, err
// }
