//go:build !no_runtime_type_checking

package cdkpipelinesgithub

import (
	"fmt"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk/v2/pipelines"
)

func (g *jsiiProxy_GitHubActionStep) validateAddDependencyFileSetParameters(fs pipelines.FileSet) error {
	if fs == nil {
		return fmt.Errorf("parameter fs is required, but nil was provided")
	}

	return nil
}

func (g *jsiiProxy_GitHubActionStep) validateAddStepDependencyParameters(step pipelines.Step) error {
	if step == nil {
		return fmt.Errorf("parameter step is required, but nil was provided")
	}

	return nil
}

func (g *jsiiProxy_GitHubActionStep) validateConfigurePrimaryOutputParameters(fs pipelines.FileSet) error {
	if fs == nil {
		return fmt.Errorf("parameter fs is required, but nil was provided")
	}

	return nil
}

func validateGitHubActionStep_SequenceParameters(steps *[]pipelines.Step) error {
	if steps == nil {
		return fmt.Errorf("parameter steps is required, but nil was provided")
	}

	return nil
}

func validateNewGitHubActionStepParameters(id *string, props *GitHubActionStepProps) error {
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

