package cdkpipelinesgithub

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/cdklabs/cdk-pipelines-github-go/cdkpipelinesgithub/jsii"

	"github.com/aws/aws-cdk-go/awscdk/v2/pipelines"
	"github.com/cdklabs/cdk-pipelines-github-go/cdkpipelinesgithub/internal"
)

// Specifies a GitHub Action as a step in the pipeline.
// Experimental.
type GitHubActionStep interface {
	pipelines.Step
	// Return the steps this step depends on, based on the FileSets it requires.
	// Experimental.
	Dependencies() *[]pipelines.Step
	// The list of FileSets consumed by this Step.
	// Experimental.
	DependencyFileSets() *[]pipelines.FileSet
	// Experimental.
	Env() *map[string]*string
	// Identifier for this step.
	// Experimental.
	Id() *string
	// Whether or not this is a Source step.
	//
	// What it means to be a Source step depends on the engine.
	// Experimental.
	IsSource() *bool
	// Experimental.
	JobSteps() *[]*JobStep
	// The primary FileSet produced by this Step.
	//
	// Not all steps produce an output FileSet--if they do
	// you can substitute the `Step` object for the `FileSet` object.
	// Experimental.
	PrimaryOutput() pipelines.FileSet
	// Add an additional FileSet to the set of file sets required by this step.
	//
	// This will lead to a dependency on the producer of that file set.
	// Experimental.
	AddDependencyFileSet(fs pipelines.FileSet)
	// Add a dependency on another step.
	// Experimental.
	AddStepDependency(step pipelines.Step)
	// Configure the given FileSet as the primary output of this step.
	// Experimental.
	ConfigurePrimaryOutput(fs pipelines.FileSet)
	// Return a string representation of this Step.
	// Experimental.
	ToString() *string
}

// The jsii proxy struct for GitHubActionStep
type jsiiProxy_GitHubActionStep struct {
	internal.Type__pipelinesStep
}

func (j *jsiiProxy_GitHubActionStep) Dependencies() *[]pipelines.Step {
	var returns *[]pipelines.Step
	_jsii_.Get(
		j,
		"dependencies",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubActionStep) DependencyFileSets() *[]pipelines.FileSet {
	var returns *[]pipelines.FileSet
	_jsii_.Get(
		j,
		"dependencyFileSets",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubActionStep) Env() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubActionStep) Id() *string {
	var returns *string
	_jsii_.Get(
		j,
		"id",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubActionStep) IsSource() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isSource",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubActionStep) JobSteps() *[]*JobStep {
	var returns *[]*JobStep
	_jsii_.Get(
		j,
		"jobSteps",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubActionStep) PrimaryOutput() pipelines.FileSet {
	var returns pipelines.FileSet
	_jsii_.Get(
		j,
		"primaryOutput",
		&returns,
	)
	return returns
}


// Experimental.
func NewGitHubActionStep(id *string, props *GitHubActionStepProps) GitHubActionStep {
	_init_.Initialize()

	if err := validateNewGitHubActionStepParameters(id, props); err != nil {
		panic(err)
	}
	j := jsiiProxy_GitHubActionStep{}

	_jsii_.Create(
		"cdk-pipelines-github.GitHubActionStep",
		[]interface{}{id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewGitHubActionStep_Override(g GitHubActionStep, id *string, props *GitHubActionStepProps) {
	_init_.Initialize()

	_jsii_.Create(
		"cdk-pipelines-github.GitHubActionStep",
		[]interface{}{id, props},
		g,
	)
}

// Define a sequence of steps to be executed in order.
// Experimental.
func GitHubActionStep_Sequence(steps *[]pipelines.Step) *[]pipelines.Step {
	_init_.Initialize()

	if err := validateGitHubActionStep_SequenceParameters(steps); err != nil {
		panic(err)
	}
	var returns *[]pipelines.Step

	_jsii_.StaticInvoke(
		"cdk-pipelines-github.GitHubActionStep",
		"sequence",
		[]interface{}{steps},
		&returns,
	)

	return returns
}

func (g *jsiiProxy_GitHubActionStep) AddDependencyFileSet(fs pipelines.FileSet) {
	if err := g.validateAddDependencyFileSetParameters(fs); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		g,
		"addDependencyFileSet",
		[]interface{}{fs},
	)
}

func (g *jsiiProxy_GitHubActionStep) AddStepDependency(step pipelines.Step) {
	if err := g.validateAddStepDependencyParameters(step); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		g,
		"addStepDependency",
		[]interface{}{step},
	)
}

func (g *jsiiProxy_GitHubActionStep) ConfigurePrimaryOutput(fs pipelines.FileSet) {
	if err := g.validateConfigurePrimaryOutputParameters(fs); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		g,
		"configurePrimaryOutput",
		[]interface{}{fs},
	)
}

func (g *jsiiProxy_GitHubActionStep) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		g,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

