package cdkpipelinesgithub

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/cdklabs/cdk-pipelines-github-go/cdkpipelinesgithub/jsii"
)

// Provides AWS credenitals to the pipeline jobs.
// Experimental.
type AwsCredentials interface {
}

// The jsii proxy struct for AwsCredentials
type jsiiProxy_AwsCredentials struct {
	_ byte // padding
}

// Experimental.
func NewAwsCredentials() AwsCredentials {
	_init_.Initialize()

	j := jsiiProxy_AwsCredentials{}

	_jsii_.Create(
		"cdk-pipelines-github.AwsCredentials",
		nil, // no parameters
		&j,
	)

	return &j
}

// Experimental.
func NewAwsCredentials_Override(a AwsCredentials) {
	_init_.Initialize()

	_jsii_.Create(
		"cdk-pipelines-github.AwsCredentials",
		nil, // no parameters
		a,
	)
}

// Reference credential secrets to authenticate with AWS.
//
// This method assumes
// that your credentials will be stored as long-lived GitHub Secrets.
// Experimental.
func AwsCredentials_FromGitHubSecrets(props *GitHubSecretsProviderProps) AwsCredentialsProvider {
	_init_.Initialize()

	if err := validateAwsCredentials_FromGitHubSecretsParameters(props); err != nil {
		panic(err)
	}
	var returns AwsCredentialsProvider

	_jsii_.StaticInvoke(
		"cdk-pipelines-github.AwsCredentials",
		"fromGitHubSecrets",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Provide AWS credentials using OpenID Connect.
// Experimental.
func AwsCredentials_FromOpenIdConnect(props *OpenIdConnectProviderProps) AwsCredentialsProvider {
	_init_.Initialize()

	if err := validateAwsCredentials_FromOpenIdConnectParameters(props); err != nil {
		panic(err)
	}
	var returns AwsCredentialsProvider

	_jsii_.StaticInvoke(
		"cdk-pipelines-github.AwsCredentials",
		"fromOpenIdConnect",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Don't provide any AWS credentials, use this if runners have preconfigured credentials.
// Experimental.
func AwsCredentials_RunnerHasPreconfiguredCreds() AwsCredentialsProvider {
	_init_.Initialize()

	var returns AwsCredentialsProvider

	_jsii_.StaticInvoke(
		"cdk-pipelines-github.AwsCredentials",
		"runnerHasPreconfiguredCreds",
		nil, // no parameters
		&returns,
	)

	return returns
}

