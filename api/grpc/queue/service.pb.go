// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.29.2
// source: api/grpc/queue/service.proto

package queue

import (
	context "context"
	pb "github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type SetQueueRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Subj string `protobuf:"bytes,3,opt,name=subj,proto3" json:"subj,omitempty"`
}

func (x *SetQueueRequest) Reset() {
	*x = SetQueueRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_grpc_queue_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetQueueRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetQueueRequest) ProtoMessage() {}

func (x *SetQueueRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_grpc_queue_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetQueueRequest.ProtoReflect.Descriptor instead.
func (*SetQueueRequest) Descriptor() ([]byte, []int) {
	return file_api_grpc_queue_service_proto_rawDescGZIP(), []int{0}
}

func (x *SetQueueRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SetQueueRequest) GetSubj() string {
	if x != nil {
		return x.Subj
	}
	return ""
}

type ReceiveMessagesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Command:
	//
	//	*ReceiveMessagesRequest_Start
	//	*ReceiveMessagesRequest_Ack
	Command isReceiveMessagesRequest_Command `protobuf_oneof:"command"`
}

func (x *ReceiveMessagesRequest) Reset() {
	*x = ReceiveMessagesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_grpc_queue_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReceiveMessagesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReceiveMessagesRequest) ProtoMessage() {}

func (x *ReceiveMessagesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_grpc_queue_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReceiveMessagesRequest.ProtoReflect.Descriptor instead.
func (*ReceiveMessagesRequest) Descriptor() ([]byte, []int) {
	return file_api_grpc_queue_service_proto_rawDescGZIP(), []int{1}
}

func (m *ReceiveMessagesRequest) GetCommand() isReceiveMessagesRequest_Command {
	if m != nil {
		return m.Command
	}
	return nil
}

func (x *ReceiveMessagesRequest) GetStart() *ReceiveMessagesCommandStart {
	if x, ok := x.GetCommand().(*ReceiveMessagesRequest_Start); ok {
		return x.Start
	}
	return nil
}

func (x *ReceiveMessagesRequest) GetAck() *ReceiveMessagesCommandAck {
	if x, ok := x.GetCommand().(*ReceiveMessagesRequest_Ack); ok {
		return x.Ack
	}
	return nil
}

type isReceiveMessagesRequest_Command interface {
	isReceiveMessagesRequest_Command()
}

type ReceiveMessagesRequest_Start struct {
	Start *ReceiveMessagesCommandStart `protobuf:"bytes,1,opt,name=start,proto3,oneof"`
}

type ReceiveMessagesRequest_Ack struct {
	Ack *ReceiveMessagesCommandAck `protobuf:"bytes,2,opt,name=ack,proto3,oneof"`
}

func (*ReceiveMessagesRequest_Start) isReceiveMessagesRequest_Command() {}

func (*ReceiveMessagesRequest_Ack) isReceiveMessagesRequest_Command() {}

type ReceiveMessagesCommandStart struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Queue     string `protobuf:"bytes,1,opt,name=queue,proto3" json:"queue,omitempty"`
	BatchSize uint32 `protobuf:"varint,2,opt,name=batchSize,proto3" json:"batchSize,omitempty"`
	Subj      string `protobuf:"bytes,3,opt,name=subj,proto3" json:"subj,omitempty"`
}

func (x *ReceiveMessagesCommandStart) Reset() {
	*x = ReceiveMessagesCommandStart{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_grpc_queue_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReceiveMessagesCommandStart) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReceiveMessagesCommandStart) ProtoMessage() {}

func (x *ReceiveMessagesCommandStart) ProtoReflect() protoreflect.Message {
	mi := &file_api_grpc_queue_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReceiveMessagesCommandStart.ProtoReflect.Descriptor instead.
func (*ReceiveMessagesCommandStart) Descriptor() ([]byte, []int) {
	return file_api_grpc_queue_service_proto_rawDescGZIP(), []int{2}
}

func (x *ReceiveMessagesCommandStart) GetQueue() string {
	if x != nil {
		return x.Queue
	}
	return ""
}

func (x *ReceiveMessagesCommandStart) GetBatchSize() uint32 {
	if x != nil {
		return x.BatchSize
	}
	return 0
}

func (x *ReceiveMessagesCommandStart) GetSubj() string {
	if x != nil {
		return x.Subj
	}
	return ""
}

type ReceiveMessagesCommandAck struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count uint32 `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *ReceiveMessagesCommandAck) Reset() {
	*x = ReceiveMessagesCommandAck{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_grpc_queue_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReceiveMessagesCommandAck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReceiveMessagesCommandAck) ProtoMessage() {}

func (x *ReceiveMessagesCommandAck) ProtoReflect() protoreflect.Message {
	mi := &file_api_grpc_queue_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReceiveMessagesCommandAck.ProtoReflect.Descriptor instead.
func (*ReceiveMessagesCommandAck) Descriptor() ([]byte, []int) {
	return file_api_grpc_queue_service_proto_rawDescGZIP(), []int{3}
}

func (x *ReceiveMessagesCommandAck) GetCount() uint32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type ReceiveMessagesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msgs []*pb.CloudEvent `protobuf:"bytes,1,rep,name=msgs,proto3" json:"msgs,omitempty"`
}

func (x *ReceiveMessagesResponse) Reset() {
	*x = ReceiveMessagesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_grpc_queue_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReceiveMessagesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReceiveMessagesResponse) ProtoMessage() {}

func (x *ReceiveMessagesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_grpc_queue_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReceiveMessagesResponse.ProtoReflect.Descriptor instead.
func (*ReceiveMessagesResponse) Descriptor() ([]byte, []int) {
	return file_api_grpc_queue_service_proto_rawDescGZIP(), []int{4}
}

func (x *ReceiveMessagesResponse) GetMsgs() []*pb.CloudEvent {
	if x != nil {
		return x.Msgs
	}
	return nil
}

var File_api_grpc_queue_service_proto protoreflect.FileDescriptor

var file_api_grpc_queue_service_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x71, 0x75, 0x65, 0x75, 0x65,
	0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d,
	0x61, 0x77, 0x61, 0x6b, 0x61, 0x72, 0x69, 0x2e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x1a, 0x1b, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65,
	0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x25, 0x61, 0x70, 0x69, 0x2f,
	0x67, 0x72, 0x70, 0x63, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73,
	0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x39, 0x0a, 0x0f, 0x53, 0x65, 0x74, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x75, 0x62, 0x6a,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73, 0x75, 0x62, 0x6a, 0x22, 0xa5, 0x01, 0x0a,
	0x16, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x42, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x61, 0x77, 0x61, 0x6b, 0x61, 0x72, 0x69,
	0x2e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x53, 0x74, 0x61,
	0x72, 0x74, 0x48, 0x00, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12, 0x3c, 0x0a, 0x03, 0x61,
	0x63, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x61, 0x77, 0x61, 0x6b, 0x61,
	0x72, 0x69, 0x2e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x41,
	0x63, 0x6b, 0x48, 0x00, 0x52, 0x03, 0x61, 0x63, 0x6b, 0x42, 0x09, 0x0a, 0x07, 0x63, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x22, 0x65, 0x0a, 0x1b, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x53, 0x74,
	0x61, 0x72, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x71, 0x75, 0x65, 0x75, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x62, 0x61, 0x74,
	0x63, 0x68, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x62, 0x61,
	0x74, 0x63, 0x68, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x75, 0x62, 0x6a, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73, 0x75, 0x62, 0x6a, 0x22, 0x31, 0x0a, 0x19, 0x52,
	0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x43, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x41, 0x63, 0x6b, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x3d,
	0x0a, 0x17, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a, 0x04, 0x6d, 0x73, 0x67,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x6c, 0x6f,
	0x75, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x04, 0x6d, 0x73, 0x67, 0x73, 0x32, 0xb3, 0x01,
	0x0a, 0x07, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x42, 0x0a, 0x08, 0x53, 0x65, 0x74,
	0x51, 0x75, 0x65, 0x75, 0x65, 0x12, 0x1e, 0x2e, 0x61, 0x77, 0x61, 0x6b, 0x61, 0x72, 0x69, 0x2e,
	0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x53, 0x65, 0x74, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x64, 0x0a,
	0x0f, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73,
	0x12, 0x25, 0x2e, 0x61, 0x77, 0x61, 0x6b, 0x61, 0x72, 0x69, 0x2e, 0x71, 0x75, 0x65, 0x75, 0x65,
	0x2e, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x61, 0x77, 0x61, 0x6b, 0x61, 0x72,
	0x69, 0x2e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28,
	0x01, 0x30, 0x01, 0x42, 0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x61, 0x77, 0x61, 0x6b, 0x61, 0x72, 0x69, 0x2f, 0x69, 0x6e, 0x74, 0x2d, 0x62, 0x6c,
	0x75, 0x65, 0x73, 0x6b, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x71,
	0x75, 0x65, 0x75, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_grpc_queue_service_proto_rawDescOnce sync.Once
	file_api_grpc_queue_service_proto_rawDescData = file_api_grpc_queue_service_proto_rawDesc
)

func file_api_grpc_queue_service_proto_rawDescGZIP() []byte {
	file_api_grpc_queue_service_proto_rawDescOnce.Do(func() {
		file_api_grpc_queue_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_grpc_queue_service_proto_rawDescData)
	})
	return file_api_grpc_queue_service_proto_rawDescData
}

var file_api_grpc_queue_service_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_api_grpc_queue_service_proto_goTypes = []interface{}{
	(*SetQueueRequest)(nil),             // 0: awakari.queue.SetQueueRequest
	(*ReceiveMessagesRequest)(nil),      // 1: awakari.queue.ReceiveMessagesRequest
	(*ReceiveMessagesCommandStart)(nil), // 2: awakari.queue.ReceiveMessagesCommandStart
	(*ReceiveMessagesCommandAck)(nil),   // 3: awakari.queue.ReceiveMessagesCommandAck
	(*ReceiveMessagesResponse)(nil),     // 4: awakari.queue.ReceiveMessagesResponse
	(*pb.CloudEvent)(nil),               // 5: pb.CloudEvent
	(*emptypb.Empty)(nil),               // 6: google.protobuf.Empty
}
var file_api_grpc_queue_service_proto_depIdxs = []int32{
	2, // 0: awakari.queue.ReceiveMessagesRequest.start:type_name -> awakari.queue.ReceiveMessagesCommandStart
	3, // 1: awakari.queue.ReceiveMessagesRequest.ack:type_name -> awakari.queue.ReceiveMessagesCommandAck
	5, // 2: awakari.queue.ReceiveMessagesResponse.msgs:type_name -> pb.CloudEvent
	0, // 3: awakari.queue.Service.SetQueue:input_type -> awakari.queue.SetQueueRequest
	1, // 4: awakari.queue.Service.ReceiveMessages:input_type -> awakari.queue.ReceiveMessagesRequest
	6, // 5: awakari.queue.Service.SetQueue:output_type -> google.protobuf.Empty
	4, // 6: awakari.queue.Service.ReceiveMessages:output_type -> awakari.queue.ReceiveMessagesResponse
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_api_grpc_queue_service_proto_init() }
func file_api_grpc_queue_service_proto_init() {
	if File_api_grpc_queue_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_grpc_queue_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetQueueRequest); i {
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
		file_api_grpc_queue_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReceiveMessagesRequest); i {
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
		file_api_grpc_queue_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReceiveMessagesCommandStart); i {
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
		file_api_grpc_queue_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReceiveMessagesCommandAck); i {
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
		file_api_grpc_queue_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReceiveMessagesResponse); i {
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
	file_api_grpc_queue_service_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*ReceiveMessagesRequest_Start)(nil),
		(*ReceiveMessagesRequest_Ack)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_grpc_queue_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_grpc_queue_service_proto_goTypes,
		DependencyIndexes: file_api_grpc_queue_service_proto_depIdxs,
		MessageInfos:      file_api_grpc_queue_service_proto_msgTypes,
	}.Build()
	File_api_grpc_queue_service_proto = out.File
	file_api_grpc_queue_service_proto_rawDesc = nil
	file_api_grpc_queue_service_proto_goTypes = nil
	file_api_grpc_queue_service_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ServiceClient interface {
	// Creates a new queue or updates the existing one's length limit.
	SetQueue(ctx context.Context, in *SetQueueRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Start receiving a messages for the certain queue.
	ReceiveMessages(ctx context.Context, opts ...grpc.CallOption) (Service_ReceiveMessagesClient, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) SetQueue(ctx context.Context, in *SetQueueRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/awakari.queue.Service/SetQueue", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) ReceiveMessages(ctx context.Context, opts ...grpc.CallOption) (Service_ReceiveMessagesClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Service_serviceDesc.Streams[0], "/awakari.queue.Service/ReceiveMessages", opts...)
	if err != nil {
		return nil, err
	}
	x := &serviceReceiveMessagesClient{stream}
	return x, nil
}

type Service_ReceiveMessagesClient interface {
	Send(*ReceiveMessagesRequest) error
	Recv() (*ReceiveMessagesResponse, error)
	grpc.ClientStream
}

type serviceReceiveMessagesClient struct {
	grpc.ClientStream
}

func (x *serviceReceiveMessagesClient) Send(m *ReceiveMessagesRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *serviceReceiveMessagesClient) Recv() (*ReceiveMessagesResponse, error) {
	m := new(ReceiveMessagesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ServiceServer is the server API for Service service.
type ServiceServer interface {
	// Creates a new queue or updates the existing one's length limit.
	SetQueue(context.Context, *SetQueueRequest) (*emptypb.Empty, error)
	// Start receiving a messages for the certain queue.
	ReceiveMessages(Service_ReceiveMessagesServer) error
}

// UnimplementedServiceServer can be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (*UnimplementedServiceServer) SetQueue(context.Context, *SetQueueRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetQueue not implemented")
}
func (*UnimplementedServiceServer) ReceiveMessages(Service_ReceiveMessagesServer) error {
	return status.Errorf(codes.Unimplemented, "method ReceiveMessages not implemented")
}

func RegisterServiceServer(s *grpc.Server, srv ServiceServer) {
	s.RegisterService(&_Service_serviceDesc, srv)
}

func _Service_SetQueue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetQueueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).SetQueue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/awakari.queue.Service/SetQueue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).SetQueue(ctx, req.(*SetQueueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_ReceiveMessages_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ServiceServer).ReceiveMessages(&serviceReceiveMessagesServer{stream})
}

type Service_ReceiveMessagesServer interface {
	Send(*ReceiveMessagesResponse) error
	Recv() (*ReceiveMessagesRequest, error)
	grpc.ServerStream
}

type serviceReceiveMessagesServer struct {
	grpc.ServerStream
}

func (x *serviceReceiveMessagesServer) Send(m *ReceiveMessagesResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *serviceReceiveMessagesServer) Recv() (*ReceiveMessagesRequest, error) {
	m := new(ReceiveMessagesRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Service_serviceDesc = grpc.ServiceDesc{
	ServiceName: "awakari.queue.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetQueue",
			Handler:    _Service_SetQueue_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ReceiveMessages",
			Handler:       _Service_ReceiveMessages_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "api/grpc/queue/service.proto",
}
