package cdkpipelinesgithub


// Job level settings applied to all jobs in the workflow.
// Experimental.
type JobSettings struct {
	// jobs.<job_id>.if.
	// See: https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#jobsjob_idif
	//
	// Experimental.
	If *string `field:"optional" json:"if" yaml:"if"`
}

