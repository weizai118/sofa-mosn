// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: envoy/api/v2/eds.proto

package v2

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import envoy_api_v2_endpoint "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
import _ "github.com/gogo/googleapis/google/api"
import _ "github.com/lyft/protoc-gen-validate/validate"
import _ "github.com/gogo/protobuf/gogoproto"

import context "golang.org/x/net/context"
import grpc "google.golang.org/grpc"

import binary "encoding/binary"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Each route from RDS will map to a single cluster or traffic split across
// clusters using weights expressed in the RDS WeightedCluster.
//
// With EDS, each cluster is treated independently from a LB perspective, with
// LB taking place between the Localities within a cluster and at a finer
// granularity between the hosts within a locality. For a given cluster, the
// effective weight of a host is its load_balancing_weight multiplied by the
// load_balancing_weight of its Locality.
type ClusterLoadAssignment struct {
	// Name of the cluster. This will be the :ref:`service_name
	// <envoy_api_field_Cluster.EdsClusterConfig.service_name>` value if specified
	// in the cluster :ref:`EdsClusterConfig
	// <envoy_api_msg_Cluster.EdsClusterConfig>`.
	ClusterName string `protobuf:"bytes,1,opt,name=cluster_name,json=clusterName,proto3" json:"cluster_name,omitempty"`
	// List of endpoints to load balance to.
	Endpoints []envoy_api_v2_endpoint.LocalityLbEndpoints `protobuf:"bytes,2,rep,name=endpoints" json:"endpoints"`
	// Load balancing policy settings.
	Policy *ClusterLoadAssignment_Policy `protobuf:"bytes,4,opt,name=policy" json:"policy,omitempty"`
}

func (m *ClusterLoadAssignment) Reset()                    { *m = ClusterLoadAssignment{} }
func (m *ClusterLoadAssignment) String() string            { return proto.CompactTextString(m) }
func (*ClusterLoadAssignment) ProtoMessage()               {}
func (*ClusterLoadAssignment) Descriptor() ([]byte, []int) { return fileDescriptorEds, []int{0} }

func (m *ClusterLoadAssignment) GetClusterName() string {
	if m != nil {
		return m.ClusterName
	}
	return ""
}

func (m *ClusterLoadAssignment) GetEndpoints() []envoy_api_v2_endpoint.LocalityLbEndpoints {
	if m != nil {
		return m.Endpoints
	}
	return nil
}

func (m *ClusterLoadAssignment) GetPolicy() *ClusterLoadAssignment_Policy {
	if m != nil {
		return m.Policy
	}
	return nil
}

// Load balancing policy settings.
type ClusterLoadAssignment_Policy struct {
	// Percentage of traffic (0-100) that should be dropped. This
	// action allows protection of upstream hosts should they unable to
	// recover from an outage or should they be unable to autoscale and hence
	// overall incoming traffic volume need to be trimmed to protect them.
	// [#v2-api-diff: This is known as maintenance mode in v1.]
	DropOverload float64 `protobuf:"fixed64,1,opt,name=drop_overload,json=dropOverload,proto3" json:"drop_overload,omitempty"`
}

func (m *ClusterLoadAssignment_Policy) Reset()         { *m = ClusterLoadAssignment_Policy{} }
func (m *ClusterLoadAssignment_Policy) String() string { return proto.CompactTextString(m) }
func (*ClusterLoadAssignment_Policy) ProtoMessage()    {}
func (*ClusterLoadAssignment_Policy) Descriptor() ([]byte, []int) {
	return fileDescriptorEds, []int{0, 0}
}

func (m *ClusterLoadAssignment_Policy) GetDropOverload() float64 {
	if m != nil {
		return m.DropOverload
	}
	return 0
}

func init() {
	proto.RegisterType((*ClusterLoadAssignment)(nil), "envoy.api.v2.ClusterLoadAssignment")
	proto.RegisterType((*ClusterLoadAssignment_Policy)(nil), "envoy.api.v2.ClusterLoadAssignment.Policy")
}
func (this *ClusterLoadAssignment) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ClusterLoadAssignment)
	if !ok {
		that2, ok := that.(ClusterLoadAssignment)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.ClusterName != that1.ClusterName {
		return false
	}
	if len(this.Endpoints) != len(that1.Endpoints) {
		return false
	}
	for i := range this.Endpoints {
		if !this.Endpoints[i].Equal(&that1.Endpoints[i]) {
			return false
		}
	}
	if !this.Policy.Equal(that1.Policy) {
		return false
	}
	return true
}
func (this *ClusterLoadAssignment_Policy) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ClusterLoadAssignment_Policy)
	if !ok {
		that2, ok := that.(ClusterLoadAssignment_Policy)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.DropOverload != that1.DropOverload {
		return false
	}
	return true
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for EndpointDiscoveryService service

type EndpointDiscoveryServiceClient interface {
	// The resource_names field in DiscoveryRequest specifies a list of clusters
	// to subscribe to updates for.
	StreamEndpoints(ctx context.Context, opts ...grpc.CallOption) (EndpointDiscoveryService_StreamEndpointsClient, error)
	FetchEndpoints(ctx context.Context, in *DiscoveryRequest, opts ...grpc.CallOption) (*DiscoveryResponse, error)
}

type endpointDiscoveryServiceClient struct {
	cc *grpc.ClientConn
}

func NewEndpointDiscoveryServiceClient(cc *grpc.ClientConn) EndpointDiscoveryServiceClient {
	return &endpointDiscoveryServiceClient{cc}
}

func (c *endpointDiscoveryServiceClient) StreamEndpoints(ctx context.Context, opts ...grpc.CallOption) (EndpointDiscoveryService_StreamEndpointsClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_EndpointDiscoveryService_serviceDesc.Streams[0], c.cc, "/envoy.api.v2.EndpointDiscoveryService/StreamEndpoints", opts...)
	if err != nil {
		return nil, err
	}
	x := &endpointDiscoveryServiceStreamEndpointsClient{stream}
	return x, nil
}

type EndpointDiscoveryService_StreamEndpointsClient interface {
	Send(*DiscoveryRequest) error
	Recv() (*DiscoveryResponse, error)
	grpc.ClientStream
}

type endpointDiscoveryServiceStreamEndpointsClient struct {
	grpc.ClientStream
}

func (x *endpointDiscoveryServiceStreamEndpointsClient) Send(m *DiscoveryRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *endpointDiscoveryServiceStreamEndpointsClient) Recv() (*DiscoveryResponse, error) {
	m := new(DiscoveryResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *endpointDiscoveryServiceClient) FetchEndpoints(ctx context.Context, in *DiscoveryRequest, opts ...grpc.CallOption) (*DiscoveryResponse, error) {
	out := new(DiscoveryResponse)
	err := grpc.Invoke(ctx, "/envoy.api.v2.EndpointDiscoveryService/FetchEndpoints", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for EndpointDiscoveryService service

type EndpointDiscoveryServiceServer interface {
	// The resource_names field in DiscoveryRequest specifies a list of clusters
	// to subscribe to updates for.
	StreamEndpoints(EndpointDiscoveryService_StreamEndpointsServer) error
	FetchEndpoints(context.Context, *DiscoveryRequest) (*DiscoveryResponse, error)
}

func RegisterEndpointDiscoveryServiceServer(s *grpc.Server, srv EndpointDiscoveryServiceServer) {
	s.RegisterService(&_EndpointDiscoveryService_serviceDesc, srv)
}

func _EndpointDiscoveryService_StreamEndpoints_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(EndpointDiscoveryServiceServer).StreamEndpoints(&endpointDiscoveryServiceStreamEndpointsServer{stream})
}

type EndpointDiscoveryService_StreamEndpointsServer interface {
	Send(*DiscoveryResponse) error
	Recv() (*DiscoveryRequest, error)
	grpc.ServerStream
}

type endpointDiscoveryServiceStreamEndpointsServer struct {
	grpc.ServerStream
}

func (x *endpointDiscoveryServiceStreamEndpointsServer) Send(m *DiscoveryResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *endpointDiscoveryServiceStreamEndpointsServer) Recv() (*DiscoveryRequest, error) {
	m := new(DiscoveryRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _EndpointDiscoveryService_FetchEndpoints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DiscoveryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointDiscoveryServiceServer).FetchEndpoints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/envoy.api.v2.EndpointDiscoveryService/FetchEndpoints",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointDiscoveryServiceServer).FetchEndpoints(ctx, req.(*DiscoveryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _EndpointDiscoveryService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "envoy.api.v2.EndpointDiscoveryService",
	HandlerType: (*EndpointDiscoveryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FetchEndpoints",
			Handler:    _EndpointDiscoveryService_FetchEndpoints_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamEndpoints",
			Handler:       _EndpointDiscoveryService_StreamEndpoints_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "envoy/api/v2/eds.proto",
}

func (m *ClusterLoadAssignment) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ClusterLoadAssignment) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.ClusterName) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintEds(dAtA, i, uint64(len(m.ClusterName)))
		i += copy(dAtA[i:], m.ClusterName)
	}
	if len(m.Endpoints) > 0 {
		for _, msg := range m.Endpoints {
			dAtA[i] = 0x12
			i++
			i = encodeVarintEds(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.Policy != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintEds(dAtA, i, uint64(m.Policy.Size()))
		n1, err := m.Policy.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	return i, nil
}

func (m *ClusterLoadAssignment_Policy) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ClusterLoadAssignment_Policy) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.DropOverload != 0 {
		dAtA[i] = 0x9
		i++
		binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.DropOverload))))
		i += 8
	}
	return i, nil
}

func encodeVarintEds(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *ClusterLoadAssignment) Size() (n int) {
	var l int
	_ = l
	l = len(m.ClusterName)
	if l > 0 {
		n += 1 + l + sovEds(uint64(l))
	}
	if len(m.Endpoints) > 0 {
		for _, e := range m.Endpoints {
			l = e.Size()
			n += 1 + l + sovEds(uint64(l))
		}
	}
	if m.Policy != nil {
		l = m.Policy.Size()
		n += 1 + l + sovEds(uint64(l))
	}
	return n
}

func (m *ClusterLoadAssignment_Policy) Size() (n int) {
	var l int
	_ = l
	if m.DropOverload != 0 {
		n += 9
	}
	return n
}

func sovEds(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozEds(x uint64) (n int) {
	return sovEds(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ClusterLoadAssignment) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEds
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ClusterLoadAssignment: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ClusterLoadAssignment: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClusterName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEds
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClusterName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Endpoints", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthEds
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Endpoints = append(m.Endpoints, envoy_api_v2_endpoint.LocalityLbEndpoints{})
			if err := m.Endpoints[len(m.Endpoints)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Policy", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthEds
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Policy == nil {
				m.Policy = &ClusterLoadAssignment_Policy{}
			}
			if err := m.Policy.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEds(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthEds
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
func (m *ClusterLoadAssignment_Policy) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEds
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Policy: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Policy: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field DropOverload", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.DropOverload = float64(math.Float64frombits(v))
		default:
			iNdEx = preIndex
			skippy, err := skipEds(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthEds
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
func skipEds(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEds
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
					return 0, ErrIntOverflowEds
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
					return 0, ErrIntOverflowEds
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
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthEds
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowEds
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
				next, err := skipEds(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
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
	ErrInvalidLengthEds = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEds   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("envoy/api/v2/eds.proto", fileDescriptorEds) }

var fileDescriptorEds = []byte{
	// 426 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x52, 0x3d, 0x8b, 0xd4, 0x40,
	0x18, 0xde, 0xd9, 0x5b, 0x16, 0x76, 0x6e, 0x55, 0x18, 0x4f, 0x6f, 0x2f, 0x2c, 0xb9, 0x25, 0x58,
	0xac, 0x87, 0x24, 0x12, 0xbb, 0x2b, 0x44, 0xe3, 0x07, 0x16, 0xcb, 0x29, 0xb9, 0x46, 0xab, 0x63,
	0x2e, 0x79, 0x89, 0x03, 0xc9, 0xbc, 0x31, 0x33, 0x1b, 0x48, 0x6b, 0x65, 0xef, 0x4f, 0xb0, 0xf1,
	0x37, 0x58, 0x59, 0x5e, 0x29, 0xd8, 0x8b, 0x04, 0x1b, 0xf1, 0x4f, 0x48, 0x3e, 0x35, 0xa8, 0x9d,
	0x53, 0xbd, 0xbc, 0xcf, 0x47, 0x9e, 0x67, 0x32, 0xf4, 0x3a, 0xc8, 0x1c, 0x0b, 0x87, 0xa7, 0xc2,
	0xc9, 0x5d, 0x07, 0x42, 0x65, 0xa7, 0x19, 0x6a, 0x64, 0xf3, 0x7a, 0x6f, 0xf3, 0x54, 0xd8, 0xb9,
	0x6b, 0x2c, 0x07, 0xac, 0x50, 0xa8, 0x00, 0x73, 0xc8, 0x8a, 0x86, 0x6b, 0xdc, 0x18, 0x7a, 0xc8,
	0x30, 0x45, 0x21, 0x75, 0x3f, 0xb4, 0xac, 0x65, 0x84, 0x18, 0xc5, 0x50, 0xd3, 0xb8, 0x94, 0xa8,
	0xb9, 0x16, 0x28, 0xdb, 0xef, 0x19, 0xfb, 0x39, 0x8f, 0x45, 0xc8, 0x35, 0x38, 0xdd, 0xd0, 0x02,
	0x7b, 0x11, 0x46, 0x58, 0x8f, 0x4e, 0x35, 0x35, 0x5b, 0xeb, 0xdd, 0x98, 0x5e, 0x7b, 0x10, 0x6f,
	0x95, 0x86, 0x6c, 0x83, 0x3c, 0xbc, 0xaf, 0x94, 0x88, 0x64, 0x02, 0x52, 0xb3, 0x5b, 0x74, 0x1e,
	0x34, 0xc0, 0x99, 0xe4, 0x09, 0x2c, 0xc8, 0x8a, 0xac, 0x67, 0xde, 0xec, 0xc3, 0xf7, 0x8f, 0x3b,
	0x93, 0x6c, 0xbc, 0x22, 0xfe, 0x6e, 0x0b, 0x9f, 0xf0, 0x04, 0xd8, 0x09, 0x9d, 0x75, 0x31, 0xd5,
	0x62, 0xbc, 0xda, 0x59, 0xef, 0xba, 0x47, 0xf6, 0xef, 0xd5, 0xed, 0xbe, 0xc5, 0x06, 0x03, 0x1e,
	0x0b, 0x5d, 0x6c, 0xce, 0x1f, 0x75, 0x0a, 0x6f, 0x72, 0xf1, 0xe5, 0x70, 0xe4, 0xff, 0xb2, 0x60,
	0x1e, 0x9d, 0xa6, 0x18, 0x8b, 0xa0, 0x58, 0x4c, 0x56, 0xe4, 0x4f, 0xb3, 0xbf, 0x46, 0xb6, 0x9f,
	0xd5, 0x0a, 0xbf, 0x55, 0x1a, 0x4f, 0xe8, 0xb4, 0xd9, 0xb0, 0xbb, 0xf4, 0x52, 0x98, 0x61, 0x7a,
	0x56, 0x5d, 0x76, 0x8c, 0x3c, 0xac, 0xcb, 0x10, 0xef, 0xa0, 0x2a, 0xb3, 0xc7, 0xd8, 0xc1, 0xa8,
	0x3e, 0x2f, 0xee, 0xdd, 0x1c, 0xb5, 0xc7, 0x9f, 0x57, 0xfc, 0xa7, 0x2d, 0xdd, 0xfd, 0x41, 0xe8,
	0xa2, 0x0b, 0xfb, 0xb0, 0xfb, 0x69, 0xa7, 0x90, 0xe5, 0x22, 0x00, 0xf6, 0x9c, 0x5e, 0x39, 0xd5,
	0x19, 0xf0, 0xa4, 0xaf, 0xc3, 0xcc, 0x61, 0xda, 0x5e, 0xe2, 0xc3, 0xab, 0x2d, 0x28, 0x6d, 0x1c,
	0xfe, 0x13, 0x57, 0x29, 0x4a, 0x05, 0xd6, 0x68, 0x4d, 0x6e, 0x13, 0xb6, 0xa5, 0x97, 0x1f, 0x83,
	0x0e, 0x5e, 0xfe, 0x47, 0x63, 0xeb, 0xf5, 0xe7, 0x6f, 0x6f, 0xc7, 0x4b, 0x6b, 0x7f, 0xf0, 0xfe,
	0x8e, 0xfb, 0x8b, 0x3f, 0x26, 0x47, 0xde, 0xd5, 0xf7, 0xa5, 0x49, 0x2e, 0x4a, 0x93, 0x7c, 0x2a,
	0x4d, 0xf2, 0xb5, 0x34, 0xc9, 0x1b, 0x42, 0xce, 0xa7, 0xf5, 0x7b, 0xb9, 0xf3, 0x33, 0x00, 0x00,
	0xff, 0xff, 0x73, 0x0b, 0x98, 0x67, 0xe8, 0x02, 0x00, 0x00,
}
