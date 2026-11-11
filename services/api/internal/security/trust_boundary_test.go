package security

import (
	"context"
	"testing"
	"time"
)

func TestTrustBoundaryTrackerCreation(t *testing.T) {
	mtlsTracker := NewMTLSStatusTracker(1)
	certTracker := NewCertificateTracker(30 * 24 * time.Hour)
	tracker := NewTrustBoundaryTracker(mtlsTracker, certTracker)

	if tracker == nil {
		t.Fatal("expected non-nil tracker")
	}
}

func TestDefineBoundary(t *testing.T) {
	mtlsTracker := NewMTLSStatusTracker(1)
	certTracker := NewCertificateTracker(30 * 24 * time.Hour)
	tracker := NewTrustBoundaryTracker(mtlsTracker, certTracker)
	ctx := context.Background()

	boundary := tracker.DefineBoundary(ctx, "test-boundary", []string{"svc-a", "svc-b"}, true, 2)

	if boundary == nil {
		t.Fatal("expected boundary")
	}
	if boundary.Name != "test-boundary" {
		t.Errorf("expected name test-boundary, got %s", boundary.Name)
	}
	if len(boundary.Services) != 2 {
		t.Errorf("expected 2 services, got %d", len(boundary.Services))
	}
	if !boundary.MTLSRequired {
		t.Error("expected MTLSRequired true")
	}
	if boundary.Phase != 2 {
		t.Errorf("expected phase 2, got %d", boundary.Phase)
	}
	if boundary.Compliant {
		t.Error("expected compliant false initially")
	}
	if len(boundary.Violations) != 0 {
		t.Errorf("expected empty violations, got %v", boundary.Violations)
	}
}

func TestDefineBoundaryPhaseClamping(t *testing.T) {
	mtlsTracker := NewMTLSStatusTracker(1)
	certTracker := NewCertificateTracker(30 * 24 * time.Hour)
	tracker := NewTrustBoundaryTracker(mtlsTracker, certTracker)
	ctx := context.Background()

	boundary := tracker.DefineBoundary(ctx, "test", []string{"a"}, true, 5)
	if boundary.Phase != 3 {
		t.Errorf("expected phase clamped to 3, got %d", boundary.Phase)
	}

	boundary = tracker.DefineBoundary(ctx, "test2", []string{"a"}, true, 0)
	if boundary.Phase != 1 {
		t.Errorf("expected phase clamped to 1, got %d", boundary.Phase)
	}
}

func TestGetBoundary(t *testing.T) {
	mtlsTracker := NewMTLSStatusTracker(1)
	certTracker := NewCertificateTracker(30 * 24 * time.Hour)
	tracker := NewTrustBoundaryTracker(mtlsTracker, certTracker)
	ctx := context.Background()

	tracker.DefineBoundary(ctx, "boundary-1", []string{"a", "b"}, true, 1)

	boundary, ok := tracker.GetBoundary("boundary-1")
	if !ok {
		t.Fatal("expected to find boundary")
	}
	if boundary.Name != "boundary-1" {
		t.Errorf("expected boundary-1, got %s", boundary.Name)
	}

	_, ok = tracker.GetBoundary("nonexistent")
	if ok {
		t.Error("expected not to find nonexistent")
	}
}

func TestGetAllBoundaries(t *testing.T) {
	mtlsTracker := NewMTLSStatusTracker(1)
	certTracker := NewCertificateTracker(30 * 24 * time.Hour)
	tracker := NewTrustBoundaryTracker(mtlsTracker, certTracker)
	ctx := context.Background()

	tracker.DefineBoundary(ctx, "b1", []string{"a"}, true, 1)
	tracker.DefineBoundary(ctx, "b2", []string{"b"}, true, 2)

	all := tracker.GetAllBoundaries()
	if len(all) != 2 {
		t.Errorf("expected 2 boundaries, got %d", len(all))
	}
}

func TestEvaluateComplianceMTLSNotRequired(t *testing.T) {
	mtlsTracker := NewMTLSStatusTracker(1)
	certTracker := NewCertificateTracker(30 * 24 * time.Hour)
	tracker := NewTrustBoundaryTracker(mtlsTracker, certTracker)
	ctx := context.Background()

	tracker.DefineBoundary(ctx, "no-mtls", []string{"a", "b"}, false, 1)
	tracker.EvaluateCompliance(ctx)

	boundary, _ := tracker.GetBoundary("no-mtls")
	if !boundary.Compliant {
		t.Error("expected compliant true when MTLS not required")
	}
	if boundary.TotalPairs != 0 {
		t.Errorf("expected 0 total pairs, got %d", boundary.TotalPairs)
	}
	if boundary.CompliantPairs != 0 {
		t.Errorf("expected 0 compliant pairs, got %d", boundary.CompliantPairs)
	}
}

func TestEvaluateComplianceNoConnections(t *testing.T) {
	mtlsTracker := NewMTLSStatusTracker(1)
	certTracker := NewCertificateTracker(30 * 24 * time.Hour)
	tracker := NewTrustBoundaryTracker(mtlsTracker, certTracker)
	ctx := context.Background()

	tracker.DefineBoundary(ctx, "empty", []string{"a", "b"}, true, 1)
	tracker.EvaluateCompliance(ctx)

	boundary, _ := tracker.GetBoundary("empty")
	if boundary.Compliant {
		t.Error("expected non-compliant when no connections observed")
	}
	if boundary.TotalPairs != 2 {
		t.Errorf("expected 2 total pairs (a->b, b->a), got %d", boundary.TotalPairs)
	}
	if boundary.CompliantPairs != 0 {
		t.Errorf("expected 0 compliant pairs, got %d", boundary.CompliantPairs)
	}
	if len(boundary.Violations) != 2 {
		t.Errorf("expected 2 violations, got %d: %v", len(boundary.Violations), boundary.Violations)
	}
}

func TestEvaluateComplianceVerifiedConnections(t *testing.T) {
	mtlsTracker := NewMTLSStatusTracker(1)
	certTracker := NewCertificateTracker(30 * 24 * time.Hour)
	tracker := NewTrustBoundaryTracker(mtlsTracker, certTracker)
	ctx := context.Background()

	tracker.DefineBoundary(ctx, "verified", []string{"a", "b"}, true, 1)

	// Record verified connections both ways
	mtlsTracker.RecordConnectionManual("a", "b", "TLS 1.3", "CIPHER", true, 60, "")
	mtlsTracker.RecordConnectionManual("b", "a", "TLS 1.3", "CIPHER", true, 60, "")

	tracker.EvaluateCompliance(ctx)

	boundary, _ := tracker.GetBoundary("verified")
	if !boundary.Compliant {
		t.Errorf("expected compliant, got violations: %v", boundary.Violations)
	}
	if boundary.TotalPairs != 2 {
		t.Errorf("expected 2 total pairs, got %d", boundary.TotalPairs)
	}
	if boundary.CompliantPairs != 2 {
		t.Errorf("expected 2 compliant pairs, got %d", boundary.CompliantPairs)
	}
}

func TestEvaluateComplianceUnverifiedConnection(t *testing.T) {
	mtlsTracker := NewMTLSStatusTracker(1)
	certTracker := NewCertificateTracker(30 * 24 * time.Hour)
	tracker := NewTrustBoundaryTracker(mtlsTracker, certTracker)
	ctx := context.Background()

	tracker.DefineBoundary(ctx, "unverified", []string{"a", "b"}, true, 1)

	// One direction unverified
	mtlsTracker.RecordConnectionManual("a", "b", "TLS 1.3", "CIPHER", true, 60, "")
	mtlsTracker.RecordConnectionManual("b", "a", "TLS 1.2", "CIPHER", false, -1, "certificate expired")

	tracker.EvaluateCompliance(ctx)

	boundary, _ := tracker.GetBoundary("unverified")
	if boundary.Compliant {
		t.Error("expected non-compliant with unverified connection")
	}
	if len(boundary.Violations) != 1 {
		t.Errorf("expected 1 violation, got %d: %v", len(boundary.Violations), boundary.Violations)
	}
}

func TestEvaluateComplianceExpiringCertificate(t *testing.T) {
	mtlsTracker := NewMTLSStatusTracker(1)
	certTracker := NewCertificateTracker(30 * 24 * time.Hour)
	tracker := NewTrustBoundaryTracker(mtlsTracker, certTracker)
	ctx := context.Background()

	tracker.DefineBoundary(ctx, "expiring", []string{"a", "b"}, true, 1)

	// Connection with cert expiring in 15 days
	mtlsTracker.RecordConnectionManual("a", "b", "TLS 1.3", "CIPHER", true, 15, "")
	mtlsTracker.RecordConnectionManual("b", "a", "TLS 1.3", "CIPHER", true, 60, "")

	tracker.EvaluateCompliance(ctx)

	boundary, _ := tracker.GetBoundary("expiring")
	if boundary.Compliant {
		t.Error("expected non-compliant with expiring certificate")
	}
	found := false
	for _, v := range boundary.Violations {
		if contains(v, "expiring soon") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected violation about expiring soon, got: %v", boundary.Violations)
	}
}

func TestEvaluateCompliancePhaseMismatch(t *testing.T) {
	mtlsTracker := NewMTLSStatusTracker(1)
	certTracker := NewCertificateTracker(30 * 24 * time.Hour)
	tracker := NewTrustBoundaryTracker(mtlsTracker, certTracker)
	ctx := context.Background()

	// Boundary requires phase 2
	tracker.DefineBoundary(ctx, "phase-mismatch", []string{"a", "b"}, true, 2)

	// Connection is phase 1
	mtlsTracker.RecordConnectionManual("a", "b", "TLS 1.3", "CIPHER", true, 60, "")
	mtlsTracker.RecordConnectionManual("b", "a", "TLS 1.3", "CIPHER", true, 60, "")

	tracker.EvaluateCompliance(ctx)

	boundary, _ := tracker.GetBoundary("phase-mismatch")
	if boundary.Compliant {
		t.Error("expected non-compliant with phase mismatch")
	}
	found := false
	for _, v := range boundary.Violations {
		if contains(v, "phase mismatch") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected violation about phase mismatch, got: %v", boundary.Violations)
	}
}

func TestGetCompliantBoundaries(t *testing.T) {
	mtlsTracker := NewMTLSStatusTracker(1)
	certTracker := NewCertificateTracker(30 * 24 * time.Hour)
	tracker := NewTrustBoundaryTracker(mtlsTracker, certTracker)
	ctx := context.Background()

	tracker.DefineBoundary(ctx, "compliant", []string{"a", "b"}, true, 1)
	tracker.DefineBoundary(ctx, "non-compliant", []string{"c", "d"}, true, 1)

	mtlsTracker.RecordConnectionManual("a", "b", "TLS 1.3", "CIPHER", true, 60, "")
	mtlsTracker.RecordConnectionManual("b", "a", "TLS 1.3", "CIPHER", true, 60, "")

	tracker.EvaluateCompliance(ctx)

	compliant := tracker.GetCompliantBoundaries()
	if len(compliant) != 1 {
		t.Errorf("expected 1 compliant boundary, got %d", len(compliant))
	}
	if compliant[0].Name != "compliant" {
		t.Errorf("expected compliant, got %s", compliant[0].Name)
	}

	nonCompliant := tracker.GetNonCompliantBoundaries()
	if len(nonCompliant) != 1 {
		t.Errorf("expected 1 non-compliant boundary, got %d", len(nonCompliant))
	}
	if nonCompliant[0].Name != "non-compliant" {
		t.Errorf("expected non-compliant, got %s", nonCompliant[0].Name)
	}
}

func TestGetBoundaryByService(t *testing.T) {
	mtlsTracker := NewMTLSStatusTracker(1)
	certTracker := NewCertificateTracker(30 * 24 * time.Hour)
	tracker := NewTrustBoundaryTracker(mtlsTracker, certTracker)
	ctx := context.Background()

	tracker.DefineBoundary(ctx, "b1", []string{"a", "b"}, true, 1)
	tracker.DefineBoundary(ctx, "b2", []string{"b", "c"}, true, 2)
	tracker.DefineBoundary(ctx, "b3", []string{"d"}, true, 3)

	boundaries := tracker.GetBoundaryByService("b")
	if len(boundaries) != 2 {
		t.Errorf("expected 2 boundaries for service b, got %d", len(boundaries))
	}

	boundaries = tracker.GetBoundaryByService("d")
	if len(boundaries) != 1 {
		t.Errorf("expected 1 boundary for service d, got %d", len(boundaries))
	}
	if boundaries[0].Name != "b3" {
		t.Errorf("expected b3, got %s", boundaries[0].Name)
	}
}

func TestUpdateBoundaryPhase(t *testing.T) {
	mtlsTracker := NewMTLSStatusTracker(1)
	certTracker := NewCertificateTracker(30 * 24 * time.Hour)
	tracker := NewTrustBoundaryTracker(mtlsTracker, certTracker)
	ctx := context.Background()

	tracker.DefineBoundary(ctx, "test", []string{"a"}, true, 1)

	ok := tracker.UpdateBoundaryPhase("test", 3)
	if !ok {
		t.Fatal("expected update to succeed")
	}

	boundary, _ := tracker.GetBoundary("test")
	if boundary.Phase != 3 {
		t.Errorf("expected phase 3, got %d", boundary.Phase)
	}

	ok = tracker.UpdateBoundaryPhase("nonexistent", 2)
	if ok {
		t.Error("expected update to fail for nonexistent")
	}
}

func TestRemoveBoundary(t *testing.T) {
	mtlsTracker := NewMTLSStatusTracker(1)
	certTracker := NewCertificateTracker(30 * 24 * time.Hour)
	tracker := NewTrustBoundaryTracker(mtlsTracker, certTracker)
	ctx := context.Background()

	tracker.DefineBoundary(ctx, "test", []string{"a"}, true, 1)

	ok := tracker.RemoveBoundary("test")
	if !ok {
		t.Fatal("expected remove to succeed")
	}

	_, ok = tracker.GetBoundary("test")
	if ok {
		t.Error("expected boundary to be removed")
	}

	ok = tracker.RemoveBoundary("nonexistent")
	if ok {
		t.Error("expected remove to fail for nonexistent")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}