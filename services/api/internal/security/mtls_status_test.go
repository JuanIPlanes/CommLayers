package security

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"testing"
	"time"
)

func TestMTLSStatusTrackerCreation(t *testing.T) {
	tracker := NewMTLSStatusTracker(1)
	if tracker == nil {
		t.Fatal("expected non-nil tracker")
	}
	if tracker.GetPhase() != 1 {
		t.Errorf("expected phase 1, got %d", tracker.GetPhase())
	}
}

func TestMTLSStatusTrackerPhaseClamping(t *testing.T) {
	tracker := NewMTLSStatusTracker(5)
	if tracker.GetPhase() != 3 {
		t.Errorf("expected phase clamped to 3, got %d", tracker.GetPhase())
	}

	tracker = NewMTLSStatusTracker(0)
	if tracker.GetPhase() != 1 {
		t.Errorf("expected phase clamped to 1, got %d", tracker.GetPhase())
	}
}

func TestSetPhase(t *testing.T) {
	tracker := NewMTLSStatusTracker(1)
	tracker.SetPhase(2)
	if tracker.GetPhase() != 2 {
		t.Errorf("expected phase 2, got %d", tracker.GetPhase())
	}

	// Add a connection first
	tracker.RecordConnectionManual("svc-a", "svc-b", "TLS 1.3", "TLS_AES_256_GCM_SHA384", true, 60, "")
	conn, _ := tracker.GetConnection("svc-a", "svc-b")
	if conn.Phase != 2 {
		t.Errorf("expected connection phase updated to 2, got %d", conn.Phase)
	}
}

func TestRecordConnectionManual(t *testing.T) {
	tracker := NewMTLSStatusTracker(1)

	tracker.RecordConnectionManual("source-svc", "target-svc", "TLS 1.3", "TLS_AES_256_GCM_SHA384", true, 60, "")

	conn, ok := tracker.GetConnection("source-svc", "target-svc")
	if !ok {
		t.Fatal("expected to find connection")
	}
	if conn.SourceService != "source-svc" {
		t.Errorf("expected source source-svc, got %s", conn.SourceService)
	}
	if conn.TargetService != "target-svc" {
		t.Errorf("expected target target-svc, got %s", conn.TargetService)
	}
	if conn.TLSVersion != "TLS 1.3" {
		t.Errorf("expected TLS 1.3, got %s", conn.TLSVersion)
	}
	if conn.CipherSuite != "TLS_AES_256_GCM_SHA384" {
		t.Errorf("expected TLS_AES_256_GCM_SHA384, got %s", conn.CipherSuite)
	}
	if !conn.Verified {
		t.Error("expected verified true")
	}
	if conn.CertExpiryDays != 60 {
		t.Errorf("expected cert expiry 60 days, got %d", conn.CertExpiryDays)
	}
	if conn.Phase != 1 {
		t.Errorf("expected phase 1, got %d", conn.Phase)
	}
	if conn.Error != "" {
		t.Errorf("expected empty error, got %s", conn.Error)
	}
}

func TestRecordConnectionManualUnverified(t *testing.T) {
	tracker := NewMTLSStatusTracker(1)

	tracker.RecordConnectionManual("source", "target", "TLS 1.2", "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384", false, -1, "certificate expired")

	conn, ok := tracker.GetConnection("source", "target")
	if !ok {
		t.Fatal("expected to find connection")
	}
	if conn.Verified {
		t.Error("expected verified false")
	}
	if conn.Error != "certificate expired" {
		t.Errorf("expected error 'certificate expired', got %s", conn.Error)
	}
}

func TestGetConnection(t *testing.T) {
	tracker := NewMTLSStatusTracker(1)
	tracker.RecordConnectionManual("a", "b", "TLS 1.3", "CIPHER", true, 30, "")

	conn, ok := tracker.GetConnection("a", "b")
	if !ok {
		t.Fatal("expected to find connection")
	}
	if conn.SourceService != "a" {
		t.Errorf("expected source a, got %s", conn.SourceService)
	}

	_, ok = tracker.GetConnection("b", "a")
	if ok {
		t.Error("expected not to find reverse connection")
	}

	_, ok = tracker.GetConnection("nonexistent", "target")
	if ok {
		t.Error("expected not to find nonexistent")
	}
}

func TestGetAllConnections(t *testing.T) {
	tracker := NewMTLSStatusTracker(1)
	tracker.RecordConnectionManual("a", "b", "TLS 1.3", "CIPHER", true, 30, "")
	tracker.RecordConnectionManual("b", "c", "TLS 1.2", "CIPHER", true, 60, "")

	all := tracker.GetAllConnections()
	if len(all) != 2 {
		t.Errorf("expected 2 connections, got %d", len(all))
	}
}

func TestGetConnectionsBySource(t *testing.T) {
	tracker := NewMTLSStatusTracker(1)
	tracker.RecordConnectionManual("svc-a", "svc-b", "TLS 1.3", "CIPHER", true, 30, "")
	tracker.RecordConnectionManual("svc-a", "svc-c", "TLS 1.2", "CIPHER", true, 60, "")
	tracker.RecordConnectionManual("svc-b", "svc-c", "TLS 1.3", "CIPHER", true, 90, "")

	fromA := tracker.GetConnectionsBySource("svc-a")
	if len(fromA) != 2 {
		t.Errorf("expected 2 connections from svc-a, got %d", len(fromA))
	}
	for _, c := range fromA {
		if c.SourceService != "svc-a" {
			t.Errorf("expected all from svc-a, got %s", c.SourceService)
		}
	}
}

func TestGetConnectionsByTarget(t *testing.T) {
	tracker := NewMTLSStatusTracker(1)
	tracker.RecordConnectionManual("svc-a", "svc-c", "TLS 1.3", "CIPHER", true, 30, "")
	tracker.RecordConnectionManual("svc-b", "svc-c", "TLS 1.2", "CIPHER", true, 60, "")

	toC := tracker.GetConnectionsByTarget("svc-c")
	if len(toC) != 2 {
		t.Errorf("expected 2 connections to svc-c, got %d", len(toC))
	}
	for _, c := range toC {
		if c.TargetService != "svc-c" {
			t.Errorf("expected all to svc-c, got %s", c.TargetService)
		}
	}
}

func TestGetUnverifiedConnections(t *testing.T) {
	tracker := NewMTLSStatusTracker(1)
	tracker.RecordConnectionManual("a", "b", "TLS 1.3", "CIPHER", true, 30, "")
	tracker.RecordConnectionManual("c", "d", "TLS 1.2", "CIPHER", false, -1, "expired")
	tracker.RecordConnectionManual("e", "f", "TLS 1.3", "CIPHER", true, 60, "")

	unverified := tracker.GetUnverifiedConnections()
	if len(unverified) != 1 {
		t.Errorf("expected 1 unverified connection, got %d", len(unverified))
	}
	if unverified[0].SourceService != "c" {
		t.Errorf("expected c->d, got %s->%s", unverified[0].SourceService, unverified[0].TargetService)
	}
}

func TestGetExpiringConnections(t *testing.T) {
	tracker := NewMTLSStatusTracker(1)
	tracker.RecordConnectionManual("a", "b", "TLS 1.3", "CIPHER", true, 15, "")
	tracker.RecordConnectionManual("c", "d", "TLS 1.2", "CIPHER", true, 45, "")
	tracker.RecordConnectionManual("e", "f", "TLS 1.3", "CIPHER", true, 5, "")

	expiring := tracker.GetExpiringConnections(30)
	if len(expiring) != 2 {
		t.Errorf("expected 2 expiring connections within 30 days, got %d", len(expiring))
	}
	for _, c := range expiring {
		if c.CertExpiryDays > 30 {
			t.Errorf("expected expiry <= 30, got %d", c.CertExpiryDays)
		}
	}
}

func TestGetConnectionMatrix(t *testing.T) {
	tracker := NewMTLSStatusTracker(1)
	tracker.RecordConnectionManual("a", "b", "TLS 1.3", "CIPHER", true, 30, "")
	tracker.RecordConnectionManual("a", "c", "TLS 1.2", "CIPHER", true, 60, "")
	tracker.RecordConnectionManual("b", "c", "TLS 1.3", "CIPHER", true, 90, "")

	matrix := tracker.GetConnectionMatrix()
	if len(matrix) != 2 {
		t.Errorf("expected 2 source services in matrix, got %d", len(matrix))
	}

	if len(matrix["a"]) != 2 {
		t.Errorf("expected 2 targets for a, got %d", len(matrix["a"]))
	}
	if len(matrix["b"]) != 1 {
		t.Errorf("expected 1 target for b, got %d", len(matrix["b"]))
	}
	if matrix["a"]["b"].TargetService != "b" {
		t.Errorf("expected a->b, got %s", matrix["a"]["b"].TargetService)
	}
}

func TestRemoveConnection(t *testing.T) {
	tracker := NewMTLSStatusTracker(1)
	tracker.RecordConnectionManual("a", "b", "TLS 1.3", "CIPHER", true, 30, "")

	ok := tracker.RemoveConnection("a", "b")
	if !ok {
		t.Fatal("expected remove to succeed")
	}

	_, ok = tracker.GetConnection("a", "b")
	if ok {
		t.Error("expected connection to be removed")
	}

	ok = tracker.RemoveConnection("nonexistent", "target")
	if ok {
		t.Error("expected remove to fail for non-existent")
	}
}

func TestTLSVersionString(t *testing.T) {
	tests := []struct {
		version    uint16
		expected   string
	}{
		{tls.VersionTLS10, "TLS 1.0"},
		{tls.VersionTLS11, "TLS 1.1"},
		{tls.VersionTLS12, "TLS 1.2"},
		{tls.VersionTLS13, "TLS 1.3"},
		{999, "Unknown"},
	}

	for _, tt := range tests {
		result := tlsVersionString(tt.version)
		if result != tt.expected {
			t.Errorf("tlsVersionString(%d) = %q, want %q", tt.version, result, tt.expected)
		}
	}
}

func TestRecordConnectionWithTLSState(t *testing.T) {
	tracker := NewMTLSStatusTracker(1)

	// Create a minimal TLS ConnectionState
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

	tracker.RecordConnection(context.Background(), "source", "target", state)

	conn, ok := tracker.GetConnection("source", "target")
	if !ok {
		t.Fatal("expected to find connection")
	}
	if conn.TLSVersion != "TLS 1.3" {
		t.Errorf("expected TLS 1.3, got %s", conn.TLSVersion)
	}
	if conn.CipherSuite != "TLS_AES_256_GCM_SHA384" {
		t.Errorf("expected TLS_AES_256_GCM_SHA384, got %s", conn.CipherSuite)
	}
	if !conn.Verified {
		t.Error("expected verified true")
	}
	if conn.CertExpiryDays < 50 || conn.CertExpiryDays > 70 {
		t.Errorf("expected cert expiry around 60 days, got %d", conn.CertExpiryDays)
	}
}