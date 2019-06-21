// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v1/enums/conversion_action_category.proto

package enums

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// The category of conversions that are associated with a ConversionAction.
type ConversionActionCategoryEnum_ConversionActionCategory int32

const (
	// Not specified.
	ConversionActionCategoryEnum_UNSPECIFIED ConversionActionCategoryEnum_ConversionActionCategory = 0
	// Used for return value only. Represents value unknown in this version.
	ConversionActionCategoryEnum_UNKNOWN ConversionActionCategoryEnum_ConversionActionCategory = 1
	// Default category.
	ConversionActionCategoryEnum_DEFAULT ConversionActionCategoryEnum_ConversionActionCategory = 2
	// User visiting a page.
	ConversionActionCategoryEnum_PAGE_VIEW ConversionActionCategoryEnum_ConversionActionCategory = 3
	// Purchase, sales, or "order placed" event.
	ConversionActionCategoryEnum_PURCHASE ConversionActionCategoryEnum_ConversionActionCategory = 4
	// Signup user action.
	ConversionActionCategoryEnum_SIGNUP ConversionActionCategoryEnum_ConversionActionCategory = 5
	// Lead-generating action.
	ConversionActionCategoryEnum_LEAD ConversionActionCategoryEnum_ConversionActionCategory = 6
	// Software download action (as for an app).
	ConversionActionCategoryEnum_DOWNLOAD ConversionActionCategoryEnum_ConversionActionCategory = 7
)

var ConversionActionCategoryEnum_ConversionActionCategory_name = map[int32]string{
	0: "UNSPECIFIED",
	1: "UNKNOWN",
	2: "DEFAULT",
	3: "PAGE_VIEW",
	4: "PURCHASE",
	5: "SIGNUP",
	6: "LEAD",
	7: "DOWNLOAD",
}

var ConversionActionCategoryEnum_ConversionActionCategory_value = map[string]int32{
	"UNSPECIFIED": 0,
	"UNKNOWN":     1,
	"DEFAULT":     2,
	"PAGE_VIEW":   3,
	"PURCHASE":    4,
	"SIGNUP":      5,
	"LEAD":        6,
	"DOWNLOAD":    7,
}

func (x ConversionActionCategoryEnum_ConversionActionCategory) String() string {
	return proto.EnumName(ConversionActionCategoryEnum_ConversionActionCategory_name, int32(x))
}

func (ConversionActionCategoryEnum_ConversionActionCategory) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_23ef90dd965d8778, []int{0, 0}
}

// Container for enum describing the category of conversions that are associated
// with a ConversionAction.
type ConversionActionCategoryEnum struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConversionActionCategoryEnum) Reset()         { *m = ConversionActionCategoryEnum{} }
func (m *ConversionActionCategoryEnum) String() string { return proto.CompactTextString(m) }
func (*ConversionActionCategoryEnum) ProtoMessage()    {}
func (*ConversionActionCategoryEnum) Descriptor() ([]byte, []int) {
	return fileDescriptor_23ef90dd965d8778, []int{0}
}

func (m *ConversionActionCategoryEnum) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConversionActionCategoryEnum.Unmarshal(m, b)
}
func (m *ConversionActionCategoryEnum) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConversionActionCategoryEnum.Marshal(b, m, deterministic)
}
func (m *ConversionActionCategoryEnum) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConversionActionCategoryEnum.Merge(m, src)
}
func (m *ConversionActionCategoryEnum) XXX_Size() int {
	return xxx_messageInfo_ConversionActionCategoryEnum.Size(m)
}
func (m *ConversionActionCategoryEnum) XXX_DiscardUnknown() {
	xxx_messageInfo_ConversionActionCategoryEnum.DiscardUnknown(m)
}

var xxx_messageInfo_ConversionActionCategoryEnum proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("google.ads.googleads.v1.enums.ConversionActionCategoryEnum_ConversionActionCategory", ConversionActionCategoryEnum_ConversionActionCategory_name, ConversionActionCategoryEnum_ConversionActionCategory_value)
	proto.RegisterType((*ConversionActionCategoryEnum)(nil), "google.ads.googleads.v1.enums.ConversionActionCategoryEnum")
}

func init() {
	proto.RegisterFile("google/ads/googleads/v1/enums/conversion_action_category.proto", fileDescriptor_23ef90dd965d8778)
}

var fileDescriptor_23ef90dd965d8778 = []byte{
	// 360 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x51, 0xc1, 0x6a, 0xea, 0x40,
	0x14, 0x7d, 0x89, 0x3e, 0xf5, 0x8d, 0xef, 0xf1, 0x86, 0xac, 0x4a, 0xd1, 0x85, 0x7e, 0xc0, 0x84,
	0xd0, 0xdd, 0x14, 0x0a, 0x63, 0x32, 0xa6, 0xa1, 0x12, 0x43, 0x6d, 0x22, 0x94, 0x80, 0xa4, 0x49,
	0x08, 0x82, 0xce, 0x48, 0x26, 0x0a, 0xfd, 0x81, 0x7e, 0x48, 0x37, 0x85, 0x7e, 0x4a, 0x3f, 0xa5,
	0xcb, 0x7e, 0x41, 0x99, 0x89, 0x66, 0x67, 0x37, 0xc9, 0x99, 0xb9, 0xe7, 0x9c, 0x7b, 0xef, 0x19,
	0x70, 0x53, 0x70, 0x5e, 0x6c, 0x72, 0x33, 0xc9, 0x84, 0x59, 0x43, 0x89, 0x0e, 0x96, 0x99, 0xb3,
	0xfd, 0x56, 0x98, 0x29, 0x67, 0x87, 0xbc, 0x14, 0x6b, 0xce, 0x56, 0x49, 0x5a, 0xc9, 0x5f, 0x9a,
	0x54, 0x79, 0xc1, 0xcb, 0x67, 0xb4, 0x2b, 0x79, 0xc5, 0x8d, 0x61, 0x2d, 0x42, 0x49, 0x26, 0x50,
	0xa3, 0x47, 0x07, 0x0b, 0x29, 0xfd, 0xe5, 0xe0, 0x64, 0xbf, 0x5b, 0x9b, 0x09, 0x63, 0xbc, 0x4a,
	0xa4, 0x89, 0xa8, 0xc5, 0xe3, 0x37, 0x0d, 0x0c, 0xec, 0xa6, 0x03, 0x51, 0x0d, 0xec, 0xa3, 0x3f,
	0x65, 0xfb, 0xed, 0xf8, 0x45, 0x03, 0x17, 0xe7, 0x08, 0xc6, 0x7f, 0xd0, 0x0f, 0xfd, 0x45, 0x40,
	0x6d, 0x6f, 0xea, 0x51, 0x07, 0xfe, 0x32, 0xfa, 0xa0, 0x1b, 0xfa, 0x77, 0xfe, 0x7c, 0xe9, 0x43,
	0x4d, 0x1e, 0x1c, 0x3a, 0x25, 0xe1, 0xec, 0x01, 0xea, 0xc6, 0x3f, 0xf0, 0x27, 0x20, 0x2e, 0x5d,
	0x45, 0x1e, 0x5d, 0xc2, 0x96, 0xf1, 0x17, 0xf4, 0x82, 0xf0, 0xde, 0xbe, 0x25, 0x0b, 0x0a, 0xdb,
	0x06, 0x00, 0x9d, 0x85, 0xe7, 0xfa, 0x61, 0x00, 0x7f, 0x1b, 0x3d, 0xd0, 0x9e, 0x51, 0xe2, 0xc0,
	0x8e, 0xe4, 0x38, 0xf3, 0xa5, 0x3f, 0x9b, 0x13, 0x07, 0x76, 0x27, 0x5f, 0x1a, 0x18, 0xa5, 0x7c,
	0x8b, 0x7e, 0xdc, 0x76, 0x32, 0x3c, 0x37, 0x6b, 0x20, 0xd7, 0x0d, 0xb4, 0xc7, 0xc9, 0x51, 0x5f,
	0xf0, 0x4d, 0xc2, 0x0a, 0xc4, 0xcb, 0xc2, 0x2c, 0x72, 0xa6, 0xc2, 0x38, 0xa5, 0xbf, 0x5b, 0x8b,
	0x33, 0x8f, 0x71, 0xad, 0xbe, 0xaf, 0x7a, 0xcb, 0x25, 0xe4, 0x5d, 0x1f, 0xba, 0xb5, 0x15, 0xc9,
	0x04, 0xaa, 0xa1, 0x44, 0x91, 0x85, 0x64, 0x70, 0xe2, 0xe3, 0x54, 0x8f, 0x49, 0x26, 0xe2, 0xa6,
	0x1e, 0x47, 0x56, 0xac, 0xea, 0x9f, 0xfa, 0xa8, 0xbe, 0xc4, 0x98, 0x64, 0x02, 0xe3, 0x86, 0x81,
	0x71, 0x64, 0x61, 0xac, 0x38, 0x4f, 0x1d, 0x35, 0xd8, 0xd5, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff,
	0xf5, 0xae, 0x28, 0xe8, 0x24, 0x02, 0x00, 0x00,
}