//runtime/audio_vad.go

package runtime

import (
	"encoding/binary"
	"fmt"
	"log"
	"strings"
	"github.com/gordonklaus/portaudio"
	"github.com/maxhawkins/go-webrtcvad"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

const (
	sampleRate = 16000
	frameSize  = 160 // 10ms @ 16kHz
)

type AudioVAD struct {
	stream      *portaudio.Stream
	vad         *webrtcvad.VAD
	frameBuf    []int16
	outChan     chan []byte
	WakeDetector *WakeWordDetector
	wakeWordOn   bool
	wakeBuffer  []int16
}


func detectWakeWord(text string) bool {
	// lowercase match for demo
	text = strings.ToLower(text)
	return strings.Contains(text, "hey system")
}

func NewAudioVAD() *AudioVAD {
	vad, err := webrtcvad.New()
if err != nil {
	log.Fatal(err)
}
	vad.SetMode(3) // aggressive

	return &AudioVAD{
		vad:      vad,
		frameBuf: make([]int16, frameSize),
		outChan:  make(chan []byte, 100),
	}
}

//===========Initialize Microphone Stream
func (a *AudioVAD) Start() error {
	if err := portaudio.Initialize(); err != nil {
		return err
	}

	stream, err := portaudio.OpenDefaultStream(
		1, // input channels
		0, // output channels
		sampleRate,
		len(a.frameBuf),
		&a.frameBuf,
	)

	if err != nil {
		return err
	}

	a.stream = stream

	if err := a.stream.Start(); err != nil {
		return err
	}
detector, err := NewWakeWordDetector("<ACCESS_KEY>", "hey_system.ppn")
if err != nil {
    log.Fatal(err)
}
	vad.WakeDetector = detector
	defer detector.Close()

	go a.loop()

	fmt.Println("[AUDIO] VAD microphone started")
	return nil
}

//===========Core VAD Loop (Speech Detection)
func (a *AudioVAD) loop() {
    wakeDetector := a.WakeDetector // new field in AudioVAD

    for {
        err := a.stream.Read()
        if err != nil {
            log.Println("[AUDIO] read error:", err)
            continue
        }

        if !a.wakeWordOn {
            // feed frame to Porcupine
            if wakeDetector.Process(a.frameBuf) {
                a.wakeWordOn = true
                fmt.Println("[WAKEWORD] Detected 'Hey system', VAD activated")
            }
            continue
        }

        // VAD processing
        buf := schema.Int16ToBytes(a.frameBuf)
		active, err := a.vad.Process(sampleRate, buf)
        if err != nil {
            continue
        }
        if active {
            select {
            case a.outChan <- audioInt16ToBytes(a.frameBuf):
            default:
            }
        }
    }
}

//=========Stop / Cleanup
func (a *AudioVAD) Stop() {
	if a.stream != nil {
		a.stream.Stop()
		a.stream.Close()
	}
	portaudio.Terminate()
}

func audioInt16ToBytes(frames []int16) []byte {
	buf := make([]byte, len(frames)*2)
	for i, v := range frames {
		binary.LittleEndian.PutUint16(buf[i*2:], uint16(v))
	}
	return buf
}

func (a *AudioVAD) ResetWakeWord() {
	a.wakeWordOn = false
	a.wakeBuffer = nil
	fmt.Println("[WAKEWORD] Reset")
}