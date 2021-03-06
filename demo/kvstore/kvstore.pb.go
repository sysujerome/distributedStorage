// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.6.1
// source: kvstore.proto

package kvstore

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

type SetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *SetRequest) Reset() {
	*x = SetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kvstore_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetRequest) ProtoMessage() {}

func (x *SetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_kvstore_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetRequest.ProtoReflect.Descriptor instead.
func (*SetRequest) Descriptor() ([]byte, []int) {
	return file_kvstore_proto_rawDescGZIP(), []int{0}
}

func (x *SetRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *SetRequest) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type SetReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Err    string `protobuf:"bytes,1,opt,name=err,proto3" json:"err,omitempty"`
	Status string `protobuf:"bytes,2,opt,name=Status,proto3" json:"Status,omitempty"`
	Result string `protobuf:"bytes,3,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *SetReply) Reset() {
	*x = SetReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kvstore_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetReply) ProtoMessage() {}

func (x *SetReply) ProtoReflect() protoreflect.Message {
	mi := &file_kvstore_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetReply.ProtoReflect.Descriptor instead.
func (*SetReply) Descriptor() ([]byte, []int) {
	return file_kvstore_proto_rawDescGZIP(), []int{1}
}

func (x *SetReply) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

func (x *SetReply) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *SetReply) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

type GetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *GetRequest) Reset() {
	*x = GetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kvstore_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRequest) ProtoMessage() {}

func (x *GetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_kvstore_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRequest.ProtoReflect.Descriptor instead.
func (*GetRequest) Descriptor() ([]byte, []int) {
	return file_kvstore_proto_rawDescGZIP(), []int{2}
}

func (x *GetRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type GetReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Err    string `protobuf:"bytes,1,opt,name=err,proto3" json:"err,omitempty"`
	Status string `protobuf:"bytes,2,opt,name=Status,proto3" json:"Status,omitempty"`
	Result string `protobuf:"bytes,3,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *GetReply) Reset() {
	*x = GetReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kvstore_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetReply) ProtoMessage() {}

func (x *GetReply) ProtoReflect() protoreflect.Message {
	mi := &file_kvstore_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetReply.ProtoReflect.Descriptor instead.
func (*GetReply) Descriptor() ([]byte, []int) {
	return file_kvstore_proto_rawDescGZIP(), []int{3}
}

func (x *GetReply) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

func (x *GetReply) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *GetReply) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

type DelRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *DelRequest) Reset() {
	*x = DelRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kvstore_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DelRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DelRequest) ProtoMessage() {}

func (x *DelRequest) ProtoReflect() protoreflect.Message {
	mi := &file_kvstore_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DelRequest.ProtoReflect.Descriptor instead.
func (*DelRequest) Descriptor() ([]byte, []int) {
	return file_kvstore_proto_rawDescGZIP(), []int{4}
}

func (x *DelRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type DelReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Err    string `protobuf:"bytes,1,opt,name=err,proto3" json:"err,omitempty"`
	Status string `protobuf:"bytes,2,opt,name=Status,proto3" json:"Status,omitempty"`
	Result int32  `protobuf:"varint,3,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *DelReply) Reset() {
	*x = DelReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kvstore_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DelReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DelReply) ProtoMessage() {}

func (x *DelReply) ProtoReflect() protoreflect.Message {
	mi := &file_kvstore_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DelReply.ProtoReflect.Descriptor instead.
func (*DelReply) Descriptor() ([]byte, []int) {
	return file_kvstore_proto_rawDescGZIP(), []int{5}
}

func (x *DelReply) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

func (x *DelReply) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *DelReply) GetResult() int32 {
	if x != nil {
		return x.Result
	}
	return 0
}

type SplitRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Port int32 `protobuf:"varint,1,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *SplitRequest) Reset() {
	*x = SplitRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kvstore_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SplitRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SplitRequest) ProtoMessage() {}

func (x *SplitRequest) ProtoReflect() protoreflect.Message {
	mi := &file_kvstore_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SplitRequest.ProtoReflect.Descriptor instead.
func (*SplitRequest) Descriptor() ([]byte, []int) {
	return file_kvstore_proto_rawDescGZIP(), []int{6}
}

func (x *SplitRequest) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

type SplitReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Err    string `protobuf:"bytes,1,opt,name=err,proto3" json:"err,omitempty"`
	Status string `protobuf:"bytes,2,opt,name=Status,proto3" json:"Status,omitempty"`
	Result int32  `protobuf:"varint,3,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *SplitReply) Reset() {
	*x = SplitReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kvstore_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SplitReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SplitReply) ProtoMessage() {}

func (x *SplitReply) ProtoReflect() protoreflect.Message {
	mi := &file_kvstore_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SplitReply.ProtoReflect.Descriptor instead.
func (*SplitReply) Descriptor() ([]byte, []int) {
	return file_kvstore_proto_rawDescGZIP(), []int{7}
}

func (x *SplitReply) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

func (x *SplitReply) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *SplitReply) GetResult() int32 {
	if x != nil {
		return x.Result
	}
	return 0
}

type ScanRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Port int32 `protobuf:"varint,1,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *ScanRequest) Reset() {
	*x = ScanRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kvstore_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScanRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScanRequest) ProtoMessage() {}

func (x *ScanRequest) ProtoReflect() protoreflect.Message {
	mi := &file_kvstore_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScanRequest.ProtoReflect.Descriptor instead.
func (*ScanRequest) Descriptor() ([]byte, []int) {
	return file_kvstore_proto_rawDescGZIP(), []int{8}
}

func (x *ScanRequest) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

type ScanReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Err    string `protobuf:"bytes,1,opt,name=err,proto3" json:"err,omitempty"`
	Status string `protobuf:"bytes,2,opt,name=Status,proto3" json:"Status,omitempty"`
	Result string `protobuf:"bytes,3,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *ScanReply) Reset() {
	*x = ScanReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kvstore_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScanReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScanReply) ProtoMessage() {}

func (x *ScanReply) ProtoReflect() protoreflect.Message {
	mi := &file_kvstore_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScanReply.ProtoReflect.Descriptor instead.
func (*ScanReply) Descriptor() ([]byte, []int) {
	return file_kvstore_proto_rawDescGZIP(), []int{9}
}

func (x *ScanReply) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

func (x *ScanReply) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *ScanReply) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

type SyncRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Next       int64 `protobuf:"varint,1,opt,name=next,proto3" json:"next,omitempty"`
	Level      int64 `protobuf:"varint,2,opt,name=level,proto3" json:"level,omitempty"`
	BeginLevel int64 `protobuf:"varint,3,opt,name=beginLevel,proto3" json:"beginLevel,omitempty"`
}

func (x *SyncRequest) Reset() {
	*x = SyncRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kvstore_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncRequest) ProtoMessage() {}

func (x *SyncRequest) ProtoReflect() protoreflect.Message {
	mi := &file_kvstore_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncRequest.ProtoReflect.Descriptor instead.
func (*SyncRequest) Descriptor() ([]byte, []int) {
	return file_kvstore_proto_rawDescGZIP(), []int{10}
}

func (x *SyncRequest) GetNext() int64 {
	if x != nil {
		return x.Next
	}
	return 0
}

func (x *SyncRequest) GetLevel() int64 {
	if x != nil {
		return x.Level
	}
	return 0
}

func (x *SyncRequest) GetBeginLevel() int64 {
	if x != nil {
		return x.BeginLevel
	}
	return 0
}

type SyncReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Err    string `protobuf:"bytes,1,opt,name=err,proto3" json:"err,omitempty"`
	Status string `protobuf:"bytes,2,opt,name=Status,proto3" json:"Status,omitempty"`
	Result string `protobuf:"bytes,3,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *SyncReply) Reset() {
	*x = SyncReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kvstore_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncReply) ProtoMessage() {}

func (x *SyncReply) ProtoReflect() protoreflect.Message {
	mi := &file_kvstore_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncReply.ProtoReflect.Descriptor instead.
func (*SyncReply) Descriptor() ([]byte, []int) {
	return file_kvstore_proto_rawDescGZIP(), []int{11}
}

func (x *SyncReply) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

func (x *SyncReply) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *SyncReply) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

var File_kvstore_proto protoreflect.FileDescriptor

var file_kvstore_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6b, 0x76, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x6b, 0x76, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x22, 0x34, 0x0a, 0x0a, 0x53, 0x65, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x4c,
	0x0a, 0x08, 0x53, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x72,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x72, 0x72, 0x12, 0x16, 0x0a, 0x06,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x1e, 0x0a, 0x0a,
	0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22, 0x4c, 0x0a, 0x08,
	0x47, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x72, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x72, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x1e, 0x0a, 0x0a, 0x44, 0x65,
	0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22, 0x4c, 0x0a, 0x08, 0x44, 0x65,
	0x6c, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x72, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x72, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x22, 0x0a, 0x0c, 0x53, 0x70, 0x6c, 0x69,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x4e, 0x0a, 0x0a,
	0x53, 0x70, 0x6c, 0x69, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x72,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x72, 0x72, 0x12, 0x16, 0x0a, 0x06,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x21, 0x0a, 0x0b,
	0x53, 0x63, 0x61, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70,
	0x6f, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x22,
	0x4d, 0x0a, 0x09, 0x53, 0x63, 0x61, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x65, 0x72, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x72, 0x72, 0x12, 0x16,
	0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x57,
	0x0a, 0x0b, 0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x6e, 0x65, 0x78,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x1e, 0x0a, 0x0a, 0x62, 0x65, 0x67, 0x69, 0x6e,
	0x4c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x62, 0x65, 0x67,
	0x69, 0x6e, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x22, 0x4d, 0x0a, 0x09, 0x53, 0x79, 0x6e, 0x63, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x72, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x65, 0x72, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16,
	0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x32, 0xbf, 0x02, 0x0a, 0x07, 0x53, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x12, 0x2f, 0x0a, 0x03, 0x53, 0x65, 0x74, 0x12, 0x13, 0x2e, 0x6b, 0x76, 0x73, 0x74,
	0x6f, 0x72, 0x65, 0x2e, 0x53, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11,
	0x2e, 0x6b, 0x76, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x53, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x22, 0x00, 0x12, 0x2f, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x13, 0x2e, 0x6b, 0x76, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x11, 0x2e, 0x6b, 0x76, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x22, 0x00, 0x12, 0x2f, 0x0a, 0x03, 0x44, 0x65, 0x6c, 0x12, 0x13, 0x2e, 0x6b, 0x76,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x44, 0x65, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x11, 0x2e, 0x6b, 0x76, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x44, 0x65, 0x6c, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x35, 0x0a, 0x05, 0x53, 0x70, 0x6c, 0x69, 0x74, 0x12, 0x15,
	0x2e, 0x6b, 0x76, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x53, 0x70, 0x6c, 0x69, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x6b, 0x76, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e,
	0x53, 0x70, 0x6c, 0x69, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x32, 0x0a, 0x04,
	0x53, 0x63, 0x61, 0x6e, 0x12, 0x14, 0x2e, 0x6b, 0x76, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x53,
	0x63, 0x61, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x6b, 0x76, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x2e, 0x53, 0x63, 0x61, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00,
	0x12, 0x36, 0x0a, 0x08, 0x53, 0x79, 0x6e, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x12, 0x14, 0x2e, 0x6b,
	0x76, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x12, 0x2e, 0x6b, 0x76, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x53, 0x79, 0x6e,
	0x63, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x0b, 0x5a, 0x09, 0x2e, 0x3b, 0x6b, 0x76,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_kvstore_proto_rawDescOnce sync.Once
	file_kvstore_proto_rawDescData = file_kvstore_proto_rawDesc
)

func file_kvstore_proto_rawDescGZIP() []byte {
	file_kvstore_proto_rawDescOnce.Do(func() {
		file_kvstore_proto_rawDescData = protoimpl.X.CompressGZIP(file_kvstore_proto_rawDescData)
	})
	return file_kvstore_proto_rawDescData
}

var file_kvstore_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_kvstore_proto_goTypes = []interface{}{
	(*SetRequest)(nil),   // 0: kvstore.SetRequest
	(*SetReply)(nil),     // 1: kvstore.SetReply
	(*GetRequest)(nil),   // 2: kvstore.GetRequest
	(*GetReply)(nil),     // 3: kvstore.GetReply
	(*DelRequest)(nil),   // 4: kvstore.DelRequest
	(*DelReply)(nil),     // 5: kvstore.DelReply
	(*SplitRequest)(nil), // 6: kvstore.SplitRequest
	(*SplitReply)(nil),   // 7: kvstore.SplitReply
	(*ScanRequest)(nil),  // 8: kvstore.ScanRequest
	(*ScanReply)(nil),    // 9: kvstore.ScanReply
	(*SyncRequest)(nil),  // 10: kvstore.SyncRequest
	(*SyncReply)(nil),    // 11: kvstore.SyncReply
}
var file_kvstore_proto_depIdxs = []int32{
	0,  // 0: kvstore.Storage.Set:input_type -> kvstore.SetRequest
	2,  // 1: kvstore.Storage.Get:input_type -> kvstore.GetRequest
	4,  // 2: kvstore.Storage.Del:input_type -> kvstore.DelRequest
	6,  // 3: kvstore.Storage.Split:input_type -> kvstore.SplitRequest
	8,  // 4: kvstore.Storage.Scan:input_type -> kvstore.ScanRequest
	10, // 5: kvstore.Storage.SyncConf:input_type -> kvstore.SyncRequest
	1,  // 6: kvstore.Storage.Set:output_type -> kvstore.SetReply
	3,  // 7: kvstore.Storage.Get:output_type -> kvstore.GetReply
	5,  // 8: kvstore.Storage.Del:output_type -> kvstore.DelReply
	7,  // 9: kvstore.Storage.Split:output_type -> kvstore.SplitReply
	9,  // 10: kvstore.Storage.Scan:output_type -> kvstore.ScanReply
	11, // 11: kvstore.Storage.SyncConf:output_type -> kvstore.SyncReply
	6,  // [6:12] is the sub-list for method output_type
	0,  // [0:6] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_kvstore_proto_init() }
func file_kvstore_proto_init() {
	if File_kvstore_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_kvstore_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetRequest); i {
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
		file_kvstore_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetReply); i {
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
		file_kvstore_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRequest); i {
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
		file_kvstore_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetReply); i {
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
		file_kvstore_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DelRequest); i {
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
		file_kvstore_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DelReply); i {
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
		file_kvstore_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SplitRequest); i {
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
		file_kvstore_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SplitReply); i {
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
		file_kvstore_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScanRequest); i {
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
		file_kvstore_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScanReply); i {
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
		file_kvstore_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SyncRequest); i {
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
		file_kvstore_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SyncReply); i {
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
			RawDescriptor: file_kvstore_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_kvstore_proto_goTypes,
		DependencyIndexes: file_kvstore_proto_depIdxs,
		MessageInfos:      file_kvstore_proto_msgTypes,
	}.Build()
	File_kvstore_proto = out.File
	file_kvstore_proto_rawDesc = nil
	file_kvstore_proto_goTypes = nil
	file_kvstore_proto_depIdxs = nil
}
