package webauthn

// func TestLogin_FinishLoginFailure(t *testing.T) {
// 	user := &defaultUser{
// 		id: []byte("123"),
// 	}
// 	session := SessionData{
// 		UserID: []byte("ABC"),
// 	}

// 	webauthn := &WebAuthn{}
// 	credential, err := webauthn.validateLogin(user, session, nil)
// 	if err == nil {
// 		t.Errorf("FinishLogin() error = nil, want %v", protocol.ErrBadRequest.Type)
// 	}
// 	if credential != nil {
// 		t.Errorf("FinishLogin() credential = %v, want nil", credential)
// 	}
// }
