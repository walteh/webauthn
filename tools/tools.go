//go:build tools
// +build tools

package tools

// Manage tool dependencies via go.mod.
//
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// https://github.com/golang/go/issues/25922
//
//nolint:all
import (
	// _ "github.com/goreleaser/goreleaser"

	_ "gotest.tools/gotestsum"

	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"

	// mockery
	_ "github.com/vektra/mockery/v2"
	// buf
	_ "github.com/bufbuild/buf/cmd/buf"
	// hello
)
