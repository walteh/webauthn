package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsbackup"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awskms"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/cdklabs/cdk-nag-go/cdknag/v2"
	"github.com/k0kubun/pp/v3"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssns"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssnssubscriptions"

	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/walteh/webauthn/constants"
)

type CdkStackProps struct {
	awscdk.StackProps
}

// func grantPriciple(stack awscdk.Stack, service *string, f func(awsiam.IGrantable) awsiam.Grant, source *string) {
// 	// princ := awsiam.NewServicePrincipal(service, &awsiam.ServicePrincipalOpts{
// 	// 	Region: stack.Region(),
// 	// })

// 	// grant := f(princ)

// 	// statements := *grant.ResourceStatements()
// 	// for i := range statements {
// 	// 	statements[i].AddAccountCondition(stack.Account())
// 	// 	statements[i].AddSourceArnCondition(source)
// 	// }
// 	// *grant.ResourceStatements() = statements
// 	// grant.AssertSuccess()

// 	// scp := awskms.NewViaServicePrincipal(service, awsiam.NewArnPrincipal(source))
// 	f(awsiam.NewArnPrincipal(source)).AssertSuccess()
// }

func NewCdkStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps

	}

	stack := awscdk.NewStack(scope, jsii.String("create-ceremony"), &sprops)

	lamkey := awskms.NewKey(stack, jsii.String("key"), &awskms.KeyProps{
		EnableKeyRotation: jsii.Bool(true),
		Enabled:           jsii.Bool(true),
	})

	backupvault := awsbackup.NewBackupVault(stack, jsii.String("backup-vault"), &awsbackup.BackupVaultProps{
		BackupVaultName: jsii.String("create-ceremony"),
		EncryptionKey:   lamkey,
	})

	lamkey.GrantEncryptDecrypt(awsiam.NewArnPrincipal(backupvault.BackupVaultArn())).AssertSuccess()

	rule := awsbackup.BackupPlanRule_Daily(backupvault)

	backupplan := awsbackup.NewBackupPlan(stack, jsii.String("backup-plan"), &awsbackup.BackupPlanProps{
		BackupPlanName: jsii.String("create-ceremony"),
		BackupPlanRules: &[]awsbackup.BackupPlanRule{
			rule,
		},
		BackupVault: backupvault,
	})

	// The code that defines your stack goes here
	tbl := awsdynamodb.NewTable(stack, jsii.String("table"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode:         awsdynamodb.BillingMode_PAY_PER_REQUEST,
		PointInTimeRecovery: jsii.Bool(true),
		TableName:           jsii.String("create-ceremony"),
		EncryptionKey:       lamkey,
		Encryption:          awsdynamodb.TableEncryption_CUSTOMER_MANAGED,
	})

	lamkey.GrantEncryptDecrypt(awsiam.NewArnPrincipal(tbl.TableArn())).AssertSuccess()

	bu := awsbackup.BackupResource_FromDynamoDbTable(tbl)

	fancy := awscdk.Arn_Format(&awscdk.ArnComponents{
		Service:      jsii.String("dynamodb"),
		Resource:     jsii.String("table"),
		ResourceName: jsii.String(*tbl.TableName() + "/backup/*"),
	}, stack)

	buboundry := awsiam.NewManagedPolicy(stack, jsii.String("boundary"), &awsiam.ManagedPolicyProps{
		Statements: &[]awsiam.PolicyStatement{
			awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
				Actions: &[]*string{
					jsii.String("dynamodb:*"),
				},
				Resources: &[]*string{
					fancy,
					tbl.TableArn(),
				},
			}),
		},
	})

	cdknag.NagSuppressions_AddResourceSuppressions(buboundry, &[]*cdknag.NagPackSuppression{
		{
			Id:     jsii.String("AwsSolutions-IAM5"),
			Reason: jsii.String("This is a boundary policy and is ok to have full access"),
		},
		{
			Id:     jsii.String("HIPAA.Security-IAMPolicyNoStatementsWithFullAccess"),
			Reason: jsii.String("This is a boundary policy and is ok to have full access"),
		},
		{
			Id:     jsii.String("PCI.DSS.321-IAMPolicyNoStatementsWithFullAccess"),
			Reason: jsii.String("This is a boundary policy and is ok to have full access"),
		},
	}, jsii.Bool(true))

	burole := awsiam.NewRole(stack, jsii.String("db-role"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("backup.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("service-role/AWSBackupServiceRolePolicyForBackup")),
		},
		PermissionsBoundary: buboundry,
	})

	cdknag.NagSuppressions_AddResourceSuppressions(burole, &[]*cdknag.NagPackSuppression{
		{
			Id:     jsii.String("AwsSolutions-IAM4"),
			Reason: jsii.String("We have applied a boundry to this aws managed policy"),
			AppliesTo: &[]any{
				jsii.String("Policy::arn:<AWS::Partition>:iam::aws:policy/service-role/AWSBackupServiceRolePolicyForBackup"),
			},
		},
	}, jsii.Bool(true))

	backupplan.AddSelection(jsii.String("table"), &awsbackup.BackupSelectionOptions{
		Resources:                  &[]awsbackup.BackupResource{bu},
		DisableDefaultBackupPolicy: jsii.Bool(true),
		Role:                       burole,
	})

	lamkey.GrantEncryptDecrypt(burole).AssertSuccess()

	lamgrp := awslogs.NewLogGroup(stack, jsii.String("log-group"), &awslogs.LogGroupProps{
		LogGroupName:  jsii.String("/aws/lambda/create-ceremony"),
		Retention:     awslogs.RetentionDays_ONE_WEEK,
		EncryptionKey: lamkey,
	})

	lamkey.GrantEncryptDecrypt(awsiam.NewArnPrincipal(lamgrp.LogGroupArn())).AssertSuccess()

	snsdlq := awssns.NewTopic(stack, jsii.String("sns-dlq"), &awssns.TopicProps{
		DisplayName: jsii.String("create-ceremony-dlq"),
		MasterKey:   lamkey,
	})

	lamkey.GrantEncryptDecrypt(awsiam.NewArnPrincipal(snsdlq.TopicArn())).AssertSuccess()

	emailsub := awssnssubscriptions.NewEmailSubscription(jsii.String("walter@nugg.xyz"), &awssnssubscriptions.EmailSubscriptionProps{
		Json: jsii.Bool(true),
	})

	snsdlq.AddSubscription(emailsub)

	lamrole := awsiam.NewLazyRole(stack, jsii.String("lambda-role"), &awsiam.LazyRoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("lambda.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
	})

	funct := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("lambda"), &awscdklambdagoalpha.GoFunctionProps{
		Entry:        jsii.String("../cmd/lambda/init/main.go"),
		FunctionName: jsii.String("create-ceremony"),
		CurrentVersionOptions: &awslambda.VersionOptions{
			ProvisionedConcurrentExecutions: jsii.Number(0),
		},
		ReservedConcurrentExecutions: jsii.Number(400),
		Architecture:                 awslambda.Architecture_ARM_64(),
		Runtime:                      awslambda.Runtime_PROVIDED_AL2023(),
		Environment: &map[string]*string{
			constants.EnvVarCeremonyTableName: tbl.TableName(),
		},
		LogGroup:        lamgrp,
		DeadLetterTopic: snsdlq,
		Role:            lamrole,
	})

	cdknag.NagSuppressions_AddResourceSuppressions(lamrole, &[]*cdknag.NagPackSuppression{
		{
			Id:     jsii.String("HIPAA.Security-IAMNoInlinePolicy"),
			Reason: jsii.String("This role is ok to have an inline policy"),
		},
		{
			Id:     jsii.String("PCI.DSS.321-IAMNoInlinePolicy"),
			Reason: jsii.String("This role is ok to have an inline policy"),
		},
	}, jsii.Bool(true))

	lamgrp.GrantWrite(funct.Role())

	funct.ApplyRemovalPolicy(awscdk.RemovalPolicy_DESTROY)

	tbl.Grant(funct, jsii.String("dynamodb:TransactWriteItems"), jsii.String("dynamodb:TransactGetItems"))

	laminteg := awsapigatewayv2integrations.NewHttpLambdaIntegration(jsii.String("integration"), funct, &awsapigatewayv2integrations.HttpLambdaIntegrationProps{
		ParameterMapping:     awsapigatewayv2.NewParameterMapping(),
		PayloadFormatVersion: awsapigatewayv2.PayloadFormatVersion_VERSION_2_0(),
	})

	api := awsapigatewayv2.NewHttpApi(funct, jsii.String("api"), &awsapigatewayv2.HttpApiProps{
		ApiName:                   jsii.String("create-ceremony"),
		DefaultIntegration:        laminteg,
		CreateDefaultStage:        jsii.Bool(true),
		DisableExecuteApiEndpoint: jsii.Bool(false),
	})

	dstg := api.DefaultStage()
	dstgcnf := dstg.Node().DefaultChild().(awsapigatewayv2.CfnStage)

	apiloggroup := awslogs.NewLogGroup(api, jsii.String("log-group"), &awslogs.LogGroupProps{
		LogGroupName:  jsii.String("/aws/apigateway/create-ceremony"),
		Retention:     awslogs.RetentionDays_ONE_WEEK,
		EncryptionKey: lamkey,
	})

	dstgcnf.SetAccessLogSettings(&awsapigatewayv2.CfnStage_AccessLogSettingsProperty{
		DestinationArn: apiloggroup.LogGroupArn(),
		Format:         jsii.String(`{"requestId":"$context.requestId","ip":"$context.identity.sourceIp","caller":"$context.identity.caller","user":"$context.identity.user","requestTime":"$context.requestTime","httpMethod":"$context.httpMethod","resourcePath":"$context.resourcePath","status":"$context.status","protocol":"$context.protocol","responseLength":"$context.responseLength"}`),
	})

	route := awsapigatewayv2.NewHttpRoute(api, jsii.String("default-route"), &awsapigatewayv2.HttpRouteProps{
		HttpApi:     api,
		RouteKey:    awsapigatewayv2.HttpRouteKey_DEFAULT(),
		Integration: laminteg,
	})

	funct.AddPermission(jsii.String("apigateway"), &awslambda.Permission{
		Principal: awsiam.NewServicePrincipal(jsii.String("apigateway.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
		SourceArn: route.RouteArn(),
	})

	// ignore nag for auth as this is the auth api
	cdknag.NagSuppressions_AddResourceSuppressions(api, &[]*cdknag.NagPackSuppression{
		{
			Id:     jsii.String("AwsSolutions-APIG4"),
			Reason: jsii.String("This is the auth api and does not require auth"),
		},
	}, jsii.Bool(true))

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

	assmbly := app.Synth(nil)
	pp.Println(assmbly.Tree())
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
