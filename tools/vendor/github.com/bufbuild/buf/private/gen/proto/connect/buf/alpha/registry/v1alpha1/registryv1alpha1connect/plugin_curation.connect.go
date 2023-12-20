// Copyright 2020-2023 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: buf/alpha/registry/v1alpha1/plugin_curation.proto

package registryv1alpha1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1alpha1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/registry/v1alpha1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_7_0

const (
	// PluginCurationServiceName is the fully-qualified name of the PluginCurationService service.
	PluginCurationServiceName = "buf.alpha.registry.v1alpha1.PluginCurationService"
	// CodeGenerationServiceName is the fully-qualified name of the CodeGenerationService service.
	CodeGenerationServiceName = "buf.alpha.registry.v1alpha1.CodeGenerationService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// PluginCurationServiceListCuratedPluginsProcedure is the fully-qualified name of the
	// PluginCurationService's ListCuratedPlugins RPC.
	PluginCurationServiceListCuratedPluginsProcedure = "/buf.alpha.registry.v1alpha1.PluginCurationService/ListCuratedPlugins"
	// PluginCurationServiceCreateCuratedPluginProcedure is the fully-qualified name of the
	// PluginCurationService's CreateCuratedPlugin RPC.
	PluginCurationServiceCreateCuratedPluginProcedure = "/buf.alpha.registry.v1alpha1.PluginCurationService/CreateCuratedPlugin"
	// PluginCurationServiceGetLatestCuratedPluginProcedure is the fully-qualified name of the
	// PluginCurationService's GetLatestCuratedPlugin RPC.
	PluginCurationServiceGetLatestCuratedPluginProcedure = "/buf.alpha.registry.v1alpha1.PluginCurationService/GetLatestCuratedPlugin"
	// PluginCurationServiceDeleteCuratedPluginProcedure is the fully-qualified name of the
	// PluginCurationService's DeleteCuratedPlugin RPC.
	PluginCurationServiceDeleteCuratedPluginProcedure = "/buf.alpha.registry.v1alpha1.PluginCurationService/DeleteCuratedPlugin"
	// CodeGenerationServiceGenerateCodeProcedure is the fully-qualified name of the
	// CodeGenerationService's GenerateCode RPC.
	CodeGenerationServiceGenerateCodeProcedure = "/buf.alpha.registry.v1alpha1.CodeGenerationService/GenerateCode"
)

// PluginCurationServiceClient is a client for the buf.alpha.registry.v1alpha1.PluginCurationService
// service.
type PluginCurationServiceClient interface {
	// ListCuratedPlugins returns all the curated plugins available.
	ListCuratedPlugins(context.Context, *connect.Request[v1alpha1.ListCuratedPluginsRequest]) (*connect.Response[v1alpha1.ListCuratedPluginsResponse], error)
	// CreateCuratedPlugin creates a new curated plugin.
	CreateCuratedPlugin(context.Context, *connect.Request[v1alpha1.CreateCuratedPluginRequest]) (*connect.Response[v1alpha1.CreateCuratedPluginResponse], error)
	// GetLatestCuratedPlugin returns the latest version of a plugin matching given parameters.
	GetLatestCuratedPlugin(context.Context, *connect.Request[v1alpha1.GetLatestCuratedPluginRequest]) (*connect.Response[v1alpha1.GetLatestCuratedPluginResponse], error)
	// DeleteCuratedPlugin deletes a curated plugin based on the given parameters.
	DeleteCuratedPlugin(context.Context, *connect.Request[v1alpha1.DeleteCuratedPluginRequest]) (*connect.Response[v1alpha1.DeleteCuratedPluginResponse], error)
}

// NewPluginCurationServiceClient constructs a client for the
// buf.alpha.registry.v1alpha1.PluginCurationService service. By default, it uses the Connect
// protocol with the binary Protobuf Codec, asks for gzipped responses, and sends uncompressed
// requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewPluginCurationServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) PluginCurationServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &pluginCurationServiceClient{
		listCuratedPlugins: connect.NewClient[v1alpha1.ListCuratedPluginsRequest, v1alpha1.ListCuratedPluginsResponse](
			httpClient,
			baseURL+PluginCurationServiceListCuratedPluginsProcedure,
			connect.WithIdempotency(connect.IdempotencyNoSideEffects),
			connect.WithClientOptions(opts...),
		),
		createCuratedPlugin: connect.NewClient[v1alpha1.CreateCuratedPluginRequest, v1alpha1.CreateCuratedPluginResponse](
			httpClient,
			baseURL+PluginCurationServiceCreateCuratedPluginProcedure,
			connect.WithIdempotency(connect.IdempotencyIdempotent),
			connect.WithClientOptions(opts...),
		),
		getLatestCuratedPlugin: connect.NewClient[v1alpha1.GetLatestCuratedPluginRequest, v1alpha1.GetLatestCuratedPluginResponse](
			httpClient,
			baseURL+PluginCurationServiceGetLatestCuratedPluginProcedure,
			connect.WithIdempotency(connect.IdempotencyNoSideEffects),
			connect.WithClientOptions(opts...),
		),
		deleteCuratedPlugin: connect.NewClient[v1alpha1.DeleteCuratedPluginRequest, v1alpha1.DeleteCuratedPluginResponse](
			httpClient,
			baseURL+PluginCurationServiceDeleteCuratedPluginProcedure,
			connect.WithIdempotency(connect.IdempotencyIdempotent),
			connect.WithClientOptions(opts...),
		),
	}
}

// pluginCurationServiceClient implements PluginCurationServiceClient.
type pluginCurationServiceClient struct {
	listCuratedPlugins     *connect.Client[v1alpha1.ListCuratedPluginsRequest, v1alpha1.ListCuratedPluginsResponse]
	createCuratedPlugin    *connect.Client[v1alpha1.CreateCuratedPluginRequest, v1alpha1.CreateCuratedPluginResponse]
	getLatestCuratedPlugin *connect.Client[v1alpha1.GetLatestCuratedPluginRequest, v1alpha1.GetLatestCuratedPluginResponse]
	deleteCuratedPlugin    *connect.Client[v1alpha1.DeleteCuratedPluginRequest, v1alpha1.DeleteCuratedPluginResponse]
}

// ListCuratedPlugins calls buf.alpha.registry.v1alpha1.PluginCurationService.ListCuratedPlugins.
func (c *pluginCurationServiceClient) ListCuratedPlugins(ctx context.Context, req *connect.Request[v1alpha1.ListCuratedPluginsRequest]) (*connect.Response[v1alpha1.ListCuratedPluginsResponse], error) {
	return c.listCuratedPlugins.CallUnary(ctx, req)
}

// CreateCuratedPlugin calls buf.alpha.registry.v1alpha1.PluginCurationService.CreateCuratedPlugin.
func (c *pluginCurationServiceClient) CreateCuratedPlugin(ctx context.Context, req *connect.Request[v1alpha1.CreateCuratedPluginRequest]) (*connect.Response[v1alpha1.CreateCuratedPluginResponse], error) {
	return c.createCuratedPlugin.CallUnary(ctx, req)
}

// GetLatestCuratedPlugin calls
// buf.alpha.registry.v1alpha1.PluginCurationService.GetLatestCuratedPlugin.
func (c *pluginCurationServiceClient) GetLatestCuratedPlugin(ctx context.Context, req *connect.Request[v1alpha1.GetLatestCuratedPluginRequest]) (*connect.Response[v1alpha1.GetLatestCuratedPluginResponse], error) {
	return c.getLatestCuratedPlugin.CallUnary(ctx, req)
}

// DeleteCuratedPlugin calls buf.alpha.registry.v1alpha1.PluginCurationService.DeleteCuratedPlugin.
func (c *pluginCurationServiceClient) DeleteCuratedPlugin(ctx context.Context, req *connect.Request[v1alpha1.DeleteCuratedPluginRequest]) (*connect.Response[v1alpha1.DeleteCuratedPluginResponse], error) {
	return c.deleteCuratedPlugin.CallUnary(ctx, req)
}

// PluginCurationServiceHandler is an implementation of the
// buf.alpha.registry.v1alpha1.PluginCurationService service.
type PluginCurationServiceHandler interface {
	// ListCuratedPlugins returns all the curated plugins available.
	ListCuratedPlugins(context.Context, *connect.Request[v1alpha1.ListCuratedPluginsRequest]) (*connect.Response[v1alpha1.ListCuratedPluginsResponse], error)
	// CreateCuratedPlugin creates a new curated plugin.
	CreateCuratedPlugin(context.Context, *connect.Request[v1alpha1.CreateCuratedPluginRequest]) (*connect.Response[v1alpha1.CreateCuratedPluginResponse], error)
	// GetLatestCuratedPlugin returns the latest version of a plugin matching given parameters.
	GetLatestCuratedPlugin(context.Context, *connect.Request[v1alpha1.GetLatestCuratedPluginRequest]) (*connect.Response[v1alpha1.GetLatestCuratedPluginResponse], error)
	// DeleteCuratedPlugin deletes a curated plugin based on the given parameters.
	DeleteCuratedPlugin(context.Context, *connect.Request[v1alpha1.DeleteCuratedPluginRequest]) (*connect.Response[v1alpha1.DeleteCuratedPluginResponse], error)
}

// NewPluginCurationServiceHandler builds an HTTP handler from the service implementation. It
// returns the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewPluginCurationServiceHandler(svc PluginCurationServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	pluginCurationServiceListCuratedPluginsHandler := connect.NewUnaryHandler(
		PluginCurationServiceListCuratedPluginsProcedure,
		svc.ListCuratedPlugins,
		connect.WithIdempotency(connect.IdempotencyNoSideEffects),
		connect.WithHandlerOptions(opts...),
	)
	pluginCurationServiceCreateCuratedPluginHandler := connect.NewUnaryHandler(
		PluginCurationServiceCreateCuratedPluginProcedure,
		svc.CreateCuratedPlugin,
		connect.WithIdempotency(connect.IdempotencyIdempotent),
		connect.WithHandlerOptions(opts...),
	)
	pluginCurationServiceGetLatestCuratedPluginHandler := connect.NewUnaryHandler(
		PluginCurationServiceGetLatestCuratedPluginProcedure,
		svc.GetLatestCuratedPlugin,
		connect.WithIdempotency(connect.IdempotencyNoSideEffects),
		connect.WithHandlerOptions(opts...),
	)
	pluginCurationServiceDeleteCuratedPluginHandler := connect.NewUnaryHandler(
		PluginCurationServiceDeleteCuratedPluginProcedure,
		svc.DeleteCuratedPlugin,
		connect.WithIdempotency(connect.IdempotencyIdempotent),
		connect.WithHandlerOptions(opts...),
	)
	return "/buf.alpha.registry.v1alpha1.PluginCurationService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case PluginCurationServiceListCuratedPluginsProcedure:
			pluginCurationServiceListCuratedPluginsHandler.ServeHTTP(w, r)
		case PluginCurationServiceCreateCuratedPluginProcedure:
			pluginCurationServiceCreateCuratedPluginHandler.ServeHTTP(w, r)
		case PluginCurationServiceGetLatestCuratedPluginProcedure:
			pluginCurationServiceGetLatestCuratedPluginHandler.ServeHTTP(w, r)
		case PluginCurationServiceDeleteCuratedPluginProcedure:
			pluginCurationServiceDeleteCuratedPluginHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedPluginCurationServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedPluginCurationServiceHandler struct{}

func (UnimplementedPluginCurationServiceHandler) ListCuratedPlugins(context.Context, *connect.Request[v1alpha1.ListCuratedPluginsRequest]) (*connect.Response[v1alpha1.ListCuratedPluginsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("buf.alpha.registry.v1alpha1.PluginCurationService.ListCuratedPlugins is not implemented"))
}

func (UnimplementedPluginCurationServiceHandler) CreateCuratedPlugin(context.Context, *connect.Request[v1alpha1.CreateCuratedPluginRequest]) (*connect.Response[v1alpha1.CreateCuratedPluginResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("buf.alpha.registry.v1alpha1.PluginCurationService.CreateCuratedPlugin is not implemented"))
}

func (UnimplementedPluginCurationServiceHandler) GetLatestCuratedPlugin(context.Context, *connect.Request[v1alpha1.GetLatestCuratedPluginRequest]) (*connect.Response[v1alpha1.GetLatestCuratedPluginResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("buf.alpha.registry.v1alpha1.PluginCurationService.GetLatestCuratedPlugin is not implemented"))
}

func (UnimplementedPluginCurationServiceHandler) DeleteCuratedPlugin(context.Context, *connect.Request[v1alpha1.DeleteCuratedPluginRequest]) (*connect.Response[v1alpha1.DeleteCuratedPluginResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("buf.alpha.registry.v1alpha1.PluginCurationService.DeleteCuratedPlugin is not implemented"))
}

// CodeGenerationServiceClient is a client for the buf.alpha.registry.v1alpha1.CodeGenerationService
// service.
type CodeGenerationServiceClient interface {
	// GenerateCode generates code using the specified remote plugins.
	GenerateCode(context.Context, *connect.Request[v1alpha1.GenerateCodeRequest]) (*connect.Response[v1alpha1.GenerateCodeResponse], error)
}

// NewCodeGenerationServiceClient constructs a client for the
// buf.alpha.registry.v1alpha1.CodeGenerationService service. By default, it uses the Connect
// protocol with the binary Protobuf Codec, asks for gzipped responses, and sends uncompressed
// requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewCodeGenerationServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) CodeGenerationServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &codeGenerationServiceClient{
		generateCode: connect.NewClient[v1alpha1.GenerateCodeRequest, v1alpha1.GenerateCodeResponse](
			httpClient,
			baseURL+CodeGenerationServiceGenerateCodeProcedure,
			opts...,
		),
	}
}

// codeGenerationServiceClient implements CodeGenerationServiceClient.
type codeGenerationServiceClient struct {
	generateCode *connect.Client[v1alpha1.GenerateCodeRequest, v1alpha1.GenerateCodeResponse]
}

// GenerateCode calls buf.alpha.registry.v1alpha1.CodeGenerationService.GenerateCode.
func (c *codeGenerationServiceClient) GenerateCode(ctx context.Context, req *connect.Request[v1alpha1.GenerateCodeRequest]) (*connect.Response[v1alpha1.GenerateCodeResponse], error) {
	return c.generateCode.CallUnary(ctx, req)
}

// CodeGenerationServiceHandler is an implementation of the
// buf.alpha.registry.v1alpha1.CodeGenerationService service.
type CodeGenerationServiceHandler interface {
	// GenerateCode generates code using the specified remote plugins.
	GenerateCode(context.Context, *connect.Request[v1alpha1.GenerateCodeRequest]) (*connect.Response[v1alpha1.GenerateCodeResponse], error)
}

// NewCodeGenerationServiceHandler builds an HTTP handler from the service implementation. It
// returns the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewCodeGenerationServiceHandler(svc CodeGenerationServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	codeGenerationServiceGenerateCodeHandler := connect.NewUnaryHandler(
		CodeGenerationServiceGenerateCodeProcedure,
		svc.GenerateCode,
		opts...,
	)
	return "/buf.alpha.registry.v1alpha1.CodeGenerationService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case CodeGenerationServiceGenerateCodeProcedure:
			codeGenerationServiceGenerateCodeHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedCodeGenerationServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedCodeGenerationServiceHandler struct{}

func (UnimplementedCodeGenerationServiceHandler) GenerateCode(context.Context, *connect.Request[v1alpha1.GenerateCodeRequest]) (*connect.Response[v1alpha1.GenerateCodeResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("buf.alpha.registry.v1alpha1.CodeGenerationService.GenerateCode is not implemented"))
}