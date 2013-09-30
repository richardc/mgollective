package mgollective

// YAML in go is a little weird.  goyaml.Unmarshal will only marshal into
// public fields (ie Agent)
// As the original YAML is from ruby and is keyed on symbols we make heavy
// use of the doc comment `yaml:":agent"` which says map that into Agent
type MessageBody struct {
	Senderagent string `yaml":senderagent",omitempty`
	Agent       string `yaml:":agent"`
	Filter      struct {
		Identity []string `yaml:"identity"`
		Agent    []string `yaml:"agent"`
		Fact     []string `yaml:"fact"`
		Compound []string `yaml:"compound"`
		Cf_class []string `yaml:"cf_class"`
	} `yaml:":filter"`
	Senderid   string `yaml:":senderid"`
	Collective string `yaml:":collective"`
	Msgtime    int64  `yaml:":msgtime"`
	Ttl        int64  `yaml:":ttl"`
	Requestid  string `yaml:":requestid"`
	Callerid   string `yaml:":callerid"`
	Body       string `yaml:":body"`
}

type Message struct {
	Target   string
	Topic    string
	Reply_to string
	Body     MessageBody
}

/////////////////
// New messages
/////////////////

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
