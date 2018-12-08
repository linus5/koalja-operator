// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/devtools/containeranalysis/v1beta1/package/package.proto

package _package // import "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/package"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Instruction set architectures supported by various package managers.
type Architecture int32

const (
	// Unknown architecture.
	Architecture_ARCHITECTURE_UNSPECIFIED Architecture = 0
	// X86 architecture.
	Architecture_X86 Architecture = 1
	// X64 architecture.
	Architecture_X64 Architecture = 2
)

var Architecture_name = map[int32]string{
	0: "ARCHITECTURE_UNSPECIFIED",
	1: "X86",
	2: "X64",
}
var Architecture_value = map[string]int32{
	"ARCHITECTURE_UNSPECIFIED": 0,
	"X86":                      1,
	"X64":                      2,
}

func (x Architecture) String() string {
	return proto.EnumName(Architecture_name, int32(x))
}
func (Architecture) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_package_98e063c9654a5d86, []int{0}
}

// Whether this is an ordinary package version or a sentinel MIN/MAX version.
type Version_VersionKind int32

const (
	// Unknown.
	Version_VERSION_KIND_UNSPECIFIED Version_VersionKind = 0
	// A standard package version, defined by the other fields.
	Version_NORMAL Version_VersionKind = 1
	// A special version representing negative infinity, other fields are
	// ignored.
	Version_MINIMUM Version_VersionKind = 2
	// A special version representing positive infinity, other fields are
	// ignored.
	Version_MAXIMUM Version_VersionKind = 3
)

var Version_VersionKind_name = map[int32]string{
	0: "VERSION_KIND_UNSPECIFIED",
	1: "NORMAL",
	2: "MINIMUM",
	3: "MAXIMUM",
}
var Version_VersionKind_value = map[string]int32{
	"VERSION_KIND_UNSPECIFIED": 0,
	"NORMAL":                   1,
	"MINIMUM":                  2,
	"MAXIMUM":                  3,
}

func (x Version_VersionKind) String() string {
	return proto.EnumName(Version_VersionKind_name, int32(x))
}
func (Version_VersionKind) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_package_98e063c9654a5d86, []int{5, 0}
}

// This represents a particular channel of distribution for a given package.
// E.g., Debian's jessie-backports dpkg mirror.
type Distribution struct {
	// The cpe_uri in [cpe format](https://cpe.mitre.org/specification/)
	// denoting the package manager version distributing a package.
	CpeUri string `protobuf:"bytes,1,opt,name=cpe_uri,json=cpeUri,proto3" json:"cpe_uri,omitempty"`
	// The CPU architecture for which packages in this distribution channel were
	// built.
	Architecture Architecture `protobuf:"varint,2,opt,name=architecture,proto3,enum=grafeas.v1beta1.package.Architecture" json:"architecture,omitempty"`
	// The latest available version of this package in this distribution
	// channel.
	LatestVersion *Version `protobuf:"bytes,3,opt,name=latest_version,json=latestVersion,proto3" json:"latest_version,omitempty"`
	// A freeform string denoting the maintainer of this package.
	Maintainer string `protobuf:"bytes,4,opt,name=maintainer,proto3" json:"maintainer,omitempty"`
	// The distribution channel-specific homepage for this package.
	Url string `protobuf:"bytes,5,opt,name=url,proto3" json:"url,omitempty"`
	// The distribution channel-specific description of this package.
	Description          string   `protobuf:"bytes,6,opt,name=description,proto3" json:"description,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Distribution) Reset()         { *m = Distribution{} }
func (m *Distribution) String() string { return proto.CompactTextString(m) }
func (*Distribution) ProtoMessage()    {}
func (*Distribution) Descriptor() ([]byte, []int) {
	return fileDescriptor_package_98e063c9654a5d86, []int{0}
}
func (m *Distribution) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Distribution.Unmarshal(m, b)
}
func (m *Distribution) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Distribution.Marshal(b, m, deterministic)
}
func (dst *Distribution) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Distribution.Merge(dst, src)
}
func (m *Distribution) XXX_Size() int {
	return xxx_messageInfo_Distribution.Size(m)
}
func (m *Distribution) XXX_DiscardUnknown() {
	xxx_messageInfo_Distribution.DiscardUnknown(m)
}

var xxx_messageInfo_Distribution proto.InternalMessageInfo

func (m *Distribution) GetCpeUri() string {
	if m != nil {
		return m.CpeUri
	}
	return ""
}

func (m *Distribution) GetArchitecture() Architecture {
	if m != nil {
		return m.Architecture
	}
	return Architecture_ARCHITECTURE_UNSPECIFIED
}

func (m *Distribution) GetLatestVersion() *Version {
	if m != nil {
		return m.LatestVersion
	}
	return nil
}

func (m *Distribution) GetMaintainer() string {
	if m != nil {
		return m.Maintainer
	}
	return ""
}

func (m *Distribution) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *Distribution) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

// An occurrence of a particular package installation found within a system's
// filesystem. E.g., glibc was found in /var/lib/dpkg/status.
type Location struct {
	// The cpe_uri in [cpe format](https://cpe.mitre.org/specification/)
	// denoting the package manager version distributing a package.
	CpeUri string `protobuf:"bytes,1,opt,name=cpe_uri,json=cpeUri,proto3" json:"cpe_uri,omitempty"`
	// The version installed at this location.
	Version *Version `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	// The path from which we gathered that this package/version is installed.
	Path                 string   `protobuf:"bytes,3,opt,name=path,proto3" json:"path,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Location) Reset()         { *m = Location{} }
func (m *Location) String() string { return proto.CompactTextString(m) }
func (*Location) ProtoMessage()    {}
func (*Location) Descriptor() ([]byte, []int) {
	return fileDescriptor_package_98e063c9654a5d86, []int{1}
}
func (m *Location) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Location.Unmarshal(m, b)
}
func (m *Location) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Location.Marshal(b, m, deterministic)
}
func (dst *Location) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Location.Merge(dst, src)
}
func (m *Location) XXX_Size() int {
	return xxx_messageInfo_Location.Size(m)
}
func (m *Location) XXX_DiscardUnknown() {
	xxx_messageInfo_Location.DiscardUnknown(m)
}

var xxx_messageInfo_Location proto.InternalMessageInfo

func (m *Location) GetCpeUri() string {
	if m != nil {
		return m.CpeUri
	}
	return ""
}

func (m *Location) GetVersion() *Version {
	if m != nil {
		return m.Version
	}
	return nil
}

func (m *Location) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

// This represents a particular package that is distributed over various
// channels. E.g., glibc (aka libc6) is distributed by many, at various
// versions.
type Package struct {
	// The name of the package.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The various channels by which a package is distributed.
	Distribution         []*Distribution `protobuf:"bytes,10,rep,name=distribution,proto3" json:"distribution,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Package) Reset()         { *m = Package{} }
func (m *Package) String() string { return proto.CompactTextString(m) }
func (*Package) ProtoMessage()    {}
func (*Package) Descriptor() ([]byte, []int) {
	return fileDescriptor_package_98e063c9654a5d86, []int{2}
}
func (m *Package) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Package.Unmarshal(m, b)
}
func (m *Package) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Package.Marshal(b, m, deterministic)
}
func (dst *Package) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Package.Merge(dst, src)
}
func (m *Package) XXX_Size() int {
	return xxx_messageInfo_Package.Size(m)
}
func (m *Package) XXX_DiscardUnknown() {
	xxx_messageInfo_Package.DiscardUnknown(m)
}

var xxx_messageInfo_Package proto.InternalMessageInfo

func (m *Package) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Package) GetDistribution() []*Distribution {
	if m != nil {
		return m.Distribution
	}
	return nil
}

// Details of a package occurrence.
type Details struct {
	// Where the package was installed.
	Installation         *Installation `protobuf:"bytes,1,opt,name=installation,proto3" json:"installation,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Details) Reset()         { *m = Details{} }
func (m *Details) String() string { return proto.CompactTextString(m) }
func (*Details) ProtoMessage()    {}
func (*Details) Descriptor() ([]byte, []int) {
	return fileDescriptor_package_98e063c9654a5d86, []int{3}
}
func (m *Details) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Details.Unmarshal(m, b)
}
func (m *Details) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Details.Marshal(b, m, deterministic)
}
func (dst *Details) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Details.Merge(dst, src)
}
func (m *Details) XXX_Size() int {
	return xxx_messageInfo_Details.Size(m)
}
func (m *Details) XXX_DiscardUnknown() {
	xxx_messageInfo_Details.DiscardUnknown(m)
}

var xxx_messageInfo_Details proto.InternalMessageInfo

func (m *Details) GetInstallation() *Installation {
	if m != nil {
		return m.Installation
	}
	return nil
}

// This represents how a particular software package may be installed on a
// system.
type Installation struct {
	// Output only. The name of the installed package.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// All of the places within the filesystem versions of this package
	// have been found.
	Location             []*Location `protobuf:"bytes,2,rep,name=location,proto3" json:"location,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Installation) Reset()         { *m = Installation{} }
func (m *Installation) String() string { return proto.CompactTextString(m) }
func (*Installation) ProtoMessage()    {}
func (*Installation) Descriptor() ([]byte, []int) {
	return fileDescriptor_package_98e063c9654a5d86, []int{4}
}
func (m *Installation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Installation.Unmarshal(m, b)
}
func (m *Installation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Installation.Marshal(b, m, deterministic)
}
func (dst *Installation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Installation.Merge(dst, src)
}
func (m *Installation) XXX_Size() int {
	return xxx_messageInfo_Installation.Size(m)
}
func (m *Installation) XXX_DiscardUnknown() {
	xxx_messageInfo_Installation.DiscardUnknown(m)
}

var xxx_messageInfo_Installation proto.InternalMessageInfo

func (m *Installation) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Installation) GetLocation() []*Location {
	if m != nil {
		return m.Location
	}
	return nil
}

// Version contains structured information about the version of a package.
type Version struct {
	// Used to correct mistakes in the version numbering scheme.
	Epoch int32 `protobuf:"varint,1,opt,name=epoch,proto3" json:"epoch,omitempty"`
	// The main part of the version name.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// The iteration of the package build from the above version.
	Revision string `protobuf:"bytes,3,opt,name=revision,proto3" json:"revision,omitempty"`
	// Distinguish between sentinel MIN/MAX versions and normal versions. If
	// kind is not NORMAL, then the other fields are ignored.
	Kind                 Version_VersionKind `protobuf:"varint,4,opt,name=kind,proto3,enum=grafeas.v1beta1.package.Version_VersionKind" json:"kind,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *Version) Reset()         { *m = Version{} }
func (m *Version) String() string { return proto.CompactTextString(m) }
func (*Version) ProtoMessage()    {}
func (*Version) Descriptor() ([]byte, []int) {
	return fileDescriptor_package_98e063c9654a5d86, []int{5}
}
func (m *Version) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Version.Unmarshal(m, b)
}
func (m *Version) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Version.Marshal(b, m, deterministic)
}
func (dst *Version) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Version.Merge(dst, src)
}
func (m *Version) XXX_Size() int {
	return xxx_messageInfo_Version.Size(m)
}
func (m *Version) XXX_DiscardUnknown() {
	xxx_messageInfo_Version.DiscardUnknown(m)
}

var xxx_messageInfo_Version proto.InternalMessageInfo

func (m *Version) GetEpoch() int32 {
	if m != nil {
		return m.Epoch
	}
	return 0
}

func (m *Version) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Version) GetRevision() string {
	if m != nil {
		return m.Revision
	}
	return ""
}

func (m *Version) GetKind() Version_VersionKind {
	if m != nil {
		return m.Kind
	}
	return Version_VERSION_KIND_UNSPECIFIED
}

func init() {
	proto.RegisterType((*Distribution)(nil), "grafeas.v1beta1.package.Distribution")
	proto.RegisterType((*Location)(nil), "grafeas.v1beta1.package.Location")
	proto.RegisterType((*Package)(nil), "grafeas.v1beta1.package.Package")
	proto.RegisterType((*Details)(nil), "grafeas.v1beta1.package.Details")
	proto.RegisterType((*Installation)(nil), "grafeas.v1beta1.package.Installation")
	proto.RegisterType((*Version)(nil), "grafeas.v1beta1.package.Version")
	proto.RegisterEnum("grafeas.v1beta1.package.Architecture", Architecture_name, Architecture_value)
	proto.RegisterEnum("grafeas.v1beta1.package.Version_VersionKind", Version_VersionKind_name, Version_VersionKind_value)
}

func init() {
	proto.RegisterFile("google/devtools/containeranalysis/v1beta1/package/package.proto", fileDescriptor_package_98e063c9654a5d86)
}

var fileDescriptor_package_98e063c9654a5d86 = []byte{
	// 575 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0xff, 0x6a, 0xd4, 0x40,
	0x10, 0x36, 0x49, 0x7b, 0x69, 0xe7, 0xce, 0x12, 0x16, 0xb1, 0x41, 0x44, 0x62, 0x40, 0x38, 0x44,
	0x12, 0x5a, 0xa5, 0x88, 0xe2, 0x8f, 0xb3, 0x77, 0xd6, 0xd0, 0xde, 0xb5, 0x6e, 0x7b, 0xa5, 0xf8,
	0xcf, 0xb1, 0xcd, 0xad, 0xb9, 0xa5, 0x69, 0x36, 0xec, 0xee, 0x1d, 0xe8, 0x4b, 0xf8, 0x0e, 0x3e,
	0x9b, 0x0f, 0x22, 0xd9, 0x24, 0x47, 0xaa, 0xb6, 0xea, 0x5f, 0x3b, 0xb3, 0x33, 0xdf, 0x37, 0xfb,
	0xcd, 0x24, 0x03, 0x6f, 0x12, 0xce, 0x93, 0x94, 0x86, 0x53, 0xba, 0x50, 0x9c, 0xa7, 0x32, 0x8c,
	0x79, 0xa6, 0x08, 0xcb, 0xa8, 0x20, 0x19, 0x49, 0xbf, 0x48, 0x26, 0xc3, 0xc5, 0xd6, 0x39, 0x55,
	0x64, 0x2b, 0xcc, 0x49, 0x7c, 0x41, 0x12, 0x5a, 0x9f, 0x41, 0x2e, 0xb8, 0xe2, 0x68, 0x33, 0x11,
	0xe4, 0x33, 0x25, 0x32, 0xa8, 0xd2, 0x82, 0x2a, 0xec, 0x7f, 0x33, 0xa1, 0xd3, 0x67, 0x52, 0x09,
	0x76, 0x3e, 0x57, 0x8c, 0x67, 0x68, 0x13, 0xec, 0x38, 0xa7, 0x93, 0xb9, 0x60, 0xae, 0xe1, 0x19,
	0xdd, 0x75, 0xdc, 0x8a, 0x73, 0x3a, 0x16, 0x0c, 0x45, 0xd0, 0x21, 0x22, 0x9e, 0x31, 0x45, 0x63,
	0x35, 0x17, 0xd4, 0x35, 0x3d, 0xa3, 0xbb, 0xb1, 0xfd, 0x28, 0xb8, 0x86, 0x39, 0xe8, 0x35, 0x92,
	0xf1, 0x15, 0x28, 0xda, 0x83, 0x8d, 0x94, 0x28, 0x2a, 0xd5, 0x64, 0x41, 0x85, 0x64, 0x3c, 0x73,
	0x2d, 0xcf, 0xe8, 0xb6, 0xb7, 0xbd, 0x6b, 0xc9, 0x4e, 0xcb, 0x3c, 0x7c, 0xbb, 0xc4, 0x55, 0x2e,
	0x7a, 0x00, 0x70, 0x49, 0x58, 0xd5, 0x0a, 0x77, 0x45, 0xbf, 0xb7, 0x71, 0x83, 0x1c, 0xb0, 0xe6,
	0x22, 0x75, 0x57, 0x75, 0xa0, 0x30, 0x91, 0x07, 0xed, 0x29, 0x95, 0xb1, 0x60, 0x79, 0xa1, 0xd6,
	0x6d, 0xe9, 0x48, 0xf3, 0xca, 0x97, 0xb0, 0x76, 0xc0, 0x63, 0x72, 0x73, 0x33, 0x5e, 0x80, 0x5d,
	0x3f, 0xdd, 0xfc, 0xc7, 0xa7, 0xd7, 0x00, 0x84, 0x60, 0x25, 0x27, 0x6a, 0xa6, 0x35, 0xaf, 0x63,
	0x6d, 0xfb, 0x33, 0xb0, 0x8f, 0xca, 0xfc, 0x22, 0x9c, 0x91, 0x4b, 0x5a, 0x15, 0xd4, 0x76, 0xd1,
	0xfb, 0x69, 0x63, 0x48, 0x2e, 0x78, 0x56, 0xb7, 0x7d, 0x43, 0xef, 0x9b, 0x13, 0xc5, 0x57, 0xa0,
	0xfe, 0x09, 0xd8, 0x7d, 0xaa, 0x08, 0x4b, 0x65, 0xc1, 0xca, 0x32, 0xa9, 0x48, 0x9a, 0x6a, 0xb5,
	0xba, 0xe2, 0x4d, 0xac, 0x51, 0x23, 0x19, 0x5f, 0x81, 0xfa, 0x04, 0x3a, 0xcd, 0xe8, 0x1f, 0x45,
	0xbc, 0x82, 0xb5, 0xb4, 0x6a, 0xac, 0x6b, 0x6a, 0x01, 0x0f, 0xaf, 0x2d, 0x55, 0x4f, 0x00, 0x2f,
	0x21, 0xfe, 0x0f, 0x03, 0xec, 0x7a, 0xee, 0x77, 0x60, 0x95, 0xe6, 0x3c, 0x9e, 0x69, 0xfe, 0x55,
	0x5c, 0x3a, 0xcb, 0xa2, 0x66, 0xa3, 0xe8, 0x3d, 0x58, 0x13, 0x74, 0xc1, 0x96, 0x1f, 0xd9, 0x3a,
	0x5e, 0xfa, 0xe8, 0x2d, 0xac, 0x5c, 0xb0, 0x6c, 0xaa, 0xbf, 0x9b, 0x8d, 0xed, 0x27, 0x7f, 0x9b,
	0x60, 0x7d, 0xee, 0xb3, 0x6c, 0x8a, 0x35, 0xd2, 0xff, 0x08, 0xed, 0xc6, 0x25, 0xba, 0x0f, 0xee,
	0xe9, 0x00, 0x1f, 0x47, 0x87, 0xa3, 0xc9, 0x7e, 0x34, 0xea, 0x4f, 0xc6, 0xa3, 0xe3, 0xa3, 0xc1,
	0x6e, 0xf4, 0x3e, 0x1a, 0xf4, 0x9d, 0x5b, 0x08, 0xa0, 0x35, 0x3a, 0xc4, 0xc3, 0xde, 0x81, 0x63,
	0xa0, 0x36, 0xd8, 0xc3, 0x68, 0x14, 0x0d, 0xc7, 0x43, 0xc7, 0xd4, 0x4e, 0xef, 0x4c, 0x3b, 0xd6,
	0xe3, 0xd7, 0xd0, 0x69, 0xfe, 0x39, 0x05, 0x67, 0x0f, 0xef, 0x7e, 0x88, 0x4e, 0x06, 0xbb, 0x27,
	0x63, 0x3c, 0xf8, 0x85, 0xd3, 0x06, 0xeb, 0xec, 0xf9, 0x8e, 0x63, 0x68, 0x63, 0xe7, 0x99, 0x63,
	0xbe, 0xfb, 0x0a, 0x77, 0x19, 0xff, 0x5d, 0xca, 0x45, 0x72, 0x64, 0x7c, 0x3a, 0x2b, 0xd7, 0x48,
	0x90, 0xf0, 0x94, 0x64, 0x49, 0xc0, 0x45, 0x12, 0x26, 0x34, 0xd3, 0x1b, 0x22, 0x2c, 0x43, 0x24,
	0x67, 0xf2, 0x3f, 0xb6, 0xcc, 0xcb, 0xea, 0xfc, 0x6e, 0x5a, 0x7b, 0xb8, 0x77, 0xde, 0xd2, 0x54,
	0x4f, 0x7f, 0x06, 0x00, 0x00, 0xff, 0xff, 0x78, 0xac, 0x5c, 0xdf, 0xaf, 0x04, 0x00, 0x00,
}
