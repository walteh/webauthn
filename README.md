# **URL** : `/auth/apple/passkey/init`

**Method** : `POST`

**Request Headers**

```
none
```

### Success Response

**Code** : `204 OK`

**Response Headers**

```http
X-Nugg-Challenge-User: "URL Base64 encoded String"
X-Nugg-Challenge-Raw: "URL Base64 encoded String"
```

### Error Response

**Code** : `500 INTERNAL SERVER ERROR`

<br>
<br>

# **URL** : `/auth/apple/passkey/register`

**Method** : `POST`

**Request Headers**

```http
X-Nugg-Webauthn-Creation: "Standard Base64 encoded JSON defined below"
```

```go
type XNuggWebauthnCreation struct {
	RawAttestationObject []byte `json:"rawAttestationObject"`
	RawClientData        []byte `json:"rawClientData"`
	CredentialID         []byte `json:"credentialID"`
}
```

### Success Response

**Code** : `204 OK`

**Response Headers**

```http
X-Nugg-Challenge-Raw: "URL Base64 encoded String"
```

### Error Response

**Code** : `500 INTERNAL SERVER ERROR`

<br>
<br>

# **URL** : `/auth/apple/passkey/login`

**Method** : `POST`

**Request Headers**

```http
X-Nugg-Webauthn-Assertion: "Standard Base64 encoded JSON defined below"
```

```go
type XNuggWebauthnAssertion struct {
	UserID               []byte `json:"userID"`
	CredentialID         []byte `json:"credentialID"`
	RawClientDataJSON    []byte `json:"rawClientDataJSON"`
	RawAuthenticatorData []byte `json:"rawAuthenticatorData"`
	Signature            []byte `json:"signature"`
	Type                 string `json:"credentialType"`
}
```

### Success Response

**Code** : `204 OK`

**Response Headers**

```http
X-Nugg-Challenge-Raw: "URL Base64 encoded String"
```

### Error Response

**Code** : `500 INTERNAL SERVER ERROR`
