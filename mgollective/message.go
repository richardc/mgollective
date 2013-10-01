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

type FactFilter struct {
	Name    string `json:"name"`
	Operand string `json:"operand"`
	Value   string `json:"value"`
}

type IdentityFilter struct {
	Operand string `json:"operand"`
	Value   string `json:"value"`
}

type ClassFilter struct {
	Operand string `json:"operand"`
	Value   string `json:"value"`
}

type RequestFilters struct {
	Facts    []FactFilter     `json:"facts,omitempty"`
	Identity []IdentityFilter `json:"identity,omitempty"`
	Cf_class []ClassFilter    `json:"cf_class,omitempty"`
}

type RequestBody struct {
	Agent  string            `json:"agent"`
	Action string            `json:"action"`
	Params map[string]string `json:"params"`
}

type RequestMessage struct {
	Headers map[string]string `json:"-"`
	Filters RequestFilters    `json:"filters,omitempty"`
	Body    RequestBody       `json:"body"`
}

type ResponseBody map[string]string

type ResponseMessage struct {
	Headers map[string]string `json:"-"`
	Body    ResponseBody      `json:"body"`
}
