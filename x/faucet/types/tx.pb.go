// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sifnode/faucet/v1/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type MsgRequestCoins struct {
	Coins     string `protobuf:"bytes,1,opt,name=coins,proto3" json:"coins,omitempty" yaml:"coins"`
	Requester string `protobuf:"bytes,2,opt,name=requester,proto3" json:"requester,omitempty" yaml:"requester"`
}

func (m *MsgRequestCoins) Reset()         { *m = MsgRequestCoins{} }
func (m *MsgRequestCoins) String() string { return proto.CompactTextString(m) }
func (*MsgRequestCoins) ProtoMessage()    {}
func (*MsgRequestCoins) Descriptor() ([]byte, []int) {
	return fileDescriptor_4fb67073b34c0d1d, []int{0}
}
func (m *MsgRequestCoins) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgRequestCoins) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRequestCoins.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgRequestCoins) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgRequestCoins.Merge(m, src)
}
func (m *MsgRequestCoins) XXX_Size() int {
	return m.Size()
}
func (m *MsgRequestCoins) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgRequestCoins.DiscardUnknown(m)
}

var xxx_messageInfo_MsgRequestCoins proto.InternalMessageInfo

func (m *MsgRequestCoins) GetCoins() string {
	if m != nil {
		return m.Coins
	}
	return ""
}

func (m *MsgRequestCoins) GetRequester() string {
	if m != nil {
		return m.Requester
	}
	return ""
}

type MsgRequestCoinsResponse struct {
}

func (m *MsgRequestCoinsResponse) Reset()         { *m = MsgRequestCoinsResponse{} }
func (m *MsgRequestCoinsResponse) String() string { return proto.CompactTextString(m) }
func (*MsgRequestCoinsResponse) ProtoMessage()    {}
func (*MsgRequestCoinsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4fb67073b34c0d1d, []int{1}
}
func (m *MsgRequestCoinsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgRequestCoinsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRequestCoinsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgRequestCoinsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgRequestCoinsResponse.Merge(m, src)
}
func (m *MsgRequestCoinsResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgRequestCoinsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgRequestCoinsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgRequestCoinsResponse proto.InternalMessageInfo

type MsgAddCoins struct {
	Signer string `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty" yaml:"signer"`
	Coins  string `protobuf:"bytes,2,opt,name=coins,proto3" json:"coins,omitempty" yaml:"coins"`
}

func (m *MsgAddCoins) Reset()         { *m = MsgAddCoins{} }
func (m *MsgAddCoins) String() string { return proto.CompactTextString(m) }
func (*MsgAddCoins) ProtoMessage()    {}
func (*MsgAddCoins) Descriptor() ([]byte, []int) {
	return fileDescriptor_4fb67073b34c0d1d, []int{2}
}
func (m *MsgAddCoins) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgAddCoins) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgAddCoins.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgAddCoins) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgAddCoins.Merge(m, src)
}
func (m *MsgAddCoins) XXX_Size() int {
	return m.Size()
}
func (m *MsgAddCoins) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgAddCoins.DiscardUnknown(m)
}

var xxx_messageInfo_MsgAddCoins proto.InternalMessageInfo

func (m *MsgAddCoins) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

func (m *MsgAddCoins) GetCoins() string {
	if m != nil {
		return m.Coins
	}
	return ""
}

type MsgAddCoinsResponse struct {
}

func (m *MsgAddCoinsResponse) Reset()         { *m = MsgAddCoinsResponse{} }
func (m *MsgAddCoinsResponse) String() string { return proto.CompactTextString(m) }
func (*MsgAddCoinsResponse) ProtoMessage()    {}
func (*MsgAddCoinsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4fb67073b34c0d1d, []int{3}
}
func (m *MsgAddCoinsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgAddCoinsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgAddCoinsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgAddCoinsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgAddCoinsResponse.Merge(m, src)
}
func (m *MsgAddCoinsResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgAddCoinsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgAddCoinsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgAddCoinsResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgRequestCoins)(nil), "sifnode.faucet.v1.MsgRequestCoins")
	proto.RegisterType((*MsgRequestCoinsResponse)(nil), "sifnode.faucet.v1.MsgRequestCoinsResponse")
	proto.RegisterType((*MsgAddCoins)(nil), "sifnode.faucet.v1.MsgAddCoins")
	proto.RegisterType((*MsgAddCoinsResponse)(nil), "sifnode.faucet.v1.MsgAddCoinsResponse")
}

func init() { proto.RegisterFile("sifnode/faucet/v1/tx.proto", fileDescriptor_4fb67073b34c0d1d) }

var fileDescriptor_4fb67073b34c0d1d = []byte{
	// 328 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2a, 0xce, 0x4c, 0xcb,
	0xcb, 0x4f, 0x49, 0xd5, 0x4f, 0x4b, 0x2c, 0x4d, 0x4e, 0x2d, 0xd1, 0x2f, 0x33, 0xd4, 0x2f, 0xa9,
	0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x84, 0xca, 0xe9, 0x41, 0xe4, 0xf4, 0xca, 0x0c,
	0xa5, 0x44, 0xd2, 0xf3, 0xd3, 0xf3, 0xc1, 0xb2, 0xfa, 0x20, 0x16, 0x44, 0xa1, 0x52, 0x2e, 0x17,
	0xbf, 0x6f, 0x71, 0x7a, 0x50, 0x6a, 0x61, 0x69, 0x6a, 0x71, 0x89, 0x73, 0x7e, 0x66, 0x5e, 0xb1,
	0x90, 0x1a, 0x17, 0x6b, 0x32, 0x88, 0x21, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0xe9, 0x24, 0xf0, 0xe9,
	0x9e, 0x3c, 0x4f, 0x65, 0x62, 0x6e, 0x8e, 0x95, 0x12, 0x58, 0x58, 0x29, 0x08, 0x22, 0x2d, 0x64,
	0xc4, 0xc5, 0x59, 0x04, 0xd1, 0x97, 0x5a, 0x24, 0xc1, 0x04, 0x56, 0x2b, 0xf2, 0xe9, 0x9e, 0xbc,
	0x00, 0x44, 0x2d, 0x5c, 0x4a, 0x29, 0x08, 0xa1, 0x4c, 0x49, 0x92, 0x4b, 0x1c, 0xcd, 0xba, 0xa0,
	0xd4, 0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54, 0xa5, 0x04, 0x2e, 0x6e, 0xdf, 0xe2, 0x74, 0xc7, 0x94,
	0x14, 0x88, 0x2b, 0x34, 0xb9, 0xd8, 0x8a, 0x33, 0xd3, 0xf3, 0x52, 0x8b, 0xa0, 0xce, 0x10, 0xfc,
	0x74, 0x4f, 0x9e, 0x17, 0x62, 0x34, 0x44, 0x5c, 0x29, 0x08, 0xaa, 0x00, 0xe1, 0x60, 0x26, 0xbc,
	0x0e, 0x56, 0x12, 0xe5, 0x12, 0x46, 0xb2, 0x01, 0x66, 0xb1, 0xd1, 0x4e, 0x46, 0x2e, 0x66, 0xdf,
	0xe2, 0x74, 0xa1, 0x38, 0x2e, 0x1e, 0x94, 0x70, 0x50, 0xd2, 0xc3, 0x08, 0x44, 0x3d, 0x34, 0xc7,
	0x4b, 0x69, 0x11, 0x56, 0x03, 0xb3, 0x47, 0x28, 0x88, 0x8b, 0x03, 0xee, 0x3b, 0x39, 0xec, 0xfa,
	0x60, 0xf2, 0x52, 0x6a, 0xf8, 0xe5, 0x61, 0x66, 0x3a, 0xb9, 0x9c, 0x78, 0x24, 0xc7, 0x78, 0xe1,
	0x91, 0x1c, 0xe3, 0x83, 0x47, 0x72, 0x8c, 0x13, 0x1e, 0xcb, 0x31, 0x5c, 0x78, 0x2c, 0xc7, 0x70,
	0xe3, 0xb1, 0x1c, 0x43, 0x94, 0x56, 0x7a, 0x66, 0x49, 0x46, 0x69, 0x92, 0x5e, 0x72, 0x7e, 0xae,
	0x7e, 0x71, 0x66, 0x5a, 0x72, 0x46, 0x62, 0x66, 0x9e, 0x3e, 0x2c, 0xc5, 0x54, 0xc0, 0xd2, 0x4c,
	0x49, 0x65, 0x41, 0x6a, 0x71, 0x12, 0x1b, 0x38, 0x2d, 0x18, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff,
	0x83, 0x6b, 0x62, 0x3a, 0x52, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	RequestCoins(ctx context.Context, in *MsgRequestCoins, opts ...grpc.CallOption) (*MsgRequestCoinsResponse, error)
	AddCoins(ctx context.Context, in *MsgAddCoins, opts ...grpc.CallOption) (*MsgAddCoinsResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) RequestCoins(ctx context.Context, in *MsgRequestCoins, opts ...grpc.CallOption) (*MsgRequestCoinsResponse, error) {
	out := new(MsgRequestCoinsResponse)
	err := c.cc.Invoke(ctx, "/sifnode.faucet.v1.Msg/RequestCoins", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) AddCoins(ctx context.Context, in *MsgAddCoins, opts ...grpc.CallOption) (*MsgAddCoinsResponse, error) {
	out := new(MsgAddCoinsResponse)
	err := c.cc.Invoke(ctx, "/sifnode.faucet.v1.Msg/AddCoins", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	RequestCoins(context.Context, *MsgRequestCoins) (*MsgRequestCoinsResponse, error)
	AddCoins(context.Context, *MsgAddCoins) (*MsgAddCoinsResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) RequestCoins(ctx context.Context, req *MsgRequestCoins) (*MsgRequestCoinsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestCoins not implemented")
}
func (*UnimplementedMsgServer) AddCoins(ctx context.Context, req *MsgAddCoins) (*MsgAddCoinsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddCoins not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_RequestCoins_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRequestCoins)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RequestCoins(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sifnode.faucet.v1.Msg/RequestCoins",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RequestCoins(ctx, req.(*MsgRequestCoins))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_AddCoins_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgAddCoins)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).AddCoins(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sifnode.faucet.v1.Msg/AddCoins",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).AddCoins(ctx, req.(*MsgAddCoins))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "sifnode.faucet.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestCoins",
			Handler:    _Msg_RequestCoins_Handler,
		},
		{
			MethodName: "AddCoins",
			Handler:    _Msg_AddCoins_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sifnode/faucet/v1/tx.proto",
}

func (m *MsgRequestCoins) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgRequestCoins) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgRequestCoins) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Requester) > 0 {
		i -= len(m.Requester)
		copy(dAtA[i:], m.Requester)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Requester)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Coins) > 0 {
		i -= len(m.Coins)
		copy(dAtA[i:], m.Coins)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Coins)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgRequestCoinsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgRequestCoinsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgRequestCoinsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgAddCoins) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgAddCoins) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgAddCoins) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Coins) > 0 {
		i -= len(m.Coins)
		copy(dAtA[i:], m.Coins)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Coins)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgAddCoinsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgAddCoinsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgAddCoinsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgRequestCoins) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Coins)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Requester)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgRequestCoinsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgAddCoins) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Coins)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgAddCoinsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgRequestCoins) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgRequestCoins: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgRequestCoins: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Coins", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Coins = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Requester", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Requester = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgRequestCoinsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgRequestCoinsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgRequestCoinsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgAddCoins) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgAddCoins: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgAddCoins: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Coins", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Coins = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgAddCoinsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgAddCoinsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgAddCoinsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)
