package cdkpipelinesgithub

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
)

// Experimental.
type GitHubStageProps struct {
	// Default AWS environment (account/region) for `Stack`s in this `Stage`.
	//
	// Stacks defined inside this `Stage` with either `region` or `account` missing
	// from its env will use the corresponding field given here.
	//
	// If either `region` or `account`is is not configured for `Stack` (either on
	// the `Stack` itself or on the containing `Stage`), the Stack will be
	// *environment-agnostic*.
	//
	// Environment-agnostic stacks can be deployed to any environment, may not be
	// able to take advantage of all features of the CDK. For example, they will
	// not be able to use environmental context lookups, will not automatically
	// translate Service Principals to the right format based on the environment's
	// AWS partition, and other such enhancements.
	//
	// Example:
	//   // Use a concrete account and region to deploy this Stage to
	//   new Stage(app, 'Stage1', {
	//     env: { account: '123456789012', region: 'us-east-1' },
	//   });
	//
	//   // Use the CLI's current credentials to determine the target environment
	//   new Stage(app, 'Stage2', {
	//     env: { account: process.env.CDK_DEFAULT_ACCOUNT, region: process.env.CDK_DEFAULT_REGION },
	//   });
	//
	// Default: - The environments should be configured on the `Stack`s.
	//
	// Experimental.
	Env *awscdk.Environment `field:"optional" json:"env" yaml:"env"`
	// The output directory into which to emit synthesized artifacts.
	//
	// Can only be specified if this stage is the root stage (the app). If this is
	// specified and this stage is nested within another stage, an error will be
	// thrown.
	// Default: - for nested stages, outdir will be determined as a relative
	// directory to the outdir of the app. For apps, if outdir is not specified, a
	// temporary directory will be created.
	//
	// Experimental.
	Outdir *string `field:"optional" json:"outdir" yaml:"outdir"`
	// Run the stage in a specific GitHub Environment.
	//
	// If specified,
	// any protection rules configured for the environment must pass
	// before the job is set to a runner. For example, if the environment
	// has a manual approval rule configured, then the workflow will
	// wait for the approval before sending the job to the runner.
	//
	// Running a workflow that references an environment that does not
	// exist will create an environment with the referenced name.
	// See: https://docs.github.com/en/actions/deployment/targeting-different-environments/using-environments-for-deployment
	//
	// Default: - no GitHub environment.
	//
	// Experimental.
	GitHubEnvironment *GitHubEnvironment `field:"optional" json:"gitHubEnvironment" yaml:"gitHubEnvironment"`
	// Job level settings that will be applied to all jobs in the stage.
	//
	// Currently the only valid setting is 'if'.
	// Experimental.
	JobSettings *JobSettings `field:"optional" json:"jobSettings" yaml:"jobSettings"`
	// In some cases, you must explicitly acknowledge that your CloudFormation stack template contains certain capabilities in order for CloudFormation to create the stack.
	//
	// If insufficiently specified, CloudFormation returns an `InsufficientCapabilities`
	// error.
	// Default: ['CAPABILITY_IAM'].
	//
	// Experimental.
	StackCapabilities *[]StackCapabilities `field:"optional" json:"stackCapabilities" yaml:"stackCapabilities"`
}

