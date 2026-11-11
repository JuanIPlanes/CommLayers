---
id: DEC-029
title: Delay and Timeout Evidence Contract
status: proposed
date: 2026-11-11
---

## Context
Extends evidence contracts (86agbv0ne - done) with delay/timeout specific evidence for async operations. This decision formalizes the delay/timeout evidence contract.

## Decision
### Extended Evidence Fields
Adds to base `Evidence` (from `evidence/contracts.go`):

```go
type AsyncEvidence struct {
    Evidence                    // Base: source, confidence, staleness, measured_at, delay_applied_ms
    QueueDelayMs       float64  `json:"queue_delay_ms,omitempty"`       // Time in async queue
    TimeoutThresholdMs float64  `json:"timeout_threshold_ms,omitempty"` // Configured timeout limit
    TimeoutRemainingMs float64  `json:"timeout_remaining_ms,omitempty"` // Time until timeout
    StateTransitionMs  float64  `json:"state_transition_ms,omitempty"`  // Transition duration
    AllowedTransition  bool     `json:"allowed_transition,omitempty"`   // Was transition allowed
    BlockedReason      string   `json:"blocked_reason,omitempty"`       // If blocked, why
}
```

### Evidence Contract for Async Operations
```go
type AsyncEvidenceContract struct {
    EvidenceContract               // Base: evidence, success, success_code, error_code, error_message
    AsyncEvidence                  // Extended fields above
}
```

### When to Use AsyncEvidenceContract
| Operation Type | Evidence Contract |
|----------------|-------------------|
| Transport polling/long-polling | `AsyncEvidenceContract` |
| SSE event delivery | `AsyncEvidenceContract` (delay_breakdown) |
| WebSocket message | `AsyncEvidenceContract` |
| Workflow stage transition | `AsyncEvidenceContract` |
| Async LRO status check | `AsyncEvidenceContract` |
| Sync session action | `AsyncEvidenceContract` |
| Benchmark runs | `BenchmarkEvidence` (DEC-009) |

### Success/Error Codes for Async (extends DEC-011)
**Additional Success Codes**:
| Code | Meaning |
|------|---------|
| `success_delayed` | Completed but exceeded warning threshold |
| `success_partial` | Completed with degraded quality (from base) |

**Additional Error Codes**:
| Code | HTTP | When |
|------|------|------|
| `timeout` | 504 | `TimeoutRemainingMs ≤ 0` |
| `transition_blocked` | 409 | `AllowedTransition == false` |
| `queue_exhausted` | 503 | Queue at capacity |
| `stale_data` | 409 | `Staleness == very_stale|reindexing|expired` |

### Delay/Timeout Taxonomy
**Delay Categories** (for `delay_applied_ms`):
| Category | Range | Label |
|----------|-------|-------|
| `instant` | 0-50ms | Near-instant |
| `fast` | 50-200ms | Fast |
| `normal` | 200-1000ms | Normal |
| `slow` | 1000-5000ms | Slow (warning) |
| `critical` | >5000ms | Critical (alert) |

**Timeout Categories** (for `TimeoutThresholdMs`):
| Category | Range | Use Case |
|----------|-------|----------|
| `quick` | 1-5s | Quick polls, health checks |
| `standard` | 5-30s | Default async operations |
| `extended` | 30-300s | Long-running workflows |
| `infinite` | 0 (disabled) | Manual/approval workflows |

### Frontend Consumption
**Delay Metadata in Responses**:
```json
{
  "delay_applied_ms": 1250,
  "delay_breakdown": {
    "queue_delay_ms": 450,
    "processing_delay_ms": 800,
    "network_delay_ms": 0
  },
  "timeout_threshold_ms": 30000,
  "timeout_remaining_ms": 28750,
  "state_transition": {
    "from": "queued",
    "to": "running",
    "transition_ms": 50,
    "allowed": true
  },
  "evidence": {
    "source": "audit",
    "confidence": 0.95,
    "staleness": "fresh",
    "measured_at": "2026-11-11T10:30:00Z",
    "delay_applied_ms": 1250,
    "queue_delay_ms": 450,
    "timeout_threshold_ms": 30000,
    "timeout_remaining_ms": 28750,
    "state_transition_ms": 50,
    "allowed_transition": true
  }
}
```

### Collection and Aggregation
**Benchmark Collector** (from DEC-009) records:
- `RecordQueueDelay(delayMs)` → `QueueDelayMs`
- `RecordTimeout(timeoutMs)` → `TimeoutThresholdMs` (enforced)
- `RecordTransition(transitionMs)` → `StateTransitionMs`

**Aggregation** (for dashboard):
- p50/p95/p99 per operation type
- Delay breakdown averages
- Timeout breach rate
- Blocked transition count

### Integration with PKI Observability (DEC-004)
When mTLS active (phase 2+):
- `network_delay_ms` includes `tls_handshake_ms` + `cert_verification_ms`
- `pki_phase` included in evidence metadata
- Critical service pair delays tracked separately

## Consequences
- AsyncEvidenceContract is mandatory for all async operations
- Delay/timeout taxonomy enables consistent UI treatment
- Extended error codes enable precise error handling
- Benchmark integration provides historical baselines
- PKI integration captures security-related latency

## Acceptance Criteria
- [ ] `AsyncEvidence` struct defined with all fields
- [ ] `AsyncEvidenceContract` used in all async endpoints
- [ ] Success/error codes match taxonomy
- [ ] Delay categories map correctly to UI badges
- [ ] Benchmark collector records all three metrics
- [ ] PKI network delay included when mTLS active
- [ ] Frontend can render delay breakdown from evidence