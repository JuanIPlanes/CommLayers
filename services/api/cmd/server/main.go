package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type app struct {
	startedAt time.Time
	jobs      sync.Map
	streams   []stream
}

type stream struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Status       string   `json:"status"`
	Dependencies []string `json:"dependencies"`
}

type comparison struct {
	Name           string   `json:"name"`
	RecommendedFor []string `json:"recommendedFor"`
	AvoidWhen      []string `json:"avoidWhen"`
	LatencyClass   string   `json:"latencyClass"`
	Notes          []string `json:"notes"`
}

type jobStatus struct {
	ID               string   `json:"id"`
	Status           string   `json:"status"`
	Progress         int      `json:"progress"`
	CurrentStep      string   `json:"currentStep"`
	DelayAppliedMS   int      `json:"delayAppliedMs"`
	TimeoutThreshold int      `json:"timeoutThresholdMs"`
	Timeline         []string `json:"timeline"`
	CreatedAt        string   `json:"createdAt"`
	UpdatedAt        string   `json:"updatedAt"`
}

type responseEnvelope struct {
	Data       any    `json:"data"`
	Source     string `json:"source"`
	Confidence string `json:"confidence"`
	Staleness  string `json:"staleness"`
	MeasuredAt string `json:"measuredAt"`
}

func main() {
	application := &app{
		startedAt: time.Now().UTC(),
		streams: []stream{
			{ID: "86agbv0k4", Title: "Backend abstraction layer matrix", Description: "Freeze normalized backend vocabulary and comparison rules.", Status: "active", Dependencies: []string{}},
			{ID: "86agbv0k3", Title: "Keycloak and Vault security bootstrap", Description: "Define security bootstrap defaults after vocabulary freeze.", Status: "active", Dependencies: []string{"86agbv0k4"}},
			{ID: "86agbv0k7", Title: "Canonical data platform", Description: "Define canonical entities, ownership, and audit lineage.", Status: "active", Dependencies: []string{"86agbv0k4"}},
			{ID: "86agbv0k9", Title: "Benchmark and observability framework", Description: "Close the first wave with benchmark and telemetry contracts.", Status: "active", Dependencies: []string{"86agbv0k4", "86agbv0k3", "86agbv0k7"}},
			{ID: "86agbv0k5", Title: "Transport family implementations", Description: "Deferred until first-wave approval.", Status: "deferred", Dependencies: []string{"86agbv0k9"}},
			{ID: "86agbv0ka", Title: "Messaging and workflow families", Description: "Deferred until first-wave approval.", Status: "deferred", Dependencies: []string{"86agbv0k9"}},
			{ID: "86agbv0kb", Title: "Sync, replication, and coordination families", Description: "Deferred until first-wave approval.", Status: "deferred", Dependencies: []string{"86agbv0k9"}},
			{ID: "86agbv0kc", Title: "Search, graph, and vector projections", Description: "Deferred until first-wave approval.", Status: "deferred", Dependencies: []string{"86agbv0k9"}},
			{ID: "86agbv0kf", Title: "Localization for frontend-facing responses", Description: "Deferred until backend-family work is stable.", Status: "held", Dependencies: []string{"86agbv0k9"}},
			{ID: "86agbv0kg", Title: "Frontend comparative visualization", Description: "Deferred until backend-family work is stable.", Status: "held", Dependencies: []string{"86agbv0k9"}},
			{ID: "86agbv0ke", Title: "Pricing and cost enrichment", Description: "Late-stage work only.", Status: "held", Dependencies: []string{"86agbv0kg"}},
			{ID: "86agbv0kh", Title: "v2 architectural paradigms groundwork", Description: "Future-only work after v1 maturity.", Status: "held", Dependencies: []string{"86agbv0kg"}},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", application.handleHealth)
	mux.HandleFunc("/api/bootstrap", application.handleBootstrap)
	mux.HandleFunc("/api/catalog", application.handleCatalog)
	mux.HandleFunc("/api/comparisons/realtime", application.handleRealtimeComparisons)
	mux.HandleFunc("/api/first-wave", application.handleFirstWave)
	mux.HandleFunc("/api/async/demo", application.handleAsyncDemo)
	mux.HandleFunc("/api/async/demo/", application.handleAsyncStatus)
	mux.HandleFunc("/api/events", application.handleEvents)

	handler := withCORS(withLogging(mux))
	server := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("CommLayers API starting on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (a *app) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"service":   "commlayers-api",
		"status":    "ok",
		"startedAt": a.startedAt.Format(time.RFC3339),
		"version":   "0.1.0-bootstrap",
	})
}

func (a *app) handleBootstrap(w http.ResponseWriter, r *http.Request) {
	writeEnvelope(w, http.StatusOK, map[string]any{
		"summary": "Backend-first bootstrap with Go API, staged first-wave execution, and deferred later waves.",
		"runtime": map[string]any{
			"api":            "Go stdlib HTTP",
			"frontend":       "TypeScript + Vite + static nginx serving",
			"mandatoryInfra": []string{"PostgreSQL", "Redis"},
			"optionalInfra":  []string{"Elasticsearch", "Neo4j", "Qdrant", "MongoDB"},
		},
		"auth": map[string]string{
			"identityProvider": "Keycloak (planned)",
			"secretBroker":     "Vault (planned)",
			"currentState":     "bootstrap-only; not implemented in first working version",
		},
	})
}

func (a *app) handleCatalog(w http.ResponseWriter, r *http.Request) {
	families := []map[string]any{
		{"name": "request_response", "status": "bootstrap-ready", "recommendedWhen": []string{"simple_crud", "bounded_request_reply"}, "avoidWhen": []string{"continuous_server_push"}},
		{"name": "polling", "status": "planned", "recommendedWhen": []string{"compatibility_fallback"}, "avoidWhen": []string{"high_frequency_updates"}},
		{"name": "long_polling", "status": "planned", "recommendedWhen": []string{"legacy_push_compatibility"}, "avoidWhen": []string{"modern_push_available"}},
		{"name": "server_sent_events", "status": "bootstrap-demo", "recommendedWhen": []string{"one_way_live_updates"}, "avoidWhen": []string{"duplex_interaction"}},
		{"name": "websocket", "status": "planned", "recommendedWhen": []string{"duplex_low_latency"}, "avoidWhen": []string{"simple_server_push"}},
		{"name": "messaging_and_sync_families", "status": "deferred", "recommendedWhen": []string{"after_first_wave_approval"}, "avoidWhen": []string{"bootstrap_scope"}},
	}

	writeEnvelope(w, http.StatusOK, map[string]any{
		"families": families,
		"notes": []string{
			"This first working version is a bootstrap foundation, not the full platform.",
			"Later families remain held until first-wave artifacts are approved.",
		},
	})
}

func (a *app) handleRealtimeComparisons(w http.ResponseWriter, r *http.Request) {
	comparisons := []comparison{
		{
			Name:           "request_response",
			RecommendedFor: []string{"bounded synchronous reads", "simple commands"},
			AvoidWhen:      []string{"continuous updates"},
			LatencyClass:   "baseline",
			Notes:          []string{"Use as the reference case for all later comparisons."},
		},
		{
			Name:           "server_sent_events",
			RecommendedFor: []string{"one-way updates", "progress streams", "status feeds"},
			AvoidWhen:      []string{"bidirectional collaboration", "binary transport needs"},
			LatencyClass:   "low with server push",
			Notes:          []string{"Default push choice in the current architecture.", "Reverse proxies must disable buffering for timely delivery."},
		},
		{
			Name:           "websocket",
			RecommendedFor: []string{"true duplex interaction", "high-frequency command loops"},
			AvoidWhen:      []string{"simple server-to-client notifications"},
			LatencyClass:   "lowest interactive overhead when justified",
			Notes:          []string{"Use only when duplex interaction is actually needed.", "Some edge platforms deeply inspect only the initial upgrade handshake."},
		},
	}

	writeEnvelope(w, http.StatusOK, map[string]any{
		"comparisons":        comparisons,
		"delayAppliedMs":     650,
		"timeoutThresholdMs": 10000,
	})
}

func (a *app) handleFirstWave(w http.ResponseWriter, r *http.Request) {
	writeEnvelope(w, http.StatusOK, map[string]any{
		"activeStreams":   filterStreams(a.streams, func(s stream) bool { return s.Status == "active" }),
		"deferredStreams": filterStreams(a.streams, func(s stream) bool { return s.Status != "active" }),
		"executionOrder":  []string{"86agbv0k4", "86agbv0k3", "86agbv0k7", "86agbv0k9"},
	})
}

func (a *app) handleAsyncDemo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
		return
	}

	jobID := fmt.Sprintf("job-%d", time.Now().UnixNano())
	now := time.Now().UTC().Format(time.RFC3339)
	status := jobStatus{
		ID:               jobID,
		Status:           "queued",
		Progress:         0,
		CurrentStep:      "queued",
		DelayAppliedMS:   1400,
		TimeoutThreshold: 10000,
		Timeline:         []string{"queued"},
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	a.jobs.Store(jobID, status)
	go a.runJob(context.Background(), jobID)

	w.Header().Set("Location", "/api/async/demo/"+jobID)
	writeEnvelope(w, http.StatusAccepted, map[string]any{
		"jobId":     jobID,
		"statusUrl": "/api/async/demo/" + jobID,
		"state":     status,
	})
}

func (a *app) handleAsyncStatus(w http.ResponseWriter, r *http.Request) {
	jobID := strings.TrimPrefix(r.URL.Path, "/api/async/demo/")
	if jobID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing_job_id"})
		return
	}

	raw, ok := a.jobs.Load(jobID)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "job_not_found"})
		return
	}

	writeEnvelope(w, http.StatusOK, raw)
}

func (a *app) handleEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "streaming_unsupported"})
		return
	}

	events := []map[string]any{
		{"step": "bootstrap_started", "delayAppliedMs": 250},
		{"step": "first_wave_loaded", "delayAppliedMs": 700},
		{"step": "sse_demo_complete", "delayAppliedMs": 1150},
	}

	for _, event := range events {
		payload, _ := json.Marshal(event)
		fmt.Fprintf(w, "event: progress\n")
		fmt.Fprintf(w, "data: %s\n\n", payload)
		flusher.Flush()
		time.Sleep(1200 * time.Millisecond)
	}
}

func (a *app) runJob(ctx context.Context, jobID string) {
	steps := []struct {
		status   string
		progress int
		step     string
	}{
		{status: "running", progress: 25, step: "warming_benchmark_context"},
		{status: "running", progress: 60, step: "collecting_realtime_comparison"},
		{status: "running", progress: 85, step: "publishing_delay_metadata"},
		{status: "succeeded", progress: 100, step: "completed"},
	}

	for _, step := range steps {
		select {
		case <-ctx.Done():
			return
		case <-time.After(1400 * time.Millisecond):
		}

		raw, ok := a.jobs.Load(jobID)
		if !ok {
			return
		}

		state := raw.(jobStatus)
		state.Status = step.status
		state.Progress = step.progress
		state.CurrentStep = step.step
		state.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
		state.Timeline = append(state.Timeline, step.step)
		a.jobs.Store(jobID, state)
	}
}

func filterStreams(streams []stream, keep func(stream) bool) []stream {
	filtered := make([]stream, 0, len(streams))
	for _, stream := range streams {
		if keep(stream) {
			filtered = append(filtered, stream)
		}
	}
	return filtered
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS_ALLOW_ORIGIN"))
		if w.Header().Get("Access-Control-Allow-Origin") == "" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("method=%s path=%s duration=%s", r.Method, r.URL.Path, time.Since(started))
	})
}

func writeEnvelope(w http.ResponseWriter, status int, data any) {
	writeJSON(w, status, responseEnvelope{
		Data:       data,
		Source:     "bootstrap",
		Confidence: "medium",
		Staleness:  "fresh",
		MeasuredAt: time.Now().UTC().Format(time.RFC3339),
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("encode error: %v", err)
	}
}
