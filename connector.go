package mgollective

type Connector interface {
	Connect() int
}

var connectorRegistry = map[string]func(*Config) Connector{}

func registerConnector(name string, connector func(*Config) Connector) {
	connectorRegistry[name] = connector
}
