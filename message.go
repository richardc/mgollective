package mgollective

// YAML in go is a little weird.  goyaml.Unmarshal will only marshal into
// public fields (ie Agent)
// As the original YAML is from ruby and is keyed on symbols we make heavy
// use of the doc comment `yaml:":agent"` which says map that into Agent
type Message struct {
	target   string
	topic    string
	reply_to string
	Agent    string `yaml:":agent"`
	Filter   struct {
		Identify []string `yaml:":identify"`
		Agent    []string `yaml:":agent"`
		Fact     []string `yaml:":fact"`
		Compound []string `yaml:":compound"`
		Cf_class []string `yaml:"cf_class"`
	} `yaml:":filter"`
	Senderid   string `yaml:":senderid"`
	Collective string `yaml:":collective"`
	Msgtime    string `yaml:":msgtime"`
	Ttl        string `yaml:":ttl"`
	Requestid  string `yaml:":requestid"`
	Callerid   string `yaml:":callerid"`
	Body       string `yaml:":body"`
}
