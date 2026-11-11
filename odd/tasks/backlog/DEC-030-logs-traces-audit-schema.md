---
id: DEC-030
title: Logs, Traces, and Audit Schema
status: proposed
date: 2026-11-11
---

## Context
Extends evidence contracts (86agbv0ne - done) and PKI audit (86agbv0ng - done) with unified log format, trace context, and audit event schema. This decision formalizes the observability schema.

## Decision
### Log Format Contract
**Structured Log Entry** (JSON, one per line):
```json
{
  "timestamp": "2026-11-11T10:30:00.123Z",
  "level": "info|warn|error|debug",
  "logger": "commlayers.api.workflow",
  "message": "Workflow stage completed",
  "fields": {
    "request_id": "req-abc123",
    "stream_id": "86agbv0ka",
    "job_id": "workflow-456",
    "phase": "stage:2",
    "result": "success",
    "duration_ms": 1250,
    "delay_applied_ms": 1250
  },
  "evidence": {
    "source": "audit",
    "confidence": 1.0,
    "staleness": "fresh",
    "measured_at": "2026-11-11T10:30:00.123Z",
    "delay_applied_ms": 1250
  }
}
```

**Required Fields** (every log line):
| Field | Type | Description |
|-------|------|-------------|
| `timestamp` | RFC3339Nano | Event timestamp |
| `level` | string | info/warn/error/debug |
| `logger` | string | Hierarchical logger name |
| `message` | string | Human-readable message |
| `fields.request_id` | string | Correlation ID (UUID) |
| `fields.stream_id` | string | First-wave stream ID |
| `fields.job_id` | string | Operation ID (workflow, session, etc.) |
| `fields.phase` | string | Current phase/state |
| `fields.result` | string | success/partial/degraded/failed |
| `evidence` | object | EvidenceContract (DEC-011) |

**Log Levels**:
- `debug`: Verbose diagnostic (queue position, cache hits)
- `info`: State transitions, completed operations
- `warn`: Degraded performance, retry attempts, threshold breaches
- `error`: Failures, circuit open, data loss

### Trace Context Contract
**Trace Header Propagation** (W3C TraceContext compatible):
```
traceparent: 00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01
tracestate: commlayers=stream_id:86agbv0ka,job_id:workflow-456
```

**Trace Fields** (in log `fields`):
| Field | Description |
|-------|-------------|
| `trace_id` | W3C trace-id (16 bytes hex) |
| `span_id` | W3C span-id (8 bytes hex) |
| `trace_flags` | W3C trace-flags (sampled=01) |
| `parent_span_id` | Parent span for child spans |

**Named Traces** (from DEC-011):
| Trace Name | Scope | When Created |
|------------|-------|--------------|
| `bootstrap.http` | HTTP request/response | Every HTTP request |
| `bootstrap.async` | Async job lifecycle | Job create → complete |
| `bootstrap.sse` | SSE connection lifecycle | Connect → disconnect |

### Audit Event Schema (extends PKI audit from 86agbv0ng)
**Base AuditEvent** (from `keycloak.go` and `pki_observability.go`):
```go
type AuditEvent struct {
    ID        string                 `json:"id"`           // audit-20261111103000.123456
    Timestamp time.Time              `json:"timestamp"`    // RFC3339
    UserID    string                 `json:"user_id,omitempty"`
    SessionID string                 `json:"session_id,omitempty"`
    Action    string                 `json:"action"`       // HTTP method or internal action
    Resource  string                 `json:"resource"`     // URL path or resource type
    Method    string                 `json:"method"`       // HTTP method or INTERNAL
    Path      string                 `json:"path"`         // Full request path
    IP        string                 `json:"ip"`           // Client IP
    UserAgent string                 `json:"user_agent"`   // Client UA
    Status    int                    `json:"status"`       // HTTP status or 0 for internal
    Metadata  map[string]interface{} `json:"metadata,omitempty"`
}
```

**Extended Metadata** (per event type):
| Event Type | Metadata Fields |
|------------|-----------------|
| `pki.certificate_recorded` | serial, subject, issuer, not_after, source, service |
| `pki.mtls_connection_observed` | source, target, tls_version, cipher_suite, verified, cert_expiry_days, phase |
| `pki.trust_boundary_evaluated` | name, compliant, total_pairs, compliant_pairs, violations |
| `pki.mtls_phase_changed` | new_phase |
| `workflow.stage_completed` | workflow_id, stage, duration_ms, retries |
| `sync.session_synced` | session_id, lag_before, lag_after, conflict_resolved |
| `async.job_completed` | job_id, duration_ms, delay_applied_ms, evidence_confidence |
| `pricing.rate_refreshed` | pair, source, old_rate, new_rate, confidence |

### Log/Trace/Audit Integration
**Middleware Chain** (from `main.go`):
```
1. CORS
2. Logging (request_id, duration, status) → LogFormat
3. Audit (AuditMiddleware) → AuditEvent
4. PKI (PKIMiddleware) → PKI AuditEvent (if PKI path)
5. Router → Handler
```

**Correlation**: All three share `request_id` / `trace_id` for end-to-end tracing.

### Output Targets
| Output | Format | Destination |
|--------|--------|-------------|
| Logs | JSON lines | stdout (container), Loki (production) |
| Traces | OTLP | Jaeger/Tempo (production) |
| Audit | JSON lines | stdout (container), dedicated audit log (production) |

### Sampling
| Component | Sample Rate |
|-----------|-------------|
| HTTP requests | 100% (info), 100% (error) |
| Async operations | 100% |
| PKI events | 100% |
| Debug logs | 10% (configurable) |
| Traces | 100% for errors, 10% for success |

### Retention
| Data Type | Retention |
|-----------|-----------|
| Logs | 30 days (hot), 1 year (cold) |
| Traces | 7 days |
| Audit | 7 years (compliance) |
| PKI Audit | 7 years (compliance) |

## Consequences
- Structured logging is mandatory: no unstructured log lines
- Every log line carries evidence contract
- Trace context propagated through all async boundaries
- Audit events cover security + operational events
- Sampling reduces volume but keeps 100% for errors/audit

## Acceptance Criteria
- [ ] Log format matches specification exactly (verified by log parser test)
- [ ] Trace headers propagated through HTTP, async, SSE
- [ ] Audit events generated for all PKI + operational events
- [ ] Correlation IDs link logs, traces, audit
- [ ] Sampling rates configurable and respected
- [ ] Retention policies documented
- [ ] Integration test verifies full chain