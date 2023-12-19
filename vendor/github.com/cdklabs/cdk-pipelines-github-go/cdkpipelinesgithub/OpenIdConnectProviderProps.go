package cdkpipelinesgithub


// Role to assume using OpenId Connect.
// Experimental.
type OpenIdConnectProviderProps struct {
	// A role that utilizes the GitHub OIDC Identity Provider in your AWS account.
	//
	// You can create your own role in the console with the necessary trust policy
	// to allow gitHub actions from your gitHub repository to assume the role, or
	// you can utilize the `GitHubActionRole` construct to create a role for you.
	// Experimental.
	GitHubActionRoleArn *string `field:"required" json:"gitHubActionRoleArn" yaml:"gitHubActionRoleArn"`
	// The role session name to use when assuming the role.
	// Default: - no role session name.
	//
	// Experimental.
	RoleSessionName *string `field:"optional" json:"roleSessionName" yaml:"roleSessionName"`
}

