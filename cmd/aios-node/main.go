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

    // 1. Hardware Watchdog (Safety Net)
    wdt := watchdog.New(watchdog.Config{ TimeoutSeconds: 5, OnFailure: "safe_park" })
    wdt.Start()

    // 2. Bootstrap the Kernel
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    kernel, err := core.Bootstrap(ctx)
    if err != nil {
        logging.Error("[FATAL] Bootstrap failed: %v", err)
        os.Exit(1)
    }

    // 3. Operational Loop
    // The Kernel doesn't "run" a loop anymore. The MODULES run the loops.
    // The Kernel just sits here holding the context.
    logging.Info("[SYSTEM] System Operational. Context: %s", kernel.Runtime.Platform.Name)
    
    // 4. Feed Watchdog
    go func() {
        ticker := time.NewTicker(1 * time.Second)
        for {
            select {
            case <-ctx.Done():
                return
            case <-ticker.C:
                // Only feed if the EventBus is healthy (example logic)
                wdt.Heartbeat()
            }
        }
    }()

    // 5. Wait for SIGTERM
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
    <-stop

    // 6. Graceful Exit
    logging.Info("[SHUTDOWN] Signal received.")
    wdt.Stop()
    kernel.Shutdown() // Cascades Stop() to all loaded modules
    logging.Info("[SHUTDOWN] Bye.")
}