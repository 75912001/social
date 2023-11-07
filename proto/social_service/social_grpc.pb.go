// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.0
// source: impl/protobuf/social/social.proto

package social_service

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

// SocialServiceClient is the client API for SocialService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SocialServiceClient interface {
	BidirectionalStreamingMethod(ctx context.Context, opts ...grpc.CallOption) (SocialService_BidirectionalStreamingMethodClient, error)
}

type socialServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSocialServiceClient(cc grpc.ClientConnInterface) SocialServiceClient {
	return &socialServiceClient{cc}
}

func (c *socialServiceClient) BidirectionalStreamingMethod(ctx context.Context, opts ...grpc.CallOption) (SocialService_BidirectionalStreamingMethodClient, error) {
	stream, err := c.cc.NewStream(ctx, &SocialService_ServiceDesc.Streams[0], "/social_service.SocialService/BidirectionalStreamingMethod", opts...)
	if err != nil {
		return nil, err
	}
	x := &socialServiceBidirectionalStreamingMethodClient{stream}
	return x, nil
}

type SocialService_BidirectionalStreamingMethodClient interface {
	Send(*CommonRequest) error
	Recv() (*CommonResponse, error)
	grpc.ClientStream
}

type socialServiceBidirectionalStreamingMethodClient struct {
	grpc.ClientStream
}

func (x *socialServiceBidirectionalStreamingMethodClient) Send(m *CommonRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *socialServiceBidirectionalStreamingMethodClient) Recv() (*CommonResponse, error) {
	m := new(CommonResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SocialServiceServer is the server API for SocialService service.
// All implementations must embed UnimplementedSocialServiceServer
// for forward compatibility
type SocialServiceServer interface {
	BidirectionalStreamingMethod(SocialService_BidirectionalStreamingMethodServer) error
	mustEmbedUnimplementedSocialServiceServer()
}

// UnimplementedSocialServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSocialServiceServer struct {
}

func (UnimplementedSocialServiceServer) BidirectionalStreamingMethod(SocialService_BidirectionalStreamingMethodServer) error {
	return status.Errorf(codes.Unimplemented, "method BidirectionalStreamingMethod not implemented")
}
func (UnimplementedSocialServiceServer) mustEmbedUnimplementedSocialServiceServer() {}

// UnsafeSocialServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SocialServiceServer will
// result in compilation errors.
type UnsafeSocialServiceServer interface {
	mustEmbedUnimplementedSocialServiceServer()
}

func RegisterSocialServiceServer(s grpc.ServiceRegistrar, srv SocialServiceServer) {
	s.RegisterService(&SocialService_ServiceDesc, srv)
}

func _SocialService_BidirectionalStreamingMethod_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SocialServiceServer).BidirectionalStreamingMethod(&socialServiceBidirectionalStreamingMethodServer{stream})
}

type SocialService_BidirectionalStreamingMethodServer interface {
	Send(*CommonResponse) error
	Recv() (*CommonRequest, error)
	grpc.ServerStream
}

type socialServiceBidirectionalStreamingMethodServer struct {
	grpc.ServerStream
}

func (x *socialServiceBidirectionalStreamingMethodServer) Send(m *CommonResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *socialServiceBidirectionalStreamingMethodServer) Recv() (*CommonRequest, error) {
	m := new(CommonRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SocialService_ServiceDesc is the grpc.ServiceDesc for SocialService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SocialService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "social_service.SocialService",
	HandlerType: (*SocialServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "BidirectionalStreamingMethod",
			Handler:       _SocialService_BidirectionalStreamingMethod_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "impl/protobuf/social/social.proto",
}
