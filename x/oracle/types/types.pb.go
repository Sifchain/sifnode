// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sifnode/oracle/v1/types.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
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

// StatusText is an enum used to represent the status of the prophecy
type StatusText int32

const (
	// Default value
	StatusText_STATUS_TEXT_UNSPECIFIED StatusText = 0
	// Pending status
	StatusText_STATUS_TEXT_PENDING StatusText = 1
	// Success status
	StatusText_STATUS_TEXT_SUCCESS StatusText = 2
	// Failed status
	StatusText_STATUS_TEXT_FAILED StatusText = 3
)

var StatusText_name = map[int32]string{
	0: "STATUS_TEXT_UNSPECIFIED",
	1: "STATUS_TEXT_PENDING",
	2: "STATUS_TEXT_SUCCESS",
	3: "STATUS_TEXT_FAILED",
}

var StatusText_value = map[string]int32{
	"STATUS_TEXT_UNSPECIFIED": 0,
	"STATUS_TEXT_PENDING":     1,
	"STATUS_TEXT_SUCCESS":     2,
	"STATUS_TEXT_FAILED":      3,
}

func (x StatusText) String() string {
	return proto.EnumName(StatusText_name, int32(x))
}

func (StatusText) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_dac1b931484f4203, []int{0}
}

type GenesisState struct {
	AddressWhitelist []string      `protobuf:"bytes,1,rep,name=address_whitelist,json=addressWhitelist,proto3" json:"address_whitelist,omitempty"`
	AdminAddress     string        `protobuf:"bytes,2,opt,name=admin_address,json=adminAddress,proto3" json:"admin_address,omitempty"`
	Prophecies       []*DBProphecy `protobuf:"bytes,3,rep,name=prophecies,proto3" json:"prophecies,omitempty"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_dac1b931484f4203, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetAddressWhitelist() []string {
	if m != nil {
		return m.AddressWhitelist
	}
	return nil
}

func (m *GenesisState) GetAdminAddress() string {
	if m != nil {
		return m.AdminAddress
	}
	return ""
}

func (m *GenesisState) GetProphecies() []*DBProphecy {
	if m != nil {
		return m.Prophecies
	}
	return nil
}

// Claim contains an arbitrary claim with arbitrary content made by a given
// validator
type Claim struct {
	Id               string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ValidatorAddress string `protobuf:"bytes,2,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
	Content          string `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
}

func (m *Claim) Reset()         { *m = Claim{} }
func (m *Claim) String() string { return proto.CompactTextString(m) }
func (*Claim) ProtoMessage()    {}
func (*Claim) Descriptor() ([]byte, []int) {
	return fileDescriptor_dac1b931484f4203, []int{1}
}
func (m *Claim) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Claim) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Claim.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Claim) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Claim.Merge(m, src)
}
func (m *Claim) XXX_Size() int {
	return m.Size()
}
func (m *Claim) XXX_DiscardUnknown() {
	xxx_messageInfo_Claim.DiscardUnknown(m)
}

var xxx_messageInfo_Claim proto.InternalMessageInfo

func (m *Claim) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Claim) GetValidatorAddress() string {
	if m != nil {
		return m.ValidatorAddress
	}
	return ""
}

func (m *Claim) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

// DBProphecy is what the prophecy becomes when being saved to the database.
//
//	Tendermint/Amino does not support maps so we must serialize those variables
//	into bytes.
type DBProphecy struct {
	Id              string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Status          Status `protobuf:"bytes,2,opt,name=status,proto3" json:"status"`
	ClaimValidators []byte `protobuf:"bytes,3,opt,name=claim_validators,json=claimValidators,proto3" json:"claim_validators,omitempty"`
	ValidatorClaims []byte `protobuf:"bytes,4,opt,name=validator_claims,json=validatorClaims,proto3" json:"validator_claims,omitempty"`
}

func (m *DBProphecy) Reset()         { *m = DBProphecy{} }
func (m *DBProphecy) String() string { return proto.CompactTextString(m) }
func (*DBProphecy) ProtoMessage()    {}
func (*DBProphecy) Descriptor() ([]byte, []int) {
	return fileDescriptor_dac1b931484f4203, []int{2}
}
func (m *DBProphecy) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DBProphecy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DBProphecy.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DBProphecy) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DBProphecy.Merge(m, src)
}
func (m *DBProphecy) XXX_Size() int {
	return m.Size()
}
func (m *DBProphecy) XXX_DiscardUnknown() {
	xxx_messageInfo_DBProphecy.DiscardUnknown(m)
}

var xxx_messageInfo_DBProphecy proto.InternalMessageInfo

func (m *DBProphecy) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *DBProphecy) GetStatus() Status {
	if m != nil {
		return m.Status
	}
	return Status{}
}

func (m *DBProphecy) GetClaimValidators() []byte {
	if m != nil {
		return m.ClaimValidators
	}
	return nil
}

func (m *DBProphecy) GetValidatorClaims() []byte {
	if m != nil {
		return m.ValidatorClaims
	}
	return nil
}

// Status is a struct that contains the status of a given prophecy
type Status struct {
	Text       StatusText `protobuf:"varint,1,opt,name=text,proto3,enum=sifnode.oracle.v1.StatusText" json:"text,omitempty"`
	FinalClaim string     `protobuf:"bytes,2,opt,name=final_claim,json=finalClaim,proto3" json:"final_claim,omitempty"`
}

func (m *Status) Reset()         { *m = Status{} }
func (m *Status) String() string { return proto.CompactTextString(m) }
func (*Status) ProtoMessage()    {}
func (*Status) Descriptor() ([]byte, []int) {
	return fileDescriptor_dac1b931484f4203, []int{3}
}
func (m *Status) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Status) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Status.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Status) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Status.Merge(m, src)
}
func (m *Status) XXX_Size() int {
	return m.Size()
}
func (m *Status) XXX_DiscardUnknown() {
	xxx_messageInfo_Status.DiscardUnknown(m)
}

var xxx_messageInfo_Status proto.InternalMessageInfo

func (m *Status) GetText() StatusText {
	if m != nil {
		return m.Text
	}
	return StatusText_STATUS_TEXT_UNSPECIFIED
}

func (m *Status) GetFinalClaim() string {
	if m != nil {
		return m.FinalClaim
	}
	return ""
}

func init() {
	proto.RegisterEnum("sifnode.oracle.v1.StatusText", StatusText_name, StatusText_value)
	proto.RegisterType((*GenesisState)(nil), "sifnode.oracle.v1.GenesisState")
	proto.RegisterType((*Claim)(nil), "sifnode.oracle.v1.Claim")
	proto.RegisterType((*DBProphecy)(nil), "sifnode.oracle.v1.DBProphecy")
	proto.RegisterType((*Status)(nil), "sifnode.oracle.v1.Status")
}

func init() { proto.RegisterFile("sifnode/oracle/v1/types.proto", fileDescriptor_dac1b931484f4203) }

var fileDescriptor_dac1b931484f4203 = []byte{
	// 485 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0x41, 0x6e, 0xda, 0x40,
	0x18, 0x85, 0x31, 0x50, 0xaa, 0xfc, 0xd0, 0xd4, 0x99, 0x56, 0x8d, 0xdb, 0x2a, 0x0e, 0xa2, 0x1b,
	0x9a, 0x4a, 0xb6, 0x48, 0x17, 0x5d, 0x75, 0x01, 0xd8, 0x89, 0x90, 0x2a, 0x84, 0x6c, 0xd3, 0x56,
	0x55, 0x55, 0x6b, 0x62, 0x0f, 0x30, 0x92, 0xf1, 0x20, 0xcf, 0x84, 0x92, 0x5b, 0xf4, 0x06, 0x3d,
	0x40, 0x2f, 0x92, 0x65, 0x96, 0x5d, 0x55, 0x15, 0x5c, 0x24, 0x62, 0x6c, 0x03, 0x4a, 0xb2, 0xb3,
	0xdf, 0xf7, 0x66, 0xde, 0x9b, 0x5f, 0x3f, 0x1c, 0x71, 0x3a, 0x8a, 0x59, 0x48, 0x4c, 0x96, 0xe0,
	0x20, 0x22, 0xe6, 0xbc, 0x65, 0x8a, 0xab, 0x19, 0xe1, 0xc6, 0x2c, 0x61, 0x82, 0xa1, 0x83, 0x0c,
	0x1b, 0x29, 0x36, 0xe6, 0xad, 0x57, 0xcf, 0xc7, 0x6c, 0xcc, 0x24, 0x35, 0xd7, 0x5f, 0xa9, 0xb1,
	0xf1, 0x5b, 0x81, 0xda, 0x39, 0x89, 0x09, 0xa7, 0xdc, 0x15, 0x58, 0x10, 0xf4, 0x0e, 0x0e, 0x70,
	0x18, 0x26, 0x84, 0x73, 0xff, 0xe7, 0x84, 0x0a, 0x12, 0x51, 0x2e, 0x34, 0xa5, 0x5e, 0x6a, 0xee,
	0x39, 0x6a, 0x06, 0xbe, 0xe4, 0x3a, 0x7a, 0x03, 0x4f, 0x70, 0x38, 0xa5, 0xb1, 0x9f, 0x11, 0xad,
	0x58, 0x57, 0x9a, 0x7b, 0x4e, 0x4d, 0x8a, 0xed, 0x54, 0x43, 0x1f, 0x01, 0x66, 0x09, 0x9b, 0x4d,
	0x48, 0x40, 0x09, 0xd7, 0x4a, 0xf5, 0x52, 0xb3, 0x7a, 0x7a, 0x64, 0xdc, 0x2b, 0x68, 0x58, 0x9d,
	0x41, 0x6a, 0xbb, 0x72, 0x76, 0x0e, 0x34, 0x7e, 0xc0, 0xa3, 0x6e, 0x84, 0xe9, 0x14, 0xed, 0x43,
	0x91, 0x86, 0x9a, 0x22, 0x13, 0x8a, 0x34, 0x5c, 0x37, 0x9d, 0xe3, 0x88, 0x86, 0x58, 0xb0, 0xe4,
	0x4e, 0x01, 0x75, 0x03, 0xf2, 0x12, 0x1a, 0x3c, 0x0e, 0x58, 0x2c, 0x48, 0x2c, 0xb4, 0x92, 0xb4,
	0xe4, 0xbf, 0x8d, 0x3f, 0x0a, 0xc0, 0x36, 0xfa, 0x5e, 0xca, 0x07, 0xa8, 0x70, 0x81, 0xc5, 0x65,
	0x7a, 0x75, 0xf5, 0xf4, 0xe5, 0x03, 0xcd, 0x5d, 0x69, 0xe8, 0x94, 0xaf, 0xff, 0x1d, 0x17, 0x9c,
	0xcc, 0x8e, 0xde, 0x82, 0x1a, 0xac, 0x7b, 0xfb, 0x9b, 0x2e, 0x5c, 0x46, 0xd7, 0x9c, 0xa7, 0x52,
	0xff, 0xbc, 0x91, 0xd7, 0xd6, 0xed, 0x4b, 0x24, 0xe4, 0x5a, 0x39, 0xb5, 0x6e, 0x74, 0x39, 0x03,
	0xde, 0xf8, 0x0e, 0x95, 0x34, 0x0d, 0xb5, 0xa0, 0x2c, 0xc8, 0x42, 0xc8, 0xaa, 0xfb, 0x0f, 0x0e,
	0x34, 0x35, 0x7a, 0x64, 0x21, 0x1c, 0x69, 0x45, 0xc7, 0x50, 0x1d, 0xd1, 0x18, 0x47, 0x69, 0x46,
	0x36, 0x2b, 0x90, 0x92, 0xbc, 0xfe, 0x84, 0x03, 0x6c, 0x0f, 0xa1, 0xd7, 0x70, 0xe8, 0x7a, 0x6d,
	0x6f, 0xe8, 0xfa, 0x9e, 0xfd, 0xd5, 0xf3, 0x87, 0x7d, 0x77, 0x60, 0x77, 0x7b, 0x67, 0x3d, 0xdb,
	0x52, 0x0b, 0xe8, 0x10, 0x9e, 0xed, 0xc2, 0x81, 0xdd, 0xb7, 0x7a, 0xfd, 0x73, 0x55, 0xb9, 0x0b,
	0xdc, 0x61, 0xb7, 0x6b, 0xbb, 0xae, 0x5a, 0x44, 0x2f, 0x00, 0xed, 0x82, 0xb3, 0x76, 0xef, 0x93,
	0x6d, 0xa9, 0xa5, 0x8e, 0x75, 0xbd, 0xd4, 0x95, 0x9b, 0xa5, 0xae, 0xfc, 0x5f, 0xea, 0xca, 0xaf,
	0x95, 0x5e, 0xb8, 0x59, 0xe9, 0x85, 0xbf, 0x2b, 0xbd, 0xf0, 0xed, 0x64, 0x4c, 0xc5, 0xe4, 0xf2,
	0xc2, 0x08, 0xd8, 0xd4, 0x74, 0xe9, 0x28, 0x98, 0x60, 0x1a, 0x9b, 0xf9, 0xe2, 0x2f, 0xf2, 0xd5,
	0x97, 0x7b, 0x7f, 0x51, 0x91, 0xfb, 0xfc, 0xfe, 0x36, 0x00, 0x00, 0xff, 0xff, 0xc7, 0x8f, 0x95,
	0x2e, 0x19, 0x03, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Prophecies) > 0 {
		for iNdEx := len(m.Prophecies) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Prophecies[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTypes(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.AdminAddress) > 0 {
		i -= len(m.AdminAddress)
		copy(dAtA[i:], m.AdminAddress)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.AdminAddress)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.AddressWhitelist) > 0 {
		for iNdEx := len(m.AddressWhitelist) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.AddressWhitelist[iNdEx])
			copy(dAtA[i:], m.AddressWhitelist[iNdEx])
			i = encodeVarintTypes(dAtA, i, uint64(len(m.AddressWhitelist[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Claim) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Claim) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Claim) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Content) > 0 {
		i -= len(m.Content)
		copy(dAtA[i:], m.Content)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Content)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.ValidatorAddress) > 0 {
		i -= len(m.ValidatorAddress)
		copy(dAtA[i:], m.ValidatorAddress)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.ValidatorAddress)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *DBProphecy) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DBProphecy) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DBProphecy) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ValidatorClaims) > 0 {
		i -= len(m.ValidatorClaims)
		copy(dAtA[i:], m.ValidatorClaims)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.ValidatorClaims)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.ClaimValidators) > 0 {
		i -= len(m.ClaimValidators)
		copy(dAtA[i:], m.ClaimValidators)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.ClaimValidators)))
		i--
		dAtA[i] = 0x1a
	}
	{
		size, err := m.Status.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Status) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Status) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Status) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.FinalClaim) > 0 {
		i -= len(m.FinalClaim)
		copy(dAtA[i:], m.FinalClaim)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.FinalClaim)))
		i--
		dAtA[i] = 0x12
	}
	if m.Text != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Text))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintTypes(dAtA []byte, offset int, v uint64) int {
	offset -= sovTypes(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.AddressWhitelist) > 0 {
		for _, s := range m.AddressWhitelist {
			l = len(s)
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	l = len(m.AdminAddress)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if len(m.Prophecies) > 0 {
		for _, e := range m.Prophecies {
			l = e.Size()
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	return n
}

func (m *Claim) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.ValidatorAddress)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.Content)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}

func (m *DBProphecy) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = m.Status.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = len(m.ClaimValidators)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.ValidatorClaims)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}

func (m *Status) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Text != 0 {
		n += 1 + sovTypes(uint64(m.Text))
	}
	l = len(m.FinalClaim)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AddressWhitelist", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AddressWhitelist = append(m.AddressWhitelist, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AdminAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AdminAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Prophecies", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Prophecies = append(m.Prophecies, &DBProphecy{})
			if err := m.Prophecies[len(m.Prophecies)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *Claim) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: Claim: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Claim: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ValidatorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Content", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Content = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *DBProphecy) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: DBProphecy: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DBProphecy: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Status.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClaimValidators", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClaimValidators = append(m.ClaimValidators[:0], dAtA[iNdEx:postIndex]...)
			if m.ClaimValidators == nil {
				m.ClaimValidators = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorClaims", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ValidatorClaims = append(m.ValidatorClaims[:0], dAtA[iNdEx:postIndex]...)
			if m.ValidatorClaims == nil {
				m.ValidatorClaims = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *Status) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: Status: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Status: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Text", wireType)
			}
			m.Text = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Text |= StatusText(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FinalClaim", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FinalClaim = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func skipTypes(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
				return 0, ErrInvalidLengthTypes
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTypes
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTypes
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTypes        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTypes          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTypes = fmt.Errorf("proto: unexpected end of group")
)
