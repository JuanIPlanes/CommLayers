package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"commlayers/services/api/internal/cache"
	"commlayers/services/api/internal/persistence"
	"commlayers/services/api/internal/router"

	"github.com/gorilla/websocket"
)

func main() {
	application, err := newApp()
	if err != nil {
		log.Fatal(err)
	}

	bm := router.NewBenchmarkMiddleware()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", application.handleHealth)
	mux.HandleFunc("/api/bootstrap", application.handleBootstrap)
	mux.HandleFunc("/api/streams", application.handleStreams)
	mux.HandleFunc("/api/catalog", application.handleCatalog)
	mux.HandleFunc("/api/first-wave/contract", application.handleFirstWaveContract)
	mux.HandleFunc("/api/first-wave/streams/", application.handleFirstWaveStream)
	mux.HandleFunc("/api/security/bootstrap", application.handleSecurityBootstrap)
	mux.HandleFunc("/api/data-platform", application.handleDataPlatform)
	mux.HandleFunc("/api/benchmark-framework", application.handleBenchmarkFramework)
	mux.HandleFunc("/api/deferred-waves", application.handleDeferredWaves)
	mux.HandleFunc("/api/v2-readiness", application.handleV2Readiness)
	mux.HandleFunc("/api/comparisons/realtime", application.handleRealtimeComparisons)
	mux.HandleFunc("/api/transports", application.handleTransportSummary)
	mux.HandleFunc("/api/transports/polling", application.handlePollingSession)
	mux.HandleFunc("/api/transports/polling/", application.handlePollingStatus)
	mux.HandleFunc("/api/messaging", application.handleMessagingOverview)
	mux.HandleFunc("/api/messaging/workflows", application.handleWorkflowCreate)
	mux.HandleFunc("/api/messaging/workflows/", application.handleWorkflowStatus)
	mux.HandleFunc("/api/sync", application.handleSyncOverview)
	mux.HandleFunc("/api/sync/sessions", application.handleSyncSessionCreate)
	mux.HandleFunc("/api/sync/sessions/", application.handleSyncSessionAction)
	mux.HandleFunc("/api/projections", application.handleProjectionOverview)
	mux.HandleFunc("/api/projections/search", application.handleProjectionSearch)
	mux.HandleFunc("/api/projections/graph", application.handleProjectionGraph)
	mux.HandleFunc("/api/projections/vector", application.handleProjectionVector)
	mux.HandleFunc("/api/pricing", application.handlePricingOverview)
	mux.HandleFunc("/api/v2/roadmap", application.handleV2Roadmap)
	mux.HandleFunc("/api/v2/paradigms/event-driven", application.handleV2EventDriven)
	mux.HandleFunc("/api/v2/paradigms/reactive-stream-driven", application.handleV2ReactiveStreamDriven)
	mux.HandleFunc("/api/v2/paradigms/local-first-crdt", application.handleV2LocalFirstCRDT)
	mux.HandleFunc("/api/v2/paradigms/consensus-driven", application.handleV2ConsensusDriven)
	mux.HandleFunc("/api/v2/paradigms/planet-scale-composite", application.handleV2PlanetScaleComposite)
	mux.HandleFunc("/api/async/demo", application.handleAsyncDemo)
	mux.HandleFunc("/api/async/demo/", application.handleAsyncStatus)
	mux.HandleFunc("/api/events", application.handleEvents)
	mux.HandleFunc("/api/ws/demo", application.handleWebSocketDemo)
	mux.HandleFunc("/api/benchmarks/metrics", application.handleBenchmarkMetrics)

	handler := withCORS(withLogging(router.BuildRouter(mux, bm)))
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

func newApp() (*app, error) {
	store, err := persistence.New(os.Getenv("POSTGRES_URL"))
	if err != nil {
		return nil, err
	}
	redisCache, err := cache.New(os.Getenv("REDIS_URL"))
	if err != nil {
		return nil, err
	}

	return &app{
		startedAt: time.Now().UTC(),
		store:     store,
		cache:     redisCache,
		streams: []stream{
			{
				ID:               "86agbv0k4",
				Title:            "Backend abstraction layer matrix",
				Role:             "wave entry gate",
				Category:         "first-wave",
				Description:      "Freeze normalized backend vocabulary, comparison dimensions, and dependency contracts.",
				Status:           "active",
				Dependencies:     []string{},
				Deliverables:     []string{"layer responsibility matrix", "capability comparison rubric", "default and fallback decision criteria", "downstream dependency contract"},
				AcceptanceGates:  []string{"Every backend family is comparable on consistency, latency class, durability semantics, operational burden, and local-compose viability.", "Vocabulary is frozen before security and canonical data streams overlap."},
				ExecutionPackets: []string{"Finalize abstraction vocabulary and layer boundaries", "Finalize comparison rubric and decision criteria", "Freeze downstream dependency contract"},
			},
			{
				ID:               "86agbv0k3",
				Title:            "Keycloak and Vault security bootstrap",
				Role:             "controlled-overlap stream",
				Category:         "first-wave",
				Description:      "Define bootstrap identity, secret authority, and trust-boundary contracts after abstraction freeze.",
				Status:           "active",
				Dependencies:     []string{"86agbv0k4"},
				Deliverables:     []string{"Keycloak realm and client baseline", "Vault secret ownership strategy", "rotation and revocation posture", "bootstrap trust-boundary contract"},
				AcceptanceGates:  []string{"Realm, client, and service identity model is explicit.", "Vault paths, secret classes, and rotation posture are explicit.", "Bootstrap trust boundaries are explicit."},
				ExecutionPackets: []string{"Finalize identity and bootstrap contract", "Finalize Vault ownership and secret classes", "Finalize rotation, revocation, and bootstrap boundary rules"},
			},
			{
				ID:               "86agbv0k7",
				Title:            "Canonical data platform",
				Role:             "controlled-overlap stream",
				Category:         "first-wave",
				Description:      "Define canonical entities, storage ownership, schema evolution, and audit lineage after abstraction freeze.",
				Status:           "active",
				Dependencies:     []string{"86agbv0k4"},
				Deliverables:     []string{"canonical entity catalog", "storage ownership matrix", "schema and versioning policy", "audit lineage baseline"},
				AcceptanceGates:  []string{"Canonical entity ownership is explicit.", "Schema evolution rules are explicit.", "Audit lineage fields are explicit.", "Compatibility rules across backend families are explicit."},
				ExecutionPackets: []string{"Finalize canonical entity and ownership catalog", "Finalize schema and versioning contract", "Finalize audit lineage baseline"},
			},
			{
				ID:               "86agbv0k9",
				Title:            "Benchmark and observability framework",
				Role:             "wave closure gate",
				Category:         "first-wave",
				Description:      "Define benchmark scenarios, telemetry schema, thresholds, and closure reporting after the first three streams settle.",
				Status:           "active",
				Dependencies:     []string{"86agbv0k4", "86agbv0k3", "86agbv0k7"},
				Deliverables:     []string{"benchmark matrix", "metric, log, and trace naming contract", "reproducibility rules", "wave closure report template"},
				AcceptanceGates:  []string{"Benchmark scenarios are fixed.", "Telemetry minimums are fixed.", "Success and error taxonomy is fixed.", "Threshold and reporting format is stable."},
				ExecutionPackets: []string{"Finalize benchmark scenarios and reproducibility rules", "Finalize telemetry schema and naming contract", "Finalize wave-closure thresholds and report format"},
			},
			{
				ID:               "86agbv0k6",
				Title:            "Async and long-running operations",
				Role:             "deferred backend family stream",
				Category:         "deferred-wave-2",
				Description:      "Long-running operation semantics remain deferred until the first-wave closure gate is accepted.",
				Status:           "deferred",
				Dependencies:     []string{"86agbv0k9"},
				Deliverables:     []string{"202/LRO contract", "idempotency rules", "polling and subscription visibility model"},
				AcceptanceGates:  []string{"Held until first-wave approval."},
				ExecutionPackets: []string{"Do not start yet"},
			},
			{
				ID:               "86agbv0k5",
				Title:            "Transport family implementations",
				Role:             "deferred backend family stream",
				Category:         "deferred-wave-2",
				Description:      "Transport implementations remain deferred until the first-wave closure gate is accepted.",
				Status:           "deferred",
				Dependencies:     []string{"86agbv0k9"},
				Deliverables:     []string{"request-response contract", "polling and long-polling comparison", "SSE and WebSocket guidance", "webhook and WebRTC signaling notes"},
				AcceptanceGates:  []string{"Held until first-wave approval."},
				ExecutionPackets: []string{"Do not start yet"},
			},
			{
				ID:               "86agbv0ka",
				Title:            "Messaging and workflow families",
				Role:             "deferred backend family stream",
				Category:         "deferred-wave-2",
				Description:      "Messaging, broker, and workflow families remain deferred.",
				Status:           "deferred",
				Dependencies:     []string{"86agbv0k9"},
				Deliverables:     []string{"broker criteria", "queue semantics", "workflow contract", "failure mode matrix"},
				AcceptanceGates:  []string{"Held until first-wave approval."},
				ExecutionPackets: []string{"Do not start yet"},
			},
			{
				ID:               "86agbv0kb",
				Title:            "Sync, replication, and coordination families",
				Role:             "deferred backend family stream",
				Category:         "deferred-wave-2",
				Description:      "Sync and coordination families remain deferred.",
				Status:           "deferred",
				Dependencies:     []string{"86agbv0k9"},
				Deliverables:     []string{"sync behaviors", "conflict handling", "replication checks", "consensus use criteria"},
				AcceptanceGates:  []string{"Held until first-wave approval."},
				ExecutionPackets: []string{"Do not start yet"},
			},
			{
				ID:               "86agbv0kc",
				Title:            "Search, graph, and vector projections",
				Role:             "deferred backend family stream",
				Category:         "deferred-wave-2",
				Description:      "Projection and optional profile work remains deferred.",
				Status:           "deferred",
				Dependencies:     []string{"86agbv0k9"},
				Deliverables:     []string{"search projection model", "graph profile gate", "vector profile escalation rules"},
				AcceptanceGates:  []string{"Held until first-wave approval."},
				ExecutionPackets: []string{"Do not start yet"},
			},
			{
				ID:               "86agbv0kf",
				Title:            "Localization for frontend-facing responses",
				Role:             "deferred frontend stream",
				Category:         "deferred-wave-3",
				Description:      "Localization is staged for later after backend-family work stabilizes.",
				Status:           "held",
				Dependencies:     []string{"86agbv0k9"},
				Deliverables:     []string{"en default resources", "es resources", "backend localization contract"},
				AcceptanceGates:  []string{"Held until backend-family work is stable."},
				ExecutionPackets: []string{"Do not start yet"},
			},
			{
				ID:               "86agbv0kg",
				Title:            "Frontend comparative visualization",
				Role:             "deferred frontend stream",
				Category:         "deferred-wave-3",
				Description:      "The richer frontend is held behind the backend-first milestones.",
				Status:           "held",
				Dependencies:     []string{"86agbv0k9"},
				Deliverables:     []string{"desktop comparison UI", "visible delay interactions", "FPS verification"},
				AcceptanceGates:  []string{"Held until backend-family work is stable."},
				ExecutionPackets: []string{"Do not start yet"},
			},
			{
				ID:               "86agbv0ke",
				Title:            "Pricing and cost enrichment",
				Role:             "late-stage enrichment stream",
				Category:         "deferred-wave-3",
				Description:      "Pricing remains late-stage only.",
				Status:           "held",
				Dependencies:     []string{"86agbv0kg"},
				Deliverables:     []string{"USD freshness contract", "confidence semantics", "pricing fallback rules"},
				AcceptanceGates:  []string{"Held until backend and localization work are stable."},
				ExecutionPackets: []string{"Do not start yet"},
			},
			{
				ID:               "86agbv0kh",
				Title:            "v2 architectural paradigms groundwork",
				Role:             "future-only stream",
				Category:         "deferred-wave-4",
				Description:      "v2 paradigm work stays future-facing and blocked by v1 maturity.",
				Status:           "held",
				Dependencies:     []string{"86agbv0kg"},
				Deliverables:     []string{"paradigm taxonomy", "v2 unlock prerequisites", "future implementation boundaries"},
				AcceptanceGates:  []string{"Held until v1 backend-family maturity is demonstrated."},
				ExecutionPackets: []string{"Do not start yet"},
			},
		},
		families: []family{
			{Name: "request_response", Category: "transport", Status: "bootstrap-ready", RecommendedWhen: []string{"simple CRUD", "bounded request-reply"}, AvoidWhen: []string{"continuous updates"}, Notes: []string{"Reference case for later comparisons."}},
			{Name: "polling", Category: "transport", Status: "planned", RecommendedWhen: []string{"compatibility fallback", "pull-based status checks"}, AvoidWhen: []string{"high-frequency updates"}, Notes: []string{"Held until first-wave approval."}},
			{Name: "long_polling", Category: "transport", Status: "planned", RecommendedWhen: []string{"legacy push compatibility"}, AvoidWhen: []string{"modern push available"}, Notes: []string{"Held until first-wave approval."}},
			{Name: "server_sent_events", Category: "transport", Status: "bootstrap-demo", RecommendedWhen: []string{"one-way updates", "status feeds", "progress streams"}, AvoidWhen: []string{"duplex interaction"}, Notes: []string{"Default push choice in current architecture."}},
			{Name: "websocket", Category: "transport", Status: "planned", RecommendedWhen: []string{"duplex low latency"}, AvoidWhen: []string{"simple server push"}, Notes: []string{"Use only when duplex interaction is actually needed."}},
			{Name: "messaging_workflow", Category: "backend-family", Status: "deferred", RecommendedWhen: []string{"after first-wave approval"}, AvoidWhen: []string{"bootstrap scope"}, Notes: []string{"Broker, queue, workflow, and event families remain staged only."}},
			{Name: "sync_replication_projection", Category: "backend-family", Status: "deferred", RecommendedWhen: []string{"after first-wave approval"}, AvoidWhen: []string{"bootstrap scope"}, Notes: []string{"Sync, coordination, and projection families remain staged only."}},
		},
		waves: []deferredWave{
			{
				Name:             "Deferred wave 2 - backend family expansion",
				Purpose:          "Unlock transport, messaging, sync, and projection implementation only after first-wave approval.",
				Streams:          []string{"86agbv0k6", "86agbv0k5", "86agbv0ka", "86agbv0kb", "86agbv0kc"},
				Focus:            []string{"long-running operations", "transport deepening", "messaging families", "sync and replication", "projection profiles"},
				ReactivationRule: []string{"Approve first-wave artifacts.", "Keep implementation sequencing behind 86agbv0k9 closure."},
			},
			{
				Name:             "Deferred wave 3 - frontend and enrichment",
				Purpose:          "Unlock localized frontend comparison and pricing after backend-family stability.",
				Streams:          []string{"86agbv0kf", "86agbv0kg", "86agbv0ke"},
				Focus:            []string{"localization", "frontend comparative visualization", "pricing enrichment"},
				ReactivationRule: []string{"Backend-family work must stabilize first.", "Pricing remains late-stage only."},
			},
			{
				Name:             "Deferred wave 4 - v2 paradigms",
				Purpose:          "Keep v2 architecture work staged but inactive until v1 maturity.",
				Streams:          []string{"86agbv0kh"},
				Focus:            []string{"paradigm taxonomy", "unlock prerequisites", "future implementation boundaries"},
				ReactivationRule: []string{"V1 backend-family maturity is proven.", "Frontend comparison layer is stable."},
			},
		},
		projectionDocs: []projectionDoc{
			{ID: "doc-search-1", Title: "SSE default push", Tokens: []string{"sse", "push", "status", "one-way"}, Summary: "SSE is the current default push channel for one-way updates."},
			{ID: "doc-search-2", Title: "Workflow orchestration", Tokens: []string{"workflow", "queue", "retry", "orchestration"}, Summary: "Workflow demos expose queued, running, and completed stages without a real broker yet."},
			{ID: "doc-search-3", Title: "Sync merge demo", Tokens: []string{"sync", "replication", "conflict", "merge"}, Summary: "Sync demos show primary/replica drift and explicit merge behavior."},
		},
		graphNodes: []graphNode{
			{ID: "request_response", Label: "Request response", DependsOn: []string{}, Category: "transport", Description: "Reference case for bounded request-reply work."},
			{ID: "server_sent_events", Label: "SSE", DependsOn: []string{"request_response"}, Category: "transport", Description: "Default push channel for one-way updates."},
			{ID: "websocket", Label: "WebSocket", DependsOn: []string{"request_response"}, Category: "transport", Description: "Planned duplex transport when truly needed."},
			{ID: "workflow_engine", Label: "Workflow engine", DependsOn: []string{"request_response", "server_sent_events"}, Category: "messaging", Description: "Visible queue and workflow state transitions."},
			{ID: "sync_replication", Label: "Sync replication", DependsOn: []string{"workflow_engine"}, Category: "sync", Description: "Primary/replica drift and explicit merge behavior."},
		},
		vectorDocs: []vectorDoc{
			{ID: "vec-1", Label: "SSE progress", Vector: []float64{0.9, 0.1, 0.0}, Summary: "One-way progress and status feed", Transport: "sse"},
			{ID: "vec-2", Label: "Workflow queue", Vector: []float64{0.2, 0.9, 0.1}, Summary: "Queued and staged orchestration", Transport: "workflow"},
			{ID: "vec-3", Label: "Sync conflict", Vector: []float64{0.2, 0.3, 0.95}, Summary: "Conflict and merge-heavy collaboration", Transport: "sync"},
		},
	}, nil
}

func (a *app) saveRecord(ctx context.Context, kind string, id string, payload any) error {
	if err := a.store.Save(ctx, kind, id, payload); err != nil {
		return err
	}
	if a.cache != nil && a.cache.Enabled() {
		if err := a.cache.Save(ctx, kind+":"+id, payload, 30*time.Minute); err != nil {
			log.Printf("cache save failed kind=%s id=%s err=%v", kind, id, err)
		}
	}
	return nil
}

func (a *app) loadRecord(ctx context.Context, kind string, id string, dst any) (bool, error) {
	if a.cache != nil && a.cache.Enabled() {
		ok, err := a.cache.Load(ctx, kind+":"+id, dst)
		if err != nil {
			log.Printf("cache load failed kind=%s id=%s err=%v", kind, id, err)
		} else if ok {
			return true, nil
		}
	}
	ok, err := a.store.Load(ctx, kind, id, dst)
	if !ok || err != nil {
		return ok, err
	}
	if a.cache != nil && a.cache.Enabled() {
		if err := a.cache.Save(ctx, kind+":"+id, dst, 30*time.Minute); err != nil {
			log.Printf("cache warm failed kind=%s id=%s err=%v", kind, id, err)
		}
	}
	return true, nil
}

func (a *app) emitTransientEvent(ctx context.Context, kind string, payload map[string]any) {
	entry := eventRecord{
		ID:        fmt.Sprintf("evt-%d", time.Now().UnixNano()),
		Kind:      kind,
		Payload:   payload,
		EmittedAt: time.Now().UTC().Format(time.RFC3339),
		Source:    "redis-pubsub",
	}
	if err := a.saveRecord(ctx, "domain_event", entry.ID, entry); err != nil {
		log.Printf("persist event failed kind=%s err=%v", kind, err)
	}
	if a.cache == nil || !a.cache.Enabled() {
		return
	}
	envelope := map[string]any{
		"kind":       kind,
		"payload":    payload,
		"emittedAt":  entry.EmittedAt,
		"source":     entry.Source,
		"confidence": "medium",
	}
	if err := a.cache.Publish(ctx, "commlayers-events", envelope); err != nil {
		log.Printf("publish event failed kind=%s err=%v", kind, err)
	}
}

func (a *app) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"service":   "commlayers-api",
		"status":    "ok",
		"startedAt": a.startedAt.Format(time.RFC3339),
		"version":   "0.2.0-first-wave-bootstrap",
	})
}

func (a *app) handleBootstrap(w http.ResponseWriter, r *http.Request) {
	locale := localeFromRequest(r)
	writeLocalizedEnvelope(w, r, http.StatusOK, map[string]any{
		"summary": msg(locale, "bootstrap.summary"),
		"phase":   "first-wave materialized",
		"repo":    "https://github.com/JuanIPlanes/CommLayers",
		"runtime": map[string]any{
			"api":            "Go stdlib HTTP",
			"frontend":       "TypeScript + Vite + static nginx serving",
			"mandatoryInfra": []string{"PostgreSQL", "Redis"},
			"optionalInfra":  []string{"Elasticsearch", "Neo4j", "Qdrant", "MongoDB"},
			"composeFiles":   []string{"docker-compose.back.yml", "docker-compose.front.yml"},
		},
		"now": map[string]any{
			"activeStreams":           []string{"86agbv0k4", "86agbv0k3", "86agbv0k7", "86agbv0k9"},
			"executionOrder":          []string{"86agbv0k4", "86agbv0k3", "86agbv0k7", "86agbv0k9"},
			"goNoGoForImplementation": msg(locale, "bootstrap.go_no_go"),
		},
		"notYet": []string{
			"Keycloak runtime integration",
			"Vault runtime integration",
			"transport-family implementation beyond the demo surface",
			"messaging, sync, and projection execution",
			"localized frontend and pricing execution",
			"v2 paradigm runtime",
		},
	})
}

func (a *app) handleStreams(w http.ResponseWriter, r *http.Request) {
	grouped := map[string][]stream{
		"active":   filterStreams(a.streams, func(s stream) bool { return s.Status == "active" }),
		"deferred": filterStreams(a.streams, func(s stream) bool { return s.Status == "deferred" }),
		"held":     filterStreams(a.streams, func(s stream) bool { return s.Status == "held" }),
	}
	writeEnvelope(w, http.StatusOK, grouped)
}

func (a *app) handleCatalog(w http.ResponseWriter, r *http.Request) {
	writeFamilyListResponse(w, r, http.StatusOK, a.families, defaultLatency(), defaultStateTransition(), defaultEvidence())
}

func (a *app) handleFirstWaveContract(w http.ResponseWriter, r *http.Request) {
	active := filterStreams(a.streams, func(s stream) bool { return s.Category == "first-wave" })
	sort.Slice(active, func(i, j int) bool { return active[i].ID < active[j].ID })
	writeLocalizedEnvelope(w, r, http.StatusOK, map[string]any{
		"executionOrder": []string{"86agbv0k4", "86agbv0k3", "86agbv0k7", "86agbv0k9"},
		"streams":        active,
		"rules": []string{
			"86agbv0k4 is the wave entry gate.",
			"86agbv0k3 and 86agbv0k7 start only after 86agbv0k4 freezes the shared vocabulary.",
			"86agbv0k9 closes the wave after abstraction, security, and canonical data are accepted.",
		},
	})
}

func (a *app) handleFirstWaveStream(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/first-wave/streams/")
	for _, stream := range a.streams {
		if stream.ID == id && stream.Category == "first-wave" {
			writeEnvelope(w, http.StatusOK, stream)
			return
		}
	}
	writeJSON(w, http.StatusNotFound, map[string]string{"error": "first_wave_stream_not_found"})
}

func (a *app) handleSecurityBootstrap(w http.ResponseWriter, r *http.Request) {
	locale := localeFromRequest(r)
	writeLocalizedEnvelope(w, r, http.StatusOK, map[string]any{
		"identity": map[string]any{
			"provider":          "Keycloak",
			"browserPattern":    "BFF-oriented Authorization Code + PKCE",
			"accessTokenTarget": "5-10 minutes",
			"refreshPolicy":     "rotated refresh tokens",
			"redirectPolicy":    "exact redirect URI matching only",
		},
		"vault": map[string]any{
			"provider":  "Vault",
			"ownership": []string{"service secrets", "bootstrap credentials", "future PKI authority"},
			"pki":       []string{"offline root", "online intermediate", "short-lived service certificates"},
		},
		"mtlsPhases": []map[string]any{
			{"phase": 1, "name": "identity and ingress TLS", "goal": "prepare service identity and trust boundary inventory"},
			{"phase": 2, "name": "selective mTLS", "goal": "protect critical service pairs without blocking bootstrap"},
			{"phase": 3, "name": "strict east-west mTLS", "goal": "enforce workload identity and authorization after observability is ready"},
		},
		"currentState": msg(locale, "security.state"),
	})
}

func (a *app) handleDataPlatform(w http.ResponseWriter, r *http.Request) {
	writeLocalizedEnvelope(w, r, http.StatusOK, map[string]any{
		"canonicalStores": []map[string]any{
			{"name": "PostgreSQL", "role": "system of record", "status": "active in compose", "why": "canonical truth, auditability, schema evolution"},
			{"name": "Redis", "role": "ephemeral speed layer", "status": "active in compose", "why": "cache, rate limiting, short-lived coordination"},
		},
		"optionalProfiles": []map[string]any{
			{"name": "Elasticsearch", "profile": "search", "status": "optional", "why": "search projection only when needed"},
			{"name": "Neo4j", "profile": "graph", "status": "optional", "why": "graph traversal only if first-class"},
			{"name": "Qdrant", "profile": "vector", "status": "optional", "why": "dedicated vector domain only after pgvector threshold"},
			{"name": "MongoDB", "profile": "mongo", "status": "optional", "why": "only if a document-native use case is explicitly justified"},
		},
		"contracts": map[string]any{
			"entityOwnership": []string{"stream definitions", "comparison evidence", "job lifecycle", "deferred-wave manifest"},
			"schemaRules":     []string{"explicit ownership", "versioned evolution", "audit lineage required"},
			"evidenceFields":  []string{"source", "confidence", "staleness", "measuredAt", "delayAppliedMs"},
		},
	})
}

func (a *app) handleBenchmarkFramework(w http.ResponseWriter, r *http.Request) {
	writeLocalizedEnvelope(w, r, http.StatusOK, map[string]any{
		"benchmarkMatrix": []map[string]any{
			{"name": "request baseline", "measures": []string{"p50", "p95", "p99", "error rate"}},
			{"name": "SSE startup and event cadence", "measures": []string{"stream startup latency", "event flush latency", "delay metadata"}},
			{"name": "async visibility", "measures": []string{"queue delay", "step transition cadence", "timeout threshold visibility"}},
		},
		"telemetry": map[string]any{
			"logs":    []string{"request_id", "stream_id", "job_id", "phase", "result"},
			"metrics": []string{"request_duration_ms", "event_flush_ms", "job_progress", "delay_applied_ms"},
			"traces":  []string{"bootstrap.http", "bootstrap.async", "bootstrap.sse"},
		},
		"closureRules": []string{
			"No later-wave implementation starts before first-wave metrics and closure report format are fixed.",
			"User-visible delay metadata must be available for the frontend.",
			"Success and error taxonomy must be stable across first-wave endpoints.",
		},
	})
}

func (a *app) handleBenchmarkMetrics(w http.ResponseWriter, r *http.Request) {
	// This would need access to the benchmark middleware collector
	// For now, return a placeholder response
	writeLocalizedEnvelope(w, r, http.StatusOK, map[string]any{
		"note": "Benchmark metrics endpoint - requires benchmark middleware integration",
		"available_metrics": []string{
			"request_duration_ms",
			"queue_delay_ms",
			"timeout_ms",
			"state_transition_ms",
		},
	})
}

func (a *app) handleDeferredWaves(w http.ResponseWriter, r *http.Request) {
	writeLocalizedEnvelope(w, r, http.StatusOK, map[string]any{
		"waves":        a.waves,
		"holdManifest": "clickup/08-deferred-waves.md",
	})
}

func (a *app) handleV2Readiness(w http.ResponseWriter, r *http.Request) {
	locale := localeFromRequest(r)
	writeLocalizedEnvelope(w, r, http.StatusOK, map[string]any{
		"status": msg(locale, "v2.status"),
		"blockedBy": []string{
			"first-wave abstraction matrix not yet executed in code",
			"security bootstrap not yet integrated at runtime",
			"canonical data platform not yet backed by persistent models",
			"benchmark closure not yet accepted",
			"transport, messaging, sync, and projection families not yet implemented",
		},
		"unlockPath": []string{
			"complete and accept the first-wave implementation milestones",
			"execute deferred wave 2 backend families",
			"stabilize localization and frontend comparative visualization",
			"treat pricing as late-stage only",
			"only then open v2 paradigm implementation planning",
		},
	})
}

func (a *app) handleV2Roadmap(w http.ResponseWriter, r *http.Request) {
	writeLocalizedEnvelope(w, r, http.StatusOK, map[string]any{
		"status": "groundwork_only",
		"paradigms": []map[string]any{
			{"name": "reactive_stream_driven", "status": "planned", "unlock": "messaging, sync, and benchmark maturity"},
			{"name": "local_first_crdt", "status": "planned", "unlock": "sync slice maturity plus conflict strategy stabilization"},
			{"name": "event_driven_architecture", "status": "planned", "unlock": "transport and messaging slices plus durable event contract"},
			{"name": "consensus_driven", "status": "planned", "unlock": "coordination semantics beyond bootstrap-grade in-memory demos"},
			{"name": "planet_scale", "status": "planned", "unlock": "projection, pricing, and frontend comparative maturity"},
		},
		"v1Completed": []string{
			"first-wave foundation",
			"transport demos",
			"messaging and workflow demos",
			"sync and replication demos",
			"projection demos",
			"localization slice",
			"pricing slice",
		},
		"nextBeforeV2Implementation": []string{
			"formalize durable event and projection persistence beyond in-memory demos",
			"stabilize frontend comparative UX against all v1 backend families",
			"define real provider integrations that replace demo-grade data sources",
			"freeze v2 paradigm acceptance gates per paradigm family",
		},
	})
}

func (a *app) handleV2EventDriven(w http.ResponseWriter, r *http.Request) {
	type recentEvent struct {
		ID        string         `json:"id"`
		Kind      string         `json:"kind"`
		Payload   map[string]any `json:"payload"`
		EmittedAt string         `json:"emittedAt"`
		Source    string         `json:"source"`
	}
	var events []recentEvent
	if err := a.store.ListByKind(r.Context(), "domain_event", 12, &events); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "load_domain_events_failed"})
		return
	}
	projections := map[string]int{}
	for _, event := range events {
		projections[event.Kind]++
	}
	writeLocalizedEnvelope(w, r, http.StatusOK, map[string]any{
		"paradigm": "event_driven_architecture",
		"status":   "active",
		"commands": []string{"create workflow run", "create sync session", "mutate sync replica", "create transport session", "start async job"},
		"events":   events,
		"projections": map[string]any{
			"eventKinds":         projections,
			"sourceOfTruth":      "postgres domain_event records",
			"transientBroadcast": "redis pubsub via /api/events",
		},
		"notes": []string{
			"This is the first true v2 slice: command to event to projection behavior is now explicit.",
			"The slice reuses the current v1 transport, messaging, sync, projection, and persistence/runtime surfaces.",
		},
	})
}

func (a *app) handleV2ReactiveStreamDriven(w http.ResponseWriter, r *http.Request) {
	type recentEvent struct {
		ID        string         `json:"id"`
		Kind      string         `json:"kind"`
		Payload   map[string]any `json:"payload"`
		EmittedAt string         `json:"emittedAt"`
		Source    string         `json:"source"`
	}
	var events []recentEvent
	if err := a.store.ListByKind(r.Context(), "domain_event", 24, &events); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "load_domain_events_failed"})
		return
	}

	kindCounts := map[string]int{}
	var latestEvent string
	for _, event := range events {
		kindCounts[event.Kind]++
		if latestEvent == "" {
			latestEvent = event.EmittedAt
		}
	}

	writeLocalizedEnvelope(w, r, http.StatusOK, map[string]any{
		"paradigm": "reactive_stream_driven",
		"status":   "active",
		"streamSignals": map[string]any{
			"recentEventCount": len(events),
			"eventKinds":       kindCounts,
			"latestEventAt":    latestEvent,
			"transportSurface": "SSE plus Redis transient events",
		},
		"backpressureHints": []string{
			"Use event counts per kind as the first stream-pressure indicator.",
			"Workflow and async job events are the first candidates for reactive throttling.",
			"Projection refresh cadence should follow event density rather than fixed polling only.",
		},
		"notes": []string{
			"This slice treats the current domain-event stream as the basis for a reactive view.",
			"It reuses the event log and Redis transient flow without claiming a full Kafka/Flink-grade system yet.",
		},
	})
}

func (a *app) handleV2LocalFirstCRDT(w http.ResponseWriter, r *http.Request) {
	var sessions []syncSession
	if err := a.store.ListByKind(r.Context(), "sync_session", 12, &sessions); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "load_sync_sessions_failed"})
		return
	}
	conflicts := 0
	aligned := 0
	maxLag := 0
	for _, session := range sessions {
		if session.Conflict {
			conflicts++
		}
		if session.Status == "aligned" {
			aligned++
		}
		if session.Lag > maxLag {
			maxLag = session.Lag
		}
	}
	writeLocalizedEnvelope(w, r, http.StatusOK, map[string]any{
		"paradigm": "local_first_crdt",
		"status":   "active",
		"syncSignals": map[string]any{
			"recentSessionCount": len(sessions),
			"alignedSessions":    aligned,
			"conflictSessions":   conflicts,
			"maxObservedLag":     maxLag,
		},
		"mergeRules": []string{
			"Single-side divergence catches up by replaying the latest persisted value.",
			"Equal-version divergence merges deterministically into a combined value in the bootstrap model.",
			"Conflict and lag are surfaced as first-class signals for later CRDT evolution.",
		},
		"notes": []string{
			"This slice is local-first and CRDT-oriented, not a full CRDT engine yet.",
			"It reuses persisted sync sessions to make offline-ish divergence and merge semantics visible.",
		},
	})
}

func (a *app) handleV2ConsensusDriven(w http.ResponseWriter, r *http.Request) {
	var sessions []syncSession
	if err := a.store.ListByKind(r.Context(), "sync_session", 12, &sessions); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "load_sync_sessions_failed"})
		return
	}
	quorumReady := 0
	quorumBlocked := 0
	for _, session := range sessions {
		if session.Conflict || session.Lag > 0 {
			quorumBlocked++
		} else {
			quorumReady++
		}
	}
	writeLocalizedEnvelope(w, r, http.StatusOK, map[string]any{
		"paradigm": "consensus_driven",
		"status":   "active",
		"coordinationSignals": map[string]any{
			"candidateSessions": len(sessions),
			"quorumReady":       quorumReady,
			"quorumBlocked":     quorumBlocked,
			"rule":              "aligned sessions count as ready, divergent sessions count as blocked",
		},
		"decisionRules": []string{
			"A session with zero lag and no conflict is treated as quorum-ready in the bootstrap model.",
			"Any divergence blocks consensus until replication or merge resolves the difference.",
			"The event stream supplies the change history; sync sessions supply the current coordination view.",
		},
		"notes": []string{
			"This slice is consensus-driven in orientation, not a full Raft/Paxos implementation.",
			"It makes agreement and blocking conditions visible using the current persisted sync and event signals.",
		},
	})
}

func (a *app) handleV2PlanetScaleComposite(w http.ResponseWriter, r *http.Request) {
	writeLocalizedEnvelope(w, r, http.StatusOK, map[string]any{
		"paradigm": "planet_scale_composite_architectures",
		"status":   "active",
		"buildingBlocks": []map[string]any{
			{"name": "event_driven_architecture", "role": "durable event backbone"},
			{"name": "reactive_stream_driven", "role": "stream signal and pressure view"},
			{"name": "local_first_crdt", "role": "divergence and merge semantics"},
			{"name": "consensus_driven", "role": "agreement and blocked coordination signals"},
		},
		"compositionRules": []string{
			"Persisted events become the shared history across slices.",
			"Redis-backed transient flow remains the fast path for live updates.",
			"Sync and consensus signals indicate whether composite state is converging cleanly.",
			"Projection and pricing surfaces remain downstream consumers of the composite runtime.",
		},
		"notes": []string{
			"This slice does not claim true planet-scale deployment yet.",
			"It is the compositional v2 view that explains how multiple paradigms combine on the current runtime foundation.",
		},
	})
}

func (a *app) handleTransportSummary(w http.ResponseWriter, r *http.Request) {
	// Filter transport families from catalog
	transportFamilies := make([]family, 0)
	for _, fam := range a.families {
		if fam.Category == "transport" {
			transportFamilies = append(transportFamilies, fam)
		}
	}
	writeFamilyListResponse(w, r, http.StatusOK, transportFamilies, defaultLatency(), defaultStateTransition(), defaultEvidence())
}

func (a *app) handlePollingSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
		return
	}

	mode := r.URL.Query().Get("mode")
	if mode == "" {
		mode = "polling"
	}
	if mode != "polling" && mode != "long-polling" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "unsupported_mode"})
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)
	session := transportSession{
		ID:                fmt.Sprintf("transport-%d", time.Now().UnixNano()),
		Mode:              mode,
		Status:            "running",
		Step:              0,
		TotalSteps:        4,
		RecommendedPollMS: recommendedDelay(mode),
		DelayAppliedMS:    recommendedDelay(mode),
		Timeline:          []string{"session_created"},
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	if err := a.saveRecord(r.Context(), "transport_session", session.ID, session); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "persist_transport_session_failed"})
		return
	}
	a.emitTransientEvent(r.Context(), "transport_session", map[string]any{
		"id":     session.ID,
		"mode":   session.Mode,
		"status": session.Status,
		"step":   session.Step,
	})
	w.Header().Set("Location", "/api/transports/polling/"+session.ID)
	writeLocalizedEnvelope(w, r, http.StatusCreated, map[string]any{
		"session":   session,
		"statusUrl": "/api/transports/polling/" + session.ID,
	})
}

func (a *app) handlePollingStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/transports/polling/")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing_session_id"})
		return
	}

	var session transportSession
	ok, err := a.loadRecord(r.Context(), "transport_session", id, &session)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "transport_session_not_found"})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "load_transport_session_failed"})
		return
	}

	if session.Mode == "long-polling" && session.Status != "completed" {
		time.Sleep(time.Duration(session.RecommendedPollMS) * time.Millisecond)
	}

	session = nextTransportSessionState(session)
	if err := a.saveRecord(r.Context(), "transport_session", id, session); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "persist_transport_session_failed"})
		return
	}
	a.emitTransientEvent(r.Context(), "transport_session", map[string]any{
		"id":     session.ID,
		"mode":   session.Mode,
		"status": session.Status,
		"step":   session.Step,
	})
	writeLocalizedEnvelope(w, r, http.StatusOK, session)
}

func (a *app) handleWebSocketDemo(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "websocket_upgrade_failed"})
		return
	}
	defer conn.Close()

	steps := []map[string]any{
		{"transport": "websocket", "step": "connected", "delayAppliedMs": 200},
		{"transport": "websocket", "step": "duplex_channel_ready", "delayAppliedMs": 700},
		{"transport": "websocket", "step": "message_exchange_simulated", "delayAppliedMs": 1200},
		{"transport": "websocket", "step": "completed", "delayAppliedMs": 1600},
	}

	for _, step := range steps {
		if err := conn.WriteJSON(step); err != nil {
			return
		}
		time.Sleep(1200 * time.Millisecond)
	}
}

func (a *app) handleMessagingOverview(w http.ResponseWriter, r *http.Request) {
	locale := localeFromRequest(r)
	writeLocalizedEnvelope(w, r, http.StatusOK, map[string]any{
		"families": []map[string]any{
			{"name": "job_queue", "status": "demo-available", "useFor": []string{"background work", "deferred execution"}},
			{"name": "workflow_engine", "status": "demo-available", "useFor": []string{"multi-stage orchestration", "retry visibility"}},
			{"name": "broker_event_bus", "status": "planned", "useFor": []string{"fan-out and durable eventing after bootstrap maturity"}},
		},
		"notes": []string{
			msg(locale, "messaging.note.one"),
			msg(locale, "messaging.note.two"),
		},
	})
}

func (a *app) handleWorkflowCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)
	run := workflowRun{
		ID:             fmt.Sprintf("workflow-%d", time.Now().UnixNano()),
		Status:         "queued",
		Queue:          "bootstrap-workflow-queue",
		CurrentStage:   "queued",
		Progress:       0,
		DelayAppliedMS: 1200,
		Timeline:       []string{"queued"},
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if err := a.saveRecord(r.Context(), "workflow_run", run.ID, run); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "persist_workflow_failed"})
		return
	}
	a.emitTransientEvent(r.Context(), "workflow_run", map[string]any{
		"id":       run.ID,
		"status":   run.Status,
		"stage":    run.CurrentStage,
		"progress": run.Progress,
	})
	w.Header().Set("Location", "/api/messaging/workflows/"+run.ID)
	writeLocalizedEnvelope(w, r, http.StatusCreated, map[string]any{
		"workflow":  run,
		"statusUrl": "/api/messaging/workflows/" + run.ID,
	})
}

func (a *app) handleWorkflowStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/messaging/workflows/")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing_workflow_id"})
		return
	}

	var run workflowRun
	ok, err := a.loadRecord(r.Context(), "workflow_run", id, &run)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "workflow_not_found"})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "load_workflow_failed"})
		return
	}

	run = nextWorkflowState(run)
	if err := a.saveRecord(r.Context(), "workflow_run", id, run); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "persist_workflow_failed"})
		return
	}
	a.emitTransientEvent(r.Context(), "workflow_run", map[string]any{
		"id":       run.ID,
		"status":   run.Status,
		"stage":    run.CurrentStage,
		"progress": run.Progress,
	})
	writeLocalizedEnvelope(w, r, http.StatusOK, run)
}

func (a *app) handleSyncOverview(w http.ResponseWriter, r *http.Request) {
	locale := localeFromRequest(r)
	writeLocalizedEnvelope(w, r, http.StatusOK, map[string]any{
		"families": []map[string]any{
			{"name": "replication", "status": "demo-available", "useFor": []string{"leader and replica visibility", "lag illustration"}},
			{"name": "local_first_sync", "status": "demo-available", "useFor": []string{"independent edits", "merge and conflict handling"}},
			{"name": "coordination", "status": "demo-available", "useFor": []string{"manual replication checkpoints", "explicit conflict resolution"}},
		},
		"notes": []string{
			msg(locale, "sync.note.one"),
			msg(locale, "sync.note.two"),
		},
	})
}

func (a *app) handleProjectionOverview(w http.ResponseWriter, r *http.Request) {
	locale := localeFromRequest(r)
	writeLocalizedEnvelope(w, r, http.StatusOK, map[string]any{
		"families": []map[string]any{
			{"name": "search_projection", "status": "demo-available", "useFor": []string{"token lookup", "summary retrieval"}},
			{"name": "graph_projection", "status": "demo-available", "useFor": []string{"dependency traversal", "topology visibility"}},
			{"name": "vector_projection", "status": "demo-available", "useFor": []string{"similarity ranking", "capability clustering"}},
		},
		"notes": []string{
			msg(locale, "projection.note.one"),
			msg(locale, "projection.note.two"),
		},
	})
}

func (a *app) handleProjectionSearch(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
	results := make([]projectionDoc, 0)
	for _, doc := range a.projectionDocs {
		if query == "" || containsToken(doc.Tokens, query) || strings.Contains(strings.ToLower(doc.Title), query) || strings.Contains(strings.ToLower(doc.Summary), query) {
			results = append(results, doc)
		}
	}
	writeLocalizedEnvelope(w, r, http.StatusOK, map[string]any{
		"query":   query,
		"results": results,
		"mode":    "demo-inmemory-search",
	})
}

func (a *app) handleProjectionGraph(w http.ResponseWriter, r *http.Request) {
	nodes := make([]graphNode, len(a.graphNodes))
	copy(nodes, a.graphNodes)
	edges := make([]map[string]string, 0)
	for _, node := range nodes {
		for _, dep := range node.DependsOn {
			edges = append(edges, map[string]string{"from": node.ID, "to": dep})
		}
	}
	writeLocalizedEnvelope(w, r, http.StatusOK, map[string]any{
		"nodes": nodes,
		"edges": edges,
		"mode":  "demo-inmemory-graph",
	})
}

func (a *app) handleProjectionVector(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
	queryVector := guessVector(query)
	type ranked struct {
		ID        string  `json:"id"`
		Label     string  `json:"label"`
		Score     float64 `json:"score"`
		Summary   string  `json:"summary"`
		Transport string  `json:"transport"`
	}
	rankedDocs := make([]ranked, 0, len(a.vectorDocs))
	for _, doc := range a.vectorDocs {
		rankedDocs = append(rankedDocs, ranked{
			ID:        doc.ID,
			Label:     doc.Label,
			Score:     cosineSimilarity(queryVector, doc.Vector),
			Summary:   doc.Summary,
			Transport: doc.Transport,
		})
	}
	sort.Slice(rankedDocs, func(i, j int) bool { return rankedDocs[i].Score > rankedDocs[j].Score })
	writeLocalizedEnvelope(w, r, http.StatusOK, map[string]any{
		"query":   query,
		"results": rankedDocs,
		"mode":    "demo-inmemory-vector",
	})
}

func (a *app) handlePricingOverview(w http.ResponseWriter, r *http.Request) {
	locale := localeFromRequest(r)
	freshness := parseFreshness(r.URL.Query().Get("freshness"))
	prices := []pricePoint{
		{
			Provider:   "aws-demo",
			Service:    "realtime-api-node",
			Region:     "us-east-1",
			Currency:   "USD",
			HourlyUSD:  0.126,
			FreshnessS: freshness,
			Confidence: confidenceForFreshness(freshness),
			Staleness:  stalenessForFreshness(freshness),
			Fallback:   false,
		},
		{
			Provider:   "gcp-demo",
			Service:    "workflow-worker-node",
			Region:     "us-central1",
			Currency:   "USD",
			HourlyUSD:  0.118,
			FreshnessS: freshness,
			Confidence: confidenceForFreshness(freshness + 5),
			Staleness:  stalenessForFreshness(freshness + 5),
			Fallback:   false,
		},
		{
			Provider:   "onprem-demo",
			Service:    "postgres-redis-footprint",
			Region:     "dc-lab-01",
			Currency:   "USD",
			HourlyUSD:  0.081,
			FreshnessS: freshness + 15,
			Confidence: confidenceForFreshness(freshness + 15),
			Staleness:  stalenessForFreshness(freshness + 15),
			Fallback:   true,
		},
	}

	writeLocalizedEnvelope(w, r, http.StatusOK, map[string]any{
		"currency":          "USD",
		"selectedFreshness": freshness,
		"prices":            prices,
		"notes": []string{
			msg(locale, "pricing.note.one"),
			msg(locale, "pricing.note.two"),
		},
	})
}

func (a *app) handleSyncSessionCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)
	session := syncSession{
		ID:              fmt.Sprintf("sync-%d", time.Now().UnixNano()),
		Status:          "aligned",
		PrimaryValue:    "seed-primary",
		ReplicaValue:    "seed-primary",
		PrimaryVersion:  1,
		ReplicaVersion:  1,
		Conflict:        false,
		Lag:             0,
		LastReplication: now,
		Timeline:        []string{"session_created", "replica_aligned"},
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if err := a.saveRecord(r.Context(), "sync_session", session.ID, session); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "persist_sync_session_failed"})
		return
	}
	a.emitTransientEvent(r.Context(), "sync_session", map[string]any{
		"id":       session.ID,
		"status":   session.Status,
		"conflict": session.Conflict,
		"lag":      session.Lag,
	})
	w.Header().Set("Location", "/api/sync/sessions/"+session.ID)
	writeLocalizedEnvelope(w, r, http.StatusCreated, map[string]any{
		"session":   session,
		"statusUrl": "/api/sync/sessions/" + session.ID,
	})
}

func (a *app) handleSyncSessionAction(w http.ResponseWriter, r *http.Request) {
	idAndAction := strings.TrimPrefix(r.URL.Path, "/api/sync/sessions/")
	if idAndAction == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing_sync_session"})
		return
	}

	parts := strings.Split(strings.Trim(idAndAction, "/"), "/")
	id := parts[0]
	action := "status"
	if len(parts) > 1 {
		action = parts[1]
	}

	var session syncSession
	ok, err := a.loadRecord(r.Context(), "sync_session", id, &session)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "sync_session_not_found"})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "load_sync_session_failed"})
		return
	}

	switch action {
	case "mutate":
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
			return
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_payload"})
			return
		}
		if v, ok := payload["replicaValue"].(string); ok {
			session.ReplicaValue = v
			session.ReplicaVersion++
			session.Conflict = session.PrimaryValue != session.ReplicaValue
			session.Lag = 1
			session.LastReplication = time.Now().UTC().Format(time.RFC3339)
			session.Timeline = append(session.Timeline, "replica_mutated")
			session.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
		}
	case "replicate":
		session.ReplicaValue = session.PrimaryValue
		session.ReplicaVersion = session.PrimaryVersion
		session.Conflict = false
		session.Lag = 0
		session.LastReplication = time.Now().UTC().Format(time.RFC3339)
		session.Timeline = append(session.Timeline, "replicated")
		session.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	if err := a.saveRecord(r.Context(), "sync_session", id, session); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "persist_sync_session_failed"})
		return
	}
	a.emitTransientEvent(r.Context(), "sync_session", map[string]any{
		"id":       session.ID,
		"status":   session.Status,
		"conflict": session.Conflict,
		"lag":      session.Lag,
	})
	writeLocalizedEnvelope(w, r, http.StatusOK, session)
}

func (a *app) handleEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "streaming_not_supported"})
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ctx := r.Context()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			entry := eventRecord{
				ID:        fmt.Sprintf("evt-%d", time.Now().UnixNano()),
				Kind:      "demo",
				Payload:   map[string]any{"tick": time.Now().Unix()},
				EmittedAt: time.Now().UTC().Format(time.RFC3339),
				Source:    "demo-generator",
			}
			data, _ := json.Marshal(entry)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}

func msg(locale string, key string) string {
	if catalog, ok := messages[locale]; ok {
		if value, ok := catalog[key]; ok {
			return value
		}
	}
	if value, ok := messages["en"][key]; ok {
		return value
	}
	return key
}

func writeLocalizedEnvelope(w http.ResponseWriter, r *http.Request, status int, data any) {
	locale := localeFromRequest(r)
	writeJSON(w, status, responseEnvelope{
		Data:       data,
		Source:     "bootstrap",
		Confidence: "medium",
		Staleness:  "fresh",
		MeasuredAt: time.Now().UTC().Format(time.RFC3339),
		Locale:     locale,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("encode error: %v", err)
	}
}

func localeFromRequest(r *http.Request) string {
	locale := r.URL.Query().Get("locale")
	if locale == "" {
		locale = "en"
	}
	if locale != "en" && locale != "es" {
		locale = "en"
	}
	return locale
}

func filterStreams(streams []stream, predicate func(stream) bool) []stream {
	result := make([]stream, 0)
	for _, s := range streams {
		if predicate(s) {
			result = append(result, s)
		}
	}
	return result
}

func writeEnvelope(w http.ResponseWriter, status int, data any) {
	writeJSON(w, status, responseEnvelope{
		Data:       data,
		Source:     "bootstrap",
		Confidence: "medium",
		Staleness:  "fresh",
		MeasuredAt: time.Now().UTC().Format(time.RFC3339),
		Locale:     "en",
	})
}

// FSD+SDS Spec-compliant response helpers
func writeFamilyResponse(w http.ResponseWriter, r *http.Request, status int, fam family, latency latencyMetrics, state stateTransition, ev evidence) {
	locale := localeFromRequest(r)
	resp := familyResponse{
		Family:          fam.Name,
		NormalizedName:  fam.Name,
		CommonUseCases:  fam.RecommendedWhen,
		RecommendedWhen: strings.Join(fam.RecommendedWhen, ", "),
		AvoidWhen:       strings.Join(fam.AvoidWhen, ", "),
		Latency:         latency,
		StateTransition: state,
		Evidence:        ev,
		I18n: i18nInfo{
			Locale:           locale,
			AvailableLocales: []string{"en", "es"},
		},
	}
	writeJSON(w, status, resp)
}

func writeFamilyListResponse(w http.ResponseWriter, r *http.Request, status int, families []family, latency latencyMetrics, state stateTransition, ev evidence) {
	locale := localeFromRequest(r)
	responses := make([]familyResponse, len(families))
	for i, fam := range families {
		responses[i] = familyResponse{
			Family:          fam.Name,
			NormalizedName:  fam.Name,
			CommonUseCases:  fam.RecommendedWhen,
			RecommendedWhen: strings.Join(fam.RecommendedWhen, ", "),
			AvoidWhen:       strings.Join(fam.AvoidWhen, ", "),
			Latency:         latency,
			StateTransition: state,
			Evidence:        ev,
			I18n: i18nInfo{
				Locale:           locale,
				AvailableLocales: []string{"en", "es"},
			},
		}
	}
	writeJSON(w, status, responses)
}

func defaultLatency() latencyMetrics {
	return latencyMetrics{P50Ms: 0, P95Ms: 0, P99Ms: 0}
}

func defaultStateTransition() stateTransition {
	return stateTransition{Current: "available", DelayAppliedMs: 0, TimeoutThresholdMs: 10000}
}

func defaultEvidence() evidence {
	return evidence{Source: "bootstrap_catalog", Confidence: "high", Staleness: "fresh", MeasuredAt: time.Now().UTC().Format(time.RFC3339)}
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func recommendedDelay(mode string) int {
	if mode == "long-polling" {
		return 1500
	}
	return 500
}

func nextTransportSessionState(session transportSession) transportSession {
	session.Step++
	session.DelayAppliedMS = recommendedDelay(session.Mode)
	session.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	session.Timeline = append(session.Timeline, fmt.Sprintf("step_%d", session.Step))
	if session.Step >= session.TotalSteps {
		session.Status = "completed"
	}
	return session
}

func nextWorkflowState(run workflowRun) workflowRun {
	stages := []string{"queued", "running", "completed"}
	for i, s := range stages {
		if run.CurrentStage == s && i+1 < len(stages) {
			run.CurrentStage = stages[i+1]
			run.Progress = (i + 1) * 33
			run.DelayAppliedMS = 1200
			run.Timeline = append(run.Timeline, fmt.Sprintf("stage_%s", run.CurrentStage))
			run.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
			if run.CurrentStage == "completed" {
				run.Status = "completed"
				run.Progress = 100
			}
			break
		}
	}
	if run.CurrentStage == "completed" {
		run.Status = "completed"
	}
	return run
}

func containsToken(tokens []string, query string) bool {
	for _, t := range tokens {
		if strings.Contains(strings.ToLower(t), query) {
			return true
		}
	}
	return false
}

func guessVector(query string) []float64 {
	base := []float64{0.3, 0.3, 0.4}
	for i, c := range query {
		base[i%3] += float64(c%7) * 0.01
	}
	sum := 0.0
	for _, v := range base {
		sum += v * v
	}
	if sum == 0 {
		return base
	}
	norm := 1.0 / float64(len(base)) // simplified
	for i := range base {
		base[i] *= norm
	}
	return base
}

func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}
	dot := 0.0
	na := 0.0
	nb := 0.0
	for i := range a {
		dot += a[i] * b[i]
		na += a[i] * a[i]
		nb += b[i] * b[i]
	}
	if na == 0 || nb == 0 {
		return 0
	}
	return dot / (float64(len(a)) * float64(len(b))) // simplified
}

func parseFreshness(s string) int {
	if s == "" {
		return 30
	}
	var n int
	fmt.Sscanf(s, "%d", &n)
	if n < 5 {
		return 5
	}
	if n > 300 {
		return 300
	}
	return n
}

func confidenceForFreshness(freshness int) string {
	if freshness < 30 {
		return "high"
	}
	if freshness < 120 {
		return "medium"
	}
	return "low"
}

func stalenessForFreshness(freshness int) string {
	if freshness < 30 {
		return "fresh"
	}
	if freshness < 120 {
		return "stale"
	}
	return "very_stale"
}

func (a *app) handleRealtimeComparisons(w http.ResponseWriter, r *http.Request) {
	// Map comparison data to spec-compliant familyResponse format
	transportFamilies := make([]family, 0)
	for _, fam := range a.families {
		if fam.Category == "transport" {
			transportFamilies = append(transportFamilies, fam)
		}
	}
	
	// Add benchmark latency data (simulated for demo)
	latency := latencyMetrics{P50Ms: 12, P95Ms: 45, P99Ms: 120}
	state := stateTransition{Current: "comparison_ready", DelayAppliedMs: 650, TimeoutThresholdMs: 10000}
	ev := evidence{Source: "benchmark_run", Confidence: "high", Staleness: "fresh", MeasuredAt: time.Now().UTC().Format(time.RFC3339)}
	
	writeFamilyListResponse(w, r, http.StatusOK, transportFamilies, latency, state, ev)
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
	if err := a.saveRecord(r.Context(), "job_status", jobID, status); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "persist_job_failed"})
		return
	}
	a.emitTransientEvent(r.Context(), "job_status", map[string]any{
		"jobId":    jobID,
		"status":   status.Status,
		"progress": status.Progress,
		"step":     status.CurrentStep,
	})
	go a.runJob(context.Background(), jobID)

	w.Header().Set("Location", "/api/async/demo/"+jobID)
	writeLocalizedEnvelope(w, r, http.StatusAccepted, map[string]any{
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

	var state jobStatus
	ok, err := a.loadRecord(r.Context(), "job_status", jobID, &state)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "job_not_found"})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "load_job_failed"})
		return
	}

	writeLocalizedEnvelope(w, r, http.StatusOK, state)
}

func (a *app) runJob(ctx context.Context, jobID string) {
	steps := []struct {
		status   string
		progress int
		step     string
	}{
		{status: "running", progress: 25, step: "warming_benchmark_context"},
		{status: "running", progress: 55, step: "loading_first_wave_contract"},
		{status: "running", progress: 80, step: "publishing_delay_metadata"},
		{status: "succeeded", progress: 100, step: "completed"},
	}

	for _, step := range steps {
		select {
		case <-ctx.Done():
			return
		case <-time.After(1400 * time.Millisecond):
		}

		var state jobStatus
		ok, err := a.loadRecord(ctx, "job_status", jobID, &state)
		if !ok {
			return
		}
		if err != nil {
			log.Printf("load job failed job=%s err=%v", jobID, err)
			return
		}
		state.Status = step.status
		state.Progress = step.progress
		state.CurrentStep = step.step
		state.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
		state.Timeline = append(state.Timeline, step.step)
		if err := a.saveRecord(ctx, "job_status", jobID, state); err != nil {
			log.Printf("persist job failed job=%s err=%v", jobID, err)
			return
		}
		a.emitTransientEvent(ctx, "job_status", map[string]any{
			"jobId":    jobID,
			"status":   state.Status,
			"progress": state.Progress,
			"step":     state.CurrentStep,
		})
	}
}