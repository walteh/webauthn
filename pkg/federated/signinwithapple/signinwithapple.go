package signinwithapple

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

var (
	ErrInvalidGrant = errors.New("invalid_grant")
)

func IsInvalidGrant(err error) bool {
	return errors.Is(err, ErrInvalidGrant)
}

func (me *Client) ValidateRegistrationCode(ctx context.Context, pk string, registrationCode string) (*ValidationResponse, error) {

	// Generate the client secret used to authenticate with Apple's validation servers
	secret, err := me.config.GenerateClientSecret(pk)
	if err != nil {
		fmt.Println("error generating secret: " + err.Error())
		return nil, err
	}

	vReq := AppValidationTokenRequest{
		ClientID:     me.config.ClientID,
		ClientSecret: secret,
		Code:         registrationCode,
	}

	// Do the verification
	resp, err := me.VerifyAppToken(ctx, vReq)
	if err != nil {
		fmt.Println("error verifying: " + err.Error())
		return nil, err
	}

	if resp.Error != "" {
		return nil, ErrInvalidGrant
	}

	return resp, nil

}

func (me *Client) ValidateRefreshToken(ctx context.Context, pk string, refreshToken string) (*RefreshResponse, error) {

	// Generate the client secret used to authenticate with Apple's validation servers
	secret, err := me.config.GenerateClientSecret(pk)
	if err != nil {
		fmt.Println("error generating secret: " + err.Error())
		return nil, err
	}

	vReq := ValidationRefreshRequest{
		ClientID:     me.config.ClientID,
		ClientSecret: secret,
		RefreshToken: refreshToken,
	}

	// Do the verification
	resp, err := me.VerifyRefreshToken(ctx, vReq)
	if err != nil {
		fmt.Println("error verifying: " + err.Error())
		return nil, err
	}

	if resp.Error != "" {
		fmt.Printf("apple returned an error: %s - %s\n", resp.Error, resp.ErrorDescription)
		return nil, fmt.Errorf("apple returned an error: %s - %s", resp.Error, resp.ErrorDescription)
	}

	return resp, nil
}

/*
This example shows you how to validate a web token for the first time
*/

func (me *Client) ValidateWebToken(ctx context.Context, pk string, authorizationCode string, redirect *url.URL) (*ValidationResponse, error) {

	// Generate the client secret used to authenticate with Apple's validation servers
	secret, err := me.config.GenerateClientSecret(pk)
	if err != nil {
		fmt.Println("error generating secret: " + err.Error())
		return nil, err
	}

	vReq := WebValidationTokenRequest{
		ClientID:     me.config.ClientID,
		ClientSecret: secret,
		Code:         authorizationCode,
		RedirectURI:  redirect.String(), // This URL must be validated with apple in your service
	}

	// Do the verification
	resp, err := me.VerifyWebToken(ctx, vReq)
	if err != nil {
		fmt.Println("error verifying: " + err.Error())
		return nil, err
	}

	if resp.Error != "" {
		fmt.Printf("apple returned an error: %s - %s\n", resp.Error, resp.ErrorDescription)
		return nil, fmt.Errorf("apple returned an error: %s - %s", resp.Error, resp.ErrorDescription)
	}

	return resp, nil
}
