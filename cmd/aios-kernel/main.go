//cmd/aios-kernel/main.go

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

    // 1. Start Monitor
    mon := monitor.NewPerformanceMonitor()
    mon.Start()
    defer mon.Stop()

    // 2. Initialize Nucleus with Monitor Link
    // Note: You'll need to update core.InitializeNucleus to accept *monitor.PerformanceMonitor
    nucleus, err := core.InitializeNucleus(mon) 
    if err != nil {
        logging.Error("[FATAL] Nucleus initialization failed: %v", err)
        os.Exit(1)
    }

    // 3. Start Lifecycle
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    go nucleus.RunLifecycle(ctx) // Use the standard naming from your kernel.go

    // 4. Signal Catching
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

    logging.Info("[SYSTEM] Nucleus operational. Monitoring hardware...")
    <-stop

    // 5. Cleanup
    cancel() // Stop the lifecycle loop first
    nucleus.Shutdown()
    logging.Info("[SYSTEM] Nucleus offline.")
}


/*
This entry point runs the system as infrastructure.

Key properties:

• No UI
• No user login
• No interactive flow
• Pure lifecycle orchestration

Its job:

Start monitor (performance telemetry)

Initialize nucleus

Run lifecycle manager

Catch signals

Shutdown cleanly

This binary is appropriate for:

Industrial embedded systems

Fleet-managed vehicle controllers

Always-on factory nodes

Edge compute deployments

Robotaxi backend nodes

It treats the kernel as an autonomous system component.

The conceptual model is:

System = Service
Not = Application

It is analogous to:

systemd-managed service

Windows service

Daemon process

This is infrastructure mode.

*/