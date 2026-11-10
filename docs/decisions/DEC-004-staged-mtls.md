# DEC-004: Staged mTLS

## Status
Accepted

## Date
2026-03-22

## Context
Mutual TLS should be introduced in phases to avoid blocking bootstrap while still providing a path to full east-west encryption.

## Decision
Implement mTLS in three phases:
1. Phase 1: Identity and ingress TLS - Prepare service identity and trust boundary inventory
2. Phase 2: Selective mTLS - Protect critical service pairs without blocking bootstrap
3. Phase 3: Strict east-west mTLS - Enforce workload identity and authorization after observability is ready

## Consequences
- Allows bootstrap to proceed without full mTLS complexity
- Provides clear progression path
- Critical services can be protected early
- Full enforcement delayed until observability is proven

## Related
- FSD+SDS Section 12: Consolidated Decisions
- Security bootstrap endpoint (/api/security/bootstrap)