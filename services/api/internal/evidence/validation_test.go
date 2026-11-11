package evidence

import (
	"errors"
	"testing"
	"time"
)

func TestValidateEvidence_Valid(t *testing.T) {
	ev := Evidence{
		Source:         SourceBenchmark,
		Confidence:     0.95,
		Staleness:      StalenessFresh,
		MeasuredAt:     time.Now(),
		DelayAppliedMs: 10,
	}

	if err := ValidateEvidence(ev); err != nil {
		t.Errorf("expected valid evidence, got error: %v", err)
	}
}

func TestValidateEvidence_InvalidConfidence(t *testing.T) {
	tests := []struct {
		name       string
		confidence float64
	}{
		{"negative", -0.1},
		{"over_one", 1.1},
		{"way_over", 2.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev := Evidence{
				Source:         SourceBenchmark,
				Confidence:     tt.confidence,
				Staleness:      StalenessFresh,
				MeasuredAt:     time.Now(),
				DelayAppliedMs: 0,
			}
			err := ValidateEvidence(ev)
			if err == nil {
				t.Error("expected error for invalid confidence")
			}
			if !errors.Is(err, ErrInvalidConfidence) {
				t.Errorf("expected ErrInvalidConfidence, got %v", err)
			}
		})
	}
}

func TestValidateEvidence_InvalidStaleness(t *testing.T) {
	ev := Evidence{
		Source:         SourceBenchmark,
		Confidence:     0.5,
		Staleness:      "invalid_staleness",
		MeasuredAt:     time.Now(),
		DelayAppliedMs: 0,
	}

	err := ValidateEvidence(ev)
	if err == nil {
		t.Error("expected error for invalid staleness")
	}
	if !errors.Is(err, ErrInvalidStaleness) {
		t.Errorf("expected ErrInvalidStaleness, got %v", err)
	}
}

func TestValidateEvidence_InvalidSource(t *testing.T) {
	ev := Evidence{
		Source:         "invalid_source",
		Confidence:     0.5,
		Staleness:      StalenessFresh,
		MeasuredAt:     time.Now(),
		DelayAppliedMs: 0,
	}

	err := ValidateEvidence(ev)
	if err == nil {
		t.Error("expected error for invalid source")
	}
	if !errors.Is(err, ErrInvalidSource) {
		t.Errorf("expected ErrInvalidSource, got %v", err)
	}
}

func TestValidateEvidence_FutureMeasuredAt(t *testing.T) {
	ev := Evidence{
		Source:         SourceBenchmark,
		Confidence:     0.5,
		Staleness:      StalenessFresh,
		MeasuredAt:     time.Now().Add(time.Hour),
		DelayAppliedMs: 0,
	}

	err := ValidateEvidence(ev)
	if err == nil {
		t.Error("expected error for future measured_at")
	}
	if !errors.Is(err, ErrFutureMeasuredAt) {
		t.Errorf("expected ErrFutureMeasuredAt, got %v", err)
	}
}

func TestValidateEvidence_NegativeDelay(t *testing.T) {
	ev := Evidence{
		Source:         SourceBenchmark,
		Confidence:     0.5,
		Staleness:      StalenessFresh,
		MeasuredAt:     time.Now(),
		DelayAppliedMs: -1,
	}

	err := ValidateEvidence(ev)
	if err == nil {
		t.Error("expected error for negative delay")
	}
	if !errors.Is(err, ErrNegativeDelay) {
		t.Errorf("expected ErrNegativeDelay, got %v", err)
	}
}

func TestValidateEvidenceContract_Success(t *testing.T) {
	ev := Evidence{
		Source:         SourceBenchmark,
		Confidence:     0.9,
		Staleness:      StalenessFresh,
		MeasuredAt:     time.Now(),
		DelayAppliedMs: 0,
	}
	ec := NewSuccessContract(ev, SuccessOK)

	if err := ValidateEvidenceContract(ec); err != nil {
		t.Errorf("expected valid contract, got error: %v", err)
	}
}

func TestValidateEvidenceContract_SuccessMissingCode(t *testing.T) {
	ev := Evidence{
		Source:         SourceBenchmark,
		Confidence:     0.9,
		Staleness:      StalenessFresh,
		MeasuredAt:     time.Now(),
		DelayAppliedMs: 0,
	}
	ec := &EvidenceContract{
		Evidence: ev,
		Success:  true,
		// Missing SuccessCode
	}

	err := ValidateEvidenceContract(ec)
	if err == nil {
		t.Error("expected error for missing success code")
	}
	if err != ErrMissingSuccessCode {
		t.Errorf("expected ErrMissingSuccessCode, got %v", err)
	}
}

func TestValidateEvidenceContract_SuccessWithErrorCode(t *testing.T) {
	ev := Evidence{
		Source:         SourceBenchmark,
		Confidence:     0.9,
		Staleness:      StalenessFresh,
		MeasuredAt:     time.Now(),
		DelayAppliedMs: 0,
	}
	ec := &EvidenceContract{
		Evidence:    ev,
		Success:     true,
		SuccessCode: SuccessOK,
		ErrorCode:   ErrorTimeout, // Should not be set for success
	}

	err := ValidateEvidenceContract(ec)
	if err == nil {
		t.Error("expected error for success with error code")
	}
	if !errors.Is(err, ErrContractMismatch) {
		t.Errorf("expected ErrContractMismatch, got %v", err)
	}
}

func TestValidateEvidenceContract_Error(t *testing.T) {
	ev := Evidence{
		Source:         SourceBenchmark,
		Confidence:     0.9,
		Staleness:      StalenessFresh,
		MeasuredAt:     time.Now(),
		DelayAppliedMs: 0,
	}
	ec := NewErrorContract(ev, ErrorTimeout, "request timed out")

	if err := ValidateEvidenceContract(ec); err != nil {
		t.Errorf("expected valid error contract, got error: %v", err)
	}
}

func TestValidateEvidenceContract_ErrorMissingCode(t *testing.T) {
	ev := Evidence{
		Source:         SourceBenchmark,
		Confidence:     0.9,
		Staleness:      StalenessFresh,
		MeasuredAt:     time.Now(),
		DelayAppliedMs: 0,
	}
	ec := &EvidenceContract{
		Evidence: ev,
		Success:  false,
		// Missing ErrorCode
	}

	err := ValidateEvidenceContract(ec)
	if err == nil {
		t.Error("expected error for missing error code")
	}
	if !errors.Is(err, ErrContractMismatch) {
		t.Errorf("expected ErrContractMismatch, got %v", err)
	}
}

func TestValidateEvidenceContract_ErrorMissingMessage(t *testing.T) {
	ev := Evidence{
		Source:         SourceBenchmark,
		Confidence:     0.9,
		Staleness:      StalenessFresh,
		MeasuredAt:     time.Now(),
		DelayAppliedMs: 0,
	}
	ec := &EvidenceContract{
		Evidence:  ev,
		Success:   false,
		ErrorCode: ErrorTimeout,
		// Missing ErrorMessage
	}

	err := ValidateEvidenceContract(ec)
	if err == nil {
		t.Error("expected error for missing error message")
	}
	if err != ErrMissingErrorMessage {
		t.Errorf("expected ErrMissingErrorMessage, got %v", err)
	}
}

func TestValidateEvidenceContract_ErrorWithSuccessCode(t *testing.T) {
	ev := Evidence{
		Source:         SourceBenchmark,
		Confidence:     0.9,
		Staleness:      StalenessFresh,
		MeasuredAt:     time.Now(),
		DelayAppliedMs: 0,
	}
	ec := &EvidenceContract{
		Evidence:     ev,
		Success:      false,
		ErrorCode:    ErrorTimeout,
		ErrorMessage: "timeout",
		SuccessCode:  SuccessOK, // Should not be set for error
	}

	err := ValidateEvidenceContract(ec)
	if err == nil {
		t.Error("expected error for error with success code")
	}
	if !errors.Is(err, ErrContractMismatch) {
		t.Errorf("expected ErrContractMismatch, got %v", err)
	}
}

func TestValidateBenchmarkEvidence_Valid(t *testing.T) {
	be := NewBenchmarkEvidence("run_123", 10.0, 100.0, 5.0, 10)

	if err := ValidateBenchmarkEvidence(be); err != nil {
		t.Errorf("expected valid benchmark evidence, got error: %v", err)
	}
}

func TestValidateBenchmarkEvidence_MissingRunID(t *testing.T) {
	be := NewBenchmarkEvidence("", 10.0, 100.0, 5.0, 10)

	err := ValidateBenchmarkEvidence(be)
	if err == nil {
		t.Error("expected error for missing run_id")
	}
}

func TestValidateBenchmarkEvidence_NegativeSampleCount(t *testing.T) {
	be := NewBenchmarkEvidence("run_123", 10.0, 100.0, 5.0, -1)

	err := ValidateBenchmarkEvidence(be)
	if err == nil {
		t.Error("expected error for negative sample_count")
	}
}

func TestIsValidConfidence(t *testing.T) {
	tests := []struct {
		value    float64
		expected bool
	}{
		{0.0, true},
		{0.5, true},
		{1.0, true},
		{-0.1, false},
		{1.1, false},
	}
	for _, tt := range tests {
		if IsValidConfidence(tt.value) != tt.expected {
			t.Errorf("IsValidConfidence(%f) = %v, want %v", tt.value, IsValidConfidence(tt.value), tt.expected)
		}
	}
}

func TestIsValidStaleness(t *testing.T) {
	if !IsValidStaleness(StalenessFresh) {
		t.Error("expected fresh to be valid")
	}
	if !IsValidStaleness(StalenessStale) {
		t.Error("expected stale to be valid")
	}
	if !IsValidStaleness(StalenessVeryStale) {
		t.Error("expected very_stale to be valid")
	}
	if !IsValidStaleness(StalenessReindexing) {
		t.Error("expected reindexing to be valid")
	}
	if IsValidStaleness("invalid") {
		t.Error("expected invalid to be invalid")
	}
}

func TestIsValidSource(t *testing.T) {
	if !IsValidSource(SourceBenchmark) {
		t.Error("expected benchmark to be valid")
	}
	if !IsValidSource(SourceAudit) {
		t.Error("expected audit to be valid")
	}
	if !IsValidSource(SourceProjection) {
		t.Error("expected projection to be valid")
	}
	if IsValidSource("invalid") {
		t.Error("expected invalid to be invalid")
	}
}

func TestIsValidSuccessCode(t *testing.T) {
	if !IsValidSuccessCode(SuccessOK) {
		t.Error("expected ok to be valid")
	}
	if !IsValidSuccessCode(SuccessPartial) {
		t.Error("expected partial to be valid")
	}
	if !IsValidSuccessCode(SuccessDegraded) {
		t.Error("expected degraded to be valid")
	}
	if IsValidSuccessCode("invalid") {
		t.Error("expected invalid to be invalid")
	}
}

func TestIsValidErrorCode(t *testing.T) {
	if !IsValidErrorCode(ErrorTimeout) {
		t.Error("expected timeout to be valid")
	}
	if !IsValidErrorCode(ErrorInternal) {
		t.Error("expected internal_error to be valid")
	}
	if IsValidErrorCode("invalid") {
		t.Error("expected invalid to be invalid")
	}
}

func TestNormalizeStaleness(t *testing.T) {
	tests := []struct {
		input    string
		expected EvidenceStaleness
	}{
		{"fresh", StalenessFresh},
		{"FRESH", StalenessFresh},
		{" stale ", StalenessStale},
		{"very_stale", StalenessVeryStale},
		{"very-stale", StalenessVeryStale},
		{"verystale", StalenessVeryStale},
		{"reindexing", StalenessReindexing},
		{"unknown", StalenessStale}, // default
		{"", StalenessStale},        // default
	}
	for _, tt := range tests {
		result := NormalizeStaleness(tt.input)
		if result != tt.expected {
			t.Errorf("NormalizeStaleness(%q) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}

func TestNormalizeSource(t *testing.T) {
	tests := []struct {
		input    string
		expected EvidenceSource
	}{
		{"benchmark", SourceBenchmark},
		{"BENCHMARK", SourceBenchmark},
		{" audit ", SourceAudit},
		{"projection", SourceProjection},
		{"unknown", SourceBenchmark}, // default
	}
	for _, tt := range tests {
		result := NormalizeSource(tt.input)
		if result != tt.expected {
			t.Errorf("NormalizeSource(%q) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}

func TestNormalizeSuccessCode(t *testing.T) {
	tests := []struct {
		input    string
		expected SuccessCode
	}{
		{"ok", SuccessOK},
		{"OK", SuccessOK},
		{"partial", SuccessPartial},
		{"degraded", SuccessDegraded},
		{"unknown", SuccessOK}, // default
	}
	for _, tt := range tests {
		result := NormalizeSuccessCode(tt.input)
		if result != tt.expected {
			t.Errorf("NormalizeSuccessCode(%q) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}

func TestNormalizeErrorCode(t *testing.T) {
	tests := []struct {
		input    string
		expected ErrorCode
	}{
		{"timeout", ErrorTimeout},
		{"unavailable", ErrorUnavailable},
		{"invalid_input", ErrorInvalidInput},
		{"invalid-input", ErrorInvalidInput},
		{"invalidinput", ErrorInvalidInput},
		{"auth_failed", ErrorAuthFailed},
		{"auth-failed", ErrorAuthFailed},
		{"authfailed", ErrorAuthFailed},
		{"rate_limited", ErrorRateLimited},
		{"rate-limited", ErrorRateLimited},
		{"ratelimited", ErrorRateLimited},
		{"internal_error", ErrorInternal},
		{"internal-error", ErrorInternal},
		{"internalerror", ErrorInternal},
		{"not_found", ErrorNotFound},
		{"not-found", ErrorNotFound},
		{"notfound", ErrorNotFound},
		{"conflict", ErrorConflict},
		{"stale_data", ErrorStaleData},
		{"stale-data", ErrorStaleData},
		{"staledata", ErrorStaleData},
		{"unknown", ErrorInternal}, // default
	}
	for _, tt := range tests {
		result := NormalizeErrorCode(tt.input)
		if result != tt.expected {
			t.Errorf("NormalizeErrorCode(%q) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}