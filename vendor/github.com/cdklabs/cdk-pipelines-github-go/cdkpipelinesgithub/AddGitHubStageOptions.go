package cdkpipelinesgithub

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/pipelines"
)

// Options to pass to `addStageWithGitHubOpts`.
// Experimental.
type AddGitHubStageOptions struct {
	// Additional steps to run after all of the stacks in the stage.
	// Default: - No additional steps.
	//
	// Experimental.
	Post *[]pipelines.Step `field:"optional" json:"post" yaml:"post"`
	// Additional steps to run before any of the stacks in the stage.
	// Default: - No additional steps.
	//
	// Experimental.
	Pre *[]pipelines.Step `field:"optional" json:"pre" yaml:"pre"`
	// Instructions for stack level steps.
	// Default: - No additional instructions.
	//
	// Experimental.
	StackSteps *[]*pipelines.StackSteps `field:"optional" json:"stackSteps" yaml:"stackSteps"`
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

