package cdkpipelinesgithub


// Locations of GitHub Secrets used to authenticate to AWS.
// Experimental.
type GitHubSecretsProviderProps struct {
	// Default: "AWS_ACCESS_KEY_ID".
	//
	// Experimental.
	AccessKeyId *string `field:"required" json:"accessKeyId" yaml:"accessKeyId"`
	// Default: "AWS_SECRET_ACCESS_KEY".
	//
	// Experimental.
	SecretAccessKey *string `field:"required" json:"secretAccessKey" yaml:"secretAccessKey"`
	// Default: - no session token is used.
	//
	// Experimental.
	SessionToken *string `field:"optional" json:"sessionToken" yaml:"sessionToken"`
}

