package backend

import (
	"encoding/json"
	"time"
)

// RedisUsageBounds defines the Redis responsibilities for cache, rate limiting,
// ephemeral coordination, and non-source-of-truth constraints.
//
// Redis is used as an auxiliary cache and coordination layer, NOT as the
// source of truth. The canonical data source is PostgreSQL (via the Store).
// Redis is used for: transient cache, rate limiting, distributed locks,
// and ephemeral state that can be rebuilt from PostgreSQL.
type RedisUsageBounds struct {
	// CacheTTL is the time-to-live for cached entries in seconds.
	CacheTTL time.Duration `json:"cache_ttl"`
	// RateLimitCapacity is the maximum number of requests in the rate limit window.
	RateLimitCapacity int `json:"rate_limit_capacity"`
	// RateLimitDuration is the duration of the rate limit window in seconds.
	RateLimitDuration time.Duration `json:"rate_limit_duration"`
	// LockTTL is the default TTL for distributed locks in seconds.
	LockTTL time.Duration `json:"lock_ttl"`
	// MaxLockAttempts is the maximum number of retry attempts for acquiring a lock.
	MaxLockAttempts int `json:"max_lock_attempts"`
	// LockBackoff is the backoff duration between lock retry attempts.
	LockBackoff time.Duration `json:"lock_backoff"`
}

// CacheEntry represents a Redis cache entry with its metadata.
type CacheEntry struct {
	// Key is the Redis key.
	Key string `json:"key"`
	// Value is the cached value.
	Value json.RawMessage `json:"value"`
	// CreatedAt is when the entry was cached.
	CreatedAt time.Time `json:"created_at"`
	// ExpiresAt is when the entry expires (CreatedAt + CacheTTL).
	ExpiresAt time.Time `json:"expires_at"`
	// AccessCount is how many times the entry has been accessed.
	AccessCount int `json:"access_count,omitempty"`
	// LastAccessedAt is the last time the entry was accessed.
	LastAccessedAt time.Time `json:"last_accessed_at,omitempty"`
}

// RateLimitResult represents the result of a rate limit check.
type RateLimitResult struct {
	// Allowed indicates if the request is allowed.
	Allowed bool `json:"allowed"`
	// Remaining is how many requests remain in the current window.
	Remaining int `json:"remaining"`
	// ResetAt is when the rate limit window resets.
	ResetAt time.Time `json:"reset_at"`
	// ResetIn is how long until the rate limit resets.
	ResetIn time.Duration `json:"reset_in"`
}

// DistributedLock represents a distributed lock acquired via Redis.
type DistributedLock struct {
	// Key is the Redis key for the lock.
	Key string `json:"key"`
	// AcquiredAt is when the lock was acquired.
	AcquiredAt time.Time `json:"acquired_at"`
	// ExpiresAt is when the lock expires.
	ExpiresAt time.Time `json:"expires_at"`
	// Owner is the identifier of the entity that owns the lock.
	Owner string `json:"owner"`
	// Token is a unique token for lock renewal/validation.
	Token string `json:"token"`
}

// RedisHealthCheck represents the result of a Redis health check.
type RedisHealthCheck struct {
	// Healthy indicates if Redis is healthy.
	Healthy bool `json:"healthy"`
	// PongResponse is the raw PONG response from Redis.
	PongResponse string `json:"pong_response,omitempty"`
	// Latency is the round-trip latency in milliseconds.
	LatencyMs float64 `json:"latency_ms,omitempty"`
	// ConnectedClients is the number of connected Redis clients.
	ConnectedClients int `json:"connected_clients,omitempty"`
	// UsedMemory is the amount of memory used by Redis.
	UsedMemory bytes `json:"used_memory,omitempty"`
	// MaxMemory is the maximum memory configured for Redis.
	MaxMemory bytes `json:"max_memory,omitempty"`
}

// Bytes represents a byte slice with JSON marshaling support.
type bytes []byte

// UnmarshalJSON implements the json.Unmarshaler interface.
func (b *bytes) UnmarshalJSON(data []byte) error {
	*b = make(bytes, 0)
	return json.Unmarshal(data, *b)
}

// MarshalJSON implements the json.Marshaler interface.
func (b bytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(b))
}

// CachePolicy defines the policy for Redis cache entries.
type CachePolicy struct {
	// DefaultTTL is the default time-to-live for cache entries.
	DefaultTTL time.Duration `json:"default_ttl"`
	// MaxEntries is the maximum number of entries in the cache.
	MaxEntries int `json:"max_entries"`
	// EvictionPolicy is the cache eviction policy (allkeys-lru, allkeys-random, etc.).
	EvictionPolicy string `json:"eviction_policy"`
	// AllowStaleReads indicates if stale reads are allowed when the entry is expired.
	AllowStaleReads bool `json:"allow_stale_reads"`
}

// RateLimitPolicy defines the policy for rate limiting.
type RateLimitPolicy struct {
	// DefaultCapacity is the default capacity for rate limit windows.
	DefaultCapacity int `json:"default_capacity"`
	// DefaultDuration is the default duration for rate limit windows.
	DefaultDuration time.Duration `json:"default_duration"`
	// BurstAllowance is the number of additional requests allowed beyond the capacity.
	BurstAllowance int `json:"burst_allowance"`
	// Strategy is the rate limiting strategy (token-bucket, leaky-bucket, sliding-window).
	Strategy string `json:"strategy"`
}

// EvictionPolicyOptions are the available Redis cache eviction policies.
var EvictionPolicyOptions = []string{
	"allkeys-lru",
	"allkeys-random",
	"volatile-lru",
	"volatile-random",
	"no-eviction",
}

// RateLimitStrategyOptions are the available rate limiting strategies.
var RateLimitStrategyOptions = []string{
	"token-bucket",
	"leaky-bucket",
	"sliding-window",
}

// RedisConfig holds the configuration for Redis usage.
type RedisConfig struct {
	// Address is the Redis server address.
	Address string `json:"address"`
	// Password is the Redis server password (empty for no auth).
	Password string `json:"password"`
	// DB is the Redis database number.
	DB int `json:"db"`
	// PoolSize is the connection pool size.
	PoolSize int `json:"pool_size"`
	// MinIdleConns is the minimum number of idle connections.
	MinIdleConns int `json:"min_idle_conns"`
	// ReadOnly indicates if the Redis connection is read-only.
	ReadOnly bool `json:"read_only"`
	// CachePolicy holds the cache policy configuration.
	CachePolicy *CachePolicy `json:"cache_policy,omitempty"`
	// RateLimitPolicy holds the rate limit policy configuration.
	RateLimitPolicy *RateLimitPolicy `json:"rate_limit_policy,omitempty"`
}