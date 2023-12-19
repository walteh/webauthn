//go:build no_runtime_type_checking

package cdkpipelinesgithub

// Building without runtime type checking enabled, so all the below just return nil

func (g *jsiiProxy_GitHubWave) validateAddStageParameters(stage awscdk.Stage, options *pipelines.AddStageOpts) error {
	return nil
}

func (g *jsiiProxy_GitHubWave) validateAddStageWithGitHubOptionsParameters(stage awscdk.Stage, options *AddGitHubStageOptions) error {
	return nil
}

func validateNewGitHubWaveParameters(id *string, pipeline GitHubWorkflow, props *pipelines.WaveProps) error {
	return nil
}

