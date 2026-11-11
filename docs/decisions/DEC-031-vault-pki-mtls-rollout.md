---
id: DEC-031
title: Vault PKI/mTLS Rollout Plan
status: proposed
date: 2026-11-11
---

## Context
Extends PKI observability (86agbv0ng - done) with a 3-phase Vault PKI/mTLS rollout as an executable plan with gates. This decision formalizes the rollout plan.

## Decision
### Three-Phase Rollout (from `keycloak.go` and `pki_observability.go`)

#### Phase 1: Identity and Ingress TLS
**Goal**: Prepare service identity and trust boundary inventory
**Duration**: 2-4 weeks
**Scope**: Ingress gateway, auth-service, public-api
**Trust Boundary**: `ingress-tls` (services: api-gateway, auth-service, public-api)

**Requirements**:
- [ ] Vault PKI backend configured (offline root, online intermediate)
- [ ] Service identity certificates issued for ingress services
- [ ] Ingress TLS termination with service certificates
- [ ] Certificate tracking operational (expiry alerts at 30 days)
- [ ] MTLSStatusTracker initialized at phase 1
- [ ] TrustBoundaryTracker defines `ingress-tls` boundary
- [ ] PKI observability endpoints return data

**Gate Criteria** (executable):
```go
func Phase1Gate() error {
    // 1. Vault PKI operational
    if !vault.PKIBackend.Healthy() {
        return fmt.Errorf("vault pki backend not healthy")
    }
    // 2. Ingress services have valid certs
    services := []string{"api-gateway", "auth-service", "public-api"}
    for _, s := range services {
        cert := pki.GetCertificateForService(s)
        if cert == nil || cert.Status != "valid" {
            return fmt.Errorf("service %s missing valid cert", s)
        }
    }
    // 3. Ingress TLS terminating with service certs
    if !ingress.TLSUsesServiceCerts() {
        return fmt.Errorf("ingress not using service certificates")
    }
    // 4. PKI observability reporting
    if !pki.IsEnabled() {
        return fmt.Errorf("pki observability not enabled")
    }
    // 5. Trust boundary defined and evaluatable
    if !pki.TrustBoundaryTracker().HasBoundary("ingress-tls") {
        return fmt.Errorf("ingress-tls boundary not defined")
    }
    return nil
}
```

#### Phase 2: Selective mTLS
**Goal**: Protect critical service pairs without blocking bootstrap
**Duration**: 4-6 weeks (after Phase 1 gate)
**Scope**: Critical service pairs: auth-service ↔ vault, auth-service ↔ database, api-gateway ↔ cache
**Trust Boundary**: `critical-services` (services: auth-service, vault, database, cache)

**Requirements**:
- [ ] mTLS enabled for critical service pairs
- [ ] Certificate rotation automated (short-lived certs, 24h TTL)
- [ ] MTLSStatusTracker records connections for critical pairs
- [ ] TrustBoundaryTracker defines `critical-services` boundary
- [ ] PKI observability shows phase 2 connections
- [ ] Failure mode: mTLS failure → fallback to TLS (not block bootstrap)

**Gate Criteria** (executable):
```go
func Phase2Gate() error {
    // 1. Critical pairs have mTLS
    criticalPairs := [][]string{
        {"auth-service", "vault"},
        {"auth-service", "database"},
        {"api-gateway", "cache"},
    }
    for _, pair := range criticalPairs {
        conn := pki.GetMTLSConnection(pair[0], pair[1])
        if conn == nil || !conn.Verified {
            return fmt.Errorf("critical pair %s-%s not mTLS verified", pair[0], pair[1])
        }
    }
    // 2. Cert rotation automated (TTL ≤ 24h)
    for _, s := range []string{"auth-service", "vault", "database", "cache"} {
        cert := pki.GetCertificateForService(s)
        if cert.TTL > 24*time.Hour {
            return fmt.Errorf("service %s cert TTL too long: %v", s, cert.TTL)
        }
    }
    // 3. Trust boundary defined
    if !pki.TrustBoundaryTracker().HasBoundary("critical-services") {
        return fmt.Errorf("critical-services boundary not defined")
    }
    // 4. mTLS failure falls back gracefully (tested)
    if !mtls.FallbackTested() {
        return fmt.Errorf("mTLS fallback not tested")
    }
    return nil
}
```

#### Phase 3: Strict East-West mTLS
**Goal**: Enforce workload identity and authorization after observability is ready
**Duration**: 6-8 weeks (after Phase 2 gate)
**Scope**: ALL internal service communication
**Trust Boundary**: `east-west-mtls` (all services: api-gateway, auth-service, vault, database, cache, messaging, workflow, sync, projection-search, projection-graph, projection-vector)

**Requirements**:
- [ ] mTLS enforced for all east-west traffic
- [ ] Workload identity via SPIFFE/SPIRE or equivalent
- [ ] Authorization policies based on identity
- [ ] MTLSStatusTracker covers all service pairs
- [ ] TrustBoundaryTracker defines `east-west-mtls` boundary
- [ ] External edge boundary (`external-edge`) remains non-mTLS
- [ ] Certificate rotation fully automated (auto-renewal)
- [ ] Compliance: 100% internal connections verified

**Gate Criteria** (executable):
```go
func Phase3Gate() error {
    // 1. All internal pairs mTLS verified
    allServices := []string{
        "api-gateway", "auth-service", "vault", "database", "cache",
        "messaging", "workflow", "sync", "projection-search",
        "projection-graph", "projection-vector",
    }
    for i, a := range allServices {
        for j := i + 1; j < len(allServices); j++ {
            b := allServices[j]
            conn := pki.GetMTLSConnection(a, b)
            if conn == nil || !conn.Verified {
                return fmt.Errorf("pair %s-%s not mTLS verified", a, b)
            }
        }
    }
    // 2. Workload identity enforced
    if !workloadIdentity.Enforced() {
        return fmt.Errorf("workload identity not enforced")
    }
    // 3. Authorization policies active
    if !authz.PoliciesActive() {
        return fmt.Errorf("authorization policies not active")
    }
    // 4. External edge remains non-mTLS
    edgeConn := pki.GetMTLSConnection("api-gateway", "cdn")
    if edgeConn != nil && edgeConn.Verified {
        return fmt.Errorf("external edge should not be mTLS")
    }
    // 5. Compliance 100%
    boundaries := pki.GetTrustBoundaries()
    for _, b := range boundaries {
        if b.Name == "east-west-mtls" && !b.Compliant {
            return fmt.Errorf("east-west-mtls boundary not compliant")
        }
    }
    return nil
}
```

### Rollout State Machine
```
phase_1_pending → phase_1_active → phase_1_gate_passed
    → phase_2_pending → phase_2_active → phase_2_gate_passed
        → phase_3_pending → phase_3_active → phase_3_gate_passed
            → rollout_complete
```

**Rollback Triggers** (any phase):
- Critical service unavailable due to mTLS
- Certificate rotation failure > 1 hour
- Compliance < 95% for 24 hours
- Security incident related to mTLS

### PKI Observability Integration (from 86agbv0ng)
**Monitoring During Rollout**:
- Certificate expiry alerts (30 days default)
- Connection verification rate (target 100% after phase gate)
- Trust boundary compliance (per phase)
- Rotation schedule adherence

**API Endpoints** (from `main.go`):
- `GET /api/security/pki/certificates` — All tracked certs
- `GET /api/security/pki/mtls-status` — Connection matrix + phase
- `GET /api/security/pki/trust-boundaries` — Boundary compliance
- `GET /api/security/pki/rotation-schedule` — Upcoming rotations
- `GET /api/security/pki/summary` — Aggregate state

### Vault PKI Configuration
```hcl
# Vault PKI Backend Config
path "pki_int/" {
  type = "pki"
  config = {
    max_lease_ttl = "8760h"  # 1 year for intermediate
    default_lease_ttl = "24h" # 24h for leaf certs
  }
}

# Roles per service
path "pki_int/role/*" {
  capabilities = ["read", "create", "update"]
  allowed_domains = "commlayers.internal"
  allow_subdomains = true
  max_ttl = "24h"
  generate_lease = true
}
```

### Evidence Integration
All PKI events produce `EvidenceContract`:
- `source`: `SourceAudit`
- `confidence`: 1.0 for verified, 0.5 for unverified
- `staleness`: From certificate/connection timestamp
- `measured_at`: Event timestamp
- `delay_applied_ms`: Operation latency

## Consequences
- Phases are sequential: Phase N+1 cannot start until Phase N gate passes
- Gates are executable: CI/CD checks before promotion
- Rollback is automatic on triggers
- PKI observability is the single source of truth for rollout status
- External edge intentionally excluded from mTLS

## Acceptance Criteria
- [ ] All three phase gate functions implemented and testable
- [ ] Vault PKI backend configured with correct TTLs
- [ ] Trust boundaries match specification exactly
- [ ] mTLS fallback tested and documented
- [ ] PKI observability endpoints return correct data per phase
- [ ] Rollback triggers documented and tested
- [ ] Workload identity integration defined for phase 3