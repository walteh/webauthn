//go:build !no_runtime_type_checking

package cdkpipelinesgithub

import (
	"fmt"
)

func validateRunner_SelfHostedParameters(labels *[]*string) error {
	if labels == nil {
		return fmt.Errorf("parameter labels is required, but nil was provided")
	}

	return nil
}

