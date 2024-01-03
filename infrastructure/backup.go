package main

// lamkey.GrantEncryptDecrypt(awsiam.NewArnPrincipal(backupvault.BackupVaultArn())).AssertSuccess()

// lamkey.GrantEncryptDecrypt(awsiam.NewArnPrincipal(awscdk.Arn_Format(&awscdk.ArnComponents{
// 	Service:      jsii.String("iam"),
// 	Resource:     jsii.String("role"),
// 	ResourceName: jsii.String("service-role/AWSBackupDefaultServiceRole"),
// }, stack))).AssertSuccess()

// backupvault := awsbackup.NewBackupVault(stack, jsii.String("backup-vault"), &awsbackup.BackupVaultProps{
// 	BackupVaultName: jsii.String("create-ceremony"),
// 	EncryptionKey:   lamkey,
// })

// rule := awsbackup.BackupPlanRule_Daily(backupvault)

// backupplan := awsbackup.NewBackupPlan(stack, jsii.String("backup-plan"), &awsbackup.BackupPlanProps{
// 	BackupPlanName: jsii.String("create-ceremony"),
// 	BackupPlanRules: &[]awsbackup.BackupPlanRule{
// 		rule,
// 	},
// 	BackupVault: backupvault,
// })

// buboundry := awsiam.NewManagedPolicy(stack, jsii.String("boundary"), &awsiam.ManagedPolicyProps{
// 	Statements: &[]awsiam.PolicyStatement{
// 		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
// 			Actions: &[]*string{
// 				jsii.String("dynamodb:*"),
// 			},
// 			Resources: &[]*string{
// 				fancy,
// 				tbl.TableArn(),
// 			},
// 		}),
// 	},
// })

// cdknag.NagSuppressions_AddResourceSuppressions(buboundry, &[]*cdknag.NagPackSuppression{
// 	{
// 		Id:     jsii.String("AwsSolutions-IAM5"),
// 		Reason: jsii.String("This is a boundary policy and is ok to have full access"),
// 	},
// 	{
// 		Id:     jsii.String("HIPAA.Security-IAMPolicyNoStatementsWithFullAccess"),
// 		Reason: jsii.String("This is a boundary policy and is ok to have full access"),
// 	},
// 	{
// 		Id:     jsii.String("PCI.DSS.321-IAMPolicyNoStatementsWithFullAccess"),
// 		Reason: jsii.String("This is a boundary policy and is ok to have full access"),
// 	},
// }, jsii.Bool(true))

// burole := awsiam.NewRole(stack, jsii.String("db-role"), &awsiam.RoleProps{
// 	AssumedBy: awsiam.NewServicePrincipal(jsii.String("backup.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
// 	ManagedPolicies: &[]awsiam.IManagedPolicy{
// 		awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("service-role/AWSBackupServiceRolePolicyForBackup")),
// 	},
// 	PermissionsBoundary: buboundry,
// })

// cdknag.NagSuppressions_AddResourceSuppressions(burole, &[]*cdknag.NagPackSuppression{
// 	{
// 		Id:     jsii.String("AwsSolutions-IAM5"),
// 		Reason: jsii.String("This * access is still restricted enough"),
// 		AppliesTo: &[]any{
// 			jsii.String("Action::kms:ReEncrypt*"),
// 			jsii.String("Action::kms:GenerateDataKey*"),
// 		},
// 	},
// 	{
// 		Id:     jsii.String("AwsSolutions-IAM4"),
// 		Reason: jsii.String("We have applied a boundry to this aws managed policy"),
// 		AppliesTo: &[]any{
// 			jsii.String("Policy::arn:<AWS::Partition>:iam::aws:policy/service-role/AWSBackupServiceRolePolicyForBackup"),
// 		},
// 	},
// 	{
// 		Id:     jsii.String("HIPAA.Security-IAMNoInlinePolicy"),
// 		Reason: jsii.String("We have a boundary policy to control access to this inline policy"),
// 	},
// 	{
// 		Id:     jsii.String("PCI.DSS.321-IAMNoInlinePolicy"),
// 		Reason: jsii.String("We have a boundary policy to control access to this inline policy"),
// 	},
// }, jsii.Bool(true))

// backupplan.AddSelection(jsii.String("table"), &awsbackup.BackupSelectionOptions{
// 	Resources:                  &[]awsbackup.BackupResource{bu},
// 	DisableDefaultBackupPolicy: jsii.Bool(true),
// 	Role:                       burole,
// })

// lamkey.GrantEncryptDecrypt(burole).AssertSuccess()

// bu := awsbackup.BackupResource_FromDynamoDbTable(tbl)

// fancy := awscdk.Arn_Format(&awscdk.ArnComponents{
// 	Service:      jsii.String("dynamodb"),
// 	Resource:     jsii.String("table"),
// 	ResourceName: jsii.String(*tbl.TableName() + "/backup/*"),
// }, stack)
