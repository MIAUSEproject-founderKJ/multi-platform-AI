// boot/probe/passive_discovery.go
package probe

import (
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/boot/platform"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

func PassiveDiscovery() (*schema.EnvConfig, error) {

	logging.Info("[PROBE] Phase 1: Passive Identity Extraction")

	machineID := buildRobustMachineID()
	hostname, _ := os.Hostname()

	hardware := collectHardwareProfile()

	env := &schema.EnvConfig{
		SchemaVersion: schema.CurrentVersion,
		GeneratedAt:   time.Now(),
		Identity: schema.MachineIdentity{
			MachineID: machineID,
			Hostname:  hostname,
			OS:        runtime.GOOS,
			Arch:      runtime.GOARCH,
		},
		Hardware: hardware,
	}

	platform.RunPlatformInference(env)

	logging.Info(
		"[PROBE] Identity confirmed: %s (%s)",
		env.Identity.MachineID,
		env.Platform.Final,
	)

	return env, nil
}

func resolveHardwareRoot() string {

	if runtime.GOOS == "windows" {
		if v := readWindowsUUID(); v != "" {
			return v
		}
	}

	if v := readTPMIdentity(); v != "" {
		return v
	}

	if v := readDMIUUID(); v != "" {
		return v
	}

	if v := readSystemSerial(); v != "" {
		return v
	}

	if v := readMACFingerprint(); v != "" {
		return v
	}

	return "soft-identity"
}

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

func readDMIUUID() string {

	if runtime.GOOS != "linux" {
		return ""
	}

	data, err := os.ReadFile("/sys/class/dmi/id/product_uuid")
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

	if len(macs) == 0 {
		return ""
	}

	return strings.Join(macs, "-")
}

func readWindowsUUID() string {

	out, err := exec.Command(
		"wmic",
		"csproduct",
		"get",
		"uuid",
	).Output()

	if err != nil {
		return ""
	}

	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		return ""
	}

	return strings.TrimSpace(lines[1])
}
