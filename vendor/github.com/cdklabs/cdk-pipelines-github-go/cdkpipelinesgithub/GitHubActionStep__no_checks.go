//go:build no_runtime_type_checking

package cdkpipelinesgithub

// Building without runtime type checking enabled, so all the below just return nil

func (g *jsiiProxy_GitHubActionStep) validateAddDependencyFileSetParameters(fs pipelines.FileSet) error {
	return nil
}

func (g *jsiiProxy_GitHubActionStep) validateAddStepDependencyParameters(step pipelines.Step) error {
	return nil
}

func (g *jsiiProxy_GitHubActionStep) validateConfigurePrimaryOutputParameters(fs pipelines.FileSet) error {
	return nil
}

func validateGitHubActionStep_SequenceParameters(steps *[]pipelines.Step) error {
	return nil
}

func validateNewGitHubActionStepParameters(id *string, props *GitHubActionStepProps) error {
	return nil
}

