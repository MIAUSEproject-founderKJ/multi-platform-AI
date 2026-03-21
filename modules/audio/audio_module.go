// modules/audio/audio_module.go
package audio

import (
	"context"
	"encoding/json"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules"
)

func (m *AudioModule) Name() string {
	return "AudioModule"
}

func (m *AudioModule) DependsOn() []string {
	return []string{"StorageModule"}
}

type AudioModule struct {
	modules.BaseModule

	ctx       *schema.RuntimeContext
	writer    *WAVWriter
	extractor *FeatureExtractor
	repo      AudioRepository
}

type AudioRepository struct{}

func NewAudioRepository() *AudioRepository {
	return &AudioRepository{}
}

func (m *AudioModule) Init(ctx *schema.RuntimeContext) error {

	m.ctx = ctx
	m.InitBase(ctx)

	m.writer = NewWAVWriter("/var/data/audio/")
	m.extractor = NewFeatureExtractor(16000)
	m.repo = NewAudioRepository(ctx.DB)

	ctx.Bus.Subscribe("audio.raw", m.Handle)

	return nil
}

func (m *AudioModule) Handle(ctx context.Context, payload []byte) error {

	if err := m.writer.AppendPCM(payload); err != nil {
		return err
	}

	features, err := m.extractor.ProcessPCM(payload)
	if err != nil {
		return err
	}

	data, _ := json.Marshal(features)

	return m.ctx.Data.Bus.Publish(ctx, "audio.features", data)
}

func (m *AudioModule) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}
