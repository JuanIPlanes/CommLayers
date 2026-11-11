package evidence

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidConfidence   = errors.New("confidence must be between 0.0 and 1.0")
	ErrInvalidStaleness    = errors.New("invalid staleness value")
	ErrInvalidSource       = errors.New("invalid evidence source")
	ErrInvalidSuccessCode  = errors.New("invalid success code")
	ErrInvalidErrorCode    = errors.New("invalid error code")
	ErrMissingErrorMessage = errors.New("error message required when error code is set")
	ErrMissingSuccessCode  = errors.New("success code required when success is true")
	ErrContractMismatch    = errors.New("contract success flag mismatches error/success code presence")
	ErrFutureMeasuredAt    = errors.New("measured_at cannot be in the future")
	ErrNegativeDelay       = errors.New("delay_applied_ms cannot be negative")
)

// ValidateEvidence validates an Evidence struct.
func ValidateEvidence(e Evidence) error {
	if e.Confidence < 0.0 || e.Confidence > 1.0 {
		return fmt.Errorf("%w: got %f", ErrInvalidConfidence, e.Confidence)
	}

	validStaleness := false
	for _, v := range ValidStalenessValues() {
		if e.Staleness == v {
			validStaleness = true
			break
		}
	}
	if !validStaleness {
		return fmt.Errorf("%w: got %s, valid: %v", ErrInvalidStaleness, e.Staleness, ValidStalenessValues())
	}

	validSource := false
	for _, v := range ValidSources() {
		if e.Source == v {
			validSource = true
			break
		}
	}
	if !validSource {
		return fmt.Errorf("%w: got %s, valid: %v", ErrInvalidSource, e.Source, ValidSources())
	}

	if e.MeasuredAt.After(time.Now().Add(time.Second)) { // Allow 1s clock skew
		return fmt.Errorf("%w: got %v", ErrFutureMeasuredAt, e.MeasuredAt)
	}

	if e.DelayAppliedMs < 0 {
		return fmt.Errorf("%w: got %d", ErrNegativeDelay, e.DelayAppliedMs)
	}

	return nil
}

// ValidateEvidenceContract validates an EvidenceContract struct.
func ValidateEvidenceContract(ec *EvidenceContract) error {
	if err := ValidateEvidence(ec.Evidence); err != nil {
		return err
	}

	if ec.Success {
		if ec.SuccessCode == "" {
			return ErrMissingSuccessCode
		}
		validSuccess := false
		for _, v := range ValidSuccessCodes() {
			if ec.SuccessCode == v {
				validSuccess = true
				break
			}
		}
		if !validSuccess {
			return fmt.Errorf("%w: got %s, valid: %v", ErrInvalidSuccessCode, ec.SuccessCode, ValidSuccessCodes())
		}
		if ec.ErrorCode != "" {
			return fmt.Errorf("%w: success=true but error_code is set", ErrContractMismatch)
		}
		if ec.ErrorMessage != "" {
			return fmt.Errorf("%w: success=true but error_message is set", ErrContractMismatch)
		}
	} else {
		if ec.ErrorCode == "" {
			return fmt.Errorf("%w: success=false but error_code is empty", ErrContractMismatch)
		}
		validError := false
		for _, v := range ValidErrorCodes() {
			if ec.ErrorCode == v {
				validError = true
				break
			}
		}
		if !validError {
			return fmt.Errorf("%w: got %s, valid: %v", ErrInvalidErrorCode, ec.ErrorCode, ValidErrorCodes())
		}
		if ec.ErrorMessage == "" {
			return ErrMissingErrorMessage
		}
		if ec.SuccessCode != "" {
			return fmt.Errorf("%w: success=false but success_code is set", ErrContractMismatch)
		}
	}

	return nil
}

// ValidateBenchmarkEvidence validates a BenchmarkEvidence struct.
func ValidateBenchmarkEvidence(be *BenchmarkEvidence) error {
	if err := ValidateEvidence(be.Evidence); err != nil {
		return err
	}

	if be.RunID == "" {
		return errors.New("run_id is required")
	}

	if be.SampleCount < 0 {
		return errors.New("sample_count cannot be negative")
	}

	return nil
}

// IsValidConfidence checks if a confidence value is valid.
func IsValidConfidence(c float64) bool {
	return c >= 0.0 && c <= 1.0
}

// IsValidStaleness checks if a staleness value is valid.
func IsValidStaleness(s EvidenceStaleness) bool {
	for _, v := range ValidStalenessValues() {
		if s == v {
			return true
		}
	}
	return false
}

// IsValidSource checks if a source value is valid.
func IsValidSource(src EvidenceSource) bool {
	for _, v := range ValidSources() {
		if src == v {
			return true
		}
	}
	return false
}

// IsValidSuccessCode checks if a success code is valid.
func IsValidSuccessCode(sc SuccessCode) bool {
	for _, v := range ValidSuccessCodes() {
		if sc == v {
			return true
		}
	}
	return false
}

// IsValidErrorCode checks if an error code is valid.
func IsValidErrorCode(ec ErrorCode) bool {
	for _, v := range ValidErrorCodes() {
		if ec == v {
			return true
		}
	}
	return false
}

// NormalizeStaleness normalizes a staleness string to a valid EvidenceStaleness.
func NormalizeStaleness(s string) EvidenceStaleness {
	normalized := strings.ToLower(strings.TrimSpace(s))
	switch normalized {
	case "fresh":
		return StalenessFresh
	case "stale":
		return StalenessStale
	case "very_stale", "very-stale", "verystale":
		return StalenessVeryStale
	case "reindexing":
		return StalenessReindexing
	default:
		return StalenessStale // default to stale for unknown values
	}
}

// NormalizeSource normalizes a source string to a valid EvidenceSource.
func NormalizeSource(s string) EvidenceSource {
	normalized := strings.ToLower(strings.TrimSpace(s))
	switch normalized {
	case "benchmark":
		return SourceBenchmark
	case "audit":
		return SourceAudit
	case "projection":
		return SourceProjection
	default:
		return SourceBenchmark // default to benchmark
	}
}

// NormalizeSuccessCode normalizes a success code string to a valid SuccessCode.
func NormalizeSuccessCode(s string) SuccessCode {
	normalized := strings.ToLower(strings.TrimSpace(s))
	switch normalized {
	case "ok":
		return SuccessOK
	case "partial":
		return SuccessPartial
	case "degraded":
		return SuccessDegraded
	default:
		return SuccessOK // default to ok
	}
}

// NormalizeErrorCode normalizes an error code string to a valid ErrorCode.
func NormalizeErrorCode(s string) ErrorCode {
	normalized := strings.ToLower(strings.TrimSpace(s))
	switch normalized {
	case "timeout":
		return ErrorTimeout
	case "unavailable":
		return ErrorUnavailable
	case "invalid_input", "invalid-input", "invalidinput":
		return ErrorInvalidInput
	case "auth_failed", "auth-failed", "authfailed":
		return ErrorAuthFailed
	case "rate_limited", "rate-limited", "ratelimited":
		return ErrorRateLimited
	case "internal_error", "internal-error", "internalerror":
		return ErrorInternal
	case "not_found", "not-found", "notfound":
		return ErrorNotFound
	case "conflict":
		return ErrorConflict
	case "stale_data", "stale-data", "staledata":
		return ErrorStaleData
	default:
		return ErrorInternal // default to internal_error
	}
}