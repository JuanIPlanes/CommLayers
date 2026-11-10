package keycloak

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

// BFFSession represents a BFF (Backend-for-Frontend) session
// This follows the BFF pattern where the backend holds the session state
// and the frontend only holds a session cookie/reference.
type BFFSession struct {
	ID           string                 `json:"id"`
	UserID       string                 `json:"user_id"`
	Email        string                 `json:"email"`
	Roles        []string               `json:"roles"`
	AccessToken  string                 `json:"-"` // Not serialized
	RefreshToken string                 `json:"-"` // Not serialized
	ExpiresAt    time.Time              `json:"expires_at"`
	CreatedAt    time.Time              `json:"created_at"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// AuditEvent represents a security audit log entry
type AuditEvent struct {
	ID        string                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	UserID    string                 `json:"user_id,omitempty"`
	SessionID string                 `json:"session_id,omitempty"`
	Action    string                 `json:"action"`
	Resource  string                 `json:"resource"`
	Method    string                 `json:"method"`
	Path      string                 `json:"path"`
	IP        string                 `json:"ip"`
	UserAgent string                 `json:"user_agent"`
	Status    int                    `json:"status"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// AuditLogger interface for pluggable audit logging
type AuditLogger interface {
	Log(ctx context.Context, event AuditEvent) error
}

// DefaultAuditLogger logs audit events to stdout
type DefaultAuditLogger struct{}

func (d *DefaultAuditLogger) Log(ctx context.Context, event AuditEvent) error {
	data, _ := json.Marshal(event)
	log.Printf("AUDIT: %s", string(data))
	return nil
}

// BFFSessionStore interface for session storage
type BFFSessionStore interface {
	Create(ctx context.Context, session *BFFSession) error
	Get(ctx context.Context, sessionID string) (*BFFSession, error)
	Delete(ctx context.Context, sessionID string) error
	Refresh(ctx context.Context, sessionID string) (*BFFSession, error)
}

// InMemorySessionStore is a simple in-memory implementation for demo/testing
type InMemorySessionStore struct {
	sessions map[string]*BFFSession
}

func NewInMemorySessionStore() *InMemorySessionStore {
	return &InMemorySessionStore{
		sessions: make(map[string]*BFFSession),
	}
}

func (s *InMemorySessionStore) Create(ctx context.Context, session *BFFSession) error {
	s.sessions[session.ID] = session
	return nil
}

func (s *InMemorySessionStore) Get(ctx context.Context, sessionID string) (*BFFSession, error) {
	session, ok := s.sessions[sessionID]
	if !ok {
		return nil, nil
	}
	if time.Now().After(session.ExpiresAt) {
		delete(s.sessions, sessionID)
		return nil, nil
	}
	return session, nil
}

func (s *InMemorySessionStore) Delete(ctx context.Context, sessionID string) error {
	delete(s.sessions, sessionID)
	return nil
}

func (s *InMemorySessionStore) Refresh(ctx context.Context, sessionID string) (*BFFSession, error) {
	session, ok := s.sessions[sessionID]
	if !ok {
		return nil, nil
	}
	if time.Now().After(session.ExpiresAt) {
		delete(s.sessions, sessionID)
		return nil, nil
	}
	session.ExpiresAt = time.Now().Add(24 * time.Hour)
	return session, nil
}

var (
	enabled         bool
	auditLogger     AuditLogger = &DefaultAuditLogger{}
	sessionStore    BFFSessionStore = NewInMemorySessionStore()
)

// InitKeycloak enables the BFF-oriented Keycloak integration surface
func InitKeycloak() {
	enabled = true
}

// IsEnabled reports if the Keycloak surface has been initialized
func IsEnabled() bool {
	return enabled
}

// SetAuditLogger sets a custom audit logger
func SetAuditLogger(logger AuditLogger) {
	if logger != nil {
		auditLogger = logger
	}
}

// SetSessionStore sets a custom session store
func SetSessionStore(store BFFSessionStore) {
	if store != nil {
		sessionStore = store
	}
}

// KeycloakMiddleware protects routes behind a Bearer token check
// In a full implementation this would validate the token with Keycloak server
func KeycloakMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !enabled {
			next.ServeHTTP(w, r)
			return
		}
		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("unauthorized: missing or invalid bearer token"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AuditMiddleware logs security-relevant events for audit trail
func AuditMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := r.Context()

		// Capture response status
		lrw := newLoggingResponseWriter(w)

		next.ServeHTTP(lrw, r)

		// Build audit event
		event := AuditEvent{
			ID:        generateAuditID(),
			Timestamp: time.Now().UTC(),
			Action:    r.Method,
			Resource:  r.URL.Path,
			Method:    r.Method,
			Path:      r.URL.Path,
			IP:        getClientIP(r),
			UserAgent: r.UserAgent(),
			Status:    lrw.statusCode,
			Metadata: map[string]interface{}{
				"duration_ms": time.Since(start).Milliseconds(),
			},
		}

		// Extract user ID from context if available
		if userID := ctx.Value("user_id"); userID != nil {
			event.UserID = userID.(string)
		}
		if sessionID := ctx.Value("session_id"); sessionID != nil {
			event.SessionID = sessionID.(string)
		}

		// Log asynchronously to not block response
		go func() {
			if err := auditLogger.Log(ctx, event); err != nil {
				log.Printf("audit log error: %v", err)
			}
		}()
	})
}

// loggingResponseWriter wraps http.ResponseWriter to capture status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func generateAuditID() string {
	return "audit-" + time.Now().UTC().Format("20060102150405.000000")
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}
	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}

// CreateSession creates a new BFF session
func CreateSession(ctx context.Context, userID, email string, roles []string) (*BFFSession, error) {
	session := &BFFSession{
		ID:        "sess-" + time.Now().UTC().Format("20060102150405.000000"),
		UserID:    userID,
		Email:     email,
		Roles:     roles,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now().UTC(),
		Metadata:  make(map[string]interface{}),
	}
	if err := sessionStore.Create(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

// GetSession retrieves a session by ID
func GetSession(ctx context.Context, sessionID string) (*BFFSession, error) {
	return sessionStore.Get(ctx, sessionID)
}

// DeleteSession deletes a session
func DeleteSession(ctx context.Context, sessionID string) error {
	return sessionStore.Delete(ctx, sessionID)
}

// RefreshSession refreshes a session's expiration
func RefreshSession(ctx context.Context, sessionID string) (*BFFSession, error) {
	return sessionStore.Refresh(ctx, sessionID)
}
