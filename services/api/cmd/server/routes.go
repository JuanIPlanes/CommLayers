package main

import routerpkg "commlayers/services/api/internal/router"

func buildRoutes(application *app) routerpkg.Routes {
	return routerpkg.Routes{
		Healthz:                application.handleHealth,
		Bootstrap:              application.handleBootstrap,
		Streams:                application.handleStreams,
		Catalog:                application.handleCatalog,
		FirstWaveContract:      application.handleFirstWaveContract,
		FirstWaveStream:        application.handleFirstWaveStream,
		SecurityBootstrap:      application.handleSecurityBootstrap,
		DataPlatform:           application.handleDataPlatform,
		BenchmarkFramework:     application.handleBenchmarkFramework,
		DeferredWaves:          application.handleDeferredWaves,
		V2Readiness:            application.handleV2Readiness,
		RealtimeComparisons:    application.handleRealtimeComparisons,
		TransportSummary:       application.handleTransportSummary,
		PollingSession:         application.handlePollingSession,
		PollingStatus:          application.handlePollingStatus,
		MessagingOverview:      application.handleMessagingOverview,
		WorkflowCreate:         application.handleWorkflowCreate,
		WorkflowStatus:         application.handleWorkflowStatus,
		SyncOverview:           application.handleSyncOverview,
		SyncSessionCreate:      application.handleSyncSessionCreate,
		SyncSessionAction:      application.handleSyncSessionAction,
		ProjectionOverview:     application.handleProjectionOverview,
		ProjectionSearch:       application.handleProjectionSearch,
		ProjectionGraph:        application.handleProjectionGraph,
		ProjectionVector:       application.handleProjectionVector,
		PricingOverview:        application.handlePricingOverview,
		V2Roadmap:              application.handleV2Roadmap,
		V2EventDriven:          application.handleV2EventDriven,
		V2ReactiveStreamDriven: application.handleV2ReactiveStreamDriven,
		V2LocalFirstCRDT:       application.handleV2LocalFirstCRDT,
		V2ConsensusDriven:      application.handleV2ConsensusDriven,
		V2PlanetScaleComposite: application.handleV2PlanetScaleComposite,
		AsyncDemo:              application.handleAsyncDemo,
		AsyncStatus:            application.handleAsyncStatus,
		Events:                 application.handleEvents,
		WebSocketDemo:          application.handleWebSocketDemo,
	}
}
