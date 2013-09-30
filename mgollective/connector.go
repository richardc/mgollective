package mgollective

type Connector interface {
	Connect()
	Subscribe()
	Unsubscribe()
	Disconnect()
	Publish(queue string, destinations []string, message WireMessage)
	RecieveLoop(chan WireMessage)
}

type ConnectorFactory func(*Mgollective) Connector

var connectorRegistry = map[string]ConnectorFactory{}

func RegisterConnector(name string, factory ConnectorFactory) {
	connectorRegistry[name] = factory
}
