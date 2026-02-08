//core/vision/engine.go
package vision
import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/memory"
)

type Object struct {
	Label      string
	Confidence float64
	BoundingBox [4]int // x, y, w, h
}

func (ve *VisionEngine) ProcessFrame(frame []byte) []Object {
	// 1. Run inference (e.g., via TensorRT or TFLite on the Edge TPU)
	detected := ve.Inference(frame)

	// 2. Feed to Cognitive Vault
	for _, obj := range detected {
		ve.Kernel.Memory.Store(obj.Label, memory.CognitiveEntry{
			Type:   "visual_memory",
			Weight: obj.Confidence,
			Data:   obj,
		})
	}
	return detected
}