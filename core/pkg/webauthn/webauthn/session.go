package webauthn

import "nugg-auth/core/pkg/webauthn/protocol"

// SessionData is the data that should be stored by the Relying Party for
// the duration of the web authentication ceremony
type SessionData struct {
	Challenge            string                               `dynamodbav:"id"                            json:"id"`
	UserID               []byte                               `dynamodbav:"user_id"                       json:"user_id"`
	AllowedCredentialIDs [][]byte                             `dynamodbav:"allowed_credentials,omitempty" json:"allowed_credentials,omitempty" `
	UserVerification     protocol.UserVerificationRequirement `dynamodbav:"user_verification"             json:"userVerification"`
	Extensions           protocol.AuthenticationExtensions    `dynamodbav:"extensions,omitempty"          json:"extensions,omitempty"`
}
