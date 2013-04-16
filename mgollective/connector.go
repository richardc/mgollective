package mgollective

type Connector interface {
	Connect()
	Subscribe()
	Publish(Message)
	Loop(chan Message)
}

var connectorRegistry = map[string]func(*Config) Connector{}

func RegisterConnector(name string, connector func(*Config) Connector) {
	connectorRegistry[name] = connector
}