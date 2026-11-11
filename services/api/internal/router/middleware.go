package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"commlayers/services/api/internal/backend"
	"commlayers/services/api/internal/evidence"
)

type BenchmarkMiddleware struct {
	collector *backend.BenchmarkCollector
	mu        sync.Mutex
}

func NewBenchmarkMiddleware() *BenchmarkMiddleware {
	return &BenchmarkMiddleware{
		collector: backend.NewBenchmarkCollector(backend.DefaultBenchmarkOptions()),
	}
}

func (bm *BenchmarkMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start).Milliseconds()

		// Record latency for benchmark
		bm.mu.Lock()
		bm.collector.RecordTransition(float64(duration))
		bm.mu.Unlock()

		// Emit evidence contract for this request
		bm.emitRequestEvidence(r, duration)

		log.Printf("method=%s path=%s duration=%dms", r.Method, r.URL.Path, duration)
	})
}

func (bm *BenchmarkMiddleware) emitRequestEvidence(r *http.Request, durationMs int64) {
	runID := fmt.Sprintf("req_%s_%d", r.Method, time.Now().UnixNano())
	benchEv := evidence.NewBenchmarkEvidence(runID, 0, 0, float64(durationMs), 1)
	
	// Determine success based on response status (would need response capture)
	// For now, assume success for all completed requests
	_ = evidence.NewSuccessContract(benchEv.Evidence, evidence.SuccessOK)
}

func (bm *BenchmarkMiddleware) Collector() *backend.BenchmarkCollector {
	return bm.collector
}

// WithEvidenceMiddleware adds evidence emission to the request chain.
func WithEvidenceMiddleware(collector *backend.BenchmarkCollector) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// Wrap ResponseWriter to capture status
			ww := &evidenceResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			
			next.ServeHTTP(ww, r)
			
			duration := time.Since(start).Milliseconds()
			
			// Record to collector
			collector.RecordTransition(float64(duration))
			
			// Emit evidence contract
			runID := fmt.Sprintf("req_%s_%d", r.Method, time.Now().UnixNano())
			benchEv := evidence.NewBenchmarkEvidence(runID, 0, 0, float64(duration), 1)
			
			var ec *evidence.EvidenceContract
			if ww.statusCode >= 200 && ww.statusCode < 400 {
				ec = evidence.NewSuccessContract(benchEv.Evidence, evidence.SuccessOK)
			} else if ww.statusCode == 429 {
				ec = evidence.NewErrorContract(benchEv.Evidence, evidence.ErrorRateLimited, "rate limited")
			} else if ww.statusCode == 401 || ww.statusCode == 403 {
				ec = evidence.NewErrorContract(benchEv.Evidence, evidence.ErrorAuthFailed, "authentication failed")
			} else if ww.statusCode >= 400 && ww.statusCode < 500 {
				ec = evidence.NewErrorContract(benchEv.Evidence, evidence.ErrorInvalidInput, "client error")
			} else {
				ec = evidence.NewErrorContract(benchEv.Evidence, evidence.ErrorInternal, "server error")
			}
			
			// Log evidence for observability
			evJSON, _ := json.Marshal(ec)
			log.Printf("evidence=%s", string(evJSON))
		})
	}
}

// evidenceResponseWriter wraps http.ResponseWriter to capture status code.
type evidenceResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *evidenceResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func WithCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS_ALLOW_ORIGIN"))
		if w.Header().Get("Access-Control-Allow-Origin") == "" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("method=%s path=%s duration=%s", r.Method, r.URL.Path, time.Since(started))
	})
}
