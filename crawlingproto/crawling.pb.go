// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: crawlingproto/crawling.proto

package crawlingproto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type UserInput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID string `protobuf:"bytes,1,opt,name=userID,proto3" json:"userID,omitempty"`
}

func (x *UserInput) Reset() {
	*x = UserInput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_crawlingproto_crawling_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserInput) ProtoMessage() {}

func (x *UserInput) ProtoReflect() protoreflect.Message {
	mi := &file_crawlingproto_crawling_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserInput.ProtoReflect.Descriptor instead.
func (*UserInput) Descriptor() ([]byte, []int) {
	return file_crawlingproto_crawling_proto_rawDescGZIP(), []int{0}
}

func (x *UserInput) GetUserID() string {
	if x != nil {
		return x.UserID
	}
	return ""
}

type BankInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BankSum    int64  `protobuf:"varint,1,opt,name=bank_sum,json=bankSum,proto3" json:"bank_sum,omitempty"`
	BankCount  int64  `protobuf:"varint,2,opt,name=bank_count,json=bankCount,proto3" json:"bank_count,omitempty"`
	BankName   string `protobuf:"bytes,3,opt,name=bank_name,json=bankName,proto3" json:"bank_name,omitempty"`
	BankAmount int64  `protobuf:"varint,4,opt,name=bank_amount,json=bankAmount,proto3" json:"bank_amount,omitempty"`
}

func (x *BankInfo) Reset() {
	*x = BankInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_crawlingproto_crawling_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BankInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BankInfo) ProtoMessage() {}

func (x *BankInfo) ProtoReflect() protoreflect.Message {
	mi := &file_crawlingproto_crawling_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BankInfo.ProtoReflect.Descriptor instead.
func (*BankInfo) Descriptor() ([]byte, []int) {
	return file_crawlingproto_crawling_proto_rawDescGZIP(), []int{1}
}

func (x *BankInfo) GetBankSum() int64 {
	if x != nil {
		return x.BankSum
	}
	return 0
}

func (x *BankInfo) GetBankCount() int64 {
	if x != nil {
		return x.BankCount
	}
	return 0
}

func (x *BankInfo) GetBankName() string {
	if x != nil {
		return x.BankName
	}
	return ""
}

func (x *BankInfo) GetBankAmount() int64 {
	if x != nil {
		return x.BankAmount
	}
	return 0
}

type CardInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CardSum    int64  `protobuf:"varint,1,opt,name=card_sum,json=cardSum,proto3" json:"card_sum,omitempty"`
	CardCount  int64  `protobuf:"varint,2,opt,name=card_count,json=cardCount,proto3" json:"card_count,omitempty"`
	CardName   string `protobuf:"bytes,3,opt,name=card_name,json=cardName,proto3" json:"card_name,omitempty"`
	CardAmount int64  `protobuf:"varint,4,opt,name=card_amount,json=cardAmount,proto3" json:"card_amount,omitempty"`
}

func (x *CardInfo) Reset() {
	*x = CardInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_crawlingproto_crawling_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CardInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CardInfo) ProtoMessage() {}

func (x *CardInfo) ProtoReflect() protoreflect.Message {
	mi := &file_crawlingproto_crawling_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CardInfo.ProtoReflect.Descriptor instead.
func (*CardInfo) Descriptor() ([]byte, []int) {
	return file_crawlingproto_crawling_proto_rawDescGZIP(), []int{2}
}

func (x *CardInfo) GetCardSum() int64 {
	if x != nil {
		return x.CardSum
	}
	return 0
}

func (x *CardInfo) GetCardCount() int64 {
	if x != nil {
		return x.CardCount
	}
	return 0
}

func (x *CardInfo) GetCardName() string {
	if x != nil {
		return x.CardName
	}
	return ""
}

func (x *CardInfo) GetCardAmount() int64 {
	if x != nil {
		return x.CardAmount
	}
	return 0
}

type DetailInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DetailCount int64  `protobuf:"varint,1,opt,name=detail_count,json=detailCount,proto3" json:"detail_count,omitempty"`
	DetailName  string `protobuf:"bytes,2,opt,name=detail_name,json=detailName,proto3" json:"detail_name,omitempty"`
	Contents    string `protobuf:"bytes,3,opt,name=contents,proto3" json:"contents,omitempty"`
	PaymentType int32  `protobuf:"varint,4,opt,name=payment_type,json=paymentType,proto3" json:"payment_type,omitempty"`
	Amount      int64  `protobuf:"varint,5,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (x *DetailInfo) Reset() {
	*x = DetailInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_crawlingproto_crawling_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DetailInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DetailInfo) ProtoMessage() {}

func (x *DetailInfo) ProtoReflect() protoreflect.Message {
	mi := &file_crawlingproto_crawling_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DetailInfo.ProtoReflect.Descriptor instead.
func (*DetailInfo) Descriptor() ([]byte, []int) {
	return file_crawlingproto_crawling_proto_rawDescGZIP(), []int{3}
}

func (x *DetailInfo) GetDetailCount() int64 {
	if x != nil {
		return x.DetailCount
	}
	return 0
}

func (x *DetailInfo) GetDetailName() string {
	if x != nil {
		return x.DetailName
	}
	return ""
}

func (x *DetailInfo) GetContents() string {
	if x != nil {
		return x.Contents
	}
	return ""
}

func (x *DetailInfo) GetPaymentType() int32 {
	if x != nil {
		return x.PaymentType
	}
	return 0
}

func (x *DetailInfo) GetAmount() int64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

type UserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserInput *UserInput `protobuf:"bytes,1,opt,name=userInput,proto3" json:"userInput,omitempty"`
	Pass      string     `protobuf:"bytes,2,opt,name=pass,proto3" json:"pass,omitempty"`
	SiteKind  int32      `protobuf:"varint,3,opt,name=site_kind,json=siteKind,proto3" json:"site_kind,omitempty"`
}

func (x *UserRequest) Reset() {
	*x = UserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_crawlingproto_crawling_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserRequest) ProtoMessage() {}

func (x *UserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_crawlingproto_crawling_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserRequest.ProtoReflect.Descriptor instead.
func (*UserRequest) Descriptor() ([]byte, []int) {
	return file_crawlingproto_crawling_proto_rawDescGZIP(), []int{4}
}

func (x *UserRequest) GetUserInput() *UserInput {
	if x != nil {
		return x.UserInput
	}
	return nil
}

func (x *UserRequest) GetPass() string {
	if x != nil {
		return x.Pass
	}
	return ""
}

func (x *UserRequest) GetSiteKind() int32 {
	if x != nil {
		return x.SiteKind
	}
	return 0
}

type ApiResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BankInfo   *BankInfo   `protobuf:"bytes,1,opt,name=bankInfo,proto3" json:"bankInfo,omitempty"`
	CardInfo   *CardInfo   `protobuf:"bytes,2,opt,name=cardInfo,proto3" json:"cardInfo,omitempty"`
	DetailInfo *DetailInfo `protobuf:"bytes,3,opt,name=detailInfo,proto3" json:"detailInfo,omitempty"`
}

func (x *ApiResponse) Reset() {
	*x = ApiResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_crawlingproto_crawling_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApiResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApiResponse) ProtoMessage() {}

func (x *ApiResponse) ProtoReflect() protoreflect.Message {
	mi := &file_crawlingproto_crawling_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApiResponse.ProtoReflect.Descriptor instead.
func (*ApiResponse) Descriptor() ([]byte, []int) {
	return file_crawlingproto_crawling_proto_rawDescGZIP(), []int{5}
}

func (x *ApiResponse) GetBankInfo() *BankInfo {
	if x != nil {
		return x.BankInfo
	}
	return nil
}

func (x *ApiResponse) GetCardInfo() *CardInfo {
	if x != nil {
		return x.CardInfo
	}
	return nil
}

func (x *ApiResponse) GetDetailInfo() *DetailInfo {
	if x != nil {
		return x.DetailInfo
	}
	return nil
}

var File_crawlingproto_crawling_proto protoreflect.FileDescriptor

var file_crawlingproto_crawling_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x63, 0x72, 0x61, 0x77, 0x6c, 0x69, 0x6e, 0x67, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x63, 0x72, 0x61, 0x77, 0x6c, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d,
	0x63, 0x72, 0x61, 0x77, 0x6c, 0x69, 0x6e, 0x67, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x23, 0x0a,
	0x09, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x44, 0x22, 0x82, 0x01, 0x0a, 0x08, 0x42, 0x61, 0x6e, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x12,
	0x19, 0x0a, 0x08, 0x62, 0x61, 0x6e, 0x6b, 0x5f, 0x73, 0x75, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x07, 0x62, 0x61, 0x6e, 0x6b, 0x53, 0x75, 0x6d, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x61,
	0x6e, 0x6b, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x62, 0x61, 0x6e, 0x6b, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x62, 0x61, 0x6e,
	0x6b, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x62, 0x61,
	0x6e, 0x6b, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x62, 0x61, 0x6e, 0x6b, 0x5f, 0x61,
	0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x62, 0x61, 0x6e,
	0x6b, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x82, 0x01, 0x0a, 0x08, 0x43, 0x61, 0x72, 0x64,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x19, 0x0a, 0x08, 0x63, 0x61, 0x72, 0x64, 0x5f, 0x73, 0x75, 0x6d,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x63, 0x61, 0x72, 0x64, 0x53, 0x75, 0x6d, 0x12,
	0x1d, 0x0a, 0x0a, 0x63, 0x61, 0x72, 0x64, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x61, 0x72, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1b,
	0x0a, 0x09, 0x63, 0x61, 0x72, 0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x63, 0x61, 0x72, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x63,
	0x61, 0x72, 0x64, 0x5f, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x0a, 0x63, 0x61, 0x72, 0x64, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0xa7, 0x01, 0x0a,
	0x0a, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x21, 0x0a, 0x0c, 0x64,
	0x65, 0x74, 0x61, 0x69, 0x6c, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x0b, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1f,
	0x0a, 0x0b, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x70,
	0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0b, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06,
	0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x76, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x36, 0x0a, 0x09, 0x75, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x70,
	0x75, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x63, 0x72, 0x61, 0x77, 0x6c,
	0x69, 0x6e, 0x67, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x70,
	0x75, 0x74, 0x52, 0x09, 0x75, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x70, 0x61, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x73,
	0x73, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x69, 0x74, 0x65, 0x5f, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x73, 0x69, 0x74, 0x65, 0x4b, 0x69, 0x6e, 0x64, 0x22, 0xb2,
	0x01, 0x0a, 0x0b, 0x41, 0x70, 0x69, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x33,
	0x0a, 0x08, 0x62, 0x61, 0x6e, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x17, 0x2e, 0x63, 0x72, 0x61, 0x77, 0x6c, 0x69, 0x6e, 0x67, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x42, 0x61, 0x6e, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x08, 0x62, 0x61, 0x6e, 0x6b, 0x49,
	0x6e, 0x66, 0x6f, 0x12, 0x33, 0x0a, 0x08, 0x63, 0x61, 0x72, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x63, 0x72, 0x61, 0x77, 0x6c, 0x69, 0x6e, 0x67,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x61, 0x72, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x08,
	0x63, 0x61, 0x72, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x39, 0x0a, 0x0a, 0x64, 0x65, 0x74, 0x61,
	0x69, 0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63,
	0x72, 0x61, 0x77, 0x6c, 0x69, 0x6e, 0x67, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x65, 0x74,
	0x61, 0x69, 0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0a, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x49,
	0x6e, 0x66, 0x6f, 0x32, 0x5a, 0x0a, 0x0f, 0x43, 0x72, 0x61, 0x77, 0x6c, 0x69, 0x6e, 0x67, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x47, 0x0a, 0x0b, 0x75, 0x73, 0x65, 0x72, 0x48, 0x61,
	0x6e, 0x64, 0x6c, 0x65, 0x72, 0x12, 0x1a, 0x2e, 0x63, 0x72, 0x61, 0x77, 0x6c, 0x69, 0x6e, 0x67,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1a, 0x2e, 0x63, 0x72, 0x61, 0x77, 0x6c, 0x69, 0x6e, 0x67, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x41, 0x70, 0x69, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42,
	0x10, 0x5a, 0x0e, 0x63, 0x72, 0x61, 0x77, 0x6c, 0x69, 0x6e, 0x67, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_crawlingproto_crawling_proto_rawDescOnce sync.Once
	file_crawlingproto_crawling_proto_rawDescData = file_crawlingproto_crawling_proto_rawDesc
)

func file_crawlingproto_crawling_proto_rawDescGZIP() []byte {
	file_crawlingproto_crawling_proto_rawDescOnce.Do(func() {
		file_crawlingproto_crawling_proto_rawDescData = protoimpl.X.CompressGZIP(file_crawlingproto_crawling_proto_rawDescData)
	})
	return file_crawlingproto_crawling_proto_rawDescData
}

var file_crawlingproto_crawling_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_crawlingproto_crawling_proto_goTypes = []interface{}{
	(*UserInput)(nil),   // 0: crawlingproto.UserInput
	(*BankInfo)(nil),    // 1: crawlingproto.BankInfo
	(*CardInfo)(nil),    // 2: crawlingproto.CardInfo
	(*DetailInfo)(nil),  // 3: crawlingproto.DetailInfo
	(*UserRequest)(nil), // 4: crawlingproto.UserRequest
	(*ApiResponse)(nil), // 5: crawlingproto.ApiResponse
}
var file_crawlingproto_crawling_proto_depIdxs = []int32{
	0, // 0: crawlingproto.UserRequest.userInput:type_name -> crawlingproto.UserInput
	1, // 1: crawlingproto.ApiResponse.bankInfo:type_name -> crawlingproto.BankInfo
	2, // 2: crawlingproto.ApiResponse.cardInfo:type_name -> crawlingproto.CardInfo
	3, // 3: crawlingproto.ApiResponse.detailInfo:type_name -> crawlingproto.DetailInfo
	4, // 4: crawlingproto.CrawlingService.userHandler:input_type -> crawlingproto.UserRequest
	5, // 5: crawlingproto.CrawlingService.userHandler:output_type -> crawlingproto.ApiResponse
	5, // [5:6] is the sub-list for method output_type
	4, // [4:5] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_crawlingproto_crawling_proto_init() }
func file_crawlingproto_crawling_proto_init() {
	if File_crawlingproto_crawling_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_crawlingproto_crawling_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserInput); i {
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
		file_crawlingproto_crawling_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BankInfo); i {
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
		file_crawlingproto_crawling_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CardInfo); i {
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
		file_crawlingproto_crawling_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DetailInfo); i {
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
		file_crawlingproto_crawling_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserRequest); i {
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
		file_crawlingproto_crawling_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApiResponse); i {
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
			RawDescriptor: file_crawlingproto_crawling_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_crawlingproto_crawling_proto_goTypes,
		DependencyIndexes: file_crawlingproto_crawling_proto_depIdxs,
		MessageInfos:      file_crawlingproto_crawling_proto_msgTypes,
	}.Build()
	File_crawlingproto_crawling_proto = out.File
	file_crawlingproto_crawling_proto_rawDesc = nil
	file_crawlingproto_crawling_proto_goTypes = nil
	file_crawlingproto_crawling_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// CrawlingServiceClient is the client API for CrawlingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CrawlingServiceClient interface {
	UserHandler(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*ApiResponse, error)
}

type crawlingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCrawlingServiceClient(cc grpc.ClientConnInterface) CrawlingServiceClient {
	return &crawlingServiceClient{cc}
}

func (c *crawlingServiceClient) UserHandler(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*ApiResponse, error) {
	out := new(ApiResponse)
	err := c.cc.Invoke(ctx, "/crawlingproto.CrawlingService/userHandler", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CrawlingServiceServer is the server API for CrawlingService service.
type CrawlingServiceServer interface {
	UserHandler(context.Context, *UserRequest) (*ApiResponse, error)
}

// UnimplementedCrawlingServiceServer can be embedded to have forward compatible implementations.
type UnimplementedCrawlingServiceServer struct {
}

func (*UnimplementedCrawlingServiceServer) UserHandler(context.Context, *UserRequest) (*ApiResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserHandler not implemented")
}

func RegisterCrawlingServiceServer(s *grpc.Server, srv CrawlingServiceServer) {
	s.RegisterService(&_CrawlingService_serviceDesc, srv)
}

func _CrawlingService_UserHandler_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CrawlingServiceServer).UserHandler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/crawlingproto.CrawlingService/UserHandler",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrawlingServiceServer).UserHandler(ctx, req.(*UserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _CrawlingService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "crawlingproto.CrawlingService",
	HandlerType: (*CrawlingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "userHandler",
			Handler:    _CrawlingService_UserHandler_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "crawlingproto/crawling.proto",
}