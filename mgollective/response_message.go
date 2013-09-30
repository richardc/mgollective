package mgollective

type ResponseMessage struct {
	Target   string
	Headers  map[string]string
	Topic    string
	Reply_to string
	Body     map[string]interface{}
}
