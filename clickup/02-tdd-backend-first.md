# CommLayers TDD - Backend-First Revision

## 1. Architectural Reframe

CommLayers v1 is a comparative backend platform with a governed visualization layer. The architecture must be designed so that multiple communication styles can be implemented, benchmarked, observed, and compared under one canonical contract surface.

The safest implementation shape is hybrid:
- one canonical backend platform with shared domain contracts and evidence model
- multiple bounded backend modules demonstrating different communication families
- one frontend visualization layer consuming normalized outputs

## 2. Canonical Layer Model

### 2.1 Transport and API Layer

Responsibilities:
- HTTP request-response endpoints
- polling and long-polling endpoints
- SSE endpoints
- WebSocket gateways
- webhook ingress and egress endpoints
- long-running operation endpoints and status polling

Rules:
- no business decisions outside transport orchestration
- explicit protocol contracts and timeout behavior
- TDD red/green gate per endpoint family

### 2.2 Application and Use-Case Layer

Responsibilities:
- orchestration of benchmark runs
- comparison flows
- async job submission and state management
- protocol selection recommendations
- localization-aware response assembly

Rules:
- protocol-agnostic use cases where possible
- idempotent submit/retry handling for async flows

### 2.3 Domain Layer

Responsibilities:
- normalized communication method catalog
- abstraction-layer taxonomy
- comparison metrics model
- state transition model
- recommendation rules
- confidence and staleness semantics

Rules:
- no storage or transport details
- canonical normalized naming only

### 2.4 Adapters and Integration Layer

Responsibilities:
- PostgreSQL persistence
- Redis integration
- Elasticsearch projection
- graph DB projection
- vector DB projection
- Keycloak integration
- Vault integration
- optional Mongo adapter only if justified

Rules:
- every adapter must prove why it exists
- no hidden datastore coupling in domain logic

### 2.5 Runtime, Operations, and Security Layer

Responsibilities:
- compose topology
- service boot order
- secrets retrieval
- staged mTLS
- metrics, logs, traces, audit
- test harness execution

Rules:
- `docker-compose.back.yml` is the only backend local entrypoint
- `docker-compose.front.yml` is the only frontend local entrypoint

### 2.6 Projections and Query Engines Layer

Responsibilities:
- search-oriented projection
- graph-oriented projection
- vector-oriented projection
- comparison query surfaces for frontend consumption

Rules:
- projections are derived, not canonical truth
- canonical truth starts in PostgreSQL-backed domain model

## 3. Language Ownership and Service Boundaries

### 3.1 Go Services

Use Go for:
- latency-sensitive transport APIs
- SSE
- WebSocket
- long-polling
- webhook receivers
- comparison read APIs

### 3.2 Java Services

Use Java for:
- workflow engine integrations
- queue orchestration where enterprise workflow semantics matter
- complex provider integrations
- later pricing and cost normalization services

### 3.3 Python Services and Tools

Use Python only where justified for:
- benchmarking harnesses
- simulation and replay tools
- ETL/projection helpers
- offline comparison analysis

### 3.4 TypeScript Services

Use TypeScript for:
- frontend application
- optional BFF/admin orchestration if a frontend-facing mediation layer is needed

## 4. Datastore and Runtime Defaults

### 4.1 Mandatory v1 Datastores

- PostgreSQL - canonical system of record
- Redis - cache, queue assist, ephemeral coordination

### 4.2 Recommended v1 Datastores

- Elasticsearch - search projection when search matters in v1
- `pgvector` inside PostgreSQL - default vector capability for v1

### 4.3 Conditionally Allowed v1 Datastores

- Graph DB - only if graph traversal is a first-class product feature in v1
- Dedicated Vector DB - only if vector workloads justify their own service boundary
- MongoDB - optional, only if a concrete v1 document-native use case is approved

### 4.4 Default Recommendations

- Graph DB default when required: Neo4j as the safer long-term default
- Vector DB default in v1: `pgvector` in PostgreSQL
- First dedicated vector DB when required: Qdrant

### 4.5 Compose Profiles

`docker-compose.back.yml` should use profiles so optional services stay off by default:
- base: PostgreSQL, Redis
- search: Elasticsearch
- graph: Neo4j
- vector: Qdrant
- mongo: MongoDB

## 5. Security and Identity Design

### 5.1 Identity

- Keycloak is mandatory
- Browser login uses Authorization Code + PKCE with S256 only
- Prefer BFF-oriented browser architecture so browser tokens stay off persistent frontend storage
- Access token lifespan target: 5-10 minutes
- Refresh tokens rotate on every use and the latest token must replace the previous one
- SSO session idle target: 30 minutes
- SSO session max target: 12-24 hours
- RP session control is enforced independently of token lifetime
- Exact redirect URI matching is mandatory

### 5.2 Secrets

- Vault is mandatory
- backend services must read sensitive credentials from Vault, not static env files
- first milestone includes dynamic DB credentials or equivalent short-lived credential flow where practical
- Vault PKI should follow offline root -> online intermediate hierarchy
- Service certificates should be short-lived and renewed automatically

### 5.3 mTLS Rollout

Staged plan:
1. ingress TLS and service identity foundations
2. selective mTLS on sensitive service pairs
3. broader east-west mTLS rollout when service count and risk justify it
4. later stricter service-mesh enforcement in a hardening milestone

Full east-west mTLS is required later, but it is not a day-one blocker for v1 runnable backend milestones.

## 6. Protocol Selection Defaults

### 6.1 Realtime

- default to SSE for one-way server push
- escalate to WebSocket for frequent client-to-server interaction, binary framing needs, or stricter interactive latency needs
- infrastructure note: reverse proxies must disable buffering for SSE delivery
- security note: raw WebSocket traffic behind common edge platforms may lose deep request inspection after the `101` upgrade and should be reserved for cases that truly need duplex interaction

### 6.2 Async Work

- if expected completion exceeds 10s, use long-running operation flow
- return `202 Accepted`
- expose operation status resource
- define `queued`, `running`, `succeeded`, `failed`, and `cancelled` states

## 7. Confidence, Staleness, and Comparison Semantics

Every frontend-facing backend response that carries comparison evidence must support:
- `confidence`
- `staleness`
- `source`
- `measured_at`
- optional `delay_applied_ms` when visualization timing is intentionally stretched for observability

## 8. Benchmark Architecture

Benchmarks are first-class architecture artifacts.

Must measure:
- request latency
- stream startup latency
- update propagation delay
- per-protocol resource overhead
- queue delay
- search, graph, and vector projection latency
- user-visible state transition delays

Initial benchmark defaults to encode now:
- request-response is baseline for simple request/reply work
- polling is fallback for environments that cannot hold push connections
- long-polling is compatibility fallback, not preferred steady-state transport
- SSE is default server-to-client push channel
- WebSocket is preferred only when duplex interaction or tighter interactivity requirements justify it

## 9. Delivery Order

1. canonical domain model and abstraction taxonomy
2. security bootstrap: Keycloak and Vault
3. PostgreSQL-backed core and normalized APIs
4. async workflow and LRO model
5. request-response, polling, long-polling, SSE, WebSocket, webhook families
6. broker, job, stream, event, actor, sync, replication, consensus families
7. projections: search, graph, vector
8. i18n and frontend-facing localization completeness
9. frontend comparative visualization and delay-aware state transitions
10. pricing and cost enrichment
11. later hardening and paradigm-v2 work

## 10. TDD Red/Green Gates

Every stream must define:
- failing test or benchmark expectation first
- minimal implementation to pass
- refactor pass without breaking comparative evidence

No backend family is considered ready without a red/green artifact trail.
