package awsrds

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/v2/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsrds/internal"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssecretsmanager"
	"github.com/aws/constructs-go/constructs/v10"
)

// RDS Database Proxy.
//
// Example:
//   var vpc vpc
//
//   cluster := rds.NewDatabaseCluster(this, jsii.String("Database"), &DatabaseClusterProps{
//   	Engine: rds.DatabaseClusterEngine_AuroraMysql(&AuroraMysqlClusterEngineProps{
//   		Version: rds.AuroraMysqlEngineVersion_VER_3_03_0(),
//   	}),
//   	Writer: rds.ClusterInstance_Provisioned(jsii.String("writer")),
//   	Vpc: Vpc,
//   })
//
//   proxy := rds.NewDatabaseProxy(this, jsii.String("Proxy"), &DatabaseProxyProps{
//   	ProxyTarget: rds.ProxyTarget_FromCluster(cluster),
//   	Secrets: []iSecret{
//   		cluster.Secret,
//   	},
//   	Vpc: Vpc,
//   })
//
//   role := iam.NewRole(this, jsii.String("DBProxyRole"), &RoleProps{
//   	AssumedBy: iam.NewAccountPrincipal(this.Account),
//   })
//   proxy.GrantConnect(role, jsii.String("admin"))
//
type DatabaseProxy interface {
	awscdk.Resource
	awsec2.IConnectable
	IDatabaseProxy
	awssecretsmanager.ISecretAttachmentTarget
	// Access to network connections.
	Connections() awsec2.Connections
	// DB Proxy ARN.
	DbProxyArn() *string
	// DB Proxy Name.
	DbProxyName() *string
	// Endpoint.
	Endpoint() *string
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
	// The stack in which this resource is defined.
	Stack() awscdk.Stack
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
	// Grant the given identity connection access to the proxy.
	GrantConnect(grantee awsiam.IGrantable, dbUser *string) awsiam.Grant
	// Returns a string representation of this construct.
	ToString() *string
}

// The jsii proxy struct for DatabaseProxy
type jsiiProxy_DatabaseProxy struct {
	internal.Type__awscdkResource
	internal.Type__awsec2IConnectable
	jsiiProxy_IDatabaseProxy
	internal.Type__awssecretsmanagerISecretAttachmentTarget
}

func (j *jsiiProxy_DatabaseProxy) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DatabaseProxy) DbProxyArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"dbProxyArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DatabaseProxy) DbProxyName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"dbProxyName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DatabaseProxy) Endpoint() *string {
	var returns *string
	_jsii_.Get(
		j,
		"endpoint",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DatabaseProxy) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DatabaseProxy) Node() constructs.Node {
	var returns constructs.Node
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DatabaseProxy) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DatabaseProxy) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


func NewDatabaseProxy(scope constructs.Construct, id *string, props *DatabaseProxyProps) DatabaseProxy {
	_init_.Initialize()

	if err := validateNewDatabaseProxyParameters(scope, id, props); err != nil {
		panic(err)
	}
	j := jsiiProxy_DatabaseProxy{}

	_jsii_.Create(
		"aws-cdk-lib.aws_rds.DatabaseProxy",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

func NewDatabaseProxy_Override(d DatabaseProxy, scope constructs.Construct, id *string, props *DatabaseProxyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"aws-cdk-lib.aws_rds.DatabaseProxy",
		[]interface{}{scope, id, props},
		d,
	)
}

// Import an existing database proxy.
func DatabaseProxy_FromDatabaseProxyAttributes(scope constructs.Construct, id *string, attrs *DatabaseProxyAttributes) IDatabaseProxy {
	_init_.Initialize()

	if err := validateDatabaseProxy_FromDatabaseProxyAttributesParameters(scope, id, attrs); err != nil {
		panic(err)
	}
	var returns IDatabaseProxy

	_jsii_.StaticInvoke(
		"aws-cdk-lib.aws_rds.DatabaseProxy",
		"fromDatabaseProxyAttributes",
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
func DatabaseProxy_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	if err := validateDatabaseProxy_IsConstructParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"aws-cdk-lib.aws_rds.DatabaseProxy",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Returns true if the construct was created by CDK, and false otherwise.
func DatabaseProxy_IsOwnedResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	if err := validateDatabaseProxy_IsOwnedResourceParameters(construct); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"aws-cdk-lib.aws_rds.DatabaseProxy",
		"isOwnedResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
func DatabaseProxy_IsResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	if err := validateDatabaseProxy_IsResourceParameters(construct); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"aws-cdk-lib.aws_rds.DatabaseProxy",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DatabaseProxy) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	if err := d.validateApplyRemovalPolicyParameters(policy); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		d,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

func (d *jsiiProxy_DatabaseProxy) AsSecretAttachmentTarget() *awssecretsmanager.SecretAttachmentTargetProps {
	var returns *awssecretsmanager.SecretAttachmentTargetProps

	_jsii_.Invoke(
		d,
		"asSecretAttachmentTarget",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DatabaseProxy) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		d,
		"generatePhysicalName",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DatabaseProxy) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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

func (d *jsiiProxy_DatabaseProxy) GetResourceNameAttribute(nameAttr *string) *string {
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

func (d *jsiiProxy_DatabaseProxy) GrantConnect(grantee awsiam.IGrantable, dbUser *string) awsiam.Grant {
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

func (d *jsiiProxy_DatabaseProxy) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		d,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}
