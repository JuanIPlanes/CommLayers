package security

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"crypto/tls"
)

// PKIObservability provides unified PKI/mTLS observability.
type PKIObservability struct {
	mu            sync.RWMutex
	enabled       bool
	certTracker   *CertificateTracker
	mtlsTracker   *MTLSStatusTracker
	trustBoundary *TrustBoundaryTracker
	auditLogger   AuditLogger
}

// AuditLogger interface for audit logging (reused from keycloak).
type AuditLogger interface {
	Log(ctx context.Context, event AuditEvent) error
}

// AuditEvent represents an audit log entry (reused from keycloak).
type AuditEvent struct {
	ID        string                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	UserID    string                 `json:"user_id,omitempty"`
	SessionID string                 `json:"session_id,omitempty"`
	Action    string                 `json:"action"`
	Resource  string                 `json:"resource"`
	Method    string                 `json:"method"`
	Path      string                 `json:"path"`
	IP        string                 `json:"ip"`
	UserAgent string                 `json:"user_agent"`
	Status    int                    `json:"status"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// PKIConfig holds configuration for PKI observability.
type PKIConfig struct {
	// ExpiryAlertDays is the number of days before expiry to alert.
	ExpiryAlertDays int
	// InitialMTLSPhase is the starting mTLS rollout phase (1-3).
	InitialMTLSPhase int
	// AuditLogger is the audit logger to use.
	AuditLogger AuditLogger
}

// DefaultPKIConfig returns default PKI configuration.
func DefaultPKIConfig() *PKIConfig {
	return &PKIConfig{
		ExpiryAlertDays:   30,
		InitialMTLSPhase:  1,
		AuditLogger:       &defaultPKIAuditLogger{},
	}
}

var (
	PKIInstance *PKIObservability
	PKIOnce     sync.Once
)

// InitPKIObservability initializes the global PKI observability instance.
func InitPKIObservability(config *PKIConfig) {
	PKIOnce.Do(func() {
		if config == nil {
			config = DefaultPKIConfig()
		}

		expiryAlert := time.Duration(config.ExpiryAlertDays) * 24 * time.Hour
		certTracker := NewCertificateTracker(expiryAlert)
		mtlsTracker := NewMTLSStatusTracker(config.InitialMTLSPhase)
		trustBoundary := NewTrustBoundaryTracker(mtlsTracker, certTracker)

		PKIInstance = &PKIObservability{
			enabled:       true,
			certTracker:   certTracker,
			mtlsTracker:   mtlsTracker,
			trustBoundary: trustBoundary,
			auditLogger:   config.AuditLogger,
		}

		// Define default trust boundaries for the 3-phase mTLS rollout
		PKIInstance.defineDefaultBoundaries()

		log.Printf("PKI Observability initialized: phase=%d, expiry_alert=%d days",
			config.InitialMTLSPhase, config.ExpiryAlertDays)
	})
}

// GetPKIObservability returns the global PKI observability instance.
func GetPKIObservability() *PKIObservability {
	return PKIInstance
}

// IsEnabled reports if PKI observability has been initialized.
func (p *PKIObservability) IsEnabled() bool {
	if p == nil {
		return false
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.enabled
}

// CertificateTracker returns the certificate tracker.
func (p *PKIObservability) CertificateTracker() *CertificateTracker {
	return p.certTracker
}

// MTLSStatusTracker returns the mTLS status tracker.
func (p *PKIObservability) MTLSStatusTracker() *MTLSStatusTracker {
	return p.mtlsTracker
}

// TrustBoundaryTracker returns the trust boundary tracker.
func (p *PKIObservability) TrustBoundaryTracker() *TrustBoundaryTracker {
	return p.trustBoundary
}

// SetMTLSPhase updates the mTLS rollout phase.
func (p *PKIObservability) SetMTLSPhase(phase int) {
	p.mtlsTracker.SetPhase(phase)
	p.auditLog(context.Background(), "pki.mtls_phase_changed", map[string]interface{}{
		"new_phase": phase,
	})
}

// RecordCertificate records a certificate from PEM data.
func (p *PKIObservability) RecordCertificate(ctx context.Context, pemData []byte, source, service string) (*CertificateInfo, error) {
	info, err := p.certTracker.TrackCertificatePEM(ctx, pemData, source, service)
	if err != nil {
		p.auditLog(ctx, "pki.certificate_record_failed", map[string]interface{}{
			"source": source,
			"service": service,
			"error": err.Error(),
		})
		return nil, err
	}

	p.auditLog(ctx, "pki.certificate_recorded", map[string]interface{}{
		"serial": info.SerialNumber,
		"subject": info.Subject,
		"issuer": info.Issuer,
		"not_after": info.NotAfter,
		"source": source,
		"service": service,
	})
	return info, nil
}

// RecordMTLSConnection records an mTLS connection observation.
func (p *PKIObservability) RecordMTLSConnection(ctx context.Context, source, target string, state *tls.ConnectionState) {
	p.mtlsTracker.RecordConnection(ctx, source, target, state)

	conn, _ := p.mtlsTracker.GetConnection(source, target)
	p.auditLog(ctx, "pki.mtls_connection_observed", map[string]interface{}{
		"source": source,
		"target": target,
		"tls_version": conn.TLSVersion,
		"cipher_suite": conn.CipherSuite,
		"verified": conn.Verified,
		"cert_expiry_days": conn.CertExpiryDays,
		"phase": conn.Phase,
	})
}

// EvaluateTrustBoundaries evaluates all trust boundary compliance.
func (p *PKIObservability) EvaluateTrustBoundaries(ctx context.Context) {
	p.trustBoundary.EvaluateCompliance(ctx)

	for _, boundary := range p.trustBoundary.GetAllBoundaries() {
		p.auditLog(ctx, "pki.trust_boundary_evaluated", map[string]interface{}{
			"name": boundary.Name,
			"compliant": boundary.Compliant,
			"total_pairs": boundary.TotalPairs,
			"compliant_pairs": boundary.CompliantPairs,
			"violations": boundary.Violations,
		})
	}
}

// GetCertificates returns all tracked certificates.
func (p *PKIObservability) GetCertificates() []*CertificateInfo {
	return p.certTracker.GetAllCertificates()
}

// GetMTLSStatus returns all mTLS connection statuses.
func (p *PKIObservability) GetMTLSStatus() []*MTLSConnection {
	return p.mtlsTracker.GetAllConnections()
}

// GetTrustBoundaries returns all trust boundaries with compliance status.
func (p *PKIObservability) GetTrustBoundaries() []*TrustBoundary {
	return p.trustBoundary.GetAllBoundaries()
}

// GetRotationSchedule returns upcoming certificate rotations.
func (p *PKIObservability) GetRotationSchedule() []*CertificateInfo {
	return p.certTracker.GetRotationSchedule()
}

// GetPKISummary returns a summary of PKI observability state.
func (p *PKIObservability) GetPKISummary() map[string]interface{} {
	certs := p.certTracker.GetAllCertificates()
	conns := p.mtlsTracker.GetAllConnections()
	boundaries := p.trustBoundary.GetAllBoundaries()

	validCerts := 0
	expiringCerts := 0
	expiredCerts := 0
	revokedCerts := 0
	for _, c := range certs {
		switch c.Status {
		case "valid":
			validCerts++
		case "expiring":
			expiringCerts++
		case "expired":
			expiredCerts++
		case "revoked":
			revokedCerts++
		}
	}

	verifiedConns := 0
	unverifiedConns := 0
	for _, c := range conns {
		if c.Verified {
			verifiedConns++
		} else {
			unverifiedConns++
		}
	}

	compliantBoundaries := 0
	nonCompliantBoundaries := 0
	for _, b := range boundaries {
		if b.Compliant {
			compliantBoundaries++
		} else {
			nonCompliantBoundaries++
		}
	}

	return map[string]interface{}{
		"certificates": map[string]interface{}{
			"total":     len(certs),
			"valid":     validCerts,
			"expiring":  expiringCerts,
			"expired":   expiredCerts,
			"revoked":   revokedCerts,
		},
		"mtls_connections": map[string]interface{}{
			"total":      len(conns),
			"verified":   verifiedConns,
			"unverified": unverifiedConns,
			"current_phase": p.mtlsTracker.GetPhase(),
		},
		"trust_boundaries": map[string]interface{}{
			"total":             len(boundaries),
			"compliant":         compliantBoundaries,
			"non_compliant":     nonCompliantBoundaries,
		},
		"rotation_schedule": len(p.certTracker.GetRotationSchedule()),
	}
}

// PKIMiddleware returns middleware that integrates PKI observability with audit logging.
func (p *PKIObservability) PKIMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !p.IsEnabled() {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		ctx := r.Context()

		lrw := newLoggingResponseWriter(w)

		next.ServeHTTP(lrw, r)

		// Log PKI-relevant events
		if isPKIPath(r.URL.Path) {
			event := AuditEvent{
				ID:        "pki-" + time.Now().UTC().Format("20060102150405.000000"),
				Timestamp: time.Now().UTC(),
				Action:    r.Method,
				Resource:  r.URL.Path,
				Method:    r.Method,
				Path:      r.URL.Path,
				IP:        getClientIP(r),
				UserAgent: r.UserAgent(),
				Status:    lrw.statusCode,
				Metadata: map[string]interface{}{
					"duration_ms": time.Since(start).Milliseconds(),
					"pki_path":    true,
				},
			}

			if userID := ctx.Value("user_id"); userID != nil {
				event.UserID = userID.(string)
			}
			if sessionID := ctx.Value("session_id"); sessionID != nil {
				event.SessionID = sessionID.(string)
			}

			go func() {
				if err := p.auditLogger.Log(ctx, event); err != nil {
					log.Printf("pki audit log error: %v", err)
				}
			}()
		}
	})
}

// defineDefaultBoundaries sets up default trust boundaries for the 3-phase rollout.
func (p *PKIObservability) defineDefaultBoundaries() {
	ctx := context.Background()

	// Phase 1: Identity and Ingress TLS boundary
	p.trustBoundary.DefineBoundary(ctx, "ingress-tls", []string{
		"api-gateway", "auth-service", "public-api",
	}, true, 1)

	// Phase 2: Selective mTLS boundary for critical services
	p.trustBoundary.DefineBoundary(ctx, "critical-services", []string{
		"auth-service", "vault", "database", "cache",
	}, true, 2)

	// Phase 3: Strict east-west mTLS boundary for all internal services
	p.trustBoundary.DefineBoundary(ctx, "east-west-mtls", []string{
		"api-gateway", "auth-service", "vault", "database", "cache",
		"messaging", "workflow", "sync", "projection-search",
		"projection-graph", "projection-vector",
	}, true, 3)

	// Non-mTLS boundary for external/internal demarcation
	p.trustBoundary.DefineBoundary(ctx, "external-edge", []string{
		"api-gateway", "cdn", "load-balancer",
	}, false, 1)
}

// isPKIPath checks if a path is PKI-relevant.
func isPKIPath(path string) bool {
	pkiPaths := []string{
		"/api/security/pki/",
		"/api/security/bootstrap",
		"/api/auth/",
		"/api/vault/",
	}
	for _, p := range pkiPaths {
		if len(path) >= len(p) && path[:len(p)] == p {
			return true
		}
	}
	return false
}

// defaultPKIAuditLogger is a simple audit logger for PKI events.
type defaultPKIAuditLogger struct{}

func (d *defaultPKIAuditLogger) Log(ctx context.Context, event AuditEvent) error {
	data, _ := json.Marshal(event)
	log.Printf("PKI AUDIT: %s", string(data))
	return nil
}

// auditLog logs a PKI audit event.
func (p *PKIObservability) auditLog(ctx context.Context, action string, metadata map[string]interface{}) {
	event := AuditEvent{
		ID:        "pki-" + time.Now().UTC().Format("20060102150405.000000"),
		Timestamp: time.Now().UTC(),
		Action:    action,
		Resource:  "pki-observability",
		Method:    "INTERNAL",
		Path:      "/internal/pki",
		Metadata:  metadata,
	}
	go func() {
		if err := p.auditLogger.Log(ctx, event); err != nil {
			log.Printf("pki audit log error: %v", err)
		}
	}()
}

// loggingResponseWriter wraps http.ResponseWriter to capture status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}
	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}