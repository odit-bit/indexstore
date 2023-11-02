package indexstore

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/odit-bit/indexstore/index"
)

func TestXxx(t *testing.T) {
	srv := Server{
		Port:    6969,
		Handler: &mockDB{},
	}
	go srv.ListenAndServe()

	client, err := ConnectIndex("localhost:6969")
	if err != nil {
		t.Fatal("failed connect to server")
	}

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
