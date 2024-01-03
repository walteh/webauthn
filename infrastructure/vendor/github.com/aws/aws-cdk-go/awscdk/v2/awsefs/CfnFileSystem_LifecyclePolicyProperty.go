package awsefs


// Describes a policy used by EFS lifecycle management and EFS Intelligent-Tiering that specifies when to transition files into and out of the file system's Infrequent Access (IA) storage class.
//
// For more information, see [EFS Intelligent‐Tiering and EFS Lifecycle Management](https://docs.aws.amazon.com/efs/latest/ug/lifecycle-management-efs.html) .
//
// > - Each `LifecyclePolicy` object can have only a single transition. This means that in a request body, `LifecyclePolicies` must be structured as an array of `LifecyclePolicy` objects, one object for each transition, `TransitionToIA` , `TransitionToPrimaryStorageClass` .
// > - See the AWS::EFS::FileSystem examples for the correct `LifecyclePolicy` structure. Do not use the syntax shown on this page.
//
// Example:
//   // The code below shows an example of how to instantiate this type.
//   // The values are placeholders you should change.
//   import "github.com/aws/aws-cdk-go/awscdk"
//
//   lifecyclePolicyProperty := &LifecyclePolicyProperty{
//   	TransitionToArchive: jsii.String("transitionToArchive"),
//   	TransitionToIa: jsii.String("transitionToIa"),
//   	TransitionToPrimaryStorageClass: jsii.String("transitionToPrimaryStorageClass"),
//   }
//
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-efs-filesystem-lifecyclepolicy.html
//
type CfnFileSystem_LifecyclePolicyProperty struct {
	// The number of days after files were last accessed in primary storage (the Standard storage class) files at which to move them to Archive storage.
	//
	// Metadata operations such as listing the contents of a directory don't count as file access events.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-efs-filesystem-lifecyclepolicy.html#cfn-efs-filesystem-lifecyclepolicy-transitiontoarchive
	//
	TransitionToArchive *string `field:"optional" json:"transitionToArchive" yaml:"transitionToArchive"`
	// The number of days after files were last accessed in primary storage (the Standard storage class) at which to move them to Infrequent Access (IA) storage.
	//
	// Metadata operations such as listing the contents of a directory don't count as file access events.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-efs-filesystem-lifecyclepolicy.html#cfn-efs-filesystem-lifecyclepolicy-transitiontoia
	//
	TransitionToIa *string `field:"optional" json:"transitionToIa" yaml:"transitionToIa"`
	// Whether to move files back to primary (Standard) storage after they are accessed in IA or Archive storage.
	//
	// Metadata operations such as listing the contents of a directory don't count as file access events.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-efs-filesystem-lifecyclepolicy.html#cfn-efs-filesystem-lifecyclepolicy-transitiontoprimarystorageclass
	//
	TransitionToPrimaryStorageClass *string `field:"optional" json:"transitionToPrimaryStorageClass" yaml:"transitionToPrimaryStorageClass"`
}

