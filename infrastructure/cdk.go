package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"

	// "github.com/aws/aws-cdk-go/awscdk/v2/aws

	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CdkStackProps struct {
	awscdk.StackProps
}

func NewCdkStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here

	tbl := awsdynamodb.NewTable(stack, jsii.String("ceremony"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode: awsdynamodb.BillingMode_PAY_PER_REQUEST,
	})

	// // example resource
	funct := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("create-ceremony"), &awscdklambdagoalpha.GoFunctionProps{
		Entry:                 jsii.String("../cmd/lambda/init/main.go"),
		CurrentVersionOptions: &awslambda.VersionOptions{},
		Architecture:          awslambda.Architecture_ARM_64(),
		Runtime:               awslambda.Runtime_PROVIDED_AL2023(),
		Environment: &map[string]*string{
			"CEREMONY_TABLE_NAME": tbl.TableName(),
		},
	})

	tbl.Grant(funct, jsii.String("dynamodb:TransactWriteItems"), jsii.String("dynamodb:TransactGetItems"))

	integ2 := awsapigatewayv2integrations.NewHttpLambdaIntegration(jsii.String("create-ceremony"), funct, &awsapigatewayv2integrations.HttpLambdaIntegrationProps{
		ParameterMapping:     awsapigatewayv2.NewParameterMapping(),
		PayloadFormatVersion: awsapigatewayv2.PayloadFormatVersion_VERSION_2_0(),
	})

	api := awsapigatewayv2.NewHttpApi(stack, jsii.String("CeremonyApi"), &awsapigatewayv2.HttpApiProps{
		ApiName:            jsii.String("CeremonyApi"),
		DefaultIntegration: integ2,
	})

	route := awsapigatewayv2.NewHttpRoute(stack, jsii.String("CreateCeremony"), &awsapigatewayv2.HttpRouteProps{
		HttpApi: api,
	})

	awsapigatewayv2.NewHttpStage(stack, jsii.String("$default"), &awsapigatewayv2.HttpStageProps{
		HttpApi:    api,
		StageName:  jsii.String("$default"),
		AutoDeploy: jsii.Bool(true),
	})

	funct.AddPermission(jsii.String("apigateway"), &awslambda.Permission{
		Principal: awsiam.NewServicePrincipal(jsii.String("apigateway.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
		SourceArn: route.RouteArn(),
	})

	api.AddStage(jsii.String("$default"), &awsapigatewayv2.HttpStageOptions{
		AutoDeploy: jsii.Bool(true),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewCdkStack(app, "CdkStack", &CdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
