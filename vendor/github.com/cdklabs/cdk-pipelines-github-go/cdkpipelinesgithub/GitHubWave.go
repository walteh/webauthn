package cdkpipelinesgithub

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/cdklabs/cdk-pipelines-github-go/cdkpipelinesgithub/jsii"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/pipelines"
	"github.com/cdklabs/cdk-pipelines-github-go/cdkpipelinesgithub/internal"
)

// Multiple stages that are deployed in parallel.
//
// A `Wave`, but with addition GitHub options
//
// Create with `GitHubWorkflow.addWave()` or `GitHubWorkflow.addGitHubWave()`.
// You should not have to instantiate a GitHubWave yourself.
// Experimental.
type GitHubWave interface {
	pipelines.Wave
	// Identifier for this Wave.
	// Experimental.
	Id() *string
	// Additional steps that are run after all of the stages in the wave.
	// Experimental.
	Post() *[]pipelines.Step
	// Additional steps that are run before any of the stages in the wave.
	// Experimental.
	Pre() *[]pipelines.Step
	// The stages that are deployed in this wave.
	// Experimental.
	Stages() *[]pipelines.StageDeployment
	// Add an additional step to run after all of the stages in this wave.
	// Experimental.
	AddPost(steps ...pipelines.Step)
	// Add an additional step to run before any of the stages in this wave.
	// Experimental.
	AddPre(steps ...pipelines.Step)
	// Add a Stage to this wave.
	//
	// It will be deployed in parallel with all other stages in this
	// wave.
	// Experimental.
	AddStage(stage awscdk.Stage, options *pipelines.AddStageOpts) pipelines.StageDeployment
	// Add a Stage to this wave.
	//
	// It will be deployed in parallel with all other stages in this
	// wave.
	// Experimental.
	AddStageWithGitHubOptions(stage awscdk.Stage, options *AddGitHubStageOptions) pipelines.StageDeployment
}

// The jsii proxy struct for GitHubWave
type jsiiProxy_GitHubWave struct {
	internal.Type__pipelinesWave
}

func (j *jsiiProxy_GitHubWave) Id() *string {
	var returns *string
	_jsii_.Get(
		j,
		"id",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubWave) Post() *[]pipelines.Step {
	var returns *[]pipelines.Step
	_jsii_.Get(
		j,
		"post",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubWave) Pre() *[]pipelines.Step {
	var returns *[]pipelines.Step
	_jsii_.Get(
		j,
		"pre",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubWave) Stages() *[]pipelines.StageDeployment {
	var returns *[]pipelines.StageDeployment
	_jsii_.Get(
		j,
		"stages",
		&returns,
	)
	return returns
}


// Create with `GitHubWorkflow.addWave()` or `GitHubWorkflow.addGitHubWave()`. You should not have to instantiate a GitHubWave yourself.
// Experimental.
func NewGitHubWave(id *string, pipeline GitHubWorkflow, props *pipelines.WaveProps) GitHubWave {
	_init_.Initialize()

	if err := validateNewGitHubWaveParameters(id, pipeline, props); err != nil {
		panic(err)
	}
	j := jsiiProxy_GitHubWave{}

	_jsii_.Create(
		"cdk-pipelines-github.GitHubWave",
		[]interface{}{id, pipeline, props},
		&j,
	)

	return &j
}

// Create with `GitHubWorkflow.addWave()` or `GitHubWorkflow.addGitHubWave()`. You should not have to instantiate a GitHubWave yourself.
// Experimental.
func NewGitHubWave_Override(g GitHubWave, id *string, pipeline GitHubWorkflow, props *pipelines.WaveProps) {
	_init_.Initialize()

	_jsii_.Create(
		"cdk-pipelines-github.GitHubWave",
		[]interface{}{id, pipeline, props},
		g,
	)
}

func (g *jsiiProxy_GitHubWave) AddPost(steps ...pipelines.Step) {
	args := []interface{}{}
	for _, a := range steps {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		g,
		"addPost",
		args,
	)
}

func (g *jsiiProxy_GitHubWave) AddPre(steps ...pipelines.Step) {
	args := []interface{}{}
	for _, a := range steps {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		g,
		"addPre",
		args,
	)
}

func (g *jsiiProxy_GitHubWave) AddStage(stage awscdk.Stage, options *pipelines.AddStageOpts) pipelines.StageDeployment {
	if err := g.validateAddStageParameters(stage, options); err != nil {
		panic(err)
	}
	var returns pipelines.StageDeployment

	_jsii_.Invoke(
		g,
		"addStage",
		[]interface{}{stage, options},
		&returns,
	)

	return returns
}

func (g *jsiiProxy_GitHubWave) AddStageWithGitHubOptions(stage awscdk.Stage, options *AddGitHubStageOptions) pipelines.StageDeployment {
	if err := g.validateAddStageWithGitHubOptionsParameters(stage, options); err != nil {
		panic(err)
	}
	var returns pipelines.StageDeployment

	_jsii_.Invoke(
		g,
		"addStageWithGitHubOptions",
		[]interface{}{stage, options},
		&returns,
	)

	return returns
}

