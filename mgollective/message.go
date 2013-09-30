package mgollective

type WireMessage struct {
	Headers map[string]string
	Body    []byte
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
	Filters RequestFilters
	Body    RequestBody
}

type ResponseMessage struct {
	Headers map[string]string
	Body    map[string]string
}
