//go:build no_runtime_type_checking

package cdkpipelinesgithub

// Building without runtime type checking enabled, so all the below just return nil

func (y *jsiiProxy_YamlFile) validateUpdateParameters(obj interface{}) error {
	return nil
}

func validateNewYamlFileParameters(filePath *string, options *YamlFileOptions) error {
	return nil
}

