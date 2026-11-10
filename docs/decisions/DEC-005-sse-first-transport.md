# DEC-005: SSE-First Transport Selection

## Status
Accepted

## Date
2026-03-22

## Context
The platform must recommend transport protocols based on use case, with clear reasoning for each recommendation.

## Decision
Prefer SSE for one-way updates, WebSocket only when bidirectional interaction is required:
1. SSE is the default push choice for one-way updates (progress streams, status feeds)
2. WebSocket is planned but not bootstrap default - only used when duplex is actually needed
3. Comparison output must explicitly state why one was preferred over another
4. Protocols shown as compatibility fallback must be clearly marked
5. Long-polling and polling are fallback options only

## Consequences
- Simplifies infrastructure (SSE works over standard HTTP)
- Clear decision criteria for transport selection
- Frontend can rely on SSE as default push mechanism
- WebSocket complexity only incurred when justified

## Related
- FSD+SDS Section 7: Realtime Selection Behavior
- /api/comparisons/realtime endpoint
- /api/transports endpoint