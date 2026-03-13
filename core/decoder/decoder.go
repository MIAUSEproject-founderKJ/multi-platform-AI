// core/decoder/decoder.go
package decoder

type Decoder interface {
	Decode([]byte) (*ExternalEvent, error)
}
