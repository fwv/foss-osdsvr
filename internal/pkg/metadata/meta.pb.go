// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.15.8
// source: internal/pkg/metadata/meta.proto

package metadata

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

// 对象元数据
type MetaData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name          string                   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`                                                                                                  // 文件名
	BucketId      string                   `protobuf:"bytes,2,opt,name=bucket_id,json=bucketId,proto3" json:"bucket_id,omitempty"`                                                                          // 桶ID
	Versions      map[int64]*VersionRecord `protobuf:"bytes,3,rep,name=versions,proto3" json:"versions,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"` // 版本记录
	LatestVersion int64                    `protobuf:"varint,4,opt,name=latest_version,json=latestVersion,proto3" json:"latest_version,omitempty"`                                                          // 最新版本号
}

func (x *MetaData) Reset() {
	*x = MetaData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pkg_metadata_meta_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetaData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetaData) ProtoMessage() {}

func (x *MetaData) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pkg_metadata_meta_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetaData.ProtoReflect.Descriptor instead.
func (*MetaData) Descriptor() ([]byte, []int) {
	return file_internal_pkg_metadata_meta_proto_rawDescGZIP(), []int{0}
}

func (x *MetaData) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *MetaData) GetBucketId() string {
	if x != nil {
		return x.BucketId
	}
	return ""
}

func (x *MetaData) GetVersions() map[int64]*VersionRecord {
	if x != nil {
		return x.Versions
	}
	return nil
}

func (x *MetaData) GetLatestVersion() int64 {
	if x != nil {
		return x.LatestVersion
	}
	return 0
}

type VersionRecord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version    int64  `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`                         // 版本号
	Hash       string `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`                                // 文件标识
	Size       int64  `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`                               // 文件大小
	UploadTime int64  `protobuf:"varint,4,opt,name=upload_time,json=uploadTime,proto3" json:"upload_time,omitempty"` // 上传时间
}

func (x *VersionRecord) Reset() {
	*x = VersionRecord{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pkg_metadata_meta_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VersionRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VersionRecord) ProtoMessage() {}

func (x *VersionRecord) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pkg_metadata_meta_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VersionRecord.ProtoReflect.Descriptor instead.
func (*VersionRecord) Descriptor() ([]byte, []int) {
	return file_internal_pkg_metadata_meta_proto_rawDescGZIP(), []int{1}
}

func (x *VersionRecord) GetVersion() int64 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *VersionRecord) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *VersionRecord) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *VersionRecord) GetUploadTime() int64 {
	if x != nil {
		return x.UploadTime
	}
	return 0
}

var File_internal_pkg_metadata_meta_proto protoreflect.FileDescriptor

var file_internal_pkg_metadata_meta_proto_rawDesc = []byte{
	0x0a, 0x20, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x22, 0xf6, 0x01, 0x0a,
	0x08, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a,
	0x09, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x49, 0x64, 0x12, 0x3c, 0x0a, 0x08, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61,
	0x2e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x6c, 0x61, 0x74, 0x65,
	0x73, 0x74, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x0d, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x1a,
	0x54, 0x0a, 0x0d, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x2d, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x17, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x72, 0x0a, 0x0d, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x68, 0x61, 0x73, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x75, 0x70, 0x6c, 0x6f,
	0x61, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x75,
	0x70, 0x6c, 0x6f, 0x61, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x42, 0x27, 0x5a, 0x25, 0x6f, 0x73, 0x64,
	0x73, 0x76, 0x72, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_pkg_metadata_meta_proto_rawDescOnce sync.Once
	file_internal_pkg_metadata_meta_proto_rawDescData = file_internal_pkg_metadata_meta_proto_rawDesc
)

func file_internal_pkg_metadata_meta_proto_rawDescGZIP() []byte {
	file_internal_pkg_metadata_meta_proto_rawDescOnce.Do(func() {
		file_internal_pkg_metadata_meta_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_pkg_metadata_meta_proto_rawDescData)
	})
	return file_internal_pkg_metadata_meta_proto_rawDescData
}

var file_internal_pkg_metadata_meta_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_internal_pkg_metadata_meta_proto_goTypes = []interface{}{
	(*MetaData)(nil),      // 0: metadata.MetaData
	(*VersionRecord)(nil), // 1: metadata.VersionRecord
	nil,                   // 2: metadata.MetaData.VersionsEntry
}
var file_internal_pkg_metadata_meta_proto_depIdxs = []int32{
	2, // 0: metadata.MetaData.versions:type_name -> metadata.MetaData.VersionsEntry
	1, // 1: metadata.MetaData.VersionsEntry.value:type_name -> metadata.VersionRecord
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_internal_pkg_metadata_meta_proto_init() }
func file_internal_pkg_metadata_meta_proto_init() {
	if File_internal_pkg_metadata_meta_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_pkg_metadata_meta_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MetaData); i {
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
		file_internal_pkg_metadata_meta_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VersionRecord); i {
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
			RawDescriptor: file_internal_pkg_metadata_meta_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_pkg_metadata_meta_proto_goTypes,
		DependencyIndexes: file_internal_pkg_metadata_meta_proto_depIdxs,
		MessageInfos:      file_internal_pkg_metadata_meta_proto_msgTypes,
	}.Build()
	File_internal_pkg_metadata_meta_proto = out.File
	file_internal_pkg_metadata_meta_proto_rawDesc = nil
	file_internal_pkg_metadata_meta_proto_goTypes = nil
	file_internal_pkg_metadata_meta_proto_depIdxs = nil
}