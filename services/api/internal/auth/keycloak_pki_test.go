package keycloak

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"commlayers/services/api/internal/security"
)

func generateTestCertificate(t *testing.T, notAfter time.Time) *x509.Certificate {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(12345),
		Subject:      pkix.Name{CommonName: "test.example.com"},
		Issuer:       pkix.Name{CommonName: "Test CA"},
		NotBefore:    time.Now().Add(-24 * time.Hour),
		NotAfter:     notAfter,
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("failed to create certificate: %v", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatalf("failed to parse certificate: %v", err)
	}
	return cert
}

func TestInitKeycloakWithPKI(t *testing.T) {
	// Reset global instances
	enabled = false
	security.PKIOnce = sync.Once{}
	security.PKIInstance = nil

	config := &PKIConfig{
		ExpiryAlertDays:  15,
		InitialMTLSPhase: 2,
	}
	InitKeycloakWithPKI(config)

	if !IsEnabled() {
		t.Error("expected Keycloak to be enabled")
	}

	pki := GetPKIObservability()
	if pki == nil {
		t.Fatal("expected PKI observability instance")
	}
	if !pki.IsEnabled() {
		t.Error("expected PKI observability to be enabled")
	}
	if pki.MTLSStatusTracker().GetPhase() != 2 {
		t.Errorf("expected MTLS phase 2, got %d", pki.MTLSStatusTracker().GetPhase())
	}
	if pki.CertificateTracker().ExpiryAlert() != 15*24*time.Hour {
		t.Errorf("expected expiry alert 15 days, got %v", pki.CertificateTracker().ExpiryAlert())
	}
}

func TestInitKeycloakWithPKINilConfig(t *testing.T) {
	enabled = false
	security.PKIOnce = sync.Once{}
	security.PKIInstance = nil

	InitKeycloakWithPKI(nil)

	pki := GetPKIObservability()
	if pki == nil {
		t.Fatal("expected PKI observability instance")
	}
	if pki.MTLSStatusTracker().GetPhase() != 1 {
		t.Errorf("expected default MTLS phase 1, got %d", pki.MTLSStatusTracker().GetPhase())
	}
	if pki.CertificateTracker().ExpiryAlert() != 30*24*time.Hour {
		t.Errorf("expected default expiry alert 30 days, got %v", pki.CertificateTracker().ExpiryAlert())
	}
}

func TestGetPKIObservability(t *testing.T) {
	enabled = false
	security.PKIOnce = sync.Once{}
	security.PKIInstance = nil

	pki := GetPKIObservability()
	if pki != nil {
		t.Error("expected nil before init")
	}

	InitKeycloakWithPKI(nil)
	pki = GetPKIObservability()
	if pki == nil {
		t.Fatal("expected PKI observability after init")
	}
}

func TestPKIMiddleware(t *testing.T) {
	enabled = false
	security.PKIOnce = sync.Once{}
	security.PKIInstance = nil

	// Test 1: Both PKI and Keycloak disabled
	handler := PKIMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Test 1: expected 200 when both disabled, got %d", w.Code)
	}

	// Test 2: PKI enabled, Keycloak disabled
	enabled = false
	security.PKIOnce = sync.Once{}
	security.PKIInstance = nil
	InitKeycloakWithPKI(&PKIConfig{ExpiryAlertDays: 30, InitialMTLSPhase: 1})
	// Now disable Keycloak but keep PKI enabled
	enabled = false

	handler = PKIMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req = httptest.NewRequest("GET", "/test", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Test 2: expected 200 when PKI enabled but Keycloak disabled, got %d", w.Code)
	}

	// Test 3: Both PKI and Keycloak enabled
	enabled = false
	security.PKIOnce = sync.Once{}
	security.PKIInstance = nil
	InitKeycloakWithPKI(nil)
	// Keycloak is now enabled (InitKeycloakWithPKI calls InitKeycloak)

	handler = PKIMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// No token - should get 401
	req = httptest.NewRequest("GET", "/test", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Test 3a: expected 401 when no token, got %d", w.Code)
	}

	// With token - should get 200
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer token")
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Test 3b: expected 200 with token, got %d", w.Code)
	}
}

func TestRecordClientCertificate(t *testing.T) {
	enabled = false
	security.PKIOnce = sync.Once{}
	security.PKIInstance = nil

	InitKeycloakWithPKI(nil)

	cert := generateTestCertificate(t, time.Now().Add(60*24*time.Hour))

	// Create request with TLS connection state
	req := httptest.NewRequest("GET", "/test", nil)
	req.TLS = &tls.ConnectionState{
		PeerCertificates: []*x509.Certificate{cert},
	}

	err := RecordClientCertificate(req, "test-service")
	if err != nil {
		t.Fatalf("RecordClientCertificate failed: %v", err)
	}

	pki := GetPKIObservability()
	certs := pki.GetCertificates()
	if len(certs) != 1 {
		t.Errorf("expected 1 certificate recorded, got %d", len(certs))
	}
	if certs[0].Service != "test-service" {
		t.Errorf("expected service test-service, got %s", certs[0].Service)
	}
}

func TestRecordClientCertificateNoTLS(t *testing.T) {
	enabled = false
	security.PKIOnce = sync.Once{}
	security.PKIInstance = nil

	InitKeycloakWithPKI(nil)

	req := httptest.NewRequest("GET", "/test", nil)
	// No TLS state

	err := RecordClientCertificate(req, "test-service")
	if err != nil {
		t.Errorf("expected no error when no TLS, got %v", err)
	}
}

func TestSetMTLSPhase(t *testing.T) {
	enabled = false
	security.PKIOnce = sync.Once{}
	security.PKIInstance = nil

	InitKeycloakWithPKI(nil)

	SetMTLSPhase(3)

	pki := GetPKIObservability()
	if pki.MTLSStatusTracker().GetPhase() != 3 {
		t.Errorf("expected phase 3, got %d", pki.MTLSStatusTracker().GetPhase())
	}

	// Test without PKI initialized
	enabled = false
	security.PKIOnce = sync.Once{}
	security.PKIInstance = nil

	SetMTLSPhase(2) // Should not panic
}