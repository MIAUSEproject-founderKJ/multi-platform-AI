//MIAUSEproject-founderKJ/multi-platform-AI/cmd/aios-kernel/main.go

package main

import (
	"multi-platform-AI/core"
	"multi-platform-AI/internal/logging"
	"multi-platform-AI/runtime/monitor"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logging.Info("--- [STRATACORE KERNEL NUCLEUS ACTIVE] ---")

	// 1. MONITOR: Initialize performance capping (VRAM/Thermal)
	mon := monitor.NewPerformanceMonitor()
	mon.Start()

	// 2. ORCHESTRATION: Load Policy & Trust Evaluators
	nucleus, err := core.InitializeNucleus()
	if err != nil {
		logging.Error("[FATAL] Nucleus initialization failed: %v", err)
		os.Exit(1)
	}

	// 3. DREAM STATE: Start Idle Simulation Logic
	// This allows the kernel to retrain when the Node reports IDLE status
	go nucleus.ManageLifecycle()

	// 4. PERSISTENCE: Block until manual stop
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	logging.Info("[SYSTEM] Nucleus operational. Managing distributed nodes.")
	<-stop

	nucleus.SyncAndHalt()
}