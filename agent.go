package mgollective

type Agent interface {
	Respond(*Message, *Connector)
}

var agentRegistry = map[string]func(*Config) Agent{}

func registerAgent(name string, agent func(*Config) Agent) {
	agentRegistry[name] = agent
}
