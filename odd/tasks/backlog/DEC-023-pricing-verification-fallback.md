---
id: DEC-023
title: Pricing Verification and Fallback
status: proposed
date: 2026-11-11
---

## Context
Pricing verification checks, fallback chain, and confidence thresholds must integrate with evidence contracts (DEC-011) and the USD freshness contract (DEC-022). This decision formalizes the verification and fallback specification.

## Decision
### Verification Checks (Executable)
Every pricing fetch runs verification checks before accepting the rate:

#### 1. Cross-Source Variance Check
```go
func VerifyCrossSourceVariance(rates map[string]float64, threshold float64) error {
    // Compare primary vs secondary vs tertiary for same pair
    for pair, primaryRate := range rates {
        if secondaryRate, ok := secondaryRates[pair]; ok {
            variance := math.Abs(primaryRate - secondaryRate) / primaryRate
            if variance > threshold { // Default: 0.02 (2%)
                return fmt.Errorf("cross-source variance %.2f%% exceeds threshold %.2f%% for %s",
                    variance*100, threshold*100, pair)
            }
        }
    }
    return nil
}
```

#### 2. Rate Sanity Check
```go
func VerifyRateSanity(pair string, rate float64) error {
    bounds := sanityBounds[pair] // e.g., USD/COP: 3000-6000
    if rate < bounds.Min || rate > bounds.Max {
        return fmt.Errorf("rate %.4f for %s outside sanity bounds [%.4f, %.4f]",
            rate, pair, bounds.Min, bounds.Max)
    }
    return nil
}
```

#### 3. Monotonicity Check (for time-series)
```go
func VerifyMonotonicity(pair string, newRate, prevRate float64, maxChangePct float64) error {
    change := math.Abs(newRate - prevRate) / prevRate
    if change > maxChangePct { // Default: 0.05 (5% per fetch interval)
        return fmt.Errorf("rate change %.2f%% exceeds max %.2f%% for %s",
            change*100, maxChangePct*100, pair)
    }
    return nil
}
```

#### 4. Freshness Check
```go
func VerifyFreshness(fetchedAt time.Time, maxAge time.Duration) error {
    age := time.Since(fetchedAt)
    if age > maxAge { // Default: 3600s (1 hour) for acceptance
        return fmt.Errorf("rate age %v exceeds max age %v", age, maxAge)
    }
    return nil
}
```

#### 5. Evidence Contract Validation
```go
func VerifyEvidenceContract(ev *evidence.EvidenceContract) error {
    if ev == nil {
        return fmt.Errorf("missing evidence contract")
    }
    if ev.Evidence.Source != evidence.SourceProjection {
        return fmt.Errorf("invalid evidence source: %s", ev.Evidence.Source)
    }
    if ev.Evidence.Confidence < 0.1 { // Absolute floor
        return fmt.Errorf("confidence %.2f below floor", ev.Evidence.Confidence)
    }
    if ev.Evidence.Staleness == evidence.StalenessExpired {
        return fmt.Errorf("evidence expired")
    }
    return nil
}
```

### Fallback Chain (Ordered, from DEC-022)
| Priority | Source | Reliability | Verification | On Failure |
|----------|--------|-------------|--------------|------------|
| 1 | Primary (Banxico) | 1.0 | All 5 checks | → Priority 2 |
| 2 | Secondary (Fixer.io) | 0.85 | All 5 checks | → Priority 3 |
| 3 | Tertiary (ExchangeRate.host) | 0.70 | Checks 2,3,4,5 | → Priority 4 |
| 4 | Manual Override | 0.50 | Checks 2,3,5 | → Priority 5 |
| 5 | Last Known Good | 0.30* | Check 5 only | → Circuit Open |

*Last Known Good reliability decays: `0.30 * (1 - age_hours/24)` minimum 0.05

### Confidence Thresholds (Gate Actions)
| Confidence | Tier | Action |
|------------|------|--------|
| ≥ 0.90 | `verified` | Accept, cache, serve normally |
| ≥ 0.75 | `provisional` | Accept, cache, flag for re-verify in 60s |
| ≥ 0.50 | `degraded` | Accept, cache, serve with warning badge, background refresh |
| ≥ 0.20 | `reference` | Serve with stale badge, hide from calculations, priority refresh |
| < 0.20 | `rejected` | Reject, circuit open, alert, require manual |

### Circuit Breaker
```go
type PricingCircuitBreaker struct {
    consecutiveFailures int
    lastFailure         time.Time
    state               CircuitState // closed, open, half_open
}

func (cb *PricingCircuitBreaker) RecordFailure() {
    cb.consecutiveFailures++
    cb.lastFailure = time.Now()
    if cb.consecutiveFailures >= 3 {
        cb.state = CircuitOpen
    }
}

func (cb *PricingCircuitBreaker) RecordSuccess() {
    cb.consecutiveFailures = 0
    cb.state = CircuitClosed
}

func (cb *PricingCircuitBreaker) AllowRequest() bool {
    if cb.state == CircuitClosed {
        return true
    }
    if cb.state == CircuitOpen {
        // Half-open after 5 minutes
        if time.Since(cb.lastFailure) > 5*time.Minute {
            cb.state = CircuitHalfOpen
            return true
        }
        return false
    }
    // Half-open: allow one request
    return true
}
```

### Integration with Evidence Contracts (DEC-011)
- All pricing responses include `EvidenceContract` with:
  - `source`: `SourceProjection`
  - `confidence`: Post-verification confidence
  - `staleness`: From freshness tier
  - `measured_at`: Fetch timestamp
  - `delay_applied_ms`: Fetch + verification latency
- Verification failures produce `ErrorContract` with:
  - `ErrorCode`: `ErrorStaleData` (freshness), `ErrorInternal` (verification), `ErrorUnavailable` (circuit open)
  - `ErrorMessage`: Specific failure reason

### API Contract
**`GET /api/pricing/verify`**: Runs all verifications on current rates
```json
{
  "verified": true,
  "checks": [
    {"name": "cross_source_variance", "passed": true, "variance_pct": 0.5},
    {"name": "rate_sanity", "passed": true},
    {"name": "monotonicity", "passed": true},
    {"name": "freshness", "passed": true, "age_seconds": 45},
    {"name": "evidence_contract", "passed": true}
  ],
  "overall_confidence": 0.9025,
  "circuit_state": "closed"
}
```

**`POST /api/pricing/refresh`**: Forces refresh with verification
```json
{
  "triggered": true,
  "source_attempted": "primary",
  "verification": { ... },
  "new_confidence": 0.95
}
```

## Consequences
- Verification runs on every fetch (not just scheduled)
- Fallback chain is strict: never skip a higher-priority source
- Confidence thresholds gate UI behavior AND API acceptance
- Circuit breaker prevents cascade failures
- Evidence contracts provide audit trail for every pricing decision

## Acceptance Criteria
- [ ] All 5 verification checks implemented and tested
- [ ] Fallback chain tested: primary→secondary→tertiary→manual→last-known
- [ ] Confidence thresholds map to correct actions
- [ ] Circuit breaker opens after 3 failures, half-opens after 5 min
- [ ] `/api/pricing/verify` returns correct check results
- [ ] Evidence contracts on all pricing responses with correct fields
- [ ] Integration with DEC-011 success/error taxonomy verified