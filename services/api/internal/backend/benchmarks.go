package backend

import (
	"fmt"
	"sync"
	"time"

	"commlayers/services/api/internal/evidence"
)

// AsyncVisibilityBenchmarks defines how async queue delay, timeout visibility,
// and delayed state transitions are measured and exposed to the frontend.
//
// These benchmarks provide the observable metrics that the frontend consumes
// to understand async operation progress, delays, and state transitions.
type AsyncVisibilityBenchmarks struct {
	// QueueDelayMetrics holds metrics about async queue delay.
	QueueDelayMetrics QueueDelayMetrics `json:"queue_delay"`
	// TimeoutMetrics holds metrics about timeout visibility.
	TimeoutMetrics TimeoutMetrics `json:"timeout"`
	// StateTransitionMetrics holds metrics about delayed state transitions.
	StateTransitionMetrics StateTransitionMetrics `json:"state_transitions"`
}

// QueueDelayMetrics holds metrics about async queue delay measurement.
type QueueDelayMetrics struct {
	// AverageDelayMs is the average queue delay in milliseconds.
	AverageDelayMs float64 `json:"average_delay_ms"`
	// P95DelayMs is the 95th percentile queue delay in milliseconds.
	P95DelayMs float64 `json:"p95_delay_ms"`
	// P99DelayMs is the 99th percentile queue delay in milliseconds.
	P99DelayMs float64 `json:"p99_delay_ms"`
	// MaxDelayMs is the maximum observed queue delay in milliseconds.
	MaxDelayMs float64 `json:"max_delay_ms"`
	// MinDelayMs is the minimum observed queue delay in milliseconds.
	MinDelayMs float64 `json:"min_delay_ms"`
	// DelaySamples is the number of samples used for the above statistics.
	DelaySamples int `json:"delay_samples"`
	// LastMeasuredAt is when the delay metrics were last measured.
	LastMeasuredAt time.Time `json:"last_measured_at"`
}

// TimeoutMetrics holds metrics about timeout visibility.
type TimeoutMetrics struct {
	// AverageTimeoutMs is the average timeout duration in milliseconds.
	AverageTimeoutMs float64 `json:"average_timeout_ms"`
	// P95TimeoutMs is the 95th percentile timeout in milliseconds.
	P95TimeoutMs float64 `json:"p95_timeout_ms"`
	// P99TimeoutMs is the 99th percentile timeout in milliseconds.
	P99TimeoutMs float64 `json:"p99_timeout_ms"`
	// MaxTimeoutMs is the maximum observed timeout in milliseconds.
	MaxTimeoutMs float64 `json:"max_timeout_ms"`
	// MinTimeoutMs is the minimum observed timeout in milliseconds.
	MinTimeoutMs float64 `json:"min_timeout_ms"`
	// TimeoutSamples is the number of samples used for the above statistics.
	TimeoutSamples int `json:"timeout_samples"`
	// LastMeasuredAt is when the timeout metrics were last measured.
	LastMeasuredAt time.Time `json:"last_measured_at"`
	// EnforcedTimeoutMs is the enforced timeout limit in milliseconds.
	EnforcedTimeoutMs float64 `json:"enforced_timeout_ms"`
}

// StateTransitionMetrics holds metrics about delayed state transitions.
type StateTransitionMetrics struct {
	// AverageTransitionMs is the average state transition duration in milliseconds.
	AverageTransitionMs float64 `json:"average_transition_ms"`
	// P95TransitionMs is the 95th percentile transition duration in milliseconds.
	P95TransitionMs float64 `json:"p95_transition_ms"`
	// P99TransitionMs is the 99th percentile transition duration in milliseconds.
	P99TransitionMs float64 `json:"p99_transition_ms"`
	// MaxTransitionMs is the maximum observed transition duration in milliseconds.
	MaxTransitionMs float64 `json:"max_transition_ms"`
	// MinTransitionMs is the minimum observed transition duration in milliseconds.
	MinTransitionMs float64 `json:"min_transition_ms"`
	// TransitionSamples is the number of samples used for the above statistics.
	TransitionSamples int `json:"transition_samples"`
	// LastMeasuredAt is when the transition metrics were last measured.
	LastMeasuredAt time.Time `json:"last_measured_at"`
	// AllowedTransitions lists the state transitions that are permitted.
	AllowedTransitions []string `json:"allowed_transitions,omitempty"`
	// BlockedTransitions lists the state transitions that are currently blocked.
	BlockedTransitions []string `json:"blocked_transitions,omitempty"`
}

// BenchmarkRun represents a single benchmark run with its results.
type BenchmarkRun struct {
	// ID is the unique identifier for the benchmark run.
	ID string `json:"id"`
	// StartedAt is when the benchmark run started.
	StartedAt time.Time `json:"started_at"`
	// CompletedAt is when the benchmark run completed.
	CompletedAt time.Time `json:"completed_at"`
	// QueueDelayMs is the queue delay observed during the run.
	QueueDelayMs float64 `json:"queue_delay_ms"`
	// TimeoutMs is the timeout observed during the run.
	TimeoutMs float64 `json:"timeout_ms"`
	// TransitionMs is the state transition duration observed during the run.
	TransitionMs float64 `json:"transition_ms"`
	// Success indicates if the benchmark run was successful.
	Success bool `json:"success"`
	// Notes contains any additional notes about the benchmark run.
	Notes string `json:"notes,omitempty"`
	// Evidence is the formal evidence contract for this run.
	Evidence *evidence.EvidenceContract `json:"evidence,omitempty"`
}

// BenchmarkSummary is a summary of multiple benchmark runs.
type BenchmarkSummary struct {
	// Runs is the list of benchmark runs included in the summary.
	Runs []BenchmarkRun `json:"runs"`
	// AverageQueueDelayMs is the average queue delay across all runs.
	AverageQueueDelayMs float64 `json:"average_queue_delay_ms"`
	// AverageTimeoutMs is the average timeout across all runs.
	AverageTimeoutMs float64 `json:"average_timeout_ms"`
	// AverageTransitionMs is the average transition duration across all runs.
	AverageTransitionMs float64 `json:"average_transition_ms"`
	// TotalRuns is the total number of benchmark runs.
	TotalRuns int `json:"total_runs"`
	// SuccessfulRuns is the number of successful benchmark runs.
	SuccessfulRuns int `json:"successful_runs"`
	// Timestamp when the summary was computed.
	ComputedAt time.Time `json:"computed_at"`
	// Evidence is the aggregate evidence contract for the summary.
	Evidence *evidence.EvidenceContract `json:"evidence,omitempty"`
}

// BenchmarkOptions configures how async visibility benchmarks are collected.
type BenchmarkOptions struct {
	// QueueDelaySampleInterval is the interval at which queue delay is sampled.
	QueueDelaySampleInterval time.Duration `json:"queue_delay_sample_interval"`
	// TimeoutSampleInterval is the interval at which timeout is sampled.
	TimeoutSampleInterval time.Duration `json:"timeout_sample_interval"`
	// StateTransitionSampleInterval is the interval at which state transitions are sampled.
	StateTransitionSampleInterval time.Duration `json:"state_transition_sample_interval"`
	// MaxSampleHistory is the maximum number of samples to retain per metric.
	MaxSampleHistory int `json:"max_sample_history"`
	// EnabledMetrics specifies which metrics to enable.
	EnabledMetrics []string `json:"enabled_metrics"`
}

// DefaultBenchmarkOptions returns the default benchmark options.
func DefaultBenchmarkOptions() BenchmarkOptions {
	return BenchmarkOptions{
		QueueDelaySampleInterval:      1 * time.Second,
		TimeoutSampleInterval:         5 * time.Second,
		StateTransitionSampleInterval: 1 * time.Second,
		MaxSampleHistory:              1000,
		EnabledMetrics:                []string{"queue_delay", "timeout", "state_transitions"},
	}
}

// BenchmarkCollector collects async visibility benchmark data.
type BenchmarkCollector struct {
	// Options configures the collector.
	Options BenchmarkOptions
	// queueDelaySamples holds the most recent queue delay samples.
	queueDelaySamples []float64
	// timeoutSamples holds the most recent timeout samples.
	timeoutSamples []float64
	// transitionSamples holds the most recent state transition samples.
	transitionSamples []float64
	// sampleMu protects the sample slices.
	sampleMu sync.Mutex
}

// NewBenchmarkCollector creates a new BenchmarkCollector with the given options.
func NewBenchmarkCollector(options BenchmarkOptions) *BenchmarkCollector {
	return &BenchmarkCollector{
		Options:      options,
		queueDelaySamples: make([]float64, 0, options.MaxSampleHistory),
		timeoutSamples:  make([]float64, 0, options.MaxSampleHistory),
		transitionSamples: make([]float64, 0, options.MaxSampleHistory),
	}
}

// RecordQueueDelay records a queue delay sample.
func (bc *BenchmarkCollector) RecordQueueDelay(delayMs float64) {
	bc.sampleMu.Lock()
	defer bc.sampleMu.Unlock()

	bc.queueDelaySamples = append(bc.queueDelaySamples, delayMs)
	if len(bc.queueDelaySamples) > bc.Options.MaxSampleHistory {
		bc.queueDelaySamples = bc.queueDelaySamples[1:]
	}
}

// RecordTimeout records a timeout sample.
func (bc *BenchmarkCollector) RecordTimeout(timeoutMs float64) {
	bc.sampleMu.Lock()
	defer bc.sampleMu.Unlock()

	bc.timeoutSamples = append(bc.timeoutSamples, timeoutMs)
	if len(bc.timeoutSamples) > bc.Options.MaxSampleHistory {
		bc.timeoutSamples = bc.timeoutSamples[1:]
	}
}

// RecordTransition records a state transition duration sample.
func (bc *BenchmarkCollector) RecordTransition(transitionMs float64) {
	bc.sampleMu.Lock()
	defer bc.sampleMu.Unlock()

	bc.transitionSamples = append(bc.transitionSamples, transitionMs)
	if len(bc.transitionSamples) > bc.Options.MaxSampleHistory {
		bc.transitionSamples = bc.transitionSamples[1:]
	}
}

// ComputeSummary computes the benchmark summary from the collected samples.
func (bc *BenchmarkCollector) ComputeSummary() *BenchmarkSummary {
	bc.sampleMu.Lock()
	defer bc.sampleMu.Unlock()

	now := time.Now()
	totalSamples := len(bc.queueDelaySamples) + len(bc.timeoutSamples) + len(bc.transitionSamples)

	summary := &BenchmarkSummary{
		Runs:        make([]BenchmarkRun, 0, totalSamples),
		TotalRuns:   totalSamples,
		ComputedAt:  now,
	}

	var runID int
	successCount := 0

	// Create runs for queue delay samples
	for _, d := range bc.queueDelaySamples {
		runID++
		run := BenchmarkRun{
			ID:           fmt.Sprintf("queue_delay_%d", runID),
			StartedAt:    now,
			CompletedAt:  now,
			QueueDelayMs: d,
			Success:      true,
			Evidence:     createBenchmarkEvidence(fmt.Sprintf("queue_delay_%d", runID), d, 0, 0, 1),
		}
		summary.Runs = append(summary.Runs, run)
		successCount++
		if summary.AverageQueueDelayMs == 0 {
			summary.AverageQueueDelayMs = d
		} else {
			summary.AverageQueueDelayMs = (summary.AverageQueueDelayMs + d) / 2
		}
	}

	// Create runs for timeout samples
	for _, t := range bc.timeoutSamples {
		runID++
		run := BenchmarkRun{
			ID:         fmt.Sprintf("timeout_%d", runID),
			StartedAt:  now,
			CompletedAt: now,
			TimeoutMs:  t,
			Success:    true,
			Evidence:   createBenchmarkEvidence(fmt.Sprintf("timeout_%d", runID), 0, t, 0, 1),
		}
		summary.Runs = append(summary.Runs, run)
		successCount++
		if summary.AverageTimeoutMs == 0 {
			summary.AverageTimeoutMs = t
		} else {
			summary.AverageTimeoutMs = (summary.AverageTimeoutMs + t) / 2
		}
	}

	// Create runs for transition samples
	for _, t := range bc.transitionSamples {
		runID++
		run := BenchmarkRun{
			ID:            fmt.Sprintf("transition_%d", runID),
			StartedAt:     now,
			CompletedAt:   now,
			TransitionMs:  t,
			Success:       true,
			Evidence:      createBenchmarkEvidence(fmt.Sprintf("transition_%d", runID), 0, 0, t, 1),
		}
		summary.Runs = append(summary.Runs, run)
		successCount++
		if summary.AverageTransitionMs == 0 {
			summary.AverageTransitionMs = t
		} else {
			summary.AverageTransitionMs = (summary.AverageTransitionMs + t) / 2
		}
	}

	summary.SuccessfulRuns = successCount

	// Create aggregate evidence for the summary
	summary.Evidence = createAggregateEvidence(summary)

	return summary
}

// createBenchmarkEvidence creates an evidence contract for a single benchmark run.
func createBenchmarkEvidence(runID string, queueDelay, timeout, transition float64, sampleCount int) *evidence.EvidenceContract {
	benchEv := evidence.NewBenchmarkEvidence(runID, queueDelay, timeout, transition, sampleCount)
	return evidence.NewSuccessContract(benchEv.Evidence, evidence.SuccessOK)
}

// createAggregateEvidence creates an aggregate evidence contract for a benchmark summary.
func createAggregateEvidence(summary *BenchmarkSummary) *evidence.EvidenceContract {
	aggEv := evidence.Evidence{
		Source:         evidence.SourceBenchmark,
		Confidence:     1.0,
		Staleness:      evidence.StalenessFresh,
		MeasuredAt:     summary.ComputedAt,
		DelayAppliedMs: 0,
	}
	return evidence.NewSuccessContract(aggEv, evidence.SuccessOK)
}