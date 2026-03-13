// modules/audio/wav_writer.go
package audio

import (
	"encoding/binary"
	"os"
	"time"
)

type WAVWriter struct {
	file       *os.File
	sampleRate int
	totalBytes int
}

func NewWAVWriter(path string) *WAVWriter {
	f, _ := os.Create(path + time.Now().Format("20060102150405") + ".wav")
	w := &WAVWriter{
		file:       f,
		sampleRate: 16000,
	}
	w.writeHeader()
	return w
}

func (w *WAVWriter) writeHeader() {
	header := make([]byte, 44)
	// RIFF header placeholder (will patch size later)
	copy(header[0:], []byte("RIFF"))
	copy(header[8:], []byte("WAVE"))
	copy(header[12:], []byte("fmt "))
	binary.LittleEndian.PutUint32(header[16:], 16)
	binary.LittleEndian.PutUint16(header[20:], 1) // PCM
	binary.LittleEndian.PutUint16(header[22:], 1) // mono
	binary.LittleEndian.PutUint32(header[24:], uint32(w.sampleRate))
	binary.LittleEndian.PutUint32(header[28:], uint32(w.sampleRate*2))
	binary.LittleEndian.PutUint16(header[32:], 2)
	binary.LittleEndian.PutUint16(header[34:], 16)
	copy(header[36:], []byte("data"))
	w.file.Write(header)
}

func (w *WAVWriter) AppendPCM(pcm []byte) error {
	w.totalBytes += len(pcm)
	_, err := w.file.Write(pcm)
	return err
}
