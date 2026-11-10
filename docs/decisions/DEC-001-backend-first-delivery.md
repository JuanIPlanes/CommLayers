# DEC-001: Backend-First Delivery

## Status
Accepted

## Date
2026-03-22

## Context
The platform must demonstrate backend communication families before building frontend visualizations. The backend must provide normalized responses that the frontend consumes without becoming the source of truth.

## Decision
Adopt a backend-first delivery approach where:
1. All backend families are implemented and tested before frontend work begins
2. API responses follow a normalized schema (familyResponse) that includes latency metrics, state transitions, evidence, and i18n support
3. Frontend development starts only after backend APIs are stable and documented

## Consequences
- Backend APIs are the contract; frontend adapts to them
- Enables parallel frontend development once API spec is frozen
- Reduces risk of frontend-driven API changes

## Related
- FSD+SDS Section 1: Functional Intent
- FSD+SDS Section 3: Required Response Model