---
id: DEC-013
title: Deferred Implementation Boundaries
status: proposed
date: 2026-11-11
---

## Context
Three deferred waves are defined in `main.go` with streams that must not begin until reactivation rules are met. This decision formalizes the boundaries for each wave and the reactivation rules as executable checks.

## Decision
### Wave 2 — Backend Family Expansion
**Streams**: 86agbv0k6 (Async/LRO), 86agbv0k5 (Transport), 86agbv0ka (Messaging/Workflow), 86agbv0kb (Sync/Replication), 86agbv0kc (Search/Graph/Vector)
**Purpose**: Unlock transport, messaging, sync, and projection implementation only after first-wave approval
**Focus Areas**: long-running operations, transport deepening, messaging families, sync and replication, projection profiles

**Reactivation Rules (Executable)**:
```go
func ReactivateWave2() error {
    // 1. Approve first-wave artifacts (all four gates passed)
    if !gate86agbv0k4.Passed() || !gate86agbv0k3.Passed() || !gate86agbv0k7.Passed() || !gate86agbv0k9.Passed() {
        return fmt.Errorf("first-wave gates not all passed")
    }
    
    // 2. Implementation sequencing behind 86agbv0k9 closure
    if !waveClosureReport.Accepted() {
        return fmt.Errorf("wave closure report not accepted")
    }
    
    // 3. Dependency contracts published and frozen
    if !dependencyContract.IsFrozen() {
        return fmt.Errorf("dependency contract not frozen")
    }
    
    // 4. Per-stream prerequisites
    // 86agbv0k6: LRO contract + idempotency rules + polling/subscription visibility model
    if !lroContract.IsDefined() || !idempotencyRules.Defined() || !visibilityModel.Defined() {
        return fmt.Errorf("86agbv0k6 prerequisites not met")
    }
    // 86agbv0k5: Request-response contract + polling comparison + SSE/WebSocket guidance + webhook notes
    if !transportContracts.Defined() {
        return fmt.Errorf("86agbv0k5 prerequisites not met")
    }
    // 86agbv0ka: Broker criteria + queue semantics + workflow contract + failure mode matrix
    if !messagingContracts.Defined() {
        return fmt.Errorf("86agbv0ka prerequisites not met")
    }
    // 86agbv0kb: Sync behaviors + conflict handling + replication checks + consensus criteria
    if !syncContracts.Defined() {
        return fmt.Errorf("86agbv0kb prerequisites not met")
    }
    // 86agbv0kc: Search projection model + graph profile gate + vector escalation rules
    if !projectionContracts.Defined() {
        return fmt.Errorf("86agbv0kc prerequisites not met")
    }
    
    return nil
}
```

### Wave 3 — Frontend and Enrichment
**Streams**: 86agbv0kf (Localization), 86agbv0kg (Frontend Comparative Visualization), 86agbv0ke (Pricing)
**Purpose**: Unlock localized frontend comparison and pricing after backend-family stability
**Focus Areas**: localization, frontend comparative visualization, pricing enrichment

**Reactivation Rules (Executable)**:
```go
func ReactivateWave3() error {
    // 1. Backend-family work must stabilize first
    if !wave2Stabilized() {
        return fmt.Errorf("wave 2 backend-family work not stabilized")
    }
    
    // 2. Pricing remains late-stage only (depends on 86agbv0kg)
    if !gate86agbv0kg.Passed() {
        return fmt.Errorf("86agbv0kg (frontend visualization) not passed")
    }
    
    // 3. Localization stability
    if !localizationStable() {
        return fmt.Errorf("localization not stable")
    }
    
    // 4. Frontend comparison UI stability
    if !frontendComparisonUIStable() {
        return fmt.Errorf("frontend comparison UI not stable")
    }
    
    return nil
}
```

### Wave 4 — v2 Paradigms
**Streams**: 86agbv0kh (v2 Architectural Paradigms Groundwork)
**Purpose**: Keep v2 architecture work staged but inactive until v1 maturity
**Focus Areas**: paradigm taxonomy, unlock prerequisites, future implementation boundaries

**Reactivation Rules (Executable)**:
```go
func ReactivateWave4() error {
    // 1. V1 backend-family maturity is proven
    if !v1BackendFamilyMaturityProven() {
        return fmt.Errorf("v1 backend-family maturity not proven")
    }
    
    // 2. Frontend comparison layer is stable
    if !frontendComparisonLayerStable() {
        return fmt.Errorf("frontend comparison layer not stable")
    }
    
    // 3. Paradigm taxonomy defined and accepted
    if !paradigmTaxonomyAccepted() {
        return fmt.Errorf("paradigm taxonomy not accepted")
    }
    
    // 4. V2 unlock prerequisites met (DEC-017)
    if !v2UnlockPrerequisitesMet() {
        return fmt.Errorf("v2 unlock prerequisites not met")
    }
    
    return nil
}
```

### Wave Stabilization Criteria
A wave is "stabilized" when:
- All its streams have passed their acceptance gates
- No critical bugs open for 14 days
- Test coverage ≥ 80% for all stream code
- Performance benchmarks within thresholds for 30 days
- No breaking API changes in last 3 releases

### Boundary Enforcement
- CI/CD pipelines for wave 2/3/4 tasks check reactivation rules before allowing work
- Task management system (ClickUp) shows wave status: `blocked` → `ready` → `active` → `stabilized`
- `/api/deferred-waves` returns current wave status and unmet prerequisites
- Manual reactivation requires signed attestation from all first-wave stream leads

## Consequences
- Wave boundaries are executable, not advisory
- No implementation work in a wave until its reactivation rules pass
- Stabilization criteria are measurable and time-bound
- Reactivation is a one-way gate (cannot revert to blocked without architectural decision)

## Acceptance Criteria
- [ ] `ReactivateWave2()`, `ReactivateWave3()`, `ReactivateWave4()` functions implemented
- [ ] CI enforces wave boundaries for all deferred tasks
- [ ] `/api/deferred-waves` returns accurate status and prerequisites
- [ ] Stabilization criteria measurable and tracked in CI
- [ ] Manual reactivation process documented with dual approval requirement