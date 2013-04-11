package mgollective

type Connector interface {
	Connect()
}

var connectorRegistry = map[string]func(*Config) Connector{}

func registerConnector(name string, connector func(*Config) Connector) {
	connectorRegistry[name] = connector
}
