---
id: DEC-015
title: Elasticsearch Projection Model
status: proposed
date: 2026-11-11
---

## Context
The `es_transport.go` defines the Elasticsearch transport provider for search projections. This decision formalizes the projection model contract: index mapping, query DSL, freshness integration, and evidence contract integration.

## Decision
### Index Mapping Contract
**Index Name Pattern**: `commlayers-{projection-type}-{environment}` (e.g., `commlayers-search-prod`)
**Document Structure**:
```json
{
  "kind": "transport|messaging|sync|search|benchmark",
  "id": "unique-within-kind",
  "title": "human-readable title",
  "description": "detailed description",
  "family": "request_response|server_sent_events|...",
  "category": "transport|backend-family|frontend|enrichment",
  "status": "bootstrap-ready|planned|deferred|held",
  "tags": ["tag1", "tag2"],
  "metrics": {"latency_p50": 120, "latency_p95": 450, ...},
  "evidence": {
    "source": "benchmark|audit|projection",
    "confidence": 0.95,
    "staleness": "fresh|stale|very_stale|reindexing",
    "measured_at": "2026-11-11T10:30:00Z",
    "delay_applied_ms": 0
  },
  "created_at": "2026-11-11T10:00:00Z",
  "updated_at": "2026-11-11T10:30:00Z"
}
```

**Field Mappings**:
| Field | Type | Indexed | Stored | Notes |
|-------|------|---------|--------|-------|
| `kind` | keyword | yes | yes | Exact match, aggregations |
| `id` | keyword | yes | yes | Unique within kind |
| `title` | text + keyword | yes | yes | Full-text + exact |
| `description` | text | yes | yes | Full-text search |
| `family` | keyword | yes | yes | Filter, aggregations |
| `category` | keyword | yes | yes | Filter, aggregations |
| `status` | keyword | yes | yes | Filter, aggregations |
| `tags` | keyword | yes | yes | Multi-value, filter |
| `metrics` | flattened | yes | yes | Dynamic metric fields |
| `evidence` | object | yes | yes | Nested evidence contract |
| `created_at` | date | yes | yes | Range queries |
| `updated_at` | date | yes | yes | Range queries |

### Query DSL Contract
**Supported Query Types**:
1. **Term/Terms Query**: Exact match on keyword fields (`kind`, `family`, `category`, `status`, `tags`)
2. **Match Query**: Full-text on `title`, `description` (standard analyzer)
3. **Range Query**: Numeric/date ranges on `metrics.*`, `created_at`, `updated_at`, `evidence.measured_at`
4. **Bool Query**: Combination of above (must/should/filter)
5. **Aggregations**: Terms on `family`, `category`, `status`; stats on `metrics.*`

**Prohibited**: Script queries, percolate queries, geo queries (not in domain)

**Pagination**: `from` + `size` (max 1000), or `search_after` for deep pagination

### Freshness Integration (from `projection/freshness.go`)
- Search projection type: `ProjectionTypeSearch`
- Freshness threshold: 60 seconds (configurable via `SearchFreshnessThresholdSeconds`)
- Critical staleness: 300 seconds (`SearchCriticalStalenessSeconds`)
- Reindex triggers:
  - Staleness ≥ 300s → `search_staleness` trigger
  - Document lag ≥ 1000 → `search_staleness` trigger
  - Error rate ≥ 5% → `search_staleness` trigger
  - Scheduled: every 15 minutes (`0 */15 * * * *`)
- Freshness status exposed via `/api/projections/freshness`

### Evidence Contract Integration
Every document indexed includes `evidence` field conforming to `evidence.Evidence`:
- `source`: `SourceProjection`
- `confidence`: Calculated from source freshness (1.0 for fresh, 0.7 for stale, 0.3 for very_stale)
- `staleness`: From `FreshnessTracker.GetFreshness(ProjectionTypeSearch)`
- `measured_at`: Document `updated_at`
- `delay_applied_ms`: 0 (projection indexing is async but measured at index time)

### Transport Provider Contract (`EsTransport` interface)
```go
type TransportProvider interface {
    Create(ctx context.Context, t *Transport) error
    GetByID(ctx context.Context, id string) (*Transport, error)
    Search(ctx context.Context, query string) ([]*Transport, error)
}
```
- `Create`: PUT `/{index}/_doc/{id}` with document; returns error if ES unavailable
- `GetByID`: GET `/{index}/_doc/{id}`; returns nil if not found
- `Search`: POST `/{index}/_search` with parsed query; returns empty slice if ES unavailable
- All methods return `ErrNoEsClient` if client not configured

### Optional Profile Gate
Elasticsearch is an **optional profile** (see `main.go` handleDataPlatform):
- Only enabled when search projection is explicitly needed
- Not required for bootstrap or first-wave closure
- Activation requires: search use case justified, ES cluster provisioned, freshness config approved

## Consequences
- Index mapping is versioned; schema changes require reindex
- Query DSL is restricted to prevent unbounded queries
- Freshness contract binds projection staleness to evidence confidence
- Evidence contracts on every document enable audit trail
- Optional profile means: no ES dependency in CI unless profile activated

## Acceptance Criteria
- [ ] Index mapping applied and verified in integration test
- [ ] Query DSL supports all 5 query types; prohibited queries rejected
- [ ] Freshness tracker updates search projection status correctly
- [ ] Evidence contracts on indexed documents match `evidence.Evidence` schema
- [ ] `EsTransport` implements `TransportProvider` interface fully
- [ ] Optional profile gate documented and enforced (no ES in bootstrap CI)