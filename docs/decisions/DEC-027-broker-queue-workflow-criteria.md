---
id: DEC-027
title: Broker Queue and Workflow Criteria
status: proposed
date: 2026-11-11
---

## Context
The messaging workflow family (stream 86agbv0ka) and failure modes (86agbv0nt - done) define broker queue semantics, retry policies, visibility, and failure mode integration. This decision formalizes the criteria as a binding contract.

## Decision
### Broker Criteria (from `messaging/failure.go` and `main.go` handleMessagingOverview)

#### 1. Queue Semantics
| Property | Workflow Engine | Job Queue | Broker Event Bus |
|----------|-----------------|-----------|------------------|
| **Delivery Guarantee** | At-least-once | At-least-once | At-least-once (exactly-once future) |
| **Ordering** | Per-workflow FIFO | Per-queue FIFO | Partitioned by key |
| **Persistence** | PostgreSQL (WAL) | PostgreSQL + Redis | Kafka/Redpanda (future) |
| **Priority** | 3 levels (high/normal/low) | 5 levels | Partition-based |
| **TTL** | 24h default | 7d default | 7d default |
| **Max Size** | 10k messages | 100k messages | 1M messages/partition |

#### 2. Retry Policy
```go
type RetryPolicy struct {
    MaxAttempts:        3,           // Workflow, Job Queue
    InitialBackoffMs:   1000,        // 1s
    MaxBackoffMs:       30000,       // 30s
    BackoffMultiplier:  2.0,         // Exponential
    Jitter:             0.1,         // ±10%
    RetryableErrors:    ["timeout", "unavailable", "internal_error", "rate_limited"],
    NonRetryableErrors: ["invalid_input", "auth_failed", "not_found", "conflict", "stale_data"],
}
```

#### 3. Visibility Contract
Every queued/running operation exposes:
```json
{
  "queue_position": 5,
  "estimated_start_ms": 12000,
  "worker_id": "worker-3",
  "stage": "stage:2",
  "progress": 45,
  "delay_applied_ms": 1200,
  "evidence": { ... }
}
```

**Visibility Guarantees**:
- Queue position accurate within ±1 (eventual consistency)
- Estimated start based on `p50` processing time × position
- Worker ID visible for debugging
- Progress 0-100 for multi-stage operations

### Workflow Contract
**Workflow Stages** (from `handleWorkflowCreate`/`handleWorkflowStatus`):
```
queued → initializing → stage:1 → stage:2 → ... → stage:N → completing → completed
                                    ↓
                              failed (on error)
                                    ↓
                              retry:N (if retries remain)
                                    ↓
                              dead_letter (if retries exhausted)
```

**Stage Definition**:
```go
type WorkflowStage struct {
    Name:           string,          // e.g., "validate", "process", "notify"
    TimeoutMs:      30000,           // Per-stage timeout
    RetryPolicy:    RetryPolicy,     // May override workflow default
    IdempotencyKey: string,          // Required for retry safety
    Compensation:   string,          // Compensation action on failure
}
```

**Workflow Definition** (persisted):
```go
type WorkflowDefinition struct {
    ID:             string,
    Name:           string,
    Version:        int,
    Stages:         []WorkflowStage,
    Concurrency:    int,             // Max parallel executions
    TimeoutMs:      int,             // Total workflow timeout
    CreatedAt:      time.Time,
    UpdatedAt:      time.Time,
}
```

### Failure Mode Integration (from `messaging/failure.go` and 86agbv0nt done)
**Failure Mode Matrix** (per family):

#### Workflow Engine (`workflow_engine`)
| Failure Mode | Warning Threshold | Critical Threshold | Remediation |
|--------------|-------------------|-------------------|-------------|
| `queue_exhaustion` | 80% capacity | 95% capacity | Scale workers, increase capacity, backpressure |
| `message_timeout` | 5000ms | 15000ms | Investigate stage, add timeout handling, async decomposition |
| `retry_exhausted` | 2 retries | 3 retries (exhausted) | Dead letter, alert, review idempotency |
| `dead_letter` | 1 message | 10 messages | Analyze patterns, fix root cause, replay if safe |
| `circuit_open` | 1 circuit | 1 circuit | Investigate downstream, fallback, manual override |
| `idempotency_violation` | 1 duplicate | 5 duplicates | Verify producer keys, check retry storms |

#### Job Queue (`job_queue`)
| Failure Mode | Warning Threshold | Critical Threshold | Remediation |
|--------------|-------------------|-------------------|-------------|
| `queue_exhaustion` | 75% | 90% | Auto-scale, priority queuing, load shedding |
| `message_timeout` | 30000ms | 60000ms | Optimize job, increase timeout, decompose |
| `retry_exhausted` | 3 retries | 5 retries | Dead letter, alert, analyze pattern |
| `payload_too_large` | 1MB | 5MB | External storage (S3), compress, chunk |
| `serialization_failure` | 1 error | 10 errors | Validate schema, check version, fix producer |

#### Broker Event Bus (`broker_event_bus`) — Future
| Failure Mode | Warning Threshold | Critical Threshold | Remediation |
|--------------|-------------------|-------------------|-------------|
| `broker_unavailable` | 1 broker down | Quorum lost | Failover, alert |
| `ordering_violation` | 1 violation | 10 violations | Check partition key, producer config |
| `consumer_lag` | 1000 messages | 10000 messages | Scale consumers, investigate processing |

### Failure Check Execution (from `RunFailureChecks`)
```go
func RunFailureChecks(matrix *FailureModeMatrix, metrics FailureMetrics) []FailureCheck {
    // Runs every 30s for active workflows/queues
    // Returns FailureCheck[] with Healthy, Severity, Details
    // Integrates with EvidenceContract:
    //   Healthy=true  → SuccessOK
    //   Severity=warning → SuccessDegraded
    //   Severity=critical → ErrorInternal (or specific)
}
```

**FailureMetrics** (runtime collected):
```go
type FailureMetrics struct {
    QueueDepthPercent:       float64,  // 0-100
    MaxStageLatencyMs:       float64,
    MaxRetriesAttempted:     int,
    DeadLetterCount:         int,
    OpenCircuitCount:        int,
    DuplicateCount:          int,
    MaxPayloadSizeBytes:     float64,
    SerializationErrorCount: int,
}
```

### Demo Metrics (from `DemoFailureMetrics`)
**Healthy Demo**:
```json
{
  "queue_depth_percent": 45,
  "max_stage_latency_ms": 2800,
  "max_retries_attempted": 1,
  "dead_letter_count": 0,
  "open_circuit_count": 0,
  "duplicate_count": 0,
  "max_payload_size_bytes": 51200,
  "serialization_error_count": 0
}
```

**Degraded Demo** (for testing failure modes):
```json
{
  "queue_depth_percent": 87,      // warning
  "max_stage_latency_ms": 18000,  // critical
  "max_retries_attempted": 3,     // critical
  "dead_letter_count": 5,         // warning
  "open_circuit_count": 1,        // critical
  "duplicate_count": 2,           // warning
  "max_payload_size_bytes": 2048000, // warning
  "serialization_error_count": 3  // warning
}
```

### API Contract
**`GET /api/messaging/failure-modes`**: Returns both matrices
**`GET /api/messaging/failure-modes/workflow`**: Workflow matrix + current checks
**`GET /api/messaging/failure-modes/job-queue`**: Job queue matrix + current checks
**`POST /api/messaging/workflows`**: Creates workflow run with failure mode tracking
**`GET /api/messaging/workflows/{id}`**: Returns workflow with current failure checks

### Evidence Integration
All workflow/job operations produce `EvidenceContract`:
- `source`: `SourceAudit`
- `confidence`: 1.0 for completed, 0.8 for running, 0.5 for queued
- `staleness`: From last status update
- `measured_at`: Status check timestamp
- `delay_applied_ms`: Queue + processing delay

## Consequences
- Broker criteria are binding: implementation must match semantics
- Retry policy is fixed: changes require gate 86agbv0ka re-approval
- Visibility contract enables frontend progress UI
- Failure modes are exhaustive: no new modes without matrix update
- Demo metrics provide known-good baseline for testing

## Acceptance Criteria
- [ ] Queue semantics match specification for all three families
- [ ] Retry policy implemented with exact parameters
- [ ] Visibility fields present on all workflow/job responses
- [ ] Failure mode matrices match `messaging/failure.go` exactly
- [ ] `RunFailureChecks` produces correct `FailureCheck` results
- [ ] Demo metrics produce expected healthy/degraded results
- [ ] API endpoints return correct failure mode data
- [ ] Evidence contracts on all messaging operations