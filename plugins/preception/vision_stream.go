//MIAUSEproject-founderKJ/multi-platform-AI/plugins/preception/vision_stream.go

package perception

import (
	"fmt"
	"image"
	"image/color"
	"sync"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/api/hmi"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/monitor"
)

type VisionStream struct {
	mu            sync.RWMutex
	currentVitals monitor.SystemVitals
	currentBoot   hmi.ProgressUpdate

	// HUD Colors (Terminator Red / Cyber Blue)
	PrimaryColor color.RGBA
}

func NewVisionStream() *VisionStream {
	return &VisionStream{
		PrimaryColor: color.RGBA{R: 255, G: 0, B: 0, A: 180}, // Classic "T-800" Red
	}
}

// UpdateVitals and UpdateProgress are called by the HMI Controller
func (vs *VisionStream) UpdateVitals(v monitor.SystemVitals) {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	vs.currentVitals = v
}

func (vs *VisionStream) UpdateProgress(p hmi.ProgressUpdate) {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	vs.currentBoot = p
}

// ProcessFrame applies the HUD overlay to a single raw camera frame
func (vs *VisionStream) ProcessFrame(rawFrame image.Image) image.Image {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	// 1. Create a canvas from the raw frame
	// (Using a theoretical DrawHUD helper)
	canvas := drawHUDBackground(rawFrame)

	// 2. Render Top-Left: System Vitals
	drawText(canvas, 20, 40, fmt.Sprintf("CPU_LOAD: %.1f%%", vs.currentVitals.CPULoad), vs.PrimaryColor)
	drawText(canvas, 20, 60, fmt.Sprintf("VRAM_USE: %dMB", vs.currentVitals.VRAMUsage/1024/1024), vs.PrimaryColor)
	drawText(canvas, 20, 80, fmt.Sprintf("SYS_TEMP: %.1fC", vs.currentVitals.Temperature), vs.PrimaryColor)

	// 3. Render Bottom-Center: Boot Progress Bar
	drawProgressBar(canvas, vs.currentBoot.Percentage, vs.currentBoot.Message, vs.PrimaryColor)

	// 4. Render Target Reticle (Center)
	drawReticle(canvas, vs.PrimaryColor)

	return canvas
}
