package indexstore

import (
	"context"

	"github.com/google/uuid"
	"github.com/odit-bit/indexstore/index"
	"github.com/odit-bit/indexstore/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ proto.IndexerServer = (*IndexServer)(nil)

type IndexServer struct {
	db index.Indexer

	proto.UnimplementedIndexerServer
}

func NewIndexServer(indexDB index.Indexer) *IndexServer {
	is := IndexServer{
		db:                         indexDB,
		UnimplementedIndexerServer: proto.UnimplementedIndexerServer{},
	}

	return &is
}

// Index implements proto.IndexerServer.
func (idx *IndexServer) Index(ctx context.Context, req *proto.Document) (*proto.Document, error) {
	doc := index.Document{
		LinkID:    uuid.UUID(req.GetLinkId()),
		URL:       req.GetUrl(),
		Title:     req.GetTitle(),
		Content:   req.GetContent(),
		IndexedAt: req.GetIndexedAt().AsTime(),
	}
	err := idx.db.Index(&doc)
	if err != nil {
		return nil, err
	}

	req.IndexedAt = timestamppb.New(doc.IndexedAt)
	return req, nil
}

// Search implements proto.IndexerServer.
func (idx *IndexServer) Search(query *proto.Query, res proto.Indexer_SearchServer) error {
	q := index.Query{
		Type:       index.QueryType(query.GetType()),
		Expression: query.GetExpression(),
		Offset:     query.GetOffset(),
	}

	iter, err := idx.db.Search(q)
	if err != nil {
		return err
	}
	defer iter.Close()

	totalCount := proto.QueryResult{
		Result: &proto.QueryResult_DocCount{
			DocCount: iter.TotalCount(),
		},
	}
	if err := res.Send(&totalCount); err != nil {
		return err
	}

	for iter.Next() {
		doc := iter.Document()
		qRes := proto.QueryResult{
			Result: &proto.QueryResult_Doc{Doc: &proto.Document{
				LinkId:    doc.LinkID[:],
				Url:       doc.URL,
				Title:     doc.Title,
				Content:   doc.Content,
				IndexedAt: timestamppb.New(doc.IndexedAt),
			}},
		}
		res.Send(&qRes)
	}

	return iter.Error()
}

// UpdateScore implements proto.IndexerServer.
func (idx *IndexServer) UpdateScore(ctx context.Context, req *proto.UpdateScoreRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, idx.db.UpdateRank(uuid.UUID(req.GetLinkId()), req.GetPageRankScore())
}
