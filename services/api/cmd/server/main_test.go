package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"commlayers/services/api/internal/cache"
	"commlayers/services/api/internal/persistence"
)

func TestHandleHealth(t *testing.T) {
	app := &app{
		startedAt: time.Now().UTC(),
	}

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	app.handleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp["service"] != "commlayers-api" {
		t.Errorf("expected service=commlayers-api, got %v", resp["service"])
	}
	if resp["status"] != "ok" {
		t.Errorf("expected status=ok, got %v", resp["status"])
	}
}

func TestHandleCatalog(t *testing.T) {
	store, err := persistence.New("postgres://commlayers:commlayers@localhost:5432/commlayers?sslmode=disable")
	if err != nil {
		t.Skipf("skipping test - database not available: %v", err)
	}

	redisCache, err := cache.New("redis://localhost:6379/0")
	if err != nil {
		t.Skipf("skipping test - redis not available: %v", err)
	}

	app := &app{
		startedAt: time.Now().UTC(),
		store:     store,
		cache:     redisCache,
		families: []family{
			{Name: "test_family", Category: "transport", RecommendedWhen: []string{"test"}, AvoidWhen: []string{"avoid"}},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/api/catalog", nil)
	w := httptest.NewRecorder()

	app.handleCatalog(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp []familyResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(resp) != 1 {
		t.Errorf("expected 1 family, got %d", len(resp))
	}

	fam := resp[0]
	if fam.Family != "test_family" {
		t.Errorf("expected family=test_family, got %s", fam.Family)
	}
	if fam.NormalizedName != "test_family" {
		t.Errorf("expected normalized_name=test_family, got %s", fam.NormalizedName)
	}
	if fam.Latency.P50Ms != 0 {
		t.Errorf("expected p50_ms=0, got %d", fam.Latency.P50Ms)
	}
	if fam.I18n.Locale != "en" {
		t.Errorf("expected locale=en, got %s", fam.I18n.Locale)
	}
	if len(fam.I18n.AvailableLocales) != 2 {
		t.Errorf("expected 2 available locales, got %d", len(fam.I18n.AvailableLocales))
	}
}

func TestHandleTransportSummary(t *testing.T) {
	store, err := persistence.New("postgres://commlayers:commlayers@localhost:5432/commlayers?sslmode=disable")
	if err != nil {
		t.Skipf("skipping test - database not available: %v", err)
	}

	redisCache, err := cache.New("redis://localhost:6379/0")
	if err != nil {
		t.Skipf("skipping test - redis not available: %v", err)
	}

	app := &app{
		startedAt: time.Now().UTC(),
		store:     store,
		cache:     redisCache,
		families: []family{
			{Name: "request_response", Category: "transport", RecommendedWhen: []string{"crud"}, AvoidWhen: []string{"streaming"}},
			{Name: "server_sent_events", Category: "transport", RecommendedWhen: []string{"updates"}, AvoidWhen: []string{"duplex"}},
			{Name: "messaging_workflow", Category: "backend-family", RecommendedWhen: []string{"async"}, AvoidWhen: []string{"sync"}},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/api/transports", nil)
	w := httptest.NewRecorder()

	app.handleTransportSummary(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp []familyResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// Should only return transport families (2)
	if len(resp) != 2 {
		t.Errorf("expected 2 transport families, got %d", len(resp))
	}
}

func TestHandleRealtimeComparisons(t *testing.T) {
	store, err := persistence.New("postgres://commlayers:commlayers@localhost:5432/commlayers?sslmode=disable")
	if err != nil {
		t.Skipf("skipping test - database not available: %v", err)
	}

	redisCache, err := cache.New("redis://localhost:6379/0")
	if err != nil {
		t.Skipf("skipping test - redis not available: %v", err)
	}

	app := &app{
		startedAt: time.Now().UTC(),
		store:     store,
		cache:     redisCache,
		families: []family{
			{Name: "request_response", Category: "transport", RecommendedWhen: []string{"crud"}, AvoidWhen: []string{"streaming"}},
			{Name: "server_sent_events", Category: "transport", RecommendedWhen: []string{"updates"}, AvoidWhen: []string{"duplex"}},
			{Name: "websocket", Category: "transport", RecommendedWhen: []string{"duplex"}, AvoidWhen: []string{"simple"}},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/api/comparisons/realtime", nil)
	w := httptest.NewRecorder()

	app.handleRealtimeComparisons(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp []familyResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(resp) != 3 {
		t.Errorf("expected 3 transport families, got %d", len(resp))
	}

	// Check benchmark latency data is present
	for _, fam := range resp {
		if fam.Latency.P50Ms != 12 {
			t.Errorf("expected p50_ms=12, got %d", fam.Latency.P50Ms)
		}
		if fam.Latency.P95Ms != 45 {
			t.Errorf("expected p95_ms=45, got %d", fam.Latency.P95Ms)
		}
		if fam.Latency.P99Ms != 120 {
			t.Errorf("expected p99_ms=120, got %d", fam.Latency.P99Ms)
		}
		if fam.Evidence.Source != "benchmark_run" {
			t.Errorf("expected evidence source=benchmark_run, got %s", fam.Evidence.Source)
		}
	}
}

func TestLocaleFromRequest(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected string
	}{
		{"default en", "", "en"},
		{"explicit en", "?locale=en", "en"},
		{"explicit es", "?locale=es", "es"},
		{"invalid defaults to en", "?locale=fr", "en"},
		{"case insensitive", "?locale=EN", "en"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test"+tt.query, nil)
			result := localeFromRequest(req)
			if result != tt.expected {
				t.Errorf("localeFromRequest(%q) = %q, want %q", tt.query, result, tt.expected)
			}
		})
	}
}

func TestMsgFunction(t *testing.T) {
	tests := []struct {
		name     string
		locale   string
		key      string
		expected string
	}{
		{"english key exists", "en", "bootstrap.summary", "Backend-first bootstrap with a concrete first-wave contract surface and deferred later waves."},
		{"spanish key exists", "es", "bootstrap.summary", "Bootstrap backend-first con una superficie contractual concreta de primera ola y olas posteriores diferidas."},
		{"fallback to english", "fr", "bootstrap.summary", "Backend-first bootstrap with a concrete first-wave contract surface and deferred later waves."},
		{"missing key returns key", "en", "nonexistent.key", "nonexistent.key"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := msg(tt.locale, tt.key)
			if result != tt.expected {
				t.Errorf("msg(%q, %q) = %q, want %q", tt.locale, tt.key, result, tt.expected)
			}
		})
	}
}

func TestFilterStreams(t *testing.T) {
	streams := []stream{
		{ID: "1", Status: "active", Category: "first-wave"},
		{ID: "2", Status: "deferred", Category: "deferred-wave-2"},
		{ID: "3", Status: "held", Category: "deferred-wave-3"},
		{ID: "4", Status: "active", Category: "first-wave"},
	}

	active := filterStreams(streams, func(s stream) bool { return s.Status == "active" })
	if len(active) != 2 {
		t.Errorf("expected 2 active streams, got %d", len(active))
	}

	deferred := filterStreams(streams, func(s stream) bool { return s.Status == "deferred" })
	if len(deferred) != 1 {
		t.Errorf("expected 1 deferred stream, got %d", len(deferred))
	}

	held := filterStreams(streams, func(s stream) bool { return s.Status == "held" })
	if len(held) != 1 {
		t.Errorf("expected 1 held stream, got %d", len(held))
	}
}

func TestContainsToken(t *testing.T) {
	tokens := []string{"sse", "push", "status", "one-way"}

	if !containsToken(tokens, "sse") {
		t.Error("expected to find 'sse'")
	}
	if !containsToken(tokens, "push") {
		t.Error("expected to find 'push'")
	}
	if containsToken(tokens, "websocket") {
		t.Error("did not expect to find 'websocket'")
	}
}

func TestNextTransportSessionState(t *testing.T) {
	session := transportSession{
		Mode:       "polling",
		Status:     "running",
		Step:       0,
		TotalSteps: 4,
	}

	// First step
	session = nextTransportSessionState(session)
	if session.Step != 1 {
		t.Errorf("expected step=1, got %d", session.Step)
	}
	if session.Status != "running" {
		t.Errorf("expected status=running, got %s", session.Status)
	}

	// Complete all steps
	for i := 0; i < 3; i++ {
		session = nextTransportSessionState(session)
	}
	if session.Step != 4 {
		t.Errorf("expected step=4, got %d", session.Step)
	}
	if session.Status != "completed" {
		t.Errorf("expected status=completed, got %s", session.Status)
	}
}

func TestNextWorkflowState(t *testing.T) {
	run := workflowRun{
		Status:       "queued",
		CurrentStage: "queued",
		Progress:     0,
	}

	run = nextWorkflowState(run)
	if run.CurrentStage != "running" {
		t.Errorf("expected stage=running, got %s", run.CurrentStage)
	}
	if run.Progress != 33 {
		t.Errorf("expected progress=33, got %d", run.Progress)
	}

	run = nextWorkflowState(run)
	if run.CurrentStage != "completed" {
		t.Errorf("expected stage=completed, got %s", run.CurrentStage)
	}
	if run.Progress != 100 {
		t.Errorf("expected progress=100, got %d", run.Progress)
	}
	if run.Status != "completed" {
		t.Errorf("expected status=completed, got %s", run.Status)
	}

	// Further calls should keep it completed
	run = nextWorkflowState(run)
	if run.Status != "completed" {
		t.Errorf("expected status=completed, got %s", run.Status)
	}
	if run.Progress != 100 {
		t.Errorf("expected progress=100, got %d", run.Progress)
	}
}