package miauth_proto

import (
	reflect "reflect"
	sync "sync"

	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// Points are represented as latitude-longitude pairs in the E7 representation
// (degrees multiplied by 10**7 and rounded to the nearest integer).
// Latitudes should be in the range +/- 90 degrees and longitude should be in
// the range +/- 180 degrees (inclusive).
type ValidationInput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccessToken string `protobuf:"bytes,1,opt,name=AccessToken,proto3" json:"AccessToken,omitempty"`
}

func (x *ValidationInput) Reset() {
	*x = ValidationInput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_demo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ValidationInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidationInput) ProtoMessage() {}

func (x *ValidationInput) ProtoReflect() protoreflect.Message {
	mi := &file_demo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidationInput.ProtoReflect.Descriptor instead.
func (*ValidationInput) Descriptor() ([]byte, []int) {
	return file_demo_proto_rawDescGZIP(), []int{0}
}

func (x *ValidationInput) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

type ValidationResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsOk         bool    `protobuf:"varint,1,opt,name=isOk,proto3" json:"isOk,omitempty"`
	UserEmail    string  `protobuf:"bytes,2,opt,name=userEmail,proto3" json:"userEmail,omitempty"`
	UserMiauthID string  `protobuf:"bytes,3,opt,name=userMiauthID,proto3" json:"userMiauthID,omitempty"`
	Exp          float64 `protobuf:"fixed64,4,opt,name=exp,proto3" json:"exp,omitempty"`
}

func (x *ValidationResult) Reset() {
	*x = ValidationResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_demo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ValidationResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidationResult) ProtoMessage() {}

func (x *ValidationResult) ProtoReflect() protoreflect.Message {
	mi := &file_demo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidationResult.ProtoReflect.Descriptor instead.
func (*ValidationResult) Descriptor() ([]byte, []int) {
	return file_demo_proto_rawDescGZIP(), []int{1}
}

func (x *ValidationResult) GetIsOk() bool {
	if x != nil {
		return x.IsOk
	}
	return false
}

func (x *ValidationResult) GetUserEmail() string {
	if x != nil {
		return x.UserEmail
	}
	return ""
}

func (x *ValidationResult) GetUserMiauthID() string {
	if x != nil {
		return x.UserMiauthID
	}
	return ""
}

func (x *ValidationResult) GetExp() float64 {
	if x != nil {
		return x.Exp
	}
	return 0
}

var File_demo_proto protoreflect.FileDescriptor

var file_demo_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x5f, 0x77, 0x69, 0x74, 0x68, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x22, 0x33,
	0x0a, 0x0f, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x70, 0x75,
	0x74, 0x12, 0x20, 0x0a, 0x0b, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x22, 0x7a, 0x0a, 0x10, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x69, 0x73, 0x4f, 0x6b, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x69, 0x73, 0x4f, 0x6b, 0x12, 0x1c, 0x0a, 0x09, 0x75,
	0x73, 0x65, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x75, 0x73, 0x65, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x22, 0x0a, 0x0c, 0x75, 0x73, 0x65,
	0x72, 0x4d, 0x69, 0x61, 0x75, 0x74, 0x68, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x75, 0x73, 0x65, 0x72, 0x4d, 0x69, 0x61, 0x75, 0x74, 0x68, 0x49, 0x44, 0x12, 0x10, 0x0a,
	0x03, 0x65, 0x78, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x65, 0x78, 0x70, 0x32,
	0x60, 0x0a, 0x06, 0x4d, 0x69, 0x41, 0x75, 0x74, 0x68, 0x12, 0x56, 0x0a, 0x0b, 0x56, 0x65, 0x72,
	0x69, 0x66, 0x79, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x21, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x5f, 0x77, 0x69, 0x74, 0x68, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x56, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x1a, 0x22, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x5f, 0x77, 0x69, 0x74, 0x68, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x56,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22,
	0x00, 0x42, 0x19, 0x5a, 0x17, 0x61, 0x6e, 0x61, 0x72, 0x67, 0x75, 0x2f, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x5f, 0x77, 0x69, 0x74, 0x68, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_demo_proto_rawDescOnce sync.Once
	file_demo_proto_rawDescData = file_demo_proto_rawDesc
)

func file_demo_proto_rawDescGZIP() []byte {
	file_demo_proto_rawDescOnce.Do(func() {
		file_demo_proto_rawDescData = protoimpl.X.CompressGZIP(file_demo_proto_rawDescData)
	})
	return file_demo_proto_rawDescData
}

var file_demo_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_demo_proto_goTypes = []interface{}{
	(*ValidationInput)(nil),  // 0: server_with_grpc.ValidationInput
	(*ValidationResult)(nil), // 1: server_with_grpc.ValidationResult
}
var file_demo_proto_depIdxs = []int32{
	0, // 0: server_with_grpc.MiAuth.VerifyToken:input_type -> server_with_grpc.ValidationInput
	1, // 1: server_with_grpc.MiAuth.VerifyToken:output_type -> server_with_grpc.ValidationResult
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_demo_proto_init() }
func file_demo_proto_init() {
	if File_demo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_demo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ValidationInput); i {
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
		file_demo_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ValidationResult); i {
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
			RawDescriptor: file_demo_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_demo_proto_goTypes,
		DependencyIndexes: file_demo_proto_depIdxs,
		MessageInfos:      file_demo_proto_msgTypes,
	}.Build()
	File_demo_proto = out.File
	file_demo_proto_rawDesc = nil
	file_demo_proto_goTypes = nil
	file_demo_proto_depIdxs = nil
}
