# CommLayers FSD + SDS - Backend-First Revision

## 1. Functional Intent

The platform must let users inspect how each backend communication family behaves, when it is appropriate, how it transitions between states, and how it compares to adjacent alternatives under realistic timing and observability rules.

## 2. Primary Functional Streams

### 2.1 Backend Family Catalog

The system shall expose a normalized catalog of backend communication families, abstraction layers, and common application patterns.

### 2.2 Comparative Execution

The system shall run comparative implementations that let users observe:
- protocol response time
- state transitions
- delays and timeouts
- reliability and fallback behavior
- operational consequences

### 2.3 Frontend Consumption

The frontend shall consume normalized backend responses and visualize comparisons without becoming the source of truth.

Frontend auth should prefer a BFF-oriented session model so frontend code does not hold persistent OAuth tokens unless a later direct-API exception is explicitly justified.

### 2.4 Localization

Every user-facing backend response shown in the frontend shall support English default and Spanish localization.

## 3. Required Response Model

Every comparable backend response returned to the frontend should support fields equivalent to:

```json
{
  "family": "sse",
  "normalized_name": "server_sent_events",
  "common_use_cases": ["progress_updates", "notifications", "streaming_status"],
  "recommended_when": "server_to_client_updates",
  "avoid_when": "high_frequency_bidirectional_exchange",
  "latency": {
    "p50_ms": 0,
    "p95_ms": 0,
    "p99_ms": 0
  },
  "state_transition": {
    "current": "running",
    "delay_applied_ms": 0,
    "timeout_threshold_ms": 10000
  },
  "evidence": {
    "source": "benchmark_run",
    "confidence": "high",
    "staleness": "fresh",
    "measured_at": "2026-03-22T00:00:00Z"
  },
  "i18n": {
    "locale": "en",
    "available_locales": ["en", "es"]
  }
}
```

## 4. Comparative Visualization Rules

The frontend must later support:
- side-by-side protocol comparisons
- explicit protocol recommendations and anti-recommendations
- visible user-observable delays where necessary to show state transitions
- confidence and staleness chips
- localized labels and backend-provided explanatory text

The backend must provide enough timing metadata so the frontend can intentionally slow non-human-visible transitions when educational visualization is required.

## 5. Common Application Matrix

Each communication family must be classified against common applications such as:
- simple CRUD request-response
- periodic status checking
- blocking wait with eventual completion
- one-way live updates
- bidirectional interaction
- webhook-driven integration
- fan-out event delivery
- durable workflow execution
- immutable event history
- continuous stream processing
- actor-style isolated concurrency
- realtime document or state sync
- replication and edge merge
- distributed coordination
- service networking and control-plane management

## 6. Long-Running Operation Behavior

If expected completion exceeds 10s:
- submit path returns `202 Accepted`
- operation resource is created
- client can poll or subscribe to status stream
- localized state labels are returned
- terminal state includes result or failure detail suitable for frontend display

## 7. Realtime Selection Behavior

The system must support both SSE and WebSocket, but recommendation logic must prefer:
- SSE for one-way update streams
- WebSocket only when interaction requires it

The comparison output must explicitly say why one was preferred over the other for a given scenario.

The comparison output must also indicate when a protocol is only being shown as a compatibility fallback rather than the recommended default.

## 8. Data and Projection Behavior

### 8.1 Canonical Truth

Canonical truth must begin in PostgreSQL-backed normalized domain models.

### 8.2 Derived Projections

Derived projections may exist for:
- search in Elasticsearch
- topology queries in graph DB when graph traversal is approved as first-class
- semantic lookup in `pgvector` first, then a dedicated vector DB only if operationally justified

### 8.3 Optional Document Store

MongoDB must not be introduced unless a documented v1 use case proves it is necessary.

## 9. Testing and TDD Requirements

Each stream must include:
- red test or benchmark first
- green minimal pass
- refactor pass
- regression suite for normalized comparison outputs

Mandatory test categories:
- unit
- integration
- protocol contract
- benchmark capture
- security and auth behavior
- localization response correctness

## 10. Stream-Level Delivery Acceptance

### Stream acceptance is not complete unless:
- implementation works
- comparative output is normalized
- benchmark evidence exists
- frontend-facing localization works
- delays and timeouts are visible where required
- security and audit hooks exist

## 11. Future v2 Separation

Architectural paradigms are future work and must remain separate from the v1 backend-family execution path, even when some groundwork is implemented now.
