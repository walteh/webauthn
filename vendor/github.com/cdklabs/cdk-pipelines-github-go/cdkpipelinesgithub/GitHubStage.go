package cdkpipelinesgithub

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/cdklabs/cdk-pipelines-github-go/cdkpipelinesgithub/jsii"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/cxapi"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/cdklabs/cdk-pipelines-github-go/cdkpipelinesgithub/internal"
)

// Experimental.
type GitHubStage interface {
	awscdk.Stage
	// The default account for all resources defined within this stage.
	// Experimental.
	Account() *string
	// Artifact ID of the assembly if it is a nested stage. The root stage (app) will return an empty string.
	//
	// Derived from the construct path.
	// Experimental.
	ArtifactId() *string
	// The cloud assembly asset output directory.
	// Experimental.
	AssetOutdir() *string
	// The tree node.
	// Experimental.
	Node() constructs.Node
	// The cloud assembly output directory.
	// Experimental.
	Outdir() *string
	// The parent stage or `undefined` if this is the app.
	//
	// *.
	// Experimental.
	ParentStage() awscdk.Stage
	// Experimental.
	Props() *GitHubStageProps
	// The default region for all resources defined within this stage.
	// Experimental.
	Region() *string
	// The name of the stage.
	//
	// Based on names of the parent stages separated by
	// hypens.
	// Experimental.
	StageName() *string
	// Synthesize this stage into a cloud assembly.
	//
	// Once an assembly has been synthesized, it cannot be modified. Subsequent
	// calls will return the same assembly.
	// Experimental.
	Synth(options *awscdk.StageSynthesisOptions) cxapi.CloudAssembly
	// Returns a string representation of this construct.
	// Experimental.
	ToString() *string
}

// The jsii proxy struct for GitHubStage
type jsiiProxy_GitHubStage struct {
	internal.Type__awscdkStage
}

func (j *jsiiProxy_GitHubStage) Account() *string {
	var returns *string
	_jsii_.Get(
		j,
		"account",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubStage) ArtifactId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"artifactId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubStage) AssetOutdir() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assetOutdir",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubStage) Node() constructs.Node {
	var returns constructs.Node
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubStage) Outdir() *string {
	var returns *string
	_jsii_.Get(
		j,
		"outdir",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubStage) ParentStage() awscdk.Stage {
	var returns awscdk.Stage
	_jsii_.Get(
		j,
		"parentStage",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubStage) Props() *GitHubStageProps {
	var returns *GitHubStageProps
	_jsii_.Get(
		j,
		"props",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubStage) Region() *string {
	var returns *string
	_jsii_.Get(
		j,
		"region",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubStage) StageName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stageName",
		&returns,
	)
	return returns
}


// Experimental.
func NewGitHubStage(scope constructs.Construct, id *string, props *GitHubStageProps) GitHubStage {
	_init_.Initialize()

	if err := validateNewGitHubStageParameters(scope, id, props); err != nil {
		panic(err)
	}
	j := jsiiProxy_GitHubStage{}

	_jsii_.Create(
		"cdk-pipelines-github.GitHubStage",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewGitHubStage_Override(g GitHubStage, scope constructs.Construct, id *string, props *GitHubStageProps) {
	_init_.Initialize()

	_jsii_.Create(
		"cdk-pipelines-github.GitHubStage",
		[]interface{}{scope, id, props},
		g,
	)
}

// Checks if `x` is a construct.
//
// Returns: true if `x` is an object created from a class which extends `Construct`.
// Deprecated: use `x instanceof Construct` instead.
func GitHubStage_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	if err := validateGitHubStage_IsConstructParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"cdk-pipelines-github.GitHubStage",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Test whether the given construct is a stage.
// Experimental.
func GitHubStage_IsStage(x interface{}) *bool {
	_init_.Initialize()

	if err := validateGitHubStage_IsStageParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"cdk-pipelines-github.GitHubStage",
		"isStage",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Return the stage this construct is contained with, if available.
//
// If called
// on a nested stage, returns its parent.
// Experimental.
func GitHubStage_Of(construct constructs.IConstruct) awscdk.Stage {
	_init_.Initialize()

	if err := validateGitHubStage_OfParameters(construct); err != nil {
		panic(err)
	}
	var returns awscdk.Stage

	_jsii_.StaticInvoke(
		"cdk-pipelines-github.GitHubStage",
		"of",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

func (g *jsiiProxy_GitHubStage) Synth(options *awscdk.StageSynthesisOptions) cxapi.CloudAssembly {
	if err := g.validateSynthParameters(options); err != nil {
		panic(err)
	}
	var returns cxapi.CloudAssembly

	_jsii_.Invoke(
		g,
		"synth",
		[]interface{}{options},
		&returns,
	)

	return returns
}

func (g *jsiiProxy_GitHubStage) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		g,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

