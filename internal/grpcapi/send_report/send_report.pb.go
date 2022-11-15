// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: send_report.proto

package send_report

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SendReportRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Payload  string                 `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	UserID   int64                  `protobuf:"varint,2,opt,name=userID,proto3" json:"userID,omitempty"`
	Currency string                 `protobuf:"bytes,3,opt,name=currency,proto3" json:"currency,omitempty"`
	Period   int32                  `protobuf:"varint,4,opt,name=period,proto3" json:"period,omitempty"`
	Format   string                 `protobuf:"bytes,5,opt,name=format,proto3" json:"format,omitempty"`
	Date     *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=date,proto3" json:"date,omitempty"`
}

func (x *SendReportRequest) Reset() {
	*x = SendReportRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_send_report_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendReportRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendReportRequest) ProtoMessage() {}

func (x *SendReportRequest) ProtoReflect() protoreflect.Message {
	mi := &file_send_report_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendReportRequest.ProtoReflect.Descriptor instead.
func (*SendReportRequest) Descriptor() ([]byte, []int) {
	return file_send_report_proto_rawDescGZIP(), []int{0}
}

func (x *SendReportRequest) GetPayload() string {
	if x != nil {
		return x.Payload
	}
	return ""
}

func (x *SendReportRequest) GetUserID() int64 {
	if x != nil {
		return x.UserID
	}
	return 0
}

func (x *SendReportRequest) GetCurrency() string {
	if x != nil {
		return x.Currency
	}
	return ""
}

func (x *SendReportRequest) GetPeriod() int32 {
	if x != nil {
		return x.Period
	}
	return 0
}

func (x *SendReportRequest) GetFormat() string {
	if x != nil {
		return x.Format
	}
	return ""
}

func (x *SendReportRequest) GetDate() *timestamppb.Timestamp {
	if x != nil {
		return x.Date
	}
	return nil
}

type SendReportResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SendReportResponse) Reset() {
	*x = SendReportResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_send_report_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendReportResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendReportResponse) ProtoMessage() {}

func (x *SendReportResponse) ProtoReflect() protoreflect.Message {
	mi := &file_send_report_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendReportResponse.ProtoReflect.Descriptor instead.
func (*SendReportResponse) Descriptor() ([]byte, []int) {
	return file_send_report_proto_rawDescGZIP(), []int{1}
}

var File_send_report_proto protoreflect.FileDescriptor

var file_send_report_proto_rawDesc = []byte{
	0x0a, 0x11, 0x73, 0x65, 0x6e, 0x64, 0x5f, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x73, 0x65, 0x6e, 0x64, 0x5f, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74,
	0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xc1, 0x01, 0x0a, 0x11, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f,
	0x61, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x75, 0x72,
	0x72, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x75, 0x72,
	0x72, 0x65, 0x6e, 0x63, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x70, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x12, 0x16, 0x0a,
	0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66,
	0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x2e, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x65, 0x22, 0x14, 0x0a, 0x12, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x70,
	0x6f, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x5f, 0x0a, 0x0c, 0x52,
	0x65, 0x70, 0x6f, 0x72, 0x74, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x4f, 0x0a, 0x0a, 0x53,
	0x65, 0x6e, 0x64, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x1e, 0x2e, 0x73, 0x65, 0x6e, 0x64,
	0x5f, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x73, 0x65, 0x6e, 0x64,
	0x5f, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x4c, 0x5a, 0x4a,
	0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x6f, 0x7a, 0x6f, 0x6e, 0x2e, 0x64, 0x65, 0x76, 0x2f,
	0x65, 0x67, 0x6f, 0x72, 0x2e, 0x6c, 0x69, 0x6e, 0x6b, 0x69, 0x6e, 0x6b, 0x65, 0x64, 0x2f, 0x6b,
	0x61, 0x72, 0x74, 0x61, 0x73, 0x68, 0x6f, 0x76, 0x2d, 0x65, 0x67, 0x6f, 0x72, 0x2f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x61, 0x70, 0x69, 0x2f, 0x73,
	0x65, 0x6e, 0x64, 0x5f, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_send_report_proto_rawDescOnce sync.Once
	file_send_report_proto_rawDescData = file_send_report_proto_rawDesc
)

func file_send_report_proto_rawDescGZIP() []byte {
	file_send_report_proto_rawDescOnce.Do(func() {
		file_send_report_proto_rawDescData = protoimpl.X.CompressGZIP(file_send_report_proto_rawDescData)
	})
	return file_send_report_proto_rawDescData
}

var file_send_report_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_send_report_proto_goTypes = []interface{}{
	(*SendReportRequest)(nil),     // 0: send_report.SendReportRequest
	(*SendReportResponse)(nil),    // 1: send_report.SendReportResponse
	(*timestamppb.Timestamp)(nil), // 2: google.protobuf.Timestamp
}
var file_send_report_proto_depIdxs = []int32{
	2, // 0: send_report.SendReportRequest.date:type_name -> google.protobuf.Timestamp
	0, // 1: send_report.ReportSender.SendReport:input_type -> send_report.SendReportRequest
	1, // 2: send_report.ReportSender.SendReport:output_type -> send_report.SendReportResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_send_report_proto_init() }
func file_send_report_proto_init() {
	if File_send_report_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_send_report_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendReportRequest); i {
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
		file_send_report_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendReportResponse); i {
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
			RawDescriptor: file_send_report_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_send_report_proto_goTypes,
		DependencyIndexes: file_send_report_proto_depIdxs,
		MessageInfos:      file_send_report_proto_msgTypes,
	}.Build()
	File_send_report_proto = out.File
	file_send_report_proto_rawDesc = nil
	file_send_report_proto_goTypes = nil
	file_send_report_proto_depIdxs = nil
}
