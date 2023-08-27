package invocation // import "go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"

// Version is the current release version of the AWS Lambda instrumentation.
func Version() string {
	return "0.77.0"
	// This string is updated by the pre_release.sh script during release
}

// SemVersion is the semantic version to be supplied to tracer/meter creation.
func SemVersion() string {
	return "semver:" + Version()
}

func Name() string {
	return "git.nugg.xyz/pkg/invocation"
}
