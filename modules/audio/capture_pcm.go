// modules/audio/capture_pcm.go
package audio

import (
	"context"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime"
	"github.com/gordonklaus/portaudio"
)

func StartMicrophoneStream(ctx context.Context, bus *runtime.MessageBus) error {

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
			raw := schema.Int16ToBytes(buffer)
			bus.Publish(runtime.Message{
				Topic: "audio.raw",
				Data:  raw,
			})
		}
	}
}
