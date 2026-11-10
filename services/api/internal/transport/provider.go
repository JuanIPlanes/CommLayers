// Package transport defines abstractions for the transport‑family data stores.
//
// The transport‑family includes graph‑based (Neo4j), vector‑based (Qdrant),
// and document‑based (Mongo/Elasticsearch) back‑ends. Each provider must be
// able to persist a Transport entity and perform queries required by the API.
//
// This file introduces a minimal interface that concrete implementations will
// satisfy. It is deliberately small – the concrete providers will add their own
// methods as needed.

package transport

import "context"

// Transport represents the domain model for a transport entity. The fields are
// intentionally generic; concrete providers may store additional metadata.
type Transport struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    Description string   `json:"description,omitempty"`
    Tags        []string `json:"tags,omitempty"`
}

// TransportProvider abstracts the operations required by the API layer.
// Implementations may back the data store with Neo4j, Qdrant, Elasticsearch, etc.
type TransportProvider interface {
    // Create stores a new Transport entity.
    Create(ctx context.Context, t *Transport) error
    // GetByID retrieves a Transport by its ID.
    GetByID(ctx context.Context, id string) (*Transport, error)
    // Search performs a generic search; concrete providers may interpret the
    // query differently (e.g., graph traversal, vector similarity, full‑text).
    Search(ctx context.Context, query string) ([]*Transport, error)
}
