// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: pkg/proto/services.proto

package shorturl

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

const (
	Services_Stats_FullMethodName                  = "/shorturl.Services/Stats"
	Services_DeleteHandler_FullMethodName          = "/shorturl.Services/DeleteHandler"
	Services_GetAllURL_FullMethodName              = "/shorturl.Services/GetAllURL"
	Services_CreateShortURLHandler_FullMethodName  = "/shorturl.Services/CreateShortURLHandler"
	Services_GetShortURLByIDHandler_FullMethodName = "/shorturl.Services/GetShortURLByIDHandler"
	Services_Batch_FullMethodName                  = "/shorturl.Services/Batch"
	Services_CreateShortURLJSON_FullMethodName     = "/shorturl.Services/CreateShortURLJSON"
	Services_Ping_FullMethodName                   = "/shorturl.Services/Ping"
)

// ServicesClient is the client API for Services service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServicesClient interface {
	Stats(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StatsResponse, error)
	DeleteHandler(ctx context.Context, in *DeleteListRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetAllURL(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*AllURLsResponse, error)
	CreateShortURLHandler(ctx context.Context, in *CreateShort, opts ...grpc.CallOption) (*CreateShortURLHandlerResponse, error)
	GetShortURLByIDHandler(ctx context.Context, in *IDParam, opts ...grpc.CallOption) (*OneString, error)
	Batch(ctx context.Context, in *PackReq, opts ...grpc.CallOption) (*PackReq, error)
	CreateShortURLJSON(ctx context.Context, in *ReqURL, opts ...grpc.CallOption) (*CreateShortURLJSONResponse, error)
	Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type servicesClient struct {
	cc grpc.ClientConnInterface
}

func NewServicesClient(cc grpc.ClientConnInterface) ServicesClient {
	return &servicesClient{cc}
}

func (c *servicesClient) Stats(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StatsResponse, error) {
	out := new(StatsResponse)
	err := c.cc.Invoke(ctx, Services_Stats_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicesClient) DeleteHandler(ctx context.Context, in *DeleteListRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Services_DeleteHandler_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicesClient) GetAllURL(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*AllURLsResponse, error) {
	out := new(AllURLsResponse)
	err := c.cc.Invoke(ctx, Services_GetAllURL_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicesClient) CreateShortURLHandler(ctx context.Context, in *CreateShort, opts ...grpc.CallOption) (*CreateShortURLHandlerResponse, error) {
	out := new(CreateShortURLHandlerResponse)
	err := c.cc.Invoke(ctx, Services_CreateShortURLHandler_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicesClient) GetShortURLByIDHandler(ctx context.Context, in *IDParam, opts ...grpc.CallOption) (*OneString, error) {
	out := new(OneString)
	err := c.cc.Invoke(ctx, Services_GetShortURLByIDHandler_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicesClient) Batch(ctx context.Context, in *PackReq, opts ...grpc.CallOption) (*PackReq, error) {
	out := new(PackReq)
	err := c.cc.Invoke(ctx, Services_Batch_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicesClient) CreateShortURLJSON(ctx context.Context, in *ReqURL, opts ...grpc.CallOption) (*CreateShortURLJSONResponse, error) {
	out := new(CreateShortURLJSONResponse)
	err := c.cc.Invoke(ctx, Services_CreateShortURLJSON_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicesClient) Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Services_Ping_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServicesServer is the server API for Services service.
// All implementations must embed UnimplementedServicesServer
// for forward compatibility
type ServicesServer interface {
	Stats(context.Context, *emptypb.Empty) (*StatsResponse, error)
	DeleteHandler(context.Context, *DeleteListRequest) (*emptypb.Empty, error)
	GetAllURL(context.Context, *emptypb.Empty) (*AllURLsResponse, error)
	CreateShortURLHandler(context.Context, *CreateShort) (*CreateShortURLHandlerResponse, error)
	GetShortURLByIDHandler(context.Context, *IDParam) (*OneString, error)
	Batch(context.Context, *PackReq) (*PackReq, error)
	CreateShortURLJSON(context.Context, *ReqURL) (*CreateShortURLJSONResponse, error)
	Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	mustEmbedUnimplementedServicesServer()
}

// UnimplementedServicesServer must be embedded to have forward compatible implementations.
type UnimplementedServicesServer struct {
}

func (UnimplementedServicesServer) Stats(context.Context, *emptypb.Empty) (*StatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stats not implemented")
}
func (UnimplementedServicesServer) DeleteHandler(context.Context, *DeleteListRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteHandler not implemented")
}
func (UnimplementedServicesServer) GetAllURL(context.Context, *emptypb.Empty) (*AllURLsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllURL not implemented")
}
func (UnimplementedServicesServer) CreateShortURLHandler(context.Context, *CreateShort) (*CreateShortURLHandlerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateShortURLHandler not implemented")
}
func (UnimplementedServicesServer) GetShortURLByIDHandler(context.Context, *IDParam) (*OneString, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetShortURLByIDHandler not implemented")
}
func (UnimplementedServicesServer) Batch(context.Context, *PackReq) (*PackReq, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Batch not implemented")
}
func (UnimplementedServicesServer) CreateShortURLJSON(context.Context, *ReqURL) (*CreateShortURLJSONResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateShortURLJSON not implemented")
}
func (UnimplementedServicesServer) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedServicesServer) mustEmbedUnimplementedServicesServer() {}

// UnsafeServicesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServicesServer will
// result in compilation errors.
type UnsafeServicesServer interface {
	mustEmbedUnimplementedServicesServer()
}

func RegisterServicesServer(s grpc.ServiceRegistrar, srv ServicesServer) {
	s.RegisterService(&Services_ServiceDesc, srv)
}

func _Services_Stats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesServer).Stats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Services_Stats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesServer).Stats(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Services_DeleteHandler_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesServer).DeleteHandler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Services_DeleteHandler_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesServer).DeleteHandler(ctx, req.(*DeleteListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Services_GetAllURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesServer).GetAllURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Services_GetAllURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesServer).GetAllURL(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Services_CreateShortURLHandler_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateShort)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesServer).CreateShortURLHandler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Services_CreateShortURLHandler_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesServer).CreateShortURLHandler(ctx, req.(*CreateShort))
	}
	return interceptor(ctx, in, info, handler)
}

func _Services_GetShortURLByIDHandler_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IDParam)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesServer).GetShortURLByIDHandler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Services_GetShortURLByIDHandler_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesServer).GetShortURLByIDHandler(ctx, req.(*IDParam))
	}
	return interceptor(ctx, in, info, handler)
}

func _Services_Batch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PackReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesServer).Batch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Services_Batch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesServer).Batch(ctx, req.(*PackReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Services_CreateShortURLJSON_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqURL)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesServer).CreateShortURLJSON(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Services_CreateShortURLJSON_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesServer).CreateShortURLJSON(ctx, req.(*ReqURL))
	}
	return interceptor(ctx, in, info, handler)
}

func _Services_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Services_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesServer).Ping(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Services_ServiceDesc is the grpc.ServiceDesc for Services service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Services_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shorturl.Services",
	HandlerType: (*ServicesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Stats",
			Handler:    _Services_Stats_Handler,
		},
		{
			MethodName: "DeleteHandler",
			Handler:    _Services_DeleteHandler_Handler,
		},
		{
			MethodName: "GetAllURL",
			Handler:    _Services_GetAllURL_Handler,
		},
		{
			MethodName: "CreateShortURLHandler",
			Handler:    _Services_CreateShortURLHandler_Handler,
		},
		{
			MethodName: "GetShortURLByIDHandler",
			Handler:    _Services_GetShortURLByIDHandler_Handler,
		},
		{
			MethodName: "Batch",
			Handler:    _Services_Batch_Handler,
		},
		{
			MethodName: "CreateShortURLJSON",
			Handler:    _Services_CreateShortURLJSON_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _Services_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/proto/services.proto",
}
