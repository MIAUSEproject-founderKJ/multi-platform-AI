//cmd/aios/run_http.go

package main

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// ============================================================
// HTTP (EDGE ONLY)
// ============================================================

func (a *App) startHTTP() {

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		status := a.supervisor.HealthStatus()

		code := http.StatusOK
		if !status.Healthy {
			code = http.StatusServiceUnavailable
		}

		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(status)
	})

	a.server = &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.Error("HTTP_ERROR", zap.Error(err))
		}
	}()
}
