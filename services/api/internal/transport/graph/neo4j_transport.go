package graph

import (
	"context"
	"errors"
	"commlayers/services/api/internal/transport"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

var (
	ErrNoNeo4jDriver = errors.New("neo4j driver not configured")
)

// Neo4jTransport implements TransportProvider using Neo4j.
type Neo4jTransport struct {
	driver neo4j.Driver
}

// New creates a Neo4jTransport and connects to Neo4j.
func New(uri, username, password string) (*Neo4jTransport, error) {
	if uri == "" {
		return &Neo4jTransport{}, nil
	}
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, err
	}
	return &Neo4jTransport{driver: driver}, nil
}

func (n *Neo4jTransport) Create(ctx context.Context, t *transport.Transport) error {
	if n.driver == nil {
		return ErrNoNeo4jDriver
	}
	session := n.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	_, err := session.Run(
		"CREATE (t:Transport {id: $id, name: $name, description: $description, tags: $tags})",
		map[string]interface{}{"id": t.ID, "name": t.Name, "description": t.Description, "tags": t.Tags},
		nil,
	)
	return err
}

func (n *Neo4jTransport) GetByID(ctx context.Context, id string) (*transport.Transport, error) {
	if n.driver == nil {
		return nil, ErrNoNeo4jDriver
	}
	session := n.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	result, err := session.Run(
		"MATCH (t:Transport {id: $id}) RETURN t.id, t.name, t.description, t.tags",
		map[string]interface{}{"id": id},
		nil,
	)
	if err != nil {
		return nil, err
	}
	record, err := result.Single()
	if err != nil {
		return nil, err
	}
	values := record.Values
	if len(values) < 4 {
		return nil, errors.New("incomplete record")
	}
	return &transport.Transport{
		ID:          values[0].(string),
		Name:        values[1].(string),
		Description: values[2].(string),
		Tags:        values[3].([]string),
	}, nil
}

func (n *Neo4jTransport) Search(ctx context.Context, query string) ([]*transport.Transport, error) {
	if n.driver == nil {
		return nil, ErrNoNeo4jDriver
	}
	session := n.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	result, err := session.Run(
		"MATCH (t:Transport) WHERE t.name CONTAINS $q OR t.description CONTAINS $q RETURN t.id, t.name, t.description, t.tags",
		map[string]interface{}{"q": query},
		nil,
	)
	if err != nil {
		return nil, err
	}
	var transports []*transport.Transport
	for result.Next() {
		record := result.Record()
		values := record.Values
		if len(values) < 4 {
			continue
		}
		transports = append(transports, &transport.Transport{
			ID:          values[0].(string),
			Name:        values[1].(string),
			Description: values[2].(string),
			Tags:        values[3].([]string),
		})
	}
	if err := result.Err(); err != nil {
		return nil, err
	}
	return transports, nil
}