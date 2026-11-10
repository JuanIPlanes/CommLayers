package vector

import (
	"testing"
	"commlayers/services/api/internal/transport"
)

func TestQdrantTransportImplementsProvider(t *testing.T) {
	var _ transport.TransportProvider = &QdrantTransport{}
}