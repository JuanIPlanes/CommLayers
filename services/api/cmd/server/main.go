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
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type app struct {
	startedAt      time.Time
	jobs           sync.Map
	sessions       sync.Map
	workflows      sync.Map
	syncRuns       sync.Map
	projectionDocs []projectionDoc
	graphNodes     []graphNode
	vectorDocs     []vectorDoc
	streams        []stream
	families       []family
	waves          []deferredWave
}

type stream struct {
	ID               string   `json:"id"`
	Title            string   `json:"title"`
	Role             string   `json:"role"`
	Category         string   `json:"category"`
	Description      string   `json:"description"`
	Status           string   `json:"status"`
	Dependencies     []string `json:"dependencies"`
	Deliverables     []string `json:"deliverables"`
	AcceptanceGates  []string `json:"acceptanceGates"`
	ExecutionPackets []string `json:"executionPackets"`
}

type family struct {
	Name            string   `json:"name"`
	Category        string   `json:"category"`
	Status          string   `json:"status"`
	RecommendedWhen []string `json:"recommendedWhen"`
	AvoidWhen       []string `json:"avoidWhen"`
	Notes           []string `json:"notes"`
}

type deferredWave struct {
	Name             string   `json:"name"`
	Purpose          string   `json:"purpose"`
	Streams          []string `json:"streams"`
	Focus            []string `json:"focus"`
	ReactivationRule []string `json:"reactivationRule"`
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
	Locale     string `json:"locale"`
}

var messages = map[string]map[string]string{
	"en": {
		"bootstrap.summary":   "Backend-first bootstrap with a concrete first-wave contract surface and deferred later waves.",
		"bootstrap.phase":     "first-wave materialized",
		"bootstrap.go_no_go":  "Hold until 86agbv0k4 is finalized, 86agbv0k3 and 86agbv0k7 overlap safely, and 86agbv0k9 closes the wave.",
		"catalog.note.one":    "This is the first working version, not the full platform.",
		"catalog.note.two":    "Later backend families remain staged until the first-wave closure gate is accepted.",
		"catalog.note.three":  "SSE is the default push example in the bootstrap; WebSocket remains planned, not fully implemented.",
		"security.state":      "documented and staged; runtime integration intentionally deferred in the first working version",
		"messaging.note.one":  "This slice uses an in-memory workflow run rather than a real broker.",
		"messaging.note.two":  "The goal is to expose queue and stage transitions clearly through the frontend before adding distributed messaging infrastructure.",
		"sync.note.one":       "This slice models sync and replication in memory to keep the bootstrap honest and observable.",
		"sync.note.two":       "Replication is explicit and conflict resolution is surfaced as visible state rather than hidden machinery.",
		"projection.note.one": "These demos are in-memory projections, not real Elasticsearch/Neo4j/Qdrant integrations.",
		"projection.note.two": "The goal is to prove the projection family surface before attaching heavyweight engines later.",
		"v2.status":           "not_started",
	},
	"es": {
		"bootstrap.summary":   "Bootstrap backend-first con una superficie contractual concreta de primera ola y olas posteriores diferidas.",
		"bootstrap.phase":     "primera ola materializada",
		"bootstrap.go_no_go":  "Esperar hasta que 86agbv0k4 quede finalizado, 86agbv0k3 y 86agbv0k7 se superpongan de forma segura, y 86agbv0k9 cierre la ola.",
		"catalog.note.one":    "Esta es la primera version funcional, no la plataforma completa.",
		"catalog.note.two":    "Las familias backend posteriores siguen en espera hasta que se acepte el cierre de la primera ola.",
		"catalog.note.three":  "SSE es el ejemplo push por defecto en el bootstrap; WebSocket sigue planificado pero no totalmente implementado.",
		"security.state":      "documentado y preparado; la integracion runtime se difiere intencionalmente en esta primera version funcional",
		"messaging.note.one":  "Esta porcion usa una ejecucion de workflow en memoria en lugar de un broker real.",
		"messaging.note.two":  "El objetivo es exponer con claridad las transiciones de cola y etapas en el frontend antes de agregar infraestructura de mensajeria distribuida.",
		"sync.note.one":       "Esta porcion modela sincronizacion y replicacion en memoria para mantener el bootstrap honesto y observable.",
		"sync.note.two":       "La replicacion es explicita y la resolucion de conflictos se muestra como estado visible en lugar de maquinaria oculta.",
		"projection.note.one": "Estas demos son proyecciones en memoria, no integraciones reales con Elasticsearch/Neo4j/Qdrant.",
		"projection.note.two": "El objetivo es probar la superficie de la familia de proyecciones antes de conectar motores pesados mas adelante.",
		"v2.status":           "no_iniciado",
	},
}

type comparison struct {
	Name            string   `json:"name"`
	RecommendedFor  []string `json:"recommendedFor"`
	AvoidWhen       []string `json:"avoidWhen"`
	LatencyClass    string   `json:"latencyClass"`
	DefaultDecision string   `json:"defaultDecision"`
	Notes           []string `json:"notes"`
}

type transportSession struct {
	ID                string   `json:"id"`
	Mode              string   `json:"mode"`
	Status            string   `json:"status"`
	Step              int      `json:"step"`
	TotalSteps        int      `json:"totalSteps"`
	RecommendedPollMS int      `json:"recommendedPollMs"`
	DelayAppliedMS    int      `json:"delayAppliedMs"`
	Timeline          []string `json:"timeline"`
	CreatedAt         string   `json:"createdAt"`
	UpdatedAt         string   `json:"updatedAt"`
}

type workflowRun struct {
	ID             string   `json:"id"`
	Status         string   `json:"status"`
	Queue          string   `json:"queue"`
	CurrentStage   string   `json:"currentStage"`
	Progress       int      `json:"progress"`
	DelayAppliedMS int      `json:"delayAppliedMs"`
	Timeline       []string `json:"timeline"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
}

type syncSession struct {
	ID              string   `json:"id"`
	Status          string   `json:"status"`
	PrimaryValue    string   `json:"primaryValue"`
	ReplicaValue    string   `json:"replicaValue"`
	PrimaryVersion  int      `json:"primaryVersion"`
	ReplicaVersion  int      `json:"replicaVersion"`
	Conflict        bool     `json:"conflict"`
	Lag             int      `json:"lag"`
	LastReplication string   `json:"lastReplication"`
	Timeline        []string `json:"timeline"`
	CreatedAt       string   `json:"createdAt"`
	UpdatedAt       string   `json:"updatedAt"`
}

type projectionDoc struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Tokens  []string `json:"tokens"`
	Summary string   `json:"summary"`
}

type graphNode struct {
	ID          string   `json:"id"`
	Label       string   `json:"label"`
	DependsOn   []string `json:"dependsOn"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
}

type vectorDoc struct {
	ID        string    `json:"id"`
	Label     string    `json:"label"`
	Vector    []float64 `json:"vector"`
	Summary   string    `json:"summary"`
	Transport string    `json:"transport"`
}

type pricePoint struct {
	Provider   string  `json:"provider"`
	Service    string  `json:"service"`
	Region     string  `json:"region"`
	Currency   string  `json:"currency"`
	HourlyUSD  float64 `json:"hourlyUsd"`
	FreshnessS int     `json:"freshnessS"`
	Confidence string  `json:"confidence"`
	Staleness  string  `json:"staleness"`
	Fallback   bool    `json:"fallback"`
}

func main() {
	application := newApp()

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
	mux.HandleFunc("/api/async/demo", application.handleAsyncDemo)
	mux.HandleFunc("/api/async/demo/", application.handleAsyncStatus)
	mux.HandleFunc("/api/events", application.handleEvents)
	mux.HandleFunc("/api/ws/demo", application.handleWebSocketDemo)

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

func newApp() *app {
	return &app{
		startedAt: time.Now().UTC(),
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
	locale := localeFromRequest(r)
	writeLocalizedEnvelope(w, r, http.StatusOK, map[string]any{
		"families": a.families,
		"notes": []string{
			msg(locale, "catalog.note.one"),
			msg(locale, "catalog.note.two"),
			msg(locale, "catalog.note.three"),
		},
	})
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

func (a *app) handleTransportSummary(w http.ResponseWriter, r *http.Request) {
	writeLocalizedEnvelope(w, r, http.StatusOK, map[string]any{
		"availableDemos": []map[string]any{
			{"name": "request_response", "endpoint": "/api/comparisons/realtime", "status": "available"},
			{"name": "polling", "endpoint": "/api/transports/polling", "status": "available"},
			{"name": "long_polling", "endpoint": "/api/transports/polling?mode=long-polling", "status": "available"},
			{"name": "server_sent_events", "endpoint": "/api/events", "status": "available"},
			{"name": "websocket", "endpoint": "/api/ws/demo", "status": "available"},
		},
		"notes": []string{
			"Polling returns immediately and advances state per fetch.",
			"Long-polling intentionally waits before returning the next state transition.",
			"SSE and WebSocket are both reachable through the frontend gateway.",
		},
	})
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
	a.sessions.Store(session.ID, session)
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

	raw, ok := a.sessions.Load(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "transport_session_not_found"})
		return
	}

	session := raw.(transportSession)
	if session.Mode == "long-polling" && session.Status != "completed" {
		time.Sleep(time.Duration(session.RecommendedPollMS) * time.Millisecond)
	}

	session = nextTransportSessionState(session)
	a.sessions.Store(id, session)
	writeEnvelope(w, http.StatusOK, session)
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
	a.workflows.Store(run.ID, run)
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

	raw, ok := a.workflows.Load(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "workflow_not_found"})
		return
	}

	run := nextWorkflowState(raw.(workflowRun))
	a.workflows.Store(id, run)
	writeEnvelope(w, http.StatusOK, run)
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
	a.syncRuns.Store(session.ID, session)
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

	raw, ok := a.syncRuns.Load(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "sync_session_not_found"})
		return
	}

	session := raw.(syncSession)
	switch action {
	case "status":
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
			return
		}
	case "mutate":
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
			return
		}
		target := r.URL.Query().Get("target")
		if target == "" {
			target = "primary"
		}
		value := r.URL.Query().Get("value")
		if value == "" {
			value = fmt.Sprintf("%s-edit-%d", target, time.Now().Unix()%1000)
		}
		session = mutateSyncSession(session, target, value)
	case "replicate":
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
			return
		}
		session = replicateSyncSession(session)
	default:
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "unsupported_sync_action"})
		return
	}

	a.syncRuns.Store(id, session)
	writeEnvelope(w, http.StatusOK, session)
}

func (a *app) handleRealtimeComparisons(w http.ResponseWriter, r *http.Request) {
	comparisons := []comparison{
		{
			Name:            "request_response",
			RecommendedFor:  []string{"bounded synchronous reads", "simple commands"},
			AvoidWhen:       []string{"continuous updates"},
			LatencyClass:    "baseline",
			DefaultDecision: "reference case",
			Notes:           []string{"Use as the baseline for future transport comparisons."},
		},
		{
			Name:            "server_sent_events",
			RecommendedFor:  []string{"one-way updates", "progress streams", "status feeds"},
			AvoidWhen:       []string{"bidirectional collaboration", "binary transport"},
			LatencyClass:    "low with server push",
			DefaultDecision: "default push choice",
			Notes:           []string{"Reverse proxies must disable buffering for timely delivery.", "Bootstrap demo uses artificial delay so the user can observe state transitions."},
		},
		{
			Name:            "websocket",
			RecommendedFor:  []string{"true duplex interaction", "high-frequency command loops"},
			AvoidWhen:       []string{"simple server-to-client notifications"},
			LatencyClass:    "lowest interactive overhead when justified",
			DefaultDecision: "planned, not bootstrap default",
			Notes:           []string{"Use only when duplex interaction is actually required.", "Some edge platforms inspect only the initial upgrade handshake."},
		},
	}

	writeLocalizedEnvelope(w, r, http.StatusOK, map[string]any{
		"comparisons":        comparisons,
		"delayAppliedMs":     650,
		"timeoutThresholdMs": 10000,
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
		{"step": "bootstrap_started", "delayAppliedMs": 250, "label": "Bootstrap started"},
		{"step": "first_wave_loaded", "delayAppliedMs": 700, "label": "First-wave contract loaded"},
		{"step": "deferred_waves_held", "delayAppliedMs": 1025, "label": "Deferred waves confirmed held"},
		{"step": "sse_demo_complete", "delayAppliedMs": 1350, "label": "SSE demo complete"},
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

func nextTransportSessionState(session transportSession) transportSession {
	steps := []struct {
		name     string
		progress int
	}{
		{name: "request_sent", progress: 25},
		{name: "server_processing", progress: 50},
		{name: "response_ready", progress: 80},
		{name: "completed", progress: 100},
	}

	if session.Status == "completed" {
		return session
	}

	index := session.Step
	if index >= len(steps) {
		session.Status = "completed"
		session.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
		return session
	}

	step := steps[index]
	session.Step++
	session.Status = "running"
	if step.name == "completed" {
		session.Status = "completed"
	}
	session.DelayAppliedMS = recommendedDelay(session.Mode)
	session.RecommendedPollMS = recommendedDelay(session.Mode)
	session.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	session.Timeline = append(session.Timeline, step.name)
	return session
}

func nextWorkflowState(run workflowRun) workflowRun {
	steps := []struct {
		stage    string
		status   string
		progress int
	}{
		{stage: "dequeued", status: "running", progress: 20},
		{stage: "handler_started", status: "running", progress: 45},
		{stage: "workflow_transition", status: "running", progress: 75},
		{stage: "completed", status: "completed", progress: 100},
	}

	if run.Status == "completed" {
		return run
	}

	index := len(run.Timeline) - 1
	if index >= len(steps) {
		run.Status = "completed"
		run.CurrentStage = "completed"
		run.Progress = 100
		run.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
		return run
	}

	step := steps[index]
	time.Sleep(time.Duration(run.DelayAppliedMS) * time.Millisecond)
	run.Status = step.status
	run.CurrentStage = step.stage
	run.Progress = step.progress
	run.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	run.Timeline = append(run.Timeline, step.stage)
	return run
}

func mutateSyncSession(session syncSession, target string, value string) syncSession {
	session.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if target == "replica" {
		session.ReplicaValue = value
		session.ReplicaVersion++
		session.Timeline = append(session.Timeline, "replica_mutated")
	} else {
		session.PrimaryValue = value
		session.PrimaryVersion++
		session.Timeline = append(session.Timeline, "primary_mutated")
	}
	session.Lag = abs(session.PrimaryVersion - session.ReplicaVersion)
	session.Conflict = session.PrimaryValue != session.ReplicaValue
	if session.Conflict {
		session.Status = "conflict"
	} else if session.Lag > 0 {
		session.Status = "replicating"
	} else {
		session.Status = "aligned"
	}
	return session
}

func replicateSyncSession(session syncSession) syncSession {
	session.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	session.LastReplication = session.UpdatedAt
	if session.PrimaryValue != session.ReplicaValue {
		if session.PrimaryVersion == session.ReplicaVersion {
			merged := session.PrimaryValue + " | " + session.ReplicaValue
			session.PrimaryValue = merged
			session.ReplicaValue = merged
			session.PrimaryVersion++
			session.ReplicaVersion = session.PrimaryVersion
			session.Timeline = append(session.Timeline, "conflict_merged")
		} else if session.PrimaryVersion > session.ReplicaVersion {
			session.ReplicaValue = session.PrimaryValue
			session.ReplicaVersion = session.PrimaryVersion
			session.Timeline = append(session.Timeline, "replica_caught_up")
		} else {
			session.PrimaryValue = session.ReplicaValue
			session.PrimaryVersion = session.ReplicaVersion
			session.Timeline = append(session.Timeline, "primary_caught_up")
		}
	}
	session.Conflict = false
	session.Lag = 0
	session.Status = "aligned"
	return session
}

func containsToken(tokens []string, query string) bool {
	for _, token := range tokens {
		if token == query {
			return true
		}
	}
	return false
}

func guessVector(query string) []float64 {
	switch {
	case strings.Contains(query, "sync") || strings.Contains(query, "merge"):
		return []float64{0.2, 0.3, 1.0}
	case strings.Contains(query, "workflow") || strings.Contains(query, "queue"):
		return []float64{0.2, 1.0, 0.1}
	default:
		return []float64{1.0, 0.1, 0.0}
	}
}

func cosineSimilarity(a []float64, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (sqrt(normA) * sqrt(normB))
}

func parseFreshness(value string) int {
	switch value {
	case "30":
		return 30
	case "60":
		return 60
	default:
		return 5
	}
}

func stalenessForFreshness(seconds int) string {
	if seconds <= 5 {
		return "fresh"
	}
	if seconds <= 30 {
		return "warning"
	}
	return "stale"
}

func confidenceForFreshness(seconds int) string {
	if seconds <= 5 {
		return "high"
	}
	if seconds <= 30 {
		return "medium"
	}
	return "low"
}

func sqrt(value float64) float64 {
	guess := value
	if guess == 0 {
		return 0
	}
	for i := 0; i < 8; i++ {
		guess = 0.5 * (guess + value/guess)
	}
	return guess
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func recommendedDelay(mode string) int {
	if mode == "long-polling" {
		return 1400
	}
	return 700
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

func localeFromRequest(r *http.Request) string {
	lang := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("lang")))
	if lang == "" {
		lang = strings.ToLower(strings.TrimSpace(r.Header.Get("Accept-Language")))
	}
	if strings.HasPrefix(lang, "es") {
		return "es"
	}
	return "en"
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
