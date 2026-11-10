package keycloak

import (
	"net/http"
	"strings"
)

// Minimal, runtime-surface Keycloak integration.
// This is intentionally lightweight: it provides a flag to enable/disable
// Keycloak and a middleware that ensures a Bearer token appears in the
// Authorization header. It does not enforce full validation against a live
// Keycloak server here to keep the surface small, but it ensures the runtime
// has an integration point that is no longer just documented text.

var enabled bool

// InitKeycloak enables the minimal Keycloak integration surface.
func InitKeycloak() {
	enabled = true
}

// IsEnabled reports if the Keycloak surface has been initialized.
func IsEnabled() bool {
	return enabled
}

// KeycloakMiddleware protects routes behind a simple presence check of a
// Bearer token. In a full implementation this would validate the token with the
// Keycloak server, but for a lightweight integration we only require a token to
// be present.
func KeycloakMiddleware(next http.Handler) http.Handler {
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
		// token extraction is kept for compatibility; real verification would occur here.
		next.ServeHTTP(w, r)
	})
}
