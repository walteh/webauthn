//go:build !no_runtime_type_checking

package cdkpipelinesgithub

import (
	"fmt"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/pipelines"
)

func (g *jsiiProxy_GitHubWave) validateAddStageParameters(stage awscdk.Stage, options *pipelines.AddStageOpts) error {
	if stage == nil {
		return fmt.Errorf("parameter stage is required, but nil was provided")
	}

	if err := _jsii_.ValidateStruct(options, func() string { return "parameter options" }); err != nil {
		return err
	}

	return nil
}

func (g *jsiiProxy_GitHubWave) validateAddStageWithGitHubOptionsParameters(stage awscdk.Stage, options *AddGitHubStageOptions) error {
	if stage == nil {
		return fmt.Errorf("parameter stage is required, but nil was provided")
	}

	if err := _jsii_.ValidateStruct(options, func() string { return "parameter options" }); err != nil {
		return err
	}

	return nil
}

func validateNewGitHubWaveParameters(id *string, pipeline GitHubWorkflow, props *pipelines.WaveProps) error {
	if id == nil {
		return fmt.Errorf("parameter id is required, but nil was provided")
	}

	if pipeline == nil {
		return fmt.Errorf("parameter pipeline is required, but nil was provided")
	}

	if err := _jsii_.ValidateStruct(props, func() string { return "parameter props" }); err != nil {
		return err
	}

	return nil
}

