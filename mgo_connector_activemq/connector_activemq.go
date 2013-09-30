package mgo_connector_activemq

import (
	"crypto/tls"
	"fmt"
	"github.com/gmallard/stompngo"
	"github.com/richardc/mgollective/mgollective"
	"log"
	"net"
	"os"
)

type ActivemqConnector struct {
	app      *mgollective.Mgollective
	client   *stompngo.Connection
	reply_to string
	channels map[string]<-chan stompngo.MessageData
}

func (a *ActivemqConnector) Connect() {
	host := a.app.GetConfig("plugin.activemq.pool.1.host", "127.0.0.1")
	port := a.app.GetConfig("plugin.activemq.pool.1.port", "61613")
	log.Println("connecting to activemq", host, port)

	connection, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		log.Fatalln(err) // Handle this ......
	}
	log.Println("connected ...")

	if a.app.GetConfig("plugin.activemq.pool.1.ssl", "0") == "1" {
		log.Println("starting TLS")
		tlsConfig := new(tls.Config)
		tlsConfig.InsecureSkipVerify = true // Do *not* check the server's certificate
		tls_conn := tls.Client(connection, tlsConfig)
		err = tls_conn.Handshake()
		if err != nil {
			log.Fatalln(err)
		}
		connection = tls_conn
		log.Println("TLS configured")
	}

	user := a.app.GetConfig("plugin.activemq.pool.1.user", "")
	connection_headers := stompngo.Headers{
		"login", user,
		"passcode", a.app.GetConfig("plugin.activemq.pool.1.password", ""),
	}

	log.Println("logging in as ", user)
	client, err := stompngo.Connect(connection, connection_headers)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("logged in")

	a.client = client
}

func (a *ActivemqConnector) Disconnect() {
	log.Println("disconnecting")
	eh := stompngo.Headers{}
	e := a.client.Disconnect(eh)
	if e != nil {
		log.Fatalln(e)
	}
}

func (a *ActivemqConnector) Subscribe() {
	log.Println("subscribing to channels")
	if a.app.IsClient() {
		a.reply_to = fmt.Sprintf("/queue/%s.reply.%s_%d", a.app.Collective(), a.app.Senderid(), os.Getpid())
		log.Println("subscribing to '" + a.reply_to + "'")
		sub := stompngo.Headers{
			"destination", a.reply_to,
		}
		channel, err := a.client.Subscribe(sub)
		if err != nil {
			log.Fatalln(sub)
		}
		a.channels = make(map[string]<-chan stompngo.MessageData)
		a.channels["reply"] = channel
	} else {
		log.Fatalln("server connections not yet supported")
	}
}

func (a *ActivemqConnector) Unsubscribe() {
	log.Println("unsubscribing from channels")
	if a.app.IsClient() {
		log.Println("Unsubscribing from '" + a.reply_to + "'")
		sub := stompngo.Headers{
			"destination", a.reply_to,
		}
		err := a.client.Unsubscribe(sub)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func (a *ActivemqConnector) Publish(msg mgollective.Message) {
	log.Println("publishing message", msg)

}

func (a *ActivemqConnector) Loop(parsed chan mgollective.Message) {
	log.Println("entering recieve loop")
}

func makeActivemqConnector(app *mgollective.Mgollective) mgollective.Connector {
	return &ActivemqConnector{
		app: app,
	}
}

func init() {
	mgollective.DeclareConfig("direct_addressing")
	mgollective.DeclareConfig("plugin.activemq.base64")
	mgollective.DeclareConfig("plugin.activemq.pool.size")
	mgollective.DeclareConfig("plugin.activemq.pool.randomize")
	mgollective.DeclareConfig("plugin.activemq.pool.*.host")
	mgollective.DeclareConfig("plugin.activemq.pool.*.port")
	mgollective.DeclareConfig("plugin.activemq.pool.*.user")
	mgollective.DeclareConfig("plugin.activemq.pool.*.password")
	mgollective.RegisterConnector("activemq", makeActivemqConnector)
}
