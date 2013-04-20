package mgollective

type Connector interface {
	Connect()
	Subscribe()
	Publish(Message)
	Loop(chan Message)
}

type ConnectorFactory func(*Mgollective) Connector

var connectorRegistry = map[string]ConnectorFactory{}

func RegisterConnector(name string, factory ConnectorFactory) {
	connectorRegistry[name] = factory
}
