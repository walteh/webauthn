package cdkpipelinesgithub

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/cdklabs/cdk-pipelines-github-go/cdkpipelinesgithub/jsii"
)

// AWS credential provider.
// Experimental.
type AwsCredentialsProvider interface {
	// Experimental.
	CredentialSteps(region *string, assumeRoleArn *string) *[]*JobStep
	// Experimental.
	JobPermission() JobPermission
}

// The jsii proxy struct for AwsCredentialsProvider
type jsiiProxy_AwsCredentialsProvider struct {
	_ byte // padding
}

// Experimental.
func NewAwsCredentialsProvider_Override(a AwsCredentialsProvider) {
	_init_.Initialize()

	_jsii_.Create(
		"cdk-pipelines-github.AwsCredentialsProvider",
		nil, // no parameters
		a,
	)
}

func (a *jsiiProxy_AwsCredentialsProvider) CredentialSteps(region *string, assumeRoleArn *string) *[]*JobStep {
	if err := a.validateCredentialStepsParameters(region); err != nil {
		panic(err)
	}
	var returns *[]*JobStep

	_jsii_.Invoke(
		a,
		"credentialSteps",
		[]interface{}{region, assumeRoleArn},
		&returns,
	)

	return returns
}

func (a *jsiiProxy_AwsCredentialsProvider) JobPermission() JobPermission {
	var returns JobPermission

	_jsii_.Invoke(
		a,
		"jobPermission",
		nil, // no parameters
		&returns,
	)

	return returns
}

