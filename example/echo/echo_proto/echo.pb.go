// Code generated by protoc-gen-go. DO NOT EDIT.
// source: echo.proto

/*
Package echo_proto is a generated protocol buffer package.

It is generated from these files:
	echo.proto

It has these top-level messages:
	EchoRequest
	EchoResponse
	UpperRequest
	UpperResponse
*/
package echo_proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type EchoRequest struct {
	Msg string `protobuf:"bytes,1,opt,name=msg" json:"msg,omitempty"`
}

func (m *EchoRequest) Reset()                    { *m = EchoRequest{} }
func (m *EchoRequest) String() string            { return proto.CompactTextString(m) }
func (*EchoRequest) ProtoMessage()               {}
func (*EchoRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *EchoRequest) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type EchoResponse struct {
	Msg string `protobuf:"bytes,1,opt,name=msg" json:"msg,omitempty"`
}

func (m *EchoResponse) Reset()                    { *m = EchoResponse{} }
func (m *EchoResponse) String() string            { return proto.CompactTextString(m) }
func (*EchoResponse) ProtoMessage()               {}
func (*EchoResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *EchoResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type UpperRequest struct {
	Msg string `protobuf:"bytes,1,opt,name=msg" json:"msg,omitempty"`
}

func (m *UpperRequest) Reset()                    { *m = UpperRequest{} }
func (m *UpperRequest) String() string            { return proto.CompactTextString(m) }
func (*UpperRequest) ProtoMessage()               {}
func (*UpperRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *UpperRequest) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type UpperResponse struct {
	Msg string `protobuf:"bytes,1,opt,name=msg" json:"msg,omitempty"`
}

func (m *UpperResponse) Reset()                    { *m = UpperResponse{} }
func (m *UpperResponse) String() string            { return proto.CompactTextString(m) }
func (*UpperResponse) ProtoMessage()               {}
func (*UpperResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *UpperResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func init() {
	proto.RegisterType((*EchoRequest)(nil), "echo_proto.EchoRequest")
	proto.RegisterType((*EchoResponse)(nil), "echo_proto.EchoResponse")
	proto.RegisterType((*UpperRequest)(nil), "echo_proto.UpperRequest")
	proto.RegisterType((*UpperResponse)(nil), "echo_proto.UpperResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for EchoService service

type EchoServiceClient interface {
	Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoResponse, error)
	// ORION:URL: GET/POST/OPTIONS /api/1.0/upper/{msg}
	Upper(ctx context.Context, in *UpperRequest, opts ...grpc.CallOption) (*UpperResponse, error)
	UpperProxy(ctx context.Context, in *UpperRequest, opts ...grpc.CallOption) (*UpperResponse, error)
}

type echoServiceClient struct {
	cc *grpc.ClientConn
}

func NewEchoServiceClient(cc *grpc.ClientConn) EchoServiceClient {
	return &echoServiceClient{cc}
}

func (c *echoServiceClient) Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoResponse, error) {
	out := new(EchoResponse)
	err := grpc.Invoke(ctx, "/echo_proto.EchoService/Echo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *echoServiceClient) Upper(ctx context.Context, in *UpperRequest, opts ...grpc.CallOption) (*UpperResponse, error) {
	out := new(UpperResponse)
	err := grpc.Invoke(ctx, "/echo_proto.EchoService/Upper", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *echoServiceClient) UpperProxy(ctx context.Context, in *UpperRequest, opts ...grpc.CallOption) (*UpperResponse, error) {
	out := new(UpperResponse)
	err := grpc.Invoke(ctx, "/echo_proto.EchoService/UpperProxy", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for EchoService service

type EchoServiceServer interface {
	Echo(context.Context, *EchoRequest) (*EchoResponse, error)
	// ORION:URL: GET/POST/OPTIONS /api/1.0/upper/{msg}
	Upper(context.Context, *UpperRequest) (*UpperResponse, error)
	UpperProxy(context.Context, *UpperRequest) (*UpperResponse, error)
}

func RegisterEchoServiceServer(s *grpc.Server, srv EchoServiceServer) {
	s.RegisterService(&_EchoService_serviceDesc, srv)

}

func _EchoService_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EchoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EchoServiceServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/echo_proto.EchoService/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EchoServiceServer).Echo(ctx, req.(*EchoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EchoService_Upper_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpperRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EchoServiceServer).Upper(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/echo_proto.EchoService/Upper",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EchoServiceServer).Upper(ctx, req.(*UpperRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EchoService_UpperProxy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpperRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EchoServiceServer).UpperProxy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/echo_proto.EchoService/UpperProxy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EchoServiceServer).UpperProxy(ctx, req.(*UpperRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _EchoService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "echo_proto.EchoService",
	HandlerType: (*EchoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _EchoService_Echo_Handler,
		},
		{
			MethodName: "Upper",
			Handler:    _EchoService_Upper_Handler,
		},
		{
			MethodName: "UpperProxy",
			Handler:    _EchoService_UpperProxy_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "echo.proto",
}

func init() { proto.RegisterFile("echo.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 171 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0x4d, 0xce, 0xc8,
	0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x02, 0xb3, 0xe3, 0xc1, 0x6c, 0x25, 0x79, 0x2e, 0x6e,
	0xd7, 0xe4, 0x8c, 0xfc, 0xa0, 0xd4, 0xc2, 0xd2, 0xd4, 0xe2, 0x12, 0x21, 0x01, 0x2e, 0xe6, 0xdc,
	0xe2, 0x74, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x10, 0x53, 0x49, 0x81, 0x8b, 0x07, 0xa2,
	0xa0, 0xb8, 0x20, 0x3f, 0xaf, 0x38, 0x15, 0xbb, 0x8a, 0xd0, 0x82, 0x82, 0xd4, 0x22, 0xdc, 0x66,
	0x28, 0x72, 0xf1, 0x42, 0x55, 0xe0, 0x32, 0xc4, 0xe8, 0x3c, 0x23, 0xc4, 0x21, 0xc1, 0xa9, 0x45,
	0x65, 0x99, 0xc9, 0xa9, 0x42, 0xd6, 0x5c, 0x2c, 0x20, 0xae, 0x90, 0xb8, 0x1e, 0xc2, 0xb1, 0x7a,
	0x48, 0x2e, 0x95, 0x92, 0xc0, 0x94, 0x80, 0x18, 0xae, 0xc4, 0x20, 0x64, 0xc7, 0xc5, 0x0a, 0xb6,
	0x4f, 0x08, 0x45, 0x11, 0xb2, 0x23, 0xa5, 0x24, 0xb1, 0xc8, 0xc0, 0xf5, 0x3b, 0x73, 0x71, 0x81,
	0x85, 0x02, 0x8a, 0xf2, 0x2b, 0x2a, 0xc9, 0x34, 0x24, 0x89, 0x0d, 0x2c, 0x6c, 0x0c, 0x08, 0x00,
	0x00, 0xff, 0xff, 0x74, 0xf4, 0xf2, 0x29, 0x7a, 0x01, 0x00, 0x00,
}
