//go:build no_runtime_type_checking

package cdkpipelinesgithub

// Building without runtime type checking enabled, so all the below just return nil

func validateDockerCredential_CustomRegistryParameters(registry *string, creds *ExternalDockerCredentialSecrets) error {
	return nil
}

func validateDockerCredential_DockerHubParameters(creds *DockerHubCredentialSecrets) error {
	return nil
}

func validateDockerCredential_EcrParameters(registry *string) error {
	return nil
}

