package graph

import (
	"testing"
	"commlayers/services/api/internal/transport"
)

func TestNeo4jTransportImplementsProvider(t *testing.T) {
	var _ transport.TransportProvider = &Neo4jTransport{}
}