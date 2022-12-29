// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

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

// TemplateClient is the client API for Template service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TemplateClient interface {
	// one message is sent and one is recieved
	Increment(ctx context.Context, in *Amount, opts ...grpc.CallOption) (*Ack, error)
	// Sends heartbeat with current amount so secondary servers can update
	ReceiveHeartbeat(ctx context.Context, in *Amount, opts ...grpc.CallOption) (*Ack, error)
	SendHeartbeat(ctx context.Context, in *Amount, opts ...grpc.CallOption) (*Ack, error)
}

type templateClient struct {
	cc grpc.ClientConnInterface
}

func NewTemplateClient(cc grpc.ClientConnInterface) TemplateClient {
	return &templateClient{cc}
}

func (c *templateClient) Increment(ctx context.Context, in *Amount, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/proto.Template/Increment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *templateClient) ReceiveHeartbeat(ctx context.Context, in *Amount, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/proto.Template/ReceiveHeartbeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *templateClient) SendHeartbeat(ctx context.Context, in *Amount, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/proto.Template/SendHeartbeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TemplateServer is the server API for Template service.
// All implementations must embed UnimplementedTemplateServer
// for forward compatibility
type TemplateServer interface {
	// one message is sent and one is recieved
	Increment(context.Context, *Amount) (*Ack, error)
	// Sends heartbeat with current amount so secondary servers can update
	ReceiveHeartbeat(context.Context, *Amount) (*Ack, error)
	SendHeartbeat(context.Context, *Amount) (*Ack, error)
	mustEmbedUnimplementedTemplateServer()
}

// UnimplementedTemplateServer must be embedded to have forward compatible implementations.
type UnimplementedTemplateServer struct {
}

func (UnimplementedTemplateServer) Increment(context.Context, *Amount) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Increment not implemented")
}
func (UnimplementedTemplateServer) ReceiveHeartbeat(context.Context, *Amount) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReceiveHeartbeat not implemented")
}
func (UnimplementedTemplateServer) SendHeartbeat(context.Context, *Amount) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendHeartbeat not implemented")
}
func (UnimplementedTemplateServer) mustEmbedUnimplementedTemplateServer() {}

// UnsafeTemplateServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TemplateServer will
// result in compilation errors.
type UnsafeTemplateServer interface {
	mustEmbedUnimplementedTemplateServer()
}

func RegisterTemplateServer(s grpc.ServiceRegistrar, srv TemplateServer) {
	s.RegisterService(&Template_ServiceDesc, srv)
}

func _Template_Increment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Amount)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TemplateServer).Increment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Template/Increment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TemplateServer).Increment(ctx, req.(*Amount))
	}
	return interceptor(ctx, in, info, handler)
}

func _Template_ReceiveHeartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Amount)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TemplateServer).ReceiveHeartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Template/ReceiveHeartbeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TemplateServer).ReceiveHeartbeat(ctx, req.(*Amount))
	}
	return interceptor(ctx, in, info, handler)
}

func _Template_SendHeartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Amount)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TemplateServer).SendHeartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Template/SendHeartbeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TemplateServer).SendHeartbeat(ctx, req.(*Amount))
	}
	return interceptor(ctx, in, info, handler)
}

// Template_ServiceDesc is the grpc.ServiceDesc for Template service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Template_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Template",
	HandlerType: (*TemplateServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Increment",
			Handler:    _Template_Increment_Handler,
		},
		{
			MethodName: "ReceiveHeartbeat",
			Handler:    _Template_ReceiveHeartbeat_Handler,
		},
		{
			MethodName: "SendHeartbeat",
			Handler:    _Template_SendHeartbeat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/template.proto",
}
