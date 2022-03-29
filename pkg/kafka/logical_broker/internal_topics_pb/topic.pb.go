// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: topic.proto

package internal_topics_pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TopicMessageValue_Action int32

const (
	TopicMessageValue_ACTION_CREATE TopicMessageValue_Action = 0
	TopicMessageValue_ACTION_UPDATE TopicMessageValue_Action = 1
	TopicMessageValue_ACTION_DELETE TopicMessageValue_Action = 2
)

// Enum value maps for TopicMessageValue_Action.
var (
	TopicMessageValue_Action_name = map[int32]string{
		0: "ACTION_CREATE",
		1: "ACTION_UPDATE",
		2: "ACTION_DELETE",
	}
	TopicMessageValue_Action_value = map[string]int32{
		"ACTION_CREATE": 0,
		"ACTION_UPDATE": 1,
		"ACTION_DELETE": 2,
	}
)

func (x TopicMessageValue_Action) Enum() *TopicMessageValue_Action {
	p := new(TopicMessageValue_Action)
	*p = x
	return p
}

func (x TopicMessageValue_Action) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TopicMessageValue_Action) Descriptor() protoreflect.EnumDescriptor {
	return file_topic_proto_enumTypes[0].Descriptor()
}

func (TopicMessageValue_Action) Type() protoreflect.EnumType {
	return &file_topic_proto_enumTypes[0]
}

func (x TopicMessageValue_Action) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TopicMessageValue_Action.Descriptor instead.
func (TopicMessageValue_Action) EnumDescriptor() ([]byte, []int) {
	return file_topic_proto_rawDescGZIP(), []int{1, 0}
}

type Topic struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name       string                             `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Partitions int32                              `protobuf:"varint,2,opt,name=partitions,proto3" json:"partitions,omitempty"`
	Config     map[string]*wrapperspb.StringValue `protobuf:"bytes,3,rep,name=config,proto3" json:"config,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Topic) Reset() {
	*x = Topic{}
	if protoimpl.UnsafeEnabled {
		mi := &file_topic_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Topic) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Topic) ProtoMessage() {}

func (x *Topic) ProtoReflect() protoreflect.Message {
	mi := &file_topic_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Topic.ProtoReflect.Descriptor instead.
func (*Topic) Descriptor() ([]byte, []int) {
	return file_topic_proto_rawDescGZIP(), []int{0}
}

func (x *Topic) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Topic) GetPartitions() int32 {
	if x != nil {
		return x.Partitions
	}
	return 0
}

func (x *Topic) GetConfig() map[string]*wrapperspb.StringValue {
	if x != nil {
		return x.Config
	}
	return nil
}

type TopicMessageValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string                   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Action  TopicMessageValue_Action `protobuf:"varint,2,opt,name=action,proto3,enum=proto.TopicMessageValue_Action" json:"action,omitempty"`
	Cluster string                   `protobuf:"bytes,3,opt,name=cluster,proto3" json:"cluster,omitempty"`
	Topic   *Topic                   `protobuf:"bytes,4,opt,name=topic,proto3" json:"topic,omitempty"`
}

func (x *TopicMessageValue) Reset() {
	*x = TopicMessageValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_topic_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TopicMessageValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TopicMessageValue) ProtoMessage() {}

func (x *TopicMessageValue) ProtoReflect() protoreflect.Message {
	mi := &file_topic_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TopicMessageValue.ProtoReflect.Descriptor instead.
func (*TopicMessageValue) Descriptor() ([]byte, []int) {
	return file_topic_proto_rawDescGZIP(), []int{1}
}

func (x *TopicMessageValue) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *TopicMessageValue) GetAction() TopicMessageValue_Action {
	if x != nil {
		return x.Action
	}
	return TopicMessageValue_ACTION_CREATE
}

func (x *TopicMessageValue) GetCluster() string {
	if x != nil {
		return x.Cluster
	}
	return ""
}

func (x *TopicMessageValue) GetTopic() *Topic {
	if x != nil {
		return x.Topic
	}
	return nil
}

var File_topic_proto protoreflect.FileDescriptor

var file_topic_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc6, 0x01, 0x0a, 0x05, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x12, 0x30, 0x0a, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x6f, 0x70, 0x69, 0x63,
	0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x1a, 0x57, 0x0a, 0x0b, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x32, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xe1, 0x01,
	0x0a, 0x11, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x37, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x54, 0x6f, 0x70, 0x69, 0x63, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x18, 0x0a, 0x07, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x22, 0x0a, 0x05, 0x74, 0x6f,
	0x70, 0x69, 0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x22, 0x41,
	0x0a, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x11, 0x0a, 0x0d, 0x41, 0x43, 0x54, 0x49,
	0x4f, 0x4e, 0x5f, 0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x41,
	0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x10, 0x01, 0x12, 0x11,
	0x0a, 0x0d, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x10,
	0x02, 0x42, 0x16, 0x5a, 0x14, 0x2e, 0x3b, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f,
	0x74, 0x6f, 0x70, 0x69, 0x63, 0x73, 0x5f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_topic_proto_rawDescOnce sync.Once
	file_topic_proto_rawDescData = file_topic_proto_rawDesc
)

func file_topic_proto_rawDescGZIP() []byte {
	file_topic_proto_rawDescOnce.Do(func() {
		file_topic_proto_rawDescData = protoimpl.X.CompressGZIP(file_topic_proto_rawDescData)
	})
	return file_topic_proto_rawDescData
}

var file_topic_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_topic_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_topic_proto_goTypes = []interface{}{
	(TopicMessageValue_Action)(0),  // 0: proto.TopicMessageValue.Action
	(*Topic)(nil),                  // 1: proto.Topic
	(*TopicMessageValue)(nil),      // 2: proto.TopicMessageValue
	nil,                            // 3: proto.Topic.ConfigEntry
	(*wrapperspb.StringValue)(nil), // 4: google.protobuf.StringValue
}
var file_topic_proto_depIdxs = []int32{
	3, // 0: proto.Topic.config:type_name -> proto.Topic.ConfigEntry
	0, // 1: proto.TopicMessageValue.action:type_name -> proto.TopicMessageValue.Action
	1, // 2: proto.TopicMessageValue.topic:type_name -> proto.Topic
	4, // 3: proto.Topic.ConfigEntry.value:type_name -> google.protobuf.StringValue
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_topic_proto_init() }
func file_topic_proto_init() {
	if File_topic_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_topic_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Topic); i {
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
		file_topic_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TopicMessageValue); i {
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
			RawDescriptor: file_topic_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_topic_proto_goTypes,
		DependencyIndexes: file_topic_proto_depIdxs,
		EnumInfos:         file_topic_proto_enumTypes,
		MessageInfos:      file_topic_proto_msgTypes,
	}.Build()
	File_topic_proto = out.File
	file_topic_proto_rawDesc = nil
	file_topic_proto_goTypes = nil
	file_topic_proto_depIdxs = nil
}
