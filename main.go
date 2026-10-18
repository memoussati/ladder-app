// ladder-app: a tiny HTTP service used throughout the CI/CD ladder.
// Endpoints: GET / (hello), GET /healthz (200 ok), GET /version (build metadata as JSON).
// LADDER_FAIL_RATE=0.5 makes half of the / requests return 500 (used for canary drills).
package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Set with -ldflags "-X main.version=... -X main.commit=... -X main.date=...".
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type versionInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
	Host    string `json:"host"`
}

func failRate() float64 {
	f, err := strconv.ParseFloat(os.Getenv("LADDER_FAIL_RATE"), 64)
	if err != nil || f < 0 || f > 1 {
		return 0
	}
	return f
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	host, _ := os.Hostname()
	rate := failRate()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("GET /version", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(versionInfo{Version: version, Commit: commit, Date: date, Host: host})
	})
	mux.HandleFunc("GET /", func(w http.ResponseWriter, _ *http.Request) {
		if rate > 0 && rand.Float64() < rate { //nolint:gosec // not security sensitive
			http.Error(w, "injected failure", http.StatusInternalServerError)
			return
		}
		w.Header().Set("X-Ladder-Version", version)
		_, _ = w.Write([]byte("hello from ladder-app " + version + "\n"))
	})
	return mux
}

func main() {
	addr := ":" + envOr("PORT", "8080")
	srv := &http.Server{Addr: addr, Handler: newMux(), ReadHeaderTimeout: 5 * time.Second}
	log.Printf("ladder-app %s (%s) listening on %s", version, commit, addr)
	log.Fatal(srv.ListenAndServe())
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
