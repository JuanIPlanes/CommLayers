package keycloak

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestBFFSessionCreation(t *testing.T) {
	ctx := context.Background()
	
	session, err := CreateSession(ctx, "user-123", "test@example.com", []string{"user", "admin"})
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}
	
	if session.ID == "" {
		t.Error("expected session ID to be set")
	}
	if session.UserID != "user-123" {
		t.Errorf("expected user_id=user-123, got %s", session.UserID)
	}
	if session.Email != "test@example.com" {
		t.Errorf("expected email=test@example.com, got %s", session.Email)
	}
	if len(session.Roles) != 2 {
		t.Errorf("expected 2 roles, got %d", len(session.Roles))
	}
	if session.ExpiresAt.Before(time.Now()) {
		t.Error("expected expiration in future")
	}
}

func TestBFFSessionRetrieval(t *testing.T) {
	ctx := context.Background()
	
	session, err := CreateSession(ctx, "user-456", "test2@example.com", []string{"user"})
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}
	
	retrieved, err := GetSession(ctx, session.ID)
	if err != nil {
		t.Fatalf("GetSession failed: %v", err)
	}
	
	if retrieved == nil {
		t.Fatal("expected session to be found")
	}
	if retrieved.UserID != session.UserID {
		t.Errorf("expected user_id=%s, got %s", session.UserID, retrieved.UserID)
	}
}

func TestBFFSessionDeletion(t *testing.T) {
	ctx := context.Background()
	
	session, err := CreateSession(ctx, "user-789", "test3@example.com", []string{"user"})
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}
	
	if err := DeleteSession(ctx, session.ID); err != nil {
		t.Fatalf("DeleteSession failed: %v", err)
	}
	
	retrieved, err := GetSession(ctx, session.ID)
	if err != nil {
		t.Fatalf("GetSession failed: %v", err)
	}
	
	if retrieved != nil {
		t.Error("expected session to be nil after deletion")
	}
}

func TestAuditEventStructure(t *testing.T) {
	event := AuditEvent{
		ID:        "audit-123",
		Timestamp: time.Now().UTC(),
		UserID:    "user-123",
		Action:    "GET",
		Resource:  "/api/catalog",
		Method:    "GET",
		Path:      "/api/catalog",
		IP:        "127.0.0.1",
		UserAgent: "test-agent",
		Status:    200,
		Metadata:  map[string]interface{}{"duration_ms": 42},
	}
	
	if event.ID != "audit-123" {
		t.Error("expected ID to be set")
	}
	if event.UserID != "user-123" {
		t.Error("expected user_id to be set")
	}
	if event.Status != 200 {
		t.Error("expected status=200")
	}
}

func TestDefaultAuditLogger(t *testing.T) {
	logger := &DefaultAuditLogger{}
	event := AuditEvent{
		ID:       "audit-test",
		Timestamp: time.Now().UTC(),
		Action:   "GET",
		Resource: "/api/test",
	}
	
	err := logger.Log(context.Background(), event)
	if err != nil {
		t.Errorf("DefaultAuditLogger.Log failed: %v", err)
	}
}

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name     string
		headers  map[string]string
		remote   string
		expected string
	}{
		{
			name:     "X-Forwarded-For",
			headers:  map[string]string{"X-Forwarded-For": "10.0.0.1, 10.0.0.2"},
			remote:   "192.168.1.1:1234",
			expected: "10.0.0.1",
		},
		{
			name:     "X-Real-IP",
			headers:  map[string]string{"X-Real-IP": "10.0.0.3"},
			remote:   "192.168.1.1:1234",
			expected: "10.0.0.3",
		},
		{
			name:     "RemoteAddr fallback",
			headers:  map[string]string{},
			remote:   "192.168.1.4:5678",
			expected: "192.168.1.4",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			req.RemoteAddr = tt.remote
			
			result := getClientIP(req)
			if result != tt.expected {
				t.Errorf("getClientIP() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestKeycloakMiddleware(t *testing.T) {
	// Test with auth disabled
	handler := KeycloakMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 when auth disabled, got %d", w.Code)
	}
	
	// Test with auth enabled but no token
	InitKeycloak()
	defer func() { enabled = false }()
	
	handler = KeycloakMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	
	req = httptest.NewRequest("GET", "/test", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 when no token, got %d", w.Code)
	}
	
	// Test with auth enabled and valid bearer token
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 with bearer token, got %d", w.Code)
	}
	
	// Test with invalid auth header
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Basic invalid")
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 with invalid auth type, got %d", w.Code)
	}
}

func TestAuditMiddleware(t *testing.T) {
	handler := AuditMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.Header.Set("User-Agent", "test-agent")
	w := httptest.NewRecorder()
	
	handler.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestInMemorySessionStore(t *testing.T) {
	store := NewInMemorySessionStore()
	ctx := context.Background()
	
	session := &BFFSession{
		ID:        "test-session",
		UserID:    "user-1",
		Email:     "test@example.com",
		Roles:     []string{"user"},
		ExpiresAt: time.Now().Add(time.Hour),
		CreatedAt: time.Now(),
	}
	
	if err := store.Create(ctx, session); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	retrieved, err := store.Get(ctx, "test-session")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if retrieved == nil {
		t.Fatal("expected session to be found")
	}
	if retrieved.UserID != "user-1" {
		t.Errorf("expected user_id=user-1, got %s", retrieved.UserID)
	}
	
	if err := store.Delete(ctx, "test-session"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	
	retrieved, err = store.Get(ctx, "test-session")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if retrieved != nil {
		t.Error("expected session to be nil after deletion")
	}
}

func TestSessionExpiration(t *testing.T) {
	store := NewInMemorySessionStore()
	ctx := context.Background()
	
	// Create expired session
	session := &BFFSession{
		ID:        "expired-session",
		UserID:    "user-1",
		Email:     "test@example.com",
		Roles:     []string{"user"},
		ExpiresAt: time.Now().Add(-time.Hour), // Expired
		CreatedAt: time.Now().Add(-2 * time.Hour),
	}
	
	if err := store.Create(ctx, session); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	// Should return nil for expired session
	retrieved, err := store.Get(ctx, "expired-session")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if retrieved != nil {
		t.Error("expected expired session to return nil")
	}
}

func TestAuditIDGeneration(t *testing.T) {
	id1 := generateAuditID()
	id2 := generateAuditID()
	
	if id1 == "" || id2 == "" {
		t.Error("expected non-empty audit IDs")
	}
	if id1 == id2 {
		t.Error("expected unique audit IDs")
	}
	if !strings.HasPrefix(id1, "audit-") {
		t.Errorf("expected audit ID to start with 'audit-', got %s", id1)
	}
}
