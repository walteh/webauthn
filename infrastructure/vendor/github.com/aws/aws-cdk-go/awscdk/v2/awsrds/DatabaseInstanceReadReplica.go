package awsrds

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/v2/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudwatch"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssecretsmanager"
	"github.com/aws/constructs-go/constructs/v10"
)

// A read replica database instance.
//
// Example:
//   var vpc vpc
//
//   var sourceInstance databaseInstance
//
//   rds.NewDatabaseInstanceFromSnapshot(this, jsii.String("Instance"), &DatabaseInstanceFromSnapshotProps{
//   	SnapshotIdentifier: jsii.String("my-snapshot"),
//   	Engine: rds.DatabaseInstanceEngine_Postgres(&PostgresInstanceEngineProps{
//   		Version: rds.PostgresEngineVersion_VER_15_2(),
//   	}),
//   	// optional, defaults to m5.large
//   	InstanceType: ec2.InstanceType_Of(ec2.InstanceClass_BURSTABLE2, ec2.InstanceSize_LARGE),
//   	Vpc: Vpc,
//   })
//   rds.NewDatabaseInstanceReadReplica(this, jsii.String("ReadReplica"), &DatabaseInstanceReadReplicaProps{
//   	SourceDatabaseInstance: sourceInstance,
//   	InstanceType: ec2.InstanceType_*Of(ec2.InstanceClass_BURSTABLE2, ec2.InstanceSize_LARGE),
//   	Vpc: Vpc,
//   })
//
type DatabaseInstanceReadReplica interface {
	DatabaseInstanceBase
	IDatabaseInstance
	// Access to network connections.
	Connections() awsec2.Connections
	// The instance endpoint address.
	DbInstanceEndpointAddress() *string
	// The instance endpoint port.
	DbInstanceEndpointPort() *string
	EnableIamAuthentication() *bool
	SetEnableIamAuthentication(val *bool)
	// The engine of this database Instance.
	//
	// May be not known for imported Instances if it wasn't provided explicitly,
	// or for read replicas.
	Engine() IInstanceEngine
	// The environment this resource belongs to.
	//
	// For resources that are created and managed by the CDK
	// (generally, those created by creating new class instances like Role, Bucket, etc.),
	// this is always the same as the environment of the stack they belong to;
	// however, for imported resources
	// (those obtained from static methods like fromRoleArn, fromBucketName, etc.),
	// that might be different than the stack they were imported into.
	Env() *awscdk.ResourceEnvironment
	// The instance arn.
	InstanceArn() *string
	// The instance endpoint.
	InstanceEndpoint() Endpoint
	// The instance identifier.
	InstanceIdentifier() *string
	// The AWS Region-unique, immutable identifier for the DB instance.
	//
	// This identifier is found in AWS CloudTrail log entries whenever the AWS KMS key for the DB instance is accessed.
	InstanceResourceId() *string
	InstanceType() awsec2.InstanceType
	NewCfnProps() *CfnDBInstanceProps
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
	// The stack in which this resource is defined.
	Stack() awscdk.Stack
	// The VPC where this database instance is deployed.
	Vpc() awsec2.IVpc
	VpcPlacement() *awsec2.SubnetSelection
	// Add a new db proxy to this instance.
	AddProxy(id *string, options *DatabaseProxyOptions) DatabaseProxy
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
	// Renders the secret attachment target specifications.
	AsSecretAttachmentTarget() *awssecretsmanager.SecretAttachmentTargetProps
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
	// Grant the given identity connection access to the database.
	GrantConnect(grantee awsiam.IGrantable, dbUser *string) awsiam.Grant
	// Return the given named metric for this DBInstance.
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The percentage of CPU utilization.
	//
	// Average over 5 minutes.
	MetricCPUUtilization(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of database connections in use.
	//
	// Average over 5 minutes.
	MetricDatabaseConnections(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The amount of available random access memory.
	//
	// Average over 5 minutes.
	MetricFreeableMemory(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The amount of available storage space.
	//
	// Average over 5 minutes.
	MetricFreeStorageSpace(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The average number of disk write I/O operations per second.
	//
	// Average over 5 minutes.
	MetricReadIOPS(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The average number of disk read I/O operations per second.
	//
	// Average over 5 minutes.
	MetricWriteIOPS(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// Defines a CloudWatch event rule which triggers for instance events.
	//
	// Use
	// `rule.addEventPattern(pattern)` to specify a filter.
	OnEvent(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	SetLogRetention()
	// Returns a string representation of this construct.
	ToString() *string
}

// The jsii proxy struct for DatabaseInstanceReadReplica
type jsiiProxy_DatabaseInstanceReadReplica struct {
	jsiiProxy_DatabaseInstanceBase
	jsiiProxy_IDatabaseInstance
}

func (j *jsiiProxy_DatabaseInstanceReadReplica) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DatabaseInstanceReadReplica) DbInstanceEndpointAddress() *string {
	var returns *string
	_jsii_.Get(
		j,
		"dbInstanceEndpointAddress",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DatabaseInstanceReadReplica) DbInstanceEndpointPort() *string {
	var returns *string
	_jsii_.Get(
		j,
		"dbInstanceEndpointPort",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DatabaseInstanceReadReplica) EnableIamAuthentication() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"enableIamAuthentication",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DatabaseInstanceReadReplica) Engine() IInstanceEngine {
	var returns IInstanceEngine
	_jsii_.Get(
		j,
		"engine",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DatabaseInstanceReadReplica) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DatabaseInstanceReadReplica) InstanceArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"instanceArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DatabaseInstanceReadReplica) InstanceEndpoint() Endpoint {
	var returns Endpoint
	_jsii_.Get(
		j,
		"instanceEndpoint",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DatabaseInstanceReadReplica) InstanceIdentifier() *string {
	var returns *string
	_jsii_.Get(
		j,
		"instanceIdentifier",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DatabaseInstanceReadReplica) InstanceResourceId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"instanceResourceId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DatabaseInstanceReadReplica) InstanceType() awsec2.InstanceType {
	var returns awsec2.InstanceType
	_jsii_.Get(
		j,
		"instanceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DatabaseInstanceReadReplica) NewCfnProps() *CfnDBInstanceProps {
	var returns *CfnDBInstanceProps
	_jsii_.Get(
		j,
		"newCfnProps",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DatabaseInstanceReadReplica) Node() constructs.Node {
	var returns constructs.Node
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DatabaseInstanceReadReplica) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DatabaseInstanceReadReplica) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DatabaseInstanceReadReplica) Vpc() awsec2.IVpc {
	var returns awsec2.IVpc
	_jsii_.Get(
		j,
		"vpc",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DatabaseInstanceReadReplica) VpcPlacement() *awsec2.SubnetSelection {
	var returns *awsec2.SubnetSelection
	_jsii_.Get(
		j,
		"vpcPlacement",
		&returns,
	)
	return returns
}


func NewDatabaseInstanceReadReplica(scope constructs.Construct, id *string, props *DatabaseInstanceReadReplicaProps) DatabaseInstanceReadReplica {
	_init_.Initialize()

	if err := validateNewDatabaseInstanceReadReplicaParameters(scope, id, props); err != nil {
		panic(err)
	}
	j := jsiiProxy_DatabaseInstanceReadReplica{}

	_jsii_.Create(
		"aws-cdk-lib.aws_rds.DatabaseInstanceReadReplica",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

func NewDatabaseInstanceReadReplica_Override(d DatabaseInstanceReadReplica, scope constructs.Construct, id *string, props *DatabaseInstanceReadReplicaProps) {
	_init_.Initialize()

	_jsii_.Create(
		"aws-cdk-lib.aws_rds.DatabaseInstanceReadReplica",
		[]interface{}{scope, id, props},
		d,
	)
}

func (j *jsiiProxy_DatabaseInstanceReadReplica)SetEnableIamAuthentication(val *bool) {
	_jsii_.Set(
		j,
		"enableIamAuthentication",
		val,
	)
}

// Import an existing database instance.
func DatabaseInstanceReadReplica_FromDatabaseInstanceAttributes(scope constructs.Construct, id *string, attrs *DatabaseInstanceAttributes) IDatabaseInstance {
	_init_.Initialize()

	if err := validateDatabaseInstanceReadReplica_FromDatabaseInstanceAttributesParameters(scope, id, attrs); err != nil {
		panic(err)
	}
	var returns IDatabaseInstance

	_jsii_.StaticInvoke(
		"aws-cdk-lib.aws_rds.DatabaseInstanceReadReplica",
		"fromDatabaseInstanceAttributes",
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
func DatabaseInstanceReadReplica_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	if err := validateDatabaseInstanceReadReplica_IsConstructParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"aws-cdk-lib.aws_rds.DatabaseInstanceReadReplica",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Returns true if the construct was created by CDK, and false otherwise.
func DatabaseInstanceReadReplica_IsOwnedResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	if err := validateDatabaseInstanceReadReplica_IsOwnedResourceParameters(construct); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"aws-cdk-lib.aws_rds.DatabaseInstanceReadReplica",
		"isOwnedResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
func DatabaseInstanceReadReplica_IsResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	if err := validateDatabaseInstanceReadReplica_IsResourceParameters(construct); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"aws-cdk-lib.aws_rds.DatabaseInstanceReadReplica",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DatabaseInstanceReadReplica) AddProxy(id *string, options *DatabaseProxyOptions) DatabaseProxy {
	if err := d.validateAddProxyParameters(id, options); err != nil {
		panic(err)
	}
	var returns DatabaseProxy

	_jsii_.Invoke(
		d,
		"addProxy",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DatabaseInstanceReadReplica) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	if err := d.validateApplyRemovalPolicyParameters(policy); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		d,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

func (d *jsiiProxy_DatabaseInstanceReadReplica) AsSecretAttachmentTarget() *awssecretsmanager.SecretAttachmentTargetProps {
	var returns *awssecretsmanager.SecretAttachmentTargetProps

	_jsii_.Invoke(
		d,
		"asSecretAttachmentTarget",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DatabaseInstanceReadReplica) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		d,
		"generatePhysicalName",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DatabaseInstanceReadReplica) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	if err := d.validateGetResourceArnAttributeParameters(arnAttr, arnComponents); err != nil {
		panic(err)
	}
	var returns *string

	_jsii_.Invoke(
		d,
		"getResourceArnAttribute",
		[]interface{}{arnAttr, arnComponents},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DatabaseInstanceReadReplica) GetResourceNameAttribute(nameAttr *string) *string {
	if err := d.validateGetResourceNameAttributeParameters(nameAttr); err != nil {
		panic(err)
	}
	var returns *string

	_jsii_.Invoke(
		d,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DatabaseInstanceReadReplica) GrantConnect(grantee awsiam.IGrantable, dbUser *string) awsiam.Grant {
	if err := d.validateGrantConnectParameters(grantee); err != nil {
		panic(err)
	}
	var returns awsiam.Grant

	_jsii_.Invoke(
		d,
		"grantConnect",
		[]interface{}{grantee, dbUser},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DatabaseInstanceReadReplica) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	if err := d.validateMetricParameters(metricName, props); err != nil {
		panic(err)
	}
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		d,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DatabaseInstanceReadReplica) MetricCPUUtilization(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	if err := d.validateMetricCPUUtilizationParameters(props); err != nil {
		panic(err)
	}
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		d,
		"metricCPUUtilization",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DatabaseInstanceReadReplica) MetricDatabaseConnections(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	if err := d.validateMetricDatabaseConnectionsParameters(props); err != nil {
		panic(err)
	}
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		d,
		"metricDatabaseConnections",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DatabaseInstanceReadReplica) MetricFreeableMemory(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	if err := d.validateMetricFreeableMemoryParameters(props); err != nil {
		panic(err)
	}
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		d,
		"metricFreeableMemory",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DatabaseInstanceReadReplica) MetricFreeStorageSpace(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	if err := d.validateMetricFreeStorageSpaceParameters(props); err != nil {
		panic(err)
	}
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		d,
		"metricFreeStorageSpace",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DatabaseInstanceReadReplica) MetricReadIOPS(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	if err := d.validateMetricReadIOPSParameters(props); err != nil {
		panic(err)
	}
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		d,
		"metricReadIOPS",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DatabaseInstanceReadReplica) MetricWriteIOPS(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	if err := d.validateMetricWriteIOPSParameters(props); err != nil {
		panic(err)
	}
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		d,
		"metricWriteIOPS",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DatabaseInstanceReadReplica) OnEvent(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	if err := d.validateOnEventParameters(id, options); err != nil {
		panic(err)
	}
	var returns awsevents.Rule

	_jsii_.Invoke(
		d,
		"onEvent",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DatabaseInstanceReadReplica) SetLogRetention() {
	_jsii_.InvokeVoid(
		d,
		"setLogRetention",
		nil, // no parameters
	)
}

func (d *jsiiProxy_DatabaseInstanceReadReplica) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		d,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

