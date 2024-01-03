package awscodedeploy

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/v2/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudwatch"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscodedeploy/internal"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
)

// Example:
//   var application lambdaApplication
//   var alias alias
//   config := codedeploy.NewLambdaDeploymentConfig(this, jsii.String("CustomConfig"), &LambdaDeploymentConfigProps{
//   	TrafficRouting: codedeploy.NewTimeBasedCanaryTrafficRouting(&TimeBasedCanaryTrafficRoutingProps{
//   		Interval: awscdk.Duration_Minutes(jsii.Number(15)),
//   		Percentage: jsii.Number(5),
//   	}),
//   })
//   deploymentGroup := codedeploy.NewLambdaDeploymentGroup(this, jsii.String("BlueGreenDeployment"), &LambdaDeploymentGroupProps{
//   	Application: Application,
//   	Alias: Alias,
//   	DeploymentConfig: config,
//   })
//
type LambdaDeploymentGroup interface {
	awscdk.Resource
	ILambdaDeploymentGroup
	// The reference to the CodeDeploy Lambda Application that this Deployment Group belongs to.
	Application() ILambdaApplication
	// The Deployment Configuration this Group uses.
	DeploymentConfig() ILambdaDeploymentConfig
	// The ARN of the Deployment Group.
	DeploymentGroupArn() *string
	// The name of the Deployment Group.
	DeploymentGroupName() *string
	// The environment this resource belongs to.
	//
	// For resources that are created and managed by the CDK
	// (generally, those created by creating new class instances like Role, Bucket, etc.),
	// this is always the same as the environment of the stack they belong to;
	// however, for imported resources
	// (those obtained from static methods like fromRoleArn, fromBucketName, etc.),
	// that might be different than the stack they were imported into.
	Env() *awscdk.ResourceEnvironment
	// The tree node.
	Node() constructs.Node
	// Returns a string-encoded token that resolves to the physical name that should be passed to the CloudFormation resource.
	//
	// This value will resolve to one of the following:
	// - a concrete value (e.g. `"my-awesome-bucket"`)
	// - `undefined`, when a name should be generated by CloudFormation
	// - a concrete name generated automatically during synthesis, in
	//   cross-environment scenarios.
	PhysicalName() *string
	// The service Role of this Deployment Group.
	Role() awsiam.IRole
	// The stack in which this resource is defined.
	Stack() awscdk.Stack
	// Associates an additional alarm with this Deployment Group.
	AddAlarm(alarm awscloudwatch.IAlarm)
	// Associate a function to run after deployment completes.
	AddPostHook(postHook awslambda.IFunction)
	// Associate a function to run before deployment begins.
	AddPreHook(preHook awslambda.IFunction)
	// Apply the given removal policy to this resource.
	//
	// The Removal Policy controls what happens to this resource when it stops
	// being managed by CloudFormation, either because you've removed it from the
	// CDK application or because you've made a change that requires the resource
	// to be replaced.
	//
	// The resource can be deleted (`RemovalPolicy.DESTROY`), or left in your AWS
	// account for data recovery and cleanup later (`RemovalPolicy.RETAIN`).
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	// Returns an environment-sensitive token that should be used for the resource's "ARN" attribute (e.g. `bucket.bucketArn`).
	//
	// Normally, this token will resolve to `arnAttr`, but if the resource is
	// referenced across environments, `arnComponents` will be used to synthesize
	// a concrete ARN with the resource's physical name. Make sure to reference
	// `this.physicalName` in `arnComponents`.
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	// Returns an environment-sensitive token that should be used for the resource's "name" attribute (e.g. `bucket.bucketName`).
	//
	// Normally, this token will resolve to `nameAttr`, but if the resource is
	// referenced across environments, it will be resolved to `this.physicalName`,
	// which will be a concrete name.
	GetResourceNameAttribute(nameAttr *string) *string
	// Grant a principal permission to codedeploy:PutLifecycleEventHookExecutionStatus on this deployment group resource.
	GrantPutLifecycleEventHookExecutionStatus(grantee awsiam.IGrantable) awsiam.Grant
	// Returns a string representation of this construct.
	ToString() *string
}

// The jsii proxy struct for LambdaDeploymentGroup
type jsiiProxy_LambdaDeploymentGroup struct {
	internal.Type__awscdkResource
	jsiiProxy_ILambdaDeploymentGroup
}

func (j *jsiiProxy_LambdaDeploymentGroup) Application() ILambdaApplication {
	var returns ILambdaApplication
	_jsii_.Get(
		j,
		"application",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LambdaDeploymentGroup) DeploymentConfig() ILambdaDeploymentConfig {
	var returns ILambdaDeploymentConfig
	_jsii_.Get(
		j,
		"deploymentConfig",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LambdaDeploymentGroup) DeploymentGroupArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"deploymentGroupArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LambdaDeploymentGroup) DeploymentGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"deploymentGroupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LambdaDeploymentGroup) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LambdaDeploymentGroup) Node() constructs.Node {
	var returns constructs.Node
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LambdaDeploymentGroup) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LambdaDeploymentGroup) Role() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"role",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LambdaDeploymentGroup) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


func NewLambdaDeploymentGroup(scope constructs.Construct, id *string, props *LambdaDeploymentGroupProps) LambdaDeploymentGroup {
	_init_.Initialize()

	if err := validateNewLambdaDeploymentGroupParameters(scope, id, props); err != nil {
		panic(err)
	}
	j := jsiiProxy_LambdaDeploymentGroup{}

	_jsii_.Create(
		"aws-cdk-lib.aws_codedeploy.LambdaDeploymentGroup",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

func NewLambdaDeploymentGroup_Override(l LambdaDeploymentGroup, scope constructs.Construct, id *string, props *LambdaDeploymentGroupProps) {
	_init_.Initialize()

	_jsii_.Create(
		"aws-cdk-lib.aws_codedeploy.LambdaDeploymentGroup",
		[]interface{}{scope, id, props},
		l,
	)
}

// Import an Lambda Deployment Group defined either outside the CDK app, or in a different AWS region.
//
// Account and region for the DeploymentGroup are taken from the application.
//
// Returns: a Construct representing a reference to an existing Deployment Group.
func LambdaDeploymentGroup_FromLambdaDeploymentGroupAttributes(scope constructs.Construct, id *string, attrs *LambdaDeploymentGroupAttributes) ILambdaDeploymentGroup {
	_init_.Initialize()

	if err := validateLambdaDeploymentGroup_FromLambdaDeploymentGroupAttributesParameters(scope, id, attrs); err != nil {
		panic(err)
	}
	var returns ILambdaDeploymentGroup

	_jsii_.StaticInvoke(
		"aws-cdk-lib.aws_codedeploy.LambdaDeploymentGroup",
		"fromLambdaDeploymentGroupAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Checks if `x` is a construct.
//
// Use this method instead of `instanceof` to properly detect `Construct`
// instances, even when the construct library is symlinked.
//
// Explanation: in JavaScript, multiple copies of the `constructs` library on
// disk are seen as independent, completely different libraries. As a
// consequence, the class `Construct` in each copy of the `constructs` library
// is seen as a different class, and an instance of one class will not test as
// `instanceof` the other class. `npm install` will not create installations
// like this, but users may manually symlink construct libraries together or
// use a monorepo tool: in those cases, multiple copies of the `constructs`
// library can be accidentally installed, and `instanceof` will behave
// unpredictably. It is safest to avoid using `instanceof`, and using
// this type-testing method instead.
//
// Returns: true if `x` is an object created from a class which extends `Construct`.
func LambdaDeploymentGroup_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	if err := validateLambdaDeploymentGroup_IsConstructParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"aws-cdk-lib.aws_codedeploy.LambdaDeploymentGroup",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Returns true if the construct was created by CDK, and false otherwise.
func LambdaDeploymentGroup_IsOwnedResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	if err := validateLambdaDeploymentGroup_IsOwnedResourceParameters(construct); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"aws-cdk-lib.aws_codedeploy.LambdaDeploymentGroup",
		"isOwnedResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
func LambdaDeploymentGroup_IsResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	if err := validateLambdaDeploymentGroup_IsResourceParameters(construct); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"aws-cdk-lib.aws_codedeploy.LambdaDeploymentGroup",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

func (l *jsiiProxy_LambdaDeploymentGroup) AddAlarm(alarm awscloudwatch.IAlarm) {
	if err := l.validateAddAlarmParameters(alarm); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		l,
		"addAlarm",
		[]interface{}{alarm},
	)
}

func (l *jsiiProxy_LambdaDeploymentGroup) AddPostHook(postHook awslambda.IFunction) {
	if err := l.validateAddPostHookParameters(postHook); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		l,
		"addPostHook",
		[]interface{}{postHook},
	)
}

func (l *jsiiProxy_LambdaDeploymentGroup) AddPreHook(preHook awslambda.IFunction) {
	if err := l.validateAddPreHookParameters(preHook); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		l,
		"addPreHook",
		[]interface{}{preHook},
	)
}

func (l *jsiiProxy_LambdaDeploymentGroup) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	if err := l.validateApplyRemovalPolicyParameters(policy); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		l,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

func (l *jsiiProxy_LambdaDeploymentGroup) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		l,
		"generatePhysicalName",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (l *jsiiProxy_LambdaDeploymentGroup) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	if err := l.validateGetResourceArnAttributeParameters(arnAttr, arnComponents); err != nil {
		panic(err)
	}
	var returns *string

	_jsii_.Invoke(
		l,
		"getResourceArnAttribute",
		[]interface{}{arnAttr, arnComponents},
		&returns,
	)

	return returns
}

func (l *jsiiProxy_LambdaDeploymentGroup) GetResourceNameAttribute(nameAttr *string) *string {
	if err := l.validateGetResourceNameAttributeParameters(nameAttr); err != nil {
		panic(err)
	}
	var returns *string

	_jsii_.Invoke(
		l,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

func (l *jsiiProxy_LambdaDeploymentGroup) GrantPutLifecycleEventHookExecutionStatus(grantee awsiam.IGrantable) awsiam.Grant {
	if err := l.validateGrantPutLifecycleEventHookExecutionStatusParameters(grantee); err != nil {
		panic(err)
	}
	var returns awsiam.Grant

	_jsii_.Invoke(
		l,
		"grantPutLifecycleEventHookExecutionStatus",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

func (l *jsiiProxy_LambdaDeploymentGroup) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		l,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

