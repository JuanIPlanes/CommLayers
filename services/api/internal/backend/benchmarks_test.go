package backend

import (
	"testing"
	"time"

	"commlayers/services/api/internal/evidence"
)

func TestBenchmarkCollector_ComputeSummary_WithEvidence(t *testing.T) {
	collector := NewBenchmarkCollector(DefaultBenchmarkOptions())

	// Record some samples
	collector.RecordQueueDelay(10.5)
	collector.RecordQueueDelay(15.2)
	collector.RecordTimeout(5000.0)
	collector.RecordTimeout(5100.0)
	collector.RecordTransition(2.3)
	collector.RecordTransition(3.1)

	summary := collector.ComputeSummary()

	if summary == nil {
		t.Fatal("expected summary, got nil")
	}

	if summary.TotalRuns != 6 {
		t.Errorf("expected 6 total runs, got %d", summary.TotalRuns)
	}

	if summary.SuccessfulRuns != 6 {
		t.Errorf("expected 6 successful runs, got %d", summary.SuccessfulRuns)
	}

	if len(summary.Runs) != 6 {
		t.Errorf("expected 6 runs, got %d", len(summary.Runs))
	}

	// Verify each run has evidence
	for i, run := range summary.Runs {
		if run.Evidence == nil {
			t.Errorf("run %d missing evidence", i)
			continue
		}

		if err := evidence.ValidateEvidenceContract(run.Evidence); err != nil {
			t.Errorf("run %d evidence validation failed: %v", i, err)
		}

		if !run.Evidence.Success {
			t.Errorf("run %d evidence should be success", i)
		}

		if run.Evidence.SuccessCode != evidence.SuccessOK {
			t.Errorf("run %d expected success_code=ok, got %s", i, run.Evidence.SuccessCode)
		}

		if run.Evidence.Evidence.Source != evidence.SourceBenchmark {
			t.Errorf("run %d expected source=benchmark, got %s", i, run.Evidence.Evidence.Source)
		}
	}

	// Verify summary evidence
	if summary.Evidence == nil {
		t.Fatal("expected summary evidence, got nil")
	}

	if err := evidence.ValidateEvidenceContract(summary.Evidence); err != nil {
		t.Errorf("summary evidence validation failed: %v", err)
	}

	if !summary.Evidence.Success {
		t.Error("summary evidence should be success")
	}

	if summary.Evidence.Evidence.Source != evidence.SourceBenchmark {
		t.Errorf("summary expected source=benchmark, got %s", summary.Evidence.Evidence.Source)
	}
}

func TestBenchmarkRun_EvidenceFields(t *testing.T) {
	now := time.Now()
	run := BenchmarkRun{
		ID:           "test_run",
		StartedAt:    now,
		CompletedAt:  now.Add(time.Second),
		QueueDelayMs: 10.0,
		Success:      true,
		Evidence:     evidence.NewSuccessContract(evidence.Evidence{Source: evidence.SourceBenchmark, Confidence: 1.0, Staleness: evidence.StalenessFresh, MeasuredAt: now, DelayAppliedMs: 0}, evidence.SuccessOK),
	}

	if run.Evidence == nil {
		t.Fatal("expected evidence")
	}
	if !run.Evidence.Success {
		t.Error("expected success")
	}
}

func TestBenchmarkSummary_EvidenceFields(t *testing.T) {
	now := time.Now()
	summary := BenchmarkSummary{
		ComputedAt: now,
		Evidence:   evidence.NewSuccessContract(evidence.Evidence{Source: evidence.SourceBenchmark, Confidence: 1.0, Staleness: evidence.StalenessFresh, MeasuredAt: now, DelayAppliedMs: 0}, evidence.SuccessOK),
	}

	if summary.Evidence == nil {
		t.Fatal("expected evidence")
	}
	if !summary.Evidence.Success {
		t.Error("expected success")
	}
}