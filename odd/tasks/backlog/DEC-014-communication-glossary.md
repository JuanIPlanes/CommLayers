---
id: DEC-014
title: Communication Family Glossary
status: proposed
date: 2026-11-11
---

## Context
The `main.go` defines 12 communication families across 4 categories (transport, backend-family, frontend, enrichment). This decision normalizes all families with canonical definitions as a binding glossary.

## Decision
### Category: Transport (6 families)

| Family | Status | Canonical Definition | Recommended When | Avoid When |
|--------|--------|---------------------|------------------|------------|
| `request_response` | bootstrap-ready | Bounded request-reply pattern; client sends request, server responds synchronously. Reference case for all comparisons. | Simple CRUD, bounded request-reply, periodic status checking | Continuous updates, bidirectional interaction |
| `polling` | planned | Client periodically requests status from server at fixed intervals. Compatibility fallback when push unavailable. | Periodic status checking, compatibility fallback, pull-based status checks | High-frequency updates, low-latency requirements |
| `long_polling` | planned | Client requests status; server holds request until status changes or timeout. Legacy push compatibility. | Blocking wait with eventual completion, legacy push compatibility, webhook-driven integration | Modern push available, high-frequency updates |
| `server_sent_events` | bootstrap-demo | Server pushes events to client over persistent HTTP connection. One-way, text-based, auto-reconnect. **Default push choice.** | One-way live updates, status feeds, progress streams, continuous stream processing | Bidirectional interaction, binary transport, client-to-server streaming |
| `websocket` | planned | Full-duplex bidirectional communication over single TCP connection. Binary/text frames. | Bidirectional interaction, duplex low latency, actor-style isolated concurrency, realtime document/state sync | Simple server push, one-way updates only |
| `webhook` | planned | Server calls client-registered HTTP endpoint on events. Outbound-only from server perspective. | Event notification to external systems, async callbacks, third-party integration | High-frequency updates, guaranteed delivery (use messaging) |

### Category: Backend-Family (6 families)

| Family | Status | Canonical Definition | Recommended When | Avoid When |
|--------|--------|---------------------|------------------|------------|
| `messaging_workflow` | deferred | Durable workflow execution with fan-out event delivery, immutable event history, distributed coordination. Includes broker, queue, workflow, event families. | Durable workflow execution, fan-out event delivery, immutable event history, distributed coordination | Bootstrap scope, simple request-response |
| `sync_replication_projection` | deferred | Replication and edge merge, realtime document/state sync, service networking and control-plane management. Includes sync, coordination, projection families. | Replication and edge merge, realtime document/state sync, service networking and control-plane management | Bootstrap scope, simple CRUD |
| `search_projection` | deferred | Full-text search projection over canonical data. Elasticsearch-backed. Optional profile. | Search-heavy workloads, faceted search, relevance ranking | Simple key-value lookups, bootstrap scope |
| `graph_projection` | deferred | Graph traversal projection over canonical data. Neo4j-backed. Optional profile. | Relationship-heavy queries, path finding, dependency analysis | Simple relational queries, bootstrap scope |
| `vector_projection` | deferred | Vector similarity projection over canonical data. Qdrant/pgvector-backed. Optional profile. | Semantic search, recommendation, embedding similarity | Exact match queries, bootstrap scope |
| `mongo_projection` | deferred | Document-native projection for document-oriented use cases. MongoDB-backed. Optional profile. | Explicitly justified document-native use cases | Relational data, bootstrap scope |

### Category: Frontend (1 family)

| Family | Status | Canonical Definition | Recommended When | Avoid When |
|--------|--------|---------------------|------------------|------------|
| `comparative_visualization` | held | Desktop comparison UI with TreeGraphView, SplitView, RightPanel components. FPS verification, visible delay interactions. | Multi-backend comparison, decision support, architecture review | Mobile-first, simple dashboards |

### Category: Enrichment (1 family)

| Family | Status | Canonical Definition | Recommended When | Avoid When |
|--------|--------|---------------------|------------------|------------|
| `pricing_cost` | held | USD freshness contract, confidence semantics, pricing fallback rules. Late-stage only. | Cost-aware decisions, budget planning, vendor comparison | Early-stage, bootstrap scope |

### Glossary Contract
Each family entry MUST include:
- `Name`: Canonical identifier (snake_case, matches code constants)
- `Category`: One of {transport, backend-family, frontend, enrichment}
- `Status`: One of {bootstrap-ready, bootstrap-demo, planned, deferred, held}
- `CanonicalDefinition`: Single authoritative sentence
- `RecommendedWhen`: Array of use cases
- `AvoidWhen`: Array of anti-patterns
- `Notes`: Implementation status, dependencies, constraints
- `Dependencies`: Other families this depends on (from graphNodes)
- `UnlockPrerequisites`: What must be true before implementation (for deferred/held)

### Vocabulary Freeze
- Family names are frozen: no additions/renames without gate 86agbv0k4 re-approval
- Categories are frozen: no new categories without architectural decision
- Status values are frozen: no new status values
- This glossary is the single source of truth for `/api/catalog` and `/api/first-wave/contract`

## Consequences
- All API responses, documentation, and code comments MUST use canonical family names
- New transport mechanisms require new family entry (not overloading existing)
- Status transitions follow: bootstrap-ready → planned/deferred → held → active
- Glossary version embedded in all evidence contracts

## Acceptance Criteria
- [ ] Glossary matches `main.go` families array exactly (12 families, 4 categories)
- [ ] `/api/catalog` returns glossary with all required fields
- [ ] Vocabulary freeze enforced: build fails on new family names without gate re-approval
- [ ] Each family has canonical definition, recommended/avoid, notes, dependencies, unlock prerequisites
- [ ] Glossary version tracked in evidence contracts