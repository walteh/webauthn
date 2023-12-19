//go:build no_runtime_type_checking

package cdkpipelinesgithub

// Building without runtime type checking enabled, so all the below just return nil

func (g *jsiiProxy_GitHubStage) validateSynthParameters(options *awscdk.StageSynthesisOptions) error {
	return nil
}

func validateGitHubStage_IsConstructParameters(x interface{}) error {
	return nil
}

func validateGitHubStage_IsStageParameters(x interface{}) error {
	return nil
}

func validateGitHubStage_OfParameters(construct constructs.IConstruct) error {
	return nil
}

func validateNewGitHubStageParameters(scope constructs.Construct, id *string, props *GitHubStageProps) error {
	return nil
}

