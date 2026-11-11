package security

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"sync"
	"time"
)

// CertificateInfo represents a tracked certificate with lifecycle metadata.
type CertificateInfo struct {
	SerialNumber      string    `json:"serial_number"`
	Subject           string    `json:"subject"`
	Issuer            string    `json:"issuer"`
	NotBefore         time.Time `json:"not_before"`
	NotAfter          time.Time `json:"not_after"`
	Status            string    `json:"status"`              // "valid" | "expiring" | "expired" | "revoked"
	RotationDue       bool      `json:"rotation_due"`
	RevocationReason  string    `json:"revocation_reason,omitempty"`
	Source            string    `json:"source"`              // "vault" | "file" | "manual"
	Service           string    `json:"service"`             // Service this cert belongs to
	LastChecked       time.Time `json:"last_checked"`
	RotationThreshold time.Duration `json:"-"`               // When to mark as rotation due (default 30 days)
}

// CertificateTracker manages certificate lifecycle tracking.
type CertificateTracker struct {
	mu           sync.RWMutex
	certificates map[string]*CertificateInfo // keyed by serial number
	expiryAlert  time.Duration
}

// NewCertificateTracker creates a new certificate tracker.
func NewCertificateTracker(expiryAlert time.Duration) *CertificateTracker {
	if expiryAlert == 0 {
		expiryAlert = 30 * 24 * time.Hour // 30 days default
	}
	return &CertificateTracker{
		certificates: make(map[string]*CertificateInfo),
		expiryAlert:  expiryAlert,
	}
}

// ExpiryAlert returns the expiry alert duration.
func (ct *CertificateTracker) ExpiryAlert() time.Duration {
	ct.mu.RLock()
	defer ct.mu.RUnlock()
	return ct.expiryAlert
}

// TrackCertificate adds or updates a certificate in the tracker.
func (ct *CertificateTracker) TrackCertificate(ctx context.Context, cert *x509.Certificate, source, service string) *CertificateInfo {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	info := ct.certificateFromX509(cert, source, service)
	ct.certificates[info.SerialNumber] = info
	return info
}

// TrackCertificatePEM parses and tracks a certificate from PEM-encoded data.
func (ct *CertificateTracker) TrackCertificatePEM(ctx context.Context, pemData []byte, source, service string) (*CertificateInfo, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, ErrInvalidCertificate
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return ct.TrackCertificate(ctx, cert, source, service), nil
}

// GetCertificate retrieves a certificate by serial number.
func (ct *CertificateTracker) GetCertificate(serialNumber string) (*CertificateInfo, bool) {
	ct.mu.RLock()
	defer ct.mu.RUnlock()
	info, ok := ct.certificates[serialNumber]
	return info, ok
}

// GetAllCertificates returns all tracked certificates.
func (ct *CertificateTracker) GetAllCertificates() []*CertificateInfo {
	ct.mu.RLock()
	defer ct.mu.RUnlock()
	result := make([]*CertificateInfo, 0, len(ct.certificates))
	for _, info := range ct.certificates {
		result = append(result, info)
	}
	return result
}

// GetCertificatesByService returns certificates for a specific service.
func (ct *CertificateTracker) GetCertificatesByService(service string) []*CertificateInfo {
	ct.mu.RLock()
	defer ct.mu.RUnlock()
	result := make([]*CertificateInfo, 0)
	for _, info := range ct.certificates {
		if info.Service == service {
			result = append(result, info)
		}
	}
	return result
}

// GetExpiringCertificates returns certificates expiring within the alert window.
func (ct *CertificateTracker) GetExpiringCertificates() []*CertificateInfo {
	ct.mu.RLock()
	defer ct.mu.RUnlock()
	now := time.Now()
	alertWindow := now.Add(ct.expiryAlert)
	result := make([]*CertificateInfo, 0)
	for _, info := range ct.certificates {
		if info.NotAfter.Before(alertWindow) && info.NotAfter.After(now) && info.Status != "revoked" {
			result = append(result, info)
		}
	}
	return result
}

// GetExpiredCertificates returns certificates that have expired.
func (ct *CertificateTracker) GetExpiredCertificates() []*CertificateInfo {
	ct.mu.RLock()
	defer ct.mu.RUnlock()
	now := time.Now()
	result := make([]*CertificateInfo, 0)
	for _, info := range ct.certificates {
		if info.NotAfter.Before(now) && info.Status != "revoked" {
			result = append(result, info)
		}
	}
	return result
}

// GetRevokedCertificates returns revoked certificates.
func (ct *CertificateTracker) GetRevokedCertificates() []*CertificateInfo {
	ct.mu.RLock()
	defer ct.mu.RUnlock()
	result := make([]*CertificateInfo, 0)
	for _, info := range ct.certificates {
		if info.Status == "revoked" {
			result = append(result, info)
		}
	}
	return result
}

// RevokeCertificate marks a certificate as revoked.
func (ct *CertificateTracker) RevokeCertificate(serialNumber, reason string) bool {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	info, ok := ct.certificates[serialNumber]
	if !ok {
		return false
	}
	info.Status = "revoked"
	info.RevocationReason = reason
	info.LastChecked = time.Now()
	return true
}

// RefreshStatuses re-evaluates all certificate statuses based on current time.
func (ct *CertificateTracker) RefreshStatuses() {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	now := time.Now()
	alertWindow := now.Add(ct.expiryAlert)
	for _, info := range ct.certificates {
		info.LastChecked = now
		if info.Status == "revoked" {
			continue
		}
		if info.NotAfter.Before(now) {
			info.Status = "expired"
			info.RotationDue = true
		} else if info.NotAfter.Before(alertWindow) {
			info.Status = "expiring"
			info.RotationDue = true
		} else {
			info.Status = "valid"
			info.RotationDue = false
		}
	}
}

// GetRotationSchedule returns upcoming certificate rotations sorted by expiry.
func (ct *CertificateTracker) GetRotationSchedule() []*CertificateInfo {
	ct.mu.RLock()
	defer ct.mu.RUnlock()
	result := make([]*CertificateInfo, 0)
	for _, info := range ct.certificates {
		if info.Status != "revoked" && info.RotationDue {
			result = append(result, info)
		}
	}
	// Sort by NotAfter ascending (soonest first)
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].NotAfter.After(result[j].NotAfter) {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	return result
}

// RemoveCertificate removes a certificate from tracking.
func (ct *CertificateTracker) RemoveCertificate(serialNumber string) bool {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	if _, ok := ct.certificates[serialNumber]; ok {
		delete(ct.certificates, serialNumber)
		return true
	}
	return false
}

// certificateFromX509 converts an x509.Certificate to CertificateInfo.
func (ct *CertificateTracker) certificateFromX509(cert *x509.Certificate, source, service string) *CertificateInfo {
	return &CertificateInfo{
		SerialNumber:      cert.SerialNumber.String(),
		Subject:           cert.Subject.String(),
		Issuer:            cert.Issuer.String(),
		NotBefore:         cert.NotBefore,
		NotAfter:          cert.NotAfter,
		Status:            "valid",
		RotationDue:       false,
		Source:            source,
		Service:           service,
		LastChecked:       time.Now(),
		RotationThreshold: ct.expiryAlert,
	}
}

// ErrInvalidCertificate is returned when certificate data is invalid.
var ErrInvalidCertificate = &CertificateError{Message: "invalid certificate data"}

// CertificateError represents a certificate tracking error.
type CertificateError struct {
	Message string
}

func (e *CertificateError) Error() string {
	return "certificate error: " + e.Message
}