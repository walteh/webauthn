package cdkpipelinesgithub


// Generic structure to supply the locations of GitHub Secrets used to authenticate to a docker registry.
// Experimental.
type ExternalDockerCredentialSecrets struct {
	// The key of the GitHub Secret containing your registry password.
	// Experimental.
	PasswordKey *string `field:"required" json:"passwordKey" yaml:"passwordKey"`
	// The key of the GitHub Secret containing your registry username.
	// Experimental.
	UsernameKey *string `field:"required" json:"usernameKey" yaml:"usernameKey"`
}

