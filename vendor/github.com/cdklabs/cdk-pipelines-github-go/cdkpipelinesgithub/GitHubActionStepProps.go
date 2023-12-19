package cdkpipelinesgithub


// Experimental.
type GitHubActionStepProps struct {
	// The Job steps.
	// Experimental.
	JobSteps *[]*JobStep `field:"required" json:"jobSteps" yaml:"jobSteps"`
	// Environment variables to set.
	// Experimental.
	Env *map[string]*string `field:"optional" json:"env" yaml:"env"`
}

