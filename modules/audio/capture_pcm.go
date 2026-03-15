// modules/audio/capture_pcm.go
package audio

import (
	"context"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/agent"
	"github.com/gordonklaus/portaudio"
)

func StartMicrophoneStream(agent *agent.AgentRuntime, opt optimization.Optimizer) error {

	portaudio.Initialize()
	defer portaudio.Terminate()

	buffer := make([]int16, 320)

	stream, err := portaudio.OpenDefaultStream(1, 0, 16000, len(buffer), &buffer)
	if err != nil {
		return err
	}

	stream.Start()

	for {
		if err := stream.Read(); err != nil {
			return err
		}

		raw := int16ToBytes(buffer)
		agent.Process(context.Background(), opt, raw)
	}
}
