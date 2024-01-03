package awsecs

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/v2/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
)

// The details of a task definition run on a Fargate cluster.
//
// Example:
//   var vpc vpc
//
//
//   cluster := ecs.NewCluster(this, jsii.String("FargateCPCluster"), &ClusterProps{
//   	Vpc: Vpc,
//   	EnableFargateCapacityProviders: jsii.Boolean(true),
//   })
//
//   taskDefinition := ecs.NewFargateTaskDefinition(this, jsii.String("TaskDef"))
//
//   taskDefinition.AddContainer(jsii.String("web"), &ContainerDefinitionOptions{
//   	Image: ecs.ContainerImage_FromRegistry(jsii.String("amazon/amazon-ecs-sample")),
//   })
//
//   ecs.NewFargateService(this, jsii.String("FargateService"), &FargateServiceProps{
//   	Cluster: Cluster,
//   	TaskDefinition: TaskDefinition,
//   	CapacityProviderStrategies: []capacityProviderStrategy{
//   		&capacityProviderStrategy{
//   			CapacityProvider: jsii.String("FARGATE_SPOT"),
//   			Weight: jsii.Number(2),
//   		},
//   		&capacityProviderStrategy{
//   			CapacityProvider: jsii.String("FARGATE"),
//   			Weight: jsii.Number(1),
//   		},
//   	},
//   })
//
type FargateTaskDefinition interface {
	TaskDefinition
	IFargateTaskDefinition
	// The task launch type compatibility requirement.
	Compatibility() Compatibility
	// The container definitions.
	Containers() *[]ContainerDefinition
	// Default container for this task.
	//
	// Load balancers will send traffic to this container. The first
	// essential container that is added to this task will become the default
	// container.
	DefaultContainer() ContainerDefinition
	SetDefaultContainer(val ContainerDefinition)
	// The environment this resource belongs to.
	//
	// For resources that are created and managed by the CDK
	// (generally, those created by creating new class instances like Role, Bucket, etc.),
	// this is always the same as the environment of the stack they belong to;
	// however, for imported resources
	// (those obtained from static methods like fromRoleArn, fromBucketName, etc.),
	// that might be different than the stack they were imported into.
	Env() *awscdk.ResourceEnvironment
	// The amount (in GiB) of ephemeral storage to be allocated to the task.
	EphemeralStorageGiB() *float64
	// Execution role for this task definition.
	ExecutionRole() awsiam.IRole
	// The name of a family that this task definition is registered to.
	//
	// A family groups multiple versions of a task definition.
	Family() *string
	// Public getter method to access list of inference accelerators attached to the instance.
	InferenceAccelerators() *[]*InferenceAccelerator
	// Return true if the task definition can be run on an EC2 cluster.
	IsEc2Compatible() *bool
	// Return true if the task definition can be run on a ECS anywhere cluster.
	IsExternalCompatible() *bool
	// Return true if the task definition can be run on a Fargate cluster.
	IsFargateCompatible() *bool
	// The Docker networking mode to use for the containers in the task.
	//
	// Fargate tasks require the awsvpc network mode.
	NetworkMode() NetworkMode
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
	// Whether this task definition has at least a container that references a specific JSON field of a secret stored in Secrets Manager.
	ReferencesSecretJsonField() *bool
	// The stack in which this resource is defined.
	Stack() awscdk.Stack
	// The full Amazon Resource Name (ARN) of the task definition.
	TaskDefinitionArn() *string
	// The name of the IAM role that grants containers in the task permission to call AWS APIs on your behalf.
	TaskRole() awsiam.IRole
	// Adds a new container to the task definition.
	AddContainer(id *string, props *ContainerDefinitionOptions) ContainerDefinition
	// Adds the specified extension to the task definition.
	//
	// Extension can be used to apply a packaged modification to
	// a task definition.
	AddExtension(extension ITaskDefinitionExtension)
	// Adds a firelens log router to the task definition.
	AddFirelensLogRouter(id *string, props *FirelensLogRouterDefinitionOptions) FirelensLogRouter
	// Adds an inference accelerator to the task definition.
	AddInferenceAccelerator(inferenceAccelerator *InferenceAccelerator)
	// Adds the specified placement constraint to the task definition.
	AddPlacementConstraint(constraint PlacementConstraint)
	// Adds a policy statement to the task execution IAM role.
	AddToExecutionRolePolicy(statement awsiam.PolicyStatement)
	// Adds a policy statement to the task IAM role.
	AddToTaskRolePolicy(statement awsiam.PolicyStatement)
	// Adds a volume to the task definition.
	AddVolume(volume *Volume)
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
	// Returns the container that match the provided containerName.
	FindContainer(containerName *string) ContainerDefinition
	// Determine the existing port mapping for the provided name.
	//
	// Returns: PortMapping for the provided name, if it exists.
	FindPortMappingByName(name *string) *PortMapping
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
	// Grants permissions to run this task definition.
	//
	// This will grant the following permissions:
	//
	//   - ecs:RunTask
	// - iam:PassRole.
	GrantRun(grantee awsiam.IGrantable) awsiam.Grant
	// Creates the task execution IAM role if it doesn't already exist.
	ObtainExecutionRole() awsiam.IRole
	// Returns a string representation of this construct.
	ToString() *string
}

// The jsii proxy struct for FargateTaskDefinition
type jsiiProxy_FargateTaskDefinition struct {
	jsiiProxy_TaskDefinition
	jsiiProxy_IFargateTaskDefinition
}

func (j *jsiiProxy_FargateTaskDefinition) Compatibility() Compatibility {
	var returns Compatibility
	_jsii_.Get(
		j,
		"compatibility",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) Containers() *[]ContainerDefinition {
	var returns *[]ContainerDefinition
	_jsii_.Get(
		j,
		"containers",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) DefaultContainer() ContainerDefinition {
	var returns ContainerDefinition
	_jsii_.Get(
		j,
		"defaultContainer",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) EphemeralStorageGiB() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"ephemeralStorageGiB",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) ExecutionRole() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"executionRole",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) Family() *string {
	var returns *string
	_jsii_.Get(
		j,
		"family",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) InferenceAccelerators() *[]*InferenceAccelerator {
	var returns *[]*InferenceAccelerator
	_jsii_.Get(
		j,
		"inferenceAccelerators",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) IsEc2Compatible() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isEc2Compatible",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) IsExternalCompatible() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isExternalCompatible",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) IsFargateCompatible() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isFargateCompatible",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) NetworkMode() NetworkMode {
	var returns NetworkMode
	_jsii_.Get(
		j,
		"networkMode",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) Node() constructs.Node {
	var returns constructs.Node
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) ReferencesSecretJsonField() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"referencesSecretJsonField",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) TaskDefinitionArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"taskDefinitionArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) TaskRole() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"taskRole",
		&returns,
	)
	return returns
}


// Constructs a new instance of the FargateTaskDefinition class.
func NewFargateTaskDefinition(scope constructs.Construct, id *string, props *FargateTaskDefinitionProps) FargateTaskDefinition {
	_init_.Initialize()

	if err := validateNewFargateTaskDefinitionParameters(scope, id, props); err != nil {
		panic(err)
	}
	j := jsiiProxy_FargateTaskDefinition{}

	_jsii_.Create(
		"aws-cdk-lib.aws_ecs.FargateTaskDefinition",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Constructs a new instance of the FargateTaskDefinition class.
func NewFargateTaskDefinition_Override(f FargateTaskDefinition, scope constructs.Construct, id *string, props *FargateTaskDefinitionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"aws-cdk-lib.aws_ecs.FargateTaskDefinition",
		[]interface{}{scope, id, props},
		f,
	)
}

func (j *jsiiProxy_FargateTaskDefinition)SetDefaultContainer(val ContainerDefinition) {
	_jsii_.Set(
		j,
		"defaultContainer",
		val,
	)
}

// Imports a task definition from the specified task definition ARN.
func FargateTaskDefinition_FromFargateTaskDefinitionArn(scope constructs.Construct, id *string, fargateTaskDefinitionArn *string) IFargateTaskDefinition {
	_init_.Initialize()

	if err := validateFargateTaskDefinition_FromFargateTaskDefinitionArnParameters(scope, id, fargateTaskDefinitionArn); err != nil {
		panic(err)
	}
	var returns IFargateTaskDefinition

	_jsii_.StaticInvoke(
		"aws-cdk-lib.aws_ecs.FargateTaskDefinition",
		"fromFargateTaskDefinitionArn",
		[]interface{}{scope, id, fargateTaskDefinitionArn},
		&returns,
	)

	return returns
}

// Import an existing Fargate task definition from its attributes.
func FargateTaskDefinition_FromFargateTaskDefinitionAttributes(scope constructs.Construct, id *string, attrs *FargateTaskDefinitionAttributes) IFargateTaskDefinition {
	_init_.Initialize()

	if err := validateFargateTaskDefinition_FromFargateTaskDefinitionAttributesParameters(scope, id, attrs); err != nil {
		panic(err)
	}
	var returns IFargateTaskDefinition

	_jsii_.StaticInvoke(
		"aws-cdk-lib.aws_ecs.FargateTaskDefinition",
		"fromFargateTaskDefinitionAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Imports a task definition from the specified task definition ARN.
//
// The task will have a compatibility of EC2+Fargate.
func FargateTaskDefinition_FromTaskDefinitionArn(scope constructs.Construct, id *string, taskDefinitionArn *string) ITaskDefinition {
	_init_.Initialize()

	if err := validateFargateTaskDefinition_FromTaskDefinitionArnParameters(scope, id, taskDefinitionArn); err != nil {
		panic(err)
	}
	var returns ITaskDefinition

	_jsii_.StaticInvoke(
		"aws-cdk-lib.aws_ecs.FargateTaskDefinition",
		"fromTaskDefinitionArn",
		[]interface{}{scope, id, taskDefinitionArn},
		&returns,
	)

	return returns
}

// Create a task definition from a task definition reference.
func FargateTaskDefinition_FromTaskDefinitionAttributes(scope constructs.Construct, id *string, attrs *TaskDefinitionAttributes) ITaskDefinition {
	_init_.Initialize()

	if err := validateFargateTaskDefinition_FromTaskDefinitionAttributesParameters(scope, id, attrs); err != nil {
		panic(err)
	}
	var returns ITaskDefinition

	_jsii_.StaticInvoke(
		"aws-cdk-lib.aws_ecs.FargateTaskDefinition",
		"fromTaskDefinitionAttributes",
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
func FargateTaskDefinition_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	if err := validateFargateTaskDefinition_IsConstructParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"aws-cdk-lib.aws_ecs.FargateTaskDefinition",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Returns true if the construct was created by CDK, and false otherwise.
func FargateTaskDefinition_IsOwnedResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	if err := validateFargateTaskDefinition_IsOwnedResourceParameters(construct); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"aws-cdk-lib.aws_ecs.FargateTaskDefinition",
		"isOwnedResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
func FargateTaskDefinition_IsResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	if err := validateFargateTaskDefinition_IsResourceParameters(construct); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"aws-cdk-lib.aws_ecs.FargateTaskDefinition",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

func (f *jsiiProxy_FargateTaskDefinition) AddContainer(id *string, props *ContainerDefinitionOptions) ContainerDefinition {
	if err := f.validateAddContainerParameters(id, props); err != nil {
		panic(err)
	}
	var returns ContainerDefinition

	_jsii_.Invoke(
		f,
		"addContainer",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

func (f *jsiiProxy_FargateTaskDefinition) AddExtension(extension ITaskDefinitionExtension) {
	if err := f.validateAddExtensionParameters(extension); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		f,
		"addExtension",
		[]interface{}{extension},
	)
}

func (f *jsiiProxy_FargateTaskDefinition) AddFirelensLogRouter(id *string, props *FirelensLogRouterDefinitionOptions) FirelensLogRouter {
	if err := f.validateAddFirelensLogRouterParameters(id, props); err != nil {
		panic(err)
	}
	var returns FirelensLogRouter

	_jsii_.Invoke(
		f,
		"addFirelensLogRouter",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

func (f *jsiiProxy_FargateTaskDefinition) AddInferenceAccelerator(inferenceAccelerator *InferenceAccelerator) {
	if err := f.validateAddInferenceAcceleratorParameters(inferenceAccelerator); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		f,
		"addInferenceAccelerator",
		[]interface{}{inferenceAccelerator},
	)
}

func (f *jsiiProxy_FargateTaskDefinition) AddPlacementConstraint(constraint PlacementConstraint) {
	if err := f.validateAddPlacementConstraintParameters(constraint); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		f,
		"addPlacementConstraint",
		[]interface{}{constraint},
	)
}

func (f *jsiiProxy_FargateTaskDefinition) AddToExecutionRolePolicy(statement awsiam.PolicyStatement) {
	if err := f.validateAddToExecutionRolePolicyParameters(statement); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		f,
		"addToExecutionRolePolicy",
		[]interface{}{statement},
	)
}

func (f *jsiiProxy_FargateTaskDefinition) AddToTaskRolePolicy(statement awsiam.PolicyStatement) {
	if err := f.validateAddToTaskRolePolicyParameters(statement); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		f,
		"addToTaskRolePolicy",
		[]interface{}{statement},
	)
}

func (f *jsiiProxy_FargateTaskDefinition) AddVolume(volume *Volume) {
	if err := f.validateAddVolumeParameters(volume); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		f,
		"addVolume",
		[]interface{}{volume},
	)
}

func (f *jsiiProxy_FargateTaskDefinition) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	if err := f.validateApplyRemovalPolicyParameters(policy); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		f,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

func (f *jsiiProxy_FargateTaskDefinition) FindContainer(containerName *string) ContainerDefinition {
	if err := f.validateFindContainerParameters(containerName); err != nil {
		panic(err)
	}
	var returns ContainerDefinition

	_jsii_.Invoke(
		f,
		"findContainer",
		[]interface{}{containerName},
		&returns,
	)

	return returns
}

func (f *jsiiProxy_FargateTaskDefinition) FindPortMappingByName(name *string) *PortMapping {
	if err := f.validateFindPortMappingByNameParameters(name); err != nil {
		panic(err)
	}
	var returns *PortMapping

	_jsii_.Invoke(
		f,
		"findPortMappingByName",
		[]interface{}{name},
		&returns,
	)

	return returns
}

func (f *jsiiProxy_FargateTaskDefinition) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		f,
		"generatePhysicalName",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (f *jsiiProxy_FargateTaskDefinition) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	if err := f.validateGetResourceArnAttributeParameters(arnAttr, arnComponents); err != nil {
		panic(err)
	}
	var returns *string

	_jsii_.Invoke(
		f,
		"getResourceArnAttribute",
		[]interface{}{arnAttr, arnComponents},
		&returns,
	)

	return returns
}

func (f *jsiiProxy_FargateTaskDefinition) GetResourceNameAttribute(nameAttr *string) *string {
	if err := f.validateGetResourceNameAttributeParameters(nameAttr); err != nil {
		panic(err)
	}
	var returns *string

	_jsii_.Invoke(
		f,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

func (f *jsiiProxy_FargateTaskDefinition) GrantRun(grantee awsiam.IGrantable) awsiam.Grant {
	if err := f.validateGrantRunParameters(grantee); err != nil {
		panic(err)
	}
	var returns awsiam.Grant

	_jsii_.Invoke(
		f,
		"grantRun",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

func (f *jsiiProxy_FargateTaskDefinition) ObtainExecutionRole() awsiam.IRole {
	var returns awsiam.IRole

	_jsii_.Invoke(
		f,
		"obtainExecutionRole",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (f *jsiiProxy_FargateTaskDefinition) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		f,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

