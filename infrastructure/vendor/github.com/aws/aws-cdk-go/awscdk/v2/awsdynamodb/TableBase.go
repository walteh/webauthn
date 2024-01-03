package awsdynamodb

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/v2/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudwatch"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb/internal"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awskms"
	"github.com/aws/constructs-go/constructs/v10"
)

type TableBase interface {
	awscdk.Resource
	ITable
	// KMS encryption key, if this table uses a customer-managed encryption key.
	EncryptionKey() awskms.IKey
	// The environment this resource belongs to.
	//
	// For resources that are created and managed by the CDK
	// (generally, those created by creating new class instances like Role, Bucket, etc.),
	// this is always the same as the environment of the stack they belong to;
	// however, for imported resources
	// (those obtained from static methods like fromRoleArn, fromBucketName, etc.),
	// that might be different than the stack they were imported into.
	Env() *awscdk.ResourceEnvironment
	HasIndex() *bool
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
	RegionalArns() *[]*string
	// The stack in which this resource is defined.
	Stack() awscdk.Stack
	// Arn of the dynamodb table.
	TableArn() *string
	// Table name of the dynamodb table.
	TableName() *string
	// ARN of the table's stream, if there is one.
	TableStreamArn() *string
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
	// Adds an IAM policy statement associated with this table to an IAM principal's policy.
	//
	// If `encryptionKey` is present, appropriate grants to the key needs to be added
	// separately using the `table.encryptionKey.grant*` methods.
	Grant(grantee awsiam.IGrantable, actions ...*string) awsiam.Grant
	// Permits all DynamoDB operations ("dynamodb:*") to an IAM principal.
	//
	// Appropriate grants will also be added to the customer-managed KMS key
	// if one was configured.
	GrantFullAccess(grantee awsiam.IGrantable) awsiam.Grant
	// Permits an IAM principal all data read operations from this table: BatchGetItem, GetRecords, GetShardIterator, Query, GetItem, Scan, DescribeTable.
	//
	// Appropriate grants will also be added to the customer-managed KMS key
	// if one was configured.
	GrantReadData(grantee awsiam.IGrantable) awsiam.Grant
	// Permits an IAM principal to all data read/write operations to this table.
	//
	// BatchGetItem, GetRecords, GetShardIterator, Query, GetItem, Scan,
	// BatchWriteItem, PutItem, UpdateItem, DeleteItem, DescribeTable
	//
	// Appropriate grants will also be added to the customer-managed KMS key
	// if one was configured.
	GrantReadWriteData(grantee awsiam.IGrantable) awsiam.Grant
	// Adds an IAM policy statement associated with this table's stream to an IAM principal's policy.
	//
	// If `encryptionKey` is present, appropriate grants to the key needs to be added
	// separately using the `table.encryptionKey.grant*` methods.
	GrantStream(grantee awsiam.IGrantable, actions ...*string) awsiam.Grant
	// Permits an IAM principal all stream data read operations for this table's stream: DescribeStream, GetRecords, GetShardIterator, ListStreams.
	//
	// Appropriate grants will also be added to the customer-managed KMS key
	// if one was configured.
	GrantStreamRead(grantee awsiam.IGrantable) awsiam.Grant
	// Permits an IAM Principal to list streams attached to current dynamodb table.
	GrantTableListStreams(grantee awsiam.IGrantable) awsiam.Grant
	// Permits an IAM principal all data write operations to this table: BatchWriteItem, PutItem, UpdateItem, DeleteItem, DescribeTable.
	//
	// Appropriate grants will also be added to the customer-managed KMS key
	// if one was configured.
	GrantWriteData(grantee awsiam.IGrantable) awsiam.Grant
	// Return the given named metric for this Table.
	//
	// By default, the metric will be calculated as a sum over a period of 5 minutes.
	// You can customize this by using the `statistic` and `period` properties.
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// Metric for the conditional check failed requests this table.
	//
	// By default, the metric will be calculated as a sum over a period of 5 minutes.
	// You can customize this by using the `statistic` and `period` properties.
	MetricConditionalCheckFailedRequests(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// Metric for the consumed read capacity units this table.
	//
	// By default, the metric will be calculated as a sum over a period of 5 minutes.
	// You can customize this by using the `statistic` and `period` properties.
	MetricConsumedReadCapacityUnits(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// Metric for the consumed write capacity units this table.
	//
	// By default, the metric will be calculated as a sum over a period of 5 minutes.
	// You can customize this by using the `statistic` and `period` properties.
	MetricConsumedWriteCapacityUnits(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// Metric for the successful request latency this table.
	//
	// By default, the metric will be calculated as an average over a period of 5 minutes.
	// You can customize this by using the `statistic` and `period` properties.
	MetricSuccessfulRequestLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// Metric for the system errors this table.
	// Deprecated: use `metricSystemErrorsForOperations`.
	MetricSystemErrors(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// Metric for the system errors this table.
	//
	// This will sum errors across all possible operations.
	// Note that by default, each individual metric will be calculated as a sum over a period of 5 minutes.
	// You can customize this by using the `statistic` and `period` properties.
	MetricSystemErrorsForOperations(props *SystemErrorsForOperationsMetricOptions) awscloudwatch.IMetric
	// How many requests are throttled on this table.
	//
	// Default: sum over 5 minutes.
	// Deprecated: Do not use this function. It returns an invalid metric. Use `metricThrottledRequestsForOperation` instead.
	MetricThrottledRequests(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// How many requests are throttled on this table, for the given operation.
	//
	// Default: sum over 5 minutes.
	MetricThrottledRequestsForOperation(operation *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// How many requests are throttled on this table.
	//
	// This will sum errors across all possible operations.
	// Note that by default, each individual metric will be calculated as a sum over a period of 5 minutes.
	// You can customize this by using the `statistic` and `period` properties.
	MetricThrottledRequestsForOperations(props *OperationsMetricOptions) awscloudwatch.IMetric
	// Metric for the user errors.
	//
	// Note that this metric reports user errors across all
	// the tables in the account and region the table resides in.
	//
	// By default, the metric will be calculated as a sum over a period of 5 minutes.
	// You can customize this by using the `statistic` and `period` properties.
	MetricUserErrors(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// Returns a string representation of this construct.
	ToString() *string
}

// The jsii proxy struct for TableBase
type jsiiProxy_TableBase struct {
	internal.Type__awscdkResource
	jsiiProxy_ITable
}

func (j *jsiiProxy_TableBase) EncryptionKey() awskms.IKey {
	var returns awskms.IKey
	_jsii_.Get(
		j,
		"encryptionKey",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TableBase) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TableBase) HasIndex() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"hasIndex",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TableBase) Node() constructs.Node {
	var returns constructs.Node
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TableBase) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TableBase) RegionalArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"regionalArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TableBase) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TableBase) TableArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"tableArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TableBase) TableName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"tableName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TableBase) TableStreamArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"tableStreamArn",
		&returns,
	)
	return returns
}


func NewTableBase_Override(t TableBase, scope constructs.Construct, id *string, props *awscdk.ResourceProps) {
	_init_.Initialize()

	_jsii_.Create(
		"aws-cdk-lib.aws_dynamodb.TableBase",
		[]interface{}{scope, id, props},
		t,
	)
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
func TableBase_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	if err := validateTableBase_IsConstructParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"aws-cdk-lib.aws_dynamodb.TableBase",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Returns true if the construct was created by CDK, and false otherwise.
func TableBase_IsOwnedResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	if err := validateTableBase_IsOwnedResourceParameters(construct); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"aws-cdk-lib.aws_dynamodb.TableBase",
		"isOwnedResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
func TableBase_IsResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	if err := validateTableBase_IsResourceParameters(construct); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"aws-cdk-lib.aws_dynamodb.TableBase",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

func (t *jsiiProxy_TableBase) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	if err := t.validateApplyRemovalPolicyParameters(policy); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		t,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

func (t *jsiiProxy_TableBase) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		t,
		"generatePhysicalName",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (t *jsiiProxy_TableBase) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	if err := t.validateGetResourceArnAttributeParameters(arnAttr, arnComponents); err != nil {
		panic(err)
	}
	var returns *string

	_jsii_.Invoke(
		t,
		"getResourceArnAttribute",
		[]interface{}{arnAttr, arnComponents},
		&returns,
	)

	return returns
}

func (t *jsiiProxy_TableBase) GetResourceNameAttribute(nameAttr *string) *string {
	if err := t.validateGetResourceNameAttributeParameters(nameAttr); err != nil {
		panic(err)
	}
	var returns *string

	_jsii_.Invoke(
		t,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

func (t *jsiiProxy_TableBase) Grant(grantee awsiam.IGrantable, actions ...*string) awsiam.Grant {
	if err := t.validateGrantParameters(grantee); err != nil {
		panic(err)
	}
	args := []interface{}{grantee}
	for _, a := range actions {
		args = append(args, a)
	}

	var returns awsiam.Grant

	_jsii_.Invoke(
		t,
		"grant",
		args,
		&returns,
	)

	return returns
}

func (t *jsiiProxy_TableBase) GrantFullAccess(grantee awsiam.IGrantable) awsiam.Grant {
	if err := t.validateGrantFullAccessParameters(grantee); err != nil {
		panic(err)
	}
	var returns awsiam.Grant

	_jsii_.Invoke(
		t,
		"grantFullAccess",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

func (t *jsiiProxy_TableBase) GrantReadData(grantee awsiam.IGrantable) awsiam.Grant {
	if err := t.validateGrantReadDataParameters(grantee); err != nil {
		panic(err)
	}
	var returns awsiam.Grant

	_jsii_.Invoke(
		t,
		"grantReadData",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

func (t *jsiiProxy_TableBase) GrantReadWriteData(grantee awsiam.IGrantable) awsiam.Grant {
	if err := t.validateGrantReadWriteDataParameters(grantee); err != nil {
		panic(err)
	}
	var returns awsiam.Grant

	_jsii_.Invoke(
		t,
		"grantReadWriteData",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

func (t *jsiiProxy_TableBase) GrantStream(grantee awsiam.IGrantable, actions ...*string) awsiam.Grant {
	if err := t.validateGrantStreamParameters(grantee); err != nil {
		panic(err)
	}
	args := []interface{}{grantee}
	for _, a := range actions {
		args = append(args, a)
	}

	var returns awsiam.Grant

	_jsii_.Invoke(
		t,
		"grantStream",
		args,
		&returns,
	)

	return returns
}

func (t *jsiiProxy_TableBase) GrantStreamRead(grantee awsiam.IGrantable) awsiam.Grant {
	if err := t.validateGrantStreamReadParameters(grantee); err != nil {
		panic(err)
	}
	var returns awsiam.Grant

	_jsii_.Invoke(
		t,
		"grantStreamRead",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

func (t *jsiiProxy_TableBase) GrantTableListStreams(grantee awsiam.IGrantable) awsiam.Grant {
	if err := t.validateGrantTableListStreamsParameters(grantee); err != nil {
		panic(err)
	}
	var returns awsiam.Grant

	_jsii_.Invoke(
		t,
		"grantTableListStreams",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

func (t *jsiiProxy_TableBase) GrantWriteData(grantee awsiam.IGrantable) awsiam.Grant {
	if err := t.validateGrantWriteDataParameters(grantee); err != nil {
		panic(err)
	}
	var returns awsiam.Grant

	_jsii_.Invoke(
		t,
		"grantWriteData",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

func (t *jsiiProxy_TableBase) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	if err := t.validateMetricParameters(metricName, props); err != nil {
		panic(err)
	}
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

func (t *jsiiProxy_TableBase) MetricConditionalCheckFailedRequests(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	if err := t.validateMetricConditionalCheckFailedRequestsParameters(props); err != nil {
		panic(err)
	}
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricConditionalCheckFailedRequests",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (t *jsiiProxy_TableBase) MetricConsumedReadCapacityUnits(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	if err := t.validateMetricConsumedReadCapacityUnitsParameters(props); err != nil {
		panic(err)
	}
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricConsumedReadCapacityUnits",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (t *jsiiProxy_TableBase) MetricConsumedWriteCapacityUnits(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	if err := t.validateMetricConsumedWriteCapacityUnitsParameters(props); err != nil {
		panic(err)
	}
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricConsumedWriteCapacityUnits",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (t *jsiiProxy_TableBase) MetricSuccessfulRequestLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	if err := t.validateMetricSuccessfulRequestLatencyParameters(props); err != nil {
		panic(err)
	}
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricSuccessfulRequestLatency",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (t *jsiiProxy_TableBase) MetricSystemErrors(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	if err := t.validateMetricSystemErrorsParameters(props); err != nil {
		panic(err)
	}
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricSystemErrors",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (t *jsiiProxy_TableBase) MetricSystemErrorsForOperations(props *SystemErrorsForOperationsMetricOptions) awscloudwatch.IMetric {
	if err := t.validateMetricSystemErrorsForOperationsParameters(props); err != nil {
		panic(err)
	}
	var returns awscloudwatch.IMetric

	_jsii_.Invoke(
		t,
		"metricSystemErrorsForOperations",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (t *jsiiProxy_TableBase) MetricThrottledRequests(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	if err := t.validateMetricThrottledRequestsParameters(props); err != nil {
		panic(err)
	}
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricThrottledRequests",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (t *jsiiProxy_TableBase) MetricThrottledRequestsForOperation(operation *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	if err := t.validateMetricThrottledRequestsForOperationParameters(operation, props); err != nil {
		panic(err)
	}
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricThrottledRequestsForOperation",
		[]interface{}{operation, props},
		&returns,
	)

	return returns
}

func (t *jsiiProxy_TableBase) MetricThrottledRequestsForOperations(props *OperationsMetricOptions) awscloudwatch.IMetric {
	if err := t.validateMetricThrottledRequestsForOperationsParameters(props); err != nil {
		panic(err)
	}
	var returns awscloudwatch.IMetric

	_jsii_.Invoke(
		t,
		"metricThrottledRequestsForOperations",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (t *jsiiProxy_TableBase) MetricUserErrors(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	if err := t.validateMetricUserErrorsParameters(props); err != nil {
		panic(err)
	}
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricUserErrors",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (t *jsiiProxy_TableBase) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		t,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

