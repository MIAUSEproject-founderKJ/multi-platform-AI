//X/probe.md
To ensure the system remains stable during the critical "First-Boot" phase, probe.go must operate as a Passive Observer before it ever becomes an Active Interrogator.

If the system blindly sends signals to unknown ports, it risks triggering a "Deadly Embrace" (where two hardware components lock up) or causing an electrical fault.

1. The Logic of core/platform/probe/probe.go
In the Multi-Platform AI architecture, this file implements a "Safety-First" discovery sequence. It follows a three-step escalation path to identify the hardware environment without risking a kernel panic.

Mechanism A: Passive Fingerprinting (The "Listen" Phase)
Before touching any bus, the system reads immutable software descriptors provided by the OS kernel.

DMI/SMBIOS: Reads the system UUID, BIOS version, and board manufacturer.

Device Tree / Sysfs: On Linux/Vehicle systems, it traverses /sys/class to see what drivers are already initialized.

PCI/USB Enumeration: It checks Vendor IDs (VID) and Product IDs (PID). If it sees a "Comet Lake" CPU and "NVIDIA" GPU, it instantly classifies itself as Desktop.

Mechanism B: Bounded Interrogation (The "Ask" Phase)
If passive checks are ambiguous, it sends Standardized Non-State-Changing Requests.

CPUID: A safe instruction that returns processor capabilities without affecting the CPU state.

OBD-II Mode 01: On a vehicle bus (CAN), it requests "Supported PIDs" (01 00). This is a read-only request that all vehicle ECUs are programmed to answer without moving any physical parts.

Mechanism C: Panic-Recovery Wrapping (The "Shield" Phase)
Every probe is wrapped in a Deferred Recovery Block. If a probe hits a "Memory-Mapped I/O" (MMIO) region that doesn't exist (which usually causes a crash), the Go runtime catches the fault, logs the failure, and moves to the next check.