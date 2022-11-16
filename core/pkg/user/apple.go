package user

import (
	"nugg-auth/core/pkg/signinwithapple"
	"nugg-auth/core/pkg/webauthn/webauthn"
	"time"
)

type Credentials struct {
	list []webauthn.Credential
}

type AppleAuthData struct {
	Id                       string                  `dynamodbav:"id" json:"id"`
	SignInWithAppleSession   *SignInWithAppleSession `dynamodbav:"signinwithapple_session,omitempty"    json:"signinwithapple_session,omitempty"`
	AppleWebAuthnCredentials *Credentials            `dynamodbav:"apple_webauthn_credentials" json:"apple_webauthn_credentials"`
	CreatedAt                int64                   `dynamodbav:"created_at"          json:"created_at"`
	UpdatedAt                int64                   `dynamodbav:"updated_at"          json:"updated_at"`
}

type SignInWithAppleSession struct {
	AccessToken  string `dynamodbav:"access_token"  json:"access_token"`
	Ttl          int64  `dynamodbav:"ttl"           json:"ttl"`
	CreatedAt    int64  `dynamodbav:"created_at"    json:"created_at"`
	UpdatedAt    int64  `dynamodbav:"updated_at"    json:"updated_at"`
	RefreshedAt  int64  `dynamodbav:"refreshed_at"  json:"refreshed_at"`
	CognitoId    string `dynamodbav:"cognito_id"    json:"cognito_id"`
	RefreshToken string `dynamodbav:"refresh_token" json:"refresh_token"`
}

func NewEmptyAppleAuthData(id string) *AppleAuthData {
	now := time.Now().Unix()

	return &AppleAuthData{
		Id:                       id,
		AppleWebAuthnCredentials: &Credentials{},
		CreatedAt:                now,
		UpdatedAt:                now,
	}
}

func NewAppleAuthData(id string, cognitoId string, abc *signinwithapple.ValidationResponse) *AppleAuthData {
	now := time.Now().Unix()

	return &AppleAuthData{
		Id: id,
		SignInWithAppleSession: &SignInWithAppleSession{
			AccessToken:  abc.AccessToken,
			RefreshToken: abc.RefreshToken,
			Ttl:          (now) + int64(abc.ExpiresIn),
			CreatedAt:    now,
			UpdatedAt:    now,
			RefreshedAt:  now,
			CognitoId:    cognitoId,
		},
		AppleWebAuthnCredentials: &Credentials{},
		CreatedAt:                now,
		UpdatedAt:                now,
	}
}

func (user *AppleAuthData) AddAppleWebAuthnCredentials(cred *webauthn.Credential) {
	user.AppleWebAuthnCredentials.list = append(user.AppleWebAuthnCredentials.list, *cred)
}

func (user *AppleAuthData) SetSignInWithAppleSession(session *signinwithapple.ValidationResponse, cognitoId string) {
	now := time.Now().Unix()

	user.SignInWithAppleSession = &SignInWithAppleSession{
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		Ttl:          (now) + int64(session.ExpiresIn),
		CreatedAt:    now,
		UpdatedAt:    now,
		RefreshedAt:  now,
		CognitoId:    cognitoId,
	}
}

func (user *User) IsAppleUser() bool {
	return user.AppleId != ""
}

func (user *User) IsAppleUserWithWebAuthn() bool {
	return user.IsAppleUser() && user.AppleAuthData.HasAppleWebAuthnCredentials()
}

func (user *AppleAuthData) HasAppleWebAuthnCredentials() bool {
	return len(user.AppleWebAuthnCredentials.list) > 0
}

type AppleWebAuthnUser struct {
	Username    string
	AppleId     string
	Credentials []webauthn.Credential
}

func (user *User) CreateAppleWebAuthnUser() *AppleWebAuthnUser {

	return &AppleWebAuthnUser{
		Username:    user.Username,
		AppleId:     user.AppleId,
		Credentials: (*user.AppleAuthData.AppleWebAuthnCredentials).list,
	}
}

func (user *AppleWebAuthnUser) WebAuthnID() []byte {
	return []byte(user.AppleId)
}

func (user *AppleWebAuthnUser) WebAuthnName() string {
	return user.Username
}

func (user *AppleWebAuthnUser) WebAuthnDisplayName() string {
	return user.Username
}

func (user *AppleWebAuthnUser) WebAuthnIcon() string {
	return "https://pics.com/avatar.png"
}

func (user *AppleWebAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	return user.Credentials
}
