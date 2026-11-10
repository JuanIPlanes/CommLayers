package main
import (
	"time"
	"sync"
	"commlayers/services/api/internal/transport/graph"
	"commlayers/services/api/internal/transport/vector"
	"commlayers/services/api/internal/cache"
	"commlayers/services/api/internal/persistence"
)


type app struct {
	startedAt      time.Time
	store          *persistence.Store
	cache          *cache.Cache
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
	neo4jTransport *graph.Neo4jTransport
	qdrantTransport *vector.QdrantTransport
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
		"catalog.note.four":   "All 14 common applications are mapped: CRUD, periodic status, blocking wait, one-way updates, bidirectional, webhooks, fan-out, durable workflows, event history, stream processing, actor concurrency, realtime sync, replication, coordination.",
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
		"catalog.note.four":   "Los 14 casos de uso comunes estan mapeados: CRUD, estado periodico, espera bloqueante, actualizaciones unidireccionales, bidireccional, webhooks, fan-out, workflows durables, historial de eventos, procesamiento de streams, concurrencia actores, sync tiempo real, replicacion, coordinacion.",
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

type eventRecord struct {
	ID        string         `json:"id"`
	Kind      string         `json:"kind"`
	Payload   map[string]any `json:"payload"`
	EmittedAt string         `json:"emittedAt"`
	Source    string         `json:"source"`
}

// FSD+SDS Spec-compliant response types
type latencyMetrics struct {
	P50Ms int `json:"p50_ms"`
	P95Ms int `json:"p95_ms"`
	P99Ms int `json:"p99_ms"`
}

type stateTransition struct {
	Current            string `json:"current"`
	DelayAppliedMs     int    `json:"delay_applied_ms"`
	TimeoutThresholdMs int    `json:"timeout_threshold_ms"`
}

type evidence struct {
	Source     string `json:"source"`
	Confidence string `json:"confidence"`
	Staleness  string `json:"staleness"`
	MeasuredAt string `json:"measured_at"`
}

type i18nInfo struct {
	Locale           string   `json:"locale"`
	AvailableLocales []string `json:"available_locales"`
}

type familyResponse struct {
	Family            string         `json:"family"`
	NormalizedName    string         `json:"normalized_name"`
	CommonUseCases    []string       `json:"common_use_cases"`
	RecommendedWhen   string         `json:"recommended_when"`
	AvoidWhen         string         `json:"avoid_when"`
	Latency           latencyMetrics `json:"latency"`
	StateTransition   stateTransition `json:"state_transition"`
	Evidence          evidence       `json:"evidence"`
	I18n              i18nInfo       `json:"i18n"`
}
