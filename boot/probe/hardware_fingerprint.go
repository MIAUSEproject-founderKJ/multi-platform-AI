//boot\probe\hardware_fingerprint.go

package probe

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
)

type HardwareFingerprint struct {
	TPM     string
	CPU     string
	DMI     string
	PCI     []string
	MAC     []string
	Storage []string
}

func collectHardwareFingerprint() HardwareFingerprint {

	return HardwareFingerprint{
		TPM:     readTPM(),
		CPU:     readCPUModel(),
		DMI:     readDMIUUID(),
		PCI:     readPCITopology(),
		MAC:     readMACs(),
		Storage: readDiskSerials(),
	}
}

func buildRobustMachineID() string {

	fp := collectHardwareFingerprint()

	// ---- Tier 1 & 2 (CORE IDENTITY) ----
	coreParts := []string{
		fp.TPM,
		fp.DMI,
		fp.CPU,
	}

	core := strings.Join(coreParts, "|")

	// ---- Tier 3 & 4 (ENTROPY ONLY) ----
	entropyParts := []string{
		strings.Join(fp.PCI, ","),
		strings.Join(fp.MAC, ","),
		strings.Join(fp.Storage, ","),
	}

	entropy := strings.Join(entropyParts, "|")

	if fp.TPM != "" {
		// Ignore unstable signals entirely
		entropy = ""
	}
	// ---- Combined canonical structure ----
	final := "core:" + core + "||entropy:" + entropy

	hash := sha256.Sum256([]byte(final))
	return hex.EncodeToString(hash[:])
}
func readCPUModel() string {

	if runtime.GOOS != "linux" {
		return ""
	}

	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return ""
	}

	for _, line := range strings.Split(string(data), "\n") {
		if strings.Contains(line, "model name") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
	}

	return ""
}

// This is extremely useful for robotics and vehicles because it fingerprints the motherboard bus layout.
func readPCITopology() []string {

	out, err := exec.Command("lspci").Output()
	if err != nil {
		return nil
	}

	var devices []string

	for _, l := range strings.Split(string(out), "\n") {
		l = normalize(l)
		if l != "" {
			devices = append(devices, l)
		}
	}

	sort.Strings(devices)
	return devices
}

// Disk serial numbers.
func readDiskSerials() []string {

	out, err := exec.Command("lsblk", "-o", "SERIAL").Output()
	if err != nil {
		return nil
	}

	var serials []string

	for _, l := range strings.Split(string(out), "\n") {
		l = normalize(l)
		if l == "" || l == "serial" {
			continue
		}

		serials = append(serials, l)
	}

	sort.Strings(serials)
	return serials
}
func readMACs() []string {

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	var macs []string

	for _, i := range ifaces {

		// Skip loopback & down interfaces
		if i.Flags&net.FlagLoopback != 0 || i.Flags&net.FlagUp == 0 {
			continue
		}

		mac := strings.TrimSpace(i.HardwareAddr.String())
		if mac == "" {
			continue
		}

		macs = append(macs, normalize(mac))
	}

	sort.Strings(macs)
	return macs
}

func readTPM() string {

	if runtime.GOOS != "linux" {
		return ""
	}

	// Preferred: tpm2-tools
	out, err := exec.Command("tpm2_getcap", "properties-fixed").Output()
	if err == nil && len(out) > 0 {
		return normalize(string(out))
	}

	// Fallback: sysfs
	paths := []string{
		"/sys/class/tpm/tpm0/device/unique_id",
		"/sys/class/tpm/tpm0/device/description",
	}

	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err == nil && len(data) > 0 {
			return normalize(string(data))
		}
	}

	return ""
}
