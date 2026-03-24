# CommLayers

CommLayers is a backend-first comparative platform for communication layers and architectural evolution.

This repository currently contains the first honest runnable bootstrap version:
- a Go API exposing the first-wave foundation and bootstrap comparison endpoints
- a TypeScript frontend that consumes that API and visualizes the current bootstrap state
- split Docker Compose entrypoints for backend and frontend
- staged docs and ClickUp/SDD planning artifacts for later waves

## What This Version Is

This is a foundation snapshot, not the full planned platform.

Implemented in this first working version:
- backend-first project skeleton
- Go API with:
- `/healthz`
- `/api/bootstrap`
- `/api/streams`
- `/api/catalog`
- `/api/first-wave/contract`
- `/api/first-wave/streams/:id`
- `/api/security/bootstrap`
- `/api/data-platform`
- `/api/benchmark-framework`
- `/api/deferred-waves`
- `/api/v2-readiness`
- `/api/comparisons/realtime`
- `/api/async/demo`
- `/api/events` (SSE demo)
- TypeScript frontend with:
- first-wave execution console
- stream contract cards
- security, data, and benchmark panels
- deferred-wave and v2 readiness panels
- catalog snapshot and realtime guidance
- SSE event feed
- async demo runner
- Compose-based runtime with PostgreSQL and Redis active by default
- real Postgres-backed persistence for workflow runs, transport sessions, sync sessions, and async job state
- Redis-backed ephemeral cache layer for persisted demo entities
- Redis-backed transient event flow for SSE updates across workflow, sync, transport, and async demo changes
- Optional backend profiles for Elasticsearch, Neo4j, Qdrant, and MongoDB

Not implemented yet:
- Keycloak and Vault runtime integration
- full transport-family implementation set
- messaging, sync, projection, and pricing execution
- frontend comparative graph UI
- v2 paradigm runtime

## Repository Layout

```text
clickup/               staged docs, task payloads, and deferred-wave manifests
frontend/              static TypeScript frontend
services/api/          Go API bootstrap service
docker-compose.back.yml
docker-compose.front.yml
```

## Running Locally

Start backend first so the shared Docker network exists:

```bash
docker compose -f docker-compose.back.yml up --build -d
```

Then start the frontend:

```bash
docker compose -f docker-compose.front.yml up --build -d
```

Open:
- frontend: `http://localhost:4173`
- health via frontend gateway: `http://localhost:4173/healthz`
- bootstrap via frontend gateway: `http://localhost:4173/api/bootstrap`

Optional profiles on the backend stack:

```bash
docker compose -f docker-compose.back.yml --profile search up -d
docker compose -f docker-compose.back.yml --profile graph up -d
docker compose -f docker-compose.back.yml --profile vector up -d
docker compose -f docker-compose.back.yml --profile mongo up -d
```

## First-Wave Execution Order

Active first-wave streams:
1. `86agbv0k4` - backend abstraction layer matrix
2. `86agbv0k3` - Keycloak and Vault security bootstrap
3. `86agbv0k7` - canonical data platform
4. `86agbv0k9` - benchmark and observability framework

Deferred waves remain staged in `clickup/08-deferred-waves.md`.

## Validation Snapshot

The current bootstrap version was validated locally with:

```bash
cd services/api && go build ./...
cd ../../frontend && npm run build
cd .. && docker compose -f docker-compose.back.yml config
docker compose -f docker-compose.front.yml config
```

Validation result at the current stage:
- Go API build passed
- frontend production build passed
- backend compose config passed
- frontend compose config passed
- full split stack startup passed through the frontend gateway on `4173`
- persisted workflow and sync records survived an API container restart
- SSE now streams Redis-backed transient events for workflow, sync, transport, and async job state changes

## Notes

- Frontend and pricing remain intentionally deferred beyond this bootstrap.
- The repo is structured so a public GitHub push can happen once the bootstrap files are committed.
- The ClickUp and SDD artifacts remain the planning source for later implementation waves.
