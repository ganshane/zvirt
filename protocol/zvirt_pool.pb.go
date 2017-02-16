// Code generated by protoc-gen-go.
// source: zvirt_pool.proto
// DO NOT EDIT!

package protocol

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

// domain state
type PoolState int32

const (
	PoolState_STORAGE_POOL_INACTIVE     PoolState = 0
	PoolState_STORAGE_POOL_BUILDING     PoolState = 1
	PoolState_STORAGE_POOL_RUNNING      PoolState = 2
	PoolState_STORAGE_POOL_DEGRADED     PoolState = 3
	PoolState_STORAGE_POOL_INACCESSIBLE PoolState = 4
)

var PoolState_name = map[int32]string{
	0: "STORAGE_POOL_INACTIVE",
	1: "STORAGE_POOL_BUILDING",
	2: "STORAGE_POOL_RUNNING",
	3: "STORAGE_POOL_DEGRADED",
	4: "STORAGE_POOL_INACCESSIBLE",
}
var PoolState_value = map[string]int32{
	"STORAGE_POOL_INACTIVE":     0,
	"STORAGE_POOL_BUILDING":     1,
	"STORAGE_POOL_RUNNING":      2,
	"STORAGE_POOL_DEGRADED":     3,
	"STORAGE_POOL_INACCESSIBLE": 4,
}

func (x PoolState) String() string {
	return proto.EnumName(PoolState_name, int32(x))
}
func (PoolState) EnumDescriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

type PoolUUID struct {
	Uuid string `protobuf:"bytes,1,opt,name=uuid" json:"uuid,omitempty"`
}

func (m *PoolUUID) Reset()                    { *m = PoolUUID{} }
func (m *PoolUUID) String() string            { return proto.CompactTextString(m) }
func (*PoolUUID) ProtoMessage()               {}
func (*PoolUUID) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *PoolUUID) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

type PoolDefineRequest struct {
	Xml string `protobuf:"bytes,1,opt,name=xml" json:"xml,omitempty"`
}

func (m *PoolDefineRequest) Reset()                    { *m = PoolDefineRequest{} }
func (m *PoolDefineRequest) String() string            { return proto.CompactTextString(m) }
func (*PoolDefineRequest) ProtoMessage()               {}
func (*PoolDefineRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *PoolDefineRequest) GetXml() string {
	if m != nil {
		return m.Xml
	}
	return ""
}

type PoolStateResponse struct {
	State PoolState `protobuf:"varint,1,opt,name=state,enum=protocol.PoolState" json:"state,omitempty"`
}

func (m *PoolStateResponse) Reset()                    { *m = PoolStateResponse{} }
func (m *PoolStateResponse) String() string            { return proto.CompactTextString(m) }
func (*PoolStateResponse) ProtoMessage()               {}
func (*PoolStateResponse) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *PoolStateResponse) GetState() PoolState {
	if m != nil {
		return m.State
	}
	return PoolState_STORAGE_POOL_INACTIVE
}

func init() {
	proto.RegisterType((*PoolUUID)(nil), "protocol.PoolUUID")
	proto.RegisterType((*PoolDefineRequest)(nil), "protocol.PoolDefineRequest")
	proto.RegisterType((*PoolStateResponse)(nil), "protocol.PoolStateResponse")
	proto.RegisterEnum("protocol.PoolState", PoolState_name, PoolState_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ZvirtPoolService service

type ZvirtPoolServiceClient interface {
	Info(ctx context.Context, in *PoolUUID, opts ...grpc.CallOption) (*PoolStateResponse, error)
	Define(ctx context.Context, in *PoolDefineRequest, opts ...grpc.CallOption) (*PoolUUID, error)
	Start(ctx context.Context, in *PoolUUID, opts ...grpc.CallOption) (*PoolStateResponse, error)
	Destroy(ctx context.Context, in *PoolUUID, opts ...grpc.CallOption) (*PoolStateResponse, error)
}

type zvirtPoolServiceClient struct {
	cc *grpc.ClientConn
}

func NewZvirtPoolServiceClient(cc *grpc.ClientConn) ZvirtPoolServiceClient {
	return &zvirtPoolServiceClient{cc}
}

func (c *zvirtPoolServiceClient) Info(ctx context.Context, in *PoolUUID, opts ...grpc.CallOption) (*PoolStateResponse, error) {
	out := new(PoolStateResponse)
	err := grpc.Invoke(ctx, "/protocol.ZvirtPoolService/info", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *zvirtPoolServiceClient) Define(ctx context.Context, in *PoolDefineRequest, opts ...grpc.CallOption) (*PoolUUID, error) {
	out := new(PoolUUID)
	err := grpc.Invoke(ctx, "/protocol.ZvirtPoolService/define", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *zvirtPoolServiceClient) Start(ctx context.Context, in *PoolUUID, opts ...grpc.CallOption) (*PoolStateResponse, error) {
	out := new(PoolStateResponse)
	err := grpc.Invoke(ctx, "/protocol.ZvirtPoolService/start", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *zvirtPoolServiceClient) Destroy(ctx context.Context, in *PoolUUID, opts ...grpc.CallOption) (*PoolStateResponse, error) {
	out := new(PoolStateResponse)
	err := grpc.Invoke(ctx, "/protocol.ZvirtPoolService/destroy", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ZvirtPoolService service

type ZvirtPoolServiceServer interface {
	Info(context.Context, *PoolUUID) (*PoolStateResponse, error)
	Define(context.Context, *PoolDefineRequest) (*PoolUUID, error)
	Start(context.Context, *PoolUUID) (*PoolStateResponse, error)
	Destroy(context.Context, *PoolUUID) (*PoolStateResponse, error)
}

func RegisterZvirtPoolServiceServer(s *grpc.Server, srv ZvirtPoolServiceServer) {
	s.RegisterService(&_ZvirtPoolService_serviceDesc, srv)
}

func _ZvirtPoolService_Info_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PoolUUID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZvirtPoolServiceServer).Info(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.ZvirtPoolService/Info",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZvirtPoolServiceServer).Info(ctx, req.(*PoolUUID))
	}
	return interceptor(ctx, in, info, handler)
}

func _ZvirtPoolService_Define_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PoolDefineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZvirtPoolServiceServer).Define(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.ZvirtPoolService/Define",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZvirtPoolServiceServer).Define(ctx, req.(*PoolDefineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ZvirtPoolService_Start_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PoolUUID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZvirtPoolServiceServer).Start(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.ZvirtPoolService/Start",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZvirtPoolServiceServer).Start(ctx, req.(*PoolUUID))
	}
	return interceptor(ctx, in, info, handler)
}

func _ZvirtPoolService_Destroy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PoolUUID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZvirtPoolServiceServer).Destroy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.ZvirtPoolService/Destroy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZvirtPoolServiceServer).Destroy(ctx, req.(*PoolUUID))
	}
	return interceptor(ctx, in, info, handler)
}

var _ZvirtPoolService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protocol.ZvirtPoolService",
	HandlerType: (*ZvirtPoolServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "info",
			Handler:    _ZvirtPoolService_Info_Handler,
		},
		{
			MethodName: "define",
			Handler:    _ZvirtPoolService_Define_Handler,
		},
		{
			MethodName: "start",
			Handler:    _ZvirtPoolService_Start_Handler,
		},
		{
			MethodName: "destroy",
			Handler:    _ZvirtPoolService_Destroy_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "zvirt_pool.proto",
}

func init() { proto.RegisterFile("zvirt_pool.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 331 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x9c, 0x92, 0x4f, 0x4f, 0xf2, 0x40,
	0x10, 0x87, 0x29, 0xf0, 0xf2, 0xc2, 0x1c, 0xc8, 0x3a, 0x6a, 0x02, 0x18, 0x8d, 0x69, 0x62, 0xa2,
	0x1e, 0x38, 0xe0, 0xc9, 0x3f, 0x31, 0xa1, 0x74, 0x43, 0x36, 0x21, 0xa5, 0xd9, 0x52, 0x0f, 0x5e,
	0x88, 0xc2, 0x92, 0x34, 0xa9, 0x2c, 0xb6, 0x0b, 0x51, 0xcf, 0x7e, 0x05, 0xbf, 0xaf, 0xd9, 0x05,
	0x31, 0x58, 0x4e, 0x9c, 0x3a, 0x99, 0xdf, 0x33, 0x9d, 0xce, 0x93, 0x02, 0xf9, 0x58, 0x44, 0x89,
	0x1a, 0xce, 0xa4, 0x8c, 0x9b, 0xb3, 0x44, 0x2a, 0x89, 0x65, 0xf3, 0x18, 0xc9, 0xd8, 0x3e, 0x81,
	0xb2, 0x2f, 0x65, 0x1c, 0x86, 0xcc, 0x45, 0x84, 0xe2, 0x7c, 0x1e, 0x8d, 0x6b, 0xd6, 0xa9, 0x75,
	0x5e, 0xe1, 0xa6, 0xb6, 0xcf, 0x60, 0x4f, 0xe7, 0xae, 0x98, 0x44, 0x53, 0xc1, 0xc5, 0xeb, 0x5c,
	0xa4, 0x0a, 0x09, 0x14, 0xde, 0x5e, 0xe2, 0x15, 0xa7, 0x4b, 0xfb, 0x7e, 0x89, 0x05, 0xea, 0x49,
	0x09, 0x2e, 0xd2, 0x99, 0x9c, 0xa6, 0x02, 0x2f, 0xe0, 0x5f, 0xaa, 0x1b, 0x06, 0xac, 0xb6, 0xf6,
	0x9b, 0x3f, 0x5b, 0x9b, 0xbf, 0xec, 0x92, 0xb8, 0xfc, 0xb2, 0xa0, 0xb2, 0x6e, 0x62, 0x1d, 0x0e,
	0x83, 0x41, 0x9f, 0xb7, 0xbb, 0x74, 0xe8, 0xf7, 0xfb, 0xbd, 0x21, 0xf3, 0xda, 0x9d, 0x01, 0x7b,
	0xa0, 0x24, 0x97, 0x89, 0x9c, 0x90, 0xf5, 0x5c, 0xe6, 0x75, 0x89, 0x85, 0x35, 0x38, 0xd8, 0x88,
	0x78, 0xe8, 0x79, 0x3a, 0xc9, 0x67, 0x86, 0x5c, 0xda, 0xe5, 0x6d, 0x97, 0xba, 0xa4, 0x80, 0xc7,
	0x50, 0xcf, 0xac, 0xea, 0xd0, 0x20, 0x60, 0x4e, 0x8f, 0x92, 0x62, 0xeb, 0x33, 0x0f, 0xe4, 0x51,
	0xdb, 0x33, 0x1f, 0x27, 0x92, 0x45, 0x34, 0x12, 0x78, 0x0d, 0xc5, 0x68, 0x3a, 0x91, 0x88, 0x9b,
	0x07, 0x69, 0x87, 0x8d, 0xa3, 0x6d, 0x47, 0xae, 0x84, 0xd8, 0x39, 0xbc, 0x85, 0xd2, 0xd8, 0xa8,
	0xc4, 0x3f, 0xe0, 0x86, 0xe0, 0xc6, 0x96, 0x37, 0xdb, 0x39, 0xbc, 0x31, 0x3e, 0x13, 0xb5, 0xcb,
	0xe2, 0x3b, 0xf8, 0x3f, 0x16, 0xa9, 0x4a, 0xe4, 0xfb, 0x0e, 0xd3, 0x8e, 0x0d, 0x55, 0xf3, 0x0f,
	0xad, 0x29, 0xa7, 0xba, 0xb6, 0xe2, 0xeb, 0x96, 0x6f, 0x3d, 0x97, 0x4c, 0x76, 0xf5, 0x1d, 0x00,
	0x00, 0xff, 0xff, 0x0a, 0x7e, 0x7e, 0xdd, 0x6e, 0x02, 0x00, 0x00,
}
