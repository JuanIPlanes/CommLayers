# DEC-002: Lean Runtime Defaults

## Status
Accepted

## Date
2026-03-22

## Context
The platform should start with minimal infrastructure and only add components when justified by real use cases.

## Decision
Adopt lean runtime defaults:
1. Core infrastructure: PostgreSQL + Redis only (mandatory)
2. Optional infrastructure: Elasticsearch, Neo4j, Qdrant, MongoDB (enabled only when justified)
3. All optional components are behind feature flags / compose profiles
4. No component is introduced without a documented v1 use case

## Consequences
- Reduces operational complexity for initial deployment
- Forces explicit justification for each infrastructure component
- Aligns with FSD+SDS Section 8: Data and Projection Behavior

## Related
- FSD+SDS Section 8.1: Canonical Truth (PostgreSQL)
- FSD+SDS Section 8.2: Derived Projections (optional)
- FSD+SDS Section 8.3: Optional Document Store (MongoDB not unless justified)