---
id: DEC-021
title: EN/ES Resource Model
status: proposed
date: 2026-11-11
---

## Context
Localization (stream 86agbv0kf) requires a formal resource bundle structure with pluralization, interpolation, and a defined fallback chain. This decision formalizes the EN/ES resource model as a binding contract.

## Decision
### Resource Bundle File Structure
```
resources/
├── en.json          # Base language (default fallback)
├── es.json          # Base Spanish
├── en-US.json       # Region: US English
├── en-GB.json       # Region: UK English
├── es-ES.json       # Region: Spain Spanish
├── es-MX.json       # Region: Mexico Spanish
├── es-AR.json       # Region: Argentina Spanish
└── ...
```

### Resource Key Schema
**Format**: `{domain}.{feature}.{key}`
**Examples**:
- `bootstrap.summary`
- `transport.sse.connecting`
- `transport.sse.connected`
- `workflow.queued`
- `workflow.running`
- `workflow.completed`
- `sync.aligned`
- `sync.conflict`
- `error.timeout`
- `error.unavailable`

**Key Rules**:
- Lowercase with dots as separators
- No spaces, no special characters except dots
- Domain = communication family or cross-cutting concern
- Feature = functional area within domain
- Key = specific message

### Value Types

#### 1. Simple String
```json
"bootstrap.summary": "CommLayers API v0.2.0 - First Wave Bootstrap"
```

#### 2. Pluralization Object (ICU MessageFormat compatible)
```json
"workflow.jobs_processed": {
  "one": "1 job processed",
  "other": "{{count}} jobs processed"
}
```
**Plural Categories** (CLDR):
- EN: `one` (n=1), `other` (all else)
- ES: `one` (n=1), `other` (all else)
- Future: `zero`, `few`, `many` per CLDR rules

#### 3. Interpolation (Double-Brace Syntax)
```json
"transport.session_created": "Session {{session_id}} created at {{timestamp}}"
"sync.conflict_detected": "Conflict in session {{session_id}}: {{conflict_type}}"
```
**Placeholder Rules**:
- Format: `{{variable_name}}`
- Variable name: `[a-z_][a-z0-9_]*` (lowercase, underscore, alphanumeric)
- No nested braces, no expressions
- Missing placeholder at render time → empty string (logged warning)

### Fallback Chain (Exact Order, Immutable)
```
1. Exact locale match (e.g., es-MX)
2. Region fallback (e.g., es-ES for es-MX)
3. Language fallback (e.g., es for es-ES)
4. Default locale (en)
5. Key literal (if key missing entirely: "bootstrap.summary")
```

**Fallback Implementation**:
```go
func ResolveMessage(locale, key string, data map[string]any) string {
    // 1. Exact
    if msg, ok := bundles[locale][key]; ok {
        return render(msg, data)
    }
    // 2. Region (strip variant, keep language-region)
    if region := extractRegion(locale); region != "" {
        if msg, ok := bundles[region][key]; ok {
            return render(msg, data)
        }
    }
    // 3. Language (strip region)
    if lang := extractLanguage(locale); lang != "" && lang != "en" {
        if msg, ok := bundles[lang][key]; ok {
            return render(msg, data)
        }
    }
    // 4. Default (en)
    if msg, ok := bundles["en"][key]; ok {
        return render(msg, data)
    }
    // 5. Key literal
    return key
}
```

### Locale Detection Priority (Exact Order)
1. `?locale=` query parameter (explicit, highest)
2. `Accept-Language` header (RFC 7231, quality values)
3. `cl_locale` cookie (persisted preference)
4. Default: `en`

### Bundle Completeness Contract
**Mandatory**: Every supported locale MUST have all keys present in `en.json`
**Verification**: `VerifyResourceCompleteness()` (DEC-016) runs in CI
**Exception**: Region-specific overrides only (e.g., `es-MX` may override `es` keys)

### Rendering Function Signature
```go
func Render(locale, key string, data map[string]any) string
```
- Returns rendered string with pluralization and interpolation applied
- Logs warning for missing keys, missing placeholders, invalid pluralization
- Never panics; returns key literal as ultimate fallback

### API Integration
All localized endpoints return `LocalizedEnvelope`:
```json
{
  "locale": "es-MX",
  "data": { ... },
  "messages": {
    "summary": "CommLayers API v0.2.0 - Primer Bootstrap de Onda"
  }
}
```

### Versioning
- Bundle version in `resources/version.json`: `{"version": "1.0.0", "base_locale": "en"}`
- Breaking changes (key rename/remove) → major version
- New keys → minor version
- Fixes/overrides → patch version
- Version embedded in `LocalizedEnvelope` for debugging

## Consequences
- Resource keys are immutable once released (only additions allowed in minor/patch)
- Pluralization follows CLDR; no custom rules
- Interpolation syntax is fixed: double-brace, simple variable names
- Fallback chain is immutable; no runtime configuration
- Locale detection priority is fixed; no middleware reordering

## Acceptance Criteria
- [ ] Bundle structure matches specification exactly
- [ ] Pluralization works for EN/ES with CLDR categories
- [ ] Interpolation renders correctly with nested data
- [ ] Fallback chain order verified in integration tests
- [ ] Locale detection priority matrix tested
- [ ] `VerifyResourceCompleteness()` passes for all locales
- [ ] Version tracking in `resources/version.json`
- [ ] All localized endpoints return `LocalizedEnvelope` with `locale` field