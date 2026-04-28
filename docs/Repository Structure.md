//X/Repository Structure.md
1. Top-Level Repository Structure (Anti-Bloat Oriented)
Layered configs prevent redundancy:
schema → Platform type → Platform instance → User → Runtime.
Only differences are written at each layer.
Platform detection first: core/platform handles first bootstrap detection and classification. The system then loads the correct platform type and instance configs.
Scalability: Adding a new platform only requires a new type or instance YAML, no need to copy common configs.

multi-platform-ai/

api
api\commands
api\commands\command_contract.go
api\gen
api\hmi
api\hmi\hmi_contract.go
api\hmi\state_update.go
api\proto
api\proto\hmi.proto
api\proto\kernel.proto
api\proto\perception.proto
cmd\aios
cmd\aios\vault
cmd\aios\vault\credentials_.json
cmd\aios\vault\machine_first_boot_marker.json
cmd\aios\vault\users_hkj.json
cmd\aios\aios.exe
cmd\aios\bootstrap.go
core
core\agent
core\agent\agent_optimization_service.go
core\agent\agent_runtime_engine.go
core\auth
core\auth\auth_gatekeeper.go
core\auth\auth_service.go
core\policy
core\policy\policy_resolver.go
core\router
core\router\command_handler_contract.go
core\router\command_router.go
core\security
core\security\decision
core\security\decision\decision_engine.go
core\security\decision\permission_deriver.go
core\security\identity
core\security\identity\token_service.go
core\security\measurement
core\security\measurement\measured_boot.go
core\security\persistence
core\security\persistence\config_store.go
core\security\persistence\golden_hash_store.go
core\security\persistence\key_manager.go
core\security\persistence\kv_store.go
core\security\persistence\marker_store.go
core\security\persistence\vault_store.go
core\security\verification
core\security\verification\verification_engine.go
drivers
drivers\audio
drivers\audio\mic_driver.go
drivers\camera
drivers\camera\camera_driver.go
drivers\can
drivers\can\can_bus_driver.go
execute_boot
execute_boot\builder
execute_boot\builder\boot_context_builder.go
execute_boot\orchestrator
execute_boot\orchestrator\boot_orchestrator.go
execute_boot\phases
execute_boot\phases\attestation_phase.go
execute_boot\phases\boot_resolution_phase.go
execute_boot\phases\capability_phase.go
execute_boot\phases\discovery_phase.go
execute_boot\phases\identity_phase.go
execute_boot\phases\interface_phase.go
execute_boot\phases\module_resolution_phase.go
execute_boot\platform
execute_boot\platform\computer\receivers
execute_boot\platform\computer\receivers\can_receiver.go
execute_boot\platform\computer\receivers\mqtt_receiver.go
execute_boot\platform\computer\receivers\receiver.go
execute_boot\platform\industrial\receivers
execute_boot\platform\industrial\receivers\can_receiver.go
execute_boot\platform\industrial\receivers\mqtt_receiver.go
execute_boot\platform\industrial\receivers\receiver.go
execute_boot\platform\vehicle\receivers
execute_boot\platform\vehicle\receivers\can_receiver.go
execute_boot\platform\vehicle\receivers\mqtt_receiver.go
execute_boot\platform\vehicle\receivers\receiver.go
execute_boot\platform\identity_resolver.go
execute_boot\platform\resolve.go
execute_boot\platform\scoring.go
execute_boot\probe
execute_boot\probe\active_discovery.go
execute_boot\probe\hardware_fingerprint.go
execute_boot\probe\hardware_profile.go
execute_boot\probe\identity_probe.go
execute_boot\probe\passive_discovery.go
execute_boot\probe\type_struct.go
execute_boot\resolver
execute_boot\resolver\boot_policy_resolver.go
execute_boot\resolver\execution_context_resolver.go
execute_boot\types
execute_boot\types\boot_modes.go
execute_boot\detect_cap.go
execute_boot\runtime_context.md
execute_boot\runtime_types.go
execute_boot\session.go
internal
internal\apppath
internal\apppath\paths.go
internal\keys
internal\keys\env_keys.go
internal\logging
internal\logging\bloat_guard.go
internal\logging\logger.go
internal\logging\reflective_logger.go
internal\logging\structured_logger.go
internal\math_convert
internal\math_convert\convert_int_byte.go
internal\math_convert\q16.go
internal\network
internal\network\network_discovery_service.go
internal\policy
internal\policy\decision.go
internal\policy\model.go
internal\policy\registry.go
internal\policy\resolver.go
internal\schema
internal\schema\bootstrap
internal\schema\bootstrap\boot_marker.go
internal\schema\bootstrap\boot_mode.go
internal\schema\bootstrap\context.go
internal\schema\environment
internal\schema\environment\attestation.go
internal\schema\environment\capabilities.go
internal\schema\environment\device_identity.go
internal\schema\environment\discovery_profile.go
internal\schema\environment\env_config.go
internal\schema\environment\hardware_profile.go
internal\schema\user
internal\schema\user\settings.go
internal\schema\verification
internal\schema\verification\permissions.go
internal\schema\verification\security_config.go
internal\schema\all_internal_schema.txt
modules
modules\audio
modules\audio\audio_feature_service.go
modules\audio\audio_module_entry.go
modules\audio\capture_pcm.go
modules\audio\wav_writer.go
modules\auth
modules\auth\session.go
modules\contracts
modules\contracts\module_contract.go
modules\file
modules\file\file_module_entry.go
modules\file\http_ingestion.go
modules\implement_unknown
modules\implement_unknown\adapter.go
modules\implement_unknown\download.go
modules\implement_unknown\manifest.json
modules\industrial
modules\industrial\industrial_module_entry.go
modules\inference
modules\inference\inference_module_entry.go
modules\registry
modules\registry\module_dependency_resolver.go
modules\registry\module_registry.go
modules\adapter.go
modules\audit_module.go
modules\autonomous_kernel.go
modules\base_module.go
modules\cognition_module.go
modules\database_sink_module.go
modules\domain_module.go
modules\domain_types.go
modules\filter_module.go
modules\industrial_protocol_module.go
modules\inference_module.go
modules\ingestion_module.go
modules\productivity_module.go
modules\telemetry_module.go
mutual_interaction
mutual_interaction\interaction_mode_controller.go
runtime
runtime\adapter
runtime\adapter\cli_adapter.go
runtime\adapter\hmi_adapter.go
runtime\bus
runtime\bus\message_bus.go
runtime\engine
runtime\engine\runtime_builder.go
runtime\session
runtime\session\session_manager.go
runtime\supervisor
runtime\supervisor\runtime_supervisor.go
runtime\audio_vad.go
runtime\voice_engine.go
runtime\wake_word.go
Audience: external users, integrators, OEMs
Trust level: high, verification-sensitive
Stability: slow-changing, conservative


Safety Inversion: The verification and policy folders are moved inside core/. This ensures that Identity and Trust are part of the kernel's memory space, making it harder for a compromised plugin to bypass safety checks.

Internalized Determinism: By moving the scheduler and watchdog to internal/, we prevent external developers from accidentally modifying the timing constraints that keep the real-time loop stable.

Simulation Isolation: simulation/ is strictly a peer to the execution layers. This enforces the rule that Simulation is a Data Generator, and code within this folder can never accidentally control physical hardware.

Policy Distillation Path: The addition of cognition/distillation/ provides a dedicated home for the logic that converts a "Heavy Teacher" (GPU-bound) into a "Fast-Load Student" (Real-time bound).


If you were to deploy this today on a tractor:core/platform/probe would see a CAN-BUS and a ZED-Camera.core/platform/classify would trigger the vehicle.yaml config.bridge/busmap would correlate a pulse on a specific wire to the hydraulic arm moving.core/policy would incrementally increase the trust score for that arm.Once trust hits $99\%$, cognition/distillation would create a tiny, fast C++ snippet to control that arm perfectly.



Below is a concise, technical review of the operating circumstances for Multi-Platform AI and a clear list of pitfalls to avoid, grounded in real-world production failures seen in autonomous systems, robotics platforms, and large modular AI stacks. The tone is deliberately corrective and implementation-focused.


1. Operating Circumstances (Including First-Boot Reality)

Multi-Platform AI operates under non-ideal, adversarial, and heterogeneous conditions. Design decisions must assume the following are always true:

On first bootstrap, the system does not yet know what it is allowed to be. It must determine which platform class it is installed on before any domain logic, intelligence, or plugins are activated. This determination is based on execution context, hardware capabilities, and attestation—not user intent or installation method.

The system will run across radically different platforms (PCs, vehicles, embedded controllers) with inconsistent hardware guarantees, partial sensor availability, and varying real-time constraints. Hardware may be replaced, misconfigured, obstructed, degraded, or temporarily unavailable between boots. Platform classification must therefore be revalidated, not assumed.

Users are not homogeneous. Identity, trust tier, intent, and competence vary continuously. Many sessions are transient (guests, passengers, technicians), and some users must be explicitly prevented from accessing system-critical capabilities despite being physically present. User identity is subordinate to platform authority, not the other way around.

Environmental perception is inherently uncertain. Vision-based depth is noisy, monocular inference is scale-ambiguous, GNSS is unreliable or absent in many domains, and all sensors can lie temporarily or fail silently. Platform classification must not depend on perception correctness.

Workloads are bursty and asymmetric. Most of the time the system should be doing very little, but when anomalies occur, it must react deterministically and immediately. Heavy AI workloads must never sit on the critical path of platform detection, safety enforcement, or control logic.

Finally, the system will evolve. New plugins, sensors, policies, and regulatory constraints will be introduced long after initial deployment. First-bootstrap and re-bootstrap platform identification must remain backward-compatible, conservative, and biased toward restriction rather than assumption.

2. What to Avoid (Critical Architectural and Conceptual Traps)

A. Avoid Deciding “What AI This Is” Before Knowing the Platform

Do not allow the system to assume it is a “desktop AI,” “vehicle AI,” or “industrial AI” based on the launcher, build target, or installation path. Platform identity must be discovered and attested first; AIOS mode is a consequence of that discovery, not an input.

B. Avoid Treating “AI” as a Monolith

Do not allow perception, reasoning, control, and policy enforcement to collapse into a single execution context. This leads to non-deterministic behavior, unbounded latency, and catastrophic failure modes. Platform detection, safety, and trust logic must be minimal, deterministic, and hostile to intelligence modules.

C. Avoid Letting High-Fidelity Models Drive Platform or Safety Decisions

Gaussian splats, neural depth maps, semantic segmentation, and generative reconstruction are non-authoritative. They must never be used to infer platform class, actuator availability, or safety envelopes. Platform detection relies on hardware facts and attestation, not inference.

D. Avoid Running Discovery or Learning During Platform Identification

Closed-loop bus discovery, adaptive learning, or exploratory probing must not run during first bootstrap or platform classification. Hardware probing must be read-only and bounded. Active signaling belongs only in explicitly authorized modes.

E. Avoid Global State and Implicit Coupling

Shared mutable state across plugins, perception stacks, or agents will eventually cause race conditions, policy bypasses, or cascading failures. Platform identity must be a read-only fact published by Layer I, never negotiated downstream.

F. Avoid Assuming Sensor Truth or Persistence

Do not assume that a detected sensor, bus, or device will remain present. Platform classification must tolerate partial availability and downgrade safely if required.

G. Avoid Over-Personalization Before Platform Lock-In

Natural language personalization, behavior adaptation, and agent autonomy must not activate before platform and trust tier are fixed. Personalization layers sit on top of a constrained, platform-approved intent schema.

H. Avoid Loading “Just in Case” Code

Before platform identification completes, only Layer I and minimal HAL probing code may load. Preloading plugins or models before platform lock-in increases attack surface and risks irreversible misclassification.

I. Avoid Treating Simulation as a Valid Platform

Simulation is a platform class, not a shortcut. Simulation environments must be explicitly detected and sandboxed. No simulated capability should ever be implicitly trusted as real hardware.

3. Structural Safeguards to Enforce

Platform identification must occur on first bootstrap inside core/bootstrap + core/attestation, before any plugin loader, runtime scheduler, or cognition module is allowed to start.

The platform result must be:

Deterministic

Attested or provable

Versioned

Downgradable but not upgradable without reprovisioning

Every capability must answer three questions before execution:
Who is requesting this?
On which verified platform?
Under what trust and safety conditions?

All expensive computation must be interruptible, preemptible, and killable without destabilizing platform identity or safety enforcement.

All mappings—from bus ports to actuators, from words to actions, from pixels to objects—must be probabilistic, versioned, reversible, and subordinate to platform constraints.

4. Final Reality Check (Revised)

The most common failure in systems like Multi-Platform AI is not model accuracy or feature completeness. It is allowing the system to decide what it is before proving where it is running.

If you maintain strict separation between:
platform identity and AI behavior,
safety and intelligence,
perception and control,
simulation and reality,
identity and capability,