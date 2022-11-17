package user

// type PassKeys struct {
// 	list []webauthn.Credential
// }

// type AppleAuthData struct {
// 	Id        string    `dynamodbav:"id"                  json:"id"`
// 	PassKeys  *PassKeys `dynamodbav:"passkeys"            json:"passkeys"`
// 	CreatedAt int64     `dynamodbav:"created_at"          json:"created_at"`
// 	UpdatedAt int64     `dynamodbav:"updated_at"          json:"updated_at"`
// }

// // SignInWithAppleSession   *SignInWithAppleSession `dynamodbav:"signinwithapple_session,omitempty"    json:"signinwithapple_session,omitempty"`

// // type SignInWithAppleSession struct {
// // 	AccessToken  string `dynamodbav:"access_token"  json:"access_token"`
// // 	Ttl          int64  `dynamodbav:"ttl"           json:"ttl"`
// // 	CreatedAt    int64  `dynamodbav:"created_at"    json:"created_at"`
// // 	UpdatedAt    int64  `dynamodbav:"updated_at"    json:"updated_at"`
// // 	RefreshedAt  int64  `dynamodbav:"refreshed_at"  json:"refreshed_at"`
// // 	CognitoId    string `dynamodbav:"cognito_id"    json:"cognito_id"`
// // 	RefreshToken string `dynamodbav:"refresh_token" json:"refresh_token"`
// // }

// func NewAppleAuthData(id string, cognitoId string, abc *signinwithapple.ValidationResponse) *AppleAuthData {
// 	now := time.Now().Unix()

// 	return &AppleAuthData{
// 		Id:        id,
// 		PassKeys:  &PassKeys{},
// 		CreatedAt: now,
// 		UpdatedAt: now,
// 	}
// }

// func (user *AppleAuthData) AddPassKey(cred *webauthn.Credential) {
// 	user.PassKeys.list = append(user.PassKeys.list, *cred)
// }

// // func (user *AppleAuthData) SetSignInWithAppleSession(session *signinwithapple.ValidationResponse, cognitoId string) {
// // 	now := time.Now().Unix()

// // 	user.SignInWithAppleSession = &SignInWithAppleSession{
// // 		AccessToken:  session.AccessToken,
// // 		RefreshToken: session.RefreshToken,
// // 		Ttl:          (now) + int64(session.ExpiresIn),
// // 		CreatedAt:    now,
// // 		UpdatedAt:    now,
// // 		RefreshedAt:  now,
// // 		CognitoId:    cognitoId,
// // 	}
// // }

// func (user *User) IsAppleUser() bool {
// 	return user.AppleId != ""
// }

// func (user *User) IsAppleUserWithWebAuthn() bool {
// 	return user.IsAppleUser() && user.AppleAuthData.HasPassKey()
// }

// func (user *AppleAuthData) HasPassKey() bool {
// 	return len(user.PassKeys.list) > 0
// }

// type AppleWebAuthnUser struct {
// 	Username    string
// 	AppleId     string
// 	Credentials []webauthn.Credential
// }

// func (user *User) CreateAppleWebAuthnUser() *AppleWebAuthnUser {
// 	return &AppleWebAuthnUser{
// 		Username:    user.Username,
// 		AppleId:     user.AppleId,
// 		Credentials: (*user.AppleAuthData.PassKeys).list,
// 	}
// }

// func (user *AppleWebAuthnUser) WebAuthnID() []byte {
// 	return []byte(user.AppleId)
// }

// func (user *AppleWebAuthnUser) WebAuthnName() string {
// 	return user.Username
// }

// func (user *AppleWebAuthnUser) WebAuthnDisplayName() string {
// 	return user.Username
// }

// func (user *AppleWebAuthnUser) WebAuthnIcon() string {
// 	return "https://pics.com/avatar.png"
// }

// func (user *AppleWebAuthnUser) WebAuthnCredentials() []webauthn.Credential {
// 	return user.Credentials
// }
