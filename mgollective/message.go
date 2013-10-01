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
	Key     string `json:"key"`
	Operand string `json:"operand"`
	Value   string `json:"value"`
}

type RequestFilters struct {
	Identity []Filter `json:"identity"`
	Agent    []Filter `json:"agent"`
	Compound []Filter `json:"compound"`
	Cf_class []Filter `json:"cf_class"`
}

type RequestBody struct {
	Agent  string            `json:"agent"`
	Action string            `json:"action"`
	Params map[string]string `json:"params"`
}

type RequestMessage struct {
	Headers map[string]string `json:"headers"`
	Filters RequestFilters    `json:"filters"`
	Body    RequestBody       `json:"body"`
}

type ResponseBody map[string]string

type ResponseMessage struct {
	Headers map[string]string `json:"headers"`
	Body    ResponseBody      `json:"body"`
}
