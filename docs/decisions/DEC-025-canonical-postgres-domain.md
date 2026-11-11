---
id: DEC-025
title: Canonical Postgres Domain Model
status: proposed
date: 2026-11-11
---

## Context
The `store.go` and `domain.go` define the PostgreSQL-first system of record with entity catalog, ownership matrix, schema evolution, and audit lineage. This decision formalizes the canonical domain model as a binding contract.

## Decision
### Entity Catalog (from `CanonicalPostgresDomainModel` and `store.go`)
All entities stored in single `app_records` table with `(kind, id)` primary key.

| Kind | Description | Ownership | Schema Version | Audit Lineage |
|------|-------------|-----------|----------------|---------------|
| `stream_definitions` | First-wave + deferred stream configs | Stream lead | v1 | Full |
| `comparison_evidence` | Benchmark results, family comparisons | Benchmark owner | v1 | Full |
| `job_lifecycle` | Workflow runs, async jobs, transport sessions | Job owner | v1 | Full |
| `deferred_wave_manifest` | Wave definitions, reactivation rules | Wave coordinator | v1 | Full |
| `domain_event` | Event-sourced commands for v2 paradigms | Event producer | v1 | Full |
| `sync_session` | Sync replication sessions | Sync owner | v1 | Full |
| `workflow_run` | Workflow execution records | Workflow owner | v1 | Full |
| `transport_session` | Polling/long-polling/SSE/WebSocket sessions | Transport owner | v1 | Full |
| `projection_freshness` | Search/Graph/Vector freshness state | Projection owner | v1 | Full |
| `pricing_rate` | USD/XXX rates with metadata | Pricing owner | v1 | Full |

### Storage Ownership Matrix (from `handleDataPlatform`)
| Store | Role | Status | Entities Owned |
|-------|------|--------|----------------|
| **PostgreSQL** | System of Record | Active (compose) | ALL kinds above |
| **Redis** | Ephemeral Speed Layer | Active (compose) | Cache entries, rate limits, locks, transient events |
| **Elasticsearch** | Search Projection | Optional | `projection_search` documents |
| **Neo4j** | Graph Projection | Optional | `projection_graph` nodes/edges |
| **Qdrant** | Vector Projection | Optional | `projection_vector` vectors |
| **MongoDB** | Document Projection | Optional | Only if explicitly justified |

**Ownership Rules**:
1. PostgreSQL is the **only** system of record — all canonical truth lives here
2. Redis data MUST be rebuildable from PostgreSQL
3. Optional profiles (ES/Neo4j/Qdrant/Mongo) are projections ONLY — never source of truth
4. No entity kind may have dual ownership

### Schema Evolution Policy (from `domain.go` and `store.go`)
**Rules**:
1. **Explicit Ownership**: Every `kind` has a declared owner (team/individual)
2. **Versioned Evolution**: Schema changes via migration scripts with version numbers
3. **Audit Lineage Required**: Every record carries `created_at`, `updated_at`, plus evidence fields
4. **Backward Compatibility**: Migrations must not break existing readers (additive only)
5. **Deprecation Process**: 
   - Mark kind deprecated in metadata
   - 90-day notice before removal
   - Migration path documented

**Migration Contract**:
```go
type SchemaMigration struct {
    Version     int       // Monotonic increasing
    Kind        string    // Affected kind
    Description string    // Human-readable
    UpSQL       string    // Forward migration
    DownSQL     string    // Rollback (required)
    Owner       string    // Approving owner
    AppliedAt   time.Time // When applied
}
```

**Current Schema Version**: 1 (initial `app_records` table)

### Audit Lineage Contract (from `handleDataPlatform` and `evidence/contracts.go`)
Every record in `app_records` implicitly carries audit fields via `EvidenceContract`:
```json
{
  "evidence": {
    "source": "benchmark|audit|projection",
    "confidence": 0.95,
    "staleness": "fresh|stale|very_stale|reindexing",
    "measured_at": "2026-11-11T10:30:00Z",
    "delay_applied_ms": 120
  },
  "success": true,
  "success_code": "ok|partial|degraded",
  "error_code": "timeout|unavailable|...",
  "error_message": "..."
}
```

**Audit Fields Required on ALL Kinds**:
| Field | Source | Description |
|-------|--------|-------------|
| `source` | `EvidenceSource` | Origin of data |
| `confidence` | float64 (0-1) | Reliability score |
| `staleness` | `EvidenceStaleness` | Freshness classification |
| `measured_at` | RFC3339 timestamp | When data was measured |
| `delay_applied_ms` | int64 | Processing latency |

### Compatibility Rules Across Backend Families (from `domain.go`)
| Family | PostgreSQL | Redis | Elasticsearch | Neo4j | Qdrant |
|--------|------------|-------|---------------|-------|--------|
| Transport | System of record | Cache/sessions | Search projection | Graph projection | Vector projection |
| Messaging | System of record | Queue state | Search projection | Graph projection | — |
| Sync | System of record | Lock coordination | — | Graph projection (topology) | — |
| Search | System of record | Cache | **Primary** | — | — |
| Benchmark | System of record | Transient metrics | — | — | — |

**Compatibility Guarantees**:
1. All families read/write PostgreSQL via `Store` interface
2. Redis operations via `Cache` interface (cache-aside pattern)
3. Optional profiles implement `TransportProvider` interface
4. No family bypasses PostgreSQL for canonical data

### Domain Event Contract (from `domain.go` and `handleV2EventDriven`)
```go
type DomainEvent struct {
    EventType   string          `json:"event_type"`   // e.g., "workflow.run.created"
    Source      string          `json:"source"`       // e.g., "workflow_engine"
    AggregateID string          `json:"aggregate_id"` // e.g., "workflow-123"
    Payload     json.RawMessage `json:"payload"`
    OccurredAt  time.Time       `json:"occurred_at"`
}
```

**Event Types** (canonical):
- `stream.definition.created|updated`
- `workflow.run.created|stage_changed|completed|failed`
- `sync.session.created|synced|conflict|resolved`
- `transport.session.created|step_changed|completed`
- `async.job.created|started|completed|failed`
- `projection.freshness.changed|reindex.triggered|reindex.completed`
- `pricing.rate.updated|verified|fallback`

**Event Store**: PostgreSQL `app_records` with `kind = "domain_event"`

### API Contract
**`GET /api/data-platform`** returns ownership matrix, optional profiles, contracts
**`GET /api/projections/freshness`** returns freshness for all projection types
**`GET /api/events`** returns recent domain events (transient via Redis pubsub)

## Consequences
- Single `app_records` table is the canonical schema — no parallel tables
- Ownership matrix is binding: no entity without declared owner
- Schema migrations are the only way to evolve — no ad-hoc ALTER TABLE
- Audit lineage is mandatory: every write produces evidence
- Optional profiles are projections only — never authoritative

## Acceptance Criteria
- [ ] All 10 entity kinds defined with ownership
- [ ] Storage ownership matrix matches `handleDataPlatform`
- [ ] Schema evolution policy enforced: migrations required, backward compatible
- [ ] Audit lineage fields present on all records (verified by grep)
- [ ] Compatibility matrix matches family-to-store mapping
- [ ] Domain event types canonical and versioned
- [ ] API endpoints return correct contract data