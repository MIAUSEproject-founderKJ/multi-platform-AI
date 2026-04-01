//boot/probe/hardware_fingerprint.go

package probe

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

// HardwareFingerprint represents low-level device identifiers.
type HardwareFingerprint struct {
	TPM     string
	CPU     string
	GPU     string
	DMI     string
	PCI     []string
	MAC     []string
	Storage []string
	Buses   map[string]bool
}

type fingerprintBuilder struct {
	mu sync.Mutex
	fp HardwareFingerprint
}

func (b *fingerprintBuilder) setString(target *string, val string) {
	b.mu.Lock()
	*target = val
	b.mu.Unlock()
}

func (b *fingerprintBuilder) setSlice(target *[]string, val []string) {
	b.mu.Lock()
	*target = val
	b.mu.Unlock()
}

// CollectHardwareFingerprint gathers basic hardware info with a timeout.
func CollectHardwareFingerprint(ctx context.Context) (HardwareFingerprint, []string) {
	var probeErrors []string
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	builder := &fingerprintBuilder{
		fp: HardwareFingerprint{
			Buses: make(map[string]bool),
		},
	}

	var wg sync.WaitGroup
	run := func(fn func(context.Context)) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fn(ctx)
		}()
	}

	// Core identifiers
	run(func(c context.Context) { builder.setString(&builder.fp.TPM, readTPM()) })
	run(func(c context.Context) { builder.setString(&builder.fp.CPU, readCPUModel(c)) })
	run(func(c context.Context) { builder.setString(&builder.fp.DMI, readDMIUUID(c)) })

	// PCI devices
	run(func(c context.Context) {
		pci := readPCITopology(c)
		builder.setSlice(&builder.fp.PCI, pci)
	})

	// MAC addresses
	run(func(c context.Context) { builder.setSlice(&builder.fp.MAC, readMACs()) })

	// Storage
	run(func(c context.Context) { builder.setSlice(&builder.fp.Storage, readDiskSerials(c)) })

	wg.Wait()

	// Detect buses
	detectBuses(&builder.fp)

	return builder.fp, probeErrors
}

// Utilities like TPM, CPU, DMI, PCI, MAC, Storage (same as your previous code)...

//
// ------------------------------------------------------------
// Command helper
// ------------------------------------------------------------
//

// MaxCommandOutput defines a safety limit (1MB) to prevent a bugged or
// malicious tool from exhausting system memory.
const MaxCommandOutput = 1024 * 1024

type limitedBuffer struct {
	buf       *bytes.Buffer
	limit     int
	written   int
	truncated bool
}

func (l *limitedBuffer) Write(p []byte) (int, error) {
	remaining := l.limit - l.written

	if remaining <= 0 {
		l.truncated = true
		return len(p), nil
	}

	if len(p) > remaining {
		l.buf.Write(p[:remaining])
		l.written += remaining
		l.truncated = true
		return len(p), nil
	}

	l.buf.Write(p)
	l.written += len(p)
	return len(p), nil
}

// runCommand executes a system binary with context awareness,
// resource limits, and detailed error reporting.
func runCommand(ctx context.Context, name string, args ...string) (string, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return "", fmt.Errorf("command %s not found: %w", name, err)
	}

	cmd := exec.CommandContext(ctx, path, args...)
	cmd.Env = os.Environ()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// Track truncation
	stdoutLimited := &limitedBuffer{buf: &stdout, limit: MaxCommandOutput}
	stderrLimited := &limitedBuffer{buf: &stderr, limit: MaxCommandOutput}

	cmd.Stdout = stdoutLimited
	cmd.Stderr = stderrLimited

	err = cmd.Run()

	// Context handling (precise)
	if ctx.Err() != nil {
		return "", fmt.Errorf("command %s cancelled: %w", name, ctx.Err())
	}

	if err != nil {
		return "", fmt.Errorf("command %s failed: %w (stderr: %s)",
			name, err, strings.TrimSpace(stderr.String()))
	}

	out := strings.TrimSpace(stdout.String())

	// Optional: annotate truncation
	if stdoutLimited.truncated {
		out += "\n[truncated]"
	}

	return out, nil
}

//
// ------------------------------------------------------------
// Probe implementations
// ------------------------------------------------------------
//

// TPM

func readTPM() string {
	if runtime.GOOS != "linux" {
		return ""
	}

	data, err := os.ReadFile("/sys/class/tpm/tpm0/device/description")
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(data))
}

// CPU

func readCPUModel(ctx context.Context) string {

	switch runtime.GOOS {
	case "linux":
		data, err := os.ReadFile("/proc/cpuinfo")
		if err != nil {
			return ""
		}
		for _, line := range strings.Split(string(data), "\n") {
			if strings.Contains(line, "model name") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					return normalize(parts[1])
				}
			}
		}
	case "windows":
		out, err := runCommand(ctx, "wmic", "cpu", "get", "name")
		if err == nil {
			lines := strings.Split(out, "\n")
			if len(lines) > 1 {
				return normalize(lines[1])
			}
		}
	}

	return ""
}

// DMI / UUID

func readDMIUUID(ctx context.Context) string {

	switch runtime.GOOS {

	case "linux":
		data, err := os.ReadFile("/sys/class/dmi/id/product_uuid")
		if err == nil {
			return normalize(string(data))
		}

	case "windows":
		out, err := runCommand(ctx, "wmic", "csproduct", "get", "uuid")
		if err == nil {
			lines := strings.Split(out, "\n")
			if len(lines) > 1 {
				return normalize(lines[1])
			}
		}
	}

	return ""
}

// PCI

func readPCITopology(ctx context.Context) []string {

	if runtime.GOOS != "linux" {
		return nil
	}

	out, err := runCommand(ctx, "lspci", "-mm")
	if err != nil {
		return nil
	}

	var devices []string

	for _, l := range strings.Split(out, "\n") {
		l = normalize(l)
		if l != "" {
			devices = append(devices, l)
		}
	}

	sort.Strings(devices)
	return devices
}

// MAC

func readMACs() []string {

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	var macs []string

	for _, i := range ifaces {

		if i.Flags&net.FlagLoopback != 0 {
			continue
		}
		if i.Flags&net.FlagUp == 0 {
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

// Storage (basic implementation)

func readDiskSerials(ctx context.Context) []string {

	if runtime.GOOS != "linux" {
		return nil
	}

	out, err := runCommand(ctx, "lsblk", "-o", "SERIAL")
	if err != nil {
		return nil
	}

	var serials []string

	for _, line := range strings.Split(out, "\n") {
		line = normalize(line)
		if line != "" && line != "SERIAL" {
			serials = append(serials, line)
		}
	}

	sort.Strings(serials)
	return serials
}

//
// ------------------------------------------------------------
// Utility
// ------------------------------------------------------------
//

func normalize(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}
