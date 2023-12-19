package cdkpipelinesgithub

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/cdklabs/cdk-pipelines-github-go/cdkpipelinesgithub/jsii"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/pipelines"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/cdklabs/cdk-pipelines-github-go/cdkpipelinesgithub/internal"
)

// CDK Pipelines for GitHub workflows.
// Experimental.
type GitHubWorkflow interface {
	pipelines.PipelineBase
	// The FileSet tha contains the cloud assembly.
	//
	// This is the primary output of the synth step.
	// Experimental.
	CloudAssemblyFileSet() pipelines.FileSet
	// The tree node.
	// Experimental.
	Node() constructs.Node
	// The build step that produces the CDK Cloud Assembly.
	// Experimental.
	Synth() pipelines.IFileSetProducer
	// The waves in this pipeline.
	// Experimental.
	Waves() *[]pipelines.Wave
	// Experimental.
	WorkflowFile() YamlFile
	// Experimental.
	WorkflowName() *string
	// Experimental.
	WorkflowPath() *string
	// Experimental.
	AddGitHubWave(id *string, options *pipelines.WaveOptions) GitHubWave
	// Deploy a single Stage by itself.
	//
	// Add a Stage to the pipeline, to be deployed in sequence with other
	// Stages added to the pipeline. All Stacks in the stage will be deployed
	// in an order automatically determined by their relative dependencies.
	// Experimental.
	AddStage(stage awscdk.Stage, options *pipelines.AddStageOpts) pipelines.StageDeployment
	// Deploy a single Stage by itself with options for further GitHub configuration.
	//
	// Add a Stage to the pipeline, to be deployed in sequence with other Stages added to the pipeline.
	// All Stacks in the stage will be deployed in an order automatically determined by their relative dependencies.
	// Experimental.
	AddStageWithGitHubOptions(stage awscdk.Stage, options *AddGitHubStageOptions) pipelines.StageDeployment
	// Add a Wave to the pipeline, for deploying multiple Stages in parallel.
	//
	// Use the return object of this method to deploy multiple stages in parallel.
	//
	// Example:
	//
	// ```ts
	// declare const pipeline: GitHubWorkflow; // assign pipeline a value
	//
	// const wave = pipeline.addWave('MyWave');
	// wave.addStage(new MyStage(this, 'Stage1'));
	// wave.addStage(new MyStage(this, 'Stage2'));
	// ```.
	// Experimental.
	AddWave(id *string, options *pipelines.WaveOptions) pipelines.Wave
	// Send the current pipeline definition to the engine, and construct the pipeline.
	//
	// It is not possible to modify the pipeline after calling this method.
	// Experimental.
	BuildPipeline()
	// Implemented by subclasses to do the actual pipeline construction.
	// Experimental.
	DoBuildPipeline()
	// Returns a string representation of this construct.
	// Experimental.
	ToString() *string
}

// The jsii proxy struct for GitHubWorkflow
type jsiiProxy_GitHubWorkflow struct {
	internal.Type__pipelinesPipelineBase
}

func (j *jsiiProxy_GitHubWorkflow) CloudAssemblyFileSet() pipelines.FileSet {
	var returns pipelines.FileSet
	_jsii_.Get(
		j,
		"cloudAssemblyFileSet",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubWorkflow) Node() constructs.Node {
	var returns constructs.Node
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubWorkflow) Synth() pipelines.IFileSetProducer {
	var returns pipelines.IFileSetProducer
	_jsii_.Get(
		j,
		"synth",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubWorkflow) Waves() *[]pipelines.Wave {
	var returns *[]pipelines.Wave
	_jsii_.Get(
		j,
		"waves",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubWorkflow) WorkflowFile() YamlFile {
	var returns YamlFile
	_jsii_.Get(
		j,
		"workflowFile",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubWorkflow) WorkflowName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"workflowName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubWorkflow) WorkflowPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"workflowPath",
		&returns,
	)
	return returns
}


// Experimental.
func NewGitHubWorkflow(scope constructs.Construct, id *string, props *GitHubWorkflowProps) GitHubWorkflow {
	_init_.Initialize()

	if err := validateNewGitHubWorkflowParameters(scope, id, props); err != nil {
		panic(err)
	}
	j := jsiiProxy_GitHubWorkflow{}

	_jsii_.Create(
		"cdk-pipelines-github.GitHubWorkflow",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewGitHubWorkflow_Override(g GitHubWorkflow, scope constructs.Construct, id *string, props *GitHubWorkflowProps) {
	_init_.Initialize()

	_jsii_.Create(
		"cdk-pipelines-github.GitHubWorkflow",
		[]interface{}{scope, id, props},
		g,
	)
}

// Checks if `x` is a construct.
//
// Returns: true if `x` is an object created from a class which extends `Construct`.
// Deprecated: use `x instanceof Construct` instead.
func GitHubWorkflow_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	if err := validateGitHubWorkflow_IsConstructParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"cdk-pipelines-github.GitHubWorkflow",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func (g *jsiiProxy_GitHubWorkflow) AddGitHubWave(id *string, options *pipelines.WaveOptions) GitHubWave {
	if err := g.validateAddGitHubWaveParameters(id, options); err != nil {
		panic(err)
	}
	var returns GitHubWave

	_jsii_.Invoke(
		g,
		"addGitHubWave",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (g *jsiiProxy_GitHubWorkflow) AddStage(stage awscdk.Stage, options *pipelines.AddStageOpts) pipelines.StageDeployment {
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

func (g *jsiiProxy_GitHubWorkflow) AddStageWithGitHubOptions(stage awscdk.Stage, options *AddGitHubStageOptions) pipelines.StageDeployment {
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

func (g *jsiiProxy_GitHubWorkflow) AddWave(id *string, options *pipelines.WaveOptions) pipelines.Wave {
	if err := g.validateAddWaveParameters(id, options); err != nil {
		panic(err)
	}
	var returns pipelines.Wave

	_jsii_.Invoke(
		g,
		"addWave",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (g *jsiiProxy_GitHubWorkflow) BuildPipeline() {
	_jsii_.InvokeVoid(
		g,
		"buildPipeline",
		nil, // no parameters
	)
}

func (g *jsiiProxy_GitHubWorkflow) DoBuildPipeline() {
	_jsii_.InvokeVoid(
		g,
		"doBuildPipeline",
		nil, // no parameters
	)
}

func (g *jsiiProxy_GitHubWorkflow) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		g,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

