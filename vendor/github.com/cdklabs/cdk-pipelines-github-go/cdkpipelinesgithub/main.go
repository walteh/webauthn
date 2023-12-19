// GitHub Workflows support for CDK Pipelines
package cdkpipelinesgithub

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.AddGitHubStageOptions",
		reflect.TypeOf((*AddGitHubStageOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"cdk-pipelines-github.AwsCredentials",
		reflect.TypeOf((*AwsCredentials)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_AwsCredentials{}
		},
	)
	_jsii_.RegisterClass(
		"cdk-pipelines-github.AwsCredentialsProvider",
		reflect.TypeOf((*AwsCredentialsProvider)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "credentialSteps", GoMethod: "CredentialSteps"},
			_jsii_.MemberMethod{JsiiMethod: "jobPermission", GoMethod: "JobPermission"},
		},
		func() interface{} {
			return &jsiiProxy_AwsCredentialsProvider{}
		},
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.AwsCredentialsSecrets",
		reflect.TypeOf((*AwsCredentialsSecrets)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.CheckRunOptions",
		reflect.TypeOf((*CheckRunOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.CheckSuiteOptions",
		reflect.TypeOf((*CheckSuiteOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.ContainerCredentials",
		reflect.TypeOf((*ContainerCredentials)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.ContainerOptions",
		reflect.TypeOf((*ContainerOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.CreateOptions",
		reflect.TypeOf((*CreateOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.CronScheduleOptions",
		reflect.TypeOf((*CronScheduleOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.DeleteOptions",
		reflect.TypeOf((*DeleteOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.DeploymentOptions",
		reflect.TypeOf((*DeploymentOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.DeploymentStatusOptions",
		reflect.TypeOf((*DeploymentStatusOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"cdk-pipelines-github.DockerCredential",
		reflect.TypeOf((*DockerCredential)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "name", GoGetter: "Name"},
			_jsii_.MemberProperty{JsiiProperty: "passwordKey", GoGetter: "PasswordKey"},
			_jsii_.MemberProperty{JsiiProperty: "registry", GoGetter: "Registry"},
			_jsii_.MemberProperty{JsiiProperty: "usernameKey", GoGetter: "UsernameKey"},
		},
		func() interface{} {
			return &jsiiProxy_DockerCredential{}
		},
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.DockerHubCredentialSecrets",
		reflect.TypeOf((*DockerHubCredentialSecrets)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.ExternalDockerCredentialSecrets",
		reflect.TypeOf((*ExternalDockerCredentialSecrets)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.ForkOptions",
		reflect.TypeOf((*ForkOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"cdk-pipelines-github.GitHubActionRole",
		reflect.TypeOf((*GitHubActionRole)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "role", GoGetter: "Role"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
		},
		func() interface{} {
			j := jsiiProxy_GitHubActionRole{}
			_jsii_.InitJsiiProxy(&j.Type__constructsConstruct)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.GitHubActionRoleProps",
		reflect.TypeOf((*GitHubActionRoleProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"cdk-pipelines-github.GitHubActionStep",
		reflect.TypeOf((*GitHubActionStep)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDependencyFileSet", GoMethod: "AddDependencyFileSet"},
			_jsii_.MemberMethod{JsiiMethod: "addStepDependency", GoMethod: "AddStepDependency"},
			_jsii_.MemberMethod{JsiiMethod: "configurePrimaryOutput", GoMethod: "ConfigurePrimaryOutput"},
			_jsii_.MemberProperty{JsiiProperty: "dependencies", GoGetter: "Dependencies"},
			_jsii_.MemberProperty{JsiiProperty: "dependencyFileSets", GoGetter: "DependencyFileSets"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberProperty{JsiiProperty: "id", GoGetter: "Id"},
			_jsii_.MemberProperty{JsiiProperty: "isSource", GoGetter: "IsSource"},
			_jsii_.MemberProperty{JsiiProperty: "jobSteps", GoGetter: "JobSteps"},
			_jsii_.MemberProperty{JsiiProperty: "primaryOutput", GoGetter: "PrimaryOutput"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
		},
		func() interface{} {
			j := jsiiProxy_GitHubActionStep{}
			_jsii_.InitJsiiProxy(&j.Type__pipelinesStep)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.GitHubActionStepProps",
		reflect.TypeOf((*GitHubActionStepProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.GitHubCommonProps",
		reflect.TypeOf((*GitHubCommonProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.GitHubEnvironment",
		reflect.TypeOf((*GitHubEnvironment)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.GitHubSecretsProviderProps",
		reflect.TypeOf((*GitHubSecretsProviderProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"cdk-pipelines-github.GitHubStage",
		reflect.TypeOf((*GitHubStage)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "account", GoGetter: "Account"},
			_jsii_.MemberProperty{JsiiProperty: "artifactId", GoGetter: "ArtifactId"},
			_jsii_.MemberProperty{JsiiProperty: "assetOutdir", GoGetter: "AssetOutdir"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "outdir", GoGetter: "Outdir"},
			_jsii_.MemberProperty{JsiiProperty: "parentStage", GoGetter: "ParentStage"},
			_jsii_.MemberProperty{JsiiProperty: "props", GoGetter: "Props"},
			_jsii_.MemberProperty{JsiiProperty: "region", GoGetter: "Region"},
			_jsii_.MemberProperty{JsiiProperty: "stageName", GoGetter: "StageName"},
			_jsii_.MemberMethod{JsiiMethod: "synth", GoMethod: "Synth"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
		},
		func() interface{} {
			j := jsiiProxy_GitHubStage{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkStage)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.GitHubStageProps",
		reflect.TypeOf((*GitHubStageProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"cdk-pipelines-github.GitHubWave",
		reflect.TypeOf((*GitHubWave)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addPost", GoMethod: "AddPost"},
			_jsii_.MemberMethod{JsiiMethod: "addPre", GoMethod: "AddPre"},
			_jsii_.MemberMethod{JsiiMethod: "addStage", GoMethod: "AddStage"},
			_jsii_.MemberMethod{JsiiMethod: "addStageWithGitHubOptions", GoMethod: "AddStageWithGitHubOptions"},
			_jsii_.MemberProperty{JsiiProperty: "id", GoGetter: "Id"},
			_jsii_.MemberProperty{JsiiProperty: "post", GoGetter: "Post"},
			_jsii_.MemberProperty{JsiiProperty: "pre", GoGetter: "Pre"},
			_jsii_.MemberProperty{JsiiProperty: "stages", GoGetter: "Stages"},
		},
		func() interface{} {
			j := jsiiProxy_GitHubWave{}
			_jsii_.InitJsiiProxy(&j.Type__pipelinesWave)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"cdk-pipelines-github.GitHubWorkflow",
		reflect.TypeOf((*GitHubWorkflow)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addGitHubWave", GoMethod: "AddGitHubWave"},
			_jsii_.MemberMethod{JsiiMethod: "addStage", GoMethod: "AddStage"},
			_jsii_.MemberMethod{JsiiMethod: "addStageWithGitHubOptions", GoMethod: "AddStageWithGitHubOptions"},
			_jsii_.MemberMethod{JsiiMethod: "addWave", GoMethod: "AddWave"},
			_jsii_.MemberMethod{JsiiMethod: "buildPipeline", GoMethod: "BuildPipeline"},
			_jsii_.MemberProperty{JsiiProperty: "cloudAssemblyFileSet", GoGetter: "CloudAssemblyFileSet"},
			_jsii_.MemberMethod{JsiiMethod: "doBuildPipeline", GoMethod: "DoBuildPipeline"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "synth", GoGetter: "Synth"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "waves", GoGetter: "Waves"},
			_jsii_.MemberProperty{JsiiProperty: "workflowFile", GoGetter: "WorkflowFile"},
			_jsii_.MemberProperty{JsiiProperty: "workflowName", GoGetter: "WorkflowName"},
			_jsii_.MemberProperty{JsiiProperty: "workflowPath", GoGetter: "WorkflowPath"},
		},
		func() interface{} {
			j := jsiiProxy_GitHubWorkflow{}
			_jsii_.InitJsiiProxy(&j.Type__pipelinesPipelineBase)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.GitHubWorkflowProps",
		reflect.TypeOf((*GitHubWorkflowProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.GollumOptions",
		reflect.TypeOf((*GollumOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.IssueCommentOptions",
		reflect.TypeOf((*IssueCommentOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.IssuesOptions",
		reflect.TypeOf((*IssuesOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.Job",
		reflect.TypeOf((*Job)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.JobDefaults",
		reflect.TypeOf((*JobDefaults)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.JobMatrix",
		reflect.TypeOf((*JobMatrix)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"cdk-pipelines-github.JobPermission",
		reflect.TypeOf((*JobPermission)(nil)).Elem(),
		map[string]interface{}{
			"READ": JobPermission_READ,
			"WRITE": JobPermission_WRITE,
			"NONE": JobPermission_NONE,
		},
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.JobPermissions",
		reflect.TypeOf((*JobPermissions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.JobSettings",
		reflect.TypeOf((*JobSettings)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.JobStep",
		reflect.TypeOf((*JobStep)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.JobStepOutput",
		reflect.TypeOf((*JobStepOutput)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.JobStrategy",
		reflect.TypeOf((*JobStrategy)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"cdk-pipelines-github.JsonPatch",
		reflect.TypeOf((*JsonPatch)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_JsonPatch{}
		},
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.LabelOptions",
		reflect.TypeOf((*LabelOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.MilestoneOptions",
		reflect.TypeOf((*MilestoneOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.OpenIdConnectProviderProps",
		reflect.TypeOf((*OpenIdConnectProviderProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.PageBuildOptions",
		reflect.TypeOf((*PageBuildOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.ProjectCardOptions",
		reflect.TypeOf((*ProjectCardOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.ProjectColumnOptions",
		reflect.TypeOf((*ProjectColumnOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.ProjectOptions",
		reflect.TypeOf((*ProjectOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.PublicOptions",
		reflect.TypeOf((*PublicOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.PullRequestOptions",
		reflect.TypeOf((*PullRequestOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.PullRequestReviewCommentOptions",
		reflect.TypeOf((*PullRequestReviewCommentOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.PullRequestReviewOptions",
		reflect.TypeOf((*PullRequestReviewOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.PullRequestTargetOptions",
		reflect.TypeOf((*PullRequestTargetOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.PushOptions",
		reflect.TypeOf((*PushOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.RegistryPackageOptions",
		reflect.TypeOf((*RegistryPackageOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.ReleaseOptions",
		reflect.TypeOf((*ReleaseOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.RepositoryDispatchOptions",
		reflect.TypeOf((*RepositoryDispatchOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.RunSettings",
		reflect.TypeOf((*RunSettings)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"cdk-pipelines-github.Runner",
		reflect.TypeOf((*Runner)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "runsOn", GoGetter: "RunsOn"},
		},
		func() interface{} {
			return &jsiiProxy_Runner{}
		},
	)
	_jsii_.RegisterEnum(
		"cdk-pipelines-github.StackCapabilities",
		reflect.TypeOf((*StackCapabilities)(nil)).Elem(),
		map[string]interface{}{
			"IAM": StackCapabilities_IAM,
			"NAMED_IAM": StackCapabilities_NAMED_IAM,
			"AUTO_EXPAND": StackCapabilities_AUTO_EXPAND,
		},
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.StatusOptions",
		reflect.TypeOf((*StatusOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.WatchOptions",
		reflect.TypeOf((*WatchOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.WorkflowDispatchOptions",
		reflect.TypeOf((*WorkflowDispatchOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.WorkflowRunOptions",
		reflect.TypeOf((*WorkflowRunOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.WorkflowTriggers",
		reflect.TypeOf((*WorkflowTriggers)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"cdk-pipelines-github.YamlFile",
		reflect.TypeOf((*YamlFile)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "commentAtTop", GoGetter: "CommentAtTop"},
			_jsii_.MemberMethod{JsiiMethod: "patch", GoMethod: "Patch"},
			_jsii_.MemberMethod{JsiiMethod: "toYaml", GoMethod: "ToYaml"},
			_jsii_.MemberMethod{JsiiMethod: "update", GoMethod: "Update"},
			_jsii_.MemberMethod{JsiiMethod: "writeFile", GoMethod: "WriteFile"},
		},
		func() interface{} {
			return &jsiiProxy_YamlFile{}
		},
	)
	_jsii_.RegisterStruct(
		"cdk-pipelines-github.YamlFileOptions",
		reflect.TypeOf((*YamlFileOptions)(nil)).Elem(),
	)
}
