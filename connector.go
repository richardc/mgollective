package mgollective

type Connector interface {
	Connect() int
}

var connectorRegistry = map[string]func() Connector{}

func registerConnector(name string, connector func() Connector) {
	connectorRegistry[name] = connector
}
