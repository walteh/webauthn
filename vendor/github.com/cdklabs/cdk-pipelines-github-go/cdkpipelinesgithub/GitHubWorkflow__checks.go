//go:build !no_runtime_type_checking

package cdkpipelinesgithub

import (
	"fmt"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/pipelines"
	"github.com/aws/constructs-go/constructs/v10"
)

func (g *jsiiProxy_GitHubWorkflow) validateAddGitHubWaveParameters(id *string, options *pipelines.WaveOptions) error {
	if id == nil {
		return fmt.Errorf("parameter id is required, but nil was provided")
	}

	if err := _jsii_.ValidateStruct(options, func() string { return "parameter options" }); err != nil {
		return err
	}

	return nil
}

func (g *jsiiProxy_GitHubWorkflow) validateAddStageParameters(stage awscdk.Stage, options *pipelines.AddStageOpts) error {
	if stage == nil {
		return fmt.Errorf("parameter stage is required, but nil was provided")
	}

	if err := _jsii_.ValidateStruct(options, func() string { return "parameter options" }); err != nil {
		return err
	}

	return nil
}

func (g *jsiiProxy_GitHubWorkflow) validateAddStageWithGitHubOptionsParameters(stage awscdk.Stage, options *AddGitHubStageOptions) error {
	if stage == nil {
		return fmt.Errorf("parameter stage is required, but nil was provided")
	}

	if err := _jsii_.ValidateStruct(options, func() string { return "parameter options" }); err != nil {
		return err
	}

	return nil
}

func (g *jsiiProxy_GitHubWorkflow) validateAddWaveParameters(id *string, options *pipelines.WaveOptions) error {
	if id == nil {
		return fmt.Errorf("parameter id is required, but nil was provided")
	}

	if err := _jsii_.ValidateStruct(options, func() string { return "parameter options" }); err != nil {
		return err
	}

	return nil
}

func validateGitHubWorkflow_IsConstructParameters(x interface{}) error {
	if x == nil {
		return fmt.Errorf("parameter x is required, but nil was provided")
	}

	return nil
}

func validateNewGitHubWorkflowParameters(scope constructs.Construct, id *string, props *GitHubWorkflowProps) error {
	if scope == nil {
		return fmt.Errorf("parameter scope is required, but nil was provided")
	}

	if id == nil {
		return fmt.Errorf("parameter id is required, but nil was provided")
	}

	if props == nil {
		return fmt.Errorf("parameter props is required, but nil was provided")
	}
	if err := _jsii_.ValidateStruct(props, func() string { return "parameter props" }); err != nil {
		return err
	}

	return nil
}

