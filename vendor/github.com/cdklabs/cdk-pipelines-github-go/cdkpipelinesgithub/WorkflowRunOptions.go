package cdkpipelinesgithub


// Workflow run options.
// Experimental.
type WorkflowRunOptions struct {
	// Which activity types to trigger on.
	// Experimental.
	Types *[]*string `field:"optional" json:"types" yaml:"types"`
}

