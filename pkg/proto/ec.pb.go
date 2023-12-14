// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.0
// source: ec.proto

package proto

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

type ERROR_CODE int32

const (
	ERROR_CODE_Success ERROR_CODE = 0 //#tag_desc:成功
	// [0x1001,0x1fff] 为 常规错误
	ERROR_CODE_AccountNonexistence ERROR_CODE = 4097 //#tag_desc:账号不存在
)

// Enum value maps for ERROR_CODE.
var (
	ERROR_CODE_name = map[int32]string{
		0:    "Success",
		4097: "AccountNonexistence",
	}
	ERROR_CODE_value = map[string]int32{
		"Success":             0,
		"AccountNonexistence": 4097,
	}
)

func (x ERROR_CODE) Enum() *ERROR_CODE {
	p := new(ERROR_CODE)
	*p = x
	return p
}

func (x ERROR_CODE) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ERROR_CODE) Descriptor() protoreflect.EnumDescriptor {
	return file_ec_proto_enumTypes[0].Descriptor()
}

func (ERROR_CODE) Type() protoreflect.EnumType {
	return &file_ec_proto_enumTypes[0]
}

func (x ERROR_CODE) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ERROR_CODE.Descriptor instead.
func (ERROR_CODE) EnumDescriptor() ([]byte, []int) {
	return file_ec_proto_rawDescGZIP(), []int{0}
}

var File_ec_proto protoreflect.FileDescriptor

var file_ec_proto_rawDesc = []byte{
	0x0a, 0x08, 0x65, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0x33, 0x0a, 0x0a, 0x45, 0x52,
	0x52, 0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x10, 0x00, 0x12, 0x18, 0x0a, 0x13, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x4e, 0x6f, 0x6e, 0x65, 0x78, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x63, 0x65, 0x10, 0x81, 0x20, 0x42,
	0x12, 0x5a, 0x10, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x6c, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ec_proto_rawDescOnce sync.Once
	file_ec_proto_rawDescData = file_ec_proto_rawDesc
)

func file_ec_proto_rawDescGZIP() []byte {
	file_ec_proto_rawDescOnce.Do(func() {
		file_ec_proto_rawDescData = protoimpl.X.CompressGZIP(file_ec_proto_rawDescData)
	})
	return file_ec_proto_rawDescData
}

var file_ec_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_ec_proto_goTypes = []interface{}{
	(ERROR_CODE)(0), // 0: ERROR_CODE
}
var file_ec_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_ec_proto_init() }
func file_ec_proto_init() {
	if File_ec_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_ec_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ec_proto_goTypes,
		DependencyIndexes: file_ec_proto_depIdxs,
		EnumInfos:         file_ec_proto_enumTypes,
	}.Build()
	File_ec_proto = out.File
	file_ec_proto_rawDesc = nil
	file_ec_proto_goTypes = nil
	file_ec_proto_depIdxs = nil
}
