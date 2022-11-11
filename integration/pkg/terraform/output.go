package terraform

import (
	"reflect"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func LoadOutput[I interface{}](t *testing.T) *I {

	opt := &terraform.Options{
		TerraformDir:    "../terraform",
		TerraformBinary: "terraform",
		Logger:          logger.Discard,
	}

	plan, _ := terraform.ShowWithStructE(t, opt)

	assert.Empty(t, plan.ResourceChangesMap, "ResourceChangesMap is not empty")

	var output I

	abc := reflect.TypeOf(output)

	for i := 0; i < abc.NumField(); i++ {
		val := terraform.Output(t, opt, abc.Field(i).Tag.Get("json"))

		reflect.ValueOf(&output).Elem().Field(i).SetString(val)
	}

	return &output
}
