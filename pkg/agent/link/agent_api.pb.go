// Code generated by protoc-gen-go. DO NOT EDIT.
// source: agent_api.proto

package link

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	context "golang.org/x/net/context"
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

func init() { proto.RegisterFile("agent_api.proto", fileDescriptor_011b677ced0569e6) }

var fileDescriptor_011b677ced0569e6 = []byte{
	// 59 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4f, 0x4c, 0x4f, 0xcd,
	0x2b, 0x89, 0x4f, 0x2c, 0xc8, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xc9, 0xc9, 0xcc,
	0xcb, 0x36, 0x62, 0xe7, 0x62, 0x75, 0x04, 0x49, 0x24, 0xb1, 0x81, 0x45, 0x8d, 0x01, 0x01, 0x00,
	0x00, 0xff, 0xff, 0xf4, 0xb0, 0x63, 0xe6, 0x28, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// AgentClient is the client API for Agent service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AgentClient interface {
}

type agentClient struct {
	cc *grpc.ClientConn
}

func NewAgentClient(cc *grpc.ClientConn) AgentClient {
	return &agentClient{cc}
}

// AgentServer is the server API for Agent service.
type AgentServer interface {
}

func RegisterAgentServer(s *grpc.Server, srv AgentServer) {
	s.RegisterService(&_Agent_serviceDesc, srv)
}

var _Agent_serviceDesc = grpc.ServiceDesc{
	ServiceName: "link.Agent",
	HandlerType: (*AgentServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "agent_api.proto",
}