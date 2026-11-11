package evidence

import (
	"encoding/json"
	"time"
)

// EvidenceSource represents the source of evidence.
type EvidenceSource string

const (
	SourceBenchmark EvidenceSource = "benchmark"
	SourceAudit     EvidenceSource = "audit"
	SourceProjection EvidenceSource = "projection"
)

// EvidenceStaleness represents the staleness classification of evidence.
type EvidenceStaleness string

const (
	StalenessFresh       EvidenceStaleness = "fresh"
	StalenessStale       EvidenceStaleness = "stale"
	StalenessVeryStale   EvidenceStaleness = "very_stale"
	StalenessReindexing  EvidenceStaleness = "reindexing"
)

// SuccessCode represents the fixed success taxonomy.
type SuccessCode string

const (
	SuccessOK        SuccessCode = "ok"
	SuccessPartial   SuccessCode = "partial"
	SuccessDegraded  SuccessCode = "degraded"
)

// ErrorCode represents the fixed error taxonomy.
type ErrorCode string

const (
	ErrorTimeout       ErrorCode = "timeout"
	ErrorUnavailable   ErrorCode = "unavailable"
	ErrorInvalidInput  ErrorCode = "invalid_input"
	ErrorAuthFailed    ErrorCode = "auth_failed"
	ErrorRateLimited   ErrorCode = "rate_limited"
	ErrorInternal      ErrorCode = "internal_error"
	ErrorNotFound      ErrorCode = "not_found"
	ErrorConflict      ErrorCode = "conflict"
	ErrorStaleData     ErrorCode = "stale_data"
)

// Evidence represents the canonical evidence structure used across all streams.
type Evidence struct {
	Source          EvidenceSource    `json:"source"`
	Confidence      float64           `json:"confidence"`
	Staleness       EvidenceStaleness `json:"staleness"`
	MeasuredAt      time.Time         `json:"measured_at"`
	DelayAppliedMs  int64             `json:"delay_applied_ms"`
}

// EvidenceContract represents a formal evidence contract with success/error taxonomy.
type EvidenceContract struct {
	Evidence     Evidence   `json:"evidence"`
	Success      bool       `json:"success"`
	SuccessCode  SuccessCode `json:"success_code,omitempty"`
	ErrorCode    ErrorCode   `json:"error_code,omitempty"`
	ErrorMessage string      `json:"error_message,omitempty"`
}

// BenchmarkEvidence extends Evidence with benchmark-specific fields.
type BenchmarkEvidence struct {
	Evidence
	RunID           string  `json:"run_id"`
	QueueDelayMs    float64 `json:"queue_delay_ms,omitempty"`
	TimeoutMs       float64 `json:"timeout_ms,omitempty"`
	TransitionMs    float64 `json:"transition_ms,omitempty"`
	SampleCount     int     `json:"sample_count"`
}

// NewBenchmarkEvidence creates a new BenchmarkEvidence with defaults.
func NewBenchmarkEvidence(runID string, delayMs, timeoutMs, transitionMs float64, sampleCount int) *BenchmarkEvidence {
	return &BenchmarkEvidence{
		Evidence: Evidence{
			Source:         SourceBenchmark,
			Confidence:     1.0,
			Staleness:      StalenessFresh,
			MeasuredAt:     time.Now(),
			DelayAppliedMs: 0,
		},
		RunID:        runID,
		QueueDelayMs: delayMs,
		TimeoutMs:    timeoutMs,
		TransitionMs: transitionMs,
		SampleCount:  sampleCount,
	}
}

// NewSuccessContract creates a successful evidence contract.
func NewSuccessContract(evidence Evidence, code SuccessCode) *EvidenceContract {
	return &EvidenceContract{
		Evidence:    evidence,
		Success:     true,
		SuccessCode: code,
	}
}

// NewErrorContract creates an error evidence contract.
func NewErrorContract(evidence Evidence, code ErrorCode, message string) *EvidenceContract {
	return &EvidenceContract{
		Evidence:     evidence,
		Success:      false,
		ErrorCode:    code,
		ErrorMessage: message,
	}
}

// MarshalJSON implements custom JSON marshaling for EvidenceContract.
func (ec *EvidenceContract) MarshalJSON() ([]byte, error) {
	type Alias EvidenceContract
	return json.Marshal(&struct {
		*Alias
		SuccessCode string `json:"success_code,omitempty"`
		ErrorCode   string `json:"error_code,omitempty"`
	}{
		Alias:       (*Alias)(ec),
		SuccessCode: string(ec.SuccessCode),
		ErrorCode:   string(ec.ErrorCode),
	})
}

// ValidSuccessCodes returns all valid success codes.
func ValidSuccessCodes() []SuccessCode {
	return []SuccessCode{SuccessOK, SuccessPartial, SuccessDegraded}
}

// ValidErrorCodes returns all valid error codes.
func ValidErrorCodes() []ErrorCode {
	return []ErrorCode{
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
}

// ValidStalenessValues returns all valid staleness values.
func ValidStalenessValues() []EvidenceStaleness {
	return []EvidenceStaleness{
		StalenessFresh,
		StalenessStale,
		StalenessVeryStale,
		StalenessReindexing,
	}
}

// ValidSources returns all valid evidence sources.
func ValidSources() []EvidenceSource {
	return []EvidenceSource{SourceBenchmark, SourceAudit, SourceProjection}
}