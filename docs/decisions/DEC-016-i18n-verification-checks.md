---
id: DEC-016
title: i18n Verification Checks
status: proposed
date: 2026-11-11
---

## Context
Localization (stream 86agbv0kf) requires formal verification of EN/ES resource bundles, locale detection, and fallback behavior. This decision formalizes the test contract for localization correctness.

## Decision
### Resource Bundle Structure
**File Layout**:
```
resources/
  en.json          # Base English (default)
  es.json          # Spanish
  en-US.json       # US English (region override)
  es-ES.json       # Spain Spanish (region override)
  es-MX.json       # Mexico Spanish (region override)
```

**Resource Key Format**: `{family}.{feature}.{key}` (e.g., `transport.sse.connecting`, `bootstrap.summary`)

**Value Types**:
- Simple string: `"Welcome to CommLayers"`
- Pluralization object: `{"one": "1 item", "other": "{count} items"}`
- Interpolation: `"Hello {{name}}, welcome to {{app}}"` (double-brace syntax)

### Locale Detection Chain
1. `Accept-Language` header (parsed per RFC 7231, quality values respected)
2. `?locale=` query parameter (explicit override)
3. Cookie `cl_locale` (persisted preference)
4. Default: `en`

**Fallback Chain** (exact order):
```
requested locale (e.g., es-MX)
    → region fallback (e.g., es-ES)
    → language fallback (e.g., en)
    → default locale (en)
    → key literal (if key missing entirely)
```

### Verification Checks (Executable)

#### 1. Resource Completeness
```go
func VerifyResourceCompleteness() error {
    baseKeys := extractKeys("resources/en.json")
    for _, locale := range []string{"es", "en-US", "es-ES", "es-MX"} {
        localeKeys := extractKeys(fmt.Sprintf("resources/%s.json", locale))
        missing := setDifference(baseKeys, localeKeys)
        if len(missing) > 0 {
            return fmt.Errorf("locale %s missing keys: %v", locale, missing)
        }
    }
    return nil
}
```

#### 2. Pluralization Rules
- EN: `one` (count == 1), `other` (all else)
- ES: `one` (count == 1), `other` (all else) — same as EN for now
- ICU MessageFormat compatible (future: `zero`, `few`, `many`)

#### 3. Interpolation Safety
```go
func VerifyInterpolationSafety() error {
    for _, locale := range allLocales {
        for key, value := range localeBundle {
            placeholders := extractPlaceholders(value) // {{...}}
            for _, ph := range placeholders {
                if !validPlaceholder(ph) {
                    return fmt.Errorf("invalid placeholder %s in %s.%s", ph, locale, key)
                }
            }
        }
    }
    return nil
}
```

#### 4. Locale Detection Correctness
| Input | Expected Output |
|-------|-----------------|
| `Accept-Language: es-MX,es;q=0.9,en;q=0.8` | `es-MX` |
| `Accept-Language: es;q=0.5, en;q=0.5` | `es` (first quality tie) |
| `?locale=fr` | `en` (fr not supported → fallback) |
| `Cookie: cl_locale=es-ES` + `Accept-Language: en` | `es-ES` (cookie wins) |

#### 5. Fallback Chain Correctness
```go
func TestFallbackChain() {
    // Key exists in es-MX → returns es-MX value
    // Key missing in es-MX, exists in es-ES → returns es-ES value
    // Key missing in es-MX, es-ES, exists in en → returns en value
    // Key missing everywhere → returns key literal
}
```

#### 6. No Hardcoded Strings in Code
- All user-facing strings MUST come from resource bundles
- Grep check: no string literals in handlers that match UI text patterns
- Exception: error codes, log messages, metric names

### API Contract
**Endpoint**: `/api/bootstrap` returns localized envelope
```json
{
  "summary": "CommLayers API v0.2.0 - First Wave Bootstrap",
  "phase": "first-wave materialized",
  "runtime": { ... },
  "now": { ... }
}
```
- `summary` value comes from `bootstrap.summary` in detected locale
- All localized responses wrap in `LocalizedEnvelope` with `locale` field

### Test Contract
**Test File**: `internal/i18n/verification_test.go`
- Runs in CI on every PR
- Fails build on any verification failure
- Generates coverage report for resource keys exercised

## Consequences
- Resource keys are the contract: changing a key breaks all locales
- New locales require full key parity with `en.json`
- Pluralization/interpolation syntax changes require all locales updated
- Fallback chain is immutable: no runtime reordering
- Locale detection priority is fixed: header → query → cookie → default

## Acceptance Criteria
- [ ] `VerifyResourceCompleteness()` passes for all supported locales
- [ ] Pluralization rules tested for EN/ES with counts 0, 1, 2, 5, 100
- [ ] Interpolation safety verified: no unclosed braces, no invalid chars
- [ ] Locale detection test matrix covers all priority combinations
- [ ] Fallback chain test passes for all 4 levels
- [ ] No hardcoded UI strings in handler code (grep check in CI)
- [ ] CI runs verification tests and fails on any violation