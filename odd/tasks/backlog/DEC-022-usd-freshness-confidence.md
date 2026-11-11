---
id: DEC-022
title: USD Freshness and Confidence Contract
status: proposed
date: 2026-11-11
---

## Context
Pricing enrichment (stream 86agbv0ke) requires a formal USD freshness contract with confidence calculation and staleness semantics. This decision formalizes the contract from the pricing demo in `main.go`.

## Decision
### Freshness Tiers
| Tier | Name | Max Age | Confidence | Use Case |
|------|------|---------|------------|----------|
| 0 | `real_time` | ≤ 1 second | 1.0 | High-frequency trading, critical decisions |
| 1 | `fresh` | ≤ 60 seconds | 0.95 | Dashboard display, comparison UI |
| 2 | `recent` | ≤ 300 seconds (5 min) | 0.80 | Batch reports, historical analysis |
| 3 | `stale` | ≤ 3600 seconds (1 hour) | 0.50 | Archive reference, non-critical |
| 4 | `very_stale` | ≤ 86400 seconds (24 hours) | 0.20 | Audit trail only |
| 5 | `expired` | > 86400 seconds | 0.00 | Must not be used; requires refresh |

### Confidence Calculation
```go
func CalculateConfidence(ageSeconds int, sourceReliability float64) float64 {
    // Base confidence from freshness tier
    var baseConfidence float64
    switch {
    case ageSeconds <= 1:
        baseConfidence = 1.0
    case ageSeconds <= 60:
        baseConfidence = 0.95
    case ageSeconds <= 300:
        baseConfidence = 0.80
    case ageSeconds <= 3600:
        baseConfidence = 0.50
    case ageSeconds <= 86400:
        baseConfidence = 0.20
    default:
        baseConfidence = 0.00
    }
    
    // Adjust for source reliability (0.0 - 1.0)
    // Source reliability: 1.0 = primary feed, 0.8 = secondary, 0.5 = tertiary
    adjusted := baseConfidence * sourceReliability
    
    // Floor at 0.0
    if adjusted < 0 {
        adjusted = 0
    }
    return adjusted
}
```

### Staleness Semantics
**Staleness Values** (aligned with `evidence.EvidenceStaleness`):
- `fresh`: Tier 0-1 (≤ 60s) — Safe for all decisions
- `stale`: Tier 2 (60-300s) — Safe for display, not for critical decisions
- `very_stale`: Tier 3-4 (300s-24h) — Reference only, flag in UI
- `reindexing`: Currently refreshing — Blocking, show spinner
- `expired`: Tier 5 (>24h) — Must not display; trigger refresh

### Pricing Data Structure
```json
{
  "currency_pair": "USD/COP",
  "rate": 4250.75,
  "source": "banxico|fixer|exchangerate_host|manual",
  "source_reliability": 0.95,
  "fetched_at": "2026-11-11T10:30:00Z",
  "freshness_tier": "fresh",
  "confidence": 0.9025,
  "staleness": "fresh",
  "age_seconds": 45,
  "evidence": {
    "source": "projection",
    "confidence": 0.9025,
    "staleness": "fresh",
    "measured_at": "2026-11-11T10:30:00Z",
    "delay_applied_ms": 120
  }
}
```

### Refresh Contract
**Refresh Triggers**:
1. **Age-based**: Tier ≥ `stale` (age > 300s) → background refresh
2. **Confidence-based**: Confidence < 0.5 → priority refresh
3. **Scheduled**: Every 60s for tier 0-1, every 300s for tier 2, every 3600s for tier 3-4
4. **Manual**: `/api/pricing/refresh` endpoint

**Refresh Behavior**:
- Async, non-blocking (returns stale data immediately, updates when fresh)
- Source priority: primary → secondary → tertiary
- On refresh failure: retain previous value, decrement confidence by 0.1, set staleness = `very_stale`
- Max 3 consecutive failures → circuit open, require manual intervention

### Fallback Chain
```
1. Primary feed (Banxico, official) — reliability 1.0
2. Secondary feed (Fixer.io) — reliability 0.85
3. Tertiary feed (ExchangeRate.host) — reliability 0.70
4. Manual override (admin) — reliability 0.50
5. Last known good — reliability 0.30 (confidence decays with age)
```

### Evidence Contract Integration
Every pricing response includes `EvidenceContract`:
- `source`: `SourceProjection` (pricing is a projection)
- `confidence`: Calculated per above
- `staleness`: Mapped from freshness tier
- `measured_at`: `fetched_at`
- `delay_applied_ms`: Time to fetch + compute (typically 50-200ms)

### API Contract
**`GET /api/pricing`**:
```json
{
  "rates": {
    "USD/COP": { ... },
    "USD/MXN": { ... }
  },
  "default_pair": "USD/COP",
  "freshness_summary": {
    "fresh_count": 8,
    "stale_count": 2,
    "very_stale_count": 0,
    "expired_count": 0
  },
  "evidence": { ... }
}
```

**`GET /api/pricing/{pair}`**: Single pair with full metadata

### Thresholds for UI
| Confidence | UI Treatment |
|------------|--------------|
| ≥ 0.9 | Green badge "Live" |
| ≥ 0.8 | Blue badge "Fresh" |
| ≥ 0.5 | Yellow badge "Delayed" |
| ≥ 0.2 | Orange badge "Stale" |
| < 0.2 | Red badge "Expired" — hide value, show refresh button |

## Consequences
- Freshness tiers are binding: UI MUST reflect tier via badges
- Confidence calculation is deterministic: same inputs = same output
- Fallback chain is ordered: never skip a higher-reliability source
- Expired data MUST NOT be displayed without explicit user action
- Evidence contracts enable audit trail for pricing decisions

## Acceptance Criteria
- [ ] Freshness tier calculation matches specification for all boundary values
- [ ] Confidence calculation verified with test matrix (age × source reliability)
- [ ] Staleness mapping to `evidence.EvidenceStaleness` correct
- [ ] Refresh triggers fire at correct intervals
- [ ] Fallback chain tested: primary failure → secondary → tertiary
- [ ] UI badges match confidence thresholds
- [ ] Evidence contracts on all pricing responses
- [ ] Circuit breaker after 3 consecutive failures