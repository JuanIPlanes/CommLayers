package router

import (
	"encoding/json"
	"errors"
	"net/http"
	"commlayers/services/api/internal/transport"
	"github.com/gorilla/mux"
)

// BuildTransportMux registers transport‑family endpoints.
// It expects a ready TransportProvider implementation.
func BuildTransportMux(provider transport.TransportProvider) *mux.Router {
    r := mux.NewRouter()
    // Create transport
    r.HandleFunc("/api/transport", func(w http.ResponseWriter, r *http.Request) {
        var t transport.Transport
        if err := decodeJSON(r, &t); err != nil {
            writeError(w, http.StatusBadRequest, err)
            return
        }
        if err := provider.Create(r.Context(), &t); err != nil {
            writeError(w, http.StatusInternalServerError, err)
            return
        }
        writeJSON(w, http.StatusCreated, t)
    }).Methods(http.MethodPost)

    // Get by ID
    r.HandleFunc("/api/transport/{id}", func(w http.ResponseWriter, r *http.Request) {
        id := mux.Vars(r)["id"]
        t, err := provider.GetByID(r.Context(), id)
        if err != nil {
            writeError(w, http.StatusInternalServerError, err)
            return
        }
        if t == nil {
            writeError(w, http.StatusNotFound, errors.New("not found"))
            return
        }
        writeJSON(w, http.StatusOK, t)
    }).Methods(http.MethodGet)

    // Search (simple query param)
    r.HandleFunc("/api/transport/search", func(w http.ResponseWriter, r *http.Request) {
        q := r.URL.Query().Get("q")
        results, err := provider.Search(r.Context(), q)
        if err != nil {
            writeError(w, http.StatusInternalServerError, err)
            return
        }
        writeJSON(w, http.StatusOK, results)
    }).Methods(http.MethodGet)

    return r
}

// Helper utilities (simplified for demo)
func decodeJSON(r *http.Request, v interface{}) error {
    return json.NewDecoder(r.Body).Decode(v)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, err error) {
    writeJSON(w, status, map[string]string{"error": err.Error()})
}
