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
// Source: buf/alpha/registry/v1alpha1/labels.proto

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
	// LabelServiceName is the fully-qualified name of the LabelService service.
	LabelServiceName = "buf.alpha.registry.v1alpha1.LabelService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// LabelServiceCreateLabelProcedure is the fully-qualified name of the LabelService's CreateLabel
	// RPC.
	LabelServiceCreateLabelProcedure = "/buf.alpha.registry.v1alpha1.LabelService/CreateLabel"
	// LabelServiceMoveLabelProcedure is the fully-qualified name of the LabelService's MoveLabel RPC.
	LabelServiceMoveLabelProcedure = "/buf.alpha.registry.v1alpha1.LabelService/MoveLabel"
	// LabelServiceGetLabelsProcedure is the fully-qualified name of the LabelService's GetLabels RPC.
	LabelServiceGetLabelsProcedure = "/buf.alpha.registry.v1alpha1.LabelService/GetLabels"
	// LabelServiceGetLabelsInNamespaceProcedure is the fully-qualified name of the LabelService's
	// GetLabelsInNamespace RPC.
	LabelServiceGetLabelsInNamespaceProcedure = "/buf.alpha.registry.v1alpha1.LabelService/GetLabelsInNamespace"
)

// LabelServiceClient is a client for the buf.alpha.registry.v1alpha1.LabelService service.
type LabelServiceClient interface {
	CreateLabel(context.Context, *connect.Request[v1alpha1.CreateLabelRequest]) (*connect.Response[v1alpha1.CreateLabelResponse], error)
	MoveLabel(context.Context, *connect.Request[v1alpha1.MoveLabelRequest]) (*connect.Response[v1alpha1.MoveLabelResponse], error)
	// GetLabels returns labels in a repository with optional label name and value filters.
	GetLabels(context.Context, *connect.Request[v1alpha1.GetLabelsRequest]) (*connect.Response[v1alpha1.GetLabelsResponse], error)
	// GetLabelsInNamespace returns labels in a given namespace, optionally matching label names.
	GetLabelsInNamespace(context.Context, *connect.Request[v1alpha1.GetLabelsInNamespaceRequest]) (*connect.Response[v1alpha1.GetLabelsInNamespaceResponse], error)
}

// NewLabelServiceClient constructs a client for the buf.alpha.registry.v1alpha1.LabelService
// service. By default, it uses the Connect protocol with the binary Protobuf Codec, asks for
// gzipped responses, and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply
// the connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewLabelServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) LabelServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &labelServiceClient{
		createLabel: connect.NewClient[v1alpha1.CreateLabelRequest, v1alpha1.CreateLabelResponse](
			httpClient,
			baseURL+LabelServiceCreateLabelProcedure,
			connect.WithIdempotency(connect.IdempotencyIdempotent),
			connect.WithClientOptions(opts...),
		),
		moveLabel: connect.NewClient[v1alpha1.MoveLabelRequest, v1alpha1.MoveLabelResponse](
			httpClient,
			baseURL+LabelServiceMoveLabelProcedure,
			opts...,
		),
		getLabels: connect.NewClient[v1alpha1.GetLabelsRequest, v1alpha1.GetLabelsResponse](
			httpClient,
			baseURL+LabelServiceGetLabelsProcedure,
			connect.WithIdempotency(connect.IdempotencyNoSideEffects),
			connect.WithClientOptions(opts...),
		),
		getLabelsInNamespace: connect.NewClient[v1alpha1.GetLabelsInNamespaceRequest, v1alpha1.GetLabelsInNamespaceResponse](
			httpClient,
			baseURL+LabelServiceGetLabelsInNamespaceProcedure,
			connect.WithIdempotency(connect.IdempotencyNoSideEffects),
			connect.WithClientOptions(opts...),
		),
	}
}

// labelServiceClient implements LabelServiceClient.
type labelServiceClient struct {
	createLabel          *connect.Client[v1alpha1.CreateLabelRequest, v1alpha1.CreateLabelResponse]
	moveLabel            *connect.Client[v1alpha1.MoveLabelRequest, v1alpha1.MoveLabelResponse]
	getLabels            *connect.Client[v1alpha1.GetLabelsRequest, v1alpha1.GetLabelsResponse]
	getLabelsInNamespace *connect.Client[v1alpha1.GetLabelsInNamespaceRequest, v1alpha1.GetLabelsInNamespaceResponse]
}

// CreateLabel calls buf.alpha.registry.v1alpha1.LabelService.CreateLabel.
func (c *labelServiceClient) CreateLabel(ctx context.Context, req *connect.Request[v1alpha1.CreateLabelRequest]) (*connect.Response[v1alpha1.CreateLabelResponse], error) {
	return c.createLabel.CallUnary(ctx, req)
}

// MoveLabel calls buf.alpha.registry.v1alpha1.LabelService.MoveLabel.
func (c *labelServiceClient) MoveLabel(ctx context.Context, req *connect.Request[v1alpha1.MoveLabelRequest]) (*connect.Response[v1alpha1.MoveLabelResponse], error) {
	return c.moveLabel.CallUnary(ctx, req)
}

// GetLabels calls buf.alpha.registry.v1alpha1.LabelService.GetLabels.
func (c *labelServiceClient) GetLabels(ctx context.Context, req *connect.Request[v1alpha1.GetLabelsRequest]) (*connect.Response[v1alpha1.GetLabelsResponse], error) {
	return c.getLabels.CallUnary(ctx, req)
}

// GetLabelsInNamespace calls buf.alpha.registry.v1alpha1.LabelService.GetLabelsInNamespace.
func (c *labelServiceClient) GetLabelsInNamespace(ctx context.Context, req *connect.Request[v1alpha1.GetLabelsInNamespaceRequest]) (*connect.Response[v1alpha1.GetLabelsInNamespaceResponse], error) {
	return c.getLabelsInNamespace.CallUnary(ctx, req)
}

// LabelServiceHandler is an implementation of the buf.alpha.registry.v1alpha1.LabelService service.
type LabelServiceHandler interface {
	CreateLabel(context.Context, *connect.Request[v1alpha1.CreateLabelRequest]) (*connect.Response[v1alpha1.CreateLabelResponse], error)
	MoveLabel(context.Context, *connect.Request[v1alpha1.MoveLabelRequest]) (*connect.Response[v1alpha1.MoveLabelResponse], error)
	// GetLabels returns labels in a repository with optional label name and value filters.
	GetLabels(context.Context, *connect.Request[v1alpha1.GetLabelsRequest]) (*connect.Response[v1alpha1.GetLabelsResponse], error)
	// GetLabelsInNamespace returns labels in a given namespace, optionally matching label names.
	GetLabelsInNamespace(context.Context, *connect.Request[v1alpha1.GetLabelsInNamespaceRequest]) (*connect.Response[v1alpha1.GetLabelsInNamespaceResponse], error)
}

// NewLabelServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewLabelServiceHandler(svc LabelServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	labelServiceCreateLabelHandler := connect.NewUnaryHandler(
		LabelServiceCreateLabelProcedure,
		svc.CreateLabel,
		connect.WithIdempotency(connect.IdempotencyIdempotent),
		connect.WithHandlerOptions(opts...),
	)
	labelServiceMoveLabelHandler := connect.NewUnaryHandler(
		LabelServiceMoveLabelProcedure,
		svc.MoveLabel,
		opts...,
	)
	labelServiceGetLabelsHandler := connect.NewUnaryHandler(
		LabelServiceGetLabelsProcedure,
		svc.GetLabels,
		connect.WithIdempotency(connect.IdempotencyNoSideEffects),
		connect.WithHandlerOptions(opts...),
	)
	labelServiceGetLabelsInNamespaceHandler := connect.NewUnaryHandler(
		LabelServiceGetLabelsInNamespaceProcedure,
		svc.GetLabelsInNamespace,
		connect.WithIdempotency(connect.IdempotencyNoSideEffects),
		connect.WithHandlerOptions(opts...),
	)
	return "/buf.alpha.registry.v1alpha1.LabelService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case LabelServiceCreateLabelProcedure:
			labelServiceCreateLabelHandler.ServeHTTP(w, r)
		case LabelServiceMoveLabelProcedure:
			labelServiceMoveLabelHandler.ServeHTTP(w, r)
		case LabelServiceGetLabelsProcedure:
			labelServiceGetLabelsHandler.ServeHTTP(w, r)
		case LabelServiceGetLabelsInNamespaceProcedure:
			labelServiceGetLabelsInNamespaceHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedLabelServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedLabelServiceHandler struct{}

func (UnimplementedLabelServiceHandler) CreateLabel(context.Context, *connect.Request[v1alpha1.CreateLabelRequest]) (*connect.Response[v1alpha1.CreateLabelResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("buf.alpha.registry.v1alpha1.LabelService.CreateLabel is not implemented"))
}

func (UnimplementedLabelServiceHandler) MoveLabel(context.Context, *connect.Request[v1alpha1.MoveLabelRequest]) (*connect.Response[v1alpha1.MoveLabelResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("buf.alpha.registry.v1alpha1.LabelService.MoveLabel is not implemented"))
}

func (UnimplementedLabelServiceHandler) GetLabels(context.Context, *connect.Request[v1alpha1.GetLabelsRequest]) (*connect.Response[v1alpha1.GetLabelsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("buf.alpha.registry.v1alpha1.LabelService.GetLabels is not implemented"))
}

func (UnimplementedLabelServiceHandler) GetLabelsInNamespace(context.Context, *connect.Request[v1alpha1.GetLabelsInNamespaceRequest]) (*connect.Response[v1alpha1.GetLabelsInNamespaceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("buf.alpha.registry.v1alpha1.LabelService.GetLabelsInNamespace is not implemented"))
}
