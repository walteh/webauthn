//go:build no_runtime_type_checking

package cdknag

// Building without runtime type checking enabled, so all the below just return nil

func (s *jsiiProxy_SuppressionIgnoreErrors) validateCreateMessageParameters(input *SuppressionIgnoreInput) error {
	return nil
}

