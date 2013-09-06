package mgo_connector_activemq

import (
	"crypto/tls"
	"fmt"
	"github.com/gmallard/stompngo"
	"github.com/richardc/mgollective/mgollective"
	"log"
	"net"
)

type ActivemqConnector struct {
	app    *mgollective.Mgollective
	client *stompngo.Connection
}

func (a *ActivemqConnector) Connect() {
	host := a.app.GetStringDefault("connector", "host", "127.0.0.1")
	port := a.app.GetStringDefault("connector", "port", "61613")

	connection, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		log.Fatalln(err) // Handle this ......
	}
	fmt.Println("connected ...")

	if false {
		// do TLS setup
		tlsConfig := new(tls.Config)
		tlsConfig.InsecureSkipVerify = true // Do *not* check the server's certificate
		tls_conn := tls.Client(connection, tlsConfig)
		err = tls_conn.Handshake()
		if err != nil {
			log.Fatalln(err)
		}
		connection = tls_conn
	}

	connection_headers := stompngo.Headers{"accept-version", "1.1"}
	client, err := stompngo.Connect(connection, connection_headers)
	a.client = client
}

func (a *ActivemqConnector) Subscribe() {
}

func (a *ActivemqConnector) Publish(msg mgollective.Message) {
}

func (a *ActivemqConnector) Loop(parsed chan mgollective.Message) {
}

func makeActivemqConnector(app *mgollective.Mgollective) mgollective.Connector {
	return &ActivemqConnector{
		app: app,
	}
}

func init() {
	mgollective.RegisterConnector("activemq", makeActivemqConnector)
}
