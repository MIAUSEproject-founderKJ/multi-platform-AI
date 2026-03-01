//modules/audio/audio_module.go

type AudioModule struct {
    writer     *WAVWriter
    extractor  *FeatureExtractor
    repo       AudioRepository
}

func (m *AudioModule) Name() string {
    return "AudioModule"
}

func (m *AudioModule) DependsOn() []string {
    return []string{"StorageModule"}
}

func (m *AudioModule) Init(ctx *runtime.ExecutionContext) error {
    m.writer = NewWAVWriter("/var/data/audio/")
    m.extractor = NewFeatureExtractor(16000)
    m.repo = NewAudioRepository(ctx.DB)
    return nil
}

func (m *AudioModule) Run(ctx context.Context) error {
    <-ctx.Done()
    return nil
}

func (m *AudioModule) Handle(ctx context.Context, payload []byte) error {

    // 1. Save raw PCM to WAV
    if err := m.writer.AppendPCM(payload); err != nil {
        return err
    }

    // 2. Extract features
    features, err := m.extractor.ProcessPCM(payload)
    if err != nil {
        return err
    }

    // 3. Store features
    return m.repo.InsertFeatures(ctx, features)
}