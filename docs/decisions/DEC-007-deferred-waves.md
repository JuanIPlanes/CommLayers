# DEC-007: Deferred Wave Architecture

## Status
Accepted

## Date
2026-03-22

## Context
Not all backend families can be implemented in the first wave. A structured deferral mechanism is needed to keep bootstrap focused while providing a clear path for future implementation.

## Decision
Organize backend families into waves with explicit gating:
- Wave 1 (First Wave): Abstraction layer, Security bootstrap, Canonical data platform, Benchmark framework
- Wave 2 (Deferred): Transport implementations, Messaging/workflow, Sync/replication, Projections
- Wave 3 (Deferred): Localization, Frontend comparative visualization, Pricing
- Wave 4 (Future): v2 architectural paradigms

Gating rules:
- No later-wave implementation starts before first-wave metrics and closure report format are fixed
- First-wave closure (benchmark acceptance) unlocks Wave 2
- Backend-family stability unlocks Wave 3
- v1 maturity unlocks Wave 4

## Consequences
- Prevents scope creep in bootstrap
- Provides clear implementation roadmap
- Ensures foundational work is solid before building on top
- Aligns with FSD+SDS Section 11: Future v2 Separation

## Related
- FSD+SDS Section 11: Future v2 Separation
- /api/deferred-waves endpoint
- Stream categories: first-wave, deferred-wave-2, deferred-wave-3, deferred-wave-4