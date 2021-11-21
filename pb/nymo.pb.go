// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.12.4
// source: nymo.proto

package pb

import (
	reflect "reflect"
	sync "sync"

	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PeerHandshake struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cohort     uint32   `protobuf:"varint,1,opt,name=cohort,proto3" json:"cohort,omitempty"`
	Pow        []byte   `protobuf:"bytes,2,opt,name=pow,proto3" json:"pow,omitempty"`
	PeerTokens [][]byte `protobuf:"bytes,3,rep,name=peerTokens,proto3" json:"peerTokens,omitempty"`
}

func (x *PeerHandshake) Reset() {
	*x = PeerHandshake{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nymo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PeerHandshake) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PeerHandshake) ProtoMessage() {}

func (x *PeerHandshake) ProtoReflect() protoreflect.Message {
	mi := &file_nymo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PeerHandshake.ProtoReflect.Descriptor instead.
func (*PeerHandshake) Descriptor() ([]byte, []int) {
	return file_nymo_proto_rawDescGZIP(), []int{0}
}

func (x *PeerHandshake) GetCohort() uint32 {
	if x != nil {
		return x.Cohort
	}
	return 0
}

func (x *PeerHandshake) GetPow() []byte {
	if x != nil {
		return x.Pow
	}
	return nil
}

func (x *PeerHandshake) GetPeerTokens() [][]byte {
	if x != nil {
		return x.PeerTokens
	}
	return nil
}

type HandshakeOK struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cohort     uint32   `protobuf:"varint,1,opt,name=cohort,proto3" json:"cohort,omitempty"`
	PeerTokens [][]byte `protobuf:"bytes,2,rep,name=peerTokens,proto3" json:"peerTokens,omitempty"`
}

func (x *HandshakeOK) Reset() {
	*x = HandshakeOK{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nymo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HandshakeOK) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HandshakeOK) ProtoMessage() {}

func (x *HandshakeOK) ProtoReflect() protoreflect.Message {
	mi := &file_nymo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HandshakeOK.ProtoReflect.Descriptor instead.
func (*HandshakeOK) Descriptor() ([]byte, []int) {
	return file_nymo_proto_rawDescGZIP(), []int{1}
}

func (x *HandshakeOK) GetCohort() uint32 {
	if x != nil {
		return x.Cohort
	}
	return 0
}

func (x *HandshakeOK) GetPeerTokens() [][]byte {
	if x != nil {
		return x.PeerTokens
	}
	return nil
}

type MsgDigest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cohort uint32 `protobuf:"varint,1,opt,name=cohort,proto3" json:"cohort,omitempty"`
	Hash   []byte `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (x *MsgDigest) Reset() {
	*x = MsgDigest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nymo_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgDigest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgDigest) ProtoMessage() {}

func (x *MsgDigest) ProtoReflect() protoreflect.Message {
	mi := &file_nymo_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgDigest.ProtoReflect.Descriptor instead.
func (*MsgDigest) Descriptor() ([]byte, []int) {
	return file_nymo_proto_rawDescGZIP(), []int{2}
}

func (x *MsgDigest) GetCohort() uint32 {
	if x != nil {
		return x.Cohort
	}
	return 0
}

func (x *MsgDigest) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

type MsgList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Messages []*MsgDigest `protobuf:"bytes,1,rep,name=messages,proto3" json:"messages,omitempty"`
}

func (x *MsgList) Reset() {
	*x = MsgList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nymo_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgList) ProtoMessage() {}

func (x *MsgList) ProtoReflect() protoreflect.Message {
	mi := &file_nymo_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgList.ProtoReflect.Descriptor instead.
func (*MsgList) Descriptor() ([]byte, []int) {
	return file_nymo_proto_rawDescGZIP(), []int{3}
}

func (x *MsgList) GetMessages() []*MsgDigest {
	if x != nil {
		return x.Messages
	}
	return nil
}

type RequestMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash []byte `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (x *RequestMsg) Reset() {
	*x = RequestMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nymo_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestMsg) ProtoMessage() {}

func (x *RequestMsg) ProtoReflect() protoreflect.Message {
	mi := &file_nymo_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestMsg.ProtoReflect.Descriptor instead.
func (*RequestMsg) Descriptor() ([]byte, []int) {
	return file_nymo_proto_rawDescGZIP(), []int{4}
}

func (x *RequestMsg) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

type RequestPeer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token  []byte `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Pow    []byte `protobuf:"bytes,2,opt,name=pow,proto3" json:"pow,omitempty"`
	Cohort uint32 `protobuf:"varint,3,opt,name=cohort,proto3" json:"cohort,omitempty"`
}

func (x *RequestPeer) Reset() {
	*x = RequestPeer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nymo_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestPeer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestPeer) ProtoMessage() {}

func (x *RequestPeer) ProtoReflect() protoreflect.Message {
	mi := &file_nymo_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestPeer.ProtoReflect.Descriptor instead.
func (*RequestPeer) Descriptor() ([]byte, []int) {
	return file_nymo_proto_rawDescGZIP(), []int{5}
}

func (x *RequestPeer) GetToken() []byte {
	if x != nil {
		return x.Token
	}
	return nil
}

func (x *RequestPeer) GetPow() []byte {
	if x != nil {
		return x.Pow
	}
	return nil
}

func (x *RequestPeer) GetCohort() uint32 {
	if x != nil {
		return x.Cohort
	}
	return 0
}

type ResponsePeer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address []byte `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

func (x *ResponsePeer) Reset() {
	*x = ResponsePeer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nymo_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponsePeer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponsePeer) ProtoMessage() {}

func (x *ResponsePeer) ProtoReflect() protoreflect.Message {
	mi := &file_nymo_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponsePeer.ProtoReflect.Descriptor instead.
func (*ResponsePeer) Descriptor() ([]byte, []int) {
	return file_nymo_proto_rawDescGZIP(), []int{6}
}

func (x *ResponsePeer) GetAddress() []byte {
	if x != nil {
		return x.Address
	}
	return nil
}

type RealMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SendTime *timestamp.Timestamp `protobuf:"bytes,1,opt,name=sendTime,proto3" json:"sendTime,omitempty"`
	Message  string               `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	SenderID []byte               `protobuf:"bytes,3,opt,name=senderID,proto3" json:"senderID,omitempty"` // sender's public key
}

func (x *RealMessage) Reset() {
	*x = RealMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nymo_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RealMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RealMessage) ProtoMessage() {}

func (x *RealMessage) ProtoReflect() protoreflect.Message {
	mi := &file_nymo_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RealMessage.ProtoReflect.Descriptor instead.
func (*RealMessage) Descriptor() ([]byte, []int) {
	return file_nymo_proto_rawDescGZIP(), []int{7}
}

func (x *RealMessage) GetSendTime() *timestamp.Timestamp {
	if x != nil {
		return x.SendTime
	}
	return nil
}

func (x *RealMessage) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *RealMessage) GetSenderID() []byte {
	if x != nil {
		return x.SenderID
	}
	return nil
}

type EncryptedMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg          *RealMessage `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	MsgSignature []byte       `protobuf:"bytes,2,opt,name=msgSignature,proto3" json:"msgSignature,omitempty"` // signed with sender's private key
}

func (x *EncryptedMessage) Reset() {
	*x = EncryptedMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nymo_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EncryptedMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EncryptedMessage) ProtoMessage() {}

func (x *EncryptedMessage) ProtoReflect() protoreflect.Message {
	mi := &file_nymo_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EncryptedMessage.ProtoReflect.Descriptor instead.
func (*EncryptedMessage) Descriptor() ([]byte, []int) {
	return file_nymo_proto_rawDescGZIP(), []int{8}
}

func (x *EncryptedMessage) GetMsg() *RealMessage {
	if x != nil {
		return x.Msg
	}
	return nil
}

func (x *EncryptedMessage) GetMsgSignature() []byte {
	if x != nil {
		return x.MsgSignature
	}
	return nil
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TargetCohort uint32               `protobuf:"varint,1,opt,name=targetCohort,proto3" json:"targetCohort,omitempty"`
	Generated    *timestamp.Timestamp `protobuf:"bytes,2,opt,name=generated,proto3" json:"generated,omitempty"`
	EncAesKey    []byte               `protobuf:"bytes,3,opt,name=encAesKey,proto3" json:"encAesKey,omitempty"`   // receiver's public key encrypted AES ECB key
	EncMessage   []byte               `protobuf:"bytes,4,opt,name=encMessage,proto3" json:"encMessage,omitempty"` // aesKey encrypted EncryptedMessage
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nymo_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_nymo_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_nymo_proto_rawDescGZIP(), []int{9}
}

func (x *Message) GetTargetCohort() uint32 {
	if x != nil {
		return x.TargetCohort
	}
	return 0
}

func (x *Message) GetGenerated() *timestamp.Timestamp {
	if x != nil {
		return x.Generated
	}
	return nil
}

func (x *Message) GetEncAesKey() []byte {
	if x != nil {
		return x.EncAesKey
	}
	return nil
}

func (x *Message) GetEncMessage() []byte {
	if x != nil {
		return x.EncMessage
	}
	return nil
}

type MessageContainer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg     *Message `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	MsgHash []byte   `protobuf:"bytes,2,opt,name=msgHash,proto3" json:"msgHash,omitempty"`
	Pow     []byte   `protobuf:"bytes,3,opt,name=pow,proto3" json:"pow,omitempty"`
}

func (x *MessageContainer) Reset() {
	*x = MessageContainer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nymo_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageContainer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageContainer) ProtoMessage() {}

func (x *MessageContainer) ProtoReflect() protoreflect.Message {
	mi := &file_nymo_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageContainer.ProtoReflect.Descriptor instead.
func (*MessageContainer) Descriptor() ([]byte, []int) {
	return file_nymo_proto_rawDescGZIP(), []int{10}
}

func (x *MessageContainer) GetMsg() *Message {
	if x != nil {
		return x.Msg
	}
	return nil
}

func (x *MessageContainer) GetMsgHash() []byte {
	if x != nil {
		return x.MsgHash
	}
	return nil
}

func (x *MessageContainer) GetPow() []byte {
	if x != nil {
		return x.Pow
	}
	return nil
}

var File_nymo_proto protoreflect.FileDescriptor

var file_nymo_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x6e, 0x79, 0x6d, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62,
	0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x59, 0x0a, 0x0d, 0x50, 0x65, 0x65, 0x72, 0x48, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61,
	0x6b, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x6f, 0x68, 0x6f, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x06, 0x63, 0x6f, 0x68, 0x6f, 0x72, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x6f,
	0x77, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x70, 0x6f, 0x77, 0x12, 0x1e, 0x0a, 0x0a,
	0x70, 0x65, 0x65, 0x72, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0c,
	0x52, 0x0a, 0x70, 0x65, 0x65, 0x72, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x22, 0x45, 0x0a, 0x0b,
	0x48, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x4f, 0x4b, 0x12, 0x16, 0x0a, 0x06, 0x63,
	0x6f, 0x68, 0x6f, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x63, 0x6f, 0x68,
	0x6f, 0x72, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x65, 0x65, 0x72, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x0a, 0x70, 0x65, 0x65, 0x72, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x73, 0x22, 0x37, 0x0a, 0x09, 0x4d, 0x73, 0x67, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74,
	0x12, 0x16, 0x0a, 0x06, 0x63, 0x6f, 0x68, 0x6f, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x06, 0x63, 0x6f, 0x68, 0x6f, 0x72, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x22, 0x34, 0x0a, 0x07,
	0x4d, 0x73, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x29, 0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x62, 0x2e, 0x4d,
	0x73, 0x67, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x52, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x73, 0x22, 0x20, 0x0a, 0x0a, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4d, 0x73, 0x67,
	0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04,
	0x68, 0x61, 0x73, 0x68, 0x22, 0x4d, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x50,
	0x65, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x6f, 0x77,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x70, 0x6f, 0x77, 0x12, 0x16, 0x0a, 0x06, 0x63,
	0x6f, 0x68, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x63, 0x6f, 0x68,
	0x6f, 0x72, 0x74, 0x22, 0x28, 0x0a, 0x0c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x50,
	0x65, 0x65, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x7b, 0x0a,
	0x0b, 0x52, 0x65, 0x61, 0x6c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x36, 0x0a, 0x08,
	0x73, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x64,
	0x54, 0x69, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x44, 0x22, 0x59, 0x0a, 0x10, 0x45, 0x6e,
	0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x21,
	0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x62,
	0x2e, 0x52, 0x65, 0x61, 0x6c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x03, 0x6d, 0x73,
	0x67, 0x12, 0x22, 0x0a, 0x0c, 0x6d, 0x73, 0x67, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x6d, 0x73, 0x67, 0x53, 0x69, 0x67, 0x6e,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0xa5, 0x01, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x22, 0x0a, 0x0c, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x43, 0x6f, 0x68, 0x6f, 0x72,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0c, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x43,
	0x6f, 0x68, 0x6f, 0x72, 0x74, 0x12, 0x38, 0x0a, 0x09, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74,
	0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x12,
	0x1c, 0x0a, 0x09, 0x65, 0x6e, 0x63, 0x41, 0x65, 0x73, 0x4b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x09, 0x65, 0x6e, 0x63, 0x41, 0x65, 0x73, 0x4b, 0x65, 0x79, 0x12, 0x1e, 0x0a,
	0x0a, 0x65, 0x6e, 0x63, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x0a, 0x65, 0x6e, 0x63, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x5d, 0x0a,
	0x10, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65,
	0x72, 0x12, 0x1d, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b,
	0x2e, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x03, 0x6d, 0x73, 0x67,
	0x12, 0x18, 0x0a, 0x07, 0x6d, 0x73, 0x67, 0x48, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x07, 0x6d, 0x73, 0x67, 0x48, 0x61, 0x73, 0x68, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x6f,
	0x77, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x70, 0x6f, 0x77, 0x42, 0x1a, 0x5a, 0x18,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x79, 0x6d, 0x6f, 0x2d,
	0x6e, 0x65, 0x74, 0x2f, 0x6e, 0x79, 0x6d, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_nymo_proto_rawDescOnce sync.Once
	file_nymo_proto_rawDescData = file_nymo_proto_rawDesc
)

func file_nymo_proto_rawDescGZIP() []byte {
	file_nymo_proto_rawDescOnce.Do(func() {
		file_nymo_proto_rawDescData = protoimpl.X.CompressGZIP(file_nymo_proto_rawDescData)
	})
	return file_nymo_proto_rawDescData
}

var file_nymo_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_nymo_proto_goTypes = []interface{}{
	(*PeerHandshake)(nil),       // 0: pb.PeerHandshake
	(*HandshakeOK)(nil),         // 1: pb.HandshakeOK
	(*MsgDigest)(nil),           // 2: pb.MsgDigest
	(*MsgList)(nil),             // 3: pb.MsgList
	(*RequestMsg)(nil),          // 4: pb.RequestMsg
	(*RequestPeer)(nil),         // 5: pb.RequestPeer
	(*ResponsePeer)(nil),        // 6: pb.ResponsePeer
	(*RealMessage)(nil),         // 7: pb.RealMessage
	(*EncryptedMessage)(nil),    // 8: pb.EncryptedMessage
	(*Message)(nil),             // 9: pb.Message
	(*MessageContainer)(nil),    // 10: pb.MessageContainer
	(*timestamp.Timestamp)(nil), // 11: google.protobuf.Timestamp
}
var file_nymo_proto_depIdxs = []int32{
	2,  // 0: pb.MsgList.messages:type_name -> pb.MsgDigest
	11, // 1: pb.RealMessage.sendTime:type_name -> google.protobuf.Timestamp
	7,  // 2: pb.EncryptedMessage.msg:type_name -> pb.RealMessage
	11, // 3: pb.Message.generated:type_name -> google.protobuf.Timestamp
	9,  // 4: pb.MessageContainer.msg:type_name -> pb.Message
	5,  // [5:5] is the sub-list for method output_type
	5,  // [5:5] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_nymo_proto_init() }
func file_nymo_proto_init() {
	if File_nymo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_nymo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PeerHandshake); i {
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
		file_nymo_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HandshakeOK); i {
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
		file_nymo_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgDigest); i {
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
		file_nymo_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgList); i {
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
		file_nymo_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestMsg); i {
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
		file_nymo_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestPeer); i {
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
		file_nymo_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponsePeer); i {
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
		file_nymo_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RealMessage); i {
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
		file_nymo_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EncryptedMessage); i {
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
		file_nymo_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
		file_nymo_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageContainer); i {
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
			RawDescriptor: file_nymo_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_nymo_proto_goTypes,
		DependencyIndexes: file_nymo_proto_depIdxs,
		MessageInfos:      file_nymo_proto_msgTypes,
	}.Build()
	File_nymo_proto = out.File
	file_nymo_proto_rawDesc = nil
	file_nymo_proto_goTypes = nil
	file_nymo_proto_depIdxs = nil
}
