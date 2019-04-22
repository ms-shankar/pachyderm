// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/pachyderm/pachyderm/src/server/pkg/storage/chunk/chunk.proto

/*
Package chunk is a generated protocol buffer package.

It is generated from these files:
	github.com/pachyderm/pachyderm/src/server/pkg/storage/chunk/chunk.proto

It has these top-level messages:
	DataRef
*/
package chunk

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

// DataRef is a reference to a chunk of data.
type DataRef struct {
	Hash string `protobuf:"bytes,1,opt,name=hash" json:"hash,omitempty"`
	// Hash of the subchunk (defined by offset and size).
	// Should be empty if this is a reference to the whole chunk.
	SubHash string `protobuf:"bytes,2,opt,name=sub_hash,json=subHash" json:"sub_hash,omitempty"`
	Offset  int64  `protobuf:"varint,3,opt,name=offset" json:"offset,omitempty"`
	Size    int64  `protobuf:"varint,4,opt,name=size" json:"size,omitempty"`
}

func (m *DataRef) Reset()                    { *m = DataRef{} }
func (m *DataRef) String() string            { return proto.CompactTextString(m) }
func (*DataRef) ProtoMessage()               {}
func (*DataRef) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *DataRef) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *DataRef) GetSubHash() string {
	if m != nil {
		return m.SubHash
	}
	return ""
}

func (m *DataRef) GetOffset() int64 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *DataRef) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func init() {
	proto.RegisterType((*DataRef)(nil), "chunk.DataRef")
}

func init() {
	proto.RegisterFile("github.com/pachyderm/pachyderm/src/server/pkg/storage/chunk/chunk.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 175 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x72, 0x4f, 0xcf, 0x2c, 0xc9,
	0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x2f, 0x48, 0x4c, 0xce, 0xa8, 0x4c, 0x49, 0x2d, 0x42,
	0x66, 0x15, 0x17, 0x25, 0xeb, 0x17, 0xa7, 0x16, 0x95, 0xa5, 0x16, 0xe9, 0x17, 0x64, 0xa7, 0xeb,
	0x17, 0x97, 0xe4, 0x17, 0x25, 0xa6, 0xa7, 0xea, 0x27, 0x67, 0x94, 0xe6, 0x65, 0x43, 0x48, 0xbd,
	0x82, 0xa2, 0xfc, 0x92, 0x7c, 0x21, 0x56, 0x30, 0x47, 0x29, 0x85, 0x8b, 0xdd, 0x25, 0xb1, 0x24,
	0x31, 0x28, 0x35, 0x4d, 0x48, 0x88, 0x8b, 0x25, 0x23, 0xb1, 0x38, 0x43, 0x82, 0x51, 0x81, 0x51,
	0x83, 0x33, 0x08, 0xcc, 0x16, 0x92, 0xe4, 0xe2, 0x28, 0x2e, 0x4d, 0x8a, 0x07, 0x8b, 0x33, 0x81,
	0xc5, 0xd9, 0x8b, 0x4b, 0x93, 0x3c, 0x40, 0x52, 0x62, 0x5c, 0x6c, 0xf9, 0x69, 0x69, 0xc5, 0xa9,
	0x25, 0x12, 0xcc, 0x0a, 0x8c, 0x1a, 0xcc, 0x41, 0x50, 0x1e, 0xc8, 0x98, 0xe2, 0xcc, 0xaa, 0x54,
	0x09, 0x16, 0xb0, 0x28, 0x98, 0xed, 0x64, 0x1b, 0x65, 0x4d, 0x81, 0xbb, 0x93, 0xd8, 0xc0, 0x4e,
	0x36, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x21, 0x55, 0x62, 0xeb, 0xfd, 0x00, 0x00, 0x00,
}
