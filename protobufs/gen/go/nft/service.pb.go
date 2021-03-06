// Code generated by protoc-gen-go. DO NOT EDIT.
// source: nft/service.proto

package nftpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
import _ "google.golang.org/genproto/googleapis/api/annotations"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ResponseHeader struct {
	TransactionId        string   `protobuf:"bytes,5,opt,name=transaction_id,json=transactionId,proto3" json:"transaction_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResponseHeader) Reset()         { *m = ResponseHeader{} }
func (m *ResponseHeader) String() string { return proto.CompactTextString(m) }
func (*ResponseHeader) ProtoMessage()    {}
func (*ResponseHeader) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_10e4a4ecba67c7da, []int{0}
}
func (m *ResponseHeader) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResponseHeader.Unmarshal(m, b)
}
func (m *ResponseHeader) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResponseHeader.Marshal(b, m, deterministic)
}
func (dst *ResponseHeader) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResponseHeader.Merge(dst, src)
}
func (m *ResponseHeader) XXX_Size() int {
	return xxx_messageInfo_ResponseHeader.Size(m)
}
func (m *ResponseHeader) XXX_DiscardUnknown() {
	xxx_messageInfo_ResponseHeader.DiscardUnknown(m)
}

var xxx_messageInfo_ResponseHeader proto.InternalMessageInfo

func (m *ResponseHeader) GetTransactionId() string {
	if m != nil {
		return m.TransactionId
	}
	return ""
}

type NFTMintRequest struct {
	// Document identifier
	Identifier string `protobuf:"bytes,1,opt,name=identifier,proto3" json:"identifier,omitempty"`
	// The contract address of the registry where the token should be minted
	RegistryAddress string   `protobuf:"bytes,2,opt,name=registry_address,json=registryAddress,proto3" json:"registry_address,omitempty"`
	DepositAddress  string   `protobuf:"bytes,3,opt,name=deposit_address,json=depositAddress,proto3" json:"deposit_address,omitempty"`
	ProofFields     []string `protobuf:"bytes,4,rep,name=proof_fields,json=proofFields,proto3" json:"proof_fields,omitempty"`
	// proof that nft is part of document
	SubmitTokenProof bool `protobuf:"varint,5,opt,name=submit_token_proof,json=submitTokenProof,proto3" json:"submit_token_proof,omitempty"`
	// proof that nft owner can access the document if nft_grant_access is true
	SubmitNftOwnerAccessProof bool `protobuf:"varint,7,opt,name=submit_nft_owner_access_proof,json=submitNftOwnerAccessProof,proto3" json:"submit_nft_owner_access_proof,omitempty"`
	// grant nft read access to the document
	GrantNftAccess       bool     `protobuf:"varint,8,opt,name=grant_nft_access,json=grantNftAccess,proto3" json:"grant_nft_access,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NFTMintRequest) Reset()         { *m = NFTMintRequest{} }
func (m *NFTMintRequest) String() string { return proto.CompactTextString(m) }
func (*NFTMintRequest) ProtoMessage()    {}
func (*NFTMintRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_10e4a4ecba67c7da, []int{1}
}
func (m *NFTMintRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NFTMintRequest.Unmarshal(m, b)
}
func (m *NFTMintRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NFTMintRequest.Marshal(b, m, deterministic)
}
func (dst *NFTMintRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NFTMintRequest.Merge(dst, src)
}
func (m *NFTMintRequest) XXX_Size() int {
	return xxx_messageInfo_NFTMintRequest.Size(m)
}
func (m *NFTMintRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_NFTMintRequest.DiscardUnknown(m)
}

var xxx_messageInfo_NFTMintRequest proto.InternalMessageInfo

func (m *NFTMintRequest) GetIdentifier() string {
	if m != nil {
		return m.Identifier
	}
	return ""
}

func (m *NFTMintRequest) GetRegistryAddress() string {
	if m != nil {
		return m.RegistryAddress
	}
	return ""
}

func (m *NFTMintRequest) GetDepositAddress() string {
	if m != nil {
		return m.DepositAddress
	}
	return ""
}

func (m *NFTMintRequest) GetProofFields() []string {
	if m != nil {
		return m.ProofFields
	}
	return nil
}

func (m *NFTMintRequest) GetSubmitTokenProof() bool {
	if m != nil {
		return m.SubmitTokenProof
	}
	return false
}

func (m *NFTMintRequest) GetSubmitNftOwnerAccessProof() bool {
	if m != nil {
		return m.SubmitNftOwnerAccessProof
	}
	return false
}

func (m *NFTMintRequest) GetGrantNftAccess() bool {
	if m != nil {
		return m.GrantNftAccess
	}
	return false
}

type NFTMintResponse struct {
	Header               *ResponseHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	TokenId              string          `protobuf:"bytes,2,opt,name=token_id,json=tokenId,proto3" json:"token_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *NFTMintResponse) Reset()         { *m = NFTMintResponse{} }
func (m *NFTMintResponse) String() string { return proto.CompactTextString(m) }
func (*NFTMintResponse) ProtoMessage()    {}
func (*NFTMintResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_10e4a4ecba67c7da, []int{2}
}
func (m *NFTMintResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NFTMintResponse.Unmarshal(m, b)
}
func (m *NFTMintResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NFTMintResponse.Marshal(b, m, deterministic)
}
func (dst *NFTMintResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NFTMintResponse.Merge(dst, src)
}
func (m *NFTMintResponse) XXX_Size() int {
	return xxx_messageInfo_NFTMintResponse.Size(m)
}
func (m *NFTMintResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_NFTMintResponse.DiscardUnknown(m)
}

var xxx_messageInfo_NFTMintResponse proto.InternalMessageInfo

func (m *NFTMintResponse) GetHeader() *ResponseHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *NFTMintResponse) GetTokenId() string {
	if m != nil {
		return m.TokenId
	}
	return ""
}

func init() {
	proto.RegisterType((*ResponseHeader)(nil), "nft.ResponseHeader")
	proto.RegisterType((*NFTMintRequest)(nil), "nft.NFTMintRequest")
	proto.RegisterType((*NFTMintResponse)(nil), "nft.NFTMintResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// NFTServiceClient is the client API for NFTService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type NFTServiceClient interface {
	MintNFT(ctx context.Context, in *NFTMintRequest, opts ...grpc.CallOption) (*NFTMintResponse, error)
}

type nFTServiceClient struct {
	cc *grpc.ClientConn
}

func NewNFTServiceClient(cc *grpc.ClientConn) NFTServiceClient {
	return &nFTServiceClient{cc}
}

func (c *nFTServiceClient) MintNFT(ctx context.Context, in *NFTMintRequest, opts ...grpc.CallOption) (*NFTMintResponse, error) {
	out := new(NFTMintResponse)
	err := c.cc.Invoke(ctx, "/nft.NFTService/MintNFT", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NFTServiceServer is the server API for NFTService service.
type NFTServiceServer interface {
	MintNFT(context.Context, *NFTMintRequest) (*NFTMintResponse, error)
}

func RegisterNFTServiceServer(s *grpc.Server, srv NFTServiceServer) {
	s.RegisterService(&_NFTService_serviceDesc, srv)
}

func _NFTService_MintNFT_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NFTMintRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NFTServiceServer).MintNFT(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nft.NFTService/MintNFT",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NFTServiceServer).MintNFT(ctx, req.(*NFTMintRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _NFTService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "nft.NFTService",
	HandlerType: (*NFTServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MintNFT",
			Handler:    _NFTService_MintNFT_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "nft/service.proto",
}

func init() { proto.RegisterFile("nft/service.proto", fileDescriptor_service_10e4a4ecba67c7da) }

var fileDescriptor_service_10e4a4ecba67c7da = []byte{
	// 484 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x92, 0xc1, 0x6e, 0x13, 0x3d,
	0x10, 0x80, 0x95, 0xe4, 0x6f, 0x93, 0x3a, 0xfd, 0x93, 0x60, 0x38, 0xa4, 0x11, 0xa0, 0x25, 0x12,
	0x10, 0xa0, 0xcd, 0x4a, 0xe5, 0x80, 0xc4, 0x89, 0x14, 0x14, 0xd1, 0x03, 0x4b, 0xb4, 0xe4, 0x02,
	0x97, 0x95, 0xb3, 0x1e, 0x2f, 0x16, 0xcd, 0x78, 0xb1, 0x27, 0x44, 0x5c, 0x91, 0x78, 0x01, 0x78,
	0x22, 0x9e, 0x81, 0x57, 0xe0, 0x41, 0xd0, 0xda, 0xdb, 0xaa, 0x11, 0xa7, 0xd5, 0x7e, 0xf3, 0xcd,
	0xd8, 0x9e, 0x19, 0x76, 0x03, 0x15, 0xc5, 0x0e, 0xec, 0x17, 0x9d, 0xc3, 0xb4, 0xb4, 0x86, 0x0c,
	0x6f, 0xa1, 0xa2, 0xd1, 0xed, 0xc2, 0x98, 0xe2, 0x02, 0x62, 0x51, 0xea, 0x58, 0x20, 0x1a, 0x12,
	0xa4, 0x0d, 0xba, 0xa0, 0x8c, 0x8e, 0xfd, 0x27, 0x3f, 0x29, 0x00, 0x4f, 0xdc, 0x56, 0x14, 0x05,
	0xd8, 0xd8, 0x94, 0xde, 0xf8, 0xd7, 0x1e, 0x3f, 0x63, 0xbd, 0x14, 0x5c, 0x69, 0xd0, 0xc1, 0x6b,
	0x10, 0x12, 0x2c, 0xbf, 0xcf, 0x7a, 0x64, 0x05, 0x3a, 0x91, 0x57, 0x5e, 0xa6, 0xe5, 0x70, 0x2f,
	0x6a, 0x4c, 0x0e, 0xd2, 0xff, 0xaf, 0xd1, 0x73, 0x39, 0xfe, 0xd5, 0x64, 0xbd, 0x64, 0xbe, 0x7c,
	0xa3, 0x91, 0x52, 0xf8, 0xbc, 0x01, 0x47, 0xfc, 0x2e, 0x63, 0x5a, 0x02, 0x92, 0x56, 0x1a, 0xec,
	0xb0, 0xe1, 0xb3, 0xae, 0x11, 0xfe, 0x88, 0x0d, 0x2c, 0x14, 0xda, 0x91, 0xfd, 0x9a, 0x09, 0x29,
	0x2d, 0x38, 0x37, 0x6c, 0x7a, 0xab, 0x7f, 0xc9, 0x67, 0x01, 0xf3, 0x87, 0xac, 0x2f, 0xa1, 0x34,
	0x4e, 0xd3, 0x95, 0xd9, 0xf2, 0x66, 0xaf, 0xc6, 0x97, 0xe2, 0x3d, 0x76, 0x58, 0x5a, 0x63, 0x54,
	0xa6, 0x34, 0x5c, 0x48, 0x37, 0xfc, 0x2f, 0x6a, 0x4d, 0x0e, 0xd2, 0xae, 0x67, 0x73, 0x8f, 0xf8,
	0x31, 0xe3, 0x6e, 0xb3, 0x5a, 0x6b, 0xca, 0xc8, 0x7c, 0x02, 0xcc, 0x7c, 0xcc, 0x3f, 0xaa, 0x93,
	0x0e, 0x42, 0x64, 0x59, 0x05, 0x16, 0x15, 0xe7, 0x2f, 0xd8, 0x9d, 0xda, 0x46, 0x45, 0x99, 0xd9,
	0x22, 0xd8, 0x4c, 0xe4, 0x39, 0x38, 0x57, 0x27, 0xb6, 0x7d, 0xe2, 0x51, 0x90, 0x12, 0x45, 0x6f,
	0x2b, 0x65, 0xe6, 0x8d, 0x50, 0x61, 0xc2, 0x06, 0x85, 0x15, 0x18, 0x0a, 0x84, 0xd4, 0x61, 0xc7,
	0x27, 0xf5, 0x3c, 0x4f, 0x14, 0x05, 0x7d, 0xfc, 0x9e, 0xf5, 0xaf, 0x5a, 0x18, 0x66, 0xc0, 0x9f,
	0xb0, 0xfd, 0x8f, 0x7e, 0x0e, 0xbe, 0x7f, 0xdd, 0xd3, 0x9b, 0x53, 0x54, 0x34, 0xdd, 0x1d, 0x51,
	0x5a, 0x2b, 0xfc, 0x88, 0x75, 0xc2, 0x93, 0xb4, 0xac, 0x1b, 0xd9, 0xf6, 0xff, 0xe7, 0xf2, 0xf4,
	0x7b, 0x83, 0xb1, 0x64, 0xbe, 0x7c, 0x17, 0xb6, 0x87, 0x6f, 0x59, 0xbb, 0x3a, 0x26, 0x99, 0x2f,
	0x79, 0xa8, 0xb8, 0x3b, 0xba, 0xd1, 0xad, 0x5d, 0x18, 0x4e, 0x1b, 0xcf, 0x7e, 0xcc, 0x26, 0xa3,
	0x07, 0x15, 0x8a, 0x04, 0x46, 0xc9, 0x7c, 0x19, 0x29, 0x6b, 0xd6, 0x91, 0x88, 0x5e, 0x02, 0x92,
	0xd5, 0x6a, 0x53, 0x40, 0xf4, 0xca, 0xe4, 0x9b, 0x35, 0x20, 0x7d, 0xfb, 0xfd, 0xe7, 0x67, 0x73,
	0x30, 0xee, 0xc6, 0xfe, 0x06, 0xf1, 0x5a, 0x23, 0x3d, 0x6f, 0x3c, 0x3e, 0x8b, 0x58, 0x3b, 0x37,
	0xeb, 0xaa, 0xfa, 0xd9, 0x61, 0x7d, 0x99, 0x45, 0xb5, 0x78, 0x8b, 0xc6, 0x87, 0x3d, 0x54, 0x54,
	0xae, 0x56, 0xfb, 0x7e, 0x11, 0x9f, 0xfe, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x7e, 0xdb, 0x53, 0x5a,
	0xee, 0x02, 0x00, 0x00,
}
