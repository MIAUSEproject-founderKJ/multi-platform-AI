//MIAUSEproject-founderKJ/multi-platform-AI/cmd/aios-node/main.go

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"multi-platform-AI/core"
	"multi-platform-AI/internal/logging"
	"multi-platform-AI/internal/watchdog"
)

func main() {
	// 1. Initial Logging
	logging.Info("[BOOT] Initializing StrataCore Node (AIOS-Node)...")

	// 2. Setup Watchdog
	// If the Bootstrap process hangs (e.g., waiting for a dead CAN-bus), 
	// the watchdog triggers the degradation protocol.
	wdt := watchdog.New(watchdog.Config{
		TimeoutSeconds: 5,
		OnFailure:      "degrade_to_safe_mode",
	})
	wdt.Start()

	// 3. Trigger the Nucleus Bootstrap
	// This initiates the entire Layer I: 
	// apppath -> probe -> classify -> attestation -> boot_manager
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	kernel, err := core.Bootstrap(ctx)
	if err != nil {
		logging.Error("[FATAL] Kernel failed to bootstrap: %v", err)
		// Hardware watchdog will likely handle the hardware side, 
		// but we exit cleanly from the software side.
		os.Exit(1)
	}

	// 4. Feed the Watchdog: Transition to Operational State
	wdt.Heartbeat()
	logging.Info("[BOOT] Kernel Active. Trust Level: %.2f", kernel.TrustLevel())

	// 5. Start HMI Lifecycle
	// The kernel now owns the HMI pipe we built earlier.
	go kernel.RunHMILoop()

	// 6. Wait for Shutdown Signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	logging.Info("[SYSTEM] Node fully operational. Standing by for instructions.")

	<-stop

	// 7. Graceful Shutdown: The "Secure Lock"
	logging.Info("[SHUTDOWN] Locking Vault and terminating Vision Streams...")
	kernel.Shutdown()
	logging.Info("[SHUTDOWN] Node offline.")
}