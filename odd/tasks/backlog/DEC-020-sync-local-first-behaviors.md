---
id: DEC-020
title: Sync and Local-First Behaviors
status: proposed
date: 2026-11-11
---

## Context
The sync demo (`handleSyncOverview`, `handleSyncSessionCreate`, `handleSyncSessionAction`) and local-first CRDT paradigm (`handleV2LocalFirstCRDT`) define sync behaviors, conflict detection, merge rules, and lag thresholds. This decision formalizes these as a binding contract.

## Decision
### Sync Session Model
**State Machine**:
```
created → active → aligned | conflict → resolved → closed
```
**Status Values**: `created`, `active`, `aligned`, `conflict`, `resolved`, `closed`

**Session Fields** (from `syncSession` in `main.go`):
| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | `sync-{timestamp}` |
| `Status` | string | Current state |
| `Primary` | string | Primary replica identifier |
| `Replica` | string | Replica identifier |
| `LastSyncedAt` | time.Time | Last successful sync |
| `Lag` | int | Divergence in milliseconds/operations |
| `Conflict` | bool | True if unresolved divergence |
| `CreatedAt` | time.Time | Session creation |
| `UpdatedAt` | time.Time | Last update |

### Conflict Detection Contract
**Detection Triggers**:
1. **Version Mismatch**: Primary and replica versions differ beyond threshold
2. **Divergent Mutations**: Both sides mutated since last sync
3. **Lag Threshold Exceeded**: `Lag > MaxAllowedLag` (default: 5000ms or 100 operations)

**Conflict Payload**:
```json
{
  "session_id": "sync-123",
  "conflict_type": "version_mismatch|divergent_mutations|lag_exceeded",
  "primary_version": 15,
  "replica_version": 12,
  "divergent_operations": ["op-1", "op-2"],
  "detected_at": "2026-11-11T10:30:00Z",
  "lag_ms": 7200
}
```

### Merge Rules Contract (from `handleV2LocalFirstCRDT`)
**Rule 1 — Single-Side Divergence**:
> "Single-side divergence catches up by replaying the latest persisted value."
- If only primary or only replica has new operations since last sync
- Resolution: Replay missing operations from the side that has them
- Result: `Status` → `aligned`, `Conflict` → `false`, `Lag` → `0`

**Rule 2 — Equal-Version Divergence**:
> "Equal-version divergence merges deterministically into a combined value in the bootstrap model."
- Both sides have same version but different values (concurrent edits)
- Resolution: Deterministic merge function (e.g., last-writer-wins by timestamp, or CRDT merge)
- Result: `Status` → `resolved`, `Conflict` → `false`, merged value persisted

**Rule 3 — Conflict and Lag as Signals**:
> "Conflict and lag are surfaced as first-class signals for later CRDT evolution."
- All conflicts recorded with full context
- Lag tracked continuously
- Signals feed `local_first_crdt` paradigm readiness

### Lag Thresholds
| Threshold | Value | Action |
|-----------|-------|--------|
| `WarningLagMs` | 2000 | Log warning, alert dashboard |
| `CriticalLagMs` | 5000 | Mark session `conflict`, block consensus (DEC-018) |
| `MaxAllowedLagMs` | 10000 | Force resolution, escalate |

### Sync API Contract
**Create Session** (`POST /api/sync/sessions`):
```json
{
  "primary": "service-a",
  "replica": "service-b",
  "initial_lag_ms": 0
}
```
→ Returns `SyncSession` with `Status: "created"`

**Session Action** (`POST /api/sync/sessions/{id}`):
```json
{ "action": "sync|resolve|force_align" }
```
- `sync`: Attempt normal sync (apply merge rules)
- `resolve`: Apply deterministic merge for equal-version divergence
- `force_align`: Force primary value (admin only)

**Get Session** (`GET /api/sync/sessions/{id}`):
Returns current `SyncSession` with updated `Lag`, `Conflict`, `Status`

### Local-First CRDT Paradigm Readiness (from `handleV2LocalFirstCRDT`)
**Signals**:
- `recentSessionCount`: Active sync sessions in window
- `alignedSessions`: Count with `Status == "aligned"`
- `conflictSessions`: Count with `Conflict == true`
- `maxObservedLag`: Maximum `Lag` across sessions

**Merge Rules** (exposed for paradigm):
1. Single-side divergence → replay latest persisted
2. Equal-version divergence → deterministic combine
3. Conflict/lag → first-class signals

**Note**: "This slice is local-first and CRDT-oriented, not a full CRDT engine yet."

### Evidence Integration
Every sync session action produces `EvidenceContract`:
- `source`: `SourceAudit` (sync is operational, not benchmark)
- `confidence`: 1.0 for aligned, 0.7 for conflict, 0.3 for forced
- `staleness`: `StalenessFresh` if `LastSyncedAt` < 60s, else `StalenessStale`
- `measured_at`: Action timestamp
- `delay_applied_ms`: Sync operation duration

## Consequences
- Sync sessions are the contract for replication state
- Conflict detection is synchronous and blocking
- Merge rules are deterministic and auditable
- Lag thresholds gate consensus-driven paradigm (DEC-018)
- Local-first CRDT paradigm reuses sync signals — no separate engine yet

## Acceptance Criteria
- [ ] Sync session state machine enforced: invalid transitions rejected
- [ ] Conflict detection triggers on all three conditions
- [ ] Merge rules 1, 2, 3 implemented and tested
- [ ] Lag thresholds enforced with correct actions
- [ ] API endpoints return correct session state
- [ ] Evidence contracts attach to all sync actions
- [ ] Local-first CRDT readiness signals match `handleV2LocalFirstCRDT` output