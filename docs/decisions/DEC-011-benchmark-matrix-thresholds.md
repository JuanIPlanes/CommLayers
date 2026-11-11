---
id: DEC-011
title: Benchmark Matrix and Thresholds
status: proposed
date: 2026-11-11
---

## Context
The benchmark framework (`backend/benchmarks.go` and `main.go` handleBenchmarkFramework) defines three benchmark scenarios with specific metrics and thresholds. This decision formalizes the benchmark matrix, success/error taxonomy, and thresholds as binding contracts.

## Decision
### Benchmark Scenarios

#### 1. Request Baseline
Measures: `request_duration_ms` (p50, p95, p99), `error_rate`
- **Thresholds**:
  - p50 ≤ 100ms
  - p95 ≤ 300ms
  - p99 ≤ 800ms
  - error_rate ≤ 0.1% (SuccessOK), ≤ 1% (SuccessPartial), > 1% (SuccessDegraded)
- **Evidence Source**: `SourceBenchmark`
- **Run Configuration**: 1000 requests, 10 concurrent, warm-up 100

#### 2. SSE Startup and Event Cadence
Measures: `stream_startup_latency_ms`, `event_flush_latency_ms`, `delay_metadata`
- **Thresholds**:
  - stream_startup_latency p50 ≤ 200ms, p95 ≤ 500ms
  - event_flush_latency p50 ≤ 50ms, p95 ≤ 150ms
  - delay_metadata: `DelayAppliedMs` present in 100% of events
- **Evidence Source**: `SourceBenchmark`
- **Run Configuration**: 100 streams, 10 events each, 1s interval

#### 3. Async Visibility
Measures: `queue_delay_ms` (p50, p95, p99), `step_transition_cadence_ms`, `timeout_threshold_visibility`
- **Thresholds**:
  - queue_delay p50 ≤ 500ms, p95 ≤ 2000ms, p99 ≤ 5000ms
  - step_transition_cadence p50 ≤ 1000ms, p95 ≤ 3000ms
  - timeout_threshold: `EnforcedTimeoutMs` visible in 100% of responses
- **Evidence Source**: `SourceBenchmark`
- **Run Configuration**: 500 async jobs, varied payload sizes

### Success/Error Taxonomy (from `evidence/contracts.go`)

#### Success Codes (fixed, no additions without gate re-approval)
| Code | Meaning | When Used |
|------|---------|-----------|
| `ok` | Fully successful | All thresholds met, no degradation |
| `partial` | Partially successful | Some thresholds met, minor degradation (e.g., p99 exceeded but p95 ok) |
| `degraded` | Degraded but functional | Significant thresholds exceeded, but service operational |

#### Error Codes (fixed, no additions without gate re-approval)
| Code | HTTP Status | When Used |
|------|-------------|-----------|
| `timeout` | 504 / 408 | Operation exceeded `EnforcedTimeoutMs` |
| `unavailable` | 503 | Dependency unavailable (DB, Redis, broker) |
| `invalid_input` | 400 | Request validation failed |
| `auth_failed` | 401 | Authentication/authorization failed |
| `rate_limited` | 429 | Rate limit exceeded (includes `Retry-After`) |
| `internal_error` | 500 | Unexpected internal failure |
| `not_found` | 404 | Resource not found |
| `conflict` | 409 | Concurrency conflict (optimistic lock, idempotency) |
| `stale_data` | 409 | Evidence staleness = `very_stale` or `reindexing` |

### Telemetry Naming Contract
| Category | Names (exact) |
|----------|---------------|
| Metrics | `request_duration_ms`, `event_flush_ms`, `job_progress`, `delay_applied_ms`, `queue_delay_ms`, `timeout_ms`, `state_transition_ms` |
| Logs | `request_id`, `stream_id`, `job_id`, `phase`, `result` |
| Traces | `bootstrap.http`, `bootstrap.async`, `bootstrap.sse` |

### Closure Rules (Executable)
```go
func ClosureRulesPassed() error {
    // 1. No later-wave implementation starts before first-wave metrics fixed
    if deferredWaveImplementationStarted() && !firstWaveMetricsFixed() {
        return fmt.Errorf("later-wave started before first-wave metrics fixed")
    }
    // 2. User-visible delay metadata available for frontend
    if !delayMetadataAvailableForFrontend() {
        return fmt.Errorf("delay metadata not available for frontend")
    }
    // 3. Success/error taxonomy stable across first-wave endpoints
    if !taxonomyStableAcrossEndpoints() {
        return fmt.Errorf("taxonomy drift detected")
    }
    return nil
}
```

## Consequences
- Thresholds are binding: CI fails if benchmarks exceed thresholds
- Taxonomy is frozen: new success/error codes require gate 86agbv0k9 re-approval
- Telemetry names are exact strings: no aliases, no variations
- Frontend consumes delay metadata via `DelayAppliedMs` field in all async responses
- Wave closure report template includes: scenario results, threshold compliance, taxonomy version, telemetry schema version

## Acceptance Criteria
- [ ] Benchmark runner produces results with p50/p95/p99 for all three scenarios
- [ ] CI pipeline runs benchmarks and fails on threshold breach
- [ ] Evidence contracts attach to all benchmark runs with correct taxonomy
- [ ] `/api/benchmarks/metrics` returns telemetry with exact naming contract
- [ ] Closure rules enforced in CI: later-wave tasks blocked until gates pass
- [ ] Taxonomy version embedded in all evidence contracts