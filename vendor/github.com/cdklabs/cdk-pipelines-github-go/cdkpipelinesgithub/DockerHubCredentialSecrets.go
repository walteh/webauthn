package cdkpipelinesgithub


// Locations of GitHub Secrets used to authenticate to DockerHub.
// Experimental.
type DockerHubCredentialSecrets struct {
	// The key of the GitHub Secret containing the DockerHub personal access token.
	// Default: 'DOCKERHUB_TOKEN'.
	//
	// Experimental.
	PersonalAccessTokenKey *string `field:"optional" json:"personalAccessTokenKey" yaml:"personalAccessTokenKey"`
	// The key of the GitHub Secret containing the DockerHub username.
	// Default: 'DOCKERHUB_USERNAME'.
	//
	// Experimental.
	UsernameKey *string `field:"optional" json:"usernameKey" yaml:"usernameKey"`
}

