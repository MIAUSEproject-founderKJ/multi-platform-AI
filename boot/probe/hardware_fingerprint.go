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

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

// ------------------------------------------------------------
// Types
// ------------------------------------------------------------

type HardwareFingerprint struct {
	TPM     string
	CPU     string
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

// ------------------------------------------------------------
// Public Entry
// ------------------------------------------------------------

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

	run(func(c context.Context) {
		builder.setString(&builder.fp.TPM, readTPM(c))
	})

	run(func(c context.Context) {
		builder.setString(&builder.fp.CPU, readCPUModel(c))
	})

	run(func(c context.Context) {
		builder.setString(&builder.fp.DMI, readDMIUUID(c))
	})

	var errMu sync.Mutex

	run(func(c context.Context) {

	run(func(c context.Context) {
	builder.mu.Lock()
	builder.fp.Buses = detectBuses(c)
	builder.mu.Unlock()
	})

	//func runProbe[T any](ctx context.Context, name string, fn func(context.Context) (T, error)) ProbeResult[T]
    res := runProbe(c, "pci_scan", func(ctx context.Context) ([]string, error) {
        return readPCITopology(ctx), nil
    })

    if res.Error != nil {
        errMu.Lock()
        probeErrors = append(probeErrors, fmt.Sprintf("%s: %v", res.Source, res.Error))
        errMu.Unlock()
    } else {
        builder.setSlice(&builder.fp.PCI, res.Value)
    }
	})

	run(func(c context.Context) {
		builder.setSlice(&builder.fp.MAC, readMACs(c))
	})

	run(func(c context.Context) {
		builder.setSlice(&builder.fp.Storage, readDiskSerials(c))
	})

	wg.Wait()
	return builder.fp, probeErrors
}

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

func readTPM(ctx context.Context) string {
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

func readMACs(ctx context.Context) []string {

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

//
// ------------------------------------------------------------
// Hardware Profile Conversion
// ------------------------------------------------------------
//

func convertFingerprintToProfile(fp HardwareFingerprint) schema.HardwareProfile {

	var buses []schema.BusCapability

	for bus := range fp.Buses {
		buses = append(buses, schema.BusCapability{
			ID:         bus + "-bus",
			Type:       bus,
			Confidence: 0.9,
			Source:     "fingerprint",
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
