// modules/audio/capture_pcm.go
package audio

import (
	"context"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/encoding"
	runtime_bus "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/bus"
	"github.com/gordonklaus/portaudio"
)

func StartMicrophoneStream(ctx context.Context, bus *runtime_bus.MessageBus) error {

	if err := portaudio.Initialize(); err != nil {
		return err
	}
	defer portaudio.Terminate()

	buffer := make([]int16, 320)

	stream, err := portaudio.OpenDefaultStream(1, 0, 16000, len(buffer), &buffer)
	if err != nil {
		return err
	}
	defer stream.Close()

	if err := stream.Start(); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := stream.Read(); err != nil {
				return err
			}
			raw := encoding.Int16ToBytes(buffer)
			bus.Publish(runtime_bus.Message{
				Topic: "audio.raw",
				Data:  raw,
			})
		}
	}
}
