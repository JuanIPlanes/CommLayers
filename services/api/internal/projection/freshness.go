package projection

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

// ProjectionType identifies the type of projection.
type ProjectionType string

const (
	ProjectionTypeSearch ProjectionType = "search"
	ProjectionTypeGraph  ProjectionType = "graph"
	ProjectionTypeVector ProjectionType = "vector"
)

// FreshnessStatus represents the freshness state of a projection.
type FreshnessStatus string

const (
	FreshnessStatusFresh     FreshnessStatus = "fresh"
	FreshnessStatusStale     FreshnessStatus = "stale"
	FreshnessStatusVeryStale FreshnessStatus = "very_stale"
	FreshnessStatusUnknown   FreshnessStatus = "unknown"
	FreshnessStatusReindexing FreshnessStatus = "reindexing"
)

// ProjectionFreshness tracks the freshness state of a single projection.
type ProjectionFreshness struct {
	// Type is the projection type.
	Type ProjectionType `json:"type"`
	// Status is the current freshness status.
	Status FreshnessStatus `json:"status"`
	// LastIndexedAt is when the projection was last successfully indexed.
	LastIndexedAt time.Time `json:"last_indexed_at,omitempty"`
	// NextReindexAt is when the next reindex is scheduled (if any).
	NextReindexAt time.Time `json:"next_reindex_at,omitempty"`
	// StalenessSeconds is how stale the projection is in seconds.
	StalenessSeconds int `json:"staleness_seconds,omitempty"`
	// FreshnessThresholdSeconds is the threshold (in seconds) after which the projection is considered stale.
	FreshnessThresholdSeconds int `json:"freshness_threshold_seconds"`
	// CriticalStalenessSeconds is the threshold after which the projection is considered very stale.
	CriticalStalenessSeconds int `json:"critical_staleness_seconds"`
	// DocumentCount is the number of documents in the projection.
	DocumentCount int `json:"document_count,omitempty"`
	// LastError is the last error encountered during indexing (if any).
	LastError string `json:"last_error,omitempty"`
	// ReindexTriggered indicates if a reindex has been triggered.
	ReindexTriggered bool `json:"reindex_triggered,omitempty"`
	// ReindexReason is the reason for the last reindex trigger.
	ReindexReason string `json:"reindex_reason,omitempty"`
	// UpdatedAt is when this freshness record was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// ReindexTrigger defines a condition that triggers a reindex.
type ReindexTrigger struct {
	// Name is the unique identifier for this trigger.
	Name string `json:"name"`
	// Type is the projection type this trigger applies to.
	Type ProjectionType `json:"type"`
	// Enabled indicates if this trigger is active.
	Enabled bool `json:"enabled"`
	// MaxStalenessSeconds is the maximum staleness before triggering reindex.
	MaxStalenessSeconds int `json:"max_staleness_seconds"`
	// MaxDocumentLag is the maximum document count difference before triggering reindex.
	MaxDocumentLag int `json:"max_document_lag,omitempty"`
	// ErrorRateThreshold is the error rate (0-1) that triggers reindex.
	ErrorRateThreshold float64 `json:"error_rate_threshold,omitempty"`
	// ScheduleCron is an optional cron expression for scheduled reindex.
	ScheduleCron string `json:"schedule_cron,omitempty"`
	// Description explains when this trigger fires.
	Description string `json:"description"`
}

// ReindexResult represents the result of a reindex operation.
type ReindexResult struct {
	// Type is the projection type that was reindexed.
	Type ProjectionType `json:"type"`
	// TriggeredAt is when the reindex was triggered.
	TriggeredAt time.Time `json:"triggered_at"`
	// CompletedAt is when the reindex completed (zero if still running or failed).
	CompletedAt time.Time `json:"completed_at,omitempty"`
	// Success indicates if the reindex succeeded.
	Success bool `json:"success"`
	// DocumentsIndexed is the number of documents indexed.
	DocumentsIndexed int `json:"documents_indexed"`
	// DurationMs is the reindex duration in milliseconds.
	DurationMs int64 `json:"duration_ms,omitempty"`
	// Error is the error message if the reindex failed.
	Error string `json:"error,omitempty"`
	// TriggerName is the name of the trigger that initiated this reindex.
	TriggerName string `json:"trigger_name"`
}

// ProjectionFreshnessConfig holds the configuration for projection freshness tracking.
type ProjectionFreshnessConfig struct {
	// SearchFreshnessThresholdSeconds is the staleness threshold for search projections.
	SearchFreshnessThresholdSeconds int `json:"search_freshness_threshold_seconds"`
	// SearchCriticalStalenessSeconds is the critical staleness threshold for search projections.
	SearchCriticalStalenessSeconds int `json:"search_critical_staleness_seconds"`
	// GraphFreshnessThresholdSeconds is the staleness threshold for graph projections.
	GraphFreshnessThresholdSeconds int `json:"graph_freshness_threshold_seconds"`
	// GraphCriticalStalenessSeconds is the critical staleness threshold for graph projections.
	GraphCriticalStalenessSeconds int `json:"graph_critical_staleness_seconds"`
	// VectorFreshnessThresholdSeconds is the staleness threshold for vector projections.
	VectorFreshnessThresholdSeconds int `json:"vector_freshness_threshold_seconds"`
	// VectorCriticalStalenessSeconds is the critical staleness threshold for vector projections.
	VectorCriticalStalenessSeconds int `json:"vector_critical_staleness_seconds"`
	// DefaultReindexTriggers are the default reindex triggers.
	DefaultReindexTriggers []ReindexTrigger `json:"default_reindex_triggers"`
}

// DefaultFreshnessConfig returns the default projection freshness configuration.
func DefaultFreshnessConfig() *ProjectionFreshnessConfig {
	return &ProjectionFreshnessConfig{
		SearchFreshnessThresholdSeconds:  60,   // 1 minute
		SearchCriticalStalenessSeconds:   300,  // 5 minutes
		GraphFreshnessThresholdSeconds:   120,  // 2 minutes
		GraphCriticalStalenessSeconds:    600,  // 10 minutes
		VectorFreshnessThresholdSeconds:  180,  // 3 minutes
		VectorCriticalStalenessSeconds:   900,  // 15 minutes
		DefaultReindexTriggers: []ReindexTrigger{
			{
				Name:                  "search_staleness",
				Type:                  ProjectionTypeSearch,
				Enabled:               true,
				MaxStalenessSeconds:   300,
				MaxDocumentLag:        1000,
				ErrorRateThreshold:    0.05,
				Description:           "Trigger reindex when search projection exceeds 5min staleness or 1000 doc lag",
			},
			{
				Name:                  "graph_staleness",
				Type:                  ProjectionTypeGraph,
				Enabled:               true,
				MaxStalenessSeconds:   600,
				MaxDocumentLag:        500,
				ErrorRateThreshold:    0.05,
				Description:           "Trigger reindex when graph projection exceeds 10min staleness or 500 doc lag",
			},
			{
				Name:                  "vector_staleness",
				Type:                  ProjectionTypeVector,
				Enabled:               true,
				MaxStalenessSeconds:   900,
				MaxDocumentLag:        2000,
				ErrorRateThreshold:    0.10,
				Description:           "Trigger reindex when vector projection exceeds 15min staleness or 2000 doc lag",
			},
			{
				Name:                  "scheduled_search_reindex",
				Type:                  ProjectionTypeSearch,
				Enabled:               true,
				MaxStalenessSeconds:   0, // schedule-based only
				ScheduleCron:          "0 */15 * * * *", // every 15 minutes
				Description:           "Scheduled reindex for search projection every 15 minutes",
			},
			{
				Name:                  "scheduled_graph_reindex",
				Type:                  ProjectionTypeGraph,
				Enabled:               true,
				MaxStalenessSeconds:   0,
				ScheduleCron:          "0 */30 * * * *", // every 30 minutes
				Description:           "Scheduled reindex for graph projection every 30 minutes",
			},
			{
				Name:                  "scheduled_vector_reindex",
				Type:                  ProjectionTypeVector,
				Enabled:               true,
				MaxStalenessSeconds:   0,
				ScheduleCron:          "0 */60 * * * *", // every hour
				Description:           "Scheduled reindex for vector projection every hour",
			},
		},
	}
}

// FreshnessTracker tracks projection freshness and evaluates reindex triggers.
type FreshnessTracker struct {
	mu             sync.RWMutex
	config         *ProjectionFreshnessConfig
	freshness      map[ProjectionType]*ProjectionFreshness
	reindexHistory []ReindexResult
	reindexMu      sync.Mutex
}

// NewFreshnessTracker creates a new FreshnessTracker with the given config.
func NewFreshnessTracker(config *ProjectionFreshnessConfig) *FreshnessTracker {
	if config == nil {
		config = DefaultFreshnessConfig()
	}

	tracker := &FreshnessTracker{
		config:         config,
		freshness:      make(map[ProjectionType]*ProjectionFreshness),
		reindexHistory: make([]ReindexResult, 0, 100),
	}

	// Initialize freshness for each projection type
	for _, pt := range []ProjectionType{ProjectionTypeSearch, ProjectionTypeGraph, ProjectionTypeVector} {
		tracker.freshness[pt] = &ProjectionFreshness{
			Type:                       pt,
			Status:                     FreshnessStatusUnknown,
			FreshnessThresholdSeconds:  tracker.thresholdForType(pt),
			CriticalStalenessSeconds:   tracker.criticalThresholdForType(pt),
			UpdatedAt:                  time.Now().UTC(),
		}
	}

	return tracker
}

// thresholdForType returns the freshness threshold for a projection type.
func (ft *FreshnessTracker) thresholdForType(pt ProjectionType) int {
	switch pt {
	case ProjectionTypeSearch:
		return ft.config.SearchFreshnessThresholdSeconds
	case ProjectionTypeGraph:
		return ft.config.GraphFreshnessThresholdSeconds
	case ProjectionTypeVector:
		return ft.config.VectorFreshnessThresholdSeconds
	default:
		return 60
	}
}

// criticalThresholdForType returns the critical staleness threshold for a projection type.
func (ft *FreshnessTracker) criticalThresholdForType(pt ProjectionType) int {
	switch pt {
	case ProjectionTypeSearch:
		return ft.config.SearchCriticalStalenessSeconds
	case ProjectionTypeGraph:
		return ft.config.GraphCriticalStalenessSeconds
	case ProjectionTypeVector:
		return ft.config.VectorCriticalStalenessSeconds
	default:
		return 300
	}
}

// UpdateFreshness updates the freshness state for a projection type.
func (ft *FreshnessTracker) UpdateFreshness(ctx context.Context, pt ProjectionType, documentCount int, lastError string) {
	ft.mu.Lock()
	defer ft.mu.Unlock()

	now := time.Now().UTC()
	f := ft.freshness[pt]
	if f == nil {
		f = &ProjectionFreshness{
			Type:                      pt,
			FreshnessThresholdSeconds: ft.thresholdForType(pt),
			CriticalStalenessSeconds:  ft.criticalThresholdForType(pt),
		}
		ft.freshness[pt] = f
	}

	f.DocumentCount = documentCount
	f.LastError = lastError
	f.UpdatedAt = now

	if !f.LastIndexedAt.IsZero() {
		f.StalenessSeconds = int(now.Sub(f.LastIndexedAt).Seconds())
		f.Status = ft.evaluateStatus(f)
	} else {
		f.Status = FreshnessStatusUnknown
		f.StalenessSeconds = -1
	}

	// Check reindex triggers
	ft.checkTriggersLocked(ctx, pt, f)
}

// evaluateStatus determines the freshness status based on staleness.
func (ft *FreshnessTracker) evaluateStatus(f *ProjectionFreshness) FreshnessStatus {
	if f.StalenessSeconds < 0 {
		return FreshnessStatusUnknown
	}
	if f.ReindexTriggered {
		return FreshnessStatusReindexing
	}
	if f.StalenessSeconds >= f.CriticalStalenessSeconds {
		return FreshnessStatusVeryStale
	}
	if f.StalenessSeconds >= f.FreshnessThresholdSeconds {
		return FreshnessStatusStale
	}
	return FreshnessStatusFresh
}

// checkTriggersLocked evaluates reindex triggers for a projection (must hold lock).
func (ft *FreshnessTracker) checkTriggersLocked(ctx context.Context, pt ProjectionType, f *ProjectionFreshness) {
	for _, trigger := range ft.config.DefaultReindexTriggers {
		if !trigger.Enabled || trigger.Type != pt {
			continue
		}

		shouldTrigger := false
		reason := ""

		// Staleness-based trigger
		if trigger.MaxStalenessSeconds > 0 && f.StalenessSeconds >= trigger.MaxStalenessSeconds {
			shouldTrigger = true
			reason = "staleness threshold exceeded"
		}

		// Document lag trigger (would need source count - simplified here)
		if trigger.MaxDocumentLag > 0 && f.DocumentCount > trigger.MaxDocumentLag {
			shouldTrigger = true
			if reason != "" {
				reason += " and "
			}
			reason += "document lag exceeded"
		}

		// Error rate trigger (simplified - would need actual error tracking)
		if trigger.ErrorRateThreshold > 0 && f.LastError != "" {
			shouldTrigger = true
			if reason != "" {
				reason += " and "
			}
			reason += "error rate threshold exceeded"
		}

		// Schedule-based triggers are handled externally (cron)

		now := time.Now().UTC()
		if shouldTrigger && !f.ReindexTriggered {
			f.ReindexTriggered = true
			f.ReindexReason = trigger.Name + ": " + reason
			f.NextReindexAt = now.Add(5 * time.Second) // Simulated async start
			// In a real implementation, this would trigger an async reindex job
			// For demo, we simulate completion
			go ft.simulateReindex(ctx, pt, trigger.Name, reason)
		}
	}
}

// simulateReindex simulates an async reindex operation for demo purposes.
func (ft *FreshnessTracker) simulateReindex(ctx context.Context, pt ProjectionType, triggerName, reason string) {
	start := time.Now().UTC()

	// Simulate reindex work
	select {
	case <-ctx.Done():
		return
	case <-time.After(2 * time.Second):
	}

	ft.reindexMu.Lock()
	defer ft.reindexMu.Unlock()

	completedAt := time.Now().UTC()
	duration := completedAt.Sub(start).Milliseconds()

	result := ReindexResult{
		Type:             pt,
		TriggeredAt:      start,
		CompletedAt:      completedAt,
		Success:          true,
		DocumentsIndexed: ft.freshness[pt].DocumentCount,
		DurationMs:       duration,
		TriggerName:      triggerName,
	}

	ft.reindexHistory = append(ft.reindexHistory, result)
	if len(ft.reindexHistory) > 100 {
		ft.reindexHistory = ft.reindexHistory[1:]
	}

	// Update freshness after reindex
	ft.mu.Lock()
	defer ft.mu.Unlock()
	if f := ft.freshness[pt]; f != nil {
		f.LastIndexedAt = completedAt
		f.ReindexTriggered = false
		f.ReindexReason = ""
		f.NextReindexAt = time.Time{}
		f.LastError = ""
		f.Status = FreshnessStatusFresh
		f.StalenessSeconds = 0
		f.UpdatedAt = completedAt
	}
}

// GetFreshness returns the current freshness state for a projection type.
func (ft *FreshnessTracker) GetFreshness(pt ProjectionType) *ProjectionFreshness {
	ft.mu.RLock()
	defer ft.mu.RUnlock()
	if f := ft.freshness[pt]; f != nil {
		// Return a copy to avoid race conditions
		copy := *f
		return &copy
	}
	return nil
}

// GetAllFreshness returns the freshness state for all projection types.
func (ft *FreshnessTracker) GetAllFreshness() map[ProjectionType]*ProjectionFreshness {
	ft.mu.RLock()
	defer ft.mu.RUnlock()
	result := make(map[ProjectionType]*ProjectionFreshness)
	for pt, f := range ft.freshness {
		copy := *f
		result[pt] = &copy
	}
	return result
}

// GetReindexHistory returns the reindex history, most recent first.
func (ft *FreshnessTracker) GetReindexHistory(limit int) []ReindexResult {
	ft.reindexMu.Lock()
	defer ft.reindexMu.Unlock()

	if limit <= 0 || limit > len(ft.reindexHistory) {
		limit = len(ft.reindexHistory)
	}

	// Return most recent first
	result := make([]ReindexResult, limit)
	for i := 0; i < limit; i++ {
		result[i] = ft.reindexHistory[len(ft.reindexHistory)-1-i]
	}
	return result
}

// TriggerManualReindex manually triggers a reindex for a projection type.
func (ft *FreshnessTracker) TriggerManualReindex(ctx context.Context, pt ProjectionType, reason string) error {
	ft.mu.Lock()
	f := ft.freshness[pt]
	if f == nil {
		ft.mu.Unlock()
		return ErrProjectionNotFound
	}
	if f.ReindexTriggered {
		ft.mu.Unlock()
		return ErrReindexInProgress
	}
	f.ReindexTriggered = true
	f.ReindexReason = "manual: " + reason
	f.NextReindexAt = time.Now().UTC().Add(5 * time.Second)
	ft.mu.Unlock()

	go ft.simulateReindex(ctx, pt, "manual", reason)
	return nil
}

// DemoFreshnessData returns simulated freshness data for demo purposes.
func DemoFreshnessData() map[ProjectionType]*ProjectionFreshness {
	now := time.Now().UTC()
	return map[ProjectionType]*ProjectionFreshness{
		ProjectionTypeSearch: {
			Type:                      ProjectionTypeSearch,
			Status:                    FreshnessStatusFresh,
			LastIndexedAt:             now.Add(-30 * time.Second),
			StalenessSeconds:          30,
			FreshnessThresholdSeconds: 60,
			CriticalStalenessSeconds:  300,
			DocumentCount:             1247,
			UpdatedAt:                 now,
		},
		ProjectionTypeGraph: {
			Type:                      ProjectionTypeGraph,
			Status:                    FreshnessStatusStale,
			LastIndexedAt:             now.Add(-180 * time.Second),
			StalenessSeconds:          180,
			FreshnessThresholdSeconds: 120,
			CriticalStalenessSeconds:  600,
			DocumentCount:             892,
			UpdatedAt:                 now,
		},
		ProjectionTypeVector: {
			Type:                      ProjectionTypeVector,
			Status:                    FreshnessStatusVeryStale,
			LastIndexedAt:             now.Add(-1200 * time.Second),
			StalenessSeconds:          1200,
			FreshnessThresholdSeconds: 180,
			CriticalStalenessSeconds:  900,
			DocumentCount:             3421,
			ReindexTriggered:          true,
			ReindexReason:             "vector_staleness: staleness threshold exceeded",
			NextReindexAt:             now.Add(5 * time.Second),
			UpdatedAt:                 now,
		},
	}
}

// DemoReindexHistory returns simulated reindex history for demo purposes.
func DemoReindexHistory() []ReindexResult {
	now := time.Now().UTC()
	return []ReindexResult{
		{
			Type:             ProjectionTypeSearch,
			TriggeredAt:      now.Add(-20 * time.Minute),
			CompletedAt:      now.Add(-19 * time.Minute).Add(30 * time.Second),
			Success:          true,
			DocumentsIndexed: 1247,
			DurationMs:       30000,
			TriggerName:      "scheduled_search_reindex",
		},
		{
			Type:             ProjectionTypeGraph,
			TriggeredAt:      now.Add(-45 * time.Minute),
			CompletedAt:      now.Add(-44 * time.Minute).Add(45 * time.Second),
			Success:          true,
			DocumentsIndexed: 892,
			DurationMs:       45000,
			TriggerName:      "scheduled_graph_reindex",
		},
		{
			Type:             ProjectionTypeVector,
			TriggeredAt:      now.Add(-2 * time.Minute),
			CompletedAt:      time.Time{}, // Still running
			Success:          false,
			DocumentsIndexed: 0,
			DurationMs:       0,
			TriggerName:      "vector_staleness",
		},
	}
}

// Error types
type projectionError string

func (e projectionError) Error() string { return string(e) }

const (
	ErrProjectionNotFound   = projectionError("projection not found")
	ErrReindexInProgress    = projectionError("reindex already in progress")
	ErrInvalidConfiguration = projectionError("invalid freshness configuration")
)

// FreshnessResponse is the API response for projection freshness.
type FreshnessResponse struct {
	Projections map[ProjectionType]*ProjectionFreshness `json:"projections"`
	ReindexHistory []ReindexResult `json:"reindex_history,omitempty"`
	Config *ProjectionFreshnessConfig `json:"config,omitempty"`
}

// FreshnessCheckResult represents the result of a freshness check for API responses.
type FreshnessCheckResult struct {
	ProjectionType string `json:"projection_type"`
	Status         string `json:"status"`
	StalenessSeconds int `json:"staleness_seconds"`
	ThresholdSeconds int `json:"threshold_seconds"`
	CriticalThresholdSeconds int `json:"critical_threshold_seconds"`
	DocumentCount int `json:"document_count"`
	LastIndexedAt string `json:"last_indexed_at,omitempty"`
	ReindexTriggered bool `json:"reindex_triggered"`
	ReindexReason string `json:"reindex_reason,omitempty"`
}

// ToCheckResult converts ProjectionFreshness to FreshnessCheckResult for API.
func (f *ProjectionFreshness) ToCheckResult() FreshnessCheckResult {
	return FreshnessCheckResult{
		ProjectionType:            string(f.Type),
		Status:                    string(f.Status),
		StalenessSeconds:          f.StalenessSeconds,
		ThresholdSeconds:          f.FreshnessThresholdSeconds,
		CriticalThresholdSeconds:  f.CriticalStalenessSeconds,
		DocumentCount:             f.DocumentCount,
		LastIndexedAt:             f.LastIndexedAt.Format(time.RFC3339),
		ReindexTriggered:          f.ReindexTriggered,
		ReindexReason:             f.ReindexReason,
	}
}

// MarshalJSON implements custom JSON marshaling for ProjectionFreshness.
func (f *ProjectionFreshness) MarshalJSON() ([]byte, error) {
	type Alias ProjectionFreshness
	return json.Marshal(&struct {
		*Alias
		LastIndexedAt string `json:"last_indexed_at,omitempty"`
		NextReindexAt string `json:"next_reindex_at,omitempty"`
		UpdatedAt     string `json:"updated_at"`
	}{
		Alias:         (*Alias)(f),
		LastIndexedAt: f.LastIndexedAt.Format(time.RFC3339),
		NextReindexAt: f.NextReindexAt.Format(time.RFC3339),
		UpdatedAt:     f.UpdatedAt.Format(time.RFC3339),
	})
}

// MarshalJSON implements custom JSON marshaling for ReindexResult.
func (r *ReindexResult) MarshalJSON() ([]byte, error) {
	type Alias ReindexResult
	return json.Marshal(&struct {
		*Alias
		TriggeredAt  string `json:"triggered_at"`
		CompletedAt  string `json:"completed_at,omitempty"`
	}{
		Alias:       (*Alias)(r),
		TriggeredAt: r.TriggeredAt.Format(time.RFC3339),
		CompletedAt: r.CompletedAt.Format(time.RFC3339),
	})
}