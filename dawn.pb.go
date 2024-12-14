package main

import (
	fmt "fmt"

	io "io"
	math "math"
	math_bits "math/bits"

	"github.com/gogo/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type Dawn struct {
	Day uint64 `protobuf:"varint,1,opt,name=day,proto3" json:"day,omitempty"`
}

func (m *Dawn) Reset()         { *m = Dawn{} }
func (m *Dawn) String() string { return proto.CompactTextString(m) }
func (*Dawn) ProtoMessage()    {}
func (*Dawn) Descriptor() ([]byte, []int) {
	return fileDescriptor_85df52da6d694a76, []int{0}
}
func (m *Dawn) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Dawn) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Dawn.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Dawn) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Dawn.Merge(m, src)
}
func (m *Dawn) XXX_Size() int {
	return m.Size()
}
func (m *Dawn) XXX_DiscardUnknown() {
	xxx_messageInfo_Dawn.DiscardUnknown(m)
}

var xxx_messageInfo_Dawn proto.InternalMessageInfo

func (m *Dawn) GetDay() uint64 {
	if m != nil {
		return m.Day
	}
	return 0
}

func init() {
	proto.RegisterType((*Dawn)(nil), "stchain.rollapp.checkin.Dawn")
}

func init() { proto.RegisterFile("rollapp/checkin/dawn.proto", fileDescriptor_85df52da6d694a76) }

var fileDescriptor_85df52da6d694a76 = []byte{
	// 153 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2a, 0xca, 0xcf, 0xc9,
	0x49, 0x2c, 0x28, 0xd0, 0x4f, 0xce, 0x48, 0x4d, 0xce, 0xce, 0xcc, 0xd3, 0x4f, 0x49, 0x2c, 0xcf,
	0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x2f, 0x2e, 0x49, 0xce, 0x48, 0xcc, 0xcc, 0xd3,
	0x83, 0xaa, 0xd1, 0x83, 0xaa, 0x51, 0x92, 0xe0, 0x62, 0x71, 0x49, 0x2c, 0xcf, 0x13, 0x12, 0xe0,
	0x62, 0x4e, 0x49, 0xac, 0x94, 0x60, 0x54, 0x60, 0xd4, 0x60, 0x09, 0x02, 0x31, 0x9d, 0x5c, 0x4f,
	0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23, 0x39, 0xc6, 0x09, 0x8f, 0xe5, 0x18,
	0x2e, 0x3c, 0x96, 0x63, 0xb8, 0xf1, 0x58, 0x8e, 0x21, 0x4a, 0x3b, 0x3d, 0xb3, 0x24, 0xa3, 0x34,
	0x49, 0x2f, 0x39, 0x3f, 0x57, 0xbf, 0xb8, 0x44, 0x17, 0x6c, 0xb0, 0x3e, 0xcc, 0xf2, 0x0a, 0xb8,
	0xf5, 0x25, 0x95, 0x05, 0xa9, 0xc5, 0x49, 0x6c, 0x60, 0x07, 0x18, 0x03, 0x02, 0x00, 0x00, 0xff,
	0xff, 0x72, 0x31, 0x81, 0x6a, 0x9e, 0x00, 0x00, 0x00,
}

func (m *Dawn) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Dawn) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Dawn) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Day != 0 {
		i = encodeVarintDawn(dAtA, i, uint64(m.Day))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintDawn(dAtA []byte, offset int, v uint64) int {
	offset -= sovDawn(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Dawn) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Day != 0 {
		n += 1 + sovDawn(uint64(m.Day))
	}
	return n
}

func sovDawn(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozDawn(x uint64) (n int) {
	return sovDawn(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Dawn) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDawn
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Dawn: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Dawn: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Day", wireType)
			}
			m.Day = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDawn
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Day |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipDawn(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDawn
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipDawn(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDawn
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowDawn
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowDawn
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthDawn
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupDawn
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthDawn
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthDawn        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDawn          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupDawn = fmt.Errorf("proto: unexpected end of group")
)
