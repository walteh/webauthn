//go:build !no_runtime_type_checking

package cdkpipelinesgithub

import (
	"fmt"
)

func (a *jsiiProxy_AwsCredentialsProvider) validateCredentialStepsParameters(region *string) error {
	if region == nil {
		return fmt.Errorf("parameter region is required, but nil was provided")
	}

	return nil
}

