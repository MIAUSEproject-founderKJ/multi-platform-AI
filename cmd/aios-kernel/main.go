//MIAUSEproject-founderKJ/multi-platform-AI/cmd/aios-kernel/main.go

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/monitor"
)

func main() {
	logging.Info("--- [STRATACORE KERNEL NUCLEUS ACTIVE] ---")

	mon := monitor.NewPerformanceMonitor()
	mon.Start()
	defer mon.Stop() // Cleanup on exit

	nucleus, err := core.InitializeNucleus()
	if err != nil {
		logging.Error("[FATAL] Nucleus initialization failed: %v", err)
		os.Exit(1)
	}

	// Runs ManageLifecycle in a separate routine to allow the main thread to block on signals
	go nucleus.ManageLifecycle()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	logging.Info("[SYSTEM] Nucleus operational. Managing distributed nodes.")
	<-stop

	nucleus.SyncAndHalt()
	logging.Info("[SYSTEM] Nucleus offline.")
}
