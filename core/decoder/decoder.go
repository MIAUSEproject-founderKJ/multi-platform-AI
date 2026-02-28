//core/decoder/decoder.go

type Decoder interface {
    Decode([]byte) (*ExternalEvent, error)
}