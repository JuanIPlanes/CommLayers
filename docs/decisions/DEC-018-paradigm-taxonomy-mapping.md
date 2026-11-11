---
id: DEC-018
title: Paradigm Taxonomy and Mapping
status: proposed
date: 2026-11-11
---

## Context
The `main.go` defines 5 v2 architectural paradigms with unlock dependencies. This decision formalizes the paradigm taxonomy, their mapping to v1 streams, and unlock dependencies as a binding contract.

## Decision
### The 5 Paradigms (Canonical Taxonomy)

#### 1. Event-Driven Architecture (`event_driven_architecture`)
**Definition**: Commands produce durable events; events drive projections; event log is the source of truth.
**Status**: planned
**v1 Stream Mapping**:
- Transport: `server_sent_events` (event broadcast), `websocket` (command submission)
- Messaging: `workflow_engine` (command → event → projection), `broker_event_bus` (durable event backbone)
- Sync: `sync_replication` (event-sourced replica state)
- Projection: All projection types consume event log
- Persistence: `domain_event` records in PostgreSQL (canonical event store)
**Unlock Dependencies**:
- Transport slices: SSE + WebSocket implemented
- Messaging slices: Workflow engine + broker event bus implemented
- Durable event contract: `DomainEvent` schema frozen, persistence guaranteed
- Evidence integration: All events carry `EvidenceContract`

**Readiness Signals** (from `handleV2EventDriven`):
- Commands: `create workflow run`, `create sync session`, `mutate sync replica`, `create transport session`, `start async job`
- Events: Persisted `domain_event` records with kind projections
- Source of Truth: PostgreSQL `domain_event` table
- Transient Broadcast: Redis pubsub via `/api/events`

#### 2. Reactive Stream Driven (`reactive_stream_driven`)
**Definition**: Continuous streams of signals (events, metrics, state changes) with backpressure-aware processing.
**Status**: planned
**v1 Stream Mapping**:
- Transport: `server_sent_events` (primary stream surface), Redis transient events (fast path)
- Messaging: Workflow/async job events as stream signals
- Benchmarks: `AsyncVisibilityBenchmarks` as stream pressure indicators
- Projection: Refresh cadence follows event density (not fixed polling)
**Unlock Dependencies**:
- Messaging maturity: Workflow engine events flowing
- Sync maturity: Sync session signals flowing
- Benchmark maturity: Queue delay, timeout, transition metrics stable
- Backpressure signals: Event counts per kind, workflow/async job event rates

**Readiness Signals** (from `handleV2ReactiveStreamDriven`):
- Recent event count, event kinds, latest event timestamp
- Transport surface: SSE + Redis transient events
- Backpressure hints: Event counts per kind, workflow/async job events as first throttling candidates

#### 3. Local-First CRDT (`local_first_crdt`)
**Definition**: Optimistic local mutation with deterministic merge; conflict and lag as first-class signals.
**Status**: planned
**v1 Stream Mapping**:
- Sync: `sync_replication` (primary/replica drift, explicit merge behavior)
- Transport: `websocket` (duplex sync), `server_sent_events` (state broadcast)
- Persistence: `sync_session` records in PostgreSQL
**Unlock Dependencies**:
- Sync slice maturity: Primary/replica demo operational
- Conflict strategy stabilization: Merge rules defined and tested
- Lag thresholds: Max observed lag tracked, alerting configured

**Readiness Signals** (from `handleV2LocalFirstCRDT`):
- Sync signals: Recent session count, aligned sessions, conflict sessions, max observed lag
- Merge rules: Single-side divergence replay, equal-version deterministic combine, conflict/lag as signals
- Note: "Local-first and CRDT-oriented, not a full CRDT engine yet"

#### 4. Consensus-Driven (`consensus_driven`)
**Definition**: Quorum-based agreement; divergent states block consensus until replication/merge resolves.
**Status**: planned
**v1 Stream Mapping**:
- Sync: `sync_replication` (aligned sessions = quorum-ready)
- Messaging: `workflow_engine` (stage completion as agreement signals)
- Events: `domain_event` stream as change history
**Unlock Dependencies**:
- Coordination semantics beyond bootstrap-grade in-memory demos
- Quorum rules defined and tested
- Blocking conditions visible and measurable

**Readiness Signals** (from `handleV2ConsensusDriven`):
- Coordination signals: Candidate sessions, quorum-ready, quorum-blocked
- Decision rules: Zero lag + no conflict = quorum-ready; any divergence blocks
- Note: "Consensus-driven in orientation, not full Raft/Paxos"

#### 5. Planet-Scale Composite (`planet_scale_composite_architectures`)
**Definition**: Composition of multiple paradigms on shared runtime foundation; events as shared history.
**Status**: planned
**v1 Stream Mapping**:
- All paradigms combined: event-driven (backbone) + reactive (signals) + local-first (divergence) + consensus (agreement)
- Projection: Search/Graph/Vector as downstream consumers
- Pricing: Cost enrichment as downstream consumer
- Frontend: Comparative visualization as downstream consumer
**Unlock Dependencies**:
- Projection maturity: Search/Graph/Vector projections stable
- Pricing maturity: USD freshness contract, confidence semantics
- Frontend comparative maturity: UI stable against all v1 backend families
- Composition rules validated

**Readiness Signals** (from `handleV2PlanetScaleComposite`):
- Building blocks: All 4 paradigms as components
- Composition rules: Persisted events = shared history, Redis = fast path, Sync/consensus = convergence signals

### Paradigm-to-v1-Stream Dependency Matrix
| Paradigm | Transport | Messaging | Sync | Projection | Benchmark | Persistence |
|----------|-----------|-----------|------|------------|-----------|-------------|
| Event-Driven | SSE, WS | Workflow, Broker | Sync | All | - | domain_event |
| Reactive Stream | SSE, Redis | Workflow, Async | - | Refresh cadence | Queue delay, transitions | - |
| Local-First CRDT | WS, SSE | - | Sync | - | - | sync_session |
| Consensus | - | Workflow | Sync | - | - | domain_event, sync_session |
| Planet-Scale | All | All | All | All | All | All |

### Taxonomy Freeze Rules
- Paradigm names are frozen: no additions/renames without architectural decision
- Status values frozen: `planned` | `active` | `implemented` | `deprecated`
- Unlock dependencies are binding: paradigm cannot move to `active` until all dependencies met
- Mapping to v1 streams is the contract: v2 paradigms MUST reuse v1 runtime surfaces

### Readiness API Contract
`/api/v2/paradigms/{paradigm}` returns:
```json
{
  "paradigm": "event_driven_architecture",
  "status": "active",
  "commands": [...],
  "events": [...],
  "projections": {...},
  "notes": [...]
}
```

## Consequences
- Paradigm taxonomy is the north star for v2 architecture
- Each paradigm's unlock dependencies are executable checks
- v2 paradigms compose on v1 runtime — no new infrastructure until planet-scale
- Readiness signals are observable metrics, not aspirations

## Acceptance Criteria
- [ ] All 5 paradigms defined with canonical definitions
- [ ] v1 stream mapping matrix complete and accurate
- [ ] Unlock dependencies executable and enforced in CI
- [ ] `/api/v2/paradigms/*` endpoints return correct readiness signals
- [ ] Taxonomy freeze enforced: build fails on new paradigms without ADR