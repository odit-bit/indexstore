package indexstore

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/odit-bit/indexstore/index"
	"github.com/odit-bit/indexstore/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestXxx(t *testing.T) {
	listen, _ := net.Listen("tcp", "localhost:6969")
	srv := setupGrpcServer()
	defer srv.Stop()
	go srv.Serve(listen)

	client := createGRPCClient("localhost:6969")
	res, err := client.Search(index.Query{
		Type:       0,
		Expression: "",
		Offset:     0,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer res.Close()

	if res.TotalCount() != 1 {
		t.Fatal("got total count ", res.TotalCount())
	}

	for res.Next() {
		doc := res.Document()
		if doc.Title != "title-1" {
			t.Fatal("doc error")
		}
	}

	if res.Error() != nil {
		t.Fatal("got error when close", res.Error())
	}
}

func createGRPCClient(address string) *IndexClient {
	// Create a gRPC client to connect to the server.
	// Return the client.

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	return NewClient(context.Background(), conn)
}

func setupGrpcServer() *grpc.Server {
	// Implement your gRPC server setup here.
	// Return a running gRPC server.

	mock := mockDB{}
	srv := NewServer(&mock)
	gSrv := grpc.NewServer()
	proto.RegisterIndexerServer(gSrv, srv)

	return gSrv
}

var _ index.Indexer = (*mockDB)(nil)

type mockDB struct {
}

// Find implements index.Indexer.
func (*mockDB) Find(linkID uuid.UUID) (*index.Document, error) {
	panic("unimplemented")
}

// Index implements index.Indexer.
func (*mockDB) Index(doc *index.Document) error {
	panic("unimplemented")
}

// Search implements index.Indexer.
func (m *mockDB) Search(query index.Query) (index.Iterator, error) {
	idxIter := mockIdxIter{
		ds: []*index.Document{{
			LinkID:    uuid.New(),
			URL:       "www.example1.com",
			Title:     "title-1",
			Content:   "content-1",
			IndexedAt: time.Time{},
			Pagerank:  0,
		}},
		idx: 0,
	}

	return &idxIter, nil
}

var _ index.Iterator = (*mockIdxIter)(nil)

type mockIdxIter struct {
	ds  []*index.Document
	idx int
}

// Close implements index.Iterator.
func (*mockIdxIter) Close() error {
	return nil
}

// Document implements index.Iterator.
func (m *mockIdxIter) Document() *index.Document {
	doc := m.ds[m.idx]
	m.idx++
	return doc
}

// Error implements index.Iterator.
func (*mockIdxIter) Error() error {
	return nil
}

// Next implements index.Iterator.
func (m *mockIdxIter) Next() bool {
	return m.idx < len(m.ds)
}

// TotalCount implements index.Iterator.
func (m *mockIdxIter) TotalCount() uint64 {
	return uint64(len(m.ds))
}

// UpdateRank implements index.Indexer.
func (*mockDB) UpdateRank(linkIID uuid.UUID, score float64) error {
	panic("unimplemented")
}
