package security

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"
)

func TestCertificateTrackerCreation(t *testing.T) {
	tracker := NewCertificateTracker(30 * 24 * time.Hour)
	if tracker == nil {
		t.Fatal("expected non-nil tracker")
	}
	if tracker.expiryAlert != 30*24*time.Hour {
		t.Errorf("expected expiry alert of 30 days, got %v", tracker.expiryAlert)
	}
}

func TestCertificateTrackerDefaultExpiryAlert(t *testing.T) {
	tracker := NewCertificateTracker(0)
	if tracker.expiryAlert != 30*24*time.Hour {
		t.Errorf("expected default expiry alert of 30 days, got %v", tracker.expiryAlert)
	}
}

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

func TestTrackCertificate(t *testing.T) {
	tracker := NewCertificateTracker(30 * 24 * time.Hour)
	ctx := context.Background()

	cert := generateTestCertificate(t, time.Now().Add(60*24*time.Hour))
	info := tracker.TrackCertificate(ctx, cert, "test", "test-service")

	if info == nil {
		t.Fatal("expected certificate info")
	}
	if info.SerialNumber != "12345" {
		t.Errorf("expected serial number 12345, got %s", info.SerialNumber)
	}
	if info.Subject != "CN=test.example.com" {
		t.Errorf("expected subject CN=test.example.com, got %s", info.Subject)
	}
	// Self-signed cert has issuer = subject
	if info.Issuer != "CN=test.example.com" {
		t.Errorf("expected issuer CN=test.example.com, got %s", info.Issuer)
	}
	if info.Source != "test" {
		t.Errorf("expected source test, got %s", info.Source)
	}
	if info.Service != "test-service" {
		t.Errorf("expected service test-service, got %s", info.Service)
	}
	if info.Status != "valid" {
		t.Errorf("expected status valid, got %s", info.Status)
	}
	if info.RotationDue {
		t.Error("expected rotation due false for future expiry")
	}
}

func TestTrackCertificatePEM(t *testing.T) {
	tracker := NewCertificateTracker(30 * 24 * time.Hour)
	ctx := context.Background()

	cert := generateTestCertificate(t, time.Now().Add(60*24*time.Hour))
	pemData := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})

	info, err := tracker.TrackCertificatePEM(ctx, pemData, "pem-test", "pem-service")
	if err != nil {
		t.Fatalf("TrackCertificatePEM failed: %v", err)
	}

	if info == nil {
		t.Fatal("expected certificate info")
	}
	if info.Source != "pem-test" {
		t.Errorf("expected source pem-test, got %s", info.Source)
	}
	if info.Service != "pem-service" {
		t.Errorf("expected service pem-service, got %s", info.Service)
	}
}

func TestTrackCertificatePEMInvalid(t *testing.T) {
	tracker := NewCertificateTracker(30 * 24 * time.Hour)
	ctx := context.Background()

	_, err := tracker.TrackCertificatePEM(ctx, []byte("invalid"), "test", "test")
	if err == nil {
		t.Error("expected error for invalid PEM")
	}
}

func TestGetCertificate(t *testing.T) {
	tracker := NewCertificateTracker(30 * 24 * time.Hour)
	ctx := context.Background()

	cert := generateTestCertificate(t, time.Now().Add(60*24*time.Hour))
	info := tracker.TrackCertificate(ctx, cert, "test", "test-service")

	retrieved, ok := tracker.GetCertificate(info.SerialNumber)
	if !ok {
		t.Fatal("expected to find certificate")
	}
	if retrieved.SerialNumber != info.SerialNumber {
		t.Errorf("expected serial %s, got %s", info.SerialNumber, retrieved.SerialNumber)
	}

	_, ok = tracker.GetCertificate("nonexistent")
	if ok {
		t.Error("expected not to find nonexistent certificate")
	}
}

func TestGetAllCertificates(t *testing.T) {
	tracker := NewCertificateTracker(30 * 24 * time.Hour)
	ctx := context.Background()

	cert1 := generateTestCertificate(t, time.Now().Add(60*24*time.Hour))
	tracker.TrackCertificate(ctx, cert1, "test", "service-1")

	cert2 := generateTestCertificate(t, time.Now().Add(90*24*time.Hour))
	cert2.SerialNumber = big.NewInt(67890)
	tracker.TrackCertificate(ctx, cert2, "test", "service-2")

	all := tracker.GetAllCertificates()
	if len(all) != 2 {
		t.Errorf("expected 2 certificates, got %d", len(all))
	}
}

func TestGetCertificatesByService(t *testing.T) {
	tracker := NewCertificateTracker(30 * 24 * time.Hour)
	ctx := context.Background()

	cert1 := generateTestCertificate(t, time.Now().Add(60*24*time.Hour))
	tracker.TrackCertificate(ctx, cert1, "test", "service-a")

	cert2 := generateTestCertificate(t, time.Now().Add(90*24*time.Hour))
	cert2.SerialNumber = big.NewInt(67890)
	tracker.TrackCertificate(ctx, cert2, "test", "service-b")

	serviceA := tracker.GetCertificatesByService("service-a")
	if len(serviceA) != 1 {
		t.Errorf("expected 1 certificate for service-a, got %d", len(serviceA))
	}
	if serviceA[0].Service != "service-a" {
		t.Errorf("expected service-a, got %s", serviceA[0].Service)
	}
}

func TestExpiringCertificates(t *testing.T) {
	tracker := NewCertificateTracker(30 * 24 * time.Hour)
	ctx := context.Background()

	// Certificate expiring in 15 days (within alert window)
	cert1 := generateTestCertificate(t, time.Now().Add(15*24*time.Hour))
	tracker.TrackCertificate(ctx, cert1, "test", "service-1")

	// Certificate expiring in 60 days (outside alert window)
	cert2 := generateTestCertificate(t, time.Now().Add(60*24*time.Hour))
	cert2.SerialNumber = big.NewInt(67890)
	tracker.TrackCertificate(ctx, cert2, "test", "service-2")

	expiring := tracker.GetExpiringCertificates()
	if len(expiring) != 1 {
		t.Errorf("expected 1 expiring certificate, got %d", len(expiring))
	}
	if expiring[0].Service != "service-1" {
		t.Errorf("expected service-1, got %s", expiring[0].Service)
	}
}

func TestExpiredCertificates(t *testing.T) {
	tracker := NewCertificateTracker(30 * 24 * time.Hour)
	ctx := context.Background()

	// Expired certificate
	cert1 := generateTestCertificate(t, time.Now().Add(-24*time.Hour))
	tracker.TrackCertificate(ctx, cert1, "test", "service-1")

	// Valid certificate
	cert2 := generateTestCertificate(t, time.Now().Add(60*24*time.Hour))
	cert2.SerialNumber = big.NewInt(67890)
	tracker.TrackCertificate(ctx, cert2, "test", "service-2")

	// Need to refresh statuses to evaluate expiry
	tracker.RefreshStatuses()

	expired := tracker.GetExpiredCertificates()
	if len(expired) != 1 {
		t.Errorf("expected 1 expired certificate, got %d", len(expired))
	}
	if expired[0].Status != "expired" {
		t.Errorf("expected status expired, got %s", expired[0].Status)
	}
}

func TestRevokeCertificate(t *testing.T) {
	tracker := NewCertificateTracker(30 * 24 * time.Hour)
	ctx := context.Background()

	cert := generateTestCertificate(t, time.Now().Add(60*24*time.Hour))
	info := tracker.TrackCertificate(ctx, cert, "test", "test-service")

	ok := tracker.RevokeCertificate(info.SerialNumber, "key compromise")
	if !ok {
		t.Fatal("expected revoke to succeed")
	}

	revoked := tracker.GetRevokedCertificates()
	if len(revoked) != 1 {
		t.Errorf("expected 1 revoked certificate, got %d", len(revoked))
	}
	if revoked[0].RevocationReason != "key compromise" {
		t.Errorf("expected revocation reason 'key compromise', got %s", revoked[0].RevocationReason)
	}
	if revoked[0].Status != "revoked" {
		t.Errorf("expected status revoked, got %s", revoked[0].Status)
	}

	// Try to revoke non-existent
	ok = tracker.RevokeCertificate("nonexistent", "test")
	if ok {
		t.Error("expected revoke to fail for non-existent")
	}
}

func TestRefreshStatuses(t *testing.T) {
	tracker := NewCertificateTracker(30 * 24 * time.Hour)
	ctx := context.Background()

	// Valid cert
	cert1 := generateTestCertificate(t, time.Now().Add(60*24*time.Hour))
	tracker.TrackCertificate(ctx, cert1, "test", "service-1")

	// Expiring cert (15 days)
	cert2 := generateTestCertificate(t, time.Now().Add(15*24*time.Hour))
	cert2.SerialNumber = big.NewInt(67890)
	tracker.TrackCertificate(ctx, cert2, "test", "service-2")

	// Expired cert
	cert3 := generateTestCertificate(t, time.Now().Add(-24*time.Hour))
	cert3.SerialNumber = big.NewInt(99999)
	tracker.TrackCertificate(ctx, cert3, "test", "service-3")

	tracker.RefreshStatuses()

	all := tracker.GetAllCertificates()
	for _, info := range all {
		switch info.Service {
		case "service-1":
			if info.Status != "valid" {
				t.Errorf("service-1: expected valid, got %s", info.Status)
			}
			if info.RotationDue {
				t.Error("service-1: expected rotation due false")
			}
		case "service-2":
			if info.Status != "expiring" {
				t.Errorf("service-2: expected expiring, got %s", info.Status)
			}
			if !info.RotationDue {
				t.Error("service-2: expected rotation due true")
			}
		case "service-3":
			if info.Status != "expired" {
				t.Errorf("service-3: expected expired, got %s", info.Status)
			}
			if !info.RotationDue {
				t.Error("service-3: expected rotation due true")
			}
		}
	}
}

func TestRotationSchedule(t *testing.T) {
	tracker := NewCertificateTracker(30 * 24 * time.Hour)
	ctx := context.Background()

	// Valid cert (no rotation due)
	cert1 := generateTestCertificate(t, time.Now().Add(60*24*time.Hour))
	tracker.TrackCertificate(ctx, cert1, "test", "service-1")

	// Expiring cert
	cert2 := generateTestCertificate(t, time.Now().Add(15*24*time.Hour))
	cert2.SerialNumber = big.NewInt(67890)
	tracker.TrackCertificate(ctx, cert2, "test", "service-2")

	// Expired cert
	cert3 := generateTestCertificate(t, time.Now().Add(-24*time.Hour))
	cert3.SerialNumber = big.NewInt(99999)
	tracker.TrackCertificate(ctx, cert3, "test", "service-3")

	// Revoked cert
	cert4 := generateTestCertificate(t, time.Now().Add(90*24*time.Hour))
	cert4.SerialNumber = big.NewInt(11111)
	info4 := tracker.TrackCertificate(ctx, cert4, "test", "service-4")
	tracker.RevokeCertificate(info4.SerialNumber, "revoked")

	tracker.RefreshStatuses()

	schedule := tracker.GetRotationSchedule()
	// Should have 2: expiring and expired (revoked excluded)
	if len(schedule) != 2 {
		t.Errorf("expected 2 certificates in rotation schedule, got %d", len(schedule))
	}

	// Check sorted by expiry (soonest first)
	if schedule[0].NotAfter.After(schedule[1].NotAfter) {
		t.Error("expected rotation schedule sorted by expiry ascending")
	}
}

func TestRemoveCertificate(t *testing.T) {
	tracker := NewCertificateTracker(30 * 24 * time.Hour)
	ctx := context.Background()

	cert := generateTestCertificate(t, time.Now().Add(60*24*time.Hour))
	info := tracker.TrackCertificate(ctx, cert, "test", "test-service")

	ok := tracker.RemoveCertificate(info.SerialNumber)
	if !ok {
		t.Fatal("expected remove to succeed")
	}

	_, ok = tracker.GetCertificate(info.SerialNumber)
	if ok {
		t.Error("expected certificate to be removed")
	}

	// Remove non-existent
	ok = tracker.RemoveCertificate("nonexistent")
	if ok {
		t.Error("expected remove to fail for non-existent")
	}
}

func TestCertificateError(t *testing.T) {
	err := ErrInvalidCertificate
	if err.Error() != "certificate error: invalid certificate data" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}