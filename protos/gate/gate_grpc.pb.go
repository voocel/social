// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.15.8
// source: protos/gate.proto

package gate

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

// GateClient is the client API for Gate service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GateClient interface {
	// 绑定用户与连接
	Bind(ctx context.Context, in *BindRequest, opts ...grpc.CallOption) (*BindReply, error)
	// 解绑用户与连接
	Unbind(ctx context.Context, in *UnbindRequest, opts ...grpc.CallOption) (*UnbindReply, error)
	// 推送消息
	Push(ctx context.Context, in *PushRequest, opts ...grpc.CallOption) (*PushReply, error)
}

type gateClient struct {
	cc grpc.ClientConnInterface
}

func NewGateClient(cc grpc.ClientConnInterface) GateClient {
	return &gateClient{cc}
}

func (c *gateClient) Bind(ctx context.Context, in *BindRequest, opts ...grpc.CallOption) (*BindReply, error) {
	out := new(BindReply)
	err := c.cc.Invoke(ctx, "/Gate.Gate/Bind", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gateClient) Unbind(ctx context.Context, in *UnbindRequest, opts ...grpc.CallOption) (*UnbindReply, error) {
	out := new(UnbindReply)
	err := c.cc.Invoke(ctx, "/Gate.Gate/Unbind", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gateClient) Push(ctx context.Context, in *PushRequest, opts ...grpc.CallOption) (*PushReply, error) {
	out := new(PushReply)
	err := c.cc.Invoke(ctx, "/Gate.Gate/Push", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GateServer is the server API for Gate service.
// All implementations must embed UnimplementedGateServer
// for forward compatibility
type GateServer interface {
	// 绑定用户与连接
	Bind(context.Context, *BindRequest) (*BindReply, error)
	// 解绑用户与连接
	Unbind(context.Context, *UnbindRequest) (*UnbindReply, error)
	// 推送消息
	Push(context.Context, *PushRequest) (*PushReply, error)
	mustEmbedUnimplementedGateServer()
}

// UnimplementedGateServer must be embedded to have forward compatible implementations.
type UnimplementedGateServer struct {
}

func (UnimplementedGateServer) Bind(context.Context, *BindRequest) (*BindReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Bind not implemented")
}
func (UnimplementedGateServer) Unbind(context.Context, *UnbindRequest) (*UnbindReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unbind not implemented")
}
func (UnimplementedGateServer) Push(context.Context, *PushRequest) (*PushReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Push not implemented")
}
func (UnimplementedGateServer) mustEmbedUnimplementedGateServer() {}

// UnsafeGateServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GateServer will
// result in compilation errors.
type UnsafeGateServer interface {
	mustEmbedUnimplementedGateServer()
}

func RegisterGateServer(s grpc.ServiceRegistrar, srv GateServer) {
	s.RegisterService(&Gate_ServiceDesc, srv)
}

func _Gate_Bind_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BindRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GateServer).Bind(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Gate.Gate/Bind",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GateServer).Bind(ctx, req.(*BindRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gate_Unbind_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnbindRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GateServer).Unbind(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Gate.Gate/Unbind",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GateServer).Unbind(ctx, req.(*UnbindRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gate_Push_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GateServer).Push(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Gate.Gate/Push",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GateServer).Push(ctx, req.(*PushRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Gate_ServiceDesc is the grpc.ServiceDesc for Gate service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Gate_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Gate.Gate",
	HandlerType: (*GateServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Bind",
			Handler:    _Gate_Bind_Handler,
		},
		{
			MethodName: "Unbind",
			Handler:    _Gate_Unbind_Handler,
		},
		{
			MethodName: "Push",
			Handler:    _Gate_Push_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/gate.proto",
}