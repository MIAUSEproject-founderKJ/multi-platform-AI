// runtime/wake_word.go
package runtime

import (
	"log"

	porcupine "github.com/Picovoice/porcupine/binding/go"
)

type WakeWordDetector struct {
	engine *porcupine.Porcupine
}

func NewWakeWordDetector(accessKey, keywordPath string) (*WakeWordDetector, error) {
	engine, err := porcupine.New(accessKey,[]string{keywordPath},[]float32{0.5},)

	if err != nil {
		return nil, err
	}
	return &WakeWordDetector{engine: engine}, nil
}

func (w *WakeWordDetector) Process(frame []int16) bool {
	idx, err := w.engine.Process(frame)
	if err != nil {
		log.Println("[WakeWord] process error:", err)
		return false
	}
	return idx == 0 // index of keyword matched
}

func (w *WakeWordDetector) Close() {
	w.engine.Delete()
}
