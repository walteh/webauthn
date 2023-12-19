//go:build !no_runtime_type_checking

package cdkpipelinesgithub

import (
	"fmt"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func (y *jsiiProxy_YamlFile) validateUpdateParameters(obj interface{}) error {
	if obj == nil {
		return fmt.Errorf("parameter obj is required, but nil was provided")
	}

	return nil
}

func validateNewYamlFileParameters(filePath *string, options *YamlFileOptions) error {
	if filePath == nil {
		return fmt.Errorf("parameter filePath is required, but nil was provided")
	}

	if err := _jsii_.ValidateStruct(options, func() string { return "parameter options" }); err != nil {
		return err
	}

	return nil
}

