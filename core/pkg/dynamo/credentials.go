package dynamo

type CredentialType string

const (
	WebAuthnCredentialType     = CredentialType("webauthn")
	ApplePassKeyCredentialType = CredentialType("apple-passkey")
)

// type DynamoCredential struct {
// 	Id               string `dynamodbav:"credential_id"       json:"credential_id"`
// 	NuggId           string `dynamodbav:"nugg_id"             json:"nugg_id"`
// 	CredentialUserId []byte `dynamodbav:"credential_user_id"  json:"credential_user_id"`
// 	Type             string `dynamodbav:"type"     json:"type"`
// 	CreatedAt        int64  `dynamodbav:"created_at"          json:"created_at"`
// 	UpdatedAt        int64  `dynamodbav:"updated_at"          json:"updated_at"`
// 	Data             []byte `dynamodbav:"dat"     json:"dat"`
// }

// func (client *Client) decodeApplePassKey(data *DynamoCredential) (userId string, credential *webauthn.Credential, err error) {
// 	if data.Type != string(ApplePassKeyCredentialType) {
// 		return "", nil, ErrInvalidCredentialType
// 	}
// 	err = json.Unmarshal(data.Data, &credential)
// 	if err != nil {
// 		return "", nil, err
// 	}

// 	return data.NuggId, credential, nil
// }

// func (c *Client) makeDynamoCredentialPut(d attributevalue.Marshaler) (*types.Put, error) {
// 	av, err := d.MarshalDynamoDBAttributeValue()
// 	if err != nil {
// 		return nil, err
// 	}

// 	r, ok := av.(map[string]types.AttributeValue)
// 	if !ok {
// 		return nil, fmt.Errorf("failed to convert to map[string]types.AttributeValue")
// 	}

// 	return &types.Put{Item: r, TableName: c.MustCredentialTableName()}, nil
// }

// func (c *Client) makeDynamoCredentialUpdate(d *DynamoCredential) *types.Update {
// 	return &types.Update{
// 		Key: map[string]types.AttributeValue{
// 			"credential_id": &types.AttributeValueMemberS{Value: d.Id},
// 		},
// 		TableName: c.MustCredentialTableName(),
// 		ExpressionAttributeValues: map[string]types.AttributeValue{
// 			":dat":        &types.AttributeValueMemberN{Value: d.S},
// 			":updated_at": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", d.UpdatedAt)},
// 		},
// 		UpdateExpression: aws.String("SET dat = :dat, updated_at = :updated_at"),
// 		// UpdateExpression: aws.String("SET data = :dat, updated_at = :updated_at"),
// 	}
// }

// func (client *Client) newCredentialFromApplePassKeyData(userId string, credentialUserId []byte, credential *webauthn.Credential) (*DynamoCredential, error) {
// 	now := time.Now()

// 	raw, err := json.Marshal(credential)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &DynamoCredential{
// 		Id:               base64.RawURLEncoding.EncodeToString(credential.ID),
// 		NuggId:           userId,
// 		CredentialUserId: credentialUserId,
// 		Type:             string(ApplePassKeyCredentialType),
// 		CreatedAt:        now.Unix(),
// 		UpdatedAt:        now.Unix(),
// 		Data:             raw,
// 	}, nil
// }

// func (client *Client) NewApplePassKeyCredentialUpdate(userId string, credentialUserId []byte, credential *webauthn.Credential) (*types.Update, error) {
// 	d, err := client.newCredentialFromApplePassKeyData(userId, credentialUserId, credential)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return client.makeDynamoCredentialUpdate(d), nil
// }

// func (client *Client) NewApplePassKeyCredentialPut(userId string, credentialUserId []byte, credential *webauthn.Credential) (*types.Put, error) {
// 	d, err := client.newCredentialFromApplePassKeyData(userId, credentialUserId, credential)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return client.makeDynamoCredentialPut(d)
// }

// func (client *Client) NewCredentialGet(challenge string) *types.Get {
// 	return &types.Get{
// 		Key: map[string]types.AttributeValue{
// 			"credential_id": &types.AttributeValueMemberS{Value: protocol.ResolveToRawURLEncoding(challenge)},
// 		},
// 		TableName: client.MustCredentialTableName(),
// 	}
// }

// func (client *Client) ParseApplePassKeyCredential(data map[string]types.AttributeValue) (userId string, credential *webauthn.Credential, err error) {
// 	cred, err := client.decodeCredentialFromDynamo(data)
// 	if err != nil {
// 		return "", nil, err
// 	}
// 	return client.decodeApplePassKey(cred)
// }

// func (client *Client) FindApplePassKeyInGetResult(result []*GetOutput) (userId string, credential *webauthn.Credential, err error) {
// 	cred, err := FindInOnePerTableGetResult[DynamoCredential](result, client.MustCredentialTableName())
// 	if err != nil {
// 		return "", nil, err
// 	}
// 	if cred == nil {
// 		return "", nil, ErrNotFound
// 	}

// 	return client.decodeApplePassKey(cred)
// }
