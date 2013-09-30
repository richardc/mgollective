package mgollective

type Encoder interface {
	Name() string
	Order() int
	Encode(map[string]string) []byte
	Decode([]byte) map[string]string
}

type EncoderFactory func(*Mgollective) Encoder

var encoderRegistry = map[string]EncoderFactory{}

func RegisterEncoder(name string, factory EncoderFactory) {
	encoderRegistry[name] = factory
}

func AcceptedEncodings() string {
	return "json"
}
