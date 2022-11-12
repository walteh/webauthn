package terraform

import (
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strings"
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

func LoadOutputE[I interface{}](t *testing.T, dir string) (*I, error) {

	opt := &terraform.Options{
		TerraformDir:    dir,
		TerraformBinary: "terraform",
		Logger:          logger.Discard,
	}

	plan, err := terraform.ShowWithStructE(t, opt)
	if err != nil {
		return nil, err
	}

	if len(plan.ResourceChangesMap) != 0 {
		return nil, fmt.Errorf("ResourceChangesMap is not empty")
	}

	var output I

	abc := reflect.TypeOf(output)

	for i := 0; i < abc.NumField(); i++ {
		val, err := terraform.OutputE(t, opt, abc.Field(i).Tag.Get("json"))
		if err != nil {
			return nil, err
		}

		reflect.ValueOf(&output).Elem().Field(i).SetString(val)
	}

	return &output, nil
}

func BuildAppsyncExample(t *testing.T, appsync_authorizer_name string) (*url.URL, *url.URL, func()) {

	opt := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir:    "../terraform/examples/appsync",
		TerraformBinary: "terraform",
		Vars: map[string]interface{}{
			"appsync_authorizer_function_name": appsync_authorizer_name,
		},
	})

	teardown := func() {
		terraform.Destroy(t, opt)
	}

	teardownWithError := func(err error) {
		teardown()
		t.Fatal(err)
	}

	std, err := terraform.InitAndApplyE(t, opt)
	if err != nil {
		log.Println(std)
		teardownWithError(err)
	}

	type Output struct {
		Endpoint string `json:"appsync_graphql_api_endpoint"`
	}

	out, err := LoadOutputE[Output](t, "../terraform/examples/appsync")
	if err != nil {
		teardownWithError(err)
	}

	endpoint, err := url.Parse(out.Endpoint)
	if err != nil {
		teardownWithError(err)
	}

	realtimeEndpoint := *endpoint
	realtimeEndpoint.Scheme = "wss"
	realtimeEndpoint.Host = strings.Replace(endpoint.Host, "appsync-api", "appsync-realtime-api", 1)

	return endpoint, &realtimeEndpoint, func() {
		terraform.Destroy(t, opt)
	}

}
