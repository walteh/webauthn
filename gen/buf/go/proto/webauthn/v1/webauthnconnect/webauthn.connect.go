// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: proto/webauthn/v1/webauthn.proto

package webauthnconnect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "github.com/walteh/webauthn/gen/buf/go/proto/webauthn/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// CreateChallengeServiceName is the fully-qualified name of the CreateChallengeService service.
	CreateChallengeServiceName = "proto.webauthn.v1.CreateChallengeService"
	// ApplePasskeyCreateServiceName is the fully-qualified name of the ApplePasskeyCreateService
	// service.
	ApplePasskeyCreateServiceName = "proto.webauthn.v1.ApplePasskeyCreateService"
	// ApplePasskeyAssertServiceName is the fully-qualified name of the ApplePasskeyAssertService
	// service.
	ApplePasskeyAssertServiceName = "proto.webauthn.v1.ApplePasskeyAssertService"
	// AppleDeviceCreateServiceName is the fully-qualified name of the AppleDeviceCreateService service.
	AppleDeviceCreateServiceName = "proto.webauthn.v1.AppleDeviceCreateService"
	// AppleDeviceAssertServiceName is the fully-qualified name of the AppleDeviceAssertService service.
	AppleDeviceAssertServiceName = "proto.webauthn.v1.AppleDeviceAssertService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// CreateChallengeServiceCreateChallengeProcedure is the fully-qualified name of the
	// CreateChallengeService's CreateChallenge RPC.
	CreateChallengeServiceCreateChallengeProcedure = "/proto.webauthn.v1.CreateChallengeService/CreateChallenge"
	// ApplePasskeyCreateServiceApplePasskeyCreateProcedure is the fully-qualified name of the
	// ApplePasskeyCreateService's ApplePasskeyCreate RPC.
	ApplePasskeyCreateServiceApplePasskeyCreateProcedure = "/proto.webauthn.v1.ApplePasskeyCreateService/ApplePasskeyCreate"
	// ApplePasskeyAssertServiceApplePasskeyAssertProcedure is the fully-qualified name of the
	// ApplePasskeyAssertService's ApplePasskeyAssert RPC.
	ApplePasskeyAssertServiceApplePasskeyAssertProcedure = "/proto.webauthn.v1.ApplePasskeyAssertService/ApplePasskeyAssert"
	// AppleDeviceCreateServiceAppleDeviceCreateProcedure is the fully-qualified name of the
	// AppleDeviceCreateService's AppleDeviceCreate RPC.
	AppleDeviceCreateServiceAppleDeviceCreateProcedure = "/proto.webauthn.v1.AppleDeviceCreateService/AppleDeviceCreate"
	// AppleDeviceAssertServiceAppleDeviceAssertProcedure is the fully-qualified name of the
	// AppleDeviceAssertService's AppleDeviceAssert RPC.
	AppleDeviceAssertServiceAppleDeviceAssertProcedure = "/proto.webauthn.v1.AppleDeviceAssertService/AppleDeviceAssert"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	createChallengeServiceServiceDescriptor                     = v1.File_proto_webauthn_v1_webauthn_proto.Services().ByName("CreateChallengeService")
	createChallengeServiceCreateChallengeMethodDescriptor       = createChallengeServiceServiceDescriptor.Methods().ByName("CreateChallenge")
	applePasskeyCreateServiceServiceDescriptor                  = v1.File_proto_webauthn_v1_webauthn_proto.Services().ByName("ApplePasskeyCreateService")
	applePasskeyCreateServiceApplePasskeyCreateMethodDescriptor = applePasskeyCreateServiceServiceDescriptor.Methods().ByName("ApplePasskeyCreate")
	applePasskeyAssertServiceServiceDescriptor                  = v1.File_proto_webauthn_v1_webauthn_proto.Services().ByName("ApplePasskeyAssertService")
	applePasskeyAssertServiceApplePasskeyAssertMethodDescriptor = applePasskeyAssertServiceServiceDescriptor.Methods().ByName("ApplePasskeyAssert")
	appleDeviceCreateServiceServiceDescriptor                   = v1.File_proto_webauthn_v1_webauthn_proto.Services().ByName("AppleDeviceCreateService")
	appleDeviceCreateServiceAppleDeviceCreateMethodDescriptor   = appleDeviceCreateServiceServiceDescriptor.Methods().ByName("AppleDeviceCreate")
	appleDeviceAssertServiceServiceDescriptor                   = v1.File_proto_webauthn_v1_webauthn_proto.Services().ByName("AppleDeviceAssertService")
	appleDeviceAssertServiceAppleDeviceAssertMethodDescriptor   = appleDeviceAssertServiceServiceDescriptor.Methods().ByName("AppleDeviceAssert")
)

// CreateChallengeServiceClient is a client for the proto.webauthn.v1.CreateChallengeService
// service.
type CreateChallengeServiceClient interface {
	CreateChallenge(context.Context, *connect.Request[v1.CreateChallengeRequest]) (*connect.Response[v1.CreateChallengeResponse], error)
}

// NewCreateChallengeServiceClient constructs a client for the
// proto.webauthn.v1.CreateChallengeService service. By default, it uses the Connect protocol with
// the binary Protobuf Codec, asks for gzipped responses, and sends uncompressed requests. To use
// the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewCreateChallengeServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) CreateChallengeServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &createChallengeServiceClient{
		createChallenge: connect.NewClient[v1.CreateChallengeRequest, v1.CreateChallengeResponse](
			httpClient,
			baseURL+CreateChallengeServiceCreateChallengeProcedure,
			connect.WithSchema(createChallengeServiceCreateChallengeMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// createChallengeServiceClient implements CreateChallengeServiceClient.
type createChallengeServiceClient struct {
	createChallenge *connect.Client[v1.CreateChallengeRequest, v1.CreateChallengeResponse]
}

// CreateChallenge calls proto.webauthn.v1.CreateChallengeService.CreateChallenge.
func (c *createChallengeServiceClient) CreateChallenge(ctx context.Context, req *connect.Request[v1.CreateChallengeRequest]) (*connect.Response[v1.CreateChallengeResponse], error) {
	return c.createChallenge.CallUnary(ctx, req)
}

// CreateChallengeServiceHandler is an implementation of the
// proto.webauthn.v1.CreateChallengeService service.
type CreateChallengeServiceHandler interface {
	CreateChallenge(context.Context, *connect.Request[v1.CreateChallengeRequest]) (*connect.Response[v1.CreateChallengeResponse], error)
}

// NewCreateChallengeServiceHandler builds an HTTP handler from the service implementation. It
// returns the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewCreateChallengeServiceHandler(svc CreateChallengeServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	createChallengeServiceCreateChallengeHandler := connect.NewUnaryHandler(
		CreateChallengeServiceCreateChallengeProcedure,
		svc.CreateChallenge,
		connect.WithSchema(createChallengeServiceCreateChallengeMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/proto.webauthn.v1.CreateChallengeService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case CreateChallengeServiceCreateChallengeProcedure:
			createChallengeServiceCreateChallengeHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedCreateChallengeServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedCreateChallengeServiceHandler struct{}

func (UnimplementedCreateChallengeServiceHandler) CreateChallenge(context.Context, *connect.Request[v1.CreateChallengeRequest]) (*connect.Response[v1.CreateChallengeResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("proto.webauthn.v1.CreateChallengeService.CreateChallenge is not implemented"))
}

// ApplePasskeyCreateServiceClient is a client for the proto.webauthn.v1.ApplePasskeyCreateService
// service.
type ApplePasskeyCreateServiceClient interface {
	ApplePasskeyCreate(context.Context, *connect.Request[v1.ApplePasskeyCreateRequest]) (*connect.Response[v1.ApplePasskeyCreateResponse], error)
}

// NewApplePasskeyCreateServiceClient constructs a client for the
// proto.webauthn.v1.ApplePasskeyCreateService service. By default, it uses the Connect protocol
// with the binary Protobuf Codec, asks for gzipped responses, and sends uncompressed requests. To
// use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or connect.WithGRPCWeb()
// options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewApplePasskeyCreateServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) ApplePasskeyCreateServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &applePasskeyCreateServiceClient{
		applePasskeyCreate: connect.NewClient[v1.ApplePasskeyCreateRequest, v1.ApplePasskeyCreateResponse](
			httpClient,
			baseURL+ApplePasskeyCreateServiceApplePasskeyCreateProcedure,
			connect.WithSchema(applePasskeyCreateServiceApplePasskeyCreateMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// applePasskeyCreateServiceClient implements ApplePasskeyCreateServiceClient.
type applePasskeyCreateServiceClient struct {
	applePasskeyCreate *connect.Client[v1.ApplePasskeyCreateRequest, v1.ApplePasskeyCreateResponse]
}

// ApplePasskeyCreate calls proto.webauthn.v1.ApplePasskeyCreateService.ApplePasskeyCreate.
func (c *applePasskeyCreateServiceClient) ApplePasskeyCreate(ctx context.Context, req *connect.Request[v1.ApplePasskeyCreateRequest]) (*connect.Response[v1.ApplePasskeyCreateResponse], error) {
	return c.applePasskeyCreate.CallUnary(ctx, req)
}

// ApplePasskeyCreateServiceHandler is an implementation of the
// proto.webauthn.v1.ApplePasskeyCreateService service.
type ApplePasskeyCreateServiceHandler interface {
	ApplePasskeyCreate(context.Context, *connect.Request[v1.ApplePasskeyCreateRequest]) (*connect.Response[v1.ApplePasskeyCreateResponse], error)
}

// NewApplePasskeyCreateServiceHandler builds an HTTP handler from the service implementation. It
// returns the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewApplePasskeyCreateServiceHandler(svc ApplePasskeyCreateServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	applePasskeyCreateServiceApplePasskeyCreateHandler := connect.NewUnaryHandler(
		ApplePasskeyCreateServiceApplePasskeyCreateProcedure,
		svc.ApplePasskeyCreate,
		connect.WithSchema(applePasskeyCreateServiceApplePasskeyCreateMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/proto.webauthn.v1.ApplePasskeyCreateService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ApplePasskeyCreateServiceApplePasskeyCreateProcedure:
			applePasskeyCreateServiceApplePasskeyCreateHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedApplePasskeyCreateServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedApplePasskeyCreateServiceHandler struct{}

func (UnimplementedApplePasskeyCreateServiceHandler) ApplePasskeyCreate(context.Context, *connect.Request[v1.ApplePasskeyCreateRequest]) (*connect.Response[v1.ApplePasskeyCreateResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("proto.webauthn.v1.ApplePasskeyCreateService.ApplePasskeyCreate is not implemented"))
}

// ApplePasskeyAssertServiceClient is a client for the proto.webauthn.v1.ApplePasskeyAssertService
// service.
type ApplePasskeyAssertServiceClient interface {
	ApplePasskeyAssert(context.Context, *connect.Request[v1.ApplePasskeyAssertRequest]) (*connect.Response[v1.ApplePasskeyAssertResponse], error)
}

// NewApplePasskeyAssertServiceClient constructs a client for the
// proto.webauthn.v1.ApplePasskeyAssertService service. By default, it uses the Connect protocol
// with the binary Protobuf Codec, asks for gzipped responses, and sends uncompressed requests. To
// use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or connect.WithGRPCWeb()
// options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewApplePasskeyAssertServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) ApplePasskeyAssertServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &applePasskeyAssertServiceClient{
		applePasskeyAssert: connect.NewClient[v1.ApplePasskeyAssertRequest, v1.ApplePasskeyAssertResponse](
			httpClient,
			baseURL+ApplePasskeyAssertServiceApplePasskeyAssertProcedure,
			connect.WithSchema(applePasskeyAssertServiceApplePasskeyAssertMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// applePasskeyAssertServiceClient implements ApplePasskeyAssertServiceClient.
type applePasskeyAssertServiceClient struct {
	applePasskeyAssert *connect.Client[v1.ApplePasskeyAssertRequest, v1.ApplePasskeyAssertResponse]
}

// ApplePasskeyAssert calls proto.webauthn.v1.ApplePasskeyAssertService.ApplePasskeyAssert.
func (c *applePasskeyAssertServiceClient) ApplePasskeyAssert(ctx context.Context, req *connect.Request[v1.ApplePasskeyAssertRequest]) (*connect.Response[v1.ApplePasskeyAssertResponse], error) {
	return c.applePasskeyAssert.CallUnary(ctx, req)
}

// ApplePasskeyAssertServiceHandler is an implementation of the
// proto.webauthn.v1.ApplePasskeyAssertService service.
type ApplePasskeyAssertServiceHandler interface {
	ApplePasskeyAssert(context.Context, *connect.Request[v1.ApplePasskeyAssertRequest]) (*connect.Response[v1.ApplePasskeyAssertResponse], error)
}

// NewApplePasskeyAssertServiceHandler builds an HTTP handler from the service implementation. It
// returns the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewApplePasskeyAssertServiceHandler(svc ApplePasskeyAssertServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	applePasskeyAssertServiceApplePasskeyAssertHandler := connect.NewUnaryHandler(
		ApplePasskeyAssertServiceApplePasskeyAssertProcedure,
		svc.ApplePasskeyAssert,
		connect.WithSchema(applePasskeyAssertServiceApplePasskeyAssertMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/proto.webauthn.v1.ApplePasskeyAssertService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ApplePasskeyAssertServiceApplePasskeyAssertProcedure:
			applePasskeyAssertServiceApplePasskeyAssertHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedApplePasskeyAssertServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedApplePasskeyAssertServiceHandler struct{}

func (UnimplementedApplePasskeyAssertServiceHandler) ApplePasskeyAssert(context.Context, *connect.Request[v1.ApplePasskeyAssertRequest]) (*connect.Response[v1.ApplePasskeyAssertResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("proto.webauthn.v1.ApplePasskeyAssertService.ApplePasskeyAssert is not implemented"))
}

// AppleDeviceCreateServiceClient is a client for the proto.webauthn.v1.AppleDeviceCreateService
// service.
type AppleDeviceCreateServiceClient interface {
	AppleDeviceCreate(context.Context, *connect.Request[v1.AppleDeviceCreateRequest]) (*connect.Response[v1.AppleDeviceCreateResponse], error)
}

// NewAppleDeviceCreateServiceClient constructs a client for the
// proto.webauthn.v1.AppleDeviceCreateService service. By default, it uses the Connect protocol with
// the binary Protobuf Codec, asks for gzipped responses, and sends uncompressed requests. To use
// the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewAppleDeviceCreateServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) AppleDeviceCreateServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &appleDeviceCreateServiceClient{
		appleDeviceCreate: connect.NewClient[v1.AppleDeviceCreateRequest, v1.AppleDeviceCreateResponse](
			httpClient,
			baseURL+AppleDeviceCreateServiceAppleDeviceCreateProcedure,
			connect.WithSchema(appleDeviceCreateServiceAppleDeviceCreateMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// appleDeviceCreateServiceClient implements AppleDeviceCreateServiceClient.
type appleDeviceCreateServiceClient struct {
	appleDeviceCreate *connect.Client[v1.AppleDeviceCreateRequest, v1.AppleDeviceCreateResponse]
}

// AppleDeviceCreate calls proto.webauthn.v1.AppleDeviceCreateService.AppleDeviceCreate.
func (c *appleDeviceCreateServiceClient) AppleDeviceCreate(ctx context.Context, req *connect.Request[v1.AppleDeviceCreateRequest]) (*connect.Response[v1.AppleDeviceCreateResponse], error) {
	return c.appleDeviceCreate.CallUnary(ctx, req)
}

// AppleDeviceCreateServiceHandler is an implementation of the
// proto.webauthn.v1.AppleDeviceCreateService service.
type AppleDeviceCreateServiceHandler interface {
	AppleDeviceCreate(context.Context, *connect.Request[v1.AppleDeviceCreateRequest]) (*connect.Response[v1.AppleDeviceCreateResponse], error)
}

// NewAppleDeviceCreateServiceHandler builds an HTTP handler from the service implementation. It
// returns the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewAppleDeviceCreateServiceHandler(svc AppleDeviceCreateServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	appleDeviceCreateServiceAppleDeviceCreateHandler := connect.NewUnaryHandler(
		AppleDeviceCreateServiceAppleDeviceCreateProcedure,
		svc.AppleDeviceCreate,
		connect.WithSchema(appleDeviceCreateServiceAppleDeviceCreateMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/proto.webauthn.v1.AppleDeviceCreateService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case AppleDeviceCreateServiceAppleDeviceCreateProcedure:
			appleDeviceCreateServiceAppleDeviceCreateHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedAppleDeviceCreateServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedAppleDeviceCreateServiceHandler struct{}

func (UnimplementedAppleDeviceCreateServiceHandler) AppleDeviceCreate(context.Context, *connect.Request[v1.AppleDeviceCreateRequest]) (*connect.Response[v1.AppleDeviceCreateResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("proto.webauthn.v1.AppleDeviceCreateService.AppleDeviceCreate is not implemented"))
}

// AppleDeviceAssertServiceClient is a client for the proto.webauthn.v1.AppleDeviceAssertService
// service.
type AppleDeviceAssertServiceClient interface {
	AppleDeviceAssert(context.Context, *connect.Request[v1.AppleDeviceAssertRequest]) (*connect.Response[v1.AppleDeviceAssertResponse], error)
}

// NewAppleDeviceAssertServiceClient constructs a client for the
// proto.webauthn.v1.AppleDeviceAssertService service. By default, it uses the Connect protocol with
// the binary Protobuf Codec, asks for gzipped responses, and sends uncompressed requests. To use
// the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewAppleDeviceAssertServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) AppleDeviceAssertServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &appleDeviceAssertServiceClient{
		appleDeviceAssert: connect.NewClient[v1.AppleDeviceAssertRequest, v1.AppleDeviceAssertResponse](
			httpClient,
			baseURL+AppleDeviceAssertServiceAppleDeviceAssertProcedure,
			connect.WithSchema(appleDeviceAssertServiceAppleDeviceAssertMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// appleDeviceAssertServiceClient implements AppleDeviceAssertServiceClient.
type appleDeviceAssertServiceClient struct {
	appleDeviceAssert *connect.Client[v1.AppleDeviceAssertRequest, v1.AppleDeviceAssertResponse]
}

// AppleDeviceAssert calls proto.webauthn.v1.AppleDeviceAssertService.AppleDeviceAssert.
func (c *appleDeviceAssertServiceClient) AppleDeviceAssert(ctx context.Context, req *connect.Request[v1.AppleDeviceAssertRequest]) (*connect.Response[v1.AppleDeviceAssertResponse], error) {
	return c.appleDeviceAssert.CallUnary(ctx, req)
}

// AppleDeviceAssertServiceHandler is an implementation of the
// proto.webauthn.v1.AppleDeviceAssertService service.
type AppleDeviceAssertServiceHandler interface {
	AppleDeviceAssert(context.Context, *connect.Request[v1.AppleDeviceAssertRequest]) (*connect.Response[v1.AppleDeviceAssertResponse], error)
}

// NewAppleDeviceAssertServiceHandler builds an HTTP handler from the service implementation. It
// returns the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewAppleDeviceAssertServiceHandler(svc AppleDeviceAssertServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	appleDeviceAssertServiceAppleDeviceAssertHandler := connect.NewUnaryHandler(
		AppleDeviceAssertServiceAppleDeviceAssertProcedure,
		svc.AppleDeviceAssert,
		connect.WithSchema(appleDeviceAssertServiceAppleDeviceAssertMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/proto.webauthn.v1.AppleDeviceAssertService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case AppleDeviceAssertServiceAppleDeviceAssertProcedure:
			appleDeviceAssertServiceAppleDeviceAssertHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedAppleDeviceAssertServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedAppleDeviceAssertServiceHandler struct{}

func (UnimplementedAppleDeviceAssertServiceHandler) AppleDeviceAssert(context.Context, *connect.Request[v1.AppleDeviceAssertRequest]) (*connect.Response[v1.AppleDeviceAssertResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("proto.webauthn.v1.AppleDeviceAssertService.AppleDeviceAssert is not implemented"))
}