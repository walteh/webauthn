package cdkpipelinesgithub

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/cdklabs/cdk-pipelines-github-go/cdkpipelinesgithub/jsii"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/cdklabs/cdk-pipelines-github-go/cdkpipelinesgithub/internal"
)

// Creates or references a GitHub OIDC provider and accompanying role that trusts the provider.
//
// This role can be used to authenticate against AWS instead of using long-lived AWS user credentials
// stored in GitHub secrets.
//
// You can do this manually in the console, or create a separate stack that uses this construct.
// You must `cdk deploy` once (with your normal AWS credentials) to have this role created for you.
//
// You can then make note of the role arn in the stack output and send it into the Github Workflow app via
// the `gitHubActionRoleArn` property. The role arn will be `arn:aws:iam::<accountId>:role/GithubActionRole`.
// See: https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/configuring-openid-connect-in-amazon-web-services
//
// Experimental.
type GitHubActionRole interface {
	constructs.Construct
	// The tree node.
	// Experimental.
	Node() constructs.Node
	// The role that gets created.
	//
	// You should use the arn of this role as input to the `gitHubActionRoleArn`
	// property in your GitHub Workflow app.
	// Experimental.
	Role() awsiam.IRole
	// Returns a string representation of this construct.
	// Experimental.
	ToString() *string
}

// The jsii proxy struct for GitHubActionRole
type jsiiProxy_GitHubActionRole struct {
	internal.Type__constructsConstruct
}

func (j *jsiiProxy_GitHubActionRole) Node() constructs.Node {
	var returns constructs.Node
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubActionRole) Role() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"role",
		&returns,
	)
	return returns
}


// Experimental.
func NewGitHubActionRole(scope constructs.Construct, id *string, props *GitHubActionRoleProps) GitHubActionRole {
	_init_.Initialize()

	if err := validateNewGitHubActionRoleParameters(scope, id, props); err != nil {
		panic(err)
	}
	j := jsiiProxy_GitHubActionRole{}

	_jsii_.Create(
		"cdk-pipelines-github.GitHubActionRole",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewGitHubActionRole_Override(g GitHubActionRole, scope constructs.Construct, id *string, props *GitHubActionRoleProps) {
	_init_.Initialize()

	_jsii_.Create(
		"cdk-pipelines-github.GitHubActionRole",
		[]interface{}{scope, id, props},
		g,
	)
}

// Reference an existing GitHub Actions provider.
//
// You do not need to pass in an arn because the arn for such
// a provider is always the same.
// Experimental.
func GitHubActionRole_ExistingGitHubActionsProvider(scope constructs.Construct) awsiam.IOpenIdConnectProvider {
	_init_.Initialize()

	if err := validateGitHubActionRole_ExistingGitHubActionsProviderParameters(scope); err != nil {
		panic(err)
	}
	var returns awsiam.IOpenIdConnectProvider

	_jsii_.StaticInvoke(
		"cdk-pipelines-github.GitHubActionRole",
		"existingGitHubActionsProvider",
		[]interface{}{scope},
		&returns,
	)

	return returns
}

// Checks if `x` is a construct.
//
// Returns: true if `x` is an object created from a class which extends `Construct`.
// Deprecated: use `x instanceof Construct` instead.
func GitHubActionRole_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	if err := validateGitHubActionRole_IsConstructParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"cdk-pipelines-github.GitHubActionRole",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func (g *jsiiProxy_GitHubActionRole) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		g,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

