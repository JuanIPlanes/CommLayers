# CommLayers Stream-Level ClickUp Task Definitions

This file stages the stream-level tasks for later bulk creation in ClickUp.

## 1. Task Creation Rules

- Create tasks in bulk, not one by one.
- Before creating any task, check whether the same action already exists in staged payloads.
- If duplicate or near-identical actions exist, merge them into one bulk action.
- Use Playwright-driven ClickUp handling if MCP auth/session or rate limits are unsafe.

## 2. Stream Tasks

| Stream ID | Task Name | Goal | Core Deliverables | Depends On | Definition of Done |
| --- | --- | --- | --- | --- | --- |
| ST-01 | Backend abstraction layer matrix | Freeze canonical backend layer model and common application taxonomy | layer matrix, dependency rules, glossary, normalized names | none | docs approved and referenced by all later streams |
| ST-02 | Security bootstrap with Keycloak and Vault | Establish identity and secrets baseline | Keycloak realm/client model, BFF-oriented auth pattern, token/session defaults, Vault bootstrap, Vault PKI, staged mTLS plan | ST-01 | auth/secrets docs and runnable bootstrap plan complete |
| ST-03 | Core canonical data platform | Establish PostgreSQL-first domain core and shared evidence model | normalized domain schema, response contracts, confidence/staleness model | ST-01 | canonical core defined and test gates staged |
| ST-04 | Async and long-running operations | Implement async rule and state model | 202/LRO model, queue semantics, status lifecycle, timeout policy | ST-03, ST-02 | async contract and benchmark gates defined |
| ST-05 | Transport family implementations | Stage all primary transport families | request-response, polling, long-polling, SSE, WebSocket, webhooks, WebRTC signaling, proxy and edge constraints | ST-03, ST-04 | protocol families defined with comparison outputs, benchmark defaults, and test plan |
| ST-06 | Messaging and workflow families | Stage broker, queue, stream, event, actor patterns | event bus, job/workflow, event store, stream processor, actor runtime definitions | ST-03, ST-04 | each family has use criteria and test gate |
| ST-07 | Sync, replication, and coordination families | Stage advanced state-sync families | realtime sync, edge sync, local-first, CRDT, distributed state, consensus | ST-03, ST-04 | advanced families have comparative role and data rules |
| ST-08 | Derived projections and query engines | Define search, graph, vector projections | Elasticsearch projection, optional Neo4j graph profile, `pgvector` first, Qdrant escalation path, optional Mongo justification gate | ST-03 | all required projections defined, Mongo decision explicit |
| ST-09 | Benchmark and observability framework | Embed benchmark capture now | latency matrix, state delay policy, tracing/logging/audit metrics | ST-03, ST-05, ST-06, ST-07, ST-08 | benchmark plan approved and wired into stream acceptance |
| ST-10 | Localization and frontend-facing content | Ensure all frontend-visible responses are localizable | English default, Spanish JSON, localized backend response fields | ST-03 | i18n appears before pricing work and covers all frontend-facing returns |
| ST-11 | Frontend comparative visualization | Build desktop consumer for backend evidence | comparison UI, delay-aware transitions, protocol comparison surfaces | ST-05, ST-06, ST-07, ST-08, ST-09, ST-10 | frontend is consumer only, not source of truth |
| ST-12 | Pricing and cost enrichment | Add cloud pricing last | USD-only pricing, 5s default freshness selector, confidence/staleness display | ST-09, ST-10, ST-11 | pricing arrives after backend core and i18n are done |
| ST-13 | v2 paradigm groundwork | Prepare future paradigm implementation path | paradigm mapping, deferred task set, dependency on v1 evidence | ST-01, ST-09 | v2 isolated from v1 execution but taskable |

## 3. Task Ordering Notes

- `ST-01` through `ST-10` are backend-first core.
- `ST-11` is frontend consumption after backend evidence is available.
- `ST-12` is late-stage enrichment.
- `ST-13` prepares v2 but does not block v1.

## 4. Stream Tags

- `backend-first`
- `sdd-revision`
- `bulk-create`
- `benchmark-required`
- `i18n-required`
- `security-foundation`
- `frontend-consumer`
- `pricing-late`
- `v2-paradigms`
