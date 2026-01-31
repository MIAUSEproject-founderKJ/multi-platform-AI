//MIAUSEproject-founderKJ/multi-platform-AI/cmd/aios-node/main.go

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"multi-platform-AI/core"
	"multi-platform-AI/internal/watchdog"
)

func main() {
	// 1. Initialize Minimalist Logging for the Boot Sequence
	log.Println("[BOOT] Initializing StrataCore Node...")

	// 2. Start Hardware/Software Watchdog
	// This ensures that if the boot process hangs for more than 5 seconds,
	// the system will auto-restart or enter Safe Mode.
	wdt := watchdog.New(watchdog.Config{
		TimeoutSeconds: 5,
		OnFailure:      "degrade_to_safe_mode",
	})
	wdt.Start()

	// 3. Hand over control to the Secure Nucleus (Layer I)
	// This is where PROBE -> CLASSIFY -> ATTEST happens.
	kernel, err := core.Bootstrap()
	if err != nil {
		log.Fatalf("[FATAL] Kernel failed to bootstrap: %v", err)
		// At this point, the Nucleus has failed; core/platform/degrade logic takes over.
		os.Exit(1)
	}

	// 4. Feed the Watchdog: System is now healthy
	wdt.Heartbeat()
	log.Println("[BOOT] Kernel started successfully. Autonomy Level:", kernel.TrustLevel())

	// 5. Wait for Shutdown Signal (SIGTERM, SIGINT)
	// This prevents the app from closing immediately and handles graceful cleanup.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("[SHUTDOWN] Terminating services and locking hardware...")
	kernel.Shutdown()
}