// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.15.8
// source: gate.proto

package pb

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

type BindReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cid int64 `protobuf:"varint,1,opt,name=cid,proto3" json:"cid,omitempty"` // 连接ID
	Uid int64 `protobuf:"varint,2,opt,name=uid,proto3" json:"uid,omitempty"` // 用户ID
}

func (x *BindReq) Reset() {
	*x = BindReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gate_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BindReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BindReq) ProtoMessage() {}

func (x *BindReq) ProtoReflect() protoreflect.Message {
	mi := &file_gate_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BindReq.ProtoReflect.Descriptor instead.
func (*BindReq) Descriptor() ([]byte, []int) {
	return file_gate_proto_rawDescGZIP(), []int{0}
}

func (x *BindReq) GetCid() int64 {
	if x != nil {
		return x.Cid
	}
	return 0
}

func (x *BindReq) GetUid() int64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

type BindReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *BindReply) Reset() {
	*x = BindReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gate_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BindReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BindReply) ProtoMessage() {}

func (x *BindReply) ProtoReflect() protoreflect.Message {
	mi := &file_gate_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BindReply.ProtoReflect.Descriptor instead.
func (*BindReply) Descriptor() ([]byte, []int) {
	return file_gate_proto_rawDescGZIP(), []int{1}
}

type UnbindReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid int64 `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"` // 用户ID
}

func (x *UnbindReq) Reset() {
	*x = UnbindReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gate_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnbindReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnbindReq) ProtoMessage() {}

func (x *UnbindReq) ProtoReflect() protoreflect.Message {
	mi := &file_gate_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnbindReq.ProtoReflect.Descriptor instead.
func (*UnbindReq) Descriptor() ([]byte, []int) {
	return file_gate_proto_rawDescGZIP(), []int{2}
}

func (x *UnbindReq) GetUid() int64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

type UnbindReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UnbindReply) Reset() {
	*x = UnbindReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gate_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnbindReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnbindReply) ProtoMessage() {}

func (x *UnbindReply) ProtoReflect() protoreflect.Message {
	mi := &file_gate_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnbindReply.ProtoReflect.Descriptor instead.
func (*UnbindReply) Descriptor() ([]byte, []int) {
	return file_gate_proto_rawDescGZIP(), []int{3}
}

type PushReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Target  int64    `protobuf:"varint,3,opt,name=target,proto3" json:"target,omitempty"`  // 推送目标
	Message *Message `protobuf:"bytes,5,opt,name=Message,proto3" json:"Message,omitempty"` // 消息
}

func (x *PushReq) Reset() {
	*x = PushReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gate_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushReq) ProtoMessage() {}

func (x *PushReq) ProtoReflect() protoreflect.Message {
	mi := &file_gate_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushReq.ProtoReflect.Descriptor instead.
func (*PushReq) Descriptor() ([]byte, []int) {
	return file_gate_proto_rawDescGZIP(), []int{4}
}

func (x *PushReq) GetTarget() int64 {
	if x != nil {
		return x.Target
	}
	return 0
}

func (x *PushReq) GetMessage() *Message {
	if x != nil {
		return x.Message
	}
	return nil
}

type PushReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PushReply) Reset() {
	*x = PushReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gate_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushReply) ProtoMessage() {}

func (x *PushReply) ProtoReflect() protoreflect.Message {
	mi := &file_gate_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushReply.ProtoReflect.Descriptor instead.
func (*PushReply) Descriptor() ([]byte, []int) {
	return file_gate_proto_rawDescGZIP(), []int{5}
}

type GetIPReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid int64 `protobuf:"varint,1,opt,name=Uid,proto3" json:"Uid,omitempty"`
}

func (x *GetIPReq) Reset() {
	*x = GetIPReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gate_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetIPReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetIPReq) ProtoMessage() {}

func (x *GetIPReq) ProtoReflect() protoreflect.Message {
	mi := &file_gate_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetIPReq.ProtoReflect.Descriptor instead.
func (*GetIPReq) Descriptor() ([]byte, []int) {
	return file_gate_proto_rawDescGZIP(), []int{6}
}

func (x *GetIPReq) GetUid() int64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

type GetIPReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IP string `protobuf:"bytes,1,opt,name=IP,proto3" json:"IP,omitempty"`
}

func (x *GetIPReply) Reset() {
	*x = GetIPReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gate_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetIPReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetIPReply) ProtoMessage() {}

func (x *GetIPReply) ProtoReflect() protoreflect.Message {
	mi := &file_gate_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetIPReply.ProtoReflect.Descriptor instead.
func (*GetIPReply) Descriptor() ([]byte, []int) {
	return file_gate_proto_rawDescGZIP(), []int{7}
}

func (x *GetIPReply) GetIP() string {
	if x != nil {
		return x.IP
	}
	return ""
}

var File_gate_proto protoreflect.FileDescriptor

var file_gate_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x67, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62,
	0x1a, 0x0d, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x2d, 0x0a, 0x07, 0x42, 0x69, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x12, 0x10, 0x0a, 0x03, 0x63, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x63, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03,
	0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x75, 0x69, 0x64, 0x22, 0x0b,
	0x0a, 0x09, 0x42, 0x69, 0x6e, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x1d, 0x0a, 0x09, 0x55,
	0x6e, 0x62, 0x69, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x75, 0x69, 0x64, 0x22, 0x0d, 0x0a, 0x0b, 0x55, 0x6e,
	0x62, 0x69, 0x6e, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x48, 0x0a, 0x07, 0x50, 0x75, 0x73,
	0x68, 0x52, 0x65, 0x71, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12, 0x25, 0x0a, 0x07,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e,
	0x70, 0x62, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x07, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x22, 0x0b, 0x0a, 0x09, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x22, 0x1c, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x49, 0x50, 0x52, 0x65, 0x71, 0x12, 0x10, 0x0a, 0x03,
	0x55, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x55, 0x69, 0x64, 0x22, 0x1c,
	0x0a, 0x0a, 0x47, 0x65, 0x74, 0x49, 0x50, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x0e, 0x0a, 0x02,
	0x49, 0x50, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x50, 0x32, 0xa7, 0x01, 0x0a,
	0x04, 0x47, 0x61, 0x74, 0x65, 0x12, 0x24, 0x0a, 0x04, 0x42, 0x69, 0x6e, 0x64, 0x12, 0x0b, 0x2e,
	0x70, 0x62, 0x2e, 0x42, 0x69, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x1a, 0x0d, 0x2e, 0x70, 0x62, 0x2e,
	0x42, 0x69, 0x6e, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x2a, 0x0a, 0x06, 0x55,
	0x6e, 0x62, 0x69, 0x6e, 0x64, 0x12, 0x0d, 0x2e, 0x70, 0x62, 0x2e, 0x55, 0x6e, 0x62, 0x69, 0x6e,
	0x64, 0x52, 0x65, 0x71, 0x1a, 0x0f, 0x2e, 0x70, 0x62, 0x2e, 0x55, 0x6e, 0x62, 0x69, 0x6e, 0x64,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x24, 0x0a, 0x04, 0x50, 0x75, 0x73, 0x68, 0x12,
	0x0b, 0x2e, 0x70, 0x62, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65, 0x71, 0x1a, 0x0d, 0x2e, 0x70,
	0x62, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x27, 0x0a,
	0x05, 0x47, 0x65, 0x74, 0x49, 0x50, 0x12, 0x0c, 0x2e, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x49,
	0x50, 0x52, 0x65, 0x71, 0x1a, 0x0e, 0x2e, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x49, 0x50, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x70, 0x62, 0x3b, 0x70,
	0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gate_proto_rawDescOnce sync.Once
	file_gate_proto_rawDescData = file_gate_proto_rawDesc
)

func file_gate_proto_rawDescGZIP() []byte {
	file_gate_proto_rawDescOnce.Do(func() {
		file_gate_proto_rawDescData = protoimpl.X.CompressGZIP(file_gate_proto_rawDescData)
	})
	return file_gate_proto_rawDescData
}

var file_gate_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_gate_proto_goTypes = []interface{}{
	(*BindReq)(nil),     // 0: pb.BindReq
	(*BindReply)(nil),   // 1: pb.BindReply
	(*UnbindReq)(nil),   // 2: pb.UnbindReq
	(*UnbindReply)(nil), // 3: pb.UnbindReply
	(*PushReq)(nil),     // 4: pb.PushReq
	(*PushReply)(nil),   // 5: pb.PushReply
	(*GetIPReq)(nil),    // 6: pb.GetIPReq
	(*GetIPReply)(nil),  // 7: pb.GetIPReply
	(*Message)(nil),     // 8: pb.Message
}
var file_gate_proto_depIdxs = []int32{
	8, // 0: pb.PushReq.Message:type_name -> pb.Message
	0, // 1: pb.Gate.Bind:input_type -> pb.BindReq
	2, // 2: pb.Gate.Unbind:input_type -> pb.UnbindReq
	4, // 3: pb.Gate.Push:input_type -> pb.PushReq
	6, // 4: pb.Gate.GetIP:input_type -> pb.GetIPReq
	1, // 5: pb.Gate.Bind:output_type -> pb.BindReply
	3, // 6: pb.Gate.Unbind:output_type -> pb.UnbindReply
	5, // 7: pb.Gate.Push:output_type -> pb.PushReply
	7, // 8: pb.Gate.GetIP:output_type -> pb.GetIPReply
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_gate_proto_init() }
func file_gate_proto_init() {
	if File_gate_proto != nil {
		return
	}
	file_message_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_gate_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BindReq); i {
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
		file_gate_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BindReply); i {
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
		file_gate_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnbindReq); i {
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
		file_gate_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnbindReply); i {
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
		file_gate_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushReq); i {
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
		file_gate_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushReply); i {
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
		file_gate_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetIPReq); i {
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
		file_gate_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetIPReply); i {
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
			RawDescriptor: file_gate_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_gate_proto_goTypes,
		DependencyIndexes: file_gate_proto_depIdxs,
		MessageInfos:      file_gate_proto_msgTypes,
	}.Build()
	File_gate_proto = out.File
	file_gate_proto_rawDesc = nil
	file_gate_proto_goTypes = nil
	file_gate_proto_depIdxs = nil
}
