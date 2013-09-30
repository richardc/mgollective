package mgollective

type Encoder interface {
	Name() string
	Order() int
	EncodeRequest(RequestMessage) []byte
	DecodeRequest([]byte) RequestMessage
}

type EncoderFactory func(*Mgollective) Encoder

var encoderRegistry = map[string]EncoderFactory{}

func RegisterEncoder(name string, factory EncoderFactory) {
	encoderRegistry[name] = factory
}
