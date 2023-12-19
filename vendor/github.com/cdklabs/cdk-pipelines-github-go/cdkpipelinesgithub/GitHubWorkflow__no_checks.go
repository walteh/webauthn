//go:build no_runtime_type_checking

package cdkpipelinesgithub

// Building without runtime type checking enabled, so all the below just return nil

func (g *jsiiProxy_GitHubWorkflow) validateAddGitHubWaveParameters(id *string, options *pipelines.WaveOptions) error {
	return nil
}

func (g *jsiiProxy_GitHubWorkflow) validateAddStageParameters(stage awscdk.Stage, options *pipelines.AddStageOpts) error {
	return nil
}

func (g *jsiiProxy_GitHubWorkflow) validateAddStageWithGitHubOptionsParameters(stage awscdk.Stage, options *AddGitHubStageOptions) error {
	return nil
}

func (g *jsiiProxy_GitHubWorkflow) validateAddWaveParameters(id *string, options *pipelines.WaveOptions) error {
	return nil
}

func validateGitHubWorkflow_IsConstructParameters(x interface{}) error {
	return nil
}

func validateNewGitHubWorkflowParameters(scope constructs.Construct, id *string, props *GitHubWorkflowProps) error {
	return nil
}

