// Code generated by protoc-gen-go. DO NOT EDIT.
// source: task.proto

package task

import (
	context "context"
	fmt "fmt"
	annotatedvalue "github.com/AljabrIO/koalja-operator/pkg/annotatedvalue"
	proto "github.com/golang/protobuf/proto"
	duration "github.com/golang/protobuf/ptypes/duration"
	empty "github.com/golang/protobuf/ptypes/empty"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type NextRequest struct {
	// WaitTimeout is the maximum amount of time that the called
	// is prepared to wait for an answer.
	WaitTimeout          *duration.Duration `protobuf:"bytes,1,opt,name=WaitTimeout,json=waitTimeout,proto3" json:"WaitTimeout,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *NextRequest) Reset()         { *m = NextRequest{} }
func (m *NextRequest) String() string { return proto.CompactTextString(m) }
func (*NextRequest) ProtoMessage()    {}
func (*NextRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{0}
}

func (m *NextRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NextRequest.Unmarshal(m, b)
}
func (m *NextRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NextRequest.Marshal(b, m, deterministic)
}
func (m *NextRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NextRequest.Merge(m, src)
}
func (m *NextRequest) XXX_Size() int {
	return xxx_messageInfo_NextRequest.Size(m)
}
func (m *NextRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_NextRequest.DiscardUnknown(m)
}

var xxx_messageInfo_NextRequest proto.InternalMessageInfo

func (m *NextRequest) GetWaitTimeout() *duration.Duration {
	if m != nil {
		return m.WaitTimeout
	}
	return nil
}

type NextResponse struct {
	// If set, Snapshot holds the next snapshot.
	Snapshot *Snapshot `protobuf:"bytes,1,opt,name=Snapshot,json=snapshot,proto3" json:"Snapshot,omitempty"`
	// If Snapshot is not set, NoSnapshotYet is set to true to indicate
	// that was is no snapshot available within the timeout specified
	// in the request.
	NoSnapshotYet        bool     `protobuf:"varint,2,opt,name=NoSnapshotYet,json=noSnapshotYet,proto3" json:"NoSnapshotYet,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NextResponse) Reset()         { *m = NextResponse{} }
func (m *NextResponse) String() string { return proto.CompactTextString(m) }
func (*NextResponse) ProtoMessage()    {}
func (*NextResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{1}
}

func (m *NextResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NextResponse.Unmarshal(m, b)
}
func (m *NextResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NextResponse.Marshal(b, m, deterministic)
}
func (m *NextResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NextResponse.Merge(m, src)
}
func (m *NextResponse) XXX_Size() int {
	return xxx_messageInfo_NextResponse.Size(m)
}
func (m *NextResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_NextResponse.DiscardUnknown(m)
}

var xxx_messageInfo_NextResponse proto.InternalMessageInfo

func (m *NextResponse) GetSnapshot() *Snapshot {
	if m != nil {
		return m.Snapshot
	}
	return nil
}

func (m *NextResponse) GetNoSnapshotYet() bool {
	if m != nil {
		return m.NoSnapshotYet
	}
	return false
}

type AckRequest struct {
	SnapshotID           string   `protobuf:"bytes,1,opt,name=SnapshotID,json=snapshotID,proto3" json:"SnapshotID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AckRequest) Reset()         { *m = AckRequest{} }
func (m *AckRequest) String() string { return proto.CompactTextString(m) }
func (*AckRequest) ProtoMessage()    {}
func (*AckRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{2}
}

func (m *AckRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AckRequest.Unmarshal(m, b)
}
func (m *AckRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AckRequest.Marshal(b, m, deterministic)
}
func (m *AckRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AckRequest.Merge(m, src)
}
func (m *AckRequest) XXX_Size() int {
	return xxx_messageInfo_AckRequest.Size(m)
}
func (m *AckRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AckRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AckRequest proto.InternalMessageInfo

func (m *AckRequest) GetSnapshotID() string {
	if m != nil {
		return m.SnapshotID
	}
	return ""
}

type ExecuteTemplateRequest struct {
	// Current snapshot.
	Snapshot *Snapshot `protobuf:"bytes,1,opt,name=Snapshot,json=snapshot,proto3" json:"Snapshot,omitempty"`
	// The template (source)
	Template             string   `protobuf:"bytes,2,opt,name=Template,json=template,proto3" json:"Template,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExecuteTemplateRequest) Reset()         { *m = ExecuteTemplateRequest{} }
func (m *ExecuteTemplateRequest) String() string { return proto.CompactTextString(m) }
func (*ExecuteTemplateRequest) ProtoMessage()    {}
func (*ExecuteTemplateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{3}
}

func (m *ExecuteTemplateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExecuteTemplateRequest.Unmarshal(m, b)
}
func (m *ExecuteTemplateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExecuteTemplateRequest.Marshal(b, m, deterministic)
}
func (m *ExecuteTemplateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExecuteTemplateRequest.Merge(m, src)
}
func (m *ExecuteTemplateRequest) XXX_Size() int {
	return xxx_messageInfo_ExecuteTemplateRequest.Size(m)
}
func (m *ExecuteTemplateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ExecuteTemplateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ExecuteTemplateRequest proto.InternalMessageInfo

func (m *ExecuteTemplateRequest) GetSnapshot() *Snapshot {
	if m != nil {
		return m.Snapshot
	}
	return nil
}

func (m *ExecuteTemplateRequest) GetTemplate() string {
	if m != nil {
		return m.Template
	}
	return ""
}

type ExecuteTemplateResponse struct {
	// Result of the template execution.
	Result               []byte   `protobuf:"bytes,1,opt,name=Result,json=result,proto3" json:"Result,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExecuteTemplateResponse) Reset()         { *m = ExecuteTemplateResponse{} }
func (m *ExecuteTemplateResponse) String() string { return proto.CompactTextString(m) }
func (*ExecuteTemplateResponse) ProtoMessage()    {}
func (*ExecuteTemplateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{4}
}

func (m *ExecuteTemplateResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExecuteTemplateResponse.Unmarshal(m, b)
}
func (m *ExecuteTemplateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExecuteTemplateResponse.Marshal(b, m, deterministic)
}
func (m *ExecuteTemplateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExecuteTemplateResponse.Merge(m, src)
}
func (m *ExecuteTemplateResponse) XXX_Size() int {
	return xxx_messageInfo_ExecuteTemplateResponse.Size(m)
}
func (m *ExecuteTemplateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ExecuteTemplateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ExecuteTemplateResponse proto.InternalMessageInfo

func (m *ExecuteTemplateResponse) GetResult() []byte {
	if m != nil {
		return m.Result
	}
	return nil
}

// Snapshot is a set of sequences of annotated values for every
// input of a task.
type Snapshot struct {
	ID                   string               `protobuf:"bytes,1,opt,name=ID,json=iD,proto3" json:"ID,omitempty"`
	Inputs               []*SnapshotInputPair `protobuf:"bytes,2,rep,name=Inputs,json=inputs,proto3" json:"Inputs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Snapshot) Reset()         { *m = Snapshot{} }
func (m *Snapshot) String() string { return proto.CompactTextString(m) }
func (*Snapshot) ProtoMessage()    {}
func (*Snapshot) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{5}
}

func (m *Snapshot) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Snapshot.Unmarshal(m, b)
}
func (m *Snapshot) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Snapshot.Marshal(b, m, deterministic)
}
func (m *Snapshot) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Snapshot.Merge(m, src)
}
func (m *Snapshot) XXX_Size() int {
	return xxx_messageInfo_Snapshot.Size(m)
}
func (m *Snapshot) XXX_DiscardUnknown() {
	xxx_messageInfo_Snapshot.DiscardUnknown(m)
}

var xxx_messageInfo_Snapshot proto.InternalMessageInfo

func (m *Snapshot) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *Snapshot) GetInputs() []*SnapshotInputPair {
	if m != nil {
		return m.Inputs
	}
	return nil
}

// SnapshotInputPair is a sequences of annotated value for a specific
// task input with given name.
type SnapshotInputPair struct {
	InputName            string                           `protobuf:"bytes,1,opt,name=InputName,json=inputName,proto3" json:"InputName,omitempty"`
	AnnotatedValues      []*annotatedvalue.AnnotatedValue `protobuf:"bytes,2,rep,name=AnnotatedValues,json=annotatedValues,proto3" json:"AnnotatedValues,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                         `json:"-"`
	XXX_unrecognized     []byte                           `json:"-"`
	XXX_sizecache        int32                            `json:"-"`
}

func (m *SnapshotInputPair) Reset()         { *m = SnapshotInputPair{} }
func (m *SnapshotInputPair) String() string { return proto.CompactTextString(m) }
func (*SnapshotInputPair) ProtoMessage()    {}
func (*SnapshotInputPair) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{6}
}

func (m *SnapshotInputPair) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SnapshotInputPair.Unmarshal(m, b)
}
func (m *SnapshotInputPair) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SnapshotInputPair.Marshal(b, m, deterministic)
}
func (m *SnapshotInputPair) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SnapshotInputPair.Merge(m, src)
}
func (m *SnapshotInputPair) XXX_Size() int {
	return xxx_messageInfo_SnapshotInputPair.Size(m)
}
func (m *SnapshotInputPair) XXX_DiscardUnknown() {
	xxx_messageInfo_SnapshotInputPair.DiscardUnknown(m)
}

var xxx_messageInfo_SnapshotInputPair proto.InternalMessageInfo

func (m *SnapshotInputPair) GetInputName() string {
	if m != nil {
		return m.InputName
	}
	return ""
}

func (m *SnapshotInputPair) GetAnnotatedValues() []*annotatedvalue.AnnotatedValue {
	if m != nil {
		return m.AnnotatedValues
	}
	return nil
}

type OutputReadyRequest struct {
	// Data (content) of the annotated value.
	AnnotatedValueData string `protobuf:"bytes,1,opt,name=AnnotatedValueData,json=annotatedValueData,proto3" json:"AnnotatedValueData,omitempty"`
	// Name of the task output that is data belongs to.
	OutputName string `protobuf:"bytes,2,opt,name=OutputName,json=outputName,proto3" json:"OutputName,omitempty"`
	// Optional snapshot to be passed for custom task executors.
	Snapshot             *Snapshot `protobuf:"bytes,3,opt,name=Snapshot,json=snapshot,proto3" json:"Snapshot,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *OutputReadyRequest) Reset()         { *m = OutputReadyRequest{} }
func (m *OutputReadyRequest) String() string { return proto.CompactTextString(m) }
func (*OutputReadyRequest) ProtoMessage()    {}
func (*OutputReadyRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{7}
}

func (m *OutputReadyRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OutputReadyRequest.Unmarshal(m, b)
}
func (m *OutputReadyRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OutputReadyRequest.Marshal(b, m, deterministic)
}
func (m *OutputReadyRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OutputReadyRequest.Merge(m, src)
}
func (m *OutputReadyRequest) XXX_Size() int {
	return xxx_messageInfo_OutputReadyRequest.Size(m)
}
func (m *OutputReadyRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_OutputReadyRequest.DiscardUnknown(m)
}

var xxx_messageInfo_OutputReadyRequest proto.InternalMessageInfo

func (m *OutputReadyRequest) GetAnnotatedValueData() string {
	if m != nil {
		return m.AnnotatedValueData
	}
	return ""
}

func (m *OutputReadyRequest) GetOutputName() string {
	if m != nil {
		return m.OutputName
	}
	return ""
}

func (m *OutputReadyRequest) GetSnapshot() *Snapshot {
	if m != nil {
		return m.Snapshot
	}
	return nil
}

type OutputReadyResponse struct {
	// Accepted is true if the task agent accepted the output.
	// If false, the executor should try to notify the task agent
	// again after a timeout.
	Accepted bool `protobuf:"varint,1,opt,name=Accepted,json=accepted,proto3" json:"Accepted,omitempty"`
	// ID of the published annotated value
	AnnotatedValueID     string   `protobuf:"bytes,2,opt,name=AnnotatedValueID,json=annotatedValueID,proto3" json:"AnnotatedValueID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OutputReadyResponse) Reset()         { *m = OutputReadyResponse{} }
func (m *OutputReadyResponse) String() string { return proto.CompactTextString(m) }
func (*OutputReadyResponse) ProtoMessage()    {}
func (*OutputReadyResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{8}
}

func (m *OutputReadyResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OutputReadyResponse.Unmarshal(m, b)
}
func (m *OutputReadyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OutputReadyResponse.Marshal(b, m, deterministic)
}
func (m *OutputReadyResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OutputReadyResponse.Merge(m, src)
}
func (m *OutputReadyResponse) XXX_Size() int {
	return xxx_messageInfo_OutputReadyResponse.Size(m)
}
func (m *OutputReadyResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_OutputReadyResponse.DiscardUnknown(m)
}

var xxx_messageInfo_OutputReadyResponse proto.InternalMessageInfo

func (m *OutputReadyResponse) GetAccepted() bool {
	if m != nil {
		return m.Accepted
	}
	return false
}

func (m *OutputReadyResponse) GetAnnotatedValueID() string {
	if m != nil {
		return m.AnnotatedValueID
	}
	return ""
}

type CreateFileURIRequest struct {
	// Name of the task output that is data belongs to.
	OutputName string `protobuf:"bytes,2,opt,name=OutputName,json=outputName,proto3" json:"OutputName,omitempty"`
	// Local path of the file/dir in the Volume
	LocalPath string `protobuf:"bytes,7,opt,name=LocalPath,json=localPath,proto3" json:"LocalPath,omitempty"`
	// IsDir indicates if the URI is for a file (false) or a directory (true)
	IsDir                bool     `protobuf:"varint,8,opt,name=IsDir,json=isDir,proto3" json:"IsDir,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateFileURIRequest) Reset()         { *m = CreateFileURIRequest{} }
func (m *CreateFileURIRequest) String() string { return proto.CompactTextString(m) }
func (*CreateFileURIRequest) ProtoMessage()    {}
func (*CreateFileURIRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{9}
}

func (m *CreateFileURIRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateFileURIRequest.Unmarshal(m, b)
}
func (m *CreateFileURIRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateFileURIRequest.Marshal(b, m, deterministic)
}
func (m *CreateFileURIRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateFileURIRequest.Merge(m, src)
}
func (m *CreateFileURIRequest) XXX_Size() int {
	return xxx_messageInfo_CreateFileURIRequest.Size(m)
}
func (m *CreateFileURIRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateFileURIRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateFileURIRequest proto.InternalMessageInfo

func (m *CreateFileURIRequest) GetOutputName() string {
	if m != nil {
		return m.OutputName
	}
	return ""
}

func (m *CreateFileURIRequest) GetLocalPath() string {
	if m != nil {
		return m.LocalPath
	}
	return ""
}

func (m *CreateFileURIRequest) GetIsDir() bool {
	if m != nil {
		return m.IsDir
	}
	return false
}

type CreateFileURIResponse struct {
	// The created URI
	URI                  string   `protobuf:"bytes,1,opt,name=URI,json=uRI,proto3" json:"URI,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateFileURIResponse) Reset()         { *m = CreateFileURIResponse{} }
func (m *CreateFileURIResponse) String() string { return proto.CompactTextString(m) }
func (*CreateFileURIResponse) ProtoMessage()    {}
func (*CreateFileURIResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{10}
}

func (m *CreateFileURIResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateFileURIResponse.Unmarshal(m, b)
}
func (m *CreateFileURIResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateFileURIResponse.Marshal(b, m, deterministic)
}
func (m *CreateFileURIResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateFileURIResponse.Merge(m, src)
}
func (m *CreateFileURIResponse) XXX_Size() int {
	return xxx_messageInfo_CreateFileURIResponse.Size(m)
}
func (m *CreateFileURIResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateFileURIResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateFileURIResponse proto.InternalMessageInfo

func (m *CreateFileURIResponse) GetURI() string {
	if m != nil {
		return m.URI
	}
	return ""
}

func init() {
	proto.RegisterType((*NextRequest)(nil), "task.NextRequest")
	proto.RegisterType((*NextResponse)(nil), "task.NextResponse")
	proto.RegisterType((*AckRequest)(nil), "task.AckRequest")
	proto.RegisterType((*ExecuteTemplateRequest)(nil), "task.ExecuteTemplateRequest")
	proto.RegisterType((*ExecuteTemplateResponse)(nil), "task.ExecuteTemplateResponse")
	proto.RegisterType((*Snapshot)(nil), "task.Snapshot")
	proto.RegisterType((*SnapshotInputPair)(nil), "task.SnapshotInputPair")
	proto.RegisterType((*OutputReadyRequest)(nil), "task.OutputReadyRequest")
	proto.RegisterType((*OutputReadyResponse)(nil), "task.OutputReadyResponse")
	proto.RegisterType((*CreateFileURIRequest)(nil), "task.CreateFileURIRequest")
	proto.RegisterType((*CreateFileURIResponse)(nil), "task.CreateFileURIResponse")
}

func init() { proto.RegisterFile("task.proto", fileDescriptor_ce5d8dd45b4a91ff) }

var fileDescriptor_ce5d8dd45b4a91ff = []byte{
	// 698 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0xdd, 0x4e, 0xdb, 0x48,
	0x14, 0x56, 0x12, 0xc8, 0x3a, 0x27, 0x40, 0xc2, 0x2c, 0x0b, 0xc1, 0xcb, 0x22, 0x64, 0xed, 0x05,
	0x8b, 0xc0, 0xd6, 0x86, 0xab, 0xaa, 0x57, 0xa1, 0xa1, 0xc2, 0xa5, 0x0a, 0x68, 0x80, 0x56, 0x54,
	0xaa, 0xc4, 0xc4, 0x19, 0xc2, 0x10, 0xc7, 0xe3, 0xda, 0x63, 0x4a, 0xd4, 0x97, 0xe8, 0x2b, 0xf5,
	0xcd, 0x2a, 0x7b, 0x66, 0x42, 0x9c, 0xa4, 0xb4, 0xbd, 0xb1, 0xe6, 0x7c, 0x67, 0xce, 0x77, 0xbe,
	0x39, 0x3f, 0x06, 0x10, 0x24, 0x1e, 0xd8, 0x61, 0xc4, 0x05, 0x47, 0x0b, 0xe9, 0xd9, 0x7c, 0xd1,
	0x67, 0xe2, 0x2e, 0xe9, 0xda, 0x1e, 0x1f, 0x3a, 0x7d, 0xee, 0x93, 0xa0, 0xef, 0x64, 0xee, 0x6e,
	0x72, 0xeb, 0x84, 0x62, 0x14, 0xd2, 0xd8, 0xe9, 0x25, 0x11, 0x11, 0x8c, 0x07, 0xe3, 0x83, 0x24,
	0x30, 0x0f, 0x7f, 0x1e, 0x4a, 0x87, 0xa1, 0x18, 0xc9, 0xaf, 0x0a, 0x3a, 0x9d, 0x08, 0x6a, 0xf9,
	0xf7, 0xa4, 0x1b, 0xb9, 0x67, 0xce, 0x80, 0x13, 0xff, 0x9e, 0x1c, 0xf0, 0x90, 0x46, 0x44, 0xf0,
	0xc8, 0x09, 0x07, 0x7d, 0x87, 0x04, 0x01, 0x17, 0x44, 0xd0, 0xde, 0x03, 0xf1, 0x13, 0x3a, 0x65,
	0x4a, 0x32, 0xeb, 0x0d, 0x54, 0x3b, 0xf4, 0x51, 0x60, 0xfa, 0x29, 0xa1, 0xb1, 0x40, 0x2f, 0xa1,
	0xfa, 0x9e, 0x30, 0x71, 0xc9, 0x86, 0x94, 0x27, 0xa2, 0x51, 0xd8, 0x29, 0xec, 0x56, 0x9b, 0x9b,
	0x76, 0x9f, 0xf3, 0xbe, 0xaf, 0x42, 0xba, 0xc9, 0xad, 0xdd, 0x56, 0xcf, 0xc0, 0xd5, 0xcf, 0x4f,
	0xb7, 0xad, 0x1b, 0x58, 0x92, 0x5c, 0x71, 0xc8, 0x83, 0x98, 0xa2, 0x3d, 0x30, 0x2e, 0x02, 0x12,
	0xc6, 0x77, 0x5c, 0x33, 0xad, 0xd8, 0x59, 0xf5, 0x34, 0x8a, 0x8d, 0x58, 0x9d, 0xd0, 0xbf, 0xb0,
	0xdc, 0xe1, 0x1a, 0xbf, 0xa6, 0xa2, 0x51, 0xdc, 0x29, 0xec, 0x1a, 0x78, 0x39, 0x98, 0x04, 0xad,
	0x7d, 0x80, 0x96, 0x37, 0xd0, 0x62, 0xb7, 0x01, 0xb4, 0xd3, 0x6d, 0x67, 0x19, 0x2a, 0x18, 0xe2,
	0x31, 0x62, 0xdd, 0xc0, 0xfa, 0xf1, 0x23, 0xf5, 0x12, 0x41, 0x2f, 0xe9, 0x30, 0xf4, 0x89, 0xa0,
	0x3a, 0xf2, 0x77, 0x94, 0x99, 0x60, 0xe8, 0xf0, 0x4c, 0x54, 0x05, 0x1b, 0x42, 0xd9, 0xd6, 0xff,
	0xb0, 0x31, 0x93, 0x41, 0x3d, 0x7e, 0x1d, 0xca, 0x98, 0xc6, 0x89, 0x2f, 0x13, 0x2c, 0xe1, 0x72,
	0x94, 0x59, 0xd6, 0xe9, 0x53, 0x6a, 0xb4, 0x02, 0xc5, 0xb1, 0xf0, 0x22, 0x6b, 0x23, 0x07, 0xca,
	0x6e, 0x10, 0x26, 0x22, 0x6e, 0x14, 0x77, 0x4a, 0xbb, 0xd5, 0xe6, 0x46, 0x5e, 0x54, 0xe6, 0x3b,
	0x27, 0x2c, 0xc2, 0x65, 0x96, 0x5d, 0xb3, 0xbe, 0xc0, 0xea, 0x8c, 0x13, 0x6d, 0x41, 0x25, 0x33,
	0x3a, 0x64, 0x48, 0x15, 0x79, 0x85, 0x69, 0x00, 0x9d, 0x40, 0xad, 0xa5, 0x07, 0xe1, 0x5d, 0x3a,
	0x08, 0x3a, 0xd9, 0xb6, 0x3d, 0x35, 0x20, 0xf9, 0x6b, 0xb8, 0x46, 0xf2, 0x61, 0xd6, 0xd7, 0x02,
	0xa0, 0xb3, 0x44, 0x84, 0x89, 0xc0, 0x94, 0xf4, 0x46, 0xba, 0xb6, 0x36, 0xa0, 0x7c, 0x64, 0x9b,
	0x08, 0xa2, 0x74, 0x20, 0x32, 0xe3, 0x49, 0xbb, 0x28, 0x59, 0x32, 0xbd, 0xb2, 0xc2, 0xc0, 0xc7,
	0x48, 0xae, 0x57, 0xa5, 0xe7, 0x7b, 0x65, 0x7d, 0x84, 0x3f, 0x73, 0x8a, 0x54, 0x2f, 0x4c, 0x30,
	0x5a, 0x9e, 0x47, 0x43, 0x41, 0x7b, 0x99, 0x10, 0x03, 0x1b, 0x44, 0xd9, 0x68, 0x0f, 0xea, 0x79,
	0xb9, 0x6e, 0x5b, 0x89, 0xa8, 0x93, 0x29, 0xdc, 0xba, 0x87, 0xb5, 0x57, 0x11, 0x25, 0x82, 0xbe,
	0x66, 0x3e, 0xbd, 0xc2, 0xee, 0xc4, 0x20, 0x3e, 0xfb, 0x84, 0x2d, 0xa8, 0xbc, 0xe5, 0x1e, 0xf1,
	0xcf, 0x89, 0xb8, 0x6b, 0xfc, 0x21, 0x3b, 0xe2, 0x6b, 0x00, 0xad, 0xc1, 0xa2, 0x1b, 0xb7, 0x59,
	0xd4, 0x30, 0x32, 0x69, 0x8b, 0x2c, 0x35, 0xac, 0xff, 0xe0, 0xaf, 0xa9, 0x5c, 0xea, 0x31, 0x75,
	0x28, 0x5d, 0x61, 0x57, 0x15, 0xb4, 0x94, 0x60, 0xb7, 0xf9, 0xad, 0x00, 0x35, 0x5d, 0x8c, 0x0b,
	0x1a, 0x3d, 0x30, 0x8f, 0xa2, 0x03, 0x58, 0x48, 0x77, 0x11, 0xad, 0xca, 0x5a, 0x4d, 0xec, 0xb8,
	0x89, 0x26, 0x21, 0x45, 0xea, 0x40, 0xa9, 0xe5, 0x0d, 0x50, 0x5d, 0xba, 0x9e, 0x76, 0xcc, 0x5c,
	0x9f, 0xd9, 0xfd, 0xe3, 0xf4, 0x57, 0x84, 0x3a, 0x50, 0x9b, 0x9a, 0x7c, 0xb4, 0x25, 0x83, 0xe7,
	0xaf, 0x9c, 0xf9, 0xcf, 0x0f, 0xbc, 0x52, 0x40, 0xf3, 0x3a, 0xd7, 0xb9, 0x0e, 0x17, 0xec, 0x96,
	0xd1, 0x08, 0x1d, 0x41, 0x75, 0x02, 0x46, 0x0d, 0x49, 0x32, 0x3b, 0x75, 0xe6, 0xe6, 0x1c, 0x8f,
	0xa2, 0xf6, 0x60, 0x43, 0xc2, 0x69, 0x25, 0x2f, 0x46, 0xb1, 0xa0, 0x43, 0x5d, 0xa5, 0x13, 0x58,
	0xce, 0x15, 0x19, 0x99, 0x92, 0x66, 0x5e, 0x97, 0xcd, 0xbf, 0xe7, 0xfa, 0x64, 0x92, 0x23, 0xfb,
	0xc3, 0xfe, 0xaf, 0xfe, 0x96, 0x53, 0x96, 0x6e, 0x39, 0xab, 0xe7, 0xe1, 0xf7, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x07, 0x2e, 0x5e, 0x82, 0x4f, 0x06, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SnapshotServiceClient is the client API for SnapshotService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SnapshotServiceClient interface {
	// Next pulls the task agent for the next available snapshot.
	Next(ctx context.Context, in *NextRequest, opts ...grpc.CallOption) (*NextResponse, error)
	// Acknowledge the processing of a snapshot
	Ack(ctx context.Context, in *AckRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	// ExecuteTemplate is invoked to parse & execute a template
	// with given snapshot.
	ExecuteTemplate(ctx context.Context, in *ExecuteTemplateRequest, opts ...grpc.CallOption) (*ExecuteTemplateResponse, error)
}

type snapshotServiceClient struct {
	cc *grpc.ClientConn
}

func NewSnapshotServiceClient(cc *grpc.ClientConn) SnapshotServiceClient {
	return &snapshotServiceClient{cc}
}

func (c *snapshotServiceClient) Next(ctx context.Context, in *NextRequest, opts ...grpc.CallOption) (*NextResponse, error) {
	out := new(NextResponse)
	err := c.cc.Invoke(ctx, "/task.SnapshotService/Next", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *snapshotServiceClient) Ack(ctx context.Context, in *AckRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/task.SnapshotService/Ack", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *snapshotServiceClient) ExecuteTemplate(ctx context.Context, in *ExecuteTemplateRequest, opts ...grpc.CallOption) (*ExecuteTemplateResponse, error) {
	out := new(ExecuteTemplateResponse)
	err := c.cc.Invoke(ctx, "/task.SnapshotService/ExecuteTemplate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SnapshotServiceServer is the server API for SnapshotService service.
type SnapshotServiceServer interface {
	// Next pulls the task agent for the next available snapshot.
	Next(context.Context, *NextRequest) (*NextResponse, error)
	// Acknowledge the processing of a snapshot
	Ack(context.Context, *AckRequest) (*empty.Empty, error)
	// ExecuteTemplate is invoked to parse & execute a template
	// with given snapshot.
	ExecuteTemplate(context.Context, *ExecuteTemplateRequest) (*ExecuteTemplateResponse, error)
}

func RegisterSnapshotServiceServer(s *grpc.Server, srv SnapshotServiceServer) {
	s.RegisterService(&_SnapshotService_serviceDesc, srv)
}

func _SnapshotService_Next_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NextRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SnapshotServiceServer).Next(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/task.SnapshotService/Next",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SnapshotServiceServer).Next(ctx, req.(*NextRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SnapshotService_Ack_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SnapshotServiceServer).Ack(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/task.SnapshotService/Ack",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SnapshotServiceServer).Ack(ctx, req.(*AckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SnapshotService_ExecuteTemplate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecuteTemplateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SnapshotServiceServer).ExecuteTemplate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/task.SnapshotService/ExecuteTemplate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SnapshotServiceServer).ExecuteTemplate(ctx, req.(*ExecuteTemplateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SnapshotService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "task.SnapshotService",
	HandlerType: (*SnapshotServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Next",
			Handler:    _SnapshotService_Next_Handler,
		},
		{
			MethodName: "Ack",
			Handler:    _SnapshotService_Ack_Handler,
		},
		{
			MethodName: "ExecuteTemplate",
			Handler:    _SnapshotService_ExecuteTemplate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "task.proto",
}

// OutputReadyNotifierClient is the client API for OutputReadyNotifier service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type OutputReadyNotifierClient interface {
	// OutputReady signals the task agent that an output is ready for publication
	// into a link.
	OutputReady(ctx context.Context, in *OutputReadyRequest, opts ...grpc.CallOption) (*OutputReadyResponse, error)
}

type outputReadyNotifierClient struct {
	cc *grpc.ClientConn
}

func NewOutputReadyNotifierClient(cc *grpc.ClientConn) OutputReadyNotifierClient {
	return &outputReadyNotifierClient{cc}
}

func (c *outputReadyNotifierClient) OutputReady(ctx context.Context, in *OutputReadyRequest, opts ...grpc.CallOption) (*OutputReadyResponse, error) {
	out := new(OutputReadyResponse)
	err := c.cc.Invoke(ctx, "/task.OutputReadyNotifier/OutputReady", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OutputReadyNotifierServer is the server API for OutputReadyNotifier service.
type OutputReadyNotifierServer interface {
	// OutputReady signals the task agent that an output is ready for publication
	// into a link.
	OutputReady(context.Context, *OutputReadyRequest) (*OutputReadyResponse, error)
}

func RegisterOutputReadyNotifierServer(s *grpc.Server, srv OutputReadyNotifierServer) {
	s.RegisterService(&_OutputReadyNotifier_serviceDesc, srv)
}

func _OutputReadyNotifier_OutputReady_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OutputReadyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OutputReadyNotifierServer).OutputReady(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/task.OutputReadyNotifier/OutputReady",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OutputReadyNotifierServer).OutputReady(ctx, req.(*OutputReadyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _OutputReadyNotifier_serviceDesc = grpc.ServiceDesc{
	ServiceName: "task.OutputReadyNotifier",
	HandlerType: (*OutputReadyNotifierServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "OutputReady",
			Handler:    _OutputReadyNotifier_OutputReady_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "task.proto",
}

// OutputFileSystemServiceClient is the client API for OutputFileSystemService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type OutputFileSystemServiceClient interface {
	// CreateFileURI creates a URI for the given file/dir
	CreateFileURI(ctx context.Context, in *CreateFileURIRequest, opts ...grpc.CallOption) (*CreateFileURIResponse, error)
}

type outputFileSystemServiceClient struct {
	cc *grpc.ClientConn
}

func NewOutputFileSystemServiceClient(cc *grpc.ClientConn) OutputFileSystemServiceClient {
	return &outputFileSystemServiceClient{cc}
}

func (c *outputFileSystemServiceClient) CreateFileURI(ctx context.Context, in *CreateFileURIRequest, opts ...grpc.CallOption) (*CreateFileURIResponse, error) {
	out := new(CreateFileURIResponse)
	err := c.cc.Invoke(ctx, "/task.OutputFileSystemService/CreateFileURI", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OutputFileSystemServiceServer is the server API for OutputFileSystemService service.
type OutputFileSystemServiceServer interface {
	// CreateFileURI creates a URI for the given file/dir
	CreateFileURI(context.Context, *CreateFileURIRequest) (*CreateFileURIResponse, error)
}

func RegisterOutputFileSystemServiceServer(s *grpc.Server, srv OutputFileSystemServiceServer) {
	s.RegisterService(&_OutputFileSystemService_serviceDesc, srv)
}

func _OutputFileSystemService_CreateFileURI_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateFileURIRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OutputFileSystemServiceServer).CreateFileURI(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/task.OutputFileSystemService/CreateFileURI",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OutputFileSystemServiceServer).CreateFileURI(ctx, req.(*CreateFileURIRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _OutputFileSystemService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "task.OutputFileSystemService",
	HandlerType: (*OutputFileSystemServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateFileURI",
			Handler:    _OutputFileSystemService_CreateFileURI_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "task.proto",
}