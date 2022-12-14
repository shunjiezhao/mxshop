// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: pk.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type FindType int32

const (
	FindType_Find_NOSRV FindType = 0
	//     1. 复仇 查对局表,id,id,time,winner; 如果全胜或者没有对局就返回错误
	FindType_Avengers FindType = 1
	//   2. 随机选择 利用 choose 接口
	FindType_Random FindType = 2
	//   3. 选择敌人 传入 userid
	FindType_Choose FindType = 3
)

// Enum value maps for FindType.
var (
	FindType_name = map[int32]string{
		0: "Find_NOSRV",
		1: "Avengers",
		2: "Random",
		3: "Choose",
	}
	FindType_value = map[string]int32{
		"Find_NOSRV": 0,
		"Avengers":   1,
		"Random":     2,
		"Choose":     3,
	}
)

func (x FindType) Enum() *FindType {
	p := new(FindType)
	*p = x
	return p
}

func (x FindType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (FindType) Descriptor() protoreflect.EnumDescriptor {
	return file_pk_proto_enumTypes[0].Descriptor()
}

func (FindType) Type() protoreflect.EnumType {
	return &file_pk_proto_enumTypes[0]
}

func (x FindType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use FindType.Descriptor instead.
func (FindType) EnumDescriptor() ([]byte, []int) {
	return file_pk_proto_rawDescGZIP(), []int{0}
}

type JoinRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"` // user_id
	FindType FindType `protobuf:"varint,2,opt,name=find_type,json=findType,proto3,enum=FindType" json:"find_type,omitempty"`
	OtherId  int32    `protobuf:"varint,3,opt,name=other_id,json=otherId,proto3" json:"other_id,omitempty"` // 如果 find type 是 3 有效
}

func (x *JoinRequest) Reset() {
	*x = JoinRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pk_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JoinRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinRequest) ProtoMessage() {}

func (x *JoinRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pk_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JoinRequest.ProtoReflect.Descriptor instead.
func (*JoinRequest) Descriptor() ([]byte, []int) {
	return file_pk_proto_rawDescGZIP(), []int{0}
}

func (x *JoinRequest) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *JoinRequest) GetFindType() FindType {
	if x != nil {
		return x.FindType
	}
	return FindType_Find_NOSRV
}

func (x *JoinRequest) GetOtherId() int32 {
	if x != nil {
		return x.OtherId
	}
	return 0
}

// 对局id
type JoinResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Question []string `protobuf:"bytes,2,rep,name=question,proto3" json:"question,omitempty"` // 问题
	Answer   []string `protobuf:"bytes,3,rep,name=answer,proto3" json:"answer,omitempty"`     // 答案 1-a 2-b 3-c 4-d
}

func (x *JoinResponse) Reset() {
	*x = JoinResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pk_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JoinResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinResponse) ProtoMessage() {}

func (x *JoinResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pk_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JoinResponse.ProtoReflect.Descriptor instead.
func (*JoinResponse) Descriptor() ([]byte, []int) {
	return file_pk_proto_rawDescGZIP(), []int{1}
}

func (x *JoinResponse) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *JoinResponse) GetQuestion() []string {
	if x != nil {
		return x.Question
	}
	return nil
}

func (x *JoinResponse) GetAnswer() []string {
	if x != nil {
		return x.Answer
	}
	return nil
}

type CreateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"` // 对局id
}

func (x *CreateResponse) Reset() {
	*x = CreateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pk_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateResponse) ProtoMessage() {}

func (x *CreateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pk_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateResponse.ProtoReflect.Descriptor instead.
func (*CreateResponse) Descriptor() ([]byte, []int) {
	return file_pk_proto_rawDescGZIP(), []int{2}
}

func (x *CreateResponse) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

// gin层传递的时候按照
// 1. 自己维护一个 并发 map ?
// 2. 直接传递 我们自己检查
type CreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id1 int32 `protobuf:"varint,1,opt,name=id1,proto3" json:"id1,omitempty"`
	Id2 int32 `protobuf:"varint,2,opt,name=id2,proto3" json:"id2,omitempty"`
}

func (x *CreateRequest) Reset() {
	*x = CreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pk_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRequest) ProtoMessage() {}

func (x *CreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pk_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRequest.ProtoReflect.Descriptor instead.
func (*CreateRequest) Descriptor() ([]byte, []int) {
	return file_pk_proto_rawDescGZIP(), []int{3}
}

func (x *CreateRequest) GetId1() int32 {
	if x != nil {
		return x.Id1
	}
	return 0
}

func (x *CreateRequest) GetId2() int32 {
	if x != nil {
		return x.Id2
	}
	return 0
}

type TakePartInRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id  int32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`   //party id
	Uid int32 `protobuf:"varint,2,opt,name=uid,proto3" json:"uid,omitempty"` //user id
}

func (x *TakePartInRequest) Reset() {
	*x = TakePartInRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pk_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TakePartInRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TakePartInRequest) ProtoMessage() {}

func (x *TakePartInRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pk_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TakePartInRequest.ProtoReflect.Descriptor instead.
func (*TakePartInRequest) Descriptor() ([]byte, []int) {
	return file_pk_proto_rawDescGZIP(), []int{4}
}

func (x *TakePartInRequest) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *TakePartInRequest) GetUid() int32 {
	if x != nil {
		return x.Uid
	}
	return 0
}

var File_pk_proto protoreflect.FileDescriptor

var file_pk_proto_rawDesc = []byte{
	0x0a, 0x08, 0x70, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x60, 0x0a, 0x0b, 0x4a, 0x6f, 0x69, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x26, 0x0a, 0x09, 0x66, 0x69, 0x6e, 0x64, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x09, 0x2e, 0x46, 0x69, 0x6e, 0x64,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x08, 0x66, 0x69, 0x6e, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x19,
	0x0a, 0x08, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x07, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x49, 0x64, 0x22, 0x52, 0x0a, 0x0c, 0x4a, 0x6f, 0x69,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x22, 0x20, 0x0a,
	0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x22,
	0x33, 0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x10, 0x0a, 0x03, 0x69, 0x64, 0x31, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x69,
	0x64, 0x31, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x64, 0x32, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x03, 0x69, 0x64, 0x32, 0x22, 0x35, 0x0a, 0x11, 0x54, 0x61, 0x6b, 0x65, 0x50, 0x61, 0x72, 0x74,
	0x49, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x75, 0x69, 0x64, 0x2a, 0x40, 0x0a, 0x08, 0x46,
	0x69, 0x6e, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0e, 0x0a, 0x0a, 0x46, 0x69, 0x6e, 0x64, 0x5f,
	0x4e, 0x4f, 0x53, 0x52, 0x56, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x41, 0x76, 0x65, 0x6e, 0x67,
	0x65, 0x72, 0x73, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x52, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x10,
	0x02, 0x12, 0x0a, 0x0a, 0x06, 0x43, 0x68, 0x6f, 0x6f, 0x73, 0x65, 0x10, 0x03, 0x32, 0x8e, 0x01,
	0x0a, 0x02, 0x50, 0x4b, 0x12, 0x38, 0x0a, 0x0a, 0x54, 0x61, 0x6b, 0x65, 0x50, 0x61, 0x72, 0x74,
	0x49, 0x6e, 0x12, 0x12, 0x2e, 0x54, 0x61, 0x6b, 0x65, 0x50, 0x61, 0x72, 0x74, 0x49, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x23,
	0x0a, 0x04, 0x4a, 0x6f, 0x69, 0x6e, 0x12, 0x0c, 0x2e, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x0e, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x09,
	0x5a, 0x07, 0x2e, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_pk_proto_rawDescOnce sync.Once
	file_pk_proto_rawDescData = file_pk_proto_rawDesc
)

func file_pk_proto_rawDescGZIP() []byte {
	file_pk_proto_rawDescOnce.Do(func() {
		file_pk_proto_rawDescData = protoimpl.X.CompressGZIP(file_pk_proto_rawDescData)
	})
	return file_pk_proto_rawDescData
}

var file_pk_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pk_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_pk_proto_goTypes = []interface{}{
	(FindType)(0),             // 0: FindType
	(*JoinRequest)(nil),       // 1: JoinRequest
	(*JoinResponse)(nil),      // 2: JoinResponse
	(*CreateResponse)(nil),    // 3: CreateResponse
	(*CreateRequest)(nil),     // 4: CreateRequest
	(*TakePartInRequest)(nil), // 5: TakePartInRequest
	(*emptypb.Empty)(nil),     // 6: google.protobuf.Empty
}
var file_pk_proto_depIdxs = []int32{
	0, // 0: JoinRequest.find_type:type_name -> FindType
	5, // 1: PK.TakePartIn:input_type -> TakePartInRequest
	1, // 2: PK.Join:input_type -> JoinRequest
	4, // 3: PK.Create:input_type -> CreateRequest
	6, // 4: PK.TakePartIn:output_type -> google.protobuf.Empty
	2, // 5: PK.Join:output_type -> JoinResponse
	3, // 6: PK.Create:output_type -> CreateResponse
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_pk_proto_init() }
func file_pk_proto_init() {
	if File_pk_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pk_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JoinRequest); i {
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
		file_pk_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JoinResponse); i {
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
		file_pk_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateResponse); i {
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
		file_pk_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateRequest); i {
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
		file_pk_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TakePartInRequest); i {
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
			RawDescriptor: file_pk_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pk_proto_goTypes,
		DependencyIndexes: file_pk_proto_depIdxs,
		EnumInfos:         file_pk_proto_enumTypes,
		MessageInfos:      file_pk_proto_msgTypes,
	}.Build()
	File_pk_proto = out.File
	file_pk_proto_rawDesc = nil
	file_pk_proto_goTypes = nil
	file_pk_proto_depIdxs = nil
}
