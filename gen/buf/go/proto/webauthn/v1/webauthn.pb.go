// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: proto/webauthn/v1/webauthn.proto

package webauthn

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CreateChallengeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionId    []byte `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	CeremonyType string `protobuf:"bytes,2,opt,name=ceremony_type,json=ceremonyType,proto3" json:"ceremony_type,omitempty"`
	CredentialId []byte `protobuf:"bytes,3,opt,name=credential_id,json=credentialId,proto3" json:"credential_id,omitempty"`
}

func (x *CreateChallengeRequest) Reset() {
	*x = CreateChallengeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_webauthn_v1_webauthn_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateChallengeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateChallengeRequest) ProtoMessage() {}

func (x *CreateChallengeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_webauthn_v1_webauthn_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateChallengeRequest.ProtoReflect.Descriptor instead.
func (*CreateChallengeRequest) Descriptor() ([]byte, []int) {
	return file_proto_webauthn_v1_webauthn_proto_rawDescGZIP(), []int{0}
}

func (x *CreateChallengeRequest) GetSessionId() []byte {
	if x != nil {
		return x.SessionId
	}
	return nil
}

func (x *CreateChallengeRequest) GetCeremonyType() string {
	if x != nil {
		return x.CeremonyType
	}
	return ""
}

func (x *CreateChallengeRequest) GetCredentialId() []byte {
	if x != nil {
		return x.CredentialId
	}
	return nil
}

type CreateChallengeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Challenge []byte `protobuf:"bytes,1,opt,name=challenge,proto3" json:"challenge,omitempty"`
	Ttl       int64  `protobuf:"varint,2,opt,name=ttl,proto3" json:"ttl,omitempty"`
}

func (x *CreateChallengeResponse) Reset() {
	*x = CreateChallengeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_webauthn_v1_webauthn_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateChallengeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateChallengeResponse) ProtoMessage() {}

func (x *CreateChallengeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_webauthn_v1_webauthn_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateChallengeResponse.ProtoReflect.Descriptor instead.
func (*CreateChallengeResponse) Descriptor() ([]byte, []int) {
	return file_proto_webauthn_v1_webauthn_proto_rawDescGZIP(), []int{1}
}

func (x *CreateChallengeResponse) GetChallenge() []byte {
	if x != nil {
		return x.Challenge
	}
	return nil
}

func (x *CreateChallengeResponse) GetTtl() int64 {
	if x != nil {
		return x.Ttl
	}
	return 0
}

type ApplePasskeyCreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionId         []byte `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	CredentialId      []byte `protobuf:"bytes,2,opt,name=credential_id,json=credentialId,proto3" json:"credential_id,omitempty"`
	ClientDataHash    []byte `protobuf:"bytes,3,opt,name=client_data_hash,json=clientDataHash,proto3" json:"client_data_hash,omitempty"`
	AuthenticatorData []byte `protobuf:"bytes,4,opt,name=authenticator_data,json=authenticatorData,proto3" json:"authenticator_data,omitempty"`
	Signature         []byte `protobuf:"bytes,5,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *ApplePasskeyCreateRequest) Reset() {
	*x = ApplePasskeyCreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_webauthn_v1_webauthn_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApplePasskeyCreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApplePasskeyCreateRequest) ProtoMessage() {}

func (x *ApplePasskeyCreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_webauthn_v1_webauthn_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApplePasskeyCreateRequest.ProtoReflect.Descriptor instead.
func (*ApplePasskeyCreateRequest) Descriptor() ([]byte, []int) {
	return file_proto_webauthn_v1_webauthn_proto_rawDescGZIP(), []int{2}
}

func (x *ApplePasskeyCreateRequest) GetSessionId() []byte {
	if x != nil {
		return x.SessionId
	}
	return nil
}

func (x *ApplePasskeyCreateRequest) GetCredentialId() []byte {
	if x != nil {
		return x.CredentialId
	}
	return nil
}

func (x *ApplePasskeyCreateRequest) GetClientDataHash() []byte {
	if x != nil {
		return x.ClientDataHash
	}
	return nil
}

func (x *ApplePasskeyCreateRequest) GetAuthenticatorData() []byte {
	if x != nil {
		return x.AuthenticatorData
	}
	return nil
}

func (x *ApplePasskeyCreateRequest) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

type ApplePasskeyCreateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *ApplePasskeyCreateResponse) Reset() {
	*x = ApplePasskeyCreateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_webauthn_v1_webauthn_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApplePasskeyCreateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApplePasskeyCreateResponse) ProtoMessage() {}

func (x *ApplePasskeyCreateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_webauthn_v1_webauthn_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApplePasskeyCreateResponse.ProtoReflect.Descriptor instead.
func (*ApplePasskeyCreateResponse) Descriptor() ([]byte, []int) {
	return file_proto_webauthn_v1_webauthn_proto_rawDescGZIP(), []int{3}
}

func (x *ApplePasskeyCreateResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type ApplePasskeyAssertRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionId         []byte `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	CredentialId      []byte `protobuf:"bytes,2,opt,name=credential_id,json=credentialId,proto3" json:"credential_id,omitempty"`
	ClientDataHash    []byte `protobuf:"bytes,3,opt,name=client_data_hash,json=clientDataHash,proto3" json:"client_data_hash,omitempty"`
	AuthenticatorData []byte `protobuf:"bytes,4,opt,name=authenticator_data,json=authenticatorData,proto3" json:"authenticator_data,omitempty"`
	Signature         []byte `protobuf:"bytes,5,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *ApplePasskeyAssertRequest) Reset() {
	*x = ApplePasskeyAssertRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_webauthn_v1_webauthn_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApplePasskeyAssertRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApplePasskeyAssertRequest) ProtoMessage() {}

func (x *ApplePasskeyAssertRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_webauthn_v1_webauthn_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApplePasskeyAssertRequest.ProtoReflect.Descriptor instead.
func (*ApplePasskeyAssertRequest) Descriptor() ([]byte, []int) {
	return file_proto_webauthn_v1_webauthn_proto_rawDescGZIP(), []int{4}
}

func (x *ApplePasskeyAssertRequest) GetSessionId() []byte {
	if x != nil {
		return x.SessionId
	}
	return nil
}

func (x *ApplePasskeyAssertRequest) GetCredentialId() []byte {
	if x != nil {
		return x.CredentialId
	}
	return nil
}

func (x *ApplePasskeyAssertRequest) GetClientDataHash() []byte {
	if x != nil {
		return x.ClientDataHash
	}
	return nil
}

func (x *ApplePasskeyAssertRequest) GetAuthenticatorData() []byte {
	if x != nil {
		return x.AuthenticatorData
	}
	return nil
}

func (x *ApplePasskeyAssertRequest) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

type ApplePasskeyAssertResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *ApplePasskeyAssertResponse) Reset() {
	*x = ApplePasskeyAssertResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_webauthn_v1_webauthn_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApplePasskeyAssertResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApplePasskeyAssertResponse) ProtoMessage() {}

func (x *ApplePasskeyAssertResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_webauthn_v1_webauthn_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApplePasskeyAssertResponse.ProtoReflect.Descriptor instead.
func (*ApplePasskeyAssertResponse) Descriptor() ([]byte, []int) {
	return file_proto_webauthn_v1_webauthn_proto_rawDescGZIP(), []int{5}
}

func (x *ApplePasskeyAssertResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type AppleDeviceCreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionId         []byte `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	CredentialId      []byte `protobuf:"bytes,2,opt,name=credential_id,json=credentialId,proto3" json:"credential_id,omitempty"`
	ClientDataHash    []byte `protobuf:"bytes,3,opt,name=client_data_hash,json=clientDataHash,proto3" json:"client_data_hash,omitempty"`
	AuthenticatorData []byte `protobuf:"bytes,4,opt,name=authenticator_data,json=authenticatorData,proto3" json:"authenticator_data,omitempty"`
	Signature         []byte `protobuf:"bytes,5,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *AppleDeviceCreateRequest) Reset() {
	*x = AppleDeviceCreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_webauthn_v1_webauthn_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppleDeviceCreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppleDeviceCreateRequest) ProtoMessage() {}

func (x *AppleDeviceCreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_webauthn_v1_webauthn_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppleDeviceCreateRequest.ProtoReflect.Descriptor instead.
func (*AppleDeviceCreateRequest) Descriptor() ([]byte, []int) {
	return file_proto_webauthn_v1_webauthn_proto_rawDescGZIP(), []int{6}
}

func (x *AppleDeviceCreateRequest) GetSessionId() []byte {
	if x != nil {
		return x.SessionId
	}
	return nil
}

func (x *AppleDeviceCreateRequest) GetCredentialId() []byte {
	if x != nil {
		return x.CredentialId
	}
	return nil
}

func (x *AppleDeviceCreateRequest) GetClientDataHash() []byte {
	if x != nil {
		return x.ClientDataHash
	}
	return nil
}

func (x *AppleDeviceCreateRequest) GetAuthenticatorData() []byte {
	if x != nil {
		return x.AuthenticatorData
	}
	return nil
}

func (x *AppleDeviceCreateRequest) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

type AppleDeviceCreateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *AppleDeviceCreateResponse) Reset() {
	*x = AppleDeviceCreateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_webauthn_v1_webauthn_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppleDeviceCreateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppleDeviceCreateResponse) ProtoMessage() {}

func (x *AppleDeviceCreateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_webauthn_v1_webauthn_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppleDeviceCreateResponse.ProtoReflect.Descriptor instead.
func (*AppleDeviceCreateResponse) Descriptor() ([]byte, []int) {
	return file_proto_webauthn_v1_webauthn_proto_rawDescGZIP(), []int{7}
}

func (x *AppleDeviceCreateResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type AppleDeviceAssertRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionId         []byte `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	CredentialId      []byte `protobuf:"bytes,2,opt,name=credential_id,json=credentialId,proto3" json:"credential_id,omitempty"`
	ClientDataHash    []byte `protobuf:"bytes,3,opt,name=client_data_hash,json=clientDataHash,proto3" json:"client_data_hash,omitempty"`
	AuthenticatorData []byte `protobuf:"bytes,4,opt,name=authenticator_data,json=authenticatorData,proto3" json:"authenticator_data,omitempty"`
	Signature         []byte `protobuf:"bytes,5,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *AppleDeviceAssertRequest) Reset() {
	*x = AppleDeviceAssertRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_webauthn_v1_webauthn_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppleDeviceAssertRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppleDeviceAssertRequest) ProtoMessage() {}

func (x *AppleDeviceAssertRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_webauthn_v1_webauthn_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppleDeviceAssertRequest.ProtoReflect.Descriptor instead.
func (*AppleDeviceAssertRequest) Descriptor() ([]byte, []int) {
	return file_proto_webauthn_v1_webauthn_proto_rawDescGZIP(), []int{8}
}

func (x *AppleDeviceAssertRequest) GetSessionId() []byte {
	if x != nil {
		return x.SessionId
	}
	return nil
}

func (x *AppleDeviceAssertRequest) GetCredentialId() []byte {
	if x != nil {
		return x.CredentialId
	}
	return nil
}

func (x *AppleDeviceAssertRequest) GetClientDataHash() []byte {
	if x != nil {
		return x.ClientDataHash
	}
	return nil
}

func (x *AppleDeviceAssertRequest) GetAuthenticatorData() []byte {
	if x != nil {
		return x.AuthenticatorData
	}
	return nil
}

func (x *AppleDeviceAssertRequest) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

type AppleDeviceAssertResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *AppleDeviceAssertResponse) Reset() {
	*x = AppleDeviceAssertResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_webauthn_v1_webauthn_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppleDeviceAssertResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppleDeviceAssertResponse) ProtoMessage() {}

func (x *AppleDeviceAssertResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_webauthn_v1_webauthn_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppleDeviceAssertResponse.ProtoReflect.Descriptor instead.
func (*AppleDeviceAssertResponse) Descriptor() ([]byte, []int) {
	return file_proto_webauthn_v1_webauthn_proto_rawDescGZIP(), []int{9}
}

func (x *AppleDeviceAssertResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

var File_proto_webauthn_v1_webauthn_proto protoreflect.FileDescriptor

var file_proto_webauthn_v1_webauthn_proto_rawDesc = []byte{
	0x0a, 0x20, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77, 0x65, 0x62, 0x61, 0x75, 0x74, 0x68, 0x6e,
	0x2f, 0x76, 0x31, 0x2f, 0x77, 0x65, 0x62, 0x61, 0x75, 0x74, 0x68, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x77, 0x65, 0x62, 0x61, 0x75, 0x74,
	0x68, 0x6e, 0x2e, 0x76, 0x31, 0x22, 0x81, 0x01, 0x0a, 0x16, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12,
	0x23, 0x0a, 0x0d, 0x63, 0x65, 0x72, 0x65, 0x6d, 0x6f, 0x6e, 0x79, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x65, 0x72, 0x65, 0x6d, 0x6f, 0x6e, 0x79,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x61, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x63, 0x72, 0x65,
	0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x49, 0x64, 0x22, 0x49, 0x0a, 0x17, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e,
	0x67, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x74, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x03, 0x74, 0x74, 0x6c, 0x22, 0xd6, 0x01, 0x0a, 0x19, 0x41, 0x70, 0x70, 0x6c, 0x65, 0x50, 0x61,
	0x73, 0x73, 0x6b, 0x65, 0x79, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49,
	0x64, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e,
	0x74, 0x69, 0x61, 0x6c, 0x49, 0x64, 0x12, 0x28, 0x0a, 0x10, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x0e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61, 0x48, 0x61, 0x73, 0x68,
	0x12, 0x2d, 0x0a, 0x12, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x6f,
	0x72, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x11, 0x61, 0x75,
	0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x44, 0x61, 0x74, 0x61, 0x12,
	0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x36, 0x0a,
	0x1a, 0x41, 0x70, 0x70, 0x6c, 0x65, 0x50, 0x61, 0x73, 0x73, 0x6b, 0x65, 0x79, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0xd6, 0x01, 0x0a, 0x19, 0x41, 0x70, 0x70, 0x6c, 0x65, 0x50,
	0x61, 0x73, 0x73, 0x6b, 0x65, 0x79, 0x41, 0x73, 0x73, 0x65, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c,
	0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x63, 0x72, 0x65, 0x64, 0x65,
	0x6e, 0x74, 0x69, 0x61, 0x6c, 0x49, 0x64, 0x12, 0x28, 0x0a, 0x10, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x0e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61, 0x48, 0x61, 0x73,
	0x68, 0x12, 0x2d, 0x0a, 0x12, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74,
	0x6f, 0x72, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x11, 0x61,
	0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x36,
	0x0a, 0x1a, 0x41, 0x70, 0x70, 0x6c, 0x65, 0x50, 0x61, 0x73, 0x73, 0x6b, 0x65, 0x79, 0x41, 0x73,
	0x73, 0x65, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0xd5, 0x01, 0x0a, 0x18, 0x41, 0x70, 0x70, 0x6c, 0x65,
	0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c,
	0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x63, 0x72, 0x65, 0x64, 0x65,
	0x6e, 0x74, 0x69, 0x61, 0x6c, 0x49, 0x64, 0x12, 0x28, 0x0a, 0x10, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x0e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61, 0x48, 0x61, 0x73,
	0x68, 0x12, 0x2d, 0x0a, 0x12, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74,
	0x6f, 0x72, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x11, 0x61,
	0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x35,
	0x0a, 0x19, 0x41, 0x70, 0x70, 0x6c, 0x65, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0xd5, 0x01, 0x0a, 0x18, 0x41, 0x70, 0x70, 0x6c, 0x65, 0x44,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x41, 0x73, 0x73, 0x65, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49,
	0x64, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e,
	0x74, 0x69, 0x61, 0x6c, 0x49, 0x64, 0x12, 0x28, 0x0a, 0x10, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x0e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61, 0x48, 0x61, 0x73, 0x68,
	0x12, 0x2d, 0x0a, 0x12, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x6f,
	0x72, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x11, 0x61, 0x75,
	0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x44, 0x61, 0x74, 0x61, 0x12,
	0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x35, 0x0a,
	0x19, 0x41, 0x70, 0x70, 0x6c, 0x65, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x41, 0x73, 0x73, 0x65,
	0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x32, 0x84, 0x01, 0x0a, 0x16, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43,
	0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x6a, 0x0a, 0x0f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e,
	0x67, 0x65, 0x12, 0x29, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x77, 0x65, 0x62, 0x61, 0x75,
	0x74, 0x68, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61,
	0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x77, 0x65, 0x62, 0x61, 0x75, 0x74, 0x68, 0x6e, 0x2e, 0x76,
	0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x32, 0x90, 0x01, 0x0a, 0x19,
	0x41, 0x70, 0x70, 0x6c, 0x65, 0x50, 0x61, 0x73, 0x73, 0x6b, 0x65, 0x79, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x73, 0x0a, 0x12, 0x41, 0x70, 0x70,
	0x6c, 0x65, 0x50, 0x61, 0x73, 0x73, 0x6b, 0x65, 0x79, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12,
	0x2c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x77, 0x65, 0x62, 0x61, 0x75, 0x74, 0x68, 0x6e,
	0x2e, 0x76, 0x31, 0x2e, 0x41, 0x70, 0x70, 0x6c, 0x65, 0x50, 0x61, 0x73, 0x73, 0x6b, 0x65, 0x79,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2d, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x77, 0x65, 0x62, 0x61, 0x75, 0x74, 0x68, 0x6e, 0x2e, 0x76,
	0x31, 0x2e, 0x41, 0x70, 0x70, 0x6c, 0x65, 0x50, 0x61, 0x73, 0x73, 0x6b, 0x65, 0x79, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x32, 0x90,
	0x01, 0x0a, 0x19, 0x41, 0x70, 0x70, 0x6c, 0x65, 0x50, 0x61, 0x73, 0x73, 0x6b, 0x65, 0x79, 0x41,
	0x73, 0x73, 0x65, 0x72, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x73, 0x0a, 0x12,
	0x41, 0x70, 0x70, 0x6c, 0x65, 0x50, 0x61, 0x73, 0x73, 0x6b, 0x65, 0x79, 0x41, 0x73, 0x73, 0x65,
	0x72, 0x74, 0x12, 0x2c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x77, 0x65, 0x62, 0x61, 0x75,
	0x74, 0x68, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x70, 0x70, 0x6c, 0x65, 0x50, 0x61, 0x73, 0x73,
	0x6b, 0x65, 0x79, 0x41, 0x73, 0x73, 0x65, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x2d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x77, 0x65, 0x62, 0x61, 0x75, 0x74, 0x68,
	0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x70, 0x70, 0x6c, 0x65, 0x50, 0x61, 0x73, 0x73, 0x6b, 0x65,
	0x79, 0x41, 0x73, 0x73, 0x65, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x32, 0x8c, 0x01, 0x0a, 0x18, 0x41, 0x70, 0x70, 0x6c, 0x65, 0x44, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x70,
	0x0a, 0x11, 0x41, 0x70, 0x70, 0x6c, 0x65, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x12, 0x2b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x77, 0x65, 0x62, 0x61,
	0x75, 0x74, 0x68, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x70, 0x70, 0x6c, 0x65, 0x44, 0x65, 0x76,
	0x69, 0x63, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x2c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x77, 0x65, 0x62, 0x61, 0x75, 0x74, 0x68,
	0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x70, 0x70, 0x6c, 0x65, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x32, 0x8c, 0x01, 0x0a, 0x18, 0x41, 0x70, 0x70, 0x6c, 0x65, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x41, 0x73, 0x73, 0x65, 0x72, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x70, 0x0a,
	0x11, 0x41, 0x70, 0x70, 0x6c, 0x65, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x41, 0x73, 0x73, 0x65,
	0x72, 0x74, 0x12, 0x2b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x77, 0x65, 0x62, 0x61, 0x75,
	0x74, 0x68, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x70, 0x70, 0x6c, 0x65, 0x44, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x41, 0x73, 0x73, 0x65, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x2c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x77, 0x65, 0x62, 0x61, 0x75, 0x74, 0x68, 0x6e,
	0x2e, 0x76, 0x31, 0x2e, 0x41, 0x70, 0x70, 0x6c, 0x65, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x41,
	0x73, 0x73, 0x65, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42,
	0x42, 0x5a, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x61,
	0x6c, 0x74, 0x65, 0x68, 0x2f, 0x77, 0x65, 0x62, 0x61, 0x75, 0x74, 0x68, 0x6e, 0x2f, 0x67, 0x65,
	0x6e, 0x2f, 0x62, 0x75, 0x66, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77,
	0x65, 0x62, 0x61, 0x75, 0x74, 0x68, 0x6e, 0x2f, 0x76, 0x31, 0x3b, 0x77, 0x65, 0x62, 0x61, 0x75,
	0x74, 0x68, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_webauthn_v1_webauthn_proto_rawDescOnce sync.Once
	file_proto_webauthn_v1_webauthn_proto_rawDescData = file_proto_webauthn_v1_webauthn_proto_rawDesc
)

func file_proto_webauthn_v1_webauthn_proto_rawDescGZIP() []byte {
	file_proto_webauthn_v1_webauthn_proto_rawDescOnce.Do(func() {
		file_proto_webauthn_v1_webauthn_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_webauthn_v1_webauthn_proto_rawDescData)
	})
	return file_proto_webauthn_v1_webauthn_proto_rawDescData
}

var file_proto_webauthn_v1_webauthn_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_proto_webauthn_v1_webauthn_proto_goTypes = []interface{}{
	(*CreateChallengeRequest)(nil),     // 0: proto.webauthn.v1.CreateChallengeRequest
	(*CreateChallengeResponse)(nil),    // 1: proto.webauthn.v1.CreateChallengeResponse
	(*ApplePasskeyCreateRequest)(nil),  // 2: proto.webauthn.v1.ApplePasskeyCreateRequest
	(*ApplePasskeyCreateResponse)(nil), // 3: proto.webauthn.v1.ApplePasskeyCreateResponse
	(*ApplePasskeyAssertRequest)(nil),  // 4: proto.webauthn.v1.ApplePasskeyAssertRequest
	(*ApplePasskeyAssertResponse)(nil), // 5: proto.webauthn.v1.ApplePasskeyAssertResponse
	(*AppleDeviceCreateRequest)(nil),   // 6: proto.webauthn.v1.AppleDeviceCreateRequest
	(*AppleDeviceCreateResponse)(nil),  // 7: proto.webauthn.v1.AppleDeviceCreateResponse
	(*AppleDeviceAssertRequest)(nil),   // 8: proto.webauthn.v1.AppleDeviceAssertRequest
	(*AppleDeviceAssertResponse)(nil),  // 9: proto.webauthn.v1.AppleDeviceAssertResponse
}
var file_proto_webauthn_v1_webauthn_proto_depIdxs = []int32{
	0, // 0: proto.webauthn.v1.CreateChallengeService.CreateChallenge:input_type -> proto.webauthn.v1.CreateChallengeRequest
	2, // 1: proto.webauthn.v1.ApplePasskeyCreateService.ApplePasskeyCreate:input_type -> proto.webauthn.v1.ApplePasskeyCreateRequest
	4, // 2: proto.webauthn.v1.ApplePasskeyAssertService.ApplePasskeyAssert:input_type -> proto.webauthn.v1.ApplePasskeyAssertRequest
	6, // 3: proto.webauthn.v1.AppleDeviceCreateService.AppleDeviceCreate:input_type -> proto.webauthn.v1.AppleDeviceCreateRequest
	8, // 4: proto.webauthn.v1.AppleDeviceAssertService.AppleDeviceAssert:input_type -> proto.webauthn.v1.AppleDeviceAssertRequest
	1, // 5: proto.webauthn.v1.CreateChallengeService.CreateChallenge:output_type -> proto.webauthn.v1.CreateChallengeResponse
	3, // 6: proto.webauthn.v1.ApplePasskeyCreateService.ApplePasskeyCreate:output_type -> proto.webauthn.v1.ApplePasskeyCreateResponse
	5, // 7: proto.webauthn.v1.ApplePasskeyAssertService.ApplePasskeyAssert:output_type -> proto.webauthn.v1.ApplePasskeyAssertResponse
	7, // 8: proto.webauthn.v1.AppleDeviceCreateService.AppleDeviceCreate:output_type -> proto.webauthn.v1.AppleDeviceCreateResponse
	9, // 9: proto.webauthn.v1.AppleDeviceAssertService.AppleDeviceAssert:output_type -> proto.webauthn.v1.AppleDeviceAssertResponse
	5, // [5:10] is the sub-list for method output_type
	0, // [0:5] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_webauthn_v1_webauthn_proto_init() }
func file_proto_webauthn_v1_webauthn_proto_init() {
	if File_proto_webauthn_v1_webauthn_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_webauthn_v1_webauthn_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateChallengeRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_webauthn_v1_webauthn_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateChallengeResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_webauthn_v1_webauthn_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApplePasskeyCreateRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_webauthn_v1_webauthn_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApplePasskeyCreateResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_webauthn_v1_webauthn_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApplePasskeyAssertRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_webauthn_v1_webauthn_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApplePasskeyAssertResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_webauthn_v1_webauthn_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppleDeviceCreateRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_webauthn_v1_webauthn_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppleDeviceCreateResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_webauthn_v1_webauthn_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppleDeviceAssertRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_webauthn_v1_webauthn_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppleDeviceAssertResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_webauthn_v1_webauthn_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   5,
		},
		GoTypes:           file_proto_webauthn_v1_webauthn_proto_goTypes,
		DependencyIndexes: file_proto_webauthn_v1_webauthn_proto_depIdxs,
		MessageInfos:      file_proto_webauthn_v1_webauthn_proto_msgTypes,
	}.Build()
	File_proto_webauthn_v1_webauthn_proto = out.File
	file_proto_webauthn_v1_webauthn_proto_rawDesc = nil
	file_proto_webauthn_v1_webauthn_proto_goTypes = nil
	file_proto_webauthn_v1_webauthn_proto_depIdxs = nil
}