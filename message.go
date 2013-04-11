package mgollective

type Message struct {
	topic   string
	agent   string
	headers interface{}
	body    interface{}
}
