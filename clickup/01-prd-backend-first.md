# CommLayers PRD - Backend-First Revision

## 1. Product Intent

CommLayers is a comparative engineering platform that implements, benchmarks, and visualizes backend communication layers and adjacent architectural patterns before presenting them in a desktop visualization experience.

The product is not only a graph UI. It is a working backend laboratory that demonstrates when each communication method is useful, how it behaves under load, what operational cost and complexity it introduces, and how a frontend may consume it safely.

## 2. Product Goal

Deliver a backend-first platform that:
- implements all in-scope communication families as working backend capabilities
- benchmarks and compares their behavior with observable state transitions
- exposes normalized contracts for a later frontend comparison experience
- stages architectural paradigms for a future v2 expansion

## 3. Scope Boundaries

### 3.1 v1 Scope - Backend-First Comparative Platform + 2D Visualization

v1 includes:
- normalized backend abstraction layers
- working implementations of all in-scope communication families
- benchmark capture and comparison logic embedded in the platform definition
- i18n foundations with English default and Spanish JSON localization
- desktop-only 2D frontend to visualize backend comparisons and delayed state transitions
- security, identity, secrets, audit, and observability foundations

v1 excludes:
- full architectural paradigms implementation as a first-class runtime mode
- production-wide service-mesh hardening on day one
- cloud pricing integrations before backend core, i18n, and comparative engines are stable

### 3.2 v2 Scope - Architectural Paradigms

v2 expands into architectural paradigms, including but not limited to:
- reactive and stream-driven systems
- local-first and CRDT-oriented systems
- event-driven architecture
- consensus-driven systems
- planet-scale and multi-region patterns

## 4. Backend Abstraction Layers

The canonical backend layer model is:
1. Transport and API Layer
2. Application and Use-Case Layer
3. Domain Layer
4. Adapters and Integration Layer
5. Runtime, Operations, and Security Layer
6. Projections and Query Engines Layer

Each implementation in CommLayers must identify which layer it belongs to, which layer it depends on, and which test gate proves it works.

## 5. Communication Families in v1

All of the following must be implemented as working comparative backend capabilities in v1:
- request-response
- polling
- long polling
- server-sent events
- WebSocket
- WebRTC signaling support
- webhooks
- message broker and event bus patterns
- job queue and workflow engine patterns
- event store pattern
- stream processor pattern
- actor-style runtime pattern
- realtime sync engine pattern
- edge sync and replication pattern
- local-first sync pattern
- CRDT-backed sync pattern
- distributed database and replicated state pattern
- consensus and coordination pattern
- service mesh and control-plane demonstration path
- orchestration and autonomous platform demonstration path
- planet-scale composite architecture path as future-facing but taskable groundwork

## 6. Common Application Comparisons Required

For each backend family, the product must show:
- when it is useful
- when it is overkill
- expected latency profile
- reliability profile
- operational complexity
- local developer complexity
- observability impact
- cost amplifier profile
- recommended stage of product adoption

The frontend later visualizes these comparisons; the backend must generate the evidence first.

## 7. Benchmark and Timing Policy

Benchmarks are part of v1 definition now, not later.

Required benchmark outputs:
- baseline request latency
- streaming startup time
- per-update latency
- p50, p95, and p99 timings where meaningful
- state transition timing visible to users
- deliberate visualization delay and timeout support so users can actually observe transitions

Governed defaults:
- async threshold: expected work over 10 seconds must move to long-running operation flow
- frontend performance target: minimum 30 FPS, ideal 60 FPS
- realtime default: SSE first, WebSocket only when bidirectional and latency-critical behavior is justified

## 8. Internationalization and Content Policy

Before cloud pricing work begins, the platform must support:
- English as default language
- Spanish via JSON localization files
- localization of every response returned to the frontend that is user-visible
- localized error, status, and comparison text returned by backend APIs where surfaced in the UI

## 9. Cloud Pricing Ordering

Cloud pricing is a late-stage backend enrichment capability.

It must not begin before:
- canonical backend implementations exist
- benchmark capture exists
- i18n foundations exist
- normalized response contracts exist

When implemented, pricing must be:
- USD only
- confidence-aware
- staleness-aware
- user-selectable for refresh, with 5s default

## 10. Identity, Security, and Compliance Baseline

Fixed choices:
- Identity Provider: Keycloak
- Secrets Broker: Vault

Required baseline:
- Authorization Code + PKCE only
- BFF-oriented browser architecture whenever possible so tokens do not live in frontend storage
- access token lifespan target: 5-10 minutes
- rotated refresh tokens with one-time replacement behavior
- SSO idle target: 30 minutes
- SSO max target: 12-24 hours
- staged mTLS rollout
- TLS 1.3 at ingress
- audit and traceability designed from day one
- implementation mapping suitable for ISO 27001 and NIST 800-53 alignment

## 11. Runtime and Environment Constraints

The entire project must be runnable with only two compose entrypoints:
- `docker-compose.back.yml`
- `docker-compose.front.yml`

Preferred backend runtime palette:
- PostgreSQL as canonical SQL system of record
- `pgvector` inside PostgreSQL as the default vector capability in v1
- Redis for cache, queue assist, and ephemeral coordination
- Elasticsearch when search is required on day one
- graph database only if graph traversal is a first-class v1 feature
- dedicated vector database only when vector workloads justify their own operational domain
- Mongo only if a concrete v1 use case proves it necessary

Recommended runtime profile for first runnable backend:
- mandatory: PostgreSQL, Redis
- recommended if search matters in v1: Elasticsearch
- optional profile: Neo4j for graph traversal
- optional profile: Qdrant if vectors outgrow `pgvector`
- optional profile: Mongo only with explicit justification

## 12. Language Ownership

Approved default ownership:
- Go: low-latency APIs, realtime and transport-sensitive paths
- Java: workflow-heavy and enterprise integration services
- Python: benchmarking, simulation, offline analysis, or ETL where justified
- TypeScript: frontend and optional BFF/admin tooling

## 13. Acceptance Gates Before Implementation Starts

Before `sdd-apply`, the following must exist:
- revised backend-first PRD, TDD, and FSD+SDS
- stream-level ClickUp tasks staged in bulk form
- backend abstraction layer ownership clearly defined
- benchmark plan embedded in docs
- frontend role reduced to consumer/visualizer of backend evidence
- v2 paradigm work explicitly separated from v1 backend-first execution
