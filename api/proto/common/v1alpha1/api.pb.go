// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        (unknown)
// source: common/v1alpha1/api.proto

package commonv1alpha1

import (
	reflect "reflect"
	sync "sync"

	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RequestResult int32

const (
	RequestResult_REQUEST_RESULT_UNSPECIFIED RequestResult = 0
	RequestResult_REQUEST_RESULT_SCHEDULED   RequestResult = 1
	RequestResult_REQUEST_RESULT_SUCCESS     RequestResult = 2
	RequestResult_REQUEST_RESULT_FAILURE     RequestResult = 3
)

// Enum value maps for RequestResult.
var (
	RequestResult_name = map[int32]string{
		0: "REQUEST_RESULT_UNSPECIFIED",
		1: "REQUEST_RESULT_SCHEDULED",
		2: "REQUEST_RESULT_SUCCESS",
		3: "REQUEST_RESULT_FAILURE",
	}
	RequestResult_value = map[string]int32{
		"REQUEST_RESULT_UNSPECIFIED": 0,
		"REQUEST_RESULT_SCHEDULED":   1,
		"REQUEST_RESULT_SUCCESS":     2,
		"REQUEST_RESULT_FAILURE":     3,
	}
)

func (x RequestResult) Enum() *RequestResult {
	p := new(RequestResult)
	*p = x
	return p
}

func (x RequestResult) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RequestResult) Descriptor() protoreflect.EnumDescriptor {
	return file_common_v1alpha1_api_proto_enumTypes[0].Descriptor()
}

func (RequestResult) Type() protoreflect.EnumType {
	return &file_common_v1alpha1_api_proto_enumTypes[0]
}

func (x RequestResult) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RequestResult.Descriptor instead.
func (RequestResult) EnumDescriptor() ([]byte, []int) {
	return file_common_v1alpha1_api_proto_rawDescGZIP(), []int{0}
}

type ScanResult int32

const (
	ScanResult_SCAN_RESULT_UNSPECIFIED ScanResult = 0
	ScanResult_SCAN_RESULT_SUCCESS     ScanResult = 1
	ScanResult_SCAN_RESULT_FAILURE     ScanResult = 2
)

// Enum value maps for ScanResult.
var (
	ScanResult_name = map[int32]string{
		0: "SCAN_RESULT_UNSPECIFIED",
		1: "SCAN_RESULT_SUCCESS",
		2: "SCAN_RESULT_FAILURE",
	}
	ScanResult_value = map[string]int32{
		"SCAN_RESULT_UNSPECIFIED": 0,
		"SCAN_RESULT_SUCCESS":     1,
		"SCAN_RESULT_FAILURE":     2,
	}
)

func (x ScanResult) Enum() *ScanResult {
	p := new(ScanResult)
	*p = x
	return p
}

func (x ScanResult) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ScanResult) Descriptor() protoreflect.EnumDescriptor {
	return file_common_v1alpha1_api_proto_enumTypes[1].Descriptor()
}

func (ScanResult) Type() protoreflect.EnumType {
	return &file_common_v1alpha1_api_proto_enumTypes[1]
}

func (x ScanResult) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ScanResult.Descriptor instead.
func (ScanResult) EnumDescriptor() ([]byte, []int) {
	return file_common_v1alpha1_api_proto_rawDescGZIP(), []int{1}
}

type PackageVersion struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *PackageVersion) Reset() {
	*x = PackageVersion{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_v1alpha1_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PackageVersion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PackageVersion) ProtoMessage() {}

func (x *PackageVersion) ProtoReflect() protoreflect.Message {
	mi := &file_common_v1alpha1_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PackageVersion.ProtoReflect.Descriptor instead.
func (*PackageVersion) Descriptor() ([]byte, []int) {
	return file_common_v1alpha1_api_proto_rawDescGZIP(), []int{0}
}

func (x *PackageVersion) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PackageVersion) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

type Condition struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type               string        `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Status             string        `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	Reason             string        `protobuf:"bytes,3,opt,name=reason,proto3" json:"reason,omitempty"`
	Message            string        `protobuf:"bytes,4,opt,name=message,proto3" json:"message,omitempty"`
	ObservedGeneration int64         `protobuf:"varint,5,opt,name=observed_generation,json=observedGeneration,proto3" json:"observed_generation,omitempty"`
	LastTransitionTime *v1.Timestamp `protobuf:"bytes,6,opt,name=last_transition_time,json=lastTransitionTime,proto3" json:"last_transition_time,omitempty"`
}

func (x *Condition) Reset() {
	*x = Condition{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_v1alpha1_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Condition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Condition) ProtoMessage() {}

func (x *Condition) ProtoReflect() protoreflect.Message {
	mi := &file_common_v1alpha1_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Condition.ProtoReflect.Descriptor instead.
func (*Condition) Descriptor() ([]byte, []int) {
	return file_common_v1alpha1_api_proto_rawDescGZIP(), []int{1}
}

func (x *Condition) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Condition) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Condition) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

func (x *Condition) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *Condition) GetObservedGeneration() int64 {
	if x != nil {
		return x.ObservedGeneration
	}
	return 0
}

func (x *Condition) GetLastTransitionTime() *v1.Timestamp {
	if x != nil {
		return x.LastTransitionTime
	}
	return nil
}

var File_common_v1alpha1_api_proto protoreflect.FileDescriptor

var file_common_v1alpha1_api_proto_rawDesc = []byte{
	0x0a, 0x19, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x31, 0x2f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x1a, 0x1b, 0x62, 0x75,
	0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x34, 0x6b, 0x38, 0x73, 0x2e, 0x69,
	0x6f, 0x2f, 0x61, 0x70, 0x69, 0x6d, 0x61, 0x63, 0x68, 0x69, 0x6e, 0x65, 0x72, 0x79, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x2f, 0x76, 0x31, 0x2f,
	0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xe2, 0x01, 0x0a, 0x0e, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x3a, 0xa1, 0x01, 0xba, 0x48, 0x9d, 0x01, 0x1a, 0x48, 0x0a, 0x14, 0x70, 0x61, 0x63, 0x6b, 0x61,
	0x67, 0x65, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x6e, 0x61, 0x6d, 0x65, 0x1a,
	0x30, 0x21, 0x68, 0x61, 0x73, 0x28, 0x74, 0x68, 0x69, 0x73, 0x2e, 0x6e, 0x61, 0x6d, 0x65, 0x29,
	0x20, 0x3f, 0x20, 0x27, 0x6e, 0x61, 0x6d, 0x65, 0x20, 0x69, 0x73, 0x20, 0x6d, 0x61, 0x6e, 0x64,
	0x61, 0x74, 0x6f, 0x72, 0x79, 0x20, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x27, 0x20, 0x3a, 0x20, 0x27,
	0x27, 0x1a, 0x51, 0x0a, 0x17, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x5f, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x1a, 0x36, 0x21, 0x68,
	0x61, 0x73, 0x28, 0x74, 0x68, 0x69, 0x73, 0x2e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x29,
	0x20, 0x3f, 0x20, 0x27, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x20, 0x69, 0x73, 0x20, 0x6d,
	0x61, 0x6e, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x79, 0x20, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x27, 0x20,
	0x3a, 0x20, 0x27, 0x27, 0x22, 0xfd, 0x01, 0x0a, 0x09, 0x43, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16,
	0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x2f, 0x0a, 0x13, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x5f, 0x67, 0x65, 0x6e,
	0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x12, 0x6f,
	0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x61, 0x0a, 0x14, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x2f, 0x2e, 0x6b, 0x38, 0x73, 0x2e, 0x69, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x6d, 0x61, 0x63, 0x68,
	0x69, 0x6e, 0x65, 0x72, 0x79, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x6d,
	0x65, 0x74, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x12, 0x6c, 0x61, 0x73, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x54, 0x69, 0x6d, 0x65, 0x2a, 0x85, 0x01, 0x0a, 0x0d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x1e, 0x0a, 0x1a, 0x52, 0x45, 0x51, 0x55, 0x45, 0x53,
	0x54, 0x5f, 0x52, 0x45, 0x53, 0x55, 0x4c, 0x54, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49,
	0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1c, 0x0a, 0x18, 0x52, 0x45, 0x51, 0x55, 0x45, 0x53,
	0x54, 0x5f, 0x52, 0x45, 0x53, 0x55, 0x4c, 0x54, 0x5f, 0x53, 0x43, 0x48, 0x45, 0x44, 0x55, 0x4c,
	0x45, 0x44, 0x10, 0x01, 0x12, 0x1a, 0x0a, 0x16, 0x52, 0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x5f,
	0x52, 0x45, 0x53, 0x55, 0x4c, 0x54, 0x5f, 0x53, 0x55, 0x43, 0x43, 0x45, 0x53, 0x53, 0x10, 0x02,
	0x12, 0x1a, 0x0a, 0x16, 0x52, 0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x5f, 0x52, 0x45, 0x53, 0x55,
	0x4c, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x55, 0x52, 0x45, 0x10, 0x03, 0x2a, 0x5b, 0x0a, 0x0a,
	0x53, 0x63, 0x61, 0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x1b, 0x0a, 0x17, 0x53, 0x43,
	0x41, 0x4e, 0x5f, 0x52, 0x45, 0x53, 0x55, 0x4c, 0x54, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43,
	0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x17, 0x0a, 0x13, 0x53, 0x43, 0x41, 0x4e, 0x5f,
	0x52, 0x45, 0x53, 0x55, 0x4c, 0x54, 0x5f, 0x53, 0x55, 0x43, 0x43, 0x45, 0x53, 0x53, 0x10, 0x01,
	0x12, 0x17, 0x0a, 0x13, 0x53, 0x43, 0x41, 0x4e, 0x5f, 0x52, 0x45, 0x53, 0x55, 0x4c, 0x54, 0x5f,
	0x46, 0x41, 0x49, 0x4c, 0x55, 0x52, 0x45, 0x10, 0x02, 0x42, 0xd0, 0x01, 0x0a, 0x13, 0x63, 0x6f,
	0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x31, 0x42, 0x08, 0x41, 0x70, 0x69, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x52, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x72, 0x6f, 0x6e, 0x63, 0x6f,
	0x72, 0x65, 0x2d, 0x64, 0x65, 0x76, 0x2f, 0x6c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65,
	0x2d, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68,
	0x61, 0x31, 0x3b, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x31, 0xa2, 0x02, 0x03, 0x43, 0x58, 0x58, 0xaa, 0x02, 0x0f, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xca, 0x02, 0x0f, 0x43, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x5c, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xe2, 0x02, 0x1b, 0x43, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x5c, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x5c, 0x47, 0x50,
	0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x10, 0x43, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x3a, 0x3a, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_v1alpha1_api_proto_rawDescOnce sync.Once
	file_common_v1alpha1_api_proto_rawDescData = file_common_v1alpha1_api_proto_rawDesc
)

func file_common_v1alpha1_api_proto_rawDescGZIP() []byte {
	file_common_v1alpha1_api_proto_rawDescOnce.Do(func() {
		file_common_v1alpha1_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_v1alpha1_api_proto_rawDescData)
	})
	return file_common_v1alpha1_api_proto_rawDescData
}

var file_common_v1alpha1_api_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_common_v1alpha1_api_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_common_v1alpha1_api_proto_goTypes = []interface{}{
	(RequestResult)(0),     // 0: common.v1alpha1.RequestResult
	(ScanResult)(0),        // 1: common.v1alpha1.ScanResult
	(*PackageVersion)(nil), // 2: common.v1alpha1.PackageVersion
	(*Condition)(nil),      // 3: common.v1alpha1.Condition
	(*v1.Timestamp)(nil),   // 4: k8s.io.apimachinery.pkg.apis.meta.v1.Timestamp
}
var file_common_v1alpha1_api_proto_depIdxs = []int32{
	4, // 0: common.v1alpha1.Condition.last_transition_time:type_name -> k8s.io.apimachinery.pkg.apis.meta.v1.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_common_v1alpha1_api_proto_init() }
func file_common_v1alpha1_api_proto_init() {
	if File_common_v1alpha1_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_v1alpha1_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PackageVersion); i {
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
		file_common_v1alpha1_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Condition); i {
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
			RawDescriptor: file_common_v1alpha1_api_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_v1alpha1_api_proto_goTypes,
		DependencyIndexes: file_common_v1alpha1_api_proto_depIdxs,
		EnumInfos:         file_common_v1alpha1_api_proto_enumTypes,
		MessageInfos:      file_common_v1alpha1_api_proto_msgTypes,
	}.Build()
	File_common_v1alpha1_api_proto = out.File
	file_common_v1alpha1_api_proto_rawDesc = nil
	file_common_v1alpha1_api_proto_goTypes = nil
	file_common_v1alpha1_api_proto_depIdxs = nil
}
