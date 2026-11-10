package evidence

import (
	"encoding/json"
	"time"
)

// EvidenceRecord represents a canonical truth record stored in the PostgreSQL-backed
// persistence layer. Each record has a kind, unique ID, and a JSON payload that can
// hold family-specific data (transport, messaging, sync, etc.).
type EvidenceRecord struct {
	Kind       string    `json:"kind"`
	ID         string    `json:"id"`
	Payload    json.RawMessage `json:"payload,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ComparisonFamily identifies which communication family an evidence record belongs to.
type ComparisonFamily string

const (
	// FamilyTransport represents the transport comparison family.
	FamilyTransport ComparisonFamily = "transport"
	// FamilyMessaging represents the messaging and workflow family.
	FamilyMessaging ComparisonFamily = "messaging"
	// FamilySync represents the sync and replication family.
	FamilySync ComparisonFamily = "sync"
	// FamilySearch represents the search, graph, and vector family.
	FamilySearch ComparisonFamily = "search"
	// FamilyBenchmark represents the benchmark and observability family.
	FamilyBenchmark ComparisonFamily = "benchmark"
)

// ComparisonContext holds the context for comparing evidence across families.
type ComparisonContext struct {
	// The family being compared.
	Family ComparisonFamily
	// Optional filter or scope for the comparison.
	Scope string
	// Timestamp indicates when the comparison was performed.
	Timestamp time.Time
}

// ComparisonResult represents the result of comparing evidence across families.
type ComparisonResult struct {
	// The family that was compared.
	Family ComparisonFamily
	// Evidence records involved in the comparison.
	Evidence []EvidenceRecord
	// Similarity score or decision metric.
	Score float64
	// Decision based on the comparison (e.g., "preferred", "alternative", "avoid").
	Decision string
	// Rationale for the decision.
	Rationale string
}

// FamilyConfig holds configuration specific to a communication family.
type FamilyConfig struct {
	// Family name.
	Name ComparisonFamily
	// Priority or ordering within the overall system.
	Priority int
	// Recommended when to use this family.
	RecommendedWhen []string
	// When to avoid this family.
	AvoidWhen []string
	// Additional notes.
	Notes []string
}

// EvidenceQuery describes a query for retrieving evidence records.
type EvidenceQuery struct {
	// Kind filter.
	Kind string
	// Limit on results.
	Limit int
	// Order by field.
	OrderBy string
	// Sort direction (asc/desc).
	SortDir string
}