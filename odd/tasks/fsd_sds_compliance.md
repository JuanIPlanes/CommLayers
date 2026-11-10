---
# FSD + SDS Compliance - Backend-First Revision

## Specs
S1: Response Model - Reshape all API responses to match the required JSON structure with family, normalized_name, common_use_cases, recommended_when, avoid_when, latency (p50/p95/p99), state_transition (current, delay_applied_ms, timeout_threshold_ms), evidence (source, confidence, staleness, measured_at), i18n (locale, available_locales)
S2: Testing - Implement TDD cycle (RED→GREEN→REFACTOR) with mandatory test categories: unit, integration, protocol contract, benchmark capture, security/auth behavior, localization response correctness
S3: Benchmark Evidence - Add benchmark capture for each comparative stream with recorded p50/p95/p99 latencies and timing metadata
S4: Security/Auth - Implement BFF-oriented session model, add audit hooks for all endpoints
S5: Decision Documentation - Create DEC-001 through DEC-007 documentation aligning with FSD+SDS
S6: Application Matrix - Complete mapping of all 14 application types to communication families
S7: Derived Projections - Integrate Elasticsearch, Graph DB, pgvector transports into response model
S8: Stream Acceptance - Ensure all streams have: implementation, normalized output, benchmark evidence, localization, visible delays, security hooks

## Tasks
T1: Response Model Refactor - Update responseEnvelope and all handlers to produce spec-compliant JSON
T2: Benchmark Infrastructure - Add benchmark capture middleware and storage for p50/p95/p99 metrics
T3: Unit Tests - Add unit tests for all handlers and core functions (RED→GREEN→REFACTOR)
T4: Integration Tests - Add integration tests for API endpoints with database
T5: Protocol Contract Tests - Add contract tests for normalized response formats
T6: Benchmark Tests - Add benchmark capture tests with recorded metrics
T7: Security/Auth Tests - Add BFF auth and audit hook tests
T8: Localization Tests - Add EN/ES localization response correctness tests
T9: Security Implementation - Implement BFF session model and audit hooks
T10: Decision Docs - Create DEC-001 through DEC-007 markdown documents
T11: Application Matrix - Complete 14 application type mappings in catalog
T12: Derived Projections Integration - Wire Elasticsearch, Neo4j, Qdrant into responses
T13: Stream Acceptance Verification - Verify all streams meet acceptance criteria

## Log
L1: User requested full FSD+SDS compliance for backend-first revision
L2: Initial assessment: ~65-70% complete, critical gaps in response model, testing, benchmark evidence, security, decision docs