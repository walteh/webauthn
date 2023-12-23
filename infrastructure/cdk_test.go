package main

import (
	"fmt"
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdklabs/cdk-nag-go/cdknag/v2"
	"github.com/stretchr/testify/assert"
)

// example tests. To run these tests, uncomment this file along with the
// example resource in cdk_test.go
func TestCdkStack(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)

	// WHEN
	stack := NewCdkStack(app, "MyStack", nil)

	// THEN
	template := assertions.Template_FromStack(stack, nil)

	template.HasResourceProperties(jsii.String("AWS::Lambda::Function"), map[string]interface{}{
		"Architectures": []interface{}{
			"arm64",
		},
	})
}

func TestNag(t *testing.T) {
	app := awscdk.NewApp(nil)

	// WHEN
	stack := NewCdkStack(app, "MyStack", nil)

	nag := cdknag.NewAwsSolutionsChecks(&cdknag.NagPackProps{
		Verbose: jsii.Bool(true),
	})

	nag2 := cdknag.NewHIPAASecurityChecks(&cdknag.NagPackProps{
		Verbose: jsii.Bool(true),
	})

	nag3 := cdknag.NewPCIDSS321Checks(&cdknag.NagPackProps{
		Verbose: jsii.Bool(true),
	})

	awscdk.Aspects_Of(stack).Add(nag)
	awscdk.Aspects_Of(stack).Add(nag2)
	awscdk.Aspects_Of(stack).Add(nag3)

	anots := assertions.Annotations_FromStack(stack)

	errs := *anots.FindError(jsii.String("*"), assertions.Match_AnyValue())
	if !assert.Empty(t, errs) {
		fmt.Println("====== ERRORS ======")
		for i, err := range errs {
			fmt.Printf("(%d) [%s]:\n %s\n", i, *err.Id, err.Entry.Data)
		}
		fmt.Println("====================")
	}

	errs = *anots.FindWarning(jsii.String("*"), assertions.Match_AnyValue())
	if !assert.Empty(t, errs) {
		fmt.Println("====== WARNINGS ======")
		for i, err := range errs {
			fmt.Printf("(%d) [%s]:\n %s\n", i, *err.Id, err.Entry.Data)
		}
		fmt.Println("======================")
	}
}
