---
id: DEC-010
title: Layer Acceptance Gates
status: proposed
date: 2026-11-11
---

## Context
The first wave consists of four streams that must be accepted in sequence before any deferred wave implementation begins. This decision formalizes the acceptance gates for each stream as executable checks.

## Decision
### Stream 86agbv0k4 — Backend Abstraction Layer Matrix (Wave Entry Gate)
**Prerequisites**: None
**Frozen Vocabulary**: `BackendFamily` enum (transport, messaging, sync, search, benchmark), comparison dimensions (consistency, latency class, durability semantics, operational burden, local-compose viability)
**Explicit Contracts**:
- Layer responsibility matrix (which family owns which capability)
- Capability comparison rubric (scoring methodology)
- Default and fallback decision criteria
- Downstream dependency contract (what later streams can assume)
**Acceptance Gates (executable)**:
```go
func Gate86agbv0k4() error {
    // 1. Every backend family has a comparison entry
    families := []BackendFamily{BackendFamilyTransport, BackendFamilyMessaging, BackendFamilySync, BackendFamilySearch, BackendFamilyBenchmark}
    for _, f := range families {
        if !comparisonRubric.HasEntry(f) {
            return fmt.Errorf("missing comparison rubric for %s", f)
        }
    }
    // 2. Vocabulary is frozen: no new BackendFamily values without gate re-approval
    if backendFamilyRegistry.IsFrozen() == false {
        return fmt.Errorf("vocabulary not frozen")
    }
    // 3. Dependency contract published
    if !dependencyContract.IsPublished() {
        return fmt.Errorf("downstream dependency contract not published")
    }
    return nil
}
```

### Stream 86agbv0k3 — Keycloak and Vault Security Bootstrap (Controlled Overlap)
**Prerequisites**: 86agbv0k4 gate passed
**Frozen Vocabulary**: `PKIConfig` (ExpiryAlertDays, InitialMTLSPhase), `MTLSPhase` (1, 2, 3), trust boundary names
**Explicit Contracts**:
- Keycloak realm and client baseline (BFF-oriented Authorization Code + PKCE)
- Vault secret ownership strategy (service secrets, bootstrap credentials, future PKI authority)
- Rotation and revocation posture (offline root, online intermediate, short-lived service certificates)
- Bootstrap trust-boundary contract (ingress-tls, critical-services, east-west-mtls, external-edge)
**Acceptance Gates (executable)**:
```go
func Gate86agbv0k3() error {
    // 1. Realm, client, service identity model explicit
    if !keycloak.IdentityModel.IsExplicit() {
        return fmt.Errorf("identity model not explicit")
    }
    // 2. Vault paths, secret classes, rotation posture explicit
    if !vault.OwnershipStrategy.IsExplicit() {
        return fmt.Errorf("vault ownership not explicit")
    }
    // 3. Bootstrap trust boundaries explicit and evaluatable
    boundaries := pki.GetTrustBoundaries()
    required := []string{"ingress-tls", "critical-services", "east-west-mtls", "external-edge"}
    for _, r := range required {
        if !boundaries.Has(r) {
            return fmt.Errorf("missing trust boundary: %s", r)
        }
    }
    return nil
}
```

### Stream 86agbv0k7 — Canonical Data Platform (Controlled Overlap)
**Prerequisites**: 86agbv0k4 gate passed
**Frozen Vocabulary**: `CanonicalPostgresDomainModel` kinds, `BackendFamily` ownership matrix, schema versioning policy
**Explicit Contracts**:
- Canonical entity catalog (stream definitions, comparison evidence, job lifecycle, deferred-wave manifest)
- Storage ownership matrix (PostgreSQL = system of record; Redis = ephemeral; optional profiles = ES/Neo4j/Qdrant/Mongo)
- Schema evolution rules (explicit ownership, versioned evolution, audit lineage required)
- Audit lineage baseline (source, confidence, staleness, measuredAt, delayAppliedMs)
**Acceptance Gates (executable)**:
```go
func Gate86agbv0k7() error {
    // 1. Canonical entity ownership explicit
    entities := []string{"stream_definitions", "comparison_evidence", "job_lifecycle", "deferred_wave_manifest"}
    for _, e := range entities {
        if !store.Ownership.IsExplicit(e) {
            return fmt.Errorf("ownership not explicit for %s", e)
        }
    }
    // 2. Schema evolution rules explicit
    if !store.SchemaRules.AreExplicit() {
        return fmt.Errorf("schema evolution rules not explicit")
    }
    // 3. Audit lineage fields explicit
    requiredFields := []string{"source", "confidence", "staleness", "measured_at", "delay_applied_ms"}
    for _, f := range requiredFields {
        if !evidence.Contract.HasField(f) {
            return fmt.Errorf("missing audit field: %s", f)
        }
    }
    // 4. Compatibility rules across backend families explicit
    if !store.CompatibilityRules.AreExplicit() {
        return fmt.Errorf("compatibility rules not explicit")
    }
    return nil
}
```

### Stream 86agbv0k9 — Benchmark and Observability Framework (Wave Closure Gate)
**Prerequisites**: 86agbv0k4, 86agbv0k3, 86agbv0k7 gates passed
**Frozen Vocabulary**: Benchmark scenario names, telemetry schema, success/error taxonomy, threshold values
**Explicit Contracts**:
- Benchmark matrix (request baseline, SSE startup/event cadence, async visibility)
- Metric, log, and trace naming contract (request_duration_ms, event_flush_ms, job_progress, delay_applied_ms, etc.)
- Reproducibility rules (fixed sample intervals, ring buffer sizes, warm-up runs)
- Wave closure report template
**Acceptance Gates (executable)**:
```go
func Gate86agbv0k9() error {
    // 1. Benchmark scenarios fixed
    scenarios := []string{"request_baseline", "sse_startup_cadence", "async_visibility"}
    for _, s := range scenarios {
        if !benchmarkMatrix.HasScenario(s) {
            return fmt.Errorf("missing benchmark scenario: %s", s)
        }
    }
    // 2. Telemetry minimums fixed
    requiredMetrics := []string{"request_duration_ms", "event_flush_ms", "job_progress", "delay_applied_ms"}
    requiredLogs := []string{"request_id", "stream_id", "job_id", "phase", "result"}
    requiredTraces := []string{"bootstrap.http", "bootstrap.async", "bootstrap.sse"}
    if !telemetryContract.HasAll(requiredMetrics, requiredLogs, requiredTraces) {
        return fmt.Errorf("telemetry minimums not met")
    }
    // 3. Success and error taxonomy fixed
    if !evidence.ValidSuccessCodes().Equal([]SuccessCode{SuccessOK, SuccessPartial, SuccessDegraded}) {
        return fmt.Errorf("success taxonomy drift")
    }
    if !evidence.ValidErrorCodes().Equal(standardErrorCodes) {
        return fmt.Errorf("error taxonomy drift")
    }
    // 4. Threshold and reporting format stable
    if !benchmarkThresholds.AreStable() {
        return fmt.Errorf("thresholds not stable")
    }
    return nil
}
```

## Consequences
- No deferred wave (wave 2, 3, 4) stream may begin implementation until its prerequisite gates pass
- Gates are executable checks, not documentation — they run in CI/CD
- Vocabulary freeze means: no new enum values, no new struct fields, no new config keys without gate re-approval
- Each gate produces a signed attestation artifact stored in the canonical data platform

## Acceptance Criteria
- [ ] All four gate functions compile and pass against current codebase
- [ ] Gate 86agbv0k4 runs first in CI pipeline; subsequent gates depend on it
- [ ] Vocabulary freeze enforced: build fails if new `BackendFamily` values added without gate re-run
- [ ] Attestation artifacts generated and stored for each gate
- [ ] `/api/first-wave/contract` returns gate statuses and execution order