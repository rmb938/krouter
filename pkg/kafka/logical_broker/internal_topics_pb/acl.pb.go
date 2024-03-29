// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: acl.proto

package internal_topics_pb

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

type ACLMessageKey_Operation int32

const (
	ACLMessageKey_OPERATION_RESERVED         ACLMessageKey_Operation = 0
	ACLMessageKey_OPERATION_ALL              ACLMessageKey_Operation = 2
	ACLMessageKey_OPERATION_READ             ACLMessageKey_Operation = 3
	ACLMessageKey_OPERATION_WRITE            ACLMessageKey_Operation = 4
	ACLMessageKey_OPERATION_CREATE           ACLMessageKey_Operation = 5
	ACLMessageKey_OPERATION_DELETE           ACLMessageKey_Operation = 6
	ACLMessageKey_OPERATION_ALTER            ACLMessageKey_Operation = 7
	ACLMessageKey_OPERATION_DESCRIBE         ACLMessageKey_Operation = 8
	ACLMessageKey_OPERATION_CLUSTER_ACTION   ACLMessageKey_Operation = 9
	ACLMessageKey_OPERATION_DESCRIBE_CONFIGS ACLMessageKey_Operation = 10
	ACLMessageKey_OPERATION_ALTER_CONFIGS    ACLMessageKey_Operation = 11
	ACLMessageKey_OPERATION_IDEMPOTENT_WRITE ACLMessageKey_Operation = 12
)

// Enum value maps for ACLMessageKey_Operation.
var (
	ACLMessageKey_Operation_name = map[int32]string{
		0:  "OPERATION_RESERVED",
		2:  "OPERATION_ALL",
		3:  "OPERATION_READ",
		4:  "OPERATION_WRITE",
		5:  "OPERATION_CREATE",
		6:  "OPERATION_DELETE",
		7:  "OPERATION_ALTER",
		8:  "OPERATION_DESCRIBE",
		9:  "OPERATION_CLUSTER_ACTION",
		10: "OPERATION_DESCRIBE_CONFIGS",
		11: "OPERATION_ALTER_CONFIGS",
		12: "OPERATION_IDEMPOTENT_WRITE",
	}
	ACLMessageKey_Operation_value = map[string]int32{
		"OPERATION_RESERVED":         0,
		"OPERATION_ALL":              2,
		"OPERATION_READ":             3,
		"OPERATION_WRITE":            4,
		"OPERATION_CREATE":           5,
		"OPERATION_DELETE":           6,
		"OPERATION_ALTER":            7,
		"OPERATION_DESCRIBE":         8,
		"OPERATION_CLUSTER_ACTION":   9,
		"OPERATION_DESCRIBE_CONFIGS": 10,
		"OPERATION_ALTER_CONFIGS":    11,
		"OPERATION_IDEMPOTENT_WRITE": 12,
	}
)

func (x ACLMessageKey_Operation) Enum() *ACLMessageKey_Operation {
	p := new(ACLMessageKey_Operation)
	*p = x
	return p
}

func (x ACLMessageKey_Operation) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ACLMessageKey_Operation) Descriptor() protoreflect.EnumDescriptor {
	return file_acl_proto_enumTypes[0].Descriptor()
}

func (ACLMessageKey_Operation) Type() protoreflect.EnumType {
	return &file_acl_proto_enumTypes[0]
}

func (x ACLMessageKey_Operation) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ACLMessageKey_Operation.Descriptor instead.
func (ACLMessageKey_Operation) EnumDescriptor() ([]byte, []int) {
	return file_acl_proto_rawDescGZIP(), []int{0, 0}
}

type ACLMessageKey_ResourceType int32

const (
	ACLMessageKey_RESOURCE_TYPE_RESERVED         ACLMessageKey_ResourceType = 0
	ACLMessageKey_RESOURCE_TYPE_TOPIC            ACLMessageKey_ResourceType = 2
	ACLMessageKey_RESOURCE_TYPE_GROUP            ACLMessageKey_ResourceType = 3
	ACLMessageKey_RESOURCE_TYPE_CLUSTER          ACLMessageKey_ResourceType = 4
	ACLMessageKey_RESOURCE_TYPE_TRANSACTIONAL_ID ACLMessageKey_ResourceType = 5
	ACLMessageKey_RESOURCE_TYPE_DELEGATION_TOKEN ACLMessageKey_ResourceType = 6
)

// Enum value maps for ACLMessageKey_ResourceType.
var (
	ACLMessageKey_ResourceType_name = map[int32]string{
		0: "RESOURCE_TYPE_RESERVED",
		2: "RESOURCE_TYPE_TOPIC",
		3: "RESOURCE_TYPE_GROUP",
		4: "RESOURCE_TYPE_CLUSTER",
		5: "RESOURCE_TYPE_TRANSACTIONAL_ID",
		6: "RESOURCE_TYPE_DELEGATION_TOKEN",
	}
	ACLMessageKey_ResourceType_value = map[string]int32{
		"RESOURCE_TYPE_RESERVED":         0,
		"RESOURCE_TYPE_TOPIC":            2,
		"RESOURCE_TYPE_GROUP":            3,
		"RESOURCE_TYPE_CLUSTER":          4,
		"RESOURCE_TYPE_TRANSACTIONAL_ID": 5,
		"RESOURCE_TYPE_DELEGATION_TOKEN": 6,
	}
)

func (x ACLMessageKey_ResourceType) Enum() *ACLMessageKey_ResourceType {
	p := new(ACLMessageKey_ResourceType)
	*p = x
	return p
}

func (x ACLMessageKey_ResourceType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ACLMessageKey_ResourceType) Descriptor() protoreflect.EnumDescriptor {
	return file_acl_proto_enumTypes[1].Descriptor()
}

func (ACLMessageKey_ResourceType) Type() protoreflect.EnumType {
	return &file_acl_proto_enumTypes[1]
}

func (x ACLMessageKey_ResourceType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ACLMessageKey_ResourceType.Descriptor instead.
func (ACLMessageKey_ResourceType) EnumDescriptor() ([]byte, []int) {
	return file_acl_proto_rawDescGZIP(), []int{0, 1}
}

type ACLMessageKey_PatternType int32

const (
	ACLMessageKey_PATTERN_TYPE_RESERVED ACLMessageKey_PatternType = 0
	ACLMessageKey_PATTERN_TYPE_MATCH    ACLMessageKey_PatternType = 2
	ACLMessageKey_PATTERN_TYPE_LITERAL  ACLMessageKey_PatternType = 3
	ACLMessageKey_PATTERN_TYPE_PREFIXED ACLMessageKey_PatternType = 4
)

// Enum value maps for ACLMessageKey_PatternType.
var (
	ACLMessageKey_PatternType_name = map[int32]string{
		0: "PATTERN_TYPE_RESERVED",
		2: "PATTERN_TYPE_MATCH",
		3: "PATTERN_TYPE_LITERAL",
		4: "PATTERN_TYPE_PREFIXED",
	}
	ACLMessageKey_PatternType_value = map[string]int32{
		"PATTERN_TYPE_RESERVED": 0,
		"PATTERN_TYPE_MATCH":    2,
		"PATTERN_TYPE_LITERAL":  3,
		"PATTERN_TYPE_PREFIXED": 4,
	}
)

func (x ACLMessageKey_PatternType) Enum() *ACLMessageKey_PatternType {
	p := new(ACLMessageKey_PatternType)
	*p = x
	return p
}

func (x ACLMessageKey_PatternType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ACLMessageKey_PatternType) Descriptor() protoreflect.EnumDescriptor {
	return file_acl_proto_enumTypes[2].Descriptor()
}

func (ACLMessageKey_PatternType) Type() protoreflect.EnumType {
	return &file_acl_proto_enumTypes[2]
}

func (x ACLMessageKey_PatternType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ACLMessageKey_PatternType.Descriptor instead.
func (ACLMessageKey_PatternType) EnumDescriptor() ([]byte, []int) {
	return file_acl_proto_rawDescGZIP(), []int{0, 2}
}

type ACLMessageValue_Permission int32

const (
	ACLMessageValue_PERMISSION_RESERVED ACLMessageValue_Permission = 0
	ACLMessageValue_PERMISSION_DENY     ACLMessageValue_Permission = 2
	ACLMessageValue_PERMISSION_ALLOW    ACLMessageValue_Permission = 3
)

// Enum value maps for ACLMessageValue_Permission.
var (
	ACLMessageValue_Permission_name = map[int32]string{
		0: "PERMISSION_RESERVED",
		2: "PERMISSION_DENY",
		3: "PERMISSION_ALLOW",
	}
	ACLMessageValue_Permission_value = map[string]int32{
		"PERMISSION_RESERVED": 0,
		"PERMISSION_DENY":     2,
		"PERMISSION_ALLOW":    3,
	}
)

func (x ACLMessageValue_Permission) Enum() *ACLMessageValue_Permission {
	p := new(ACLMessageValue_Permission)
	*p = x
	return p
}

func (x ACLMessageValue_Permission) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ACLMessageValue_Permission) Descriptor() protoreflect.EnumDescriptor {
	return file_acl_proto_enumTypes[3].Descriptor()
}

func (ACLMessageValue_Permission) Type() protoreflect.EnumType {
	return &file_acl_proto_enumTypes[3]
}

func (x ACLMessageValue_Permission) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ACLMessageValue_Permission.Descriptor instead.
func (ACLMessageValue_Permission) EnumDescriptor() ([]byte, []int) {
	return file_acl_proto_rawDescGZIP(), []int{1, 0}
}

type ACLMessageKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Operation    ACLMessageKey_Operation    `protobuf:"varint,1,opt,name=operation,proto3,enum=proto.ACLMessageKey_Operation" json:"operation,omitempty"`
	ResourceType ACLMessageKey_ResourceType `protobuf:"varint,2,opt,name=resource_type,json=resourceType,proto3,enum=proto.ACLMessageKey_ResourceType" json:"resource_type,omitempty"`
	PatternType  ACLMessageKey_PatternType  `protobuf:"varint,3,opt,name=pattern_type,json=patternType,proto3,enum=proto.ACLMessageKey_PatternType" json:"pattern_type,omitempty"`
	ResourceName string                     `protobuf:"bytes,4,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	Principal    string                     `protobuf:"bytes,5,opt,name=principal,proto3" json:"principal,omitempty"`
}

func (x *ACLMessageKey) Reset() {
	*x = ACLMessageKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_acl_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ACLMessageKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ACLMessageKey) ProtoMessage() {}

func (x *ACLMessageKey) ProtoReflect() protoreflect.Message {
	mi := &file_acl_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ACLMessageKey.ProtoReflect.Descriptor instead.
func (*ACLMessageKey) Descriptor() ([]byte, []int) {
	return file_acl_proto_rawDescGZIP(), []int{0}
}

func (x *ACLMessageKey) GetOperation() ACLMessageKey_Operation {
	if x != nil {
		return x.Operation
	}
	return ACLMessageKey_OPERATION_RESERVED
}

func (x *ACLMessageKey) GetResourceType() ACLMessageKey_ResourceType {
	if x != nil {
		return x.ResourceType
	}
	return ACLMessageKey_RESOURCE_TYPE_RESERVED
}

func (x *ACLMessageKey) GetPatternType() ACLMessageKey_PatternType {
	if x != nil {
		return x.PatternType
	}
	return ACLMessageKey_PATTERN_TYPE_RESERVED
}

func (x *ACLMessageKey) GetResourceName() string {
	if x != nil {
		return x.ResourceName
	}
	return ""
}

func (x *ACLMessageKey) GetPrincipal() string {
	if x != nil {
		return x.Principal
	}
	return ""
}

type ACLMessageValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Permission ACLMessageValue_Permission `protobuf:"varint,1,opt,name=permission,proto3,enum=proto.ACLMessageValue_Permission" json:"permission,omitempty"`
}

func (x *ACLMessageValue) Reset() {
	*x = ACLMessageValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_acl_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ACLMessageValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ACLMessageValue) ProtoMessage() {}

func (x *ACLMessageValue) ProtoReflect() protoreflect.Message {
	mi := &file_acl_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ACLMessageValue.ProtoReflect.Descriptor instead.
func (*ACLMessageValue) Descriptor() ([]byte, []int) {
	return file_acl_proto_rawDescGZIP(), []int{1}
}

func (x *ACLMessageValue) GetPermission() ACLMessageValue_Permission {
	if x != nil {
		return x.Permission
	}
	return ACLMessageValue_PERMISSION_RESERVED
}

var File_acl_proto protoreflect.FileDescriptor

var file_acl_proto_rawDesc = []byte{
	0x0a, 0x09, 0x61, 0x63, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x9e, 0x07, 0x0a, 0x0d, 0x41, 0x43, 0x4c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x4b, 0x65, 0x79, 0x12, 0x3c, 0x0a, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x41, 0x43, 0x4c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4b, 0x65, 0x79, 0x2e, 0x4f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x46, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x21, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x41, 0x43, 0x4c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4b, 0x65, 0x79, 0x2e,
	0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0c, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x43, 0x0a, 0x0c, 0x70, 0x61,
	0x74, 0x74, 0x65, 0x72, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x43, 0x4c, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x4b, 0x65, 0x79, 0x2e, 0x50, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x0b, 0x70, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x72, 0x69, 0x6e, 0x63, 0x69, 0x70, 0x61,
	0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x72, 0x69, 0x6e, 0x63, 0x69, 0x70,
	0x61, 0x6c, 0x22, 0xb9, 0x02, 0x0a, 0x09, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x16, 0x0a, 0x12, 0x4f, 0x50, 0x45, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x52, 0x45,
	0x53, 0x45, 0x52, 0x56, 0x45, 0x44, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x4f, 0x50, 0x45, 0x52,
	0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x41, 0x4c, 0x4c, 0x10, 0x02, 0x12, 0x12, 0x0a, 0x0e, 0x4f,
	0x50, 0x45, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x52, 0x45, 0x41, 0x44, 0x10, 0x03, 0x12,
	0x13, 0x0a, 0x0f, 0x4f, 0x50, 0x45, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x57, 0x52, 0x49,
	0x54, 0x45, 0x10, 0x04, 0x12, 0x14, 0x0a, 0x10, 0x4f, 0x50, 0x45, 0x52, 0x41, 0x54, 0x49, 0x4f,
	0x4e, 0x5f, 0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x10, 0x05, 0x12, 0x14, 0x0a, 0x10, 0x4f, 0x50,
	0x45, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x10, 0x06,
	0x12, 0x13, 0x0a, 0x0f, 0x4f, 0x50, 0x45, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x41, 0x4c,
	0x54, 0x45, 0x52, 0x10, 0x07, 0x12, 0x16, 0x0a, 0x12, 0x4f, 0x50, 0x45, 0x52, 0x41, 0x54, 0x49,
	0x4f, 0x4e, 0x5f, 0x44, 0x45, 0x53, 0x43, 0x52, 0x49, 0x42, 0x45, 0x10, 0x08, 0x12, 0x1c, 0x0a,
	0x18, 0x4f, 0x50, 0x45, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x43, 0x4c, 0x55, 0x53, 0x54,
	0x45, 0x52, 0x5f, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x09, 0x12, 0x1e, 0x0a, 0x1a, 0x4f,
	0x50, 0x45, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x44, 0x45, 0x53, 0x43, 0x52, 0x49, 0x42,
	0x45, 0x5f, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x47, 0x53, 0x10, 0x0a, 0x12, 0x1b, 0x0a, 0x17, 0x4f,
	0x50, 0x45, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x41, 0x4c, 0x54, 0x45, 0x52, 0x5f, 0x43,
	0x4f, 0x4e, 0x46, 0x49, 0x47, 0x53, 0x10, 0x0b, 0x12, 0x1e, 0x0a, 0x1a, 0x4f, 0x50, 0x45, 0x52,
	0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x49, 0x44, 0x45, 0x4d, 0x50, 0x4f, 0x54, 0x45, 0x4e, 0x54,
	0x5f, 0x57, 0x52, 0x49, 0x54, 0x45, 0x10, 0x0c, 0x22, 0x04, 0x08, 0x01, 0x10, 0x01, 0x22, 0xc5,
	0x01, 0x0a, 0x0c, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x1a, 0x0a, 0x16, 0x52, 0x45, 0x53, 0x4f, 0x55, 0x52, 0x43, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x52, 0x45, 0x53, 0x45, 0x52, 0x56, 0x45, 0x44, 0x10, 0x00, 0x12, 0x17, 0x0a, 0x13, 0x52,
	0x45, 0x53, 0x4f, 0x55, 0x52, 0x43, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x54, 0x4f, 0x50,
	0x49, 0x43, 0x10, 0x02, 0x12, 0x17, 0x0a, 0x13, 0x52, 0x45, 0x53, 0x4f, 0x55, 0x52, 0x43, 0x45,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x47, 0x52, 0x4f, 0x55, 0x50, 0x10, 0x03, 0x12, 0x19, 0x0a,
	0x15, 0x52, 0x45, 0x53, 0x4f, 0x55, 0x52, 0x43, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x43,
	0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x10, 0x04, 0x12, 0x22, 0x0a, 0x1e, 0x52, 0x45, 0x53, 0x4f,
	0x55, 0x52, 0x43, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x41,
	0x43, 0x54, 0x49, 0x4f, 0x4e, 0x41, 0x4c, 0x5f, 0x49, 0x44, 0x10, 0x05, 0x12, 0x22, 0x0a, 0x1e,
	0x52, 0x45, 0x53, 0x4f, 0x55, 0x52, 0x43, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x44, 0x45,
	0x4c, 0x45, 0x47, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x54, 0x4f, 0x4b, 0x45, 0x4e, 0x10, 0x06,
	0x22, 0x04, 0x08, 0x01, 0x10, 0x01, 0x22, 0x7b, 0x0a, 0x0b, 0x50, 0x61, 0x74, 0x74, 0x65, 0x72,
	0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x19, 0x0a, 0x15, 0x50, 0x41, 0x54, 0x54, 0x45, 0x52, 0x4e,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x52, 0x45, 0x53, 0x45, 0x52, 0x56, 0x45, 0x44, 0x10, 0x00,
	0x12, 0x16, 0x0a, 0x12, 0x50, 0x41, 0x54, 0x54, 0x45, 0x52, 0x4e, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x4d, 0x41, 0x54, 0x43, 0x48, 0x10, 0x02, 0x12, 0x18, 0x0a, 0x14, 0x50, 0x41, 0x54, 0x54,
	0x45, 0x52, 0x4e, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4c, 0x49, 0x54, 0x45, 0x52, 0x41, 0x4c,
	0x10, 0x03, 0x12, 0x19, 0x0a, 0x15, 0x50, 0x41, 0x54, 0x54, 0x45, 0x52, 0x4e, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x50, 0x52, 0x45, 0x46, 0x49, 0x58, 0x45, 0x44, 0x10, 0x04, 0x22, 0x04, 0x08,
	0x01, 0x10, 0x01, 0x22, 0xac, 0x01, 0x0a, 0x0f, 0x41, 0x43, 0x4c, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x41, 0x0a, 0x0a, 0x70, 0x65, 0x72, 0x6d, 0x69,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x21, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x43, 0x4c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x2e, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x0a,
	0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x56, 0x0a, 0x0a, 0x50, 0x65,
	0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x17, 0x0a, 0x13, 0x50, 0x45, 0x52, 0x4d,
	0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x52, 0x45, 0x53, 0x45, 0x52, 0x56, 0x45, 0x44, 0x10,
	0x00, 0x12, 0x13, 0x0a, 0x0f, 0x50, 0x45, 0x52, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f,
	0x44, 0x45, 0x4e, 0x59, 0x10, 0x02, 0x12, 0x14, 0x0a, 0x10, 0x50, 0x45, 0x52, 0x4d, 0x49, 0x53,
	0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x41, 0x4c, 0x4c, 0x4f, 0x57, 0x10, 0x03, 0x22, 0x04, 0x08, 0x01,
	0x10, 0x01, 0x42, 0x16, 0x5a, 0x14, 0x2e, 0x3b, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x5f, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x73, 0x5f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_acl_proto_rawDescOnce sync.Once
	file_acl_proto_rawDescData = file_acl_proto_rawDesc
)

func file_acl_proto_rawDescGZIP() []byte {
	file_acl_proto_rawDescOnce.Do(func() {
		file_acl_proto_rawDescData = protoimpl.X.CompressGZIP(file_acl_proto_rawDescData)
	})
	return file_acl_proto_rawDescData
}

var file_acl_proto_enumTypes = make([]protoimpl.EnumInfo, 4)
var file_acl_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_acl_proto_goTypes = []interface{}{
	(ACLMessageKey_Operation)(0),    // 0: proto.ACLMessageKey.Operation
	(ACLMessageKey_ResourceType)(0), // 1: proto.ACLMessageKey.ResourceType
	(ACLMessageKey_PatternType)(0),  // 2: proto.ACLMessageKey.PatternType
	(ACLMessageValue_Permission)(0), // 3: proto.ACLMessageValue.Permission
	(*ACLMessageKey)(nil),           // 4: proto.ACLMessageKey
	(*ACLMessageValue)(nil),         // 5: proto.ACLMessageValue
}
var file_acl_proto_depIdxs = []int32{
	0, // 0: proto.ACLMessageKey.operation:type_name -> proto.ACLMessageKey.Operation
	1, // 1: proto.ACLMessageKey.resource_type:type_name -> proto.ACLMessageKey.ResourceType
	2, // 2: proto.ACLMessageKey.pattern_type:type_name -> proto.ACLMessageKey.PatternType
	3, // 3: proto.ACLMessageValue.permission:type_name -> proto.ACLMessageValue.Permission
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_acl_proto_init() }
func file_acl_proto_init() {
	if File_acl_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_acl_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ACLMessageKey); i {
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
		file_acl_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ACLMessageValue); i {
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
			RawDescriptor: file_acl_proto_rawDesc,
			NumEnums:      4,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_acl_proto_goTypes,
		DependencyIndexes: file_acl_proto_depIdxs,
		EnumInfos:         file_acl_proto_enumTypes,
		MessageInfos:      file_acl_proto_msgTypes,
	}.Build()
	File_acl_proto = out.File
	file_acl_proto_rawDesc = nil
	file_acl_proto_goTypes = nil
	file_acl_proto_depIdxs = nil
}
