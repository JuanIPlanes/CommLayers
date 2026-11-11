---
id: DEC-026
title: Visible Delay and State Transition Interactions
status: proposed
date: 2026-11-11
---

## Context
Frontend needs delay metadata from benchmarks (DEC-009, DEC-011) and PKI observability (DEC-004) to visualize state transitions. This decision formalizes the delay metadata contract for frontend consumption.

## Decision
### Delay Metadata Contract
Every async operation response MUST include delay metadata:

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
  "evidence": { ... }
}
```

### Delay Breakdown Fields
| Field | Source | Description |
|-------|--------|-------------|
| `queue_delay_ms` | Benchmarks (DEC-009) | Time waiting in async queue |
| `processing_delay_ms` | Benchmarks (DEC-009) | Actual processing time |
| `network_delay_ms` | PKI Observability (DEC-004) | mTLS handshake + network latency |
| `delay_applied_ms` | Sum of above | Total user-visible delay |

### State Transition Visualization Rules
**Allowed Transitions** (from `AsyncVisibilityBenchmarks.AllowedTransitions`):
| Family | Transitions |
|--------|-------------|
| Transport (polling) | `created → running → completed` |
| Transport (long-polling) | `created → running → completed` |
| Transport (SSE) | `connecting → connected → event → completed` |
| Transport (WebSocket) | `connecting → connected → message → completed` |
| Workflow | `queued → running → stage:N → completed|failed` |
| Async LRO | `pending → in_progress → completed|failed` |
| Sync | `created → active → aligned|conflict → resolved → closed` |

**Blocked Transitions** (from `AsyncVisibilityBenchmarks.BlockedTransitions`):
- Any transition not in allowed list
- Transitions gated by unmet prerequisites (e.g., consensus blocked by lag)
- Transitions during circuit open (failure mode)

### Frontend Visualization Contract
**Delay Badge Component**:
```typescript
interface DelayBadgeProps {
  delayMs: number;
  thresholdMs: number;
  breakdown?: DelayBreakdown;
  size: 'sm' | 'md' | 'lg';
}

// Visual rules:
// delayMs < 0.5 * thresholdMs → Green ("Fast")
// delayMs < thresholdMs → Yellow ("Normal")
// delayMs ≥ thresholdMs → Red ("Slow")
// delayMs ≥ 2 * thresholdMs → Red + Pulse ("Critical")
```

**State Transition Diagram**:
- Nodes = states, Edges = allowed transitions
- Edge label = `transition_ms` (p50 from benchmarks)
- Blocked transitions shown as dashed red edges
- Current state highlighted
- Animation on transition (300ms)

**Timeline Component**:
```typescript
interface TimelineEvent {
  timestamp: string;        // RFC3339
  state: string;            // Current state
  delayMs: number;          // delay_applied_ms
  evidence: EvidenceContract;
  metadata?: {
    queueDepth?: number;
    pkiPhase?: number;      // From PKI observability
    tlsHandshakeMs?: number;
  };
}
```
- Horizontal timeline with events as dots
- Hover shows delay breakdown + evidence
- Color by state (green=success, yellow=pending, red=error)
- Zoom/pan for long-running operations

### PKI Observability Integration (from DEC-004)
**Delay Metadata from PKI**:
| Field | Source | When Present |
|-------|--------|--------------|
| `tls_handshake_ms` | `MTLSConnection.handshake_duration_ms` | mTLS connections (phase 2+) |
| `cert_verification_ms` | `MTLSConnection.cert_verification_duration_ms` | mTLS connections |
| `pki_phase` | `MTLSStatusTracker.GetPhase()` | Always |

**Frontend PKI Indicators**:
- Phase 1: No delay indicators (identity + ingress TLS only)
- Phase 2: Show `tls_handshake_ms` for critical service pairs
- Phase 3: Show `tls_handshake_ms` + `cert_verification_ms` for all internal calls

### Benchmark Integration (from DEC-009, DEC-011)
**Thresholds for Visualization**:
| Metric | Warning | Critical | Source |
|--------|---------|----------|--------|
| `queue_delay_ms` | p95 > 2000ms | p99 > 5000ms | `QueueDelayMetrics` |
| `timeout_ms` | > 15000ms | > 30000ms | `TimeoutMetrics.EnforcedTimeoutMs` |
| `state_transition_ms` | p95 > 3000ms | p99 > 5000ms | `StateTransitionMetrics` |

**Real-time Updates**:
- Frontend polls `/api/benchmarks/metrics` every 30s
- WebSocket push for live delay updates (when WebSocket transport active)
- SSE fallback for delay metadata stream

### API Contract
**`GET /api/async/demo/{id}/delays`**:
```json
{
  "session_id": "async-123",
  "current_delay_ms": 1250,
  "delay_history": [
    {"timestamp": "2026-11-11T10:30:00Z", "delay_ms": 1100, "state": "queued"},
    {"timestamp": "2026-11-11T10:30:01Z", "delay_ms": 1250, "state": "running"}
  ],
  "thresholds": {
    "warning_ms": 2000,
    "critical_ms": 5000
  },
  "pki_metadata": {
    "phase": 2,
    "tls_handshake_ms": 45,
    "critical_pairs": ["auth-service", "vault"]
  }
}
```

**`GET /api/messaging/workflows/{id}/delays`**: Similar for workflow runs

### Evidence Integration
All delay responses include `EvidenceContract`:
- `source`: `SourceBenchmark` (for benchmark-derived) or `SourceAudit` (for runtime)
- `confidence`: 1.0 for measured, 0.8 for estimated
- `staleness`: `StalenessFresh` if < 60s old
- `measured_at`: Timestamp of measurement
- `delay_applied_ms`: The delay being reported (meta-delay)

## Consequences
- Delay metadata is mandatory on ALL async responses
- Breakdown fields enable root-cause analysis in UI
- PKI phase gates which delay fields are shown
- Benchmark thresholds drive visual warning levels
- Evidence contracts enable audit of delay measurements

## Acceptance Criteria
- [ ] All async endpoints return `delay_applied_ms` and `delay_breakdown`
- [ ] State transition visualization matches allowed/blocked lists
- [ ] Delay badge component implements visual rules correctly
- [ ] Timeline component renders with hover details
- [ ] PKI metadata shown correctly per phase
- [ ] Benchmark thresholds drive warning/critical colors
- [ ] Evidence contracts on all delay responses