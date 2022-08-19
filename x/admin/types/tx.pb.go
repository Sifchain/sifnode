// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sifnode/admin/v1/tx.proto

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

type MsgAddAccount struct {
	Signer  string        `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	Account *AdminAccount `protobuf:"bytes,2,opt,name=account,proto3" json:"account,omitempty"`
}

func (m *MsgAddAccount) Reset()         { *m = MsgAddAccount{} }
func (m *MsgAddAccount) String() string { return proto.CompactTextString(m) }
func (*MsgAddAccount) ProtoMessage()    {}
func (*MsgAddAccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_600acd904f18192e, []int{0}
}
func (m *MsgAddAccount) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgAddAccount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgAddAccount.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgAddAccount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgAddAccount.Merge(m, src)
}
func (m *MsgAddAccount) XXX_Size() int {
	return m.Size()
}
func (m *MsgAddAccount) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgAddAccount.DiscardUnknown(m)
}

var xxx_messageInfo_MsgAddAccount proto.InternalMessageInfo

func (m *MsgAddAccount) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

func (m *MsgAddAccount) GetAccount() *AdminAccount {
	if m != nil {
		return m.Account
	}
	return nil
}

type MsgAddAccountResponse struct {
}

func (m *MsgAddAccountResponse) Reset()         { *m = MsgAddAccountResponse{} }
func (m *MsgAddAccountResponse) String() string { return proto.CompactTextString(m) }
func (*MsgAddAccountResponse) ProtoMessage()    {}
func (*MsgAddAccountResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_600acd904f18192e, []int{1}
}
func (m *MsgAddAccountResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgAddAccountResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgAddAccountResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgAddAccountResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgAddAccountResponse.Merge(m, src)
}
func (m *MsgAddAccountResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgAddAccountResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgAddAccountResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgAddAccountResponse proto.InternalMessageInfo

type MsgRemoveAccount struct {
	Signer  string        `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	Account *AdminAccount `protobuf:"bytes,2,opt,name=account,proto3" json:"account,omitempty"`
}

func (m *MsgRemoveAccount) Reset()         { *m = MsgRemoveAccount{} }
func (m *MsgRemoveAccount) String() string { return proto.CompactTextString(m) }
func (*MsgRemoveAccount) ProtoMessage()    {}
func (*MsgRemoveAccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_600acd904f18192e, []int{2}
}
func (m *MsgRemoveAccount) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgRemoveAccount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRemoveAccount.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgRemoveAccount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgRemoveAccount.Merge(m, src)
}
func (m *MsgRemoveAccount) XXX_Size() int {
	return m.Size()
}
func (m *MsgRemoveAccount) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgRemoveAccount.DiscardUnknown(m)
}

var xxx_messageInfo_MsgRemoveAccount proto.InternalMessageInfo

func (m *MsgRemoveAccount) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

func (m *MsgRemoveAccount) GetAccount() *AdminAccount {
	if m != nil {
		return m.Account
	}
	return nil
}

type MsgRemoveAccountResponse struct {
}

func (m *MsgRemoveAccountResponse) Reset()         { *m = MsgRemoveAccountResponse{} }
func (m *MsgRemoveAccountResponse) String() string { return proto.CompactTextString(m) }
func (*MsgRemoveAccountResponse) ProtoMessage()    {}
func (*MsgRemoveAccountResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_600acd904f18192e, []int{3}
}
func (m *MsgRemoveAccountResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgRemoveAccountResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRemoveAccountResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgRemoveAccountResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgRemoveAccountResponse.Merge(m, src)
}
func (m *MsgRemoveAccountResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgRemoveAccountResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgRemoveAccountResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgRemoveAccountResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgAddAccount)(nil), "sifnode.admin.v1.MsgAddAccount")
	proto.RegisterType((*MsgAddAccountResponse)(nil), "sifnode.admin.v1.MsgAddAccountResponse")
	proto.RegisterType((*MsgRemoveAccount)(nil), "sifnode.admin.v1.MsgRemoveAccount")
	proto.RegisterType((*MsgRemoveAccountResponse)(nil), "sifnode.admin.v1.MsgRemoveAccountResponse")
}

func init() { proto.RegisterFile("sifnode/admin/v1/tx.proto", fileDescriptor_600acd904f18192e) }

var fileDescriptor_600acd904f18192e = []byte{
	// 293 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2c, 0xce, 0x4c, 0xcb,
	0xcb, 0x4f, 0x49, 0xd5, 0x4f, 0x4c, 0xc9, 0xcd, 0xcc, 0xd3, 0x2f, 0x33, 0xd4, 0x2f, 0xa9, 0xd0,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x80, 0x4a, 0xe9, 0x81, 0xa5, 0xf4, 0xca, 0x0c, 0xa5,
	0x44, 0xd2, 0xf3, 0xd3, 0xf3, 0xc1, 0x92, 0xfa, 0x20, 0x16, 0x44, 0x9d, 0x94, 0x0c, 0xa6, 0x11,
	0x95, 0x05, 0xa9, 0xc5, 0x10, 0x59, 0xa5, 0x44, 0x2e, 0x5e, 0xdf, 0xe2, 0x74, 0xc7, 0x94, 0x14,
	0xc7, 0xe4, 0xe4, 0xfc, 0xd2, 0xbc, 0x12, 0x21, 0x31, 0x2e, 0xb6, 0xe2, 0xcc, 0xf4, 0xbc, 0xd4,
	0x22, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x28, 0x4f, 0xc8, 0x82, 0x8b, 0x3d, 0x11, 0xa2,
	0x44, 0x82, 0x49, 0x81, 0x51, 0x83, 0xdb, 0x48, 0x4e, 0x0f, 0xdd, 0x01, 0x7a, 0x8e, 0x20, 0x06,
	0xd4, 0xa0, 0x20, 0x98, 0x72, 0x25, 0x71, 0x2e, 0x51, 0x14, 0x2b, 0x82, 0x52, 0x8b, 0x0b, 0xf2,
	0xf3, 0x8a, 0x53, 0x95, 0x52, 0xb8, 0x04, 0x7c, 0x8b, 0xd3, 0x83, 0x52, 0x73, 0xf3, 0xcb, 0x52,
	0x69, 0x67, 0xbd, 0x14, 0x97, 0x04, 0xba, 0x2d, 0x30, 0x17, 0x18, 0x1d, 0x62, 0xe4, 0x62, 0xf6,
	0x2d, 0x4e, 0x17, 0x8a, 0xe0, 0xe2, 0x42, 0x0a, 0x02, 0x79, 0x4c, 0xa3, 0x51, 0x3c, 0x20, 0xa5,
	0x4e, 0x40, 0x01, 0xdc, 0x87, 0x0c, 0x42, 0x89, 0x5c, 0xbc, 0xa8, 0x1e, 0x54, 0xc2, 0xaa, 0x17,
	0x45, 0x8d, 0x94, 0x16, 0x61, 0x35, 0x08, 0x2b, 0x9c, 0x9c, 0x4f, 0x3c, 0x92, 0x63, 0xbc, 0xf0,
	0x48, 0x8e, 0xf1, 0xc1, 0x23, 0x39, 0xc6, 0x09, 0x8f, 0xe5, 0x18, 0x2e, 0x3c, 0x96, 0x63, 0xb8,
	0xf1, 0x58, 0x8e, 0x21, 0x4a, 0x33, 0x3d, 0xb3, 0x24, 0xa3, 0x34, 0x49, 0x2f, 0x39, 0x3f, 0x57,
	0x3f, 0x38, 0x33, 0x2d, 0x39, 0x23, 0x31, 0x33, 0x4f, 0x1f, 0x96, 0x1c, 0x2a, 0xa0, 0x09, 0x02,
	0x9c, 0x1a, 0x92, 0xd8, 0xc0, 0xc9, 0xc1, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xe2, 0xdf, 0x39,
	0xfd, 0x71, 0x02, 0x00, 0x00,
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
	AddAccount(ctx context.Context, in *MsgAddAccount, opts ...grpc.CallOption) (*MsgAddAccountResponse, error)
	RemoveAccount(ctx context.Context, in *MsgRemoveAccount, opts ...grpc.CallOption) (*MsgRemoveAccountResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) AddAccount(ctx context.Context, in *MsgAddAccount, opts ...grpc.CallOption) (*MsgAddAccountResponse, error) {
	out := new(MsgAddAccountResponse)
	err := c.cc.Invoke(ctx, "/sifnode.admin.v1.Msg/AddAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) RemoveAccount(ctx context.Context, in *MsgRemoveAccount, opts ...grpc.CallOption) (*MsgRemoveAccountResponse, error) {
	out := new(MsgRemoveAccountResponse)
	err := c.cc.Invoke(ctx, "/sifnode.admin.v1.Msg/RemoveAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	AddAccount(context.Context, *MsgAddAccount) (*MsgAddAccountResponse, error)
	RemoveAccount(context.Context, *MsgRemoveAccount) (*MsgRemoveAccountResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) AddAccount(ctx context.Context, req *MsgAddAccount) (*MsgAddAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddAccount not implemented")
}
func (*UnimplementedMsgServer) RemoveAccount(ctx context.Context, req *MsgRemoveAccount) (*MsgRemoveAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveAccount not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_AddAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgAddAccount)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).AddAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sifnode.admin.v1.Msg/AddAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).AddAccount(ctx, req.(*MsgAddAccount))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_RemoveAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRemoveAccount)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RemoveAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sifnode.admin.v1.Msg/RemoveAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RemoveAccount(ctx, req.(*MsgRemoveAccount))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "sifnode.admin.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddAccount",
			Handler:    _Msg_AddAccount_Handler,
		},
		{
			MethodName: "RemoveAccount",
			Handler:    _Msg_RemoveAccount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sifnode/admin/v1/tx.proto",
}

func (m *MsgAddAccount) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgAddAccount) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgAddAccount) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Account != nil {
		{
			size, err := m.Account.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTx(dAtA, i, uint64(size))
		}
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

func (m *MsgAddAccountResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgAddAccountResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgAddAccountResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgRemoveAccount) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgRemoveAccount) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgRemoveAccount) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Account != nil {
		{
			size, err := m.Account.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTx(dAtA, i, uint64(size))
		}
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

func (m *MsgRemoveAccountResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgRemoveAccountResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgRemoveAccountResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
func (m *MsgAddAccount) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if m.Account != nil {
		l = m.Account.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgAddAccountResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgRemoveAccount) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if m.Account != nil {
		l = m.Account.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgRemoveAccountResponse) Size() (n int) {
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
func (m *MsgAddAccount) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgAddAccount: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgAddAccount: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Account", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Account == nil {
				m.Account = &AdminAccount{}
			}
			if err := m.Account.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
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
func (m *MsgAddAccountResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgAddAccountResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgAddAccountResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *MsgRemoveAccount) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgRemoveAccount: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgRemoveAccount: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Account", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Account == nil {
				m.Account = &AdminAccount{}
			}
			if err := m.Account.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
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
func (m *MsgRemoveAccountResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgRemoveAccountResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgRemoveAccountResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
