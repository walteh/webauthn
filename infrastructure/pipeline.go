package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/pipelines"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func NewPipelineStack(scope constructs.Construct, id string, opts *CdkStackProps) awscdk.Stack {
	// cdkpipelinesgithub.New()
	pipeline := pipelines.NewCodePipeline(scope, jsii.String("Pipeline"), &pipelines.CodePipelineProps{
		Synth: pipelines.NewShellStep(jsii.String("Synth"), &pipelines.ShellStepProps{
			Input:                  pipelines.CodePipelineSource_GitHub(jsii.String("myorg/repo1"), jsii.String("main"), &pipelines.GitHubSourceOptions{}),
			PrimaryOutputDirectory: jsii.String("./build"),
			Commands: &[]*string{
				jsii.String("./build.sh"),
			},
		}),
	})

	return pipeline.Pipeline().Stack()
}
