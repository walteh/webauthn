package cdkpipelinesgithub

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
)

// Properties for the GitHubActionRole construct.
// Experimental.
type GitHubActionRoleProps struct {
	// The GitHub OpenId Connect Provider. Must have provider url `https://token.actions.githubusercontent.com`. The audience must be `sts:amazonaws.com`.
	//
	// Only one such provider can be defined per account, so if you already
	// have a provider with the same url, a new provider cannot be created for you.
	// Default: - a provider is created for you.
	//
	// Experimental.
	Provider awsiam.IOpenIdConnectProvider `field:"optional" json:"provider" yaml:"provider"`
	// A list of GitHub repositories you want to be able to access the IAM role.
	//
	// Each entry should be your GitHub username and repository passed in as a
	// single string.
	// An entry `owner/repo` is equivalent to the subjectClaim `repo:owner/repo:*`.
	//
	// For example, `['owner/repo1', 'owner/repo2'].
	// Experimental.
	Repos *[]*string `field:"optional" json:"repos" yaml:"repos"`
	// The name of the Oidc role.
	// Default: 'GitHubActionRole'.
	//
	// Experimental.
	RoleName *string `field:"optional" json:"roleName" yaml:"roleName"`
	// A list of subject claims allowed to access the IAM role.
	//
	// See https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/about-security-hardening-with-openid-connect
	// A subject claim can include `*` and `?` wildcards according to the `StringLike`
	// condition operator.
	//
	// For example, `['repo:owner/repo1:ref:refs/heads/branch1', 'repo:owner/repo1:environment:prod']`.
	// Experimental.
	SubjectClaims *[]*string `field:"optional" json:"subjectClaims" yaml:"subjectClaims"`
	// Thumbprints of GitHub's certificates.
	//
	// Every time GitHub rotates their certificates, this value will need to be updated.
	//
	// Default value is up-to-date to June 27, 2023 as per
	// https://github.blog/changelog/2023-06-27-github-actions-update-on-oidc-integration-with-aws/
	// Default: - Use built-in keys.
	//
	// Experimental.
	Thumbprints *[]*string `field:"optional" json:"thumbprints" yaml:"thumbprints"`
}

