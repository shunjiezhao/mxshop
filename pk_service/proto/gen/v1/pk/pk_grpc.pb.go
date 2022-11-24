// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: pk.proto

package proto

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

// PKClient is the client API for PK service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PKClient interface {
	TakePartIn(ctx context.Context, in *TakePartInRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// 用户参加 -> 返回题目
	Join(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*JoinResponse, error)
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
}

type pKClient struct {
	cc grpc.ClientConnInterface
}

func NewPKClient(cc grpc.ClientConnInterface) PKClient {
	return &pKClient{cc}
}

func (c *pKClient) TakePartIn(ctx context.Context, in *TakePartInRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/PK/TakePartIn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pKClient) Join(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*JoinResponse, error) {
	out := new(JoinResponse)
	err := c.cc.Invoke(ctx, "/PK/Join", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pKClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, "/PK/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PKServer is the server API for PK service.
// All implementations must embed UnimplementedPKServer
// for forward compatibility
type PKServer interface {
	TakePartIn(context.Context, *TakePartInRequest) (*emptypb.Empty, error)
	// 用户参加 -> 返回题目
	Join(context.Context, *JoinRequest) (*JoinResponse, error)
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	mustEmbedUnimplementedPKServer()
}

// UnimplementedPKServer must be embedded to have forward compatible implementations.
type UnimplementedPKServer struct {
}

func (UnimplementedPKServer) TakePartIn(context.Context, *TakePartInRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TakePartIn not implemented")
}
func (UnimplementedPKServer) Join(context.Context, *JoinRequest) (*JoinResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Join not implemented")
}
func (UnimplementedPKServer) Create(context.Context, *CreateRequest) (*CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedPKServer) mustEmbedUnimplementedPKServer() {}

// UnsafePKServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PKServer will
// result in compilation errors.
type UnsafePKServer interface {
	mustEmbedUnimplementedPKServer()
}

func RegisterPKServer(s grpc.ServiceRegistrar, srv PKServer) {
	s.RegisterService(&PK_ServiceDesc, srv)
}

func _PK_TakePartIn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TakePartInRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PKServer).TakePartIn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/PK/TakePartIn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PKServer).TakePartIn(ctx, req.(*TakePartInRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PK_Join_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PKServer).Join(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/PK/Join",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PKServer).Join(ctx, req.(*JoinRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PK_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PKServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/PK/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PKServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PK_ServiceDesc is the grpc.ServiceDesc for PK service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PK_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "PK",
	HandlerType: (*PKServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TakePartIn",
			Handler:    _PK_TakePartIn_Handler,
		},
		{
			MethodName: "Join",
			Handler:    _PK_Join_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _PK_Create_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pk.proto",
}
