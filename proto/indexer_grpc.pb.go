// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: proto/indexer.proto

package proto

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// IndexerClient is the client API for Indexer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IndexerClient interface {
	Index(ctx context.Context, in *Document, opts ...grpc.CallOption) (*Document, error)
	// Search the index for a particular query and stream the results back to
	// the client. The first response will include the total result count while
	// all subsequent responses will include documents from the resultset.
	Search(ctx context.Context, in *Query, opts ...grpc.CallOption) (Indexer_SearchClient, error)
	// UpdateScore updates the PageRank score for a document with the specified
	// link ID.
	UpdateScore(ctx context.Context, in *UpdateScoreRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type indexerClient struct {
	cc grpc.ClientConnInterface
}

func NewIndexerClient(cc grpc.ClientConnInterface) IndexerClient {
	return &indexerClient{cc}
}

func (c *indexerClient) Index(ctx context.Context, in *Document, opts ...grpc.CallOption) (*Document, error) {
	out := new(Document)
	err := c.cc.Invoke(ctx, "/proto.Indexer/Index", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexerClient) Search(ctx context.Context, in *Query, opts ...grpc.CallOption) (Indexer_SearchClient, error) {
	stream, err := c.cc.NewStream(ctx, &Indexer_ServiceDesc.Streams[0], "/proto.Indexer/Search", opts...)
	if err != nil {
		return nil, err
	}
	x := &indexerSearchClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Indexer_SearchClient interface {
	Recv() (*QueryResult, error)
	grpc.ClientStream
}

type indexerSearchClient struct {
	grpc.ClientStream
}

func (x *indexerSearchClient) Recv() (*QueryResult, error) {
	m := new(QueryResult)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *indexerClient) UpdateScore(ctx context.Context, in *UpdateScoreRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/proto.Indexer/UpdateScore", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IndexerServer is the server API for Indexer service.
// All implementations must embed UnimplementedIndexerServer
// for forward compatibility
type IndexerServer interface {
	Index(context.Context, *Document) (*Document, error)
	// Search the index for a particular query and stream the results back to
	// the client. The first response will include the total result count while
	// all subsequent responses will include documents from the resultset.
	Search(*Query, Indexer_SearchServer) error
	// UpdateScore updates the PageRank score for a document with the specified
	// link ID.
	UpdateScore(context.Context, *UpdateScoreRequest) (*empty.Empty, error)
	mustEmbedUnimplementedIndexerServer()
}

// UnimplementedIndexerServer must be embedded to have forward compatible implementations.
type UnimplementedIndexerServer struct {
}

func (UnimplementedIndexerServer) Index(context.Context, *Document) (*Document, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Index not implemented")
}
func (UnimplementedIndexerServer) Search(*Query, Indexer_SearchServer) error {
	return status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (UnimplementedIndexerServer) UpdateScore(context.Context, *UpdateScoreRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateScore not implemented")
}
func (UnimplementedIndexerServer) mustEmbedUnimplementedIndexerServer() {}

// UnsafeIndexerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IndexerServer will
// result in compilation errors.
type UnsafeIndexerServer interface {
	mustEmbedUnimplementedIndexerServer()
}

func RegisterIndexerServer(s grpc.ServiceRegistrar, srv IndexerServer) {
	s.RegisterService(&Indexer_ServiceDesc, srv)
}

func _Indexer_Index_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Document)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexerServer).Index(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Indexer/Index",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexerServer).Index(ctx, req.(*Document))
	}
	return interceptor(ctx, in, info, handler)
}

func _Indexer_Search_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Query)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(IndexerServer).Search(m, &indexerSearchServer{stream})
}

type Indexer_SearchServer interface {
	Send(*QueryResult) error
	grpc.ServerStream
}

type indexerSearchServer struct {
	grpc.ServerStream
}

func (x *indexerSearchServer) Send(m *QueryResult) error {
	return x.ServerStream.SendMsg(m)
}

func _Indexer_UpdateScore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateScoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexerServer).UpdateScore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Indexer/UpdateScore",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexerServer).UpdateScore(ctx, req.(*UpdateScoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Indexer_ServiceDesc is the grpc.ServiceDesc for Indexer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Indexer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Indexer",
	HandlerType: (*IndexerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Index",
			Handler:    _Indexer_Index_Handler,
		},
		{
			MethodName: "UpdateScore",
			Handler:    _Indexer_UpdateScore_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Search",
			Handler:       _Indexer_Search_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/indexer.proto",
}