---
id: DEC-017
title: v2 Unlock Prerequisites
status: proposed
date: 2026-11-11
---

## Context
The `/api/v2-readiness` endpoint in `main.go` defines 5 prerequisites for v2 paradigm implementation. This decision formalizes these as executable checks with clear pass/fail criteria.

## Decision
### The 5 Prerequisites (from `handleV2Readiness`)

#### 1. First-Wave Abstraction Matrix Executed in Code
**Check**: All four first-wave streams (86agbv0k4, 86agbv0k3, 86agbv0k7, 86agbv0k9) have passed their acceptance gates AND their deliverables are implemented in code (not just documented).
```go
func CheckAbstractionMatrixExecuted() error {
    gates := []string{"86agbv0k4", "86agbv0k3", "86agbv0k7", "86agbv0k9"}
    for _, g := range gates {
        if !gateRegistry.Passed(g) {
            return fmt.Errorf("gate %s not passed", g)
        }
    }
    // Verify implementation exists (not just docs)
    if !implementationRegistry.Has("backend_abstraction_matrix") {
        return fmt.Errorf("backend abstraction matrix not implemented")
    }
    if !implementationRegistry.Has("security_bootstrap") {
        return fmt.Errorf("security bootstrap not implemented")
    }
    if !implementationRegistry.Has("canonical_data_platform") {
        return fmt.Errorf("canonical data platform not implemented")
    }
    if !implementationRegistry.Has("benchmark_framework") {
        return fmt.Errorf("benchmark framework not implemented")
    }
    return nil
}
```

#### 2. Security Bootstrap Integrated at Runtime
**Check**: Keycloak and Vault are initialized and operational at runtime (not just config).
```go
func CheckSecurityBootstrapRuntime() error {
    if !keycloak.IsEnabled() {
        return fmt.Errorf("keycloak not enabled at runtime")
    }
    pki := security.GetPKIObservability()
    if pki == nil || !pki.IsEnabled() {
        return fmt.Errorf("PKI observability not enabled at runtime")
    }
    // Verify trust boundaries evaluatable
    boundaries := pki.GetTrustBoundaries()
    if len(boundaries) < 4 {
        return fmt.Errorf("insufficient trust boundaries configured")
    }
    return nil
}
```

#### 3. Canonical Data Platform Backed by Persistent Models
**Check**: PostgreSQL store is operational and all canonical entity kinds are persisted.
```go
func CheckCanonicalDataPlatformPersistent() error {
    store := persistence.GetStore()
    if store == nil {
        return fmt.Errorf("postgres store not initialized")
    }
    // Verify all entity kinds exist and are queryable
    kinds := []string{"stream_definitions", "comparison_evidence", "job_lifecycle", "deferred_wave_manifest"}
    for _, k := range kinds {
        if !store.KindExists(k) {
            return fmt.Errorf("canonical kind %s not persisted", k)
        }
    }
    // Verify audit lineage fields on all records
    if !store.AuditLineageComplete() {
        return fmt.Errorf("audit lineage incomplete")
    }
    return nil
}
```

#### 4. Benchmark Closure Accepted
**Check**: Benchmark matrix executed, thresholds met, closure report accepted.
```go
func CheckBenchmarkClosureAccepted() error {
    if !benchmarkClosureReport.Accepted() {
        return fmt.Errorf("benchmark closure report not accepted")
    }
    // Verify all three scenarios executed and within thresholds
    scenarios := []string{"request_baseline", "sse_startup_cadence", "async_visibility"}
    for _, s := range scenarios {
        result := benchmarkResults.Latest(s)
        if result == nil || !result.WithinThresholds() {
            return fmt.Errorf("benchmark %s not within thresholds", s)
        }
    }
    // Verify telemetry schema stable
    if !telemetrySchema.Stable() {
        return fmt.Errorf("telemetry schema not stable")
    }
    return nil
}
```

#### 5. Transport, Messaging, Sync, and Projection Families Implemented
**Check**: All deferred wave 2 families have at least bootstrap implementations.
```go
func CheckBackendFamiliesImplemented() error {
    families := []string{
        "transport",        // 86agbv0k5
        "messaging_workflow", // 86agbv0ka
        "sync_replication",   // 86agbv0kb
        "search_projection",  // 86agbv0kc (search)
        "graph_projection",   // 86agbv0kc (graph)
        "vector_projection",  // 86agbv0kc (vector)
    }
    for _, f := range families {
        if !familyRegistry.Implemented(f) {
            return fmt.Errorf("family %s not implemented", f)
        }
    }
    return nil
}
```

### Unlock Path (Executable Sequence)
```go
func V2UnlockPath() []string {
    var steps []string
    if err := CheckAbstractionMatrixExecuted(); err != nil {
        steps = append(steps, "Complete and accept first-wave implementation milestones")
    }
    if err := CheckSecurityBootstrapRuntime(); err != nil {
        steps = append(steps, "Execute deferred wave 2 backend families")
    }
    if err := CheckCanonicalDataPlatformPersistent(); err != nil {
        steps = append(steps, "Stabilize localization and frontend comparative visualization")
    }
    if err := CheckBenchmarkClosureAccepted(); err != nil {
        steps = append(steps, "Treat pricing as late-stage only")
    }
    if err := CheckBackendFamiliesImplemented(); err != nil {
        steps = append(steps, "Open v2 paradigm implementation planning")
    }
    return steps
}
```

### V2 Paradigm Status (from `handleV2Roadmap`)
Each paradigm has explicit unlock dependencies:
| Paradigm | Status | Unlock Dependencies |
|----------|--------|---------------------|
| `reactive_stream_driven` | planned | messaging, sync, and benchmark maturity |
| `local_first_crdt` | planned | sync slice maturity + conflict strategy stabilization |
| `event_driven_architecture` | planned | transport and messaging slices + durable event contract |
| `consensus_driven` | planned | coordination semantics beyond bootstrap-grade demos |
| `planet_scale_composite` | planned | projection, pricing, and frontend comparative maturity |

### Gate Enforcement
- `/api/v2-readiness` returns current status and blockedBy list
- CI/CD blocks v2 implementation tasks until all 5 prerequisites pass
- Each prerequisite check is independently testable and auditable
- Manual override requires architectural decision record (ADR)

## Consequences
- v2 work cannot begin until all 5 prerequisites are demonstrably met
- Prerequisites are executable checks, not documentation
- Each paradigm has its own unlock dependencies (not just the 5 global ones)
- V1 completion list is the baseline: v2 builds on implemented v1 slices

## Acceptance Criteria
- [ ] All 5 check functions implemented and return correct status
- [ ] `/api/v2-readiness` accurately reflects prerequisite status
- [ ] CI blocks v2 tasks when prerequisites unmet
- [ ] Implementation registry tracks code-level deliverables (not just docs)
- [ ] Paradigm unlock dependencies documented and enforced per paradigm