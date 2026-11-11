---
id: DEC-008
title: Redis Usage Boundaries
status: proposed
date: 2026-11-11
---

## Context
Redis is used as an auxiliary cache and coordination layer in the CommLayers platform. The canonical data source is PostgreSQL (via the Store). Redis must not become a source of truth. This decision formalizes the RedisUsageBounds from `backend/redis.go` as a binding contract to prevent architectural drift.

## Decision
The RedisUsageBounds contract defines the following binding boundaries:

### Cache TTL
- **CacheTTL**: 30 minutes (1800 seconds) — default TTL for cached entries
- Entries expire automatically; stale reads are NOT allowed by default (`AllowStaleReads: false`)
- Maximum cache entries: configurable via `CachePolicy.MaxEntries` (default unbounded)
- Eviction policy: `allkeys-lru` (configurable via `EvictionPolicyOptions`)

### Rate Limiting
- **RateLimitCapacity**: Configurable default capacity per window (via `RateLimitPolicy.DefaultCapacity`)
- **RateLimitDuration**: Configurable default window duration (via `RateLimitPolicy.DefaultDuration`)
- Strategy: `token-bucket` (configurable via `RateLimitStrategyOptions`: token-bucket, leaky-bucket, sliding-window)
- Burst allowance: configurable via `RateLimitPolicy.BurstAllowance`

### Distributed Locks
- **LockTTL**: Default TTL for distributed locks (must be set per use case, typically 10-30 seconds)
- **MaxLockAttempts**: Maximum retry attempts for acquiring a lock (default: 3)
- **LockBackoff**: Backoff duration between retry attempts (default: 100ms)
- Lock token: Unique token for lock renewal/validation (prevents premature release by another holder)

### Non-Source-of-Truth Constraint
- Redis is explicitly designated as: transient cache, rate limiting, distributed locks, ephemeral state
- All data in Redis MUST be rebuildable from PostgreSQL
- Cache warm-on-load pattern: on cache miss, load from PostgreSQL, then populate Redis
- Cache save-on-write: on PostgreSQL write, update Redis asynchronously (best-effort)

### Health Checks
- RedisHealthCheck exposes: connectivity (PONG), latency (ms), connected clients, used/max memory
- Health checks are advisory; they do not gate PostgreSQL operations

## Consequences
- Any service using Redis MUST declare its usage category (cache/rate-limit/lock/ephemeral) in its initialization
- Cache TTL changes require explicit configuration update; no implicit TTL inheritance
- Lock acquisition failures MUST fall back to PostgreSQL-based coordination or return explicit error
- Rate limit exceeded responses MUST include `Retry-After` header derived from `ResetIn`
- Monitoring must track: cache hit rate, lock contention, rate limit rejection rate

## Acceptance Criteria
- [ ] `RedisUsageBounds` struct is the single source of truth for Redis configuration
- [ ] All Redis operations in codebase reference `RedisUsageBounds` fields (no magic numbers)
- [ ] Integration test verifies: cache entries expire at `CacheTTL`, locks expire at `LockTTL`, rate limits enforce `RateLimitCapacity`/`RateLimitDuration`
- [ ] No code path treats Redis as source of truth (verified by grep for "redis" as primary read path)
- [ ] Cache warm-on-load and save-on-write patterns are documented and tested