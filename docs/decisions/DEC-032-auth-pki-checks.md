---
id: DEC-032
title: Auth and PKI Checks
status: proposed
date: 2026-11-11
---

## Context
Additional authentication and PKI checks beyond observability (86agbv0ng - done). This decision formalizes token validation, scope enforcement, and auth-specific checks distinct from PKI observability.

## Decision
### Auth Checks (from `keycloak.go` and `security.go`)

#### 1. Token Validation Checks
```go
type TokenValidationChecks struct {
    // Signature verification
    VerifySignature:      bool,    // JWS signature valid
    VerifyExpiration:     bool,    // exp claim not past
    VerifyNotBefore:      bool,    // nbf claim not future
    VerifyIssuer:         bool,    // iss matches expected
    VerifyAudience:       bool,    // aud matches client
    VerifyKeyID:          bool,    // kid exists in JWKS
    
    // Token type
    RequireAccessToken:   bool,    // Must be access token (not refresh)
    RequireBearerFormat:  bool,    // Authorization: Bearer <token>
    
    // PKCE (for Authorization Code flow)
    VerifyPKCE:           bool,    // code_verifier matches challenge
    
    // Claims
    RequiredClaims:       []string, // e.g., ["sub", "email", "roles"]
    AllowedClockSkew:     time.Duration, // Default: 30s
}
```

**Default Configuration**:
```go
func DefaultTokenValidationChecks() TokenValidationChecks {
    return TokenValidationChecks{
        VerifySignature:      true,
        VerifyExpiration:     true,
        VerifyNotBefore:      true,
        VerifyIssuer:         true,
        VerifyAudience:       true,
        VerifyKeyID:          true,
        RequireAccessToken:   true,
        RequireBearerFormat:  true,
        VerifyPKCE:           true,
        RequiredClaims:       []string{"sub", "email", "roles"},
        AllowedClockSkew:     30 * time.Second,
    }
}
```

#### 2. Scope Enforcement Checks
```go
type ScopeEnforcementChecks struct {
    // Required scopes per endpoint
    EndpointScopes: map[string][]string, // e.g., "/api/workflows" → ["workflow:read", "workflow:write"]
    
    // Scope hierarchy (parent implies children)
    ScopeHierarchy: map[string][]string, // e.g., "admin" → ["*"]
    
    // Default scopes if none specified
    DefaultScopes: []string, // e.g., ["user:read"]
    
    // Wildcard support
    AllowWildcard: bool, // "workflow:*" matches "workflow:read", "workflow:write"
}
```

**Default Scope Map** (from `keycloak.go` roles):
```go
var DefaultEndpointScopes = map[string][]string{
    "/api/healthz":                    {},
    "/api/bootstrap":                  {"user:read"},
    "/api/streams":                    {"stream:read"},
    "/api/catalog":                    {"catalog:read"},
    "/api/first-wave/contract":        {"stream:read"},
    "/api/first-wave/streams/":        {"stream:read"},
    "/api/security/bootstrap":         {"security:read"},
    "/api/security/pki/":              {"pki:read"},
    "/api/data-platform":              {"data:read"},
    "/api/benchmark-framework":        {"benchmark:read"},
    "/api/deferred-waves":             {"stream:read"},
    "/api/v2-readiness":               {"stream:read"},
    "/api/comparisons/realtime":       {"comparison:read"},
    "/api/transports/":                {"transport:read"},
    "/api/messaging/":                 {"messaging:read"},
    "/api/messaging/workflows":        {"workflow:read", "workflow:write"},
    "/api/messaging/failure-modes":    {"messaging:read"},
    "/api/sync/":                      {"sync:read"},
    "/api/sync/sessions":              {"sync:read", "sync:write"},
    "/api/projections/":               {"projection:read"},
    "/api/pricing":                    {"pricing:read"},
    "/api/v2/":                        {"v2:read"},
    "/api/async/":                     {"async:read", "async:write"},
    "/api/events":                     {"events:read"},
    "/api/ws/demo":                    {"transport:read"},
    "/api/benchmarks/metrics":         {"benchmark:read"},
}
```

#### 3. Session Validation Checks
```go
type SessionValidationChecks struct {
    // Session store checks
    VerifySessionExists:    bool,  // Session ID in store
    VerifySessionNotExpired: bool, // ExpiresAt > now
    VerifySessionActive:    bool,  // Status == active
    
    // Binding checks
    VerifyUserIDMatch:      bool,  // Token sub == Session UserID
    VerifyClientIPMatch:    bool,  // Token IP == Request IP (optional)
    VerifyUserAgentMatch:   bool,  // Token UA == Request UA (optional)
    
    // Refresh checks
    AllowRefresh:           bool,  // Allow refresh token rotation
    MaxRefreshes:           int,   // Max refreshes per session (default: 10)
    RefreshWindow:          time.Duration, // Refresh allowed within (default: 24h)
}
```

#### 4. Rate Limit Auth Checks
```go
type RateLimitAuthChecks struct {
    // Per-user limits
    UserLimit:       int,           // Requests per window
    UserWindow:      time.Duration, // Window duration
    
    // Per-client (IP) limits
    ClientLimit:     int,
    ClientWindow:    time.Duration,
    
    // Per-endpoint limits
    EndpointLimits:  map[string]RateLimitConfig,
    
    // Burst allowance
    BurstMultiplier: float64,       // e.g., 1.5x for burst
}

type RateLimitConfig struct {
    Limit:  int
    Window: time.Duration
}
```

### PKI Checks (Beyond Observability)
These are **distinct** from PKI observability (86agbv0ng) which monitors certificates and connections. These checks **enforce** PKI requirements.

#### 1. Certificate Validation Checks
```go
type CertificateValidationChecks struct {
    // Chain validation
    VerifyChain:           bool,  // Full chain to trusted root
    VerifyRootTrusted:     bool,  // Root in trust store
    VerifyIntermediate:    bool,  // Intermediate valid
    
    // Leaf validation
    VerifyLeafNotExpired:  bool,  // NotAfter > now
    VerifyLeafNotRevoked:  bool,  // Not in CRL/OCSP
    VerifyLeafKeyUsage:    bool,  // KeyUsage matches purpose
    VerifyLeafExtKeyUsage: bool,  // ExtKeyUsage matches (client/server auth)
    
    // SAN validation
    VerifySANMatches:      bool,  // DNS/IP matches expected
    AllowedSANs:           []string, // Expected SANs per service
    
    // Policy
    RequireCTLogs:         bool,  // Certificate Transparency logs
    MinKeySize:            int,   // e.g., 2048 for RSA, 256 for ECDSA
    AllowedAlgorithms:     []string, // e.g., ["RSA-2048", "ECDSA-P256"]
}
```

#### 2. mTLS Enforcement Checks
```go
type MTLSEnforcementChecks struct {
    // Per-connection requirements
    RequireClientCert:     bool,  // Server requires client cert
    RequireServerCert:     bool,  // Client verifies server cert
    VerifyMutualAuth:      bool,  // Both sides verified
    
    // Policy per trust boundary (from PKI observability)
    BoundaryRequirements: map[string]MTLSBoundaryRequirement,
    
    // Failure handling
    OnVerificationFailure: FailureAction, // reject | fallback | log
    FallbackToTLS:         bool,          // Allow TLS fallback (phase 1-2 only)
}

type MTLSBoundaryRequirement struct {
    Required:       bool,   // mTLS required for this boundary
    MinTLSVersion:  string, // "1.2" or "1.3"
    AllowedCiphers: []string,
    ClientAuthType: string, // "require", "request", "none"
}
```

#### 3. Certificate Rotation Checks
```go
type RotationChecks struct {
    // Proactive rotation
    RotateBeforeExpiry:   time.Duration, // e.g., 30 days (matches PKI observability)
    RotationWindow:       time.Duration, // Time to complete rotation
    
    // Validation post-rotation
    VerifyNewCertValid:   bool,
    VerifyOldCertRevoked: bool,
    VerifyClientsUpdated: bool, // All clients using new cert
    
    // Emergency rotation
    AllowEmergencyRotation: bool,
    EmergencyTrigger:       []string, // e.g., ["key_compromise", "ca_compromise"]
}
```

### Integration with Evidence Contracts
All auth/PKI checks produce `EvidenceContract`:
```go
func AuthCheckEvidence(checkName string, passed bool, details map[string]any) *EvidenceContract {
    ev := Evidence{
        Source:         SourceAudit,
        Confidence:     1.0,
        Staleness:      StalenessFresh,
        MeasuredAt:     time.Now(),
        DelayAppliedMs: 0,
    }
    if passed {
        return NewSuccessContract(ev, SuccessOK)
    }
    return NewErrorContract(ev, ErrorAuthFailed, fmt.Sprintf("auth check %s failed: %v", checkName, details))
}
```

### API Contract
**`GET /api/security/auth/checks`**: Returns all auth check configurations
**`GET /api/security/pki/checks`**: Returns all PKI check configurations
**`POST /api/security/auth/validate`**: Validates token + session + scopes
```json
{
  "token": "eyJhbGciOiJSUzI1NiIs...",
  "required_scopes": ["workflow:write"],
  "client_ip": "192.168.1.1"
}
```
→ Returns validation result with evidence

### Distinction from PKI Observability (86agbv0ng)
| Aspect | PKI Observability (86agbv0ng) | Auth/PKI Checks (DEC-032) |
|--------|-------------------------------|---------------------------|
| Purpose | Monitor, alert, visualize | Enforce, validate, gate |
| Scope | Certificates, connections, boundaries | Tokens, scopes, sessions, cert validation |
| Output | Metrics, dashboards, schedules | Pass/fail decisions, evidence |
| Timing | Continuous background | Per-request / per-connection |
| Failure Mode | Alert, log | Reject request, audit |

## Consequences
- Auth checks run on every protected request (middleware)
- PKI checks run on every mTLS connection (handshake)
- Scope enforcement is mandatory for all API endpoints
- Evidence contracts enable audit of every auth decision
- Distinct from observability: these are enforcement, not monitoring

## Acceptance Criteria
- [ ] Token validation checks implemented with default config
- [ ] Scope enforcement map covers all API endpoints
- [ ] Session validation checks integrated in middleware
- [ ] Rate limit auth checks operational
- [ ] Certificate validation checks enforce chain + leaf
- [ ] mTLS enforcement checks per trust boundary
- [ ] Rotation checks trigger proactive renewal
- [ ] All checks produce EvidenceContract
- [ ] Clear distinction documented from PKI observability