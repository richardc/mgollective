package mgollective

type Connector interface {
	Connect()
	Subscribe()
	Publish(map[string]string)
	Loop(chan Message)
}

var connectorRegistry = map[string]func(*Config) Connector{}

func registerConnector(name string, connector func(*Config) Connector) {
	connectorRegistry[name] = connector
}
