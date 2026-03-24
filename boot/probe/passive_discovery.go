// boot/probe/passive_discovery.go
package probe

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/mathutil"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

// -----------------------------
// Public: Passive Discovery
// -----------------------------
func PassiveDiscovery() (*schema.EnvConfig, error) {
	log.Println("[PROBE] Phase 1: Passive Identity Extraction")

	// Collect low-level hardware fingerprint
	fp := collectHardwareFingerprint()

	// Convert fingerprint to a hardware profile
	hardware := convertFingerprintToProfile(fp)

	// Gather system info
	hostname, _ := os.Hostname()
	osName := runtime.GOOS
	arch := runtime.GOARCH

	// Build EnvConfig
	env := &schema.EnvConfig{
		SchemaVersion: schema.CurrentVersion,
		GeneratedAt:   time.Now(),
		Identity: schema.MachineIdentity{
			MachineID:    buildRobustMachineID(fp),
			Hostname:     hostname,
			OS:           osName,
			Arch:         arch,
			Hardware:     hardware,
			PlatformType: schema.PlatformComputer, // default
		},

		Hardware: hardware,
	}

	// Run platform scoring & inference
	runPlatformInference(env, fp)

	log.Println("[PROBE] Identity confirmed:", env.Identity.MachineID, "Platform:", env.Platform.Final)
	return env, nil
}

// -----------------------------
// Convert HardwareFingerprint to HardwareProfile
// -----------------------------
func convertFingerprintToProfile(fp HardwareFingerprint) schema.HardwareProfile {
	buses := []schema.BusCapability{}

	if len(fp.PCI) > 0 {
		buses = append(buses, schema.BusCapability{
			ID:         "pci-root",
			Type:       "pci",
			Confidence: 0.6,
			Source:     "pci-scan",
		})
	}
	if len(fp.MAC) > 0 {
		buses = append(buses, schema.BusCapability{
			ID:         "ethernet",
			Type:       "network",
			Confidence: 0.6,
			Source:     "net-iface",
		})
	}

	return schema.HardwareProfile{
		Processors: []schema.Processor{
			{Type: "CPU", Count: runtime.NumCPU(), Version: 1.0},
		},
		Buses:      buses,
		HasBattery: detectBattery(),
	}
}

// -----------------------------
// Build robust machine ID
// -----------------------------
func buildRobustMachineID(fp HardwareFingerprint) string {
	core := strings.Join([]string{fp.TPM, fp.DMI, fp.CPU}, "|")
	entropy := strings.Join([]string{
		strings.Join(fp.PCI, ","),
		strings.Join(fp.MAC, ","),
		strings.Join(fp.Storage, ","),
	}, "|")

	// If TPM exists, ignore volatile signals
	if fp.TPM != "" {
		entropy = ""
	}

	final := "core:" + core + "||entropy:" + entropy
	hash := sha256.Sum256([]byte(final))
	return hex.EncodeToString(hash[:])
}

// -----------------------------
// Platform scoring / inference
// -----------------------------
func runPlatformInference(env *schema.EnvConfig, fp HardwareFingerprint) {
	scores := map[schema.PlatformClass]*schema.PlatformScore{}
	ensure := func(class schema.PlatformClass, max float64) *schema.PlatformScore {
		if scores[class] == nil {
			scores[class] = &schema.PlatformScore{Type: class, MaxScore: max}
		}
		return scores[class]
	}

	osName := strings.ToLower(env.Identity.OS)

	// Vehicle / robotic detection
	if hasBus(env.Hardware, "can") || osName == "qnx" || osName == "autosar" {
		s := ensure(schema.PlatformVehicle, 1.5)
		s.Score += 1.0
		s.Signals = append(s.Signals, "CAN bus / automotive RTOS detected")
	}
	if hasBus(env.Hardware, "i2c") && hasBus(env.Hardware, "spi") {
		s := ensure(schema.PlatformRobot, 1.2)
		s.Score += 0.4
		s.Signals = append(s.Signals, "I2C+SPI sensors detected")
	}

	// Desktop/Laptop scoring
	desktop := collectDesktopSignals(fp, env)
	scores[schema.PlatformComputer] = desktop

	// Convert to confidence and pick best
	var best schema.PlatformClass = schema.PlatformUnknown
	highConf := mathutil.Q16(0)

	candidates := []schema.PlatformScore{}

	if len(scores) == 0 {
		logging.Warn("[IDENTITY] No platform signals detected, defaulting to UNKNOWN")

		env.Platform.Final = schema.PlatformUnknown
		env.Platform.Locked = false
		return
	}

	for _, s := range scores {
		s.Confidence = mathutil.Q16(mathutil.FromFloat64(s.Score / s.MaxScore))
		candidates = append(candidates, *s)
		if s.Confidence > highConf {
			highConf = s.Confidence
			best = s.Type
		}
	}
	logging.Info("[DEBUG] Passive Discovery 1 RunPlatformInference: best=%s, highConf=%d%%", best, highConf.Percentage())

	env.Platform.Candidates = candidates
	env.Platform.Final = best
	env.Platform.Locked = true
	env.Platform.ResolvedAt = time.Now()

	logging.Info("[DEBUG] Passive Discovery 2 RunPlatformInference: best=%s, highConf=%d%%", best, highConf.Percentage())
	logging.Info("[IDENTITY]  Passive Discovery Resolution: %s (Conf: %d%%)", best, highConf.Percentage())
}

// Desktop/Laptop scoring helper
func collectDesktopSignals(fp HardwareFingerprint, env *schema.EnvConfig) *schema.PlatformScore {
	cpu := runtime.NumCPU()

	score := 0.2
	maxScore := 1.5

	score += 0.1 * float64(cpu)
	score += 0.05 * float64(len(fp.PCI))
	score += 0.05 * float64(len(fp.MAC))

	if env.Hardware.HasBattery {
		score += 0.3
	}

	if score > maxScore {
		score = maxScore
	}

	return &schema.PlatformScore{
		Type:     schema.PlatformComputer,
		Score:    score,
		MaxScore: maxScore,
		Signals: []string{
			fmt.Sprintf("CPU cores: %d", cpu),
			fmt.Sprintf("Battery: %v", env.Hardware.HasBattery),
			fmt.Sprintf("PCI devices: %d", len(fp.PCI)),
			fmt.Sprintf("MAC addresses: %d", len(fp.MAC)),
		},
	}
}

// -----------------------------
// Low-level helpers
// -----------------------------

func readTPMIdentity() string {
	if runtime.GOOS != "linux" {
		return ""
	}
	data, err := os.ReadFile("/sys/class/tpm/tpm0/device/description")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func readSystemSerial() string {
	if runtime.GOOS != "linux" {
		return ""
	}
	data, err := os.ReadFile("/sys/class/dmi/id/product_serial")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func readMACFingerprint() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	var macs []string
	for _, i := range ifaces {
		if i.HardwareAddr != nil {
			macs = append(macs, i.HardwareAddr.String())
		}
	}
	return strings.Join(macs, "-")
}

func readWindowsUUID() string {
	out, err := exec.Command("wmic", "csproduct", "get", "uuid").Output()
	if err != nil {
		return ""
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		return ""
	}
	return strings.TrimSpace(lines[1])
}
