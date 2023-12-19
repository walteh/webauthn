//go:build no_runtime_type_checking

package cdkpipelinesgithub

// Building without runtime type checking enabled, so all the below just return nil

func validateGitHubActionRole_ExistingGitHubActionsProviderParameters(scope constructs.Construct) error {
	return nil
}

func validateGitHubActionRole_IsConstructParameters(x interface{}) error {
	return nil
}

func validateNewGitHubActionRoleParameters(scope constructs.Construct, id *string, props *GitHubActionRoleProps) error {
	return nil
}

