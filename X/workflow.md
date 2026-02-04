Workflow:
When the system is booted/device is started. 
1. cmd/aios-kernel/main.go initiated (focuses on policy orchestration and "Dream State", manages the high-latency Cognition and Simulation layers) and cmd/aios-node/main.go (focuses on physical HAL and HMI)

2. from cmd/aios-node/main.go, watchdog is set, then call core.Bootstrap in core/kernel.go (manages the flow between Platform Identity, Security, and Trust). 
3. from core/kernel.go, call platform.RunBootSequence in core/platform/boot.go, firstly, check "FirstBootMarker" by calling v.IsMissingMarker in core/security/vault.go. 
4. from core/security/vault.go, return the value of os.IsNotExist(err) to core/platform/boot.go.
5. log os.IsNotExist(err), then call PassiveScan() in core/platform/probe/passive.go.
6. from passive.go, call getMachineUUID, switch runtime.GOOS ("windows/linux/default unknown"), then call classifyPlatform() to infer which platform class 
7. then back to core/platform/boot.go, 

8. from core/platform/boot.go, receive the PlatformClass and call ManageBoot() in core/platform/boot_manager.go to handle the transition logic.

9. from boot_manager.go, evaluate the FirstBootMarker result:

    A. If Missing (Cold Boot): Initialize perception.NewVisionStream() and monitor.NewVitalsMonitor() to start the Reflective HUD immediately.

    B. If Present (Fast Boot): Attempt a Vault resume for low-latency startup.

10. from boot_manager.go, call probe.AggressiveScan() in core/platform/probe/aggressive.go to activate high-power hardware (Lidar, CAN-bus controllers, and Specialized AI chips).

11. from aggressive.go, return a FullProfile to the BootManager, which then triggers security.VerifyEnvironment() in core/security/attestation.go.

12. from attestation.go, calculate the Binary Self-Hash and compare it against the "Golden Hash" in the Vault; if valid, return a "Strong" attestation signal back to the BootManager.

13. from boot_manager.go, pass the FullProfile and Attestation results to policy.Evaluate() in core/policy/trust_bayesian.go.

14. from trust_bayesian.go, calculate the Q16 Trust Score by merging hardware health, signal latency (Pulse), and security validity; return this to the BootManager.

15. from boot_manager.go, finalize the boot:

16. Call navigation. InitializeSLAM() to begin 3D mapping.

17. Update the Reflective HUD to display "TRUST_LEVEL: SECURE" and the "Dream State" wireframe.

18. Call v.WriteMarker("FirstBootMarker") if this was a cold boot to seal the identity.

19. back to core/kernel.go, return the fully initialized Kernel object to cmd/aios-node/main.go, signaling that the Nucleus is now operational.

20. from aios-node/main.go, feed the Watchdog for the first time to clear the boot-timeout and enter the main operational loop.

21. from core/kernel.go (Background), monitor for IDLE state; if detected, signal simulation/digital_twin/voxel_world.go to enter the "Dream State" for neural rehearsal and policy distillation.