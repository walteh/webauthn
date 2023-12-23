package awsrds


// Metadata assigned to a DB proxy endpoint consisting of a key-value pair.
//
// Example:
//   // The code below shows an example of how to instantiate this type.
//   // The values are placeholders you should change.
//   import "github.com/aws/aws-cdk-go/awscdk"
//
//   tagFormatProperty := &TagFormatProperty{
//   	Key: jsii.String("key"),
//   	Value: jsii.String("value"),
//   }
//
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-rds-dbproxyendpoint-tagformat.html
//
type CfnDBProxyEndpoint_TagFormatProperty struct {
	// A value is the optional value of the tag.
	//
	// The string value can be 1-256 Unicode characters in length and can't be prefixed with `aws:` . The string can contain only the set of Unicode letters, digits, white-space, '_', '.', '/', '=', '+', '-' (Java regex: "^([\\p{L}\\p{Z}\\p{N}_.:/=+\\-]*)$").
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-rds-dbproxyendpoint-tagformat.html#cfn-rds-dbproxyendpoint-tagformat-key
	//
	Key *string `field:"optional" json:"key" yaml:"key"`
	// Metadata assigned to a DB instance consisting of a key-value pair.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-rds-dbproxyendpoint-tagformat.html#cfn-rds-dbproxyendpoint-tagformat-value
	//
	Value *string `field:"optional" json:"value" yaml:"value"`
}

