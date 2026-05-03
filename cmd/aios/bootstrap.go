// =========================
// PRODUCTION-GRADE MAIN.GO - cmd/aios/bootstrap.go
// Boot → ExecutionContext → RuntimeContext → Modules → Supervisor → Interfaces
// =========================
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// ============================================================
// ENTRYPOINT
// ============================================================

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log, _ := zap.NewProduction()
	defer log.Sync()

	sys, err := buildSystemContext()
	if err != nil {
		log.Fatal("BOOT_FAILED", zap.Error(err))
	}

	app, err := buildApp(log, sys)
	if err != nil {
		log.Fatal("APP_BUILD_FAILED", zap.Error(err))
	}

	if err := app.Start(ctx); err != nil {
		log.Fatal("START_FAILED", zap.Error(err))
	}

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_ = app.Stop(shutdownCtx)

	log.Info("SYSTEM_EXIT", zap.String("user", sys.Session.Identity.Username))
}
