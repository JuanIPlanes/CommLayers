package evidence

import (
	"encoding/json"
	"testing"
	"time"
)

func TestEvidenceContract_Success(t *testing.T) {
	ev := Evidence{
		Source:         SourceBenchmark,
		Confidence:     0.95,
		Staleness:      StalenessFresh,
		MeasuredAt:     time.Now(),
		DelayAppliedMs: 10,
	}

	ec := NewSuccessContract(ev, SuccessOK)

	if !ec.Success {
		t.Error("expected success=true")
	}
	if ec.SuccessCode != SuccessOK {
		t.Errorf("expected success_code=ok, got %s", ec.SuccessCode)
	}
	if ec.ErrorCode != "" {
		t.Errorf("expected empty error_code, got %s", ec.ErrorCode)
	}
	if ec.ErrorMessage != "" {
		t.Errorf("expected empty error_message, got %s", ec.ErrorMessage)
	}
}

func TestEvidenceContract_Error(t *testing.T) {
	ev := Evidence{
		Source:         SourceAudit,
		Confidence:     0.8,
		Staleness:      StalenessStale,
		MeasuredAt:     time.Now(),
		DelayAppliedMs: 0,
	}

	ec := NewErrorContract(ev, ErrorTimeout, "request timed out")

	if ec.Success {
		t.Error("expected success=false")
	}
	if ec.ErrorCode != ErrorTimeout {
		t.Errorf("expected error_code=timeout, got %s", ec.ErrorCode)
	}
	if ec.ErrorMessage != "request timed out" {
		t.Errorf("expected error_message='request timed out', got %s", ec.ErrorMessage)
	}
	if ec.SuccessCode != "" {
		t.Errorf("expected empty success_code, got %s", ec.SuccessCode)
	}
}

func TestBenchmarkEvidence(t *testing.T) {
	be := NewBenchmarkEvidence("run_123", 10.5, 5000.0, 2.3, 42)

	if be.RunID != "run_123" {
		t.Errorf("expected run_id=run_123, got %s", be.RunID)
	}
	if be.QueueDelayMs != 10.5 {
		t.Errorf("expected queue_delay_ms=10.5, got %f", be.QueueDelayMs)
	}
	if be.TimeoutMs != 5000.0 {
		t.Errorf("expected timeout_ms=5000.0, got %f", be.TimeoutMs)
	}
	if be.TransitionMs != 2.3 {
		t.Errorf("expected transition_ms=2.3, got %f", be.TransitionMs)
	}
	if be.SampleCount != 42 {
		t.Errorf("expected sample_count=42, got %d", be.SampleCount)
	}
	if be.Source != SourceBenchmark {
		t.Errorf("expected source=benchmark, got %s", be.Source)
	}
	if be.Confidence != 1.0 {
		t.Errorf("expected confidence=1.0, got %f", be.Confidence)
	}
	if be.Staleness != StalenessFresh {
		t.Errorf("expected staleness=fresh, got %s", be.Staleness)
	}
}

func TestValidSuccessCodes(t *testing.T) {
	codes := ValidSuccessCodes()
	expected := []SuccessCode{SuccessOK, SuccessPartial, SuccessDegraded}
	if len(codes) != len(expected) {
		t.Errorf("expected %d codes, got %d", len(expected), len(codes))
	}
	for i, c := range codes {
		if c != expected[i] {
			t.Errorf("expected %s at index %d, got %s", expected[i], i, c)
		}
	}
}

func TestValidErrorCodes(t *testing.T) {
	codes := ValidErrorCodes()
	expected := []ErrorCode{
		ErrorTimeout,
		ErrorUnavailable,
		ErrorInvalidInput,
		ErrorAuthFailed,
		ErrorRateLimited,
		ErrorInternal,
		ErrorNotFound,
		ErrorConflict,
		ErrorStaleData,
	}
	if len(codes) != len(expected) {
		t.Errorf("expected %d codes, got %d", len(expected), len(codes))
	}
	for i, c := range codes {
		if c != expected[i] {
			t.Errorf("expected %s at index %d, got %s", expected[i], i, c)
		}
	}
}

func TestValidStalenessValues(t *testing.T) {
	values := ValidStalenessValues()
	expected := []EvidenceStaleness{StalenessFresh, StalenessStale, StalenessVeryStale, StalenessReindexing}
	if len(values) != len(expected) {
		t.Errorf("expected %d values, got %d", len(expected), len(values))
	}
	for i, v := range values {
		if v != expected[i] {
			t.Errorf("expected %s at index %d, got %s", expected[i], i, v)
		}
	}
}

func TestValidSources(t *testing.T) {
	sources := ValidSources()
	expected := []EvidenceSource{SourceBenchmark, SourceAudit, SourceProjection}
	if len(sources) != len(expected) {
		t.Errorf("expected %d sources, got %d", len(expected), len(sources))
	}
	for i, s := range sources {
		if s != expected[i] {
			t.Errorf("expected %s at index %d, got %s", expected[i], i, s)
		}
	}
}

func TestEvidenceContract_MarshalJSON(t *testing.T) {
	ev := Evidence{
		Source:         SourceBenchmark,
		Confidence:     0.9,
		Staleness:      StalenessFresh,
		MeasuredAt:     time.Now(),
		DelayAppliedMs: 5,
	}

	ec := NewSuccessContract(ev, SuccessPartial)
	data, err := ec.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if parsed["success"] != true {
		t.Error("expected success=true in JSON")
	}
	if parsed["success_code"] != "partial" {
		t.Errorf("expected success_code=partial, got %v", parsed["success_code"])
	}
	if parsed["error_code"] != nil && parsed["error_code"] != "" {
		t.Errorf("expected empty error_code, got %v", parsed["error_code"])
	}
}