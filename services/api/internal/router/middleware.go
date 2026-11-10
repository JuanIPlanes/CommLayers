package router

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"commlayers/services/api/internal/backend"
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

		log.Printf("method=%s path=%s duration=%dms", r.Method, r.URL.Path, duration)
	})
}

func (bm *BenchmarkMiddleware) Collector() *backend.BenchmarkCollector {
	return bm.collector
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
