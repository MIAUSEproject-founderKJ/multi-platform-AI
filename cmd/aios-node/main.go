// cmd/aios-node/main.go

package main

import (
    "context"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core"
    "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
    "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/watchdog"
    
    // Import modules for side-effect registration
    _ "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/autonomous"
    _ "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/hmi"
    _ "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/simulation"
)

func main() {
	logging.Info("[BOOT] Initializing AIOS Nucleus...")

	// 1. Hardware Watchdog
	wdt := watchdog.New(watchdog.Config{TimeoutSeconds: 5, OnFailure: "safe_park"})
	wdt.Start()

	// 2. Bootstrap Kernel (Machine Level)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kernel, err := core.Bootstrap(ctx)
	if err != nil {
		logging.Error("[FATAL] Bootstrap failed: %v", err)
		os.Exit(1)
	}

	// 3. User Interface Initialization (Session Level)
	ui := hmi.NewTerminal(kernel)

	logging.Info("[SYSTEM] Kernel Ready. Handing control to User Interface.")
	
	// --- BLOCKING LOGIN FLOW ---
	if err := ui.StartLoginFlow(); err != nil {
		logging.Error("Login aborted: %v", err)
		kernel.Shutdown()
		os.Exit(0)
	}

	// 4. Operational Command Loop
	// Watchdog is fed here or via a separate ticker
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				wdt.Heartbeat()
			}
		}
	}()

	ui.RunCommandLoop()

	// 5. Shutdown
	logging.Info("[SHUTDOWN] User initiated exit.")
	wdt.Stop()
	kernel.Shutdown()
}

/*
This binary:

Starts watchdog (safety layer)

Bootstraps kernel (machine level)

Initializes HMI

Forces login flow

Runs command loop

Feeds watchdog

Handles shutdown

This is:

Application Mode
User-Facing Node

It assumes:

A human or operator exists

There is session state

There is identity negotiation

Commands are interactive

Use cases:

Workstation

Maintenance terminal

Robotaxi passenger HMI

Developer console

Simulation environment

It is not a daemon. It is an agent.
*/