---
id: DEC-012
title: Pricing Dependency Gate
status: proposed
date: 2026-11-11
---

## Context
Pricing and cost enrichment (stream 86agbv0ke) is a late-stage enrichment stream that must not begin until backend-family work (86agbv0kg - frontend comparative visualization) and localization (86agbv0kf) are stable. This decision formalizes the dependency gate as an executable check.

## Decision
### Dependency Chain
```
86agbv0k9 (benchmark closure) 
    → 86agbv0kg (frontend comparative visualization) 
        → 86agbv0ke (pricing and cost enrichment)
             ↑
        86agbv0kf (localization) ───┘
```

### Gate Definition
```go
func Gate86agbv0ke() error {
    // 1. Backend-family work stable (86agbv0kg gate passed)
    if !gate86agbv0kg.Passed() {
        return fmt.Errorf("prerequisite 86agbv0kg (frontend comparative visualization) not passed")
    }
    
    // 2. Localization stable (86agbv0kf gate passed)
    if !gate86agbv0kf.Passed() {
        return fmt.Errorf("prerequisite 86agbv0kf (localization) not passed")
    }
    
    // 3. Benchmark closure accepted (86agbv0k9 gate passed)
    if !gate86agbv0k9.Passed() {
        return fmt.Errorf("prerequisite 86agbv0k9 (benchmark closure) not passed")
    }
    
    // 4. Pricing-specific stability checks
    if !localizationResourcesStable() {
        return fmt.Errorf("localization resources (en/es) not stable")
    }
    if !frontendComparisonUIStable() {
        return fmt.Errorf("frontend comparison UI not stable")
    }
    
    return nil
}
```

### Stability Criteria

#### Localization Stability (86agbv0kf)
- EN resource bundle complete and versioned
- ES resource bundle complete and versioned
- Fallback chain implemented: locale → region → default → key
- Pluralization rules defined for both locales
- Interpolation syntax standardized (e.g., `{{variable}}`)
- No breaking changes to resource keys in last 3 releases

#### Frontend Comparative Visualization Stability (86agbv0kg)
- Desktop comparison UI component APIs frozen (TreeGraphView, SplitView, RightPanel)
- FPS verification integrated (60fps target, 30fps minimum)
- Visible delay interactions implemented (delay metadata from benchmarks + PKI)
- Component contract tests passing
- No breaking API changes in last 3 releases

### Pricing Contract Prerequisites
Once gate passes, pricing work (86agbv0ke) may begin with these contracts:
1. **USD Freshness Contract** (DEC-022): Freshness tiers, confidence calculation, staleness semantics
2. **Pricing Verification/Fallback** (DEC-023): Verification checks, fallback chain, confidence thresholds
3. **Pricing Enrichment Integration**: Evidence contracts from DEC-011, PKI observability from DEC-004

### Enforcement
- CI/CD pipeline blocks 86agbv0ke tasks until `Gate86agbv0ke()` passes
- Gate status exposed via `/api/v2-readiness` (blockedBy includes pricing prerequisites)
- Manual override requires signed attestation from stream leads of 86agbv0kg and 86agbv0kf

## Consequences
- Pricing work cannot start until both backend visualization and localization are demonstrably stable
- Gate is executable, not advisory — CI enforces it
- Stability criteria are measurable (release count, test pass rate, API freeze)
- Manual override is audited and requires dual approval

## Acceptance Criteria
- [ ] `Gate86agbv0ke()` function implemented and returns correct status
- [ ] CI pipeline blocks pricing tasks when gate fails
- [ ] `/api/v2-readiness` shows pricing as blocked with specific unmet prerequisites
- [ ] Localization stability verified: EN/ES bundles versioned, fallback chain tested
- [ ] Frontend UI stability verified: component APIs frozen, FPS tests passing
- [ ] Manual override process documented and audited