---
id: DEC-028
title: Distributed State and Consensus Behaviors
status: proposed
date: 2026-11-11
---

## Context
The consensus-driven paradigm demo (`handleV2ConsensusDriven`) and sync sessions define consensus behaviors: quorum, blocking, agreement. This decision formalizes the consensus paradigm contract.

## Decision
### Consensus Paradigm Model (from `handleV2ConsensusDriven`)

#### Coordination Signals
```go
type CoordinationSignals struct {
    CandidateSessions: int,      // Total sync sessions in window
    QuorumReady:       int,      // Sessions with zero lag + no conflict
    QuorumBlocked:     int,      // Sessions with lag > 0 or conflict
    Rule:              string,   // "aligned sessions count as ready, divergent sessions count as blocked"
}
```

#### Quorum Rules
**Ready (Quorum Member)**:
- `Status == "aligned"` AND `Conflict == false` AND `Lag == 0`
- Represents a replica that has caught up and has no divergence

**Blocked (Non-Member)**:
- `Conflict == true` OR `Lag > 0`
- Represents a replica that cannot participate in consensus

**Quorum Calculation**:
```
QuorumThreshold = ceil(CandidateSessions / 2) + 1  // Majority
ConsensusAchieved = QuorumReady >= QuorumThreshold
```

#### Decision Rules (from `handleV2ConsensusDriven`)
1. **Ready Rule**: "A session with zero lag and no conflict is treated as quorum-ready in the bootstrap model."
2. **Blocking Rule**: "Any divergence blocks consensus until replication or merge resolves the difference."
3. **History Rule**: "The event stream supplies the change history; sync sessions supply the current coordination view."

#### Agreement Protocol (Bootstrap Model)
**Not a full Raft/Paxos** — this is consensus-oriented for visibility:

```
Phase 1: Collect Signals
  → Query all sync sessions
  → Classify each as Ready or Blocked
  
Phase 2: Evaluate Quorum
  → Count Ready vs Blocked
  → Apply QuorumThreshold
  
Phase 3: Signal Result
  → ConsensusAchieved = true/false
  → Blocked sessions listed with reasons
```

#### Sync Session as Consensus Indicator (from `handleV2ConsensusDriven` and DEC-020)
| Session State | Consensus Signal |
|---------------|------------------|
| `aligned`, no conflict, lag=0 | **Ready** — votes yes |
| `aligned`, no conflict, lag>0 | **Blocked** — lag prevents vote |
| `conflict` | **Blocked** — divergence prevents vote |
| `resolved` | **Ready** — if lag=0 after resolve |
| `active` | **Blocked** — in progress |

### State Machine for Consensus Sessions
```
created → evaluating → consensus_achieved | consensus_blocked
                    ↓
              re-evaluate (on session change)
```

### API Contract
**`GET /api/v2/paradigms/consensus-driven`**:
```json
{
  "paradigm": "consensus_driven",
  "status": "active",
  "coordination_signals": {
    "candidate_sessions": 12,
    "quorum_ready": 8,
    "quorum_blocked": 4,
    "rule": "aligned sessions count as ready, divergent sessions count as blocked",
    "quorum_threshold": 7,
    "consensus_achieved": true
  },
  "decision_rules": [
    "A session with zero lag and no conflict is treated as quorum-ready in the bootstrap model.",
    "Any divergence blocks consensus until replication or merge resolves the difference.",
    "The event stream supplies the change history; sync sessions supply the current coordination view."
  ],
  "blocked_sessions": [
    {"id": "sync-3", "reason": "lag: 2500ms"},
    {"id": "sync-7", "reason": "conflict: divergent_mutations"},
    {"id": "sync-11", "reason": "lag: 800ms"},
    {"id": "sync-12", "reason": "conflict: version_mismatch"}
  ]
}
```

### Evidence Integration
Consensus evaluation produces `EvidenceContract`:
- `source`: `SourceAudit` (operational, not benchmark)
- `confidence`: 1.0 if consensus achieved, 0.7 if blocked
- `staleness`: From latest session update
- `measured_at`: Evaluation timestamp
- `delay_applied_ms`: Evaluation latency

### Integration with Sync (DEC-020) and Local-First CRDT (DEC-018)
- Consensus reuses `sync_session` records — no separate consensus engine
- `quorum_ready` ↔ `alignedSessions` from local-first CRDT
- `quorum_blocked` ↔ `conflictSessions` + laggy sessions
- Event stream (`domain_event`) provides change history for agreement

### Thresholds
| Threshold | Value | Action |
|-----------|-------|--------|
| `QuorumThreshold` | Majority (N/2 + 1) | Minimum ready for consensus |
| `MaxLagForReady` | 0ms | Any lag blocks |
| `MaxConflictForReady` | 0 | Any conflict blocks |
| `ReevaluateInterval` | 5s | How often to re-check |

### Future Evolution (from `handleV2ConsensusDriven` notes)
> "This slice is consensus-driven in orientation, not a full Raft/Paxos implementation. It makes agreement and blocking conditions visible using the current persisted sync and event signals."

**Evolution Path**:
1. Current: Visibility-only (bootstrap model)
2. Next: Quorum-based write coordination (require ready quorum for mutations)
3. Future: Full Raft/Paxos for critical coordination

## Consequences
- Consensus is visibility-only in current model — no blocking writes
- Sync sessions are the consensus indicators — no separate quorum nodes
- Blocking is deterministic: any lag/conflict = blocked
- Evidence contracts enable audit of consensus decisions
- Evolution path documented but not implemented

## Acceptance Criteria
- [ ] Coordination signals match `handleV2ConsensusDriven` output
- [ ] Quorum calculation uses majority threshold
- [ ] Ready/Blocked classification matches rules exactly
- [ ] API returns correct consensus state with blocked session details
- [ ] Evidence contracts on consensus evaluations
- [ ] Integration with sync sessions verified (shared data)
- [ ] Evolution path documented in decision rules