//MIAUSEproject-founderKJ/multi-platform-AI/internal/logger/bloat_guard.go

package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"project/internal/apppath"
	"sync"

	"gopkg.in/natefinch/lumberjack.v2" // The industry standard for rotation
)

type BloatGuard struct {
	mu         sync.Mutex
	MaxSizeMB  int
	MaxBackups int
	Compress   bool
}

// Initialize structured logging with Anti-Bloat guards
func Init(guard BloatGuard) {
	logDir := apppath.GetLogDir()
	
	// Ensure log directory exists
	_ = os.MkdirAll(logDir, 0755)

	rotator := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, "system.log"),
		MaxSize:    guard.MaxSizeMB,  // e.g., 50MB
		MaxBackups: guard.MaxBackups, // e.g., 3 files
		MaxAge:     7,                // Days
		Compress:   guard.Compress,   // Gzip old logs
	}

	// Multi-writer: Logs go to Terminal (if dev) and the Rotator (always)
	var writers []io.Writer
	writers = append(writers, rotator)
	if os.Getenv("AIOS_ENV") == "development" {
		writers = append(writers, os.Stdout)
	}

	// Set as the global logger for the entire Nucleus
	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(writers...), nil))
	slog.SetDefault(logger)
}