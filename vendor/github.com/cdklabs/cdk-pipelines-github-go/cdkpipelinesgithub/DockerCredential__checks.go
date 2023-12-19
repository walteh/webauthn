//go:build !no_runtime_type_checking

package cdkpipelinesgithub

import (
	"fmt"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func validateDockerCredential_CustomRegistryParameters(registry *string, creds *ExternalDockerCredentialSecrets) error {
	if registry == nil {
		return fmt.Errorf("parameter registry is required, but nil was provided")
	}

	if creds == nil {
		return fmt.Errorf("parameter creds is required, but nil was provided")
	}
	if err := _jsii_.ValidateStruct(creds, func() string { return "parameter creds" }); err != nil {
		return err
	}

	return nil
}

func validateDockerCredential_DockerHubParameters(creds *DockerHubCredentialSecrets) error {
	if err := _jsii_.ValidateStruct(creds, func() string { return "parameter creds" }); err != nil {
		return err
	}

	return nil
}

func validateDockerCredential_EcrParameters(registry *string) error {
	if registry == nil {
		return fmt.Errorf("parameter registry is required, but nil was provided")
	}

	return nil
}

