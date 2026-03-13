//boot\probe\hardware_fingerprint.go

package probe

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"os"
	"os/exec"
	"runtime"
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

	data := strings.Join([]string{
		fp.TPM,
		fp.CPU,
		fp.DMI,
		strings.Join(fp.PCI, ","),
		strings.Join(fp.MAC, ","),
		strings.Join(fp.Storage, ","),
	}, "|")

	hash := sha256.Sum256([]byte(data))
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

	lines := strings.Split(string(out), "\n")

	var devices []string

	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			devices = append(devices, l)
		}
	}

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
		l = strings.TrimSpace(l)
		if l != "" && l != "SERIAL" {
			serials = append(serials, l)
		}
	}

	return serials
}

func readMACs() []string {

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	var macs []string

	for _, i := range ifaces {
		if len(i.HardwareAddr) > 0 {
			macs = append(macs, i.HardwareAddr.String())
		}
	}

	return macs
}
