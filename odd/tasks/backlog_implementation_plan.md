# Backlog Implementation Plan — 26 Definition/Contract Tasks

**Generated**: 2026-11-10  
**Status**: Planning phase — design workshops, DEC docs, contract formalization

---

## Overview

This document outlines the implementation approach for 26 backlog tasks that require **design decisions, DEC documents, or contract formalization** rather than code implementation. These tasks define the architectural boundaries, vocabularies, and acceptance criteria that subsequent implementation waves will follow.

**Key Principle**: Most tasks can proceed in parallel as independent design work. Only tasks with explicit dependencies on partial task implementations (86agbv0ne, 86agbv0ng) must wait.

---

## Dependency Map

```
PARTIAL TASKS (must complete first):
├── 86agbv0ne — Evidence contracts (in progress)
│   ├── BLOCKS: 86agbv0p4 (delay/timeout evidence)
│   └── BLOCKS: 86agbv0p6 (logs/traces/audit schema)
│
└── 86agbv0ng — PKI observability (in progress)
    ├── BLOCKS: 86agbv0nc (Vault PKI/mTLS)
    └── BLOCKS: 86agbv0ng (auth/PKI checks)

INDEPENDENT (can start NOW): 24 tasks
```

---

## Task Details

### Group A: Redis & Infrastructure Boundaries (Independent)

| Task ID | Title | Approach | Dependencies | Effort | Can Start |
|---------|-------|----------|--------------|--------|-----------|
| **86agbv0n8** | Redis usage boundaries | Formalize `RedisUsageBounds` in `backend/redis.go` into DEC; define cache/rate-limit/lock contracts; add runtime validation | None | 2 days | **Immediate** |

**Deliverable**: DEC document + runtime enforcement in `cache/redis.go`  
**Details**: Current `backend/redis.go` defines types but no enforcement. Need to codify: Redis is **never** source of truth; TTL bounds; max entries; eviction policy; lock semantics.

---

### Group B: Benchmark & Observability Framework (Independent)

| Task ID | Title | Approach | Dependencies | Effort | Can Start |
|---------|-------|----------|--------------|--------|-----------|
| **86agbv0np** | Async visibility benchmarks | Activate `BenchmarkCollector` in middleware; wire `/api/benchmarks/metrics` to return real data; define benchmark scenarios in config | None | 3 days | **Immediate** |
| **86agbv0p8** | Benchmark matrix and thresholds | Extend `benchmarks.go` with YAML/JSON config for scenarios, thresholds, pass/fail criteria; add CI integration | None | 3 days | **Immediate** |
| **86agbv0ph** | Visible delay/state transition interactions | Define frontend-backend contract for delay metadata (already in `delayAppliedMs`); add visualization spec for state transitions | None | 2 days | **Immediate** |

**Deliverable**: Working benchmark pipeline + threshold config + delay visualization spec

---

### Group C: First-Wave Gates & Contracts (Independent)

| Task ID | Title | Approach | Dependencies | Effort | Can Start |
|---------|-------|----------|--------------|--------|-----------|
| **86agbv0n9** | Layer acceptance gates | Formalize acceptance gates from stream definitions (`main.go` streams) into executable checks; add gate evaluation endpoint | None | 2 days | **Immediate** |
| **86agbv0pr** | Deferred implementation boundaries | Document wave reactivation rules as executable preconditions; add `/api/deferred-waves/readiness` | None | 2 days | **Immediate** |
| **86agbv0pu** | v2 unlock prerequisites | Convert `v2Readiness` blockedBy/unlockPath into formal prerequisites with measurable criteria | None | 2 days | **Immediate** |

**Deliverable**: Gate evaluation endpoints + wave reactivation preconditions + v2 prerequisite checklist

---

### Group D: Messaging & Workflow Family (Independent)

| Task ID | Title | Approach | Dependencies | Effort | Can Start |
|---------|-------|----------|--------------|--------|-----------|
| **86agbv0n7** | Communication family glossary | Extract `families` slice from `main.go` into standalone glossary DEC with normative definitions | None | 1 day | **Immediate** |
| **86agbv0nk** | Broker queue/workflow criteria | Formalize `messaging_workflow` family deliverables (broker criteria, queue semantics, workflow contract, failure mode matrix) | None | 3 days | **Immediate** |
| **86agbv0nt** (PARTIAL) | Failure mode checks | **IMPLEMENTED** — `internal/messaging/failure.go` with `FailureModeMatrix`, `RunFailureChecks`, demo endpoints | — | — | **Done** |

**Note**: 86agbv0nt was a "missing task" — now implemented with failure mode matrix, checks, and API endpoints.

---

### Group E: Sync & Replication (Independent)

| Task ID | Title | Approach | Dependencies | Effort | Can Start |
|---------|-------|----------|--------------|--------|-----------|
| **86agbv0nw** | Sync and local-first behaviors | Formalize `sync_replication_projection` family; define sync behaviors, conflict handling, replication checks, consensus criteria | None | 3 days | **Immediate** |
| **86agbv0nv** | Distributed state/consensus behaviors | Define consensus-driven paradigm acceptance gates; quorum rules; divergence detection | None | 3 days | **Immediate** |

**Deliverable**: Sync family DEC + consensus behavior specification

---

### Group F: Projection Family (Partial — One Implemented)

| Task ID | Title | Approach | Dependencies | Effort | Can Start |
|---------|-------|----------|--------------|--------|-----------|
| **86agbv0p5** | Elasticsearch projection model | Extend `es_transport.go` stub with full search projection contract; define index mappings, query DSL, refresh policy | None | 3 days | **Immediate** |
| **86agbv0p2** (PARTIAL) | Projection freshness/reindex | **IMPLEMENTED** — `internal/projection/freshness.go` with `FreshnessTracker`, reindex triggers, demo endpoints | — | — | **Done** |

**Note**: 86agbv0p2 was a "missing task" — now implemented with freshness tracking, staleness thresholds, reindex triggers, and API endpoints.

---

### Group G: Localization & i18n (Independent)

| Task ID | Title | Approach | Dependencies | Effort | Can Start |
|---------|-------|----------|--------------|--------|-----------|
| **86agbv0p1** | i18n verification checks | Add locale validation to all endpoints; verify `en`/`es` parity; add missing translation keys | None | 2 days | **Immediate** |
| **86agbv0pb** | en/es resource model | Extract message catalogs to external resource files; define resource loading contract; add locale negotiation | None | 2 days | **Immediate** |

**Deliverable**: i18n test suite + externalized resource bundles + locale negotiation spec

---

### Group H: Pricing & Cost Enrichment (Independent)

| Task ID | Title | Approach | Dependencies | Effort | Can Start |
|---------|-------|----------|--------------|--------|-----------|
| **86agbv0pp** | Pricing dependency gate | Formalize pricing as late-stage gate; define dependency on backend-family stability + localization | None | 1 day | **Immediate** |
| **86agbv0pk** | USD freshness/confidence contract | Extend `pricePoint` with formal freshness semantics; define confidence intervals; add staleness escalation | None | 2 days | **Immediate** |
| **86agbv0pq** | Pricing verification/fallback | Add fallback chain validation; verify provider APIs; define circuit breaker for pricing | None | 2 days | **Immediate** |

**Deliverable**: Pricing gate DEC + freshness contract + fallback verification suite

---

### Group I: Canonical Data & Domain Model (Independent)

| Task ID | Title | Approach | Dependencies | Effort | Can Start |
|---------|-------|----------|--------------|--------|-----------|
| **86agbv0nd** | Canonical Postgres domain model | Formalize `store.go` entities into DEC; define ownership, schema evolution, audit lineage, migration strategy | None | 3 days | **Immediate** |
| **86agbv0nb** | Map applications and dependencies | Convert `graphNodes` + `families` from `main.go` into application dependency map with 14 use cases | None | 2 days | **Immediate** |
| **86agbv0pm** | Paradigm taxonomy and mapping | Formalize 5 paradigms from `v2/paradigms` into taxonomy with mapping to transport/backend families | None | 2 days | **Immediate** |

**Deliverable**: Canonical domain DEC + application dependency map + paradigm taxonomy

---

### Group J: Frontend Comparative Visualization (Independent)

| Task ID | Title | Approach | Dependencies | Effort | Can Start |
|---------|-------|----------|--------------|--------|-----------|
| **86agbv0pj** | Desktop comparison UI contract | Define comparison UI requirements: side-by-side, matrix, radar chart; API contract for comparison data | None | 3 days | **Immediate** |
| **86agbv0pn** (PARTIAL) | Frontend FPS/comparison verification | **IMPLEMENTED** — `TreeGraphView.tsx` with `useFPSMonitor`, `useComparisonVerification`, live overlays | — | — | **Done** |

**Note**: 86agbv0pn was a "missing task" — now implemented with FPS monitoring, comparison verification, and live overlay indicators.

---

### Group K: Dependent on 86agbv0ne (Evidence Contracts)

| Task ID | Title | Approach | Dependencies | Effort | Can Start |
|---------|-------|----------|--------------|--------|-----------|
| **86agbv0p4** | Delay/timeout evidence contracts | Extend `evidence.go` with delay/timeout evidence types; define source, confidence, staleness fields | **86agbv0ne** | 2 days | Wait for 86agbv0ne |
| **86agbv0p6** | Logs/traces/audit schema | Define telemetry schema in `evidence.go`; log format, trace context, audit fields; add OpenTelemetry mapping | **86agbv0ne** | 3 days | Wait for 86agbv0ne |

---

### Group L: Dependent on 86agbv0ng (PKI Observability)

| Task ID | Title | Approach | Dependencies | Effort | Can Start |
|---------|-------|----------|--------------|--------|-----------|
| **86agbv0nc** | Vault PKI/mTLS | Define Vault PKI integration; certificate lifecycle; mTLS phases 1-3 from `handleSecurityBootstrap` | **86agbv0ng** | 5 days | Wait for 86agbv0ng |
| **86agbv0ng** | Auth/PKI observability checks | Add PKI health endpoints; cert expiry monitoring; mTLS connection verification | **86agbv0ng** (self) | 3 days | Wait for 86agbv0ng |

---

## Summary: Start Now vs Wait

### ✅ CAN START IMMEDIATELY (24 tasks)
| Task ID | Category |
|---------|----------|
| 86agbv0n8 | Redis boundaries |
| 86agbv0np | Async benchmarks |
| 86agbv0p8 | Benchmark matrix |
| 86agbv0ph | Delay visualization |
| 86agbv0n9 | Layer gates |
| 86agbv0pr | Deferred boundaries |
| 86agbv0pu | v2 prerequisites |
| 86agbv0n7 | Family glossary |
| 86agbv0nk | Broker/workflow criteria |
| 86agbv0nw | Sync behaviors |
| 86agbv0nv | Consensus behaviors |
| 86agbv0p5 | ES projection model |
| 86agbv0p1 | i18n verification |
| 86agbv0pb | en/es resource model |
| 86agbv0pp | Pricing gate |
| 86agbv0pk | USD freshness contract |
| 86agbv0pq | Pricing fallback |
| 86agbv0pj | Comparison UI contract |
| 86agbv0nd | Canonical domain model |
| 86agbv0nb | Application map |
| 86agbv0pm | Paradigm taxonomy |
| **86agbv0nt** | **Failure modes (DONE)** |
| **86agbv0p2** | **Projection freshness (DONE)** |
| **86agbv0pn** | **FPS verification (DONE)** |

### ⏳ MUST WAIT (2 tasks)
| Task ID | Blocked By |
|---------|------------|
| 86agbv0p4 | 86agbv0ne (evidence contracts) |
| 86agbv0p6 | 86agbv0ne (evidence contracts) |
| 86agbv0nc | 86agbv0ng (PKI observability) |
| 86agbv0ng | 86agbv0ng (PKI observability) |

---

## Recommended Execution Order

### Phase 1: Foundation (Week 1-2) — Parallelizable
1. **86agbv0n8** Redis boundaries — enables cache/rate-limit contracts
2. **86agbv0n7** Family glossary — vocabulary for all later work
3. **86agbv0nd** Canonical domain model — storage ownership for all families
4. **86agbv0nb** Application map — defines 14 use cases for comparison
5. **86agbv0pm** Paradigm taxonomy — v2 unlock criteria

### Phase 2: Benchmark & Gates (Week 2-3) — Parallelizable
6. **86agbv0np** Async benchmarks — activates existing collector
7. **86agbv0p8** Benchmark matrix — config-driven thresholds
8. **86agbv0n9** Layer gates — executable acceptance criteria
9. **86agbv0pr** Deferred boundaries — wave reactivation preconditions
10. **86agbv0pu** v2 prerequisites — measurable unlock criteria

### Phase 3: Family Contracts (Week 3-4) — Parallelizable
11. **86agbv0nk** Broker/workflow criteria
12. **86agbv0nw** Sync behaviors
13. **86agbv0nv** Consensus behaviors
14. **86agbv0p5** ES projection model
15. **86agbv0p1** + **86agbv0pb** i18n verification + resource model
16. **86agbv0pp** + **86agbv0pk** + **86agbv0pq** Pricing gate + freshness + fallback
17. **86agbv0ph** Delay visualization spec
18. **86agbv0pj** Comparison UI contract

### Phase 4: Dependent Work (After 86agbv0ne, 86agbv0ng complete)
19. **86agbv0p4** Delay/timeout evidence
20. **86agbv0p6** Logs/traces/audit schema
21. **86agbv0nc** Vault PKI/mTLS
22. **86agbv0ng** Auth/PKI observability

---

## Definition of Done for Each Task

Each task produces **one of**:
- **DEC document** in `docs/decisions/DEC-XXX-<slug>.md`
- **Contract specification** as Go types + JSON Schema in `internal/<domain>/contracts.go`
- **Executable validation** as HTTP endpoint returning pass/fail with evidence
- **Test suite** verifying contract compliance

All deliverables must be:
1. Reviewed in PR with at least one approver
2. Added to `/api/catalog` or relevant discovery endpoint
3. Documented in `odd/tasks/` with implementation notes

---

## Notes

- The 3 "missing tasks" (86agbv0pn, 86agbv0nt, 86agbv0p2) are **implemented and tested** — see code changes in:
  - `frontend/src/components/TreeGraphView.tsx` (FPS + comparison verification)
  - `services/api/internal/messaging/failure.go` (failure mode matrix + checks)
  - `services/api/internal/projection/freshness.go` (freshness tracker + reindex triggers)
  - `services/api/cmd/server/main.go` (new endpoints integrated)

- Tasks marked "Independent" have no code dependencies on partial implementations and can be assigned to parallel design tracks immediately.

- The dependent tasks (86agbv0p4, 86agbv0p6, 86agbv0nc, 86agbv0ng) should be scheduled immediately after their prerequisites land.