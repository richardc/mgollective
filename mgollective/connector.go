package mgollective

type Connector interface {
	Connect()
	Subscribe()
	Unsubscribe()
	Disconnect()
	PublishRequest(message WireMessage)
	PublishResponse(message WireMessage)

	RecieveLoop(chan WireMessage)
}

type ConnectorFactory func(*Mgollective) Connector

var connectorRegistry = map[string]ConnectorFactory{}

func RegisterConnector(name string, factory ConnectorFactory) {
	connectorRegistry[name] = factory
}
