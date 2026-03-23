# CommLayers Research-Backed Defaults

This file captures the external findings that informed the backend-first document revisions.

## 1. Realtime Transport Defaults

Recommended default order:
1. request-response for basic request/reply work
2. polling or long-polling only as compatibility fallback
3. SSE as default server-to-client push channel
4. WebSocket only when duplex interaction, binary framing, or tighter interactivity requirements justify it

Important implementation notes:
- reverse proxies must disable buffering for SSE so events flush immediately
- some edge platforms inspect only the initial WebSocket upgrade request, not all traffic after the protocol upgrade
- long-polling is the lowest preference steady-state transport where SSE or WebSocket are possible

Practical source set used in delegated research:
- WHATWG/MDN/EventSource and WebSocket browser behavior
- NGINX reverse proxy guidance for response buffering
- Cloudflare WebSocket behavior guidance
- Microsoft SignalR and Azure SignalR transport ordering and performance notes

## 2. Auth and Session Defaults

Recommended default:
- Authorization Code + PKCE with S256 only
- prefer BFF-oriented browser architecture
- keep OAuth tokens out of persistent frontend storage
- access token lifespan target: 5-10 minutes
- rotate refresh tokens and replace them on every refresh
- SSO idle target: 30 minutes
- SSO max target: 12-24 hours

Hard constraints carried into TDD/FSD:
- no implicit flow
- no resource owner password credentials flow
- exact redirect URI matching only
- audience and scope restriction required

Practical source set used in delegated research:
- RFC 9700
- RFC 7636
- OAuth browser-based applications BCP draft
- Keycloak token/session guidance

## 3. mTLS and PKI Defaults

Recommended staged rollout:
1. ingress TLS plus service identity foundations
2. Vault PKI hierarchy with offline root and online intermediate
3. selective permissive mTLS on critical service pairs
4. strict east-west mTLS after inventory, observability, and authz are ready

Hard constraints:
- workload identity naming must be explicit before strict mesh-wide mTLS
- certificate expiry, renewal failures, handshake failures, and authz denials must be observable before strict enforcement
- service certificates should be short-lived and automatically renewed

Practical source set used in delegated research:
- NIST SP 800-207
- NIST SP 800-204B
- Vault PKI guidance
- SPIFFE identity guidance
- Istio authn and PeerAuthentication guidance

## 4. Datastore Defaults

Lean polyglot default for v1:
- mandatory: PostgreSQL, Redis
- recommended if search matters day one: Elasticsearch
- vector default: `pgvector` in PostgreSQL first
- graph default only if graph traversal is first-class: Neo4j
- dedicated vector DB only when vectors justify a separate operational boundary: Qdrant
- Mongo optional only with concrete document-native justification

Operational recommendation:
- keep stateful backend infrastructure in `docker-compose.back.yml`
- keep UI/web tooling in `docker-compose.front.yml`
- use optional compose profiles for search, graph, vector, and Mongo

Practical source set used in delegated research:
- `pgvector` docs
- Redis docs
- Elasticsearch Docker docs
- Neo4j Docker Compose docs
- Qdrant quickstart and operations docs
- Docker Compose networking and volumes docs

## 5. CommLayers-Specific Impact

These defaults support the user-approved direction:
- backend-first delivery
- all communication families visible and comparable in v1
- frontend as consumer/visualizer, not source of truth
- benchmark thresholds encoded now
- pricing delayed until backend core, observability, and i18n are in place
