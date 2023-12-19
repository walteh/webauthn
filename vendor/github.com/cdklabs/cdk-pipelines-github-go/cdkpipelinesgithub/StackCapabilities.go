package cdkpipelinesgithub


// Acknowledge IAM resources in AWS CloudFormation templates.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-iam-template.html#capabilities
//
// Experimental.
type StackCapabilities string

const (
	// Acknowledge your stack includes IAM resources.
	// Experimental.
	StackCapabilities_IAM StackCapabilities = "IAM"
	// Acknowledge your stack includes custom names for IAM resources.
	// Experimental.
	StackCapabilities_NAMED_IAM StackCapabilities = "NAMED_IAM"
	// Acknowledge your stack contains one or more macros.
	// Experimental.
	StackCapabilities_AUTO_EXPAND StackCapabilities = "AUTO_EXPAND"
)

