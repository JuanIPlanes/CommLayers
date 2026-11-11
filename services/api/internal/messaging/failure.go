package messaging

import (
	"time"
)

// FailureMode represents a category of failure in messaging/workflow systems.
type FailureMode string

const (
	// FailureModeQueueExhaustion indicates the queue is full or unavailable.
	FailureModeQueueExhaustion FailureMode = "queue_exhaustion"
	// FailureModeBrokerUnavailable indicates the message broker is unreachable.
	FailureModeBrokerUnavailable FailureMode = "broker_unavailable"
	// FailureModeMessageTimeout indicates a message processing timeout.
	FailureModeMessageTimeout FailureMode = "message_timeout"
	// FailureModeRetryExhausted indicates all retry attempts have been exhausted.
	FailureModeRetryExhausted FailureMode = "retry_exhausted"
	// FailureModeDeadLetter indicates the message was moved to a dead letter queue.
	FailureModeDeadLetter FailureMode = "dead_letter"
	// FailureModeSerializationFailure indicates message serialization/deserialization failed.
	FailureModeSerializationFailure FailureMode = "serialization_failure"
	// FailureModeIdempotencyViolation indicates duplicate detection or idempotency key conflict.
	FailureModeIdempotencyViolation FailureMode = "idempotency_violation"
	// FailureModeCircuitOpen indicates the circuit breaker is open.
	FailureModeCircuitOpen FailureMode = "circuit_open"
	// FailureModeOrderingViolation indicates message ordering guarantees were violated.
	FailureModeOrderingViolation FailureMode = "ordering_violation"
	// FailureModePayloadTooLarge indicates the message exceeds size limits.
	FailureModePayloadTooLarge FailureMode = "payload_too_large"
)

// FailureCheck represents a single failure mode check with its result.
type FailureCheck struct {
	// Mode is the failure mode being checked.
	Mode FailureMode `json:"mode"`
	// CheckedAt is when the check was performed.
	CheckedAt time.Time `json:"checked_at"`
	// Healthy indicates if the check passed (no failure detected).
	Healthy bool `json:"healthy"`
	// Details provides additional context about the check result.
	Details string `json:"details,omitempty"`
	// Severity indicates the severity if unhealthy (warning, critical).
	Severity string `json:"severity,omitempty"`
	// MetricValue is the numeric value associated with the check (e.g., queue depth, latency).
	MetricValue float64 `json:"metric_value,omitempty"`
	// Threshold is the threshold that triggered the failure (if unhealthy).
	Threshold float64 `json:"threshold,omitempty"`
}

// FailureModeMatrix defines the complete set of failure modes and their check configurations
// for a messaging/workflow family. This is used to validate demo and production behavior.
type FailureModeMatrix struct {
	// Family identifies the messaging family (e.g., "workflow_engine", "job_queue", "broker_event_bus").
	Family string `json:"family"`
	// Checks defines the failure mode checks for this family.
	Checks []FailureModeCheckConfig `json:"checks"`
	// UpdatedAt is when this matrix was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// FailureModeCheckConfig defines the configuration for a single failure mode check.
type FailureModeCheckConfig struct {
	// Mode is the failure mode to check.
	Mode FailureMode `json:"mode"`
	// Enabled indicates if this check is active.
	Enabled bool `json:"enabled"`
	// WarningThreshold is the threshold for a warning severity.
	WarningThreshold float64 `json:"warning_threshold,omitempty"`
	// CriticalThreshold is the threshold for a critical severity.
	CriticalThreshold float64 `json:"critical_threshold,omitempty"`
	// Description explains what this check validates.
	Description string `json:"description"`
	// Remediation suggests how to address this failure mode.
	Remediation string `json:"remediation,omitempty"`
}

// DefaultWorkflowFailureMatrix returns the default failure mode matrix for workflow demos.
func DefaultWorkflowFailureMatrix() *FailureModeMatrix {
	return &FailureModeMatrix{
		Family:    "workflow_engine",
		UpdatedAt: time.Now().UTC(),
		Checks: []FailureModeCheckConfig{
			{
				Mode:              FailureModeQueueExhaustion,
				Enabled:           true,
				WarningThreshold:  80,  // 80% capacity
				CriticalThreshold: 95,  // 95% capacity
				Description:       "Queue depth exceeds configured capacity limits",
				Remediation:       "Scale workers, increase queue capacity, or implement backpressure",
			},
			{
				Mode:              FailureModeMessageTimeout,
				Enabled:           true,
				WarningThreshold:  5000,  // 5 seconds
				CriticalThreshold: 15000, // 15 seconds
				Description:       "Workflow stage exceeds expected processing time",
				Remediation:       "Investigate stage logic, add timeout handling, consider async decomposition",
			},
			{
				Mode:              FailureModeRetryExhausted,
				Enabled:           true,
				WarningThreshold:  2, // 2 retries
				CriticalThreshold: 3, // 3 retries (exhausted)
				Description:       "Maximum retry attempts reached for a workflow stage",
				Remediation:       "Implement dead letter handling, alert on-call, review idempotency",
			},
			{
				Mode:              FailureModeDeadLetter,
				Enabled:           true,
				WarningThreshold:  1, // any dead letter
				CriticalThreshold: 10,
				Description:       "Messages routed to dead letter queue",
				Remediation:       "Analyze dead letter patterns, fix root cause, replay if safe",
			},
			{
				Mode:              FailureModeCircuitOpen,
				Enabled:           true,
				WarningThreshold:  1, // any open circuit
				CriticalThreshold: 1,
				Description:       "Circuit breaker preventing calls to downstream dependency",
				Remediation:       "Investigate downstream health, implement fallback, manual override if needed",
			},
			{
				Mode:              FailureModeIdempotencyViolation,
				Enabled:           true,
				WarningThreshold:  1,
				CriticalThreshold: 5,
				Description:       "Duplicate message detected via idempotency key",
				Remediation:       "Verify producer idempotency keys, check for client retry storms",
			},
		},
	}
}

// DefaultJobQueueFailureMatrix returns the default failure mode matrix for job queue demos.
func DefaultJobQueueFailureMatrix() *FailureModeMatrix {
	return &FailureModeMatrix{
		Family:    "job_queue",
		UpdatedAt: time.Now().UTC(),
		Checks: []FailureModeCheckConfig{
			{
				Mode:              FailureModeQueueExhaustion,
				Enabled:           true,
				WarningThreshold:  75,
				CriticalThreshold: 90,
				Description:       "Job queue depth approaching worker capacity",
				Remediation:       "Auto-scale workers, implement priority queuing, shed load",
			},
			{
				Mode:              FailureModeMessageTimeout,
				Enabled:           true,
				WarningThreshold:  30000, // 30 seconds
				CriticalThreshold: 60000, // 60 seconds
				Description:       "Job execution exceeds timeout threshold",
				Remediation:       "Optimize job logic, increase timeout, decompose into smaller jobs",
			},
			{
				Mode:              FailureModeRetryExhausted,
				Enabled:           true,
				WarningThreshold:  3,
				CriticalThreshold: 5,
				Description:       "Job failed after all retry attempts",
				Remediation:       "Move to dead letter, alert, analyze failure pattern",
			},
			{
				Mode:              FailureModePayloadTooLarge,
				Enabled:           true,
				WarningThreshold:  1024 * 1024, // 1MB
				CriticalThreshold: 5 * 1024 * 1024, // 5MB
				Description:       "Job payload exceeds size limits",
				Remediation:       "Use external storage (S3), compress payload, chunk large jobs",
			},
			{
				Mode:              FailureModeSerializationFailure,
				Enabled:           true,
				WarningThreshold:  1,
				CriticalThreshold: 10,
				Description:       "Job payload failed to serialize or deserialize",
				Remediation:       "Validate schema, check version compatibility, fix producer",
			},
		},
	}
}

// RunFailureChecks executes all enabled failure mode checks against the provided metrics
// and returns the results. This is used by demos and production monitoring.
func RunFailureChecks(matrix *FailureModeMatrix, metrics FailureMetrics) []FailureCheck {
	results := make([]FailureCheck, 0, len(matrix.Checks))
	now := time.Now().UTC()

	for _, check := range matrix.Checks {
		if !check.Enabled {
			continue
		}

		result := FailureCheck{
			Mode:      check.Mode,
			CheckedAt: now,
			Healthy:   true,
		}

		var value float64
		switch check.Mode {
		case FailureModeQueueExhaustion:
			value = metrics.QueueDepthPercent
			result.MetricValue = value
			result.Threshold = check.CriticalThreshold
			if value >= check.CriticalThreshold {
				result.Healthy = false
				result.Severity = "critical"
				result.Details = "Queue depth exceeds critical threshold"
			} else if value >= check.WarningThreshold {
				result.Healthy = false
				result.Severity = "warning"
				result.Details = "Queue depth exceeds warning threshold"
			}
		case FailureModeMessageTimeout:
			value = metrics.MaxStageLatencyMs
			result.MetricValue = value
			result.Threshold = check.CriticalThreshold
			if value >= check.CriticalThreshold {
				result.Healthy = false
				result.Severity = "critical"
				result.Details = "Stage latency exceeds critical threshold"
			} else if value >= check.WarningThreshold {
				result.Healthy = false
				result.Severity = "warning"
				result.Details = "Stage latency exceeds warning threshold"
			}
		case FailureModeRetryExhausted:
			value = float64(metrics.MaxRetriesAttempted)
			result.MetricValue = value
			result.Threshold = check.CriticalThreshold
			if value >= check.CriticalThreshold {
				result.Healthy = false
				result.Severity = "critical"
				result.Details = "Maximum retries exhausted"
			} else if value >= check.WarningThreshold {
				result.Healthy = false
				result.Severity = "warning"
				result.Details = "Approaching retry limit"
			}
		case FailureModeDeadLetter:
			value = float64(metrics.DeadLetterCount)
			result.MetricValue = value
			result.Threshold = check.CriticalThreshold
			if value >= check.CriticalThreshold {
				result.Healthy = false
				result.Severity = "critical"
				result.Details = "High dead letter volume"
			} else if value >= check.WarningThreshold {
				result.Healthy = false
				result.Severity = "warning"
				result.Details = "Dead letters detected"
			}
		case FailureModeCircuitOpen:
			value = float64(metrics.OpenCircuitCount)
			result.MetricValue = value
			result.Threshold = check.CriticalThreshold
			if value >= check.CriticalThreshold {
				result.Healthy = false
				result.Severity = "critical"
				result.Details = "Circuit breaker open"
			}
		case FailureModeIdempotencyViolation:
			value = float64(metrics.DuplicateCount)
			result.MetricValue = value
			result.Threshold = check.CriticalThreshold
			if value >= check.CriticalThreshold {
				result.Healthy = false
				result.Severity = "critical"
				result.Details = "Multiple idempotency violations"
			} else if value >= check.WarningThreshold {
				result.Healthy = false
				result.Severity = "warning"
				result.Details = "Duplicate message detected"
			}
		case FailureModePayloadTooLarge:
			value = metrics.MaxPayloadSizeBytes
			result.MetricValue = value
			result.Threshold = check.CriticalThreshold
			if value >= check.CriticalThreshold {
				result.Healthy = false
				result.Severity = "critical"
				result.Details = "Payload exceeds critical size limit"
			} else if value >= check.WarningThreshold {
				result.Healthy = false
				result.Severity = "warning"
				result.Details = "Payload exceeds warning size limit"
			}
		case FailureModeSerializationFailure:
			value = float64(metrics.SerializationErrorCount)
			result.MetricValue = value
			result.Threshold = check.CriticalThreshold
			if value >= check.CriticalThreshold {
				result.Healthy = false
				result.Severity = "critical"
				result.Details = "High serialization error rate"
			} else if value >= check.WarningThreshold {
				result.Healthy = false
				result.Severity = "warning"
				result.Details = "Serialization errors detected"
			}
		}

		results = append(results, result)
	}

	return results
}

// FailureMetrics holds the runtime metrics used to evaluate failure mode checks.
type FailureMetrics struct {
	// QueueDepthPercent is the current queue depth as a percentage of capacity (0-100).
	QueueDepthPercent float64 `json:"queue_depth_percent"`
	// MaxStageLatencyMs is the maximum observed stage latency in milliseconds.
	MaxStageLatencyMs float64 `json:"max_stage_latency_ms"`
	// MaxRetriesAttempted is the maximum retries attempted for any single message.
	MaxRetriesAttempted int `json:"max_retries_attempted"`
	// DeadLetterCount is the number of messages in the dead letter queue.
	DeadLetterCount int `json:"dead_letter_count"`
	// OpenCircuitCount is the number of open circuit breakers.
	OpenCircuitCount int `json:"open_circuit_count"`
	// DuplicateCount is the number of idempotency key collisions detected.
	DuplicateCount int `json:"duplicate_count"`
	// MaxPayloadSizeBytes is the largest payload size observed in bytes.
	MaxPayloadSizeBytes float64 `json:"max_payload_size_bytes"`
	// SerializationErrorCount is the number of serialization/deserialization errors.
	SerializationErrorCount int `json:"serialization_error_count"`
}

// DemoFailureMetrics returns simulated failure metrics for demo purposes.
func DemoFailureMetrics() FailureMetrics {
	return FailureMetrics{
		QueueDepthPercent:       45.0,
		MaxStageLatencyMs:       2800.0,
		MaxRetriesAttempted:     1,
		DeadLetterCount:         0,
		OpenCircuitCount:        0,
		DuplicateCount:          0,
		MaxPayloadSizeBytes:     51200.0,
		SerializationErrorCount: 0,
	}
}

// DemoFailureMetricsWithIssues returns simulated failure metrics with some issues for demo purposes.
func DemoFailureMetricsWithIssues() FailureMetrics {
	return FailureMetrics{
		QueueDepthPercent:       87.0, // warning
		MaxStageLatencyMs:       18000.0, // critical
		MaxRetriesAttempted:     3, // critical
		DeadLetterCount:         5, // warning
		OpenCircuitCount:        1, // critical
		DuplicateCount:          2, // warning
		MaxPayloadSizeBytes:     2048000.0, // warning
		SerializationErrorCount: 3, // warning
	}
}