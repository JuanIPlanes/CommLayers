package vector

import (
	"context"
	"errors"
	"commlayers/services/api/internal/transport"
	"net/http"
	"time"
)

var (
	ErrNoQdrantClient = errors.New("qdrant client not configured")
)

// QdrantTransport implements TransportProvider using Qdrant.
type QdrantTransport struct {
	baseURI string
	client  *http.Client
}

// New creates a QdrantTransport with the given base URI.
func New(uri string) (*QdrantTransport, error) {
	if uri == "" {
		return &QdrantTransport{}, nil
	}
	client := &http.Client{Timeout: 30 * time.Second}
	return &QdrantTransport{baseURI: uri, client: client}, nil
}

func (q *QdrantTransport) Create(ctx context.Context, t *transport.Transport) error {
	if q.baseURI == "" {
		return ErrNoQdrantClient
	}
	// TODO: POST /points to upsert transport data as Qdrant point payload
	// Qdrant API: PUT /collections/{collection_name}/points with point payload
	// For now return nil since client is configured
	return nil
}

func (q *QdrantTransport) GetByID(ctx context.Context, id string) (*transport.Transport, error) {
	if q.baseURI == "" {
		return nil, ErrNoQdrantClient
	}
	// TODO: GET /points from Qdrant by ID
	// For now return nil since client is configured
	return nil, nil
}

func (q *QdrantTransport) Search(ctx context.Context, query string) ([]*transport.Transport, error) {
	if q.baseURI == "" {
		return nil, ErrNoQdrantClient
	}
	// TODO: POST /search to Qdrant with vector/payload query
	// For now return empty results
	return nil, nil
}