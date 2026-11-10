---
# Transport Family Feature

## Specs
1. Implement transport-family scaffolding with provider interface (already created). 
2. Add Neo4j (graph) concrete provider.
3. Add Qdrant (vector) concrete provider.
4. Wire providers into API router under `/api/transport/*`.
5. Provide unit tests for both providers.

## Tasks
- T1: Create `services/api/internal/transport/graph/neo4j_transport.go` (implementation stub).
- T2: Create `services/api/internal/transport/vector/qdrant_transport.go` (implementation stub).
- T3: Create `services/api/internal/router/transport_routes.go` (router stub).
- T4: Update `services/api/cmd/server/main.go` to initialize providers based on compose profiles and register routes.
- T5: Add unit tests for Neo4j and Qdrant providers.

## Log
- L1: User requested to finish all 2.* tasks with ODD gentle‑ai.
