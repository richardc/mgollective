package mgollective

type Message struct {
	topic   string
	headers interface{}
	body    interface{}
}
