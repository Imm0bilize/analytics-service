// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: api/proto/v1/analyticsService.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AnalyticsClient is the client API for Analytics service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AnalyticsClient interface {
	CreateTask(ctx context.Context, in *NewTask, opts ...grpc.CallOption) (*emptypb.Empty, error)
	SetTimeStart(ctx context.Context, in *TimeStart, opts ...grpc.CallOption) (*emptypb.Empty, error)
	SetTimeEnd(ctx context.Context, in *TimeEnd, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type analyticsClient struct {
	cc grpc.ClientConnInterface
}

func NewAnalyticsClient(cc grpc.ClientConnInterface) AnalyticsClient {
	return &analyticsClient{cc}
}

func (c *analyticsClient) CreateTask(ctx context.Context, in *NewTask, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/interpretedProto.Analytics/CreateTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *analyticsClient) SetTimeStart(ctx context.Context, in *TimeStart, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/interpretedProto.Analytics/SetTimeStart", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *analyticsClient) SetTimeEnd(ctx context.Context, in *TimeEnd, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/interpretedProto.Analytics/SetTimeEnd", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AnalyticsServer is the server API for Analytics service.
// All implementations must embed UnimplementedAnalyticsServer
// for forward compatibility
type AnalyticsServer interface {
	CreateTask(context.Context, *NewTask) (*emptypb.Empty, error)
	SetTimeStart(context.Context, *TimeStart) (*emptypb.Empty, error)
	SetTimeEnd(context.Context, *TimeEnd) (*emptypb.Empty, error)
	mustEmbedUnimplementedAnalyticsServer()
}

// UnimplementedAnalyticsServer must be embedded to have forward compatible implementations.
type UnimplementedAnalyticsServer struct {
}

func (UnimplementedAnalyticsServer) CreateTask(context.Context, *NewTask) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTask not implemented")
}
func (UnimplementedAnalyticsServer) SetTimeStart(context.Context, *TimeStart) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetTimeStart not implemented")
}
func (UnimplementedAnalyticsServer) SetTimeEnd(context.Context, *TimeEnd) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetTimeEnd not implemented")
}
func (UnimplementedAnalyticsServer) mustEmbedUnimplementedAnalyticsServer() {}

// UnsafeAnalyticsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AnalyticsServer will
// result in compilation errors.
type UnsafeAnalyticsServer interface {
	mustEmbedUnimplementedAnalyticsServer()
}

func RegisterAnalyticsServer(s grpc.ServiceRegistrar, srv AnalyticsServer) {
	s.RegisterService(&Analytics_ServiceDesc, srv)
}

func _Analytics_CreateTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewTask)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalyticsServer).CreateTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/interpretedProto.Analytics/CreateTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalyticsServer).CreateTask(ctx, req.(*NewTask))
	}
	return interceptor(ctx, in, info, handler)
}

func _Analytics_SetTimeStart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TimeStart)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalyticsServer).SetTimeStart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/interpretedProto.Analytics/SetTimeStart",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalyticsServer).SetTimeStart(ctx, req.(*TimeStart))
	}
	return interceptor(ctx, in, info, handler)
}

func _Analytics_SetTimeEnd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TimeEnd)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalyticsServer).SetTimeEnd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/interpretedProto.Analytics/SetTimeEnd",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalyticsServer).SetTimeEnd(ctx, req.(*TimeEnd))
	}
	return interceptor(ctx, in, info, handler)
}

// Analytics_ServiceDesc is the grpc.ServiceDesc for Analytics service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Analytics_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "interpretedProto.Analytics",
	HandlerType: (*AnalyticsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateTask",
			Handler:    _Analytics_CreateTask_Handler,
		},
		{
			MethodName: "SetTimeStart",
			Handler:    _Analytics_SetTimeStart_Handler,
		},
		{
			MethodName: "SetTimeEnd",
			Handler:    _Analytics_SetTimeEnd_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/v1/analyticsService.proto",
}