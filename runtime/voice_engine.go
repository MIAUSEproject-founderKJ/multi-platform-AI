//runtime/voice_engine.go
package runtimectx

import (
	"fmt"
	"sync"
	"time"
)

type VoiceEngine struct {
	audio *AudioVAD

	inputChan     chan string
	outputChan    chan string
	interruptChan chan struct{}

	speaking bool
	mu       sync.Mutex
}

func NewVoiceEngine() *VoiceEngine {
	return &VoiceEngine{
		audio:         NewAudioVAD(),
		inputChan:     make(chan string),
		outputChan:    make(chan string),
		interruptChan: make(chan struct{}, 1),
	}
}

//============Start Full Duplex Loop
func (ve *VoiceEngine) Start() {
	go ve.listenLoop()
	go ve.processLoop()
	go ve.speakLoop()
}

//====Simulated Microphone (replace later with real STT)
func (ve *VoiceEngine) listenLoop() {
	err := ve.audio.Start()
	if err != nil {
		fmt.Println("[VOICE] Audio init failed:", err)
		return
	}

	buffer := make([]byte, 0, 16000)

	for frame := range ve.audio.outChan {

		// interrupt if speaking
		ve.mu.Lock()
		if ve.speaking {
			select {
			case ve.interruptChan <- struct{}{}:
			default:
			}
		}
		ve.mu.Unlock()

		buffer = append(buffer, frame...)

		// simple segmentation heuristic
		if len(buffer) > sampleRate*2 { // ~1 second speech
			text := fakeSTT(buffer)
			ve.inputChan <- text
			buffer = buffer[:0]
		}
	}
}

func fakeSTT(audio []byte) string {
	return "[speech detected]"
}

//========Processing Layer
func (ve *VoiceEngine) processLoop() {
	for input := range ve.inputChan {

		// simulate AI processing latency
		time.Sleep(200 * time.Millisecond)

		response := "You said: " + input

		ve.outputChan <- response
	}
}
//=======Interruptible TTS (Core Feature)

func (ve *VoiceEngine) speakLoop() {
	for msg := range ve.outputChan {

		ve.mu.Lock()
		ve.speaking = true
		ve.mu.Unlock()

		fmt.Println("[VOICE] Speaking:", msg)

		// simulate streaming speech
		for i := 0; i < len(msg); i++ {
			select {
			case <-ve.interruptChan:
				fmt.Println("\n[VOICE] Interrupted")
				goto END
			default:
				fmt.Printf("%c", msg[i])
				time.Sleep(30 * time.Millisecond)
			}
		}

		fmt.Println()

	END:
		ve.mu.Lock()
		ve.speaking = false
		ve.mu.Unlock()
	}
}

