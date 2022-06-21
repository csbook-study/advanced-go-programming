// Code generated by protoc-gen-go. DO NOT EDIT.
// source: hello.proto

package _

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import "net/rpc"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type String struct {
	Value                string   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *String) Reset()         { *m = String{} }
func (m *String) String() string { return proto.CompactTextString(m) }
func (*String) ProtoMessage()    {}
func (*String) Descriptor() ([]byte, []int) {
	return fileDescriptor_61ef911816e0a8ce, []int{0}
}

func (m *String) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_String.Unmarshal(m, b)
}
func (m *String) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_String.Marshal(b, m, deterministic)
}
func (m *String) XXX_Merge(src proto.Message) {
	xxx_messageInfo_String.Merge(m, src)
}
func (m *String) XXX_Size() int {
	return xxx_messageInfo_String.Size(m)
}
func (m *String) XXX_DiscardUnknown() {
	xxx_messageInfo_String.DiscardUnknown(m)
}

var xxx_messageInfo_String proto.InternalMessageInfo

func (m *String) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func init() {
	proto.RegisterType((*String)(nil), "hello.String")
}

func init() { proto.RegisterFile("hello.proto", fileDescriptor_61ef911816e0a8ce) }

var fileDescriptor_61ef911816e0a8ce = []byte{
	// 108 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xce, 0x48, 0xcd, 0xc9,
	0xc9, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x73, 0x94, 0xe4, 0xb8, 0xd8, 0x82,
	0x4b, 0x8a, 0x32, 0xf3, 0xd2, 0x85, 0x44, 0xb8, 0x58, 0xcb, 0x12, 0x73, 0x4a, 0x53, 0x25, 0x18,
	0x15, 0x18, 0x35, 0x38, 0x83, 0x20, 0x1c, 0x23, 0x53, 0x2e, 0x1e, 0x0f, 0x90, 0xc2, 0xe0, 0xd4,
	0xa2, 0xb2, 0xcc, 0xe4, 0x54, 0x21, 0x55, 0x2e, 0x56, 0x30, 0x5f, 0x88, 0x57, 0x0f, 0x62, 0x1a,
	0x44, 0xb7, 0x14, 0x2a, 0xd7, 0x89, 0x39, 0x8a, 0x51, 0x2f, 0x89, 0x0d, 0x6c, 0x93, 0x31, 0x20,
	0x00, 0x00, 0xff, 0xff, 0x64, 0x83, 0x6f, 0x86, 0x78, 0x00, 0x00, 0x00,
}

type HelloServiceInterface interface {
	Hello(*String, *String) error
}

func RegisterHelloService(
	srv *rpc.Server, x HelloServiceInterface,
) error {
	if err := srv.RegisterName("HelloService", x); err != nil {
		return err
	}
	return nil
}

type HelloServiceClient struct {
	*rpc.Client
}

var _ HelloServiceInterface = (*HelloServiceClient)(nil)

func DialHelloService(network, address string) (
	*HelloServiceClient, error,
) {
	c, err := rpc.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &HelloServiceClient{Client: c}, nil
}

func (p *HelloServiceClient) Hello(
	in *String, out *String,
) error {
	return p.Client.Call("HelloService.Hello", in, out)
}
