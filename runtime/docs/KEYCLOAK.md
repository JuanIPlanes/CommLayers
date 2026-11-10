# Keycloak integration surface for CommLayers backend

Overview:
- Added a minimal Keycloak integration surface so runtime includes an actual auth provider instead of documentation-only references.
- A lightweight Keycloak middleware validates the presence of a Bearer token. Full token validation against Keycloak is not implemented in this iteration to keep scope small.

What changed:
- docker-compose.yaml now includes a Keycloak service wired to a Postgres backend.
- New Go package at services/api/internal/auth/keycloak.go exposing InitKeycloak(), IsEnabled(), and KeycloakMiddleware.
- The runtime router/bootstrap file provides a minimal wiring point for HTTP endpoints while preserving current behavior.

How to use:
- Call InitKeycloak() at startup if you want to enable runtime auth surface.
- Wrap protected routes with keycloak.KeycloakMiddleware in your HTTP server setup.

Notes:
- This is a minimal surface consistent with the requirement to avoid a full auth enforcement everywhere. It is enough to show an actual runtime component and to enable/disable in code.
