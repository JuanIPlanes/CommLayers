# DEC-006: Normalized Response Model

## Status
Accepted

## Date
2026-03-22

## Context
All comparable backend responses must follow a consistent schema that includes latency metrics, state transitions, evidence metadata, and localization support.

## Decision
Standardize on the familyResponse schema for all comparable backend responses:
```json
{
  "family": "sse",
  "normalized_name": "server_sent_events",
  "common_use_cases": ["progress_updates", "notifications"],
  "recommended_when": "server_to_client_updates",
  "avoid_when": "high_frequency_bidirectional_exchange",
  "latency": {"p50_ms": 12, "p95_ms": 45, "p99_ms": 120},
  "state_transition": {"current": "available", "delay_applied_ms": 0, "timeout_threshold_ms": 10000},
  "evidence": {"source": "benchmark_run", "confidence": "high", "staleness": "fresh", "measured_at": "2026-03-22T00:00:00Z"},
  "i18n": {"locale": "en", "available_locales": ["en", "es"]}
}
```

## Consequences
- Frontend can reliably consume and visualize any family response
- Benchmark evidence is embedded in responses
- Localization is built into the response model
- State transition metadata enables educational visualization

## Related
- FSD+SDS Section 3: Required Response Model
- FSD+SDS Section 4: Comparative Visualization Rules
- Implemented in handlers: handleCatalog, handleTransportSummary, handleRealtimeComparisons