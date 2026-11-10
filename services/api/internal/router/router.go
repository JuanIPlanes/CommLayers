package router

import (
	authpkg "commlayers/services/api/internal/auth"
	"net/http"
)

type Routes struct {
	Healthz                http.HandlerFunc
	Bootstrap              http.HandlerFunc
	Streams                http.HandlerFunc
	Catalog                http.HandlerFunc
	FirstWaveContract      http.HandlerFunc
	FirstWaveStream        http.HandlerFunc
	SecurityBootstrap      http.HandlerFunc
	DataPlatform           http.HandlerFunc
	BenchmarkFramework     http.HandlerFunc
	DeferredWaves          http.HandlerFunc
	V2Readiness            http.HandlerFunc
	RealtimeComparisons    http.HandlerFunc
	TransportSummary       http.HandlerFunc
	PollingSession         http.HandlerFunc
	PollingStatus          http.HandlerFunc
	MessagingOverview      http.HandlerFunc
	WorkflowCreate         http.HandlerFunc
	WorkflowStatus         http.HandlerFunc
	SyncOverview           http.HandlerFunc
	SyncSessionCreate      http.HandlerFunc
	SyncSessionAction      http.HandlerFunc
	ProjectionOverview     http.HandlerFunc
	ProjectionSearch       http.HandlerFunc
	ProjectionGraph        http.HandlerFunc
	ProjectionVector       http.HandlerFunc
	PricingOverview        http.HandlerFunc
	V2Roadmap              http.HandlerFunc
	V2EventDriven          http.HandlerFunc
	V2ReactiveStreamDriven http.HandlerFunc
	V2LocalFirstCRDT       http.HandlerFunc
	V2ConsensusDriven      http.HandlerFunc
	V2PlanetScaleComposite http.HandlerFunc
	AsyncDemo              http.HandlerFunc
	AsyncStatus            http.HandlerFunc
	Events                 http.HandlerFunc
	WebSocketDemo          http.HandlerFunc
	BenchmarkMetrics       http.HandlerFunc
}

func BuildMux(routes Routes) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", routes.Healthz)
	mux.HandleFunc("/api/bootstrap", routes.Bootstrap)
	mux.HandleFunc("/api/streams", routes.Streams)
	mux.HandleFunc("/api/catalog", routes.Catalog)
	mux.HandleFunc("/api/first-wave/contract", routes.FirstWaveContract)
	mux.HandleFunc("/api/first-wave/streams/", routes.FirstWaveStream)
	mux.HandleFunc("/api/security/bootstrap", routes.SecurityBootstrap)
	mux.HandleFunc("/api/data-platform", routes.DataPlatform)
	mux.HandleFunc("/api/benchmark-framework", routes.BenchmarkFramework)
	mux.HandleFunc("/api/deferred-waves", routes.DeferredWaves)
	mux.HandleFunc("/api/v2-readiness", routes.V2Readiness)
	mux.HandleFunc("/api/comparisons/realtime", routes.RealtimeComparisons)
	mux.HandleFunc("/api/transports", routes.TransportSummary)
	mux.HandleFunc("/api/transports/polling", routes.PollingSession)
	mux.HandleFunc("/api/transports/polling/", routes.PollingStatus)
	mux.HandleFunc("/api/messaging", routes.MessagingOverview)
	mux.HandleFunc("/api/messaging/workflows", routes.WorkflowCreate)
	mux.HandleFunc("/api/messaging/workflows/", routes.WorkflowStatus)
	mux.HandleFunc("/api/sync", routes.SyncOverview)
	mux.HandleFunc("/api/sync/sessions", routes.SyncSessionCreate)
	mux.HandleFunc("/api/sync/sessions/", routes.SyncSessionAction)
	mux.HandleFunc("/api/projections", routes.ProjectionOverview)
	mux.HandleFunc("/api/projections/search", routes.ProjectionSearch)
	mux.HandleFunc("/api/projections/graph", routes.ProjectionGraph)
	mux.HandleFunc("/api/projections/vector", routes.ProjectionVector)
	mux.HandleFunc("/api/pricing", routes.PricingOverview)
	mux.HandleFunc("/api/v2/roadmap", routes.V2Roadmap)
	mux.HandleFunc("/api/v2/paradigms/event-driven", routes.V2EventDriven)
	mux.HandleFunc("/api/v2/paradigms/reactive-stream-driven", routes.V2ReactiveStreamDriven)
	mux.HandleFunc("/api/v2/paradigms/local-first-crdt", routes.V2LocalFirstCRDT)
	mux.HandleFunc("/api/v2/paradigms/consensus-driven", routes.V2ConsensusDriven)
	mux.HandleFunc("/api/v2/paradigms/planet-scale-composite", routes.V2PlanetScaleComposite)
	mux.HandleFunc("/api/async/demo", routes.AsyncDemo)
	mux.HandleFunc("/api/async/demo/", routes.AsyncStatus)
	mux.HandleFunc("/api/events", routes.Events)
	mux.HandleFunc("/api/ws/demo", routes.WebSocketDemo)
	mux.HandleFunc("/api/benchmarks/metrics", routes.BenchmarkMetrics)
	return mux
}

func BuildRouter(base http.Handler, bm *BenchmarkMiddleware) http.Handler {
	// Apply benchmark middleware first to capture all requests
	wrapped := bm.Handler(base)
	if authpkg.IsEnabled() {
		return authpkg.KeycloakMiddleware(wrapped)
	}
	return wrapped
}
