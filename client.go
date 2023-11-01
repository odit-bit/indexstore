package indexstore

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/odit-bit/indexstore/index"
	"github.com/odit-bit/indexstore/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IndexClient struct {
	ctx context.Context
	rpc proto.IndexerClient
}

func NewClient(ctx context.Context, conn grpc.ClientConnInterface) *IndexClient {
	cli := proto.NewIndexerClient(conn)
	idx := IndexClient{
		ctx: ctx,
		rpc: cli,
	}

	return &idx
}

// Index implements index.Indexer.
func (idx *IndexClient) Index(doc *index.Document) error {
	idxCtx, cancel := context.WithCancel(idx.ctx)
	defer cancel()

	rpcDoc := proto.Document{
		LinkId:    doc.LinkID[:],
		Url:       doc.URL,
		Title:     doc.Title,
		Content:   doc.Content,
		IndexedAt: timestamppb.New(doc.IndexedAt),
	}
	res, err := idx.rpc.Index(idxCtx, &rpcDoc)
	if err != nil {
		return err
	}

	doc.IndexedAt = res.IndexedAt.AsTime()
	return nil
}

// UpdateRank implements index.Indexer.
func (idx *IndexClient) UpdateRank(linkID uuid.UUID, score float64) error {
	idxCtx, cancel := context.WithCancel(idx.ctx)
	defer cancel()

	req := proto.UpdateScoreRequest{
		LinkId:        linkID[:],
		PageRankScore: score,
	}
	_, err := idx.rpc.UpdateScore(idxCtx, &req)
	if err != nil {
		return err
	}
	return nil
}

// Search implements index.Indexer.
func (idx *IndexClient) Search(query index.Query) (index.Iterator, error) {
	idxCtx, cancel := context.WithCancel(idx.ctx)

	rpcQuery := proto.Query{
		Type:       proto.Type(query.Type),
		Expression: query.Expression,
		Offset:     query.Offset,
	}
	rpcIter, err := idx.rpc.Search(idxCtx, &rpcQuery)
	if err != nil {
		cancel()
		return nil, err
	}

	res, err := rpcIter.Recv()
	if err != nil {
		cancel()
		return nil, err
	}

	if res.GetDoc() != nil {
		cancel()
		return nil, fmt.Errorf("expected server send doc total count before any document ")
	}

	iter := searchIterator{
		totalCount:  res.GetDocCount(),
		rpcIterator: rpcIter,
		doc:         nil,
		lastErr:     nil,
		cancelFn:    cancel,
	}

	return &iter, nil
}

var _ index.Iterator = (*searchIterator)(nil)

type searchIterator struct {
	totalCount  uint64
	rpcIterator proto.Indexer_SearchClient

	// last received document
	doc *index.Document

	lastErr error

	//
	cancelFn context.CancelFunc
}

// Close implements index.Iterator.
func (si *searchIterator) Close() error {
	si.cancelFn()
	return si.rpcIterator.CloseSend()
}

// Document implements index.Iterator.
func (si *searchIterator) Document() *index.Document {
	return si.doc
}

// Error implements index.Iterator.
func (si *searchIterator) Error() error {

	return si.lastErr
}

// Next implements index.Iterator.
func (si *searchIterator) Next() bool {
	res, err := si.rpcIterator.Recv()
	if err != nil {
		si.lastErr = err
		si.cancelFn()
		return false
	}

	rpcDoc := res.GetDoc()
	doc := &index.Document{
		LinkID:    uuid.UUID(rpcDoc.GetLinkId()),
		URL:       rpcDoc.Url,
		Title:     rpcDoc.Title,
		Content:   rpcDoc.Content,
		IndexedAt: rpcDoc.IndexedAt.AsTime(),
	}

	si.doc = doc
	return true
}

// TotalCount implements index.Iterator.
func (si *searchIterator) TotalCount() uint64 {
	return si.totalCount
}
