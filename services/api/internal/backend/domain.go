package backend

import (
	"encoding/json"
	"time"
)

// CanonicalPostgresDomainModel represents the PostgreSQL-first system of record,
// shared evidence model, and core domain schema for comparative backend flows.
//
// This model defines the canonical truth layer that all communication families
// (transport, messaging, sync, search, benchmark) read from and write to.
type CanonicalPostgresDomainModel struct {
	// Kind categorizes the type of data (transport, messaging, sync, search, benchmark).
	Kind string `json:"kind"`
	// ID is the unique identifier within the kind.
	ID string `json:"id"`
	// Payload is the JSON-encoded domain-specific data.
	Payload json.RawMessage `json:"payload,omitempty"`
	// CreatedAt is when the record was first persisted.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the record was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// DomainEvent represents a domain event that can be published across communication families.
type DomainEvent struct {
	// EventType identifies the kind of domain event.
	EventType string `json:"event_type"`
	// Source identifies which family originated the event.
	Source string `json:"source"`
	// AggregateID is the ID of the aggregate the event pertains to.
	AggregateID string `json:"aggregate_id"`
	// Payload is the event-specific data.
	Payload json.RawMessage `json:"payload,omitempty"`
	// OccurredAt is when the event occurred.
	OccurredAt time.Time `json:"occurred_at"`
}

// FamilyComparativeData holds the comparative data structure for backend family analysis.
type FamilyComparativeData struct {
	// Family is the communication family (transport, messaging, sync, search, benchmark).
	Family string `json:"family"`
	// Metrics holds family-specific metrics for comparison.
	Metrics map[string]float64 `json:"metrics,omitempty"`
	// Status holds the current status of the family comparison.
	Status string `json:"status"`
	// Decision holds the recommended decision based on the comparison.
	Decision string `json:"decision,omitempty"`
	// Rationale holds the rationale for the decision.
	Rationale string `json:"rationale,omitempty"`
}

// BackendFamily identifies which communication family a domain model belongs to.
type BackendFamily string

const (
	// BackendFamilyTransport represents the transport comparison family.
	BackendFamilyTransport BackendFamily = "transport"
	// BackendFamilyMessaging represents the messaging and workflow family.
	BackendFamilyMessaging BackendFamily = "messaging"
	// BackendFamilySync represents the sync and replication family.
	BackendFamilySync BackendFamily = "sync"
	// BackendFamilySearch represents the search, graph, and vector family.
	BackendFamilySearch BackendFamily = "search"
	// BackendFamilyBenchmark represents the benchmark and observability family.
	BackendFamilyBenchmark BackendFamily = "benchmark"
)

// FamilyComparisonResult holds the result of comparing backends across families.
type FamilyComparisonResult struct {
	// ComparedFamilies lists the families that were compared.
	ComparedFamilies []BackendFamily `json:"compared_families"`
	// Scores holds the comparison scores per family.
	Scores map[BackendFamily]float64 `json:"scores"`
	// Recommendations holds recommendations per family.
	Recommendations map[BackendFamily]string `json:"recommendations"`
	// OverallDecision holds the overall decision across all families.
	OverallDecision string `json:"overall_decision"`
	// Timestamp when the comparison was performed.
	Timestamp time.Time `json:"timestamp"`
}

// EvidenceRecordWithKind is an EvidenceRecord with an explicit Kind field
// for type-safe domain-driven operations.
type EvidenceRecordWithKind struct {
	Kind       string    `json:"kind"`
	ID         string    `json:"id"`
	Payload    json.RawMessage `json:"payload,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	// Additional domain-specific fields can be added per family
}

// InitDomainModels initializes the default domain model set.
// This is called during application bootstrap to establish the canonical model.
func InitDomainModels() {
	// Domain model initialization is handled by the persistence layer
	// via the Store's Init method which creates the app_records table.
}

// DomainModelOption represents an option for configuring a domain model.
type DomainModelOption func(*CanonicalPostgresDomainModel)

// WithKind sets the kind for a domain model.
func WithKind(kind string) DomainModelOption {
	return func(m *CanonicalPostgresDomainModel) {
		m.Kind = kind
	}
}

// WithPayload sets the payload for a domain model.
func WithPayload(payload json.RawMessage) DomainModelOption {
	return func(m *CanonicalPostgresDomainModel) {
		m.Payload = payload
	}
}

// WithTimestamps sets the timestamps for a domain model.
func WithTimestamps(createdAt, updatedAt time.Time) DomainModelOption {
	return func(m *CanonicalPostgresDomainModel) {
		m.CreatedAt = createdAt
		m.UpdatedAt = updatedAt
	}
}