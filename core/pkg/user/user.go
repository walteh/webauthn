package user

import (
	"nugg-auth/core/pkg/signinwithapple"
	"time"
)

type User struct {
	Id        string `dynamodbav:"id"                  json:"id"`
	Username  string `dynamodbav:"username"            json:"username"`
	CreatedAt int64  `dynamodbav:"created_at"          json:"created_at"`
	UpdatedAt int64  `dynamodbav:"updated_at"          json:"updated_at"`

	// Apple
	AppleId       string         `dynamodbav:"apple_id,omitempty"                   json:"apple_id,omitempty"`
	AppleAuthData *AppleAuthData `dynamodbav:"apple_auth_data,omitempty"            json:"apple_auth_data,omitempty"`
}

func NewAppleUser(newId string, username string, appleId string, cognitoId string, abc *signinwithapple.ValidationResponse) *User {
	now := time.Now().Unix()

	return &User{
		Id:            newId,
		Username:      username,
		AppleId:       appleId,
		AppleAuthData: NewAppleAuthData(appleId, cognitoId, abc),
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}
