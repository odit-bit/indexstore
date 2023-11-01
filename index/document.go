package index

import (
	"time"

	"github.com/google/uuid"
)

// represent text data of html document
type Document struct {
	LinkID  uuid.UUID
	URL     string
	Title   string
	Content string

	IndexedAt time.Time
	Pagerank  float64
}

//

type Indexer interface {
	// index will upsert the doc
	Index(doc *Document) error

	//find document by ID
	Find(linkID uuid.UUID) (*Document, error)

	// perform full-text query search
	Search(query Query) (Iterator, error)

	//update pagerank score for particular document
	UpdateRank(linkIID uuid.UUID, score float64) error
}

type QueryType int

const (
	QueryTypeMatch = iota
	QueryTypePhrase
)

// represent full-text query for search
type Query struct {
	Type QueryType

	Expression string

	Offset uint64
}

// implement by object that can paginated the result
type Iterator interface {
	Close() error

	Next() bool

	Error() error

	Document() *Document

	TotalCount() uint64
}
