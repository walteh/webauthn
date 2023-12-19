package cdkpipelinesgithub


// Names of secrets for AWS credentials.
// Experimental.
type AwsCredentialsSecrets struct {
	// Default: "AWS_ACCESS_KEY_ID".
	//
	// Experimental.
	AccessKeyId *string `field:"optional" json:"accessKeyId" yaml:"accessKeyId"`
	// Default: "AWS_SECRET_ACCESS_KEY".
	//
	// Experimental.
	SecretAccessKey *string `field:"optional" json:"secretAccessKey" yaml:"secretAccessKey"`
	// Default: - no session token is used.
	//
	// Experimental.
	SessionToken *string `field:"optional" json:"sessionToken" yaml:"sessionToken"`
}

