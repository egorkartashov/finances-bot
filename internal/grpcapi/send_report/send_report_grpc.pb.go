// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: send_report.proto

package send_report

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ReportSenderClient is the client API for ReportSender service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReportSenderClient interface {
	SendReport(ctx context.Context, in *SendReportRequest, opts ...grpc.CallOption) (*SendReportResponse, error)
}

type reportSenderClient struct {
	cc grpc.ClientConnInterface
}

func NewReportSenderClient(cc grpc.ClientConnInterface) ReportSenderClient {
	return &reportSenderClient{cc}
}

func (c *reportSenderClient) SendReport(ctx context.Context, in *SendReportRequest, opts ...grpc.CallOption) (*SendReportResponse, error) {
	out := new(SendReportResponse)
	err := c.cc.Invoke(ctx, "/send_report.ReportSender/SendReport", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReportSenderServer is the server API for ReportSender service.
// All implementations must embed UnimplementedReportSenderServer
// for forward compatibility
type ReportSenderServer interface {
	SendReport(context.Context, *SendReportRequest) (*SendReportResponse, error)
	mustEmbedUnimplementedReportSenderServer()
}

// UnimplementedReportSenderServer must be embedded to have forward compatible implementations.
type UnimplementedReportSenderServer struct {
}

func (UnimplementedReportSenderServer) SendReport(context.Context, *SendReportRequest) (*SendReportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendReport not implemented")
}
func (UnimplementedReportSenderServer) mustEmbedUnimplementedReportSenderServer() {}

// UnsafeReportSenderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReportSenderServer will
// result in compilation errors.
type UnsafeReportSenderServer interface {
	mustEmbedUnimplementedReportSenderServer()
}

func RegisterReportSenderServer(s grpc.ServiceRegistrar, srv ReportSenderServer) {
	s.RegisterService(&ReportSender_ServiceDesc, srv)
}

func _ReportSender_SendReport_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReportSenderServer).SendReport(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/send_report.ReportSender/SendReport",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReportSenderServer).SendReport(ctx, req.(*SendReportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ReportSender_ServiceDesc is the grpc.ServiceDesc for ReportSender service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ReportSender_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "send_report.ReportSender",
	HandlerType: (*ReportSenderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendReport",
			Handler:    _ReportSender_SendReport_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "send_report.proto",
}
