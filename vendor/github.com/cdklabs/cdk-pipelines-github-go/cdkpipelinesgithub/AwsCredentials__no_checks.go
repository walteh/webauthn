//go:build no_runtime_type_checking

package cdkpipelinesgithub

// Building without runtime type checking enabled, so all the below just return nil

func validateAwsCredentials_FromGitHubSecretsParameters(props *GitHubSecretsProviderProps) error {
	return nil
}

func validateAwsCredentials_FromOpenIdConnectParameters(props *OpenIdConnectProviderProps) error {
	return nil
}

