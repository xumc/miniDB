// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api.proto

package server

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	math "math"
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

type SqlRequest struct {
	Sql                  []string `protobuf:"bytes,1,rep,name=sql,proto3" json:"sql,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SqlRequest) Reset()         { *m = SqlRequest{} }
func (m *SqlRequest) String() string { return proto.CompactTextString(m) }
func (*SqlRequest) ProtoMessage()    {}
func (*SqlRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{0}
}

func (m *SqlRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SqlRequest.Unmarshal(m, b)
}
func (m *SqlRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SqlRequest.Marshal(b, m, deterministic)
}
func (m *SqlRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SqlRequest.Merge(m, src)
}
func (m *SqlRequest) XXX_Size() int {
	return xxx_messageInfo_SqlRequest.Size(m)
}
func (m *SqlRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SqlRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SqlRequest proto.InternalMessageInfo

func (m *SqlRequest) GetSql() []string {
	if m != nil {
		return m.Sql
	}
	return nil
}

type SqlReply struct {
	Records              []*SqlReplyRecord `protobuf:"bytes,1,rep,name=records,proto3" json:"records,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *SqlReply) Reset()         { *m = SqlReply{} }
func (m *SqlReply) String() string { return proto.CompactTextString(m) }
func (*SqlReply) ProtoMessage()    {}
func (*SqlReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{1}
}

func (m *SqlReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SqlReply.Unmarshal(m, b)
}
func (m *SqlReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SqlReply.Marshal(b, m, deterministic)
}
func (m *SqlReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SqlReply.Merge(m, src)
}
func (m *SqlReply) XXX_Size() int {
	return xxx_messageInfo_SqlReply.Size(m)
}
func (m *SqlReply) XXX_DiscardUnknown() {
	xxx_messageInfo_SqlReply.DiscardUnknown(m)
}

var xxx_messageInfo_SqlReply proto.InternalMessageInfo

func (m *SqlReply) GetRecords() []*SqlReplyRecord {
	if m != nil {
		return m.Records
	}
	return nil
}

type SqlReplyRecord struct {
	Columns              []string `protobuf:"bytes,1,rep,name=columns,proto3" json:"columns,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SqlReplyRecord) Reset()         { *m = SqlReplyRecord{} }
func (m *SqlReplyRecord) String() string { return proto.CompactTextString(m) }
func (*SqlReplyRecord) ProtoMessage()    {}
func (*SqlReplyRecord) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{1, 0}
}

func (m *SqlReplyRecord) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SqlReplyRecord.Unmarshal(m, b)
}
func (m *SqlReplyRecord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SqlReplyRecord.Marshal(b, m, deterministic)
}
func (m *SqlReplyRecord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SqlReplyRecord.Merge(m, src)
}
func (m *SqlReplyRecord) XXX_Size() int {
	return xxx_messageInfo_SqlReplyRecord.Size(m)
}
func (m *SqlReplyRecord) XXX_DiscardUnknown() {
	xxx_messageInfo_SqlReplyRecord.DiscardUnknown(m)
}

var xxx_messageInfo_SqlReplyRecord proto.InternalMessageInfo

func (m *SqlReplyRecord) GetColumns() []string {
	if m != nil {
		return m.Columns
	}
	return nil
}

func init() {
	proto.RegisterType((*SqlRequest)(nil), "server.SqlRequest")
	proto.RegisterType((*SqlReply)(nil), "server.SqlReply")
	proto.RegisterType((*SqlReplyRecord)(nil), "server.SqlReply.record")
}

func init() { proto.RegisterFile("api.proto", fileDescriptor_00212fb1f9d3bf1c) }

var fileDescriptor_00212fb1f9d3bf1c = []byte{
	// 165 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4c, 0x2c, 0xc8, 0xd4,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2b, 0x4e, 0x2d, 0x2a, 0x4b, 0x2d, 0x52, 0x92, 0xe3,
	0xe2, 0x0a, 0x2e, 0xcc, 0x09, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x11, 0x12, 0xe0, 0x62, 0x2e,
	0x2e, 0xcc, 0x91, 0x60, 0x54, 0x60, 0xd6, 0xe0, 0x0c, 0x02, 0x31, 0x95, 0x12, 0xb9, 0x38, 0xc0,
	0xf2, 0x05, 0x39, 0x95, 0x42, 0x86, 0x5c, 0xec, 0x45, 0xa9, 0xc9, 0xf9, 0x45, 0x29, 0xc5, 0x60,
	0x15, 0xdc, 0x46, 0xe2, 0x7a, 0x10, 0x53, 0xf4, 0x60, 0x4a, 0xf4, 0x20, 0xf2, 0x41, 0x30, 0x75,
	0x52, 0x4a, 0x5c, 0x6c, 0x10, 0xa6, 0x90, 0x04, 0x17, 0x7b, 0x72, 0x7e, 0x4e, 0x69, 0x6e, 0x5e,
	0x31, 0xd4, 0x78, 0x18, 0xd7, 0xc8, 0x96, 0x8b, 0xc3, 0xb5, 0x22, 0x35, 0xb9, 0xb4, 0x24, 0xbf,
	0x08, 0x64, 0x05, 0x84, 0x9d, 0x2a, 0x24, 0x84, 0x62, 0x38, 0xd8, 0x7d, 0x52, 0x02, 0xe8, 0x16,
	0x2a, 0x31, 0x24, 0xb1, 0x81, 0x3d, 0x64, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x26, 0xf0, 0x45,
	0xa2, 0xdd, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ExecutorClient is the client API for Executor service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ExecutorClient interface {
	Execute(ctx context.Context, in *SqlRequest, opts ...grpc.CallOption) (*SqlReply, error)
}

type executorClient struct {
	cc *grpc.ClientConn
}

func NewExecutorClient(cc *grpc.ClientConn) ExecutorClient {
	return &executorClient{cc}
}

func (c *executorClient) Execute(ctx context.Context, in *SqlRequest, opts ...grpc.CallOption) (*SqlReply, error) {
	out := new(SqlReply)
	err := c.cc.Invoke(ctx, "/server.Executor/Execute", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExecutorServer is the server API for Executor service.
type ExecutorServer interface {
	Execute(context.Context, *SqlRequest) (*SqlReply, error)
}

func RegisterExecutorServer(s *grpc.Server, srv ExecutorServer) {
	s.RegisterService(&_Executor_serviceDesc, srv)
}

func _Executor_Execute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SqlRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).Execute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.Executor/Execute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).Execute(ctx, req.(*SqlRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Executor_serviceDesc = grpc.ServiceDesc{
	ServiceName: "server.Executor",
	HandlerType: (*ExecutorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Execute",
			Handler:    _Executor_Execute_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}
