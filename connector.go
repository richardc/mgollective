package mgollective

type Connector interface {
	Connect()
	Subscribe(*Config)
	Loop(chan Message)
}

type Message struct {
	topic   string
	headers interface{}
	body    interface{}
}

var connectorRegistry = map[string]func(*Config) Connector{}

func registerConnector(name string, connector func(*Config) Connector) {
	connectorRegistry[name] = connector
}
