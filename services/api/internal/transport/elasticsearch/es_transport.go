package elasticsearch

import (
	"context"
	"errors"
	"commlayers/services/api/internal/transport"
	"github.com/elastic/go-elasticsearch/v8"
)

var (
	ErrNoEsClient = errors.New("elasticsearch client not configured")
)

// EsTransport implements TransportProvider using Elasticsearch.
type EsTransport struct {
	client *elasticsearch.Client
}

// New creates an EsTransport and connects to Elasticsearch.
func New(uri string) (*EsTransport, error) {
	if uri == "" {
		return &EsTransport{}, nil
	}
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, err
	}
	return &EsTransport{client: client}, nil
}

func (e *EsTransport) Create(ctx context.Context, t *transport.Transport) error {
	if e.client == nil {
		return ErrNoEsClient
	}
	// TODO: Index transport document into ES
	// ES API: PUT /{index}/_doc/{id} with transport payload
	// For now return nil since client is configured
	return nil
}

func (e *EsTransport) GetByID(ctx context.Context, id string) (*transport.Transport, error) {
	if e.client == nil {
		return nil, ErrNoEsClient
	}
	// TODO: GET /{index}/_doc/{id} from ES
	// For now return nil since client is configured
	return nil, nil
}

func (e *EsTransport) Search(ctx context.Context, query string) ([]*transport.Transport, error) {
	if e.client == nil {
		return nil, ErrNoEsClient
	}
	// TODO: POST /_search to ES with query matching transport fields
	// For now return empty results
	return nil, nil
}