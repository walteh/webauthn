package awsrds

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
)

// Properties for defining a `CfnEventSubscription`.
//
// Example:
//   // The code below shows an example of how to instantiate this type.
//   // The values are placeholders you should change.
//   import "github.com/aws/aws-cdk-go/awscdk"
//
//   cfnEventSubscriptionProps := &CfnEventSubscriptionProps{
//   	SnsTopicArn: jsii.String("snsTopicArn"),
//
//   	// the properties below are optional
//   	Enabled: jsii.Boolean(false),
//   	EventCategories: []*string{
//   		jsii.String("eventCategories"),
//   	},
//   	SourceIds: []*string{
//   		jsii.String("sourceIds"),
//   	},
//   	SourceType: jsii.String("sourceType"),
//   	SubscriptionName: jsii.String("subscriptionName"),
//   	Tags: []cfnTag{
//   		&cfnTag{
//   			Key: jsii.String("key"),
//   			Value: jsii.String("value"),
//   		},
//   	},
//   }
//
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-rds-eventsubscription.html
//
type CfnEventSubscriptionProps struct {
	// The Amazon Resource Name (ARN) of the SNS topic created for event notification.
	//
	// The ARN is created by Amazon SNS when you create a topic and subscribe to it.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-rds-eventsubscription.html#cfn-rds-eventsubscription-snstopicarn
	//
	SnsTopicArn *string `field:"required" json:"snsTopicArn" yaml:"snsTopicArn"`
	// Specifies whether to activate the subscription.
	//
	// If the event notification subscription isn't activated, the subscription is created but not active.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-rds-eventsubscription.html#cfn-rds-eventsubscription-enabled
	//
	// Default: - true.
	//
	Enabled interface{} `field:"optional" json:"enabled" yaml:"enabled"`
	// A list of event categories for a particular source type ( `SourceType` ) that you want to subscribe to.
	//
	// You can see a list of the categories for a given source type in the "Amazon RDS event categories and event messages" section of the [*Amazon RDS User Guide*](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_Events.Messages.html) or the [*Amazon Aurora User Guide*](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/USER_Events.Messages.html) . You can also see this list by using the `DescribeEventCategories` operation.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-rds-eventsubscription.html#cfn-rds-eventsubscription-eventcategories
	//
	EventCategories *[]*string `field:"optional" json:"eventCategories" yaml:"eventCategories"`
	// The list of identifiers of the event sources for which events are returned.
	//
	// If not specified, then all sources are included in the response. An identifier must begin with a letter and must contain only ASCII letters, digits, and hyphens. It can't end with a hyphen or contain two consecutive hyphens.
	//
	// Constraints:
	//
	// - If a `SourceIds` value is supplied, `SourceType` must also be provided.
	// - If the source type is a DB instance, a `DBInstanceIdentifier` value must be supplied.
	// - If the source type is a DB cluster, a `DBClusterIdentifier` value must be supplied.
	// - If the source type is a DB parameter group, a `DBParameterGroupName` value must be supplied.
	// - If the source type is a DB security group, a `DBSecurityGroupName` value must be supplied.
	// - If the source type is a DB snapshot, a `DBSnapshotIdentifier` value must be supplied.
	// - If the source type is a DB cluster snapshot, a `DBClusterSnapshotIdentifier` value must be supplied.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-rds-eventsubscription.html#cfn-rds-eventsubscription-sourceids
	//
	SourceIds *[]*string `field:"optional" json:"sourceIds" yaml:"sourceIds"`
	// The type of source that is generating the events.
	//
	// For example, if you want to be notified of events generated by a DB instance, set this parameter to `db-instance` . If this value isn't specified, all events are returned.
	//
	// Valid values: `db-instance` | `db-cluster` | `db-parameter-group` | `db-security-group` | `db-snapshot` | `db-cluster-snapshot`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-rds-eventsubscription.html#cfn-rds-eventsubscription-sourcetype
	//
	SourceType *string `field:"optional" json:"sourceType" yaml:"sourceType"`
	// The name of the subscription.
	//
	// Constraints: The name must be less than 255 characters.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-rds-eventsubscription.html#cfn-rds-eventsubscription-subscriptionname
	//
	SubscriptionName *string `field:"optional" json:"subscriptionName" yaml:"subscriptionName"`
	// An optional array of key-value pairs to apply to this subscription.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-rds-eventsubscription.html#cfn-rds-eventsubscription-tags
	//
	Tags *[]*awscdk.CfnTag `field:"optional" json:"tags" yaml:"tags"`
}

