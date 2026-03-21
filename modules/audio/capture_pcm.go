// modules/audio/capture_pcm.go
package audio

import (
	"context"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
	"github.com/gordonklaus/portaudio"
)

func StartMicrophoneStream(ctx context.Context, bus *schema.MessageBus) error {

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

			raw := schema.int16ToBytes(buffer)

			_ = bus.Publish(ctx, "audio.raw", raw)
		}
	}
}
