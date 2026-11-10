# DEC-003: BFF-Oriented Authentication

## Status
Accepted

## Date
2026-03-22

## Context
Frontend authentication should not hold persistent OAuth tokens. The backend-for-frontend (BFF) pattern should manage session state.

## Decision
Implement BFF-oriented authentication:
1. Backend manages session state (BFFSession with ID, user info, roles, tokens)
2. Frontend receives only session reference (cookie or header)
3. Access/refresh tokens stored only in backend, never exposed to frontend
4. Session lifecycle: create, retrieve, refresh, delete
5. Integration point for Keycloak (enabled via InitKeycloak)
6. Audit trail for all authenticated requests

## Consequences
- Reduces frontend attack surface (no token storage)
- Centralizes token management and rotation
- Enables server-side session revocation
- Aligns with FSD+SDS Section 2.3: Frontend Consumption

## Related
- FSD+SDS Section 2.3: Frontend Consumption
- Keycloak integration in auth/keycloak.go