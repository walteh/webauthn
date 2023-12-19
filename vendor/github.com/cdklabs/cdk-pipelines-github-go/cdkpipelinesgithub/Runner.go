package cdkpipelinesgithub

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/cdklabs/cdk-pipelines-github-go/cdkpipelinesgithub/jsii"
)

// The type of runner to run the job on.
//
// Can be GitHub or Self-hosted.
// In case of self-hosted, a list of labels can be supplied.
// Experimental.
type Runner interface {
	// Experimental.
	RunsOn() interface{}
}

// The jsii proxy struct for Runner
type jsiiProxy_Runner struct {
	_ byte // padding
}

func (j *jsiiProxy_Runner) RunsOn() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"runsOn",
		&returns,
	)
	return returns
}


// Creates a runner instance that sets runsOn to `self-hosted`.
//
// Additional labels can be supplied. There is no need to supply `self-hosted` as a label explicitly.
// Experimental.
func Runner_SelfHosted(labels *[]*string) Runner {
	_init_.Initialize()

	if err := validateRunner_SelfHostedParameters(labels); err != nil {
		panic(err)
	}
	var returns Runner

	_jsii_.StaticInvoke(
		"cdk-pipelines-github.Runner",
		"selfHosted",
		[]interface{}{labels},
		&returns,
	)

	return returns
}

func Runner_MACOS_LATEST() Runner {
	_init_.Initialize()
	var returns Runner
	_jsii_.StaticGet(
		"cdk-pipelines-github.Runner",
		"MACOS_LATEST",
		&returns,
	)
	return returns
}

func Runner_UBUNTU_LATEST() Runner {
	_init_.Initialize()
	var returns Runner
	_jsii_.StaticGet(
		"cdk-pipelines-github.Runner",
		"UBUNTU_LATEST",
		&returns,
	)
	return returns
}

func Runner_WINDOWS_LATEST() Runner {
	_init_.Initialize()
	var returns Runner
	_jsii_.StaticGet(
		"cdk-pipelines-github.Runner",
		"WINDOWS_LATEST",
		&returns,
	)
	return returns
}

