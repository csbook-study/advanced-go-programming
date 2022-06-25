// Code generated by protoc-gen-go. DO NOT EDIT.
// source: hello.proto

package hello

import (
	fmt "fmt"
	_ "github.com/chai2010/pbgo"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	encoding_json "encoding/json"
	io "io"
	io_ioutil "io/ioutil"
	net_http "net/http"
	net_rpc "net/rpc"
	regexp "regexp"
	strings "strings"

	github_com_chai2010_pbgo "github.com/chai2010/pbgo"
	github_com_julienschmidt_httprouter "github.com/julienschmidt/httprouter"
)

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
	// 150 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xce, 0x48, 0xcd, 0xc9,
	0xc9, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x73, 0xa4, 0x94, 0xd3, 0x33, 0x4b,
	0x32, 0x4a, 0x93, 0xf4, 0x92, 0xf3, 0x73, 0xf5, 0x93, 0x33, 0x12, 0x33, 0x8d, 0x0c, 0x0c, 0x0d,
	0xf4, 0x0b, 0x92, 0xd2, 0xf3, 0xc1, 0x04, 0x44, 0xad, 0x92, 0x1c, 0x17, 0x5b, 0x70, 0x49, 0x51,
	0x66, 0x5e, 0xba, 0x90, 0x08, 0x17, 0x6b, 0x59, 0x62, 0x4e, 0x69, 0xaa, 0x04, 0xa3, 0x02, 0xa3,
	0x06, 0x67, 0x10, 0x84, 0x63, 0xe4, 0xcd, 0xc5, 0xe3, 0x01, 0x32, 0x2d, 0x38, 0xb5, 0xa8, 0x2c,
	0x33, 0x39, 0x55, 0xc8, 0x9a, 0x8b, 0x15, 0xcc, 0x17, 0xe2, 0xd5, 0x83, 0x58, 0x09, 0xd1, 0x2d,
	0x85, 0xca, 0x55, 0x12, 0xb9, 0xf5, 0xee, 0xaf, 0x0f, 0x3f, 0x17, 0xaf, 0x3e, 0x58, 0x54, 0xdf,
	0x0a, 0x6c, 0x58, 0x12, 0x1b, 0xd8, 0x4e, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x95, 0xb7,
	0xed, 0x50, 0xae, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ = context.Background
var _ = encoding_json.Marshal
var _ = net_rpc.Server{}
var _ = net_http.ListenAndServe
var _ = io.EOF
var _ = io_ioutil.ReadAll
var _ = regexp.Match
var _ = strings.Split
var _ = github_com_chai2010_pbgo.PopulateFieldFromPath
var _ = github_com_julienschmidt_httprouter.New

type HelloServiceInterface interface {
	Hello(in *String, out *String) error
}

type HelloServiceGrpcInterface interface {
	Hello(ctx context.Context, in *String) (out *String, err error)
}

func RegisterHelloService(srv *net_rpc.Server, x HelloServiceInterface) error {
	if _, ok := x.(*HelloServiceValidator); !ok {
		x = &HelloServiceValidator{HelloServiceInterface: x}
	}

	if err := srv.RegisterName("HelloService", x); err != nil {
		return err
	}
	return nil
}

type HelloServiceValidator struct {
	HelloServiceInterface
}

func (p *HelloServiceValidator) Hello(in *String, out *String) error {
	if x, ok := proto.Message(in).(interface{ Validate() error }); ok {
		if err := x.Validate(); err != nil {
			return err
		}
	}

	if err := p.HelloServiceInterface.Hello(in, out); err != nil {
		return err
	}

	if x, ok := proto.Message(out).(interface{ Validate() error }); ok {
		if err := x.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type HelloServiceClient struct {
	*net_rpc.Client
}

func DialHelloService(network, address string) (*HelloServiceClient, error) {
	c, err := net_rpc.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &HelloServiceClient{Client: c}, nil
}

func (p *HelloServiceClient) Hello(in *String) (*String, error) {
	if x, ok := proto.Message(in).(interface{ Validate() error }); ok {
		if err := x.Validate(); err != nil {
			return nil, err
		}
	}

	var out = new(String)
	if err := p.Client.Call("HelloService.Hello", in, out); err != nil {
		return nil, err
	}

	if x, ok := proto.Message(out).(interface{ Validate() error }); ok {
		if err := x.Validate(); err != nil {
			return nil, err
		}
	}

	return out, nil
}
func (p *HelloServiceClient) AsyncHello(in *String, out *String, done chan *net_rpc.Call) *net_rpc.Call {
	if x, ok := proto.Message(in).(interface{ Validate() error }); ok {
		if err := x.Validate(); err != nil {
			call := &net_rpc.Call{
				ServiceMethod: "HelloService.Hello",
				Args:          in,
				Reply:         out,
				Error:         err,
				Done:          make(chan *net_rpc.Call, 10),
			}
			call.Done <- call
			return call
		}
	}

	return p.Go(
		"HelloService.Hello",
		in, out,
		done,
	)
}

type HelloServiceHttpClient struct {
	c       *net_http.Client
	baseurl string
}

func NewHelloServiceHttpClient(baseurl string, c ...*net_http.Client) *HelloServiceHttpClient {
	p := &HelloServiceHttpClient{
		c:       net_http.DefaultClient,
		baseurl: baseurl,
	}
	if len(c) != 0 && c[0] != nil {
		p.c = c[0]
	}
	return p
}

func (p *HelloServiceHttpClient) httpDoRequest(method, urlpath string, in interface{}) (mimeType string, content []byte, err error) {
	req, err := github_com_chai2010_pbgo.NewHttpRequest(method, urlpath, in)
	if err != nil {
		return "", nil, err
	}
	resp, err := p.c.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	mimeType = resp.Header.Get("Content-Type")
	content, err = io_ioutil.ReadAll(resp.Body)
	return
}

func (p *HelloServiceHttpClient) Hello(in *String, method ...string) (out *String, err error) {
	if len(method) == 0 {
		method = []string{"GET"}
	}
	if len(method) != 1 {
		return nil, fmt.Errorf("invalid method: %v", method)
	}

	var re = regexp.MustCompile("(\\*|\\:)(\\w|\\.)+")
	_ = re

	out = new(String)
	if method[0] == "GET" {
		urlpath := p.baseurl + fmt.Sprintf("/hello/%v", in.Value)
		err = github_com_chai2010_pbgo.HttpDo(method[0], urlpath, in, out)
		return out, err
	}

	return nil, fmt.Errorf("invalid method: %v", method)
}

func HelloServiceHandler(svc HelloServiceInterface) net_http.Handler {
	var router = github_com_julienschmidt_httprouter.New()

	var re = regexp.MustCompile("(\\*|\\:)(\\w|\\.)+")
	_ = re

	router.Handle("GET", "/hello/:value",
		func(w net_http.ResponseWriter, r *net_http.Request, ps github_com_julienschmidt_httprouter.Params) {
			var (
				protoReq   String
				protoReply String
			)

			for _, fieldPath := range re.FindAllString("/hello/:value", -1) {
				fieldPath := strings.TrimLeft(fieldPath, ":*")
				err := github_com_chai2010_pbgo.PopulateFieldFromPath(&protoReq, fieldPath, ps.ByName(fieldPath))
				if err != nil {
					net_http.Error(w, err.Error(), net_http.StatusBadRequest)
					return
				}
			}

			if err := github_com_chai2010_pbgo.PopulateQueryParameters(&protoReq, r.URL.Query()); err != nil {
				net_http.Error(w, err.Error(), net_http.StatusBadRequest)
				return
			}

			if x, ok := proto.Message(&protoReq).(interface{ Validate() error }); ok {
				if err := x.Validate(); err != nil {
					net_http.Error(w, err.Error(), net_http.StatusBadRequest)
					return
				}
			}

			if err := svc.Hello(&protoReq, &protoReply); err != nil {
				if pbgoErr, ok := err.(github_com_chai2010_pbgo.Error); ok {
					net_http.Error(w, pbgoErr.Text(), pbgoErr.HttpStatus())
					return
				} else {
					net_http.Error(w, err.Error(), net_http.StatusInternalServerError)
					return
				}
			}

			if x, ok := proto.Message(&protoReply).(interface{ Validate() error }); ok {
				if err := x.Validate(); err != nil {
					net_http.Error(w, err.Error(), net_http.StatusInternalServerError)
					return
				}
			}

			if strings.Contains(r.Header.Get("Accept"), "application/json") {
				w.Header().Set("Content-Type", "application/json")
			} else {
				w.Header().Set("Content-Type", "text/plain")
			}

			if err := encoding_json.NewEncoder(w).Encode(&protoReply); err != nil {
				net_http.Error(w, err.Error(), net_http.StatusInternalServerError)
				return
			}
		},
	)

	return router
}

func HelloServiceGrpcHandler(
	ctx context.Context, svc HelloServiceGrpcInterface,
	fnAnnotateContext func(ctx context.Context, req *net_http.Request, methodName string) (context.Context, error),
) net_http.Handler {
	var router = github_com_julienschmidt_httprouter.New()

	var re = regexp.MustCompile("(\\*|\\:)(\\w|\\.)+")
	_ = re

	router.Handle("GET", "/hello/:value",
		func(w net_http.ResponseWriter, r *net_http.Request, ps github_com_julienschmidt_httprouter.Params) {
			var (
				protoReq   String
				protoReply *String
				err        error
			)

			for _, fieldPath := range re.FindAllString("/hello/:value", -1) {
				fieldPath := strings.TrimLeft(fieldPath, ":*")
				err := github_com_chai2010_pbgo.PopulateFieldFromPath(&protoReq, fieldPath, ps.ByName(fieldPath))
				if err != nil {
					net_http.Error(w, err.Error(), net_http.StatusBadRequest)
					return
				}
			}

			if err := github_com_chai2010_pbgo.PopulateQueryParameters(&protoReq, r.URL.Query()); err != nil {
				net_http.Error(w, err.Error(), net_http.StatusBadRequest)
				return
			}

			if x, ok := proto.Message(&protoReq).(interface{ Validate() error }); ok {
				if err := x.Validate(); err != nil {
					net_http.Error(w, err.Error(), net_http.StatusBadRequest)
					return
				}
			}

			if fnAnnotateContext != nil {
				var err error
				ctx, err = fnAnnotateContext(ctx, r, "HelloService.Hello")
				if err != nil {
					net_http.Error(w, err.Error(), net_http.StatusBadRequest)
					return
				}
			}

			if protoReply, err = svc.Hello(ctx, &protoReq); err != nil {
				if pbgoErr, ok := err.(github_com_chai2010_pbgo.Error); ok {
					net_http.Error(w, pbgoErr.Text(), pbgoErr.HttpStatus())
					return
				} else {
					net_http.Error(w, err.Error(), net_http.StatusInternalServerError)
					return
				}
			}

			if x, ok := proto.Message(protoReply).(interface{ Validate() error }); ok {
				if err := x.Validate(); err != nil {
					net_http.Error(w, err.Error(), net_http.StatusInternalServerError)
					return
				}
			}

			if strings.Contains(r.Header.Get("Accept"), "application/json") {
				w.Header().Set("Content-Type", "application/json")
			} else {
				w.Header().Set("Content-Type", "text/plain")
			}

			if err := encoding_json.NewEncoder(w).Encode(&protoReply); err != nil {
				net_http.Error(w, err.Error(), net_http.StatusInternalServerError)
				return
			}
		},
	)

	return router
}