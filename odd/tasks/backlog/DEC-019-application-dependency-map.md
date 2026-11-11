---
id: DEC-019
title: Application Dependency Map
status: proposed
date: 2026-11-11
---

## Context
The `main.go` defines 14 application types via `graphNodes` and 12 communication families. This decision formalizes the mapping of application types to communication families as a binding dependency map.

## Decision
### Graph Nodes (Application Types) from `main.go`
| Node ID | Label | Category | Description | Depends On |
|---------|-------|----------|-------------|------------|
| `request_response` | Request response | transport | Reference case for bounded request-reply work | — |
| `server_sent_events` | SSE | transport | Default push channel for one-way updates | `request_response` |
| `websocket` | WebSocket | transport | Planned duplex transport when truly needed | `request_response` |
| `workflow_engine` | Workflow engine | messaging | Visible queue and workflow state transitions | `request_response`, `server_sent_events` |
| `sync_replication` | Sync replication | sync | Primary/replica drift and explicit merge behavior | `workflow_engine` |
| `search_projection` | Search projection | backend-family | Full-text search over canonical data | `workflow_engine` |
| `graph_projection` | Graph projection | backend-family | Graph traversal over canonical data | `workflow_engine` |
| `vector_projection` | Vector projection | backend-family | Vector similarity over canonical data | `workflow_engine` |
| `messaging_broker` | Message broker | backend-family | Durable event backbone for fan-out | `workflow_engine` |
| `async_lro` | Async LRO | backend-family | Long-running operations with polling/subscription | `request_response` |
| `transport_polling` | Transport polling | transport | Polling/long-polling fallback | `request_response` |
| `localization_frontend` | Localization frontend | frontend | EN/ES resource bundles, locale detection | `request_response` |
| `pricing_enrichment` | Pricing enrichment | enrichment | USD freshness, confidence, fallback | `localization_frontend`, `comparative_visualization` |
| `comparative_visualization` | Comparative visualization | frontend | Desktop comparison UI, FPS verification | `workflow_engine`, `sync_replication`, `search_projection` |

### Communication Families (from DEC-014)
Transport: `request_response`, `polling`, `long_polling`, `server_sent_events`, `websocket`, `webhook`
Backend-Family: `messaging_workflow`, `sync_replication_projection`, `search_projection`, `graph_projection`, `vector_projection`, `mongo_projection`
Frontend: `comparative_visualization`
Enrichment: `pricing_cost`

### Application → Family Dependency Map
| Application Type | Primary Family | Secondary Families | Required For |
|------------------|----------------|-------------------|--------------|
| `request_response` | `request_response` | — | All other transport families |
| `server_sent_events` | `server_sent_events` | `request_response` | `workflow_engine`, `sync_replication`, `comparative_visualization` |
| `websocket` | `websocket` | `request_response` | `sync_replication` (duplex sync), `comparative_visualization` |
| `workflow_engine` | `messaging_workflow` | `server_sent_events` | `sync_replication`, `messaging_broker`, `async_lro`, `comparative_visualization` |
| `sync_replication` | `sync_replication_projection` | `workflow_engine`, `websocket` | `comparative_visualization`, `local_first_crdt` paradigm |
| `search_projection` | `search_projection` | `workflow_engine` | `comparative_visualization`, `planet_scale` paradigm |
| `graph_projection` | `graph_projection` | `workflow_engine` | `comparative_visualization`, `planet_scale` paradigm |
| `vector_projection` | `vector_projection` | `workflow_engine` | `comparative_visualization`, `planet_scale` paradigm |
| `messaging_broker` | `messaging_workflow` | `workflow_engine` | `event_driven` paradigm, `reactive_stream` paradigm |
| `async_lro` | `messaging_workflow` | `request_response`, `server_sent_events` | `reactive_stream` paradigm |
| `transport_polling` | `polling`, `long_polling` | `request_response` | Compatibility fallback |
| `localization_frontend` | — (cross-cutting) | `request_response` | `pricing_enrichment`, `comparative_visualization` |
| `pricing_enrichment` | `pricing_cost` | `localization_frontend`, `comparative_visualization` | `planet_scale` paradigm |
| `comparative_visualization` | `comparative_visualization` | All backend families | `planet_scale` paradigm |

### Dependency Graph Rules
1. **Acyclic**: Dependency graph must remain acyclic (validated in CI)
2. **Transitive Closure**: If A→B and B→C, then A transitively depends on C
3. **Family Activation**: A family is "active" when its primary application type is implemented
4. **Wave Boundaries**: 
   - First wave: `request_response`, `server_sent_events` (bootstrap)
   - Wave 2: `websocket`, `workflow_engine`, `sync_replication`, `search_projection`, `graph_projection`, `vector_projection`, `messaging_broker`, `async_lro`, `transport_polling`
   - Wave 3: `localization_frontend`, `comparative_visualization`
   - Wave 4: `pricing_enrichment`

### Executable Dependency Check
```go
func ValidateDependencyMap() error {
    // 1. Check acyclic
    if !dependencyGraph.IsAcyclic() {
        return fmt.Errorf("dependency graph has cycles")
    }
    
    // 2. Check all declared families exist in glossary (DEC-014)
    for _, app := range applicationTypes {
        for _, fam := range app.Families {
            if !glossary.HasFamily(fam) {
                return fmt.Errorf("application %s references unknown family %s", app.ID, fam)
            }
        }
    }
    
    // 3. Check wave boundaries respected
    for _, app := range applicationTypes {
        if app.Wave == 2 && !wave1Complete() {
            return fmt.Errorf("wave 2 app %s cannot activate before wave 1 complete", app.ID)
        }
        if app.Wave == 3 && !wave2Stabilized() {
            return fmt.Errorf("wave 3 app %s cannot activate before wave 2 stabilized", app.ID)
        }
        if app.Wave == 4 && !wave3Stabilized() {
            return fmt.Errorf("wave 4 app %s cannot activate before wave 3 stabilized", app.ID)
        }
    }
    
    return nil
}
```

### API Contract
`/api/application-map` returns:
```json
{
  "applications": [
    {
      "id": "workflow_engine",
      "label": "Workflow engine",
      "category": "messaging",
      "primary_family": "messaging_workflow",
      "secondary_families": ["server_sent_events"],
      "depends_on": ["request_response", "server_sent_events"],
      "required_for": ["sync_replication", "messaging_broker", "async_lro", "comparative_visualization"],
      "wave": 2,
      "status": "deferred"
    }
  ],
  "families": { ... },
  "validation": { "acyclic": true, "wave_boundaries_respected": true }
}
```

## Consequences
- Application dependency map is the single source of truth for wave sequencing
- New application types require dependency map update and acyclicity validation
- Family activation gates wave progression
- Map version embedded in evidence contracts

## Acceptance Criteria
- [ ] All 14 application types mapped to communication families
- [ ] Dependency graph validated acyclic in CI
- [ ] Wave boundaries enforced: wave N+1 cannot activate before wave N stabilized
- [ ] `/api/application-map` returns complete map with validation status
- [ ] Glossary cross-reference: all referenced families exist in DEC-014