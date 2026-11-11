package security

import (
	"context"
	"sync"
	"time"
)

// TrustBoundary represents a trust boundary with mTLS compliance tracking.
type TrustBoundary struct {
	Name           string    `json:"name"`
	Services       []string  `json:"services"`
	MTLSRequired   bool      `json:"mtls_required"`
	Phase          int       `json:"phase"`
	Compliant      bool      `json:"compliant"`
	Violations     []string  `json:"violations,omitempty"`
	LastEvaluated  time.Time `json:"last_evaluated"`
	TotalPairs     int       `json:"total_pairs"`
	CompliantPairs int       `json:"compliant_pairs"`
}

// TrustBoundaryTracker manages trust boundary compliance tracking.
type TrustBoundaryTracker struct {
	mu             sync.RWMutex
	boundaries     map[string]*TrustBoundary // keyed by boundary name
	mtlsTracker    *MTLSStatusTracker
	certTracker    *CertificateTracker
}

// NewTrustBoundaryTracker creates a new trust boundary tracker.
func NewTrustBoundaryTracker(mtlsTracker *MTLSStatusTracker, certTracker *CertificateTracker) *TrustBoundaryTracker {
	return &TrustBoundaryTracker{
		boundaries:  make(map[string]*TrustBoundary),
		mtlsTracker: mtlsTracker,
		certTracker: certTracker,
	}
}

// DefineBoundary defines a new trust boundary.
func (tb *TrustBoundaryTracker) DefineBoundary(ctx context.Context, name string, services []string, mtlsRequired bool, phase int) *TrustBoundary {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	if phase < 1 {
		phase = 1
	}
	if phase > 3 {
		phase = 3
	}

	boundary := &TrustBoundary{
		Name:         name,
		Services:     services,
		MTLSRequired: mtlsRequired,
		Phase:        phase,
		Compliant:    false,
		Violations:   []string{},
	}

	tb.boundaries[name] = boundary
	return boundary
}

// GetBoundary retrieves a trust boundary by name.
func (tb *TrustBoundaryTracker) GetBoundary(name string) (*TrustBoundary, bool) {
	tb.mu.RLock()
	defer tb.mu.RUnlock()
	boundary, ok := tb.boundaries[name]
	return boundary, ok
}

// GetAllBoundaries returns all trust boundaries.
func (tb *TrustBoundaryTracker) GetAllBoundaries() []*TrustBoundary {
	tb.mu.RLock()
	defer tb.mu.RUnlock()
	result := make([]*TrustBoundary, 0, len(tb.boundaries))
	for _, boundary := range tb.boundaries {
		result = append(result, boundary)
	}
	return result
}

// EvaluateCompliance evaluates mTLS compliance for all boundaries.
func (tb *TrustBoundaryTracker) EvaluateCompliance(ctx context.Context) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	for _, boundary := range tb.boundaries {
		tb.evaluateBoundary(boundary)
	}
}

// evaluateBoundary evaluates a single boundary's compliance.
func (tb *TrustBoundaryTracker) evaluateBoundary(boundary *TrustBoundary) {
	boundary.LastEvaluated = time.Now()
	boundary.Violations = []string{}

	if !boundary.MTLSRequired {
		boundary.Compliant = true
		boundary.TotalPairs = 0
		boundary.CompliantPairs = 0
		return
	}

	// Check all service pairs within the boundary
	services := boundary.Services
	totalPairs := 0
	compliantPairs := 0

	for _, source := range services {
		for _, target := range services {
			if source == target {
				continue
			}
			totalPairs++

			conn, ok := tb.mtlsTracker.GetConnection(source, target)
			if !ok {
				boundary.Violations = append(boundary.Violations,
					"no mTLS connection observed: "+source+"->"+target)
				continue
			}

			if !conn.Verified {
				boundary.Violations = append(boundary.Violations,
					"mTLS verification failed: "+source+"->"+target+" ("+conn.Error+")")
				continue
			}

			// Check certificate expiry
			if conn.CertExpiryDays >= 0 && conn.CertExpiryDays <= 30 {
				boundary.Violations = append(boundary.Violations,
					"certificate expiring soon: "+source+"->"+target+" ("+string(rune(conn.CertExpiryDays))+" days)")
			}

			if conn.Phase < boundary.Phase {
				boundary.Violations = append(boundary.Violations,
					"phase mismatch: "+source+"->"+target+" is phase "+string(rune(conn.Phase))+" but boundary requires phase "+string(rune(boundary.Phase)))
			}

			compliantPairs++
		}
	}

	boundary.TotalPairs = totalPairs
	boundary.CompliantPairs = compliantPairs
	boundary.Compliant = len(boundary.Violations) == 0
}

// GetCompliantBoundaries returns boundaries that are fully compliant.
func (tb *TrustBoundaryTracker) GetCompliantBoundaries() []*TrustBoundary {
	tb.mu.RLock()
	defer tb.mu.RUnlock()
	result := make([]*TrustBoundary, 0)
	for _, boundary := range tb.boundaries {
		if boundary.Compliant {
			result = append(result, boundary)
		}
	}
	return result
}

// GetNonCompliantBoundaries returns boundaries with violations.
func (tb *TrustBoundaryTracker) GetNonCompliantBoundaries() []*TrustBoundary {
	tb.mu.RLock()
	defer tb.mu.RUnlock()
	result := make([]*TrustBoundary, 0)
	for _, boundary := range tb.boundaries {
		if !boundary.Compliant {
			result = append(result, boundary)
		}
	}
	return result
}

// GetBoundaryByService returns boundaries that contain a service.
func (tb *TrustBoundaryTracker) GetBoundaryByService(service string) []*TrustBoundary {
	tb.mu.RLock()
	defer tb.mu.RUnlock()
	result := make([]*TrustBoundary, 0)
	for _, boundary := range tb.boundaries {
		for _, svc := range boundary.Services {
			if svc == service {
				result = append(result, boundary)
				break
			}
		}
	}
	return result
}

// UpdateBoundaryPhase updates a boundary's required phase.
func (tb *TrustBoundaryTracker) UpdateBoundaryPhase(name string, phase int) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	boundary, ok := tb.boundaries[name]
	if !ok {
		return false
	}

	if phase < 1 {
		phase = 1
	}
	if phase > 3 {
		phase = 3
	}

	boundary.Phase = phase
	return true
}

// RemoveBoundary removes a trust boundary.
func (tb *TrustBoundaryTracker) RemoveBoundary(name string) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	if _, ok := tb.boundaries[name]; ok {
		delete(tb.boundaries, name)
		return true
	}
	return false
}