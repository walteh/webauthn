package signinwithapple

import (
	"context"
	"fmt"
)

func ValidateAppTokenAndObtainID() {
	// Your 10-character Team ID
	teamID := "XXXXXXXXXX"

	// ClientID is the "Services ID" value that you get when navigating to your "sign in with Apple"-enabled service ID
	clientID := "com.your.app"

	// Find the 10-char Key ID value from the portal
	keyID := "XXXXXXXXXX"

	// The contents of the p8 file/key you downloaded when you made the key in the portal
	secret := `-----BEGIN PRIVATE KEY-----
YOUR_SECRET_PRIVATE_KEY
-----END PRIVATE KEY-----`

	// Generate the client secret used to authenticate with Apple's validation servers
	secret, err := GenerateClientSecret(secret, teamID, clientID, keyID)
	if err != nil {
		fmt.Println("error generating secret: " + err.Error())
		return
	}

	// Generate a new validation client
	client := New()

	vReq := AppValidationTokenRequest{
		ClientID:     clientID,
		ClientSecret: secret,
		Code:         "the_authorization_code_to_validate",
	}

	var resp ValidationResponse

	// Do the verification
	err = client.VerifyAppToken(context.Background(), vReq, &resp)
	if err != nil {
		fmt.Println("error verifying: " + err.Error())
		return
	}

	if resp.Error != "" {
		fmt.Printf("apple returned an error: %s - %s\n", resp.Error, resp.ErrorDescription)
		return
	}

	// Get the unique user ID
	unique, err := GetUniqueID(resp.IDToken)
	if err != nil {
		fmt.Println("failed to get unique ID: " + err.Error())
		return
	}

	// Get the email
	claim, err := GetClaims(resp.IDToken)
	if err != nil {
		fmt.Println("failed to get claims: " + err.Error())
		return
	}

	email := (*claim)["email"]
	emailVerified := (*claim)["email_verified"]
	isPrivateEmail := (*claim)["is_private_email"]

	// Voila!
	fmt.Println(unique)
	fmt.Println(email)
	fmt.Println(emailVerified)
	fmt.Println(isPrivateEmail)
}

func ValidateRefreshToken() {
	// Your 10-character Team ID
	teamID := "XXXXXXXXXX"

	// ClientID is the "Services ID" value that you get when navigating to your "sign in with Apple"-enabled service ID
	clientID := "com.your.app"

	// Find the 10-char Key ID value from the portal
	keyID := "XXXXXXXXXX"

	// The contents of the p8 file/key you downloaded when you made the key in the portal
	secret := `-----BEGIN PRIVATE KEY-----
YOUR_SECRET_PRIVATE_KEY
-----END PRIVATE KEY-----`

	// Generate the client secret used to authenticate with Apple's validation servers
	secret, err := GenerateClientSecret(secret, teamID, clientID, keyID)
	if err != nil {
		fmt.Println("error generating secret: " + err.Error())
		return
	}

	// Generate a new validation client
	client := New()

	vReq := ValidationRefreshRequest{
		ClientID:     clientID,
		ClientSecret: secret,
		RefreshToken: "the_refresh_code_to_validate",
	}

	var resp RefreshResponse

	// Do the verification
	err = client.VerifyRefreshToken(context.Background(), vReq, &resp)
	if err != nil {
		fmt.Println("error verifying: " + err.Error())
		return
	}

	if resp.Error != "" {
		fmt.Printf("apple returned an error: %s - %s\n", resp.Error, resp.ErrorDescription)
		return
	}

	// Voila!
	fmt.Println(resp)
}

/*
This example shows you how to validate a web token for the first time
*/

func ValidateWebTokenAndObtainID() {
	// Your 10-character Team ID
	teamID := "XXXXXXXXXX"

	// ClientID is the "Services ID" value that you get when navigating to your "sign in with Apple"-enabled service ID
	clientID := "com.your.app"

	// Find the 10-char Key ID value from the portal
	keyID := "XXXXXXXXXX"

	// The contents of the p8 file/key you downloaded when you made the key in the portal
	secret := `-----BEGIN PRIVATE KEY-----
YOUR_SECRET_PRIVATE_KEY
-----END PRIVATE KEY-----`

	// Generate the client secret used to authenticate with Apple's validation servers
	secret, err := GenerateClientSecret(secret, teamID, clientID, keyID)
	if err != nil {
		fmt.Println("error generating secret: " + err.Error())
		return
	}

	// Generate a new validation client
	client := New()

	vReq := WebValidationTokenRequest{
		ClientID:     clientID,
		ClientSecret: secret,
		Code:         "the_authorization_code_to_validate",
		RedirectURI:  "https://example.com", // This URL must be validated with apple in your service
	}

	var resp ValidationResponse

	// Do the verification
	err = client.VerifyWebToken(context.Background(), vReq, &resp)
	if err != nil {
		fmt.Println("error verifying: " + err.Error())
		return
	}

	if resp.Error != "" {
		fmt.Printf("apple returned an error: %s - %s\n", resp.Error, resp.ErrorDescription)
		return
	}

	// Get the unique user ID
	unique, err := GetUniqueID(resp.IDToken)
	if err != nil {
		fmt.Println("failed to get unique ID: " + err.Error())
		return
	}

	// Get the email
	claim, err := GetClaims(resp.IDToken)
	if err != nil {
		fmt.Println("failed to get claims: " + err.Error())
		return
	}

	email := (*claim)["email"]
	emailVerified := (*claim)["email_verified"]
	isPrivateEmail := (*claim)["is_private_email"]

	// Voila!
	fmt.Println(unique)
	fmt.Println(email)
	fmt.Println(emailVerified)
	fmt.Println(isPrivateEmail)
}
