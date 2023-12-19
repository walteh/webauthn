package cdkpipelinesgithub


// Github environment with name and url.
// See: https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#jobsjob_idenvironment
//
// Experimental.
type GitHubEnvironment struct {
	// Name of the environment.
	// See: https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#example-using-environment-name-and-url
	//
	// Experimental.
	Name *string `field:"required" json:"name" yaml:"name"`
	// The url for the environment.
	// See: https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#example-using-environment-name-and-url
	//
	// Experimental.
	Url *string `field:"optional" json:"url" yaml:"url"`
}

