// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.12.4
// source: pkg/proto/master.proto

package mpb

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

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_master_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_master_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_pkg_proto_master_proto_rawDescGZIP(), []int{0}
}

type TaskStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskId     int32   `protobuf:"varint,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	WorkerId   string  `protobuf:"bytes,2,opt,name=worker_id,json=workerId,proto3" json:"worker_id,omitempty"`
	Status     bool    `protobuf:"varint,3,opt,name=status,proto3" json:"status,omitempty"`
	Partitions []int32 `protobuf:"varint,4,rep,packed,name=partitions,proto3" json:"partitions,omitempty"`
}

func (x *TaskStatus) Reset() {
	*x = TaskStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_master_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskStatus) ProtoMessage() {}

func (x *TaskStatus) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_master_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskStatus.ProtoReflect.Descriptor instead.
func (*TaskStatus) Descriptor() ([]byte, []int) {
	return file_pkg_proto_master_proto_rawDescGZIP(), []int{1}
}

func (x *TaskStatus) GetTaskId() int32 {
	if x != nil {
		return x.TaskId
	}
	return 0
}

func (x *TaskStatus) GetWorkerId() string {
	if x != nil {
		return x.WorkerId
	}
	return ""
}

func (x *TaskStatus) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

func (x *TaskStatus) GetPartitions() []int32 {
	if x != nil {
		return x.Partitions
	}
	return nil
}

type WorkerStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WorkerId string `protobuf:"bytes,1,opt,name=worker_id,json=workerId,proto3" json:"worker_id,omitempty"`
}

func (x *WorkerStatus) Reset() {
	*x = WorkerStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_master_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorkerStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkerStatus) ProtoMessage() {}

func (x *WorkerStatus) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_master_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorkerStatus.ProtoReflect.Descriptor instead.
func (*WorkerStatus) Descriptor() ([]byte, []int) {
	return file_pkg_proto_master_proto_rawDescGZIP(), []int{2}
}

func (x *WorkerStatus) GetWorkerId() string {
	if x != nil {
		return x.WorkerId
	}
	return ""
}

type ReadyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NumReducers int32 `protobuf:"varint,1,opt,name=NumReducers,proto3" json:"NumReducers,omitempty"` // used for mapper
}

func (x *ReadyResponse) Reset() {
	*x = ReadyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_master_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadyResponse) ProtoMessage() {}

func (x *ReadyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_master_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadyResponse.ProtoReflect.Descriptor instead.
func (*ReadyResponse) Descriptor() ([]byte, []int) {
	return file_pkg_proto_master_proto_rawDescGZIP(), []int{3}
}

func (x *ReadyResponse) GetNumReducers() int32 {
	if x != nil {
		return x.NumReducers
	}
	return 0
}

type EdgeInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Component int32 `protobuf:"varint,1,opt,name=component,proto3" json:"component,omitempty"`
	U         int32 `protobuf:"varint,2,opt,name=u,proto3" json:"u,omitempty"`
	V         int32 `protobuf:"varint,3,opt,name=v,proto3" json:"v,omitempty"`
	W         int32 `protobuf:"varint,4,opt,name=w,proto3" json:"w,omitempty"`
}

func (x *EdgeInfo) Reset() {
	*x = EdgeInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_master_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EdgeInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EdgeInfo) ProtoMessage() {}

func (x *EdgeInfo) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_master_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EdgeInfo.ProtoReflect.Descriptor instead.
func (*EdgeInfo) Descriptor() ([]byte, []int) {
	return file_pkg_proto_master_proto_rawDescGZIP(), []int{4}
}

func (x *EdgeInfo) GetComponent() int32 {
	if x != nil {
		return x.Component
	}
	return 0
}

func (x *EdgeInfo) GetU() int32 {
	if x != nil {
		return x.U
	}
	return 0
}

func (x *EdgeInfo) GetV() int32 {
	if x != nil {
		return x.V
	}
	return 0
}

func (x *EdgeInfo) GetW() int32 {
	if x != nil {
		return x.W
	}
	return 0
}

var File_pkg_proto_master_proto protoreflect.FileDescriptor

var file_pkg_proto_master_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x61, 0x73, 0x74,
	0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x6d, 0x5f, 0x6d, 0x61, 0x70, 0x72,
	0x65, 0x64, 0x75, 0x63, 0x65, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x7a,
	0x0a, 0x0a, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x17, 0x0a, 0x07,
	0x74, 0x61, 0x73, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x74,
	0x61, 0x73, 0x6b, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x61,
	0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x05, 0x52, 0x0a,
	0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x2b, 0x0a, 0x0c, 0x57, 0x6f,
	0x72, 0x6b, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x77, 0x6f,
	0x72, 0x6b, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x77,
	0x6f, 0x72, 0x6b, 0x65, 0x72, 0x49, 0x64, 0x22, 0x31, 0x0a, 0x0d, 0x52, 0x65, 0x61, 0x64, 0x79,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x4e, 0x75, 0x6d, 0x52,
	0x65, 0x64, 0x75, 0x63, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x4e,
	0x75, 0x6d, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x72, 0x73, 0x22, 0x52, 0x0a, 0x08, 0x45, 0x64,
	0x67, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e,
	0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x63, 0x6f, 0x6d, 0x70, 0x6f,
	0x6e, 0x65, 0x6e, 0x74, 0x12, 0x0c, 0x0a, 0x01, 0x75, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x01, 0x75, 0x12, 0x0c, 0x0a, 0x01, 0x76, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01, 0x76,
	0x12, 0x0c, 0x0a, 0x01, 0x77, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01, 0x77, 0x32, 0x82,
	0x02, 0x0a, 0x0d, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x3d, 0x0a, 0x0c, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x61, 0x73, 0x6b,
	0x12, 0x17, 0x2e, 0x6d, 0x5f, 0x6d, 0x61, 0x70, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x2e, 0x54,
	0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x1a, 0x12, 0x2e, 0x6d, 0x5f, 0x6d, 0x61,
	0x70, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12,
	0x40, 0x0a, 0x05, 0x52, 0x65, 0x61, 0x64, 0x79, 0x12, 0x19, 0x2e, 0x6d, 0x5f, 0x6d, 0x61, 0x70,
	0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x1a, 0x1a, 0x2e, 0x6d, 0x5f, 0x6d, 0x61, 0x70, 0x72, 0x65, 0x64, 0x75, 0x63,
	0x65, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x3a, 0x0a, 0x0b, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x69, 0x6e, 0x45, 0x64, 0x67, 0x65,
	0x12, 0x15, 0x2e, 0x6d, 0x5f, 0x6d, 0x61, 0x70, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x2e, 0x45,
	0x64, 0x67, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x12, 0x2e, 0x6d, 0x5f, 0x6d, 0x61, 0x70, 0x72,
	0x65, 0x64, 0x75, 0x63, 0x65, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x34, 0x0a,
	0x08, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x12, 0x2e, 0x6d, 0x5f, 0x6d, 0x61,
	0x70, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x12, 0x2e,
	0x6d, 0x5f, 0x6d, 0x61, 0x70, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x00, 0x42, 0x16, 0x5a, 0x14, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x3b, 0x6d, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_pkg_proto_master_proto_rawDescOnce sync.Once
	file_pkg_proto_master_proto_rawDescData = file_pkg_proto_master_proto_rawDesc
)

func file_pkg_proto_master_proto_rawDescGZIP() []byte {
	file_pkg_proto_master_proto_rawDescOnce.Do(func() {
		file_pkg_proto_master_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_proto_master_proto_rawDescData)
	})
	return file_pkg_proto_master_proto_rawDescData
}

var file_pkg_proto_master_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_pkg_proto_master_proto_goTypes = []any{
	(*Empty)(nil),         // 0: m_mapreduce.Empty
	(*TaskStatus)(nil),    // 1: m_mapreduce.TaskStatus
	(*WorkerStatus)(nil),  // 2: m_mapreduce.WorkerStatus
	(*ReadyResponse)(nil), // 3: m_mapreduce.ReadyResponse
	(*EdgeInfo)(nil),      // 4: m_mapreduce.EdgeInfo
}
var file_pkg_proto_master_proto_depIdxs = []int32{
	1, // 0: m_mapreduce.MasterService.CompleteTask:input_type -> m_mapreduce.TaskStatus
	2, // 1: m_mapreduce.MasterService.Ready:input_type -> m_mapreduce.WorkerStatus
	4, // 2: m_mapreduce.MasterService.SendMinEdge:input_type -> m_mapreduce.EdgeInfo
	0, // 3: m_mapreduce.MasterService.Complete:input_type -> m_mapreduce.Empty
	0, // 4: m_mapreduce.MasterService.CompleteTask:output_type -> m_mapreduce.Empty
	3, // 5: m_mapreduce.MasterService.Ready:output_type -> m_mapreduce.ReadyResponse
	0, // 6: m_mapreduce.MasterService.SendMinEdge:output_type -> m_mapreduce.Empty
	0, // 7: m_mapreduce.MasterService.Complete:output_type -> m_mapreduce.Empty
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_proto_master_proto_init() }
func file_pkg_proto_master_proto_init() {
	if File_pkg_proto_master_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_proto_master_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Empty); i {
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
		file_pkg_proto_master_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*TaskStatus); i {
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
		file_pkg_proto_master_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*WorkerStatus); i {
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
		file_pkg_proto_master_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*ReadyResponse); i {
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
		file_pkg_proto_master_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*EdgeInfo); i {
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
			RawDescriptor: file_pkg_proto_master_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_proto_master_proto_goTypes,
		DependencyIndexes: file_pkg_proto_master_proto_depIdxs,
		MessageInfos:      file_pkg_proto_master_proto_msgTypes,
	}.Build()
	File_pkg_proto_master_proto = out.File
	file_pkg_proto_master_proto_rawDesc = nil
	file_pkg_proto_master_proto_goTypes = nil
	file_pkg_proto_master_proto_depIdxs = nil
}
