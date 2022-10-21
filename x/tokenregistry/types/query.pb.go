// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: tokenregistry/v1/query.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type QueryEntriesResponse struct {
	Registry *Registry `protobuf:"bytes,1,opt,name=registry,proto3" json:"registry,omitempty"`
}

func (m *QueryEntriesResponse) Reset()         { *m = QueryEntriesResponse{} }
func (m *QueryEntriesResponse) String() string { return proto.CompactTextString(m) }
func (*QueryEntriesResponse) ProtoMessage()    {}
func (*QueryEntriesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8091a4dbfcdf886a, []int{0}
}
func (m *QueryEntriesResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryEntriesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryEntriesResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryEntriesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryEntriesResponse.Merge(m, src)
}
func (m *QueryEntriesResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryEntriesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryEntriesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryEntriesResponse proto.InternalMessageInfo

func (m *QueryEntriesResponse) GetRegistry() *Registry {
	if m != nil {
		return m.Registry
	}
	return nil
}

type QueryEntriesRequest struct {
}

func (m *QueryEntriesRequest) Reset()         { *m = QueryEntriesRequest{} }
func (m *QueryEntriesRequest) String() string { return proto.CompactTextString(m) }
func (*QueryEntriesRequest) ProtoMessage()    {}
func (*QueryEntriesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8091a4dbfcdf886a, []int{1}
}
func (m *QueryEntriesRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryEntriesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryEntriesRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryEntriesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryEntriesRequest.Merge(m, src)
}
func (m *QueryEntriesRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryEntriesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryEntriesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryEntriesRequest proto.InternalMessageInfo

func init() {
	proto.RegisterType((*QueryEntriesResponse)(nil), "sifnode.tokenregistry.v1.QueryEntriesResponse")
	proto.RegisterType((*QueryEntriesRequest)(nil), "sifnode.tokenregistry.v1.QueryEntriesRequest")
}

func init() { proto.RegisterFile("tokenregistry/v1/query.proto", fileDescriptor_8091a4dbfcdf886a) }

var fileDescriptor_8091a4dbfcdf886a = []byte{
	// 291 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x29, 0xc9, 0xcf, 0x4e,
	0xcd, 0x2b, 0x4a, 0x4d, 0xcf, 0x2c, 0x2e, 0x29, 0xaa, 0xd4, 0x2f, 0x33, 0xd4, 0x2f, 0x2c, 0x4d,
	0x2d, 0xaa, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x92, 0x28, 0xce, 0x4c, 0xcb, 0xcb, 0x4f,
	0x49, 0xd5, 0x43, 0x51, 0xa5, 0x57, 0x66, 0x28, 0x25, 0x92, 0x9e, 0x9f, 0x9e, 0x0f, 0x56, 0xa4,
	0x0f, 0x62, 0x41, 0xd4, 0x4b, 0xc9, 0xa4, 0xe7, 0xe7, 0xa7, 0xe7, 0xa4, 0xea, 0x27, 0x16, 0x64,
	0xea, 0x27, 0xe6, 0xe5, 0xe5, 0x97, 0x24, 0x96, 0x64, 0xe6, 0xe7, 0x15, 0xc3, 0x64, 0x31, 0xec,
	0x2a, 0xa9, 0x2c, 0x48, 0x85, 0xca, 0x2a, 0x85, 0x71, 0x89, 0x04, 0x82, 0xac, 0x76, 0xcd, 0x2b,
	0x29, 0xca, 0x4c, 0x2d, 0x0e, 0x4a, 0x2d, 0x2e, 0xc8, 0xcf, 0x2b, 0x4e, 0x15, 0xb2, 0xe3, 0xe2,
	0x80, 0x69, 0x91, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x36, 0x52, 0xd2, 0xc3, 0xe5, 0x2c, 0xbd, 0x20,
	0x28, 0x3b, 0x08, 0xae, 0x47, 0x49, 0x94, 0x4b, 0x18, 0xd5, 0xdc, 0xc2, 0xd2, 0xd4, 0xe2, 0x12,
	0xa3, 0xc5, 0x8c, 0x5c, 0xac, 0x60, 0x71, 0xa1, 0x99, 0x8c, 0x5c, 0xec, 0x50, 0x49, 0x21, 0x5d,
	0xdc, 0x46, 0x63, 0x31, 0x44, 0x4a, 0x8f, 0x58, 0xe5, 0x10, 0xbf, 0x28, 0xe9, 0x37, 0x5d, 0x7e,
	0x32, 0x99, 0x49, 0x53, 0x48, 0x5d, 0xbf, 0x38, 0x33, 0x2d, 0x39, 0x23, 0x31, 0x33, 0x4f, 0x1f,
	0x3d, 0x4c, 0x92, 0x52, 0x4b, 0x12, 0x0d, 0xf5, 0x53, 0x21, 0x1a, 0x9d, 0xbc, 0x4f, 0x3c, 0x92,
	0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23, 0x39, 0xc6, 0x09, 0x8f, 0xe5, 0x18, 0x2e, 0x3c,
	0x96, 0x63, 0xb8, 0xf1, 0x58, 0x8e, 0x21, 0xca, 0x30, 0x3d, 0xb3, 0x24, 0xa3, 0x34, 0x49, 0x2f,
	0x39, 0x3f, 0x57, 0x3f, 0x18, 0x66, 0x18, 0xd4, 0x35, 0xfa, 0x15, 0x68, 0xc6, 0x82, 0xc3, 0x39,
	0x89, 0x0d, 0x1c, 0xd0, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x8c, 0x51, 0xa9, 0x1c, 0xf4,
	0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	Entries(ctx context.Context, in *QueryEntriesRequest, opts ...grpc.CallOption) (*QueryEntriesResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Entries(ctx context.Context, in *QueryEntriesRequest, opts ...grpc.CallOption) (*QueryEntriesResponse, error) {
	out := new(QueryEntriesResponse)
	err := c.cc.Invoke(ctx, "/sifnode.tokenregistry.v1.Query/Entries", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	Entries(context.Context, *QueryEntriesRequest) (*QueryEntriesResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Entries(ctx context.Context, req *QueryEntriesRequest) (*QueryEntriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Entries not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_Entries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryEntriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Entries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sifnode.tokenregistry.v1.Query/Entries",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Entries(ctx, req.(*QueryEntriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "sifnode.tokenregistry.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Entries",
			Handler:    _Query_Entries_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tokenregistry/v1/query.proto",
}

func (m *QueryEntriesResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryEntriesResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryEntriesResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Registry != nil {
		{
			size, err := m.Registry.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryEntriesRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryEntriesRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryEntriesRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryEntriesResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Registry != nil {
		l = m.Registry.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryEntriesRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryEntriesResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryEntriesResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryEntriesResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Registry", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Registry == nil {
				m.Registry = &Registry{}
			}
			if err := m.Registry.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryEntriesRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryEntriesRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryEntriesRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
