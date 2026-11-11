package security

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestPKIObservabilityInit(t *testing.T) {
	// Reset global instance for test
	PKIOnce = sync.Once{}
	PKIInstance = nil

	config := DefaultPKIConfig()
	InitPKIObservability(config)

	pki := GetPKIObservability()
	if pki == nil {
		t.Fatal("expected PKI observability instance")
	}
	if !pki.IsEnabled() {
		t.Error("expected PKI observability to be enabled")
	}

	// Second init should be no-op
	InitPKIObservability(config)
	pki2 := GetPKIObservability()
	if pki != pki2 {
		t.Error("expected same instance on second init")
	}
}

func TestPKIObservabilityIsEnabledWhenNotInit(t *testing.T) {
	PKIOnce = sync.Once{}
	PKIInstance = nil

	pki := GetPKIObservability()
	if pki != nil {
		t.Error("expected nil when not initialized")
	}
}

func TestSetMTLSPhase(t *testing.T) {
	PKIOnce = sync.Once{}
	PKIInstance = nil

	config := DefaultPKIConfig()
	InitPKIObservability(config)

	pki := GetPKIObservability()
	pki.SetMTLSPhase(3)

	if pki.MTLSStatusTracker().GetPhase() != 3 {
		t.Errorf("expected phase 3, got %d", pki.MTLSStatusTracker().GetPhase())
	}
}

func TestRecordCertificate(t *testing.T) {
	PKIOnce = sync.Once{}
	PKIInstance = nil

	config := DefaultPKIConfig()
	InitPKIObservability(config)

	pki := GetPKIObservability()
	ctx := context.Background()

	cert := generateTestCertificate(t, time.Now().Add(60*24*time.Hour))
	pemData := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})

	info, err := pki.RecordCertificate(ctx, pemData, "test-source", "test-service")
	if err != nil {
		t.Fatalf("RecordCertificate failed: %v", err)
	}

	if info == nil {
		t.Fatal("expected certificate info")
	}
	if info.Source != "test-source" {
		t.Errorf("expected source test-source, got %s", info.Source)
	}
	if info.Service != "test-service" {
		t.Errorf("expected service test-service, got %s", info.Service)
	}
}

func TestRecordMTLSConnection(t *testing.T) {
	PKIOnce = sync.Once{}
	PKIInstance = nil

	config := DefaultPKIConfig()
	InitPKIObservability(config)

	pki := GetPKIObservability()
	ctx := context.Background()

	state := &tls.ConnectionState{
		Version:     tls.VersionTLS13,
		CipherSuite: tls.TLS_AES_256_GCM_SHA384,
		VerifiedChains: [][]*x509.Certificate{
			{
				{
					NotAfter: time.Now().Add(60 * 24 * time.Hour),
				},
			},
		},
		PeerCertificates: []*x509.Certificate{
			{
				NotAfter: time.Now().Add(60 * 24 * time.Hour),
			},
		},
	}

	pki.RecordMTLSConnection(ctx, "source-svc", "target-svc", state)

	conn, ok := pki.MTLSStatusTracker().GetConnection("source-svc", "target-svc")
	if !ok {
		t.Fatal("expected to find connection")
	}
	if conn.TLSVersion != "TLS 1.3" {
		t.Errorf("expected TLS 1.3, got %s", conn.TLSVersion)
	}
	if !conn.Verified {
		t.Error("expected verified true")
	}
}

func TestEvaluateTrustBoundaries(t *testing.T) {
	PKIOnce = sync.Once{}
	PKIInstance = nil

	config := DefaultPKIConfig()
	InitPKIObservability(config)

	pki := GetPKIObservability()
	ctx := context.Background()

	// Record some mTLS connections
	pki.MTLSStatusTracker().RecordConnectionManual("api-gateway", "auth-service", "TLS 1.3", "CIPHER", true, 60, "")
	pki.MTLSStatusTracker().RecordConnectionManual("auth-service", "api-gateway", "TLS 1.3", "CIPHER", true, 60, "")

	pki.EvaluateTrustBoundaries(ctx)

	boundaries := pki.GetTrustBoundaries()
	if len(boundaries) == 0 {
		t.Fatal("expected trust boundaries")
	}

	// Check ingress-tls boundary
	var ingressBoundary *TrustBoundary
	for _, b := range boundaries {
		if b.Name == "ingress-tls" {
			ingressBoundary = b
			break
		}
	}
	if ingressBoundary == nil {
		t.Fatal("expected ingress-tls boundary")
	}
	// Should have violations since not all service pairs have connections
}

func TestGetCertificates(t *testing.T) {
	PKIOnce = sync.Once{}
	PKIInstance = nil

	config := DefaultPKIConfig()
	InitPKIObservability(config)

	pki := GetPKIObservability()
	ctx := context.Background()

	cert := generateTestCertificate(t, time.Now().Add(60*24*time.Hour))
	pemData := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
	pki.RecordCertificate(ctx, pemData, "test", "svc")

	certs := pki.GetCertificates()
	if len(certs) != 1 {
		t.Errorf("expected 1 certificate, got %d", len(certs))
	}
}

func TestGetMTLSStatus(t *testing.T) {
	PKIOnce = sync.Once{}
	PKIInstance = nil

	config := DefaultPKIConfig()
	InitPKIObservability(config)

	pki := GetPKIObservability()
	pki.MTLSStatusTracker().RecordConnectionManual("a", "b", "TLS 1.3", "CIPHER", true, 30, "")

	status := pki.GetMTLSStatus()
	if len(status) != 1 {
		t.Errorf("expected 1 connection, got %d", len(status))
	}
}

func TestGetTrustBoundaries(t *testing.T) {
	PKIOnce = sync.Once{}
	PKIInstance = nil

	config := DefaultPKIConfig()
	InitPKIObservability(config)

	pki := GetPKIObservability()
	boundaries := pki.GetTrustBoundaries()

	// Should have default boundaries
	if len(boundaries) < 4 {
		t.Errorf("expected at least 4 default boundaries, got %d", len(boundaries))
	}

	// Check for expected default boundaries
	names := make(map[string]bool)
	for _, b := range boundaries {
		names[b.Name] = true
	}
	expected := []string{"ingress-tls", "critical-services", "east-west-mtls", "external-edge"}
	for _, exp := range expected {
		if !names[exp] {
			t.Errorf("expected boundary %s not found", exp)
		}
	}
}

func TestGetRotationSchedule(t *testing.T) {
	PKIOnce = sync.Once{}
	PKIInstance = nil

	config := DefaultPKIConfig()
	InitPKIObservability(config)

	pki := GetPKIObservability()
	ctx := context.Background()

	// Add expiring cert
	cert := generateTestCertificate(t, time.Now().Add(15*24*time.Hour))
	pemData := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
	pki.RecordCertificate(ctx, pemData, "test", "svc")

	// Need to refresh statuses to evaluate rotation due
	pki.CertificateTracker().RefreshStatuses()

	schedule := pki.GetRotationSchedule()
	if len(schedule) != 1 {
		t.Errorf("expected 1 cert in rotation schedule, got %d", len(schedule))
	}
}

func TestGetPKISummary(t *testing.T) {
	PKIOnce = sync.Once{}
	PKIInstance = nil

	config := DefaultPKIConfig()
	InitPKIObservability(config)

	pki := GetPKIObservability()
	ctx := context.Background()

	// Add some test data
	cert := generateTestCertificate(t, time.Now().Add(60*24*time.Hour))
	pemData := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
	pki.RecordCertificate(ctx, pemData, "test", "svc")

	pki.MTLSStatusTracker().RecordConnectionManual("a", "b", "TLS 1.3", "CIPHER", true, 30, "")

	summary := pki.GetPKISummary()

	certs := summary["certificates"].(map[string]interface{})
	if certs["total"].(int) != 1 {
		t.Errorf("expected 1 total cert, got %v", certs["total"])
	}

	conns := summary["mtls_connections"].(map[string]interface{})
	if conns["total"].(int) != 1 {
		t.Errorf("expected 1 total connection, got %v", conns["total"])
	}

	boundaries := summary["trust_boundaries"].(map[string]interface{})
	if boundaries["total"].(int) < 4 {
		t.Errorf("expected at least 4 boundaries, got %v", boundaries["total"])
	}
}

func TestPKIMiddleware(t *testing.T) {
	PKIOnce = sync.Once{}
	PKIInstance = nil

	config := DefaultPKIConfig()
	InitPKIObservability(config)

	pki := GetPKIObservability()

	handler := pki.PKIMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Test non-PKI path
	req := httptest.NewRequest("GET", "/api/catalog", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	// Test PKI path
	req = httptest.NewRequest("GET", "/api/security/pki/certificates", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for PKI path, got %d", w.Code)
	}
}

func TestIsPKIPath(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/api/security/pki/certificates", true},
		{"/api/security/pki/mtls-status", true},
		{"/api/security/bootstrap", true},
		{"/api/auth/login", true},
		{"/api/vault/secrets", true},
		{"/api/catalog", false},
		{"/healthz", false},
	}

	for _, tt := range tests {
		result := isPKIPath(tt.path)
		if result != tt.expected {
			t.Errorf("isPKIPath(%q) = %v, want %v", tt.path, result, tt.expected)
		}
	}
}

func TestDefaultPKIAuditLogger(t *testing.T) {
	logger := &defaultPKIAuditLogger{}
	event := AuditEvent{
		ID:       "test-1",
		Timestamp: time.Now().UTC(),
		Action:   "test",
		Resource: "test",
	}

	err := logger.Log(context.Background(), event)
	if err != nil {
		t.Errorf("DefaultAuditLogger.Log failed: %v", err)
	}
}