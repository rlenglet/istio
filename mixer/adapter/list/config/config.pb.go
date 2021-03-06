// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mixer/adapter/list/config/config.proto

// The `list` adapter makes it possible to perform simple whitelist or blacklist
// checks. You can configure the adapter with the list to check, or you can point
// it to a URL from where the list should be fetched. Lists can be simple strings,
// IP addresses, or regex patterns.
//
// This adapter supports the [listentry template](https://istio.io/docs/reference/config/policy-and-telemetry/templates/listentry/).

package config

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/types"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strconv "strconv"
	strings "strings"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// Determines the type of list that the adapter is consulting.
type Params_ListEntryType int32

const (
	// List entries are treated as plain strings.
	STRINGS Params_ListEntryType = 0
	// List entries are treated as case-insensitive strings.
	CASE_INSENSITIVE_STRINGS Params_ListEntryType = 1
	// List entries are treated as IP addresses and ranges.
	IP_ADDRESSES Params_ListEntryType = 2
	// List entries are treated as re2 regexp. See [here](https://github.com/google/re2/wiki/Syntax) for the supported syntax.
	REGEX Params_ListEntryType = 3
)

var Params_ListEntryType_name = map[int32]string{
	0: "STRINGS",
	1: "CASE_INSENSITIVE_STRINGS",
	2: "IP_ADDRESSES",
	3: "REGEX",
}

var Params_ListEntryType_value = map[string]int32{
	"STRINGS":                  0,
	"CASE_INSENSITIVE_STRINGS": 1,
	"IP_ADDRESSES":             2,
	"REGEX":                    3,
}

func (Params_ListEntryType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_46cb3210745c49e8, []int{0, 0}
}

// Configuration format for the `list` adapter.
type Params struct {
	// Where to find the list to check against. This may be omitted for a completely local list.
	ProviderUrl string `protobuf:"bytes,1,opt,name=provider_url,json=providerUrl,proto3" json:"provider_url,omitempty"`
	// Determines how often the provider is polled for
	// an updated list
	RefreshInterval time.Duration `protobuf:"bytes,2,opt,name=refresh_interval,json=refreshInterval,proto3,stdduration" json:"refresh_interval"`
	// Indicates how long to keep a list before discarding it.
	// Typically, the TTL value should be set to noticeably longer (> 2x) than the
	// refresh interval to ensure continued operation in the face of transient
	// server outages.
	Ttl time.Duration `protobuf:"bytes,3,opt,name=ttl,proto3,stdduration" json:"ttl"`
	// Indicates the amount of time a caller of this adapter can cache an answer
	// before it should ask the adapter again.
	CachingInterval time.Duration `protobuf:"bytes,4,opt,name=caching_interval,json=cachingInterval,proto3,stdduration" json:"caching_interval"`
	// Indicates the number of times a caller of this adapter can use a cached answer
	// before it should ask the adapter again.
	CachingUseCount int32 `protobuf:"varint,5,opt,name=caching_use_count,json=cachingUseCount,proto3" json:"caching_use_count,omitempty"`
	// List entries that are consulted first, before the list from the server
	Overrides []string `protobuf:"bytes,6,rep,name=overrides,proto3" json:"overrides,omitempty"`
	// Determines the kind of list entry and overrides.
	EntryType Params_ListEntryType `protobuf:"varint,7,opt,name=entry_type,json=entryType,proto3,enum=adapter.list.config.Params_ListEntryType" json:"entry_type,omitempty"`
	// Whether the list operates as a blacklist or a whitelist.
	Blacklist bool `protobuf:"varint,8,opt,name=blacklist,proto3" json:"blacklist,omitempty"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_46cb3210745c49e8, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("adapter.list.config.Params_ListEntryType", Params_ListEntryType_name, Params_ListEntryType_value)
	proto.RegisterType((*Params)(nil), "adapter.list.config.Params")
}

func init() {
	proto.RegisterFile("mixer/adapter/list/config/config.proto", fileDescriptor_46cb3210745c49e8)
}

var fileDescriptor_46cb3210745c49e8 = []byte{
	// 473 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0xb1, 0x6e, 0xd3, 0x40,
	0x1c, 0xc6, 0xef, 0x9a, 0x26, 0x8d, 0x2f, 0x05, 0xcc, 0xc1, 0x60, 0xaa, 0xea, 0x6a, 0x3a, 0x20,
	0xc3, 0x60, 0x4b, 0x45, 0xec, 0xb4, 0x8d, 0x55, 0x2c, 0xa1, 0xa8, 0xb2, 0x53, 0x40, 0x2c, 0x96,
	0xe3, 0x5c, 0xdc, 0x13, 0xae, 0xcf, 0x3a, 0x9f, 0x23, 0xb2, 0xf1, 0x02, 0x48, 0x8c, 0x3c, 0x02,
	0x8f, 0x92, 0x31, 0x63, 0x27, 0x20, 0xce, 0xc2, 0xd8, 0x47, 0x40, 0x8e, 0xed, 0x46, 0x48, 0x0c,
	0x30, 0xf9, 0xaf, 0xef, 0xbe, 0xdf, 0x7d, 0x9f, 0xff, 0x87, 0x9e, 0x5c, 0xb1, 0x8f, 0x54, 0x58,
	0xc1, 0x38, 0x48, 0x25, 0x15, 0x56, 0xcc, 0x32, 0x69, 0x85, 0x3c, 0x99, 0xb0, 0xa8, 0xfe, 0x98,
	0xa9, 0xe0, 0x92, 0xe3, 0x07, 0xb5, 0xc3, 0x2c, 0x1d, 0x66, 0x75, 0xb4, 0x47, 0x22, 0xce, 0xa3,
	0x98, 0x5a, 0x6b, 0xcb, 0x28, 0x9f, 0x58, 0xe3, 0x5c, 0x04, 0x92, 0xf1, 0xa4, 0x82, 0xf6, 0x1e,
	0x46, 0x3c, 0xe2, 0xeb, 0xd1, 0x2a, 0xa7, 0x4a, 0x3d, 0xfc, 0xbc, 0x8d, 0x3a, 0xe7, 0x81, 0x08,
	0xae, 0x32, 0xfc, 0x18, 0xed, 0xa6, 0x82, 0x4f, 0xd9, 0x98, 0x0a, 0x3f, 0x17, 0xb1, 0x06, 0x75,
	0x68, 0x28, 0x6e, 0xaf, 0xd1, 0x2e, 0x44, 0x8c, 0x07, 0x48, 0x15, 0x74, 0x22, 0x68, 0x76, 0xe9,
	0xb3, 0x44, 0x52, 0x31, 0x0d, 0x62, 0x6d, 0x4b, 0x87, 0x46, 0xef, 0xe8, 0x91, 0x59, 0xc5, 0x9b,
	0x4d, 0xbc, 0xd9, 0xaf, 0xe3, 0x4f, 0xba, 0xf3, 0xef, 0x07, 0xe0, 0xeb, 0x8f, 0x03, 0xe8, 0xde,
	0xab, 0x61, 0xa7, 0x66, 0xf1, 0x0b, 0xd4, 0x92, 0x32, 0xd6, 0x5a, 0xff, 0x7e, 0x45, 0xe9, 0x2f,
	0x6b, 0x84, 0x41, 0x78, 0xc9, 0x92, 0x68, 0x53, 0x63, 0xfb, 0x3f, 0x6a, 0xd4, 0xf0, 0x6d, 0x8d,
	0x67, 0xe8, 0x7e, 0x73, 0x5f, 0x9e, 0x51, 0x3f, 0xe4, 0x79, 0x22, 0xb5, 0xb6, 0x0e, 0x8d, 0xf6,
	0xad, 0xf7, 0x22, 0xa3, 0xa7, 0xa5, 0x8c, 0xf7, 0x91, 0xc2, 0xa7, 0x54, 0x08, 0x36, 0xa6, 0x99,
	0xd6, 0xd1, 0x5b, 0x86, 0xe2, 0x6e, 0x04, 0xfc, 0x0a, 0x21, 0x9a, 0x48, 0x31, 0xf3, 0xe5, 0x2c,
	0xa5, 0xda, 0x8e, 0x0e, 0x8d, 0xbb, 0x47, 0x4f, 0xcd, 0xbf, 0x3c, 0x97, 0x59, 0x2d, 0xdd, 0x7c,
	0xcd, 0x32, 0x69, 0x97, 0xc4, 0x70, 0x96, 0x52, 0x57, 0xa1, 0xcd, 0x58, 0xe6, 0x8c, 0xe2, 0x20,
	0xfc, 0x50, 0x32, 0x5a, 0x57, 0x87, 0x46, 0xd7, 0xdd, 0x08, 0x87, 0x6f, 0xd1, 0x9d, 0x3f, 0x48,
	0xdc, 0x43, 0x3b, 0xde, 0xd0, 0x75, 0x06, 0x67, 0x9e, 0x0a, 0xf0, 0x3e, 0xd2, 0x4e, 0x8f, 0x3d,
	0xdb, 0x77, 0x06, 0x9e, 0x3d, 0xf0, 0x9c, 0xa1, 0xf3, 0xc6, 0xf6, 0x9b, 0x53, 0x88, 0x55, 0xb4,
	0xeb, 0x9c, 0xfb, 0xc7, 0xfd, 0xbe, 0x6b, 0x7b, 0x9e, 0xed, 0xa9, 0x5b, 0x58, 0x41, 0x6d, 0xd7,
	0x3e, 0xb3, 0xdf, 0xa9, 0xad, 0x93, 0x97, 0xf3, 0x25, 0x01, 0x8b, 0x25, 0x01, 0xd7, 0x4b, 0x02,
	0x6e, 0x96, 0x04, 0x7c, 0x2a, 0x08, 0xfc, 0x56, 0x10, 0x30, 0x2f, 0x08, 0x5c, 0x14, 0x04, 0xfe,
	0x2c, 0x08, 0xfc, 0x55, 0x10, 0x70, 0x53, 0x10, 0xf8, 0x65, 0x45, 0xc0, 0x62, 0x45, 0xc0, 0xf5,
	0x8a, 0x80, 0xf7, 0x9d, 0xea, 0xc7, 0x46, 0x9d, 0xf5, 0xea, 0x9f, 0xff, 0x0e, 0x00, 0x00, 0xff,
	0xff, 0x1f, 0xed, 0x14, 0xa2, 0xcd, 0x02, 0x00, 0x00,
}

func (x Params_ListEntryType) String() string {
	s, ok := Params_ListEntryType_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Blacklist {
		i--
		if m.Blacklist {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x40
	}
	if m.EntryType != 0 {
		i = encodeVarintConfig(dAtA, i, uint64(m.EntryType))
		i--
		dAtA[i] = 0x38
	}
	if len(m.Overrides) > 0 {
		for iNdEx := len(m.Overrides) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Overrides[iNdEx])
			copy(dAtA[i:], m.Overrides[iNdEx])
			i = encodeVarintConfig(dAtA, i, uint64(len(m.Overrides[iNdEx])))
			i--
			dAtA[i] = 0x32
		}
	}
	if m.CachingUseCount != 0 {
		i = encodeVarintConfig(dAtA, i, uint64(m.CachingUseCount))
		i--
		dAtA[i] = 0x28
	}
	n1, err1 := github_com_gogo_protobuf_types.StdDurationMarshalTo(m.CachingInterval, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdDuration(m.CachingInterval):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintConfig(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x22
	n2, err2 := github_com_gogo_protobuf_types.StdDurationMarshalTo(m.Ttl, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdDuration(m.Ttl):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintConfig(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x1a
	n3, err3 := github_com_gogo_protobuf_types.StdDurationMarshalTo(m.RefreshInterval, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdDuration(m.RefreshInterval):])
	if err3 != nil {
		return 0, err3
	}
	i -= n3
	i = encodeVarintConfig(dAtA, i, uint64(n3))
	i--
	dAtA[i] = 0x12
	if len(m.ProviderUrl) > 0 {
		i -= len(m.ProviderUrl)
		copy(dAtA[i:], m.ProviderUrl)
		i = encodeVarintConfig(dAtA, i, uint64(len(m.ProviderUrl)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintConfig(dAtA []byte, offset int, v uint64) int {
	offset -= sovConfig(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ProviderUrl)
	if l > 0 {
		n += 1 + l + sovConfig(uint64(l))
	}
	l = github_com_gogo_protobuf_types.SizeOfStdDuration(m.RefreshInterval)
	n += 1 + l + sovConfig(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdDuration(m.Ttl)
	n += 1 + l + sovConfig(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdDuration(m.CachingInterval)
	n += 1 + l + sovConfig(uint64(l))
	if m.CachingUseCount != 0 {
		n += 1 + sovConfig(uint64(m.CachingUseCount))
	}
	if len(m.Overrides) > 0 {
		for _, s := range m.Overrides {
			l = len(s)
			n += 1 + l + sovConfig(uint64(l))
		}
	}
	if m.EntryType != 0 {
		n += 1 + sovConfig(uint64(m.EntryType))
	}
	if m.Blacklist {
		n += 2
	}
	return n
}

func sovConfig(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozConfig(x uint64) (n int) {
	return sovConfig(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Params) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Params{`,
		`ProviderUrl:` + fmt.Sprintf("%v", this.ProviderUrl) + `,`,
		`RefreshInterval:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.RefreshInterval), "Duration", "types.Duration", 1), `&`, ``, 1) + `,`,
		`Ttl:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.Ttl), "Duration", "types.Duration", 1), `&`, ``, 1) + `,`,
		`CachingInterval:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.CachingInterval), "Duration", "types.Duration", 1), `&`, ``, 1) + `,`,
		`CachingUseCount:` + fmt.Sprintf("%v", this.CachingUseCount) + `,`,
		`Overrides:` + fmt.Sprintf("%v", this.Overrides) + `,`,
		`EntryType:` + fmt.Sprintf("%v", this.EntryType) + `,`,
		`Blacklist:` + fmt.Sprintf("%v", this.Blacklist) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringConfig(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowConfig
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProviderUrl", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ProviderUrl = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RefreshInterval", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthConfig
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(&m.RefreshInterval, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ttl", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthConfig
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(&m.Ttl, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CachingInterval", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthConfig
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(&m.CachingInterval, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CachingUseCount", wireType)
			}
			m.CachingUseCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CachingUseCount |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Overrides", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Overrides = append(m.Overrides, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EntryType", wireType)
			}
			m.EntryType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EntryType |= Params_ListEntryType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Blacklist", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Blacklist = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipConfig(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthConfig
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthConfig
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
func skipConfig(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowConfig
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
					return 0, ErrIntOverflowConfig
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowConfig
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
				return 0, ErrInvalidLengthConfig
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthConfig
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowConfig
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipConfig(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthConfig
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthConfig = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowConfig   = fmt.Errorf("proto: integer overflow")
)
