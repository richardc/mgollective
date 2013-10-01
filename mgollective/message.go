package mgollective

type MessageType int

const (
	_                   = iota
	Request MessageType = 1
	Response
)

type WireMessageHeaders map[string]string

type WireMessage struct {
	Headers     WireMessageHeaders
	Body        []byte
	Type        MessageType
	Identity    string
	Target      string
	Destination []string
}

type Filter struct {
	Key     string
	Operand string
	Value   string
}

type RequestFilters struct {
	Identity []Filter
	Agent    []Filter
	Compound []Filter
	Cf_class []Filter
}

type RequestBody struct {
	Agent  string
	Action string
	Params map[string]string
}

type RequestMessage struct {
	Headers map[string]string
	Filters RequestFilters
	Body    RequestBody
}

type ResponseBody map[string]string

type ResponseMessage struct {
	Headers map[string]string
	Body    ResponseBody
}
