package mgollective

type Agent interface {
	Respond(Message, Connector)
}

type AgentFactory func(*Mgollective) Agent

var agentRegistry = map[string]AgentFactory{}

func registerAgent(name string, factory AgentFactory) {
	agentRegistry[name] = factory
}
