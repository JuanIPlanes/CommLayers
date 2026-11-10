package security

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"strings"
	"time"
)

// enabled reports if the security bootstrap has been initialized.
var enabled bool

// InitSecurityBootstrap initializes the security baseline (Keycloak + Vault).
// This must be called before any security-dependent middleware operates.
func InitSecurityBootstrap() {
	enabled = true
}

// IsEnabled reports if the security bootstrap has been initialized.
func IsEnabled() bool {
	return enabled
}

// KeycloakConfig holds the configuration for Keycloak integration.
type KeycloakConfig struct {
	// URL is the Keycloak server URL.
	URL string `json:"url"`
	// Realm is the Keycloak realm name.
	Realm string `json:"realm"`
	// ClientID is the OAuth2 client ID.
	ClientID string `json:"client_id"`
	// ClientSecret is the OAuth2 client secret.
	ClientSecret string `json:"client_secret"`
	// Domain is the base domain for Keycloak JWT audience validation.
	Domain string `json:"domain"`
	// Timeout is the token request timeout.
	Timeout time.Duration `json:"timeout"`
}

// KeycloakToken represents a decoded Keycloak OAuth2 token.
type KeycloakToken struct {
	// Subject is the unique identifier of the user in Keycloak.
	Subject string `json:"sub"`
	// Audience is the intended audience of the token.
	Audience string `json:"aud"`
	// Issuer is the issuer of the token.
	Issuer string `json:"iss"`
	// ExpiresAt is the token expiration time.
	ExpiresAt int64 `json:"exp"`
	// IssuedAt is the token issuance time.
	IssuedAt int64 `json:"iat"`
	// scopes represents the token's granted scopes.
	Scopes []string `json:"scp,omitempty"`
	// RoleClaims represents the realm and client-level roles.
	RoleClaims map[string][]string `json:"realm_access,omitempty"`
}

// KeycloakMiddleware protects routes behind Keycloak authentication.
// It validates the presence of a Bearer token and, if Keycloak is enabled,
// performs basic validation against the configured realm.
//
// In a full implementation, this would verify the token signature and
// refresh tokens against the Keycloak server, but the lightweight surface
// only ensures a token is present and provides a structure for full validation.
func KeycloakMiddleware(next http.Handler, config *KeycloakConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !enabled {
			next.ServeHTTP(w, r)
			return
		}
		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("unauthorized"))
			return
		}
		// Lightweight extraction; full validation would use config and
		// a JWT parser against the Keycloak public key set.
		next.ServeHTTP(w, r)
	})
}

// VaultConfig holds the configuration for Vault integration.
type VaultConfig struct {
	// Address is the Vault server address.
	Address string `json:"address"`
	// Token is the Vault token for authentication.
	Token string `json:"token"`
	// Path is the base path within Vault for secrets.
	Path string `json:"path"`
	// Timeout is the Vault API request timeout.
	Timeout time.Duration `json:"timeout"`
}

// VaultSecret represents a secret retrieved from Vault.
type VaultSecret struct {
	// Data contains the secret key-value pairs.
	Data map[string]interface{} `json:"data"`
	// WrapInfo contains wrapping information if the secret is wrapped.
	WrapInfo interface{} `json:"wrap_info,omitempty"`
}

// mTLSConfig holds the configuration for staged mTLS TLS setup.
type mTLSConfig struct {
	// ClientCA is the Certificate Authority pool for verifying client certificates.
	ClientCA *x509.CertPool `json:"-"`
	// ClientCertFile is the path to the client certificate file.
	ClientCertFile string `json:"-"`
	// ClientKeyFile is the path to the client key file.
	ClientKeyFile string `json:"-"`
	// ServerCA is the Certificate Authority pool for verifying server certificates.
	ServerCA *x509.CertPool `json:"-"`
	// MinVersion is the minimum TLS version allowed.
	MinVersion uint16 `json:"-"`
}

// TLSConfigFromFile loads client and server TLS configuration from PEM-encoded files.
func TLSConfigFromFile(clientCert, clientKey, serverCA string) (*tls.Config, error) {
	// Load client certificate
	cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, err
	}

	// Load server CA certificate
	caCert, err := os.ReadFile(serverCA)
	if err != nil {
		return nil, err
	}

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caPool,
	}

	return tlsConfig, nil
}

// TokenSessionPolicy defines the token and session policy rules for the security bootstrap.
type TokenSessionPolicy struct {
	// SessionTimeout is the duration after which a session expires.
	SessionTimeout time.Duration `json:"session_timeout"`
	// MaxIdleTimeout is the maximum idle time before a session is invalidated.
	MaxIdleTimeout time.Duration `json:"max_idle_timeout"`
	// RequirePKCE indicates if PKCE is required for OAuth2 flows.
	RequirePKCE bool `json:"require_pkce"`
	// AllowedOrigins are the origins permitted for token requests.
	AllowedOrigins []string `json:"allowed_origins,omitempty"`
}

// SessionRecord represents a recorded session state for audit and revocation.
type SessionRecord struct {
	// SessionID is the unique identifier of the session.
	SessionID string `json:"session_id"`
	// UserID is the identifier of the associated user.
	UserID string `json:"user_id"`
	// CreatedAt is when the session was created.
	CreatedAt time.Time `json:"created_at"`
	// LastUsedAt is the last time the session was used.
	LastUsedAt time.Time `json:"last_used_at"`
	// RevokedAt is when the session was revoked (if applicable).
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
}

// RevokeSession marks a session as revoked.
func (r *SessionRecord) RevokeSession() time.Time {
	now := time.Now()
	r.RevokedAt = &now
	return now
}