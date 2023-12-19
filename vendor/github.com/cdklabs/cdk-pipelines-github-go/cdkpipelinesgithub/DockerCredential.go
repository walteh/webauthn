package cdkpipelinesgithub

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/cdklabs/cdk-pipelines-github-go/cdkpipelinesgithub/jsii"
)

// Represents a credential used to authenticate to a docker registry.
//
// Uses the official Docker Login GitHub Action to authenticate.
// See: https://github.com/marketplace/actions/docker-login
//
// Experimental.
type DockerCredential interface {
	// Experimental.
	Name() *string
	// Experimental.
	PasswordKey() *string
	// Experimental.
	Registry() *string
	// Experimental.
	UsernameKey() *string
}

// The jsii proxy struct for DockerCredential
type jsiiProxy_DockerCredential struct {
	_ byte // padding
}

func (j *jsiiProxy_DockerCredential) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerCredential) PasswordKey() *string {
	var returns *string
	_jsii_.Get(
		j,
		"passwordKey",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerCredential) Registry() *string {
	var returns *string
	_jsii_.Get(
		j,
		"registry",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerCredential) UsernameKey() *string {
	var returns *string
	_jsii_.Get(
		j,
		"usernameKey",
		&returns,
	)
	return returns
}


// Create a credential for a custom registry.
//
// This method assumes that you will have long-lived
// GitHub Secrets stored under the usernameKey and passwordKey that will authenticate to the
// registry you provide.
// See: https://github.com/marketplace/actions/docker-login
//
// Experimental.
func DockerCredential_CustomRegistry(registry *string, creds *ExternalDockerCredentialSecrets) DockerCredential {
	_init_.Initialize()

	if err := validateDockerCredential_CustomRegistryParameters(registry, creds); err != nil {
		panic(err)
	}
	var returns DockerCredential

	_jsii_.StaticInvoke(
		"cdk-pipelines-github.DockerCredential",
		"customRegistry",
		[]interface{}{registry, creds},
		&returns,
	)

	return returns
}

// Reference credential secrets to authenticate to DockerHub.
//
// This method assumes
// that your credentials will be stored as long-lived GitHub Secrets under the
// usernameKey and personalAccessTokenKey.
//
// The default for usernameKey is `DOCKERHUB_USERNAME`. The default for personalAccessTokenKey
// is `DOCKERHUB_TOKEN`. If you do not set these values, your credentials should be
// found in your GitHub Secrets under these default keys.
// Experimental.
func DockerCredential_DockerHub(creds *DockerHubCredentialSecrets) DockerCredential {
	_init_.Initialize()

	if err := validateDockerCredential_DockerHubParameters(creds); err != nil {
		panic(err)
	}
	var returns DockerCredential

	_jsii_.StaticInvoke(
		"cdk-pipelines-github.DockerCredential",
		"dockerHub",
		[]interface{}{creds},
		&returns,
	)

	return returns
}

// Create a credential for ECR.
//
// This method will reuse your AWS credentials to log in to AWS.
// Your AWS credentials are already used to deploy your CDK stacks. It can be supplied via
// GitHub Secrets or using an IAM role that trusts the GitHub OIDC identity provider.
//
// NOTE - All ECR repositories in the same account and region share a domain name
// (e.g., 0123456789012.dkr.ecr.eu-west-1.amazonaws.com), and can only have one associated
// set of credentials (and DockerCredential). Attempting to associate one set of credentials
// with one ECR repo and another with another ECR repo in the same account and region will
// result in failures when using these credentials in the pipeline.
// Experimental.
func DockerCredential_Ecr(registry *string) DockerCredential {
	_init_.Initialize()

	if err := validateDockerCredential_EcrParameters(registry); err != nil {
		panic(err)
	}
	var returns DockerCredential

	_jsii_.StaticInvoke(
		"cdk-pipelines-github.DockerCredential",
		"ecr",
		[]interface{}{registry},
		&returns,
	)

	return returns
}

