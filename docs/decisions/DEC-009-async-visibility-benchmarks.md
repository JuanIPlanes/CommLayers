---
id: DEC-009
title: Async Visibility Benchmarks
status: proposed
date: 2026-11-11
---

## Context
The frontend needs observable metrics to display async operation progress, delays, and state transitions. The `AsyncVisibilityBenchmarks` struct in `backend/benchmarks.go` defines the contract for these metrics. This decision formalizes it as a binding specification.

## Decision
The AsyncVisibilityBenchmarks contract defines three metric categories with percentile thresholds (p50/p95/p99):

### Queue Delay Metrics
| Metric | Type | Description |
|--------|------|-------------|
| `AverageDelayMs` | float64 | Arithmetic mean of queue delay samples |
| `P95DelayMs` | float64 | 95th percentile queue delay |
| `P99DelayMs` | float64 | 99th percentile queue delay |
| `MaxDelayMs` | float64 | Maximum observed queue delay |
| `MinDelayMs` | float64 | Minimum observed queue delay |
| `DelaySamples` | int | Number of samples in the measurement window |
| `LastMeasuredAt` | time.Time | Timestamp of last measurement |

### Timeout Metrics
| Metric | Type | Description |
|--------|------|-------------|
| `AverageTimeoutMs` | float64 | Arithmetic mean of timeout durations |
| `P95TimeoutMs` | float64 | 95th percentile timeout |
| `P99TimeoutMs` | float64 | 99th percentile timeout |
| `MaxTimeoutMs` | float64 | Maximum observed timeout |
| `MinTimeoutMs` | float64 | Minimum observed timeout |
| `TimeoutSamples` | int | Number of timeout samples |
| `LastMeasuredAt` | time.Time | Timestamp of last measurement |
| `EnforcedTimeoutMs` | float64 | Configured timeout limit (contract boundary) |

### State Transition Metrics
| Metric | Type | Description |
|--------|------|-------------|
| `AverageTransitionMs` | float64 | Arithmetic mean of transition durations |
| `P95TransitionMs` | float64 | 95th percentile transition duration |
| `P99TransitionMs` | float64 | 99th percentile transition duration |
| `MaxTransitionMs` | float64 | Maximum observed transition duration |
| `MinTransitionMs` | float64 | Minimum observed transition duration |
| `TransitionSamples` | int | Number of transition samples |
| `LastMeasuredAt` | time.Time | Timestamp of last measurement |
| `AllowedTransitions` | []string | List of permitted state transitions (e.g., "queued→running", "running→completed") |
| `BlockedTransitions` | []string | List of currently blocked transitions |

### Collection Configuration
- `QueueDelaySampleInterval`: 1 second (default)
- `TimeoutSampleInterval`: 5 seconds (default)
- `StateTransitionSampleInterval`: 1 second (default)
- `MaxSampleHistory`: 1000 samples per metric (ring buffer)
- `EnabledMetrics`: ["queue_delay", "timeout", "state_transitions"] (all enabled by default)

### Evidence Integration
Each benchmark run produces an `EvidenceContract` (from `evidence/contracts.go`):
- Source: `SourceBenchmark`
- Confidence: 1.0
- Staleness: `StalenessFresh`
- MeasuredAt: sample timestamp
- DelayAppliedMs: 0 (benchmark runs are synchronous)
- SuccessCode: `SuccessOK` for successful runs
- ErrorCode: from fixed taxonomy (timeout, unavailable, internal_error, etc.)

### Frontend Consumption Contract
The frontend consumes `/api/benchmarks/metrics` which returns:
```json
{
  "queue_delay": { "p50": 120, "p95": 450, "p99": 1200, "samples": 847 },
  "timeout": { "p50": 5000, "p95": 15000, "p99": 28000, "enforced": 30000 },
  "state_transitions": { "p50": 200, "p95": 800, "p99": 2100, "allowed": ["queued→running", "running→completed"], "blocked": [] }
}
```

## Consequences
- Benchmark collector uses ring buffers (max 1000 samples); older samples are dropped
- Percentile calculations use the Go standard library approach (nearest-rank)
- `EnforcedTimeoutMs` is a contract boundary — operations exceeding it MUST return `ErrorTimeout`
- `AllowedTransitions`/`BlockedTransitions` are derived from the async state machine (LRO, workflow, transport)
- All metrics are tagged with `stream_id`, `job_id`, `phase` for correlation

## Acceptance Criteria
- [ ] `BenchmarkCollector` produces `BenchmarkSummary` with correct percentile calculations
- [ ] Integration test: inject known samples, verify p50/p95/p99 match expected values
- [ ] Evidence contracts attach to every `BenchmarkRun` and `BenchmarkSummary`
- [ ] Frontend can render delay/timeout/transition charts from `/api/benchmarks/metrics`
- [ ] `EnforcedTimeoutMs` breach returns `ErrorTimeout` with `DelayAppliedMs` populated
- [ ] Ring buffer eviction verified: 1001st sample drops the oldest