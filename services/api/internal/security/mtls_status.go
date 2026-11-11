package security

import (
	"context"
	"crypto/tls"
	"sync"
	"time"
)

// MTLSConnection represents an mTLS connection status between two services.
type MTLSConnection struct {
	SourceService  string    `json:"source_service"`
	TargetService  string    `json:"target_service"`
	TLSVersion     string    `json:"tls_version"`
	CipherSuite    string    `json:"cipher_suite"`
	Verified       bool      `json:"verified"`
	CertExpiryDays int       `json:"cert_expiry_days"`
	Phase          int       `json:"phase"` // 1, 2, or 3 per mTLS rollout
	LastObserved   time.Time `json:"last_observed"`
	Error          string    `json:"error,omitempty"`
}

// MTLSStatusTracker tracks mTLS connection statuses across service pairs.
type MTLSStatusTracker struct {
	mu         sync.RWMutex
	connections map[string]*MTLSConnection // keyed by "source->target"
	phase      int                        // current mTLS rollout phase (1, 2, or 3)
}

// NewMTLSStatusTracker creates a new mTLS status tracker.
func NewMTLSStatusTracker(initialPhase int) *MTLSStatusTracker {
	if initialPhase < 1 {
		initialPhase = 1
	}
	if initialPhase > 3 {
		initialPhase = 3
	}
	return &MTLSStatusTracker{
		connections: make(map[string]*MTLSConnection),
		phase:       initialPhase,
	}
}

// SetPhase updates the current mTLS rollout phase.
func (mt *MTLSStatusTracker) SetPhase(phase int) {
	if phase < 1 {
		phase = 1
	}
	if phase > 3 {
		phase = 3
	}
	mt.mu.Lock()
	defer mt.mu.Unlock()
	mt.phase = phase
	// Update phase for all existing connections
	for _, conn := range mt.connections {
		conn.Phase = phase
	}
}

// GetPhase returns the current mTLS rollout phase.
func (mt *MTLSStatusTracker) GetPhase() int {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	return mt.phase
}

// RecordConnection records an mTLS connection observation from a TLS connection state.
func (mt *MTLSStatusTracker) RecordConnection(ctx context.Context, source, target string, state *tls.ConnectionState) {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	key := source + "->" + target
	conn := &MTLSConnection{
		SourceService: source,
		TargetService: target,
		TLSVersion:    tlsVersionString(state.Version),
		CipherSuite:   tls.CipherSuiteName(state.CipherSuite),
		Verified:      len(state.VerifiedChains) > 0,
		Phase:         mt.phase,
		LastObserved:  time.Now(),
	}

	// Calculate cert expiry days from peer certificates
	if len(state.PeerCertificates) > 0 {
		expiry := state.PeerCertificates[0].NotAfter
		days := int(time.Until(expiry).Hours() / 24)
		conn.CertExpiryDays = days
		if days < 0 {
			conn.Verified = false
			conn.Error = "certificate expired"
		} else if days < 30 {
			conn.Error = "certificate expiring soon"
		}
	}

	mt.connections[key] = conn
}

// RecordConnectionManual records an mTLS connection with manual parameters.
func (mt *MTLSStatusTracker) RecordConnectionManual(source, target, tlsVersion, cipherSuite string, verified bool, certExpiryDays int, err string) {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	key := source + "->" + target
	conn := &MTLSConnection{
		SourceService:  source,
		TargetService:  target,
		TLSVersion:     tlsVersion,
		CipherSuite:    cipherSuite,
		Verified:       verified,
		CertExpiryDays: certExpiryDays,
		Phase:          mt.phase,
		LastObserved:   time.Now(),
		Error:          err,
	}

	mt.connections[key] = conn
}

// GetConnection retrieves an mTLS connection status by source and target.
func (mt *MTLSStatusTracker) GetConnection(source, target string) (*MTLSConnection, bool) {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	key := source + "->" + target
	conn, ok := mt.connections[key]
	return conn, ok
}

// GetAllConnections returns all tracked mTLS connections.
func (mt *MTLSStatusTracker) GetAllConnections() []*MTLSConnection {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	result := make([]*MTLSConnection, 0, len(mt.connections))
	for _, conn := range mt.connections {
		result = append(result, conn)
	}
	return result
}

// GetConnectionsBySource returns connections originating from a service.
func (mt *MTLSStatusTracker) GetConnectionsBySource(source string) []*MTLSConnection {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	result := make([]*MTLSConnection, 0)
	for _, conn := range mt.connections {
		if conn.SourceService == source {
			result = append(result, conn)
		}
	}
	return result
}

// GetConnectionsByTarget returns connections targeting a service.
func (mt *MTLSStatusTracker) GetConnectionsByTarget(target string) []*MTLSConnection {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	result := make([]*MTLSConnection, 0)
	for _, conn := range mt.connections {
		if conn.TargetService == target {
			result = append(result, conn)
		}
	}
	return result
}

// GetUnverifiedConnections returns connections that failed verification.
func (mt *MTLSStatusTracker) GetUnverifiedConnections() []*MTLSConnection {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	result := make([]*MTLSConnection, 0)
	for _, conn := range mt.connections {
		if !conn.Verified {
			result = append(result, conn)
		}
	}
	return result
}

// GetExpiringConnections returns connections with certificates expiring within days.
func (mt *MTLSStatusTracker) GetExpiringConnections(withinDays int) []*MTLSConnection {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	result := make([]*MTLSConnection, 0)
	for _, conn := range mt.connections {
		if conn.CertExpiryDays >= 0 && conn.CertExpiryDays <= withinDays {
			result = append(result, conn)
		}
	}
	return result
}

// GetConnectionMatrix returns a matrix view of all service pairs.
func (mt *MTLSStatusTracker) GetConnectionMatrix() map[string]map[string]*MTLSConnection {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	matrix := make(map[string]map[string]*MTLSConnection)
	for _, conn := range mt.connections {
		if matrix[conn.SourceService] == nil {
			matrix[conn.SourceService] = make(map[string]*MTLSConnection)
		}
		matrix[conn.SourceService][conn.TargetService] = conn
	}
	return matrix
}

// RemoveConnection removes a connection from tracking.
func (mt *MTLSStatusTracker) RemoveConnection(source, target string) bool {
	mt.mu.Lock()
	defer mt.mu.Unlock()
	key := source + "->" + target
	if _, ok := mt.connections[key]; ok {
		delete(mt.connections, key)
		return true
	}
	return false
}

// tlsVersionString converts a TLS version constant to a string.
func tlsVersionString(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return "Unknown"
	}
}