package mgo_connector_activemq

import (
	"crypto/tls"
	"fmt"
	"github.com/gmallard/stompngo"
	"github.com/golang/glog"
	"github.com/richardc/mgollective/mgollective"
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
	glog.Info("connecting to activemq", host, port)

	connection, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		glog.Fatal(err) // Handle this ......
	}
	glog.Info("connected ...")

	if a.app.GetConfig("plugin.activemq.pool.1.ssl", "0") == "1" {
		glog.Info("starting TLS")
		tlsConfig := new(tls.Config)
		tlsConfig.InsecureSkipVerify = true // Do *not* check the server's certificate
		tls_conn := tls.Client(connection, tlsConfig)
		err = tls_conn.Handshake()
		if err != nil {
			glog.Fatal(err)
		}
		connection = tls_conn
		glog.Info("TLS configured")
	}

	user := a.app.GetConfig("plugin.activemq.pool.1.user", "")
	connection_headers := stompngo.Headers{
		"login", user,
		"passcode", a.app.GetConfig("plugin.activemq.pool.1.password", ""),
	}

	glog.Info("logging in as ", user)
	client, err := stompngo.Connect(connection, connection_headers)
	if err != nil {
		glog.Fatal(err)
	}
	glog.Info("logged in")

	a.client = client
}

func (a *ActivemqConnector) Disconnect() {
	glog.Info("disconnecting")
	eh := stompngo.Headers{}
	e := a.client.Disconnect(eh)
	if e != nil {
		glog.Fatal(e)
	}
}

func (a *ActivemqConnector) Subscribe() {
	glog.Info("subscribing to channels")
	if a.app.IsClient() {
		a.reply_to = fmt.Sprintf("/queue/%s.reply.%s_%d", a.app.Collective(), a.app.Senderid(), os.Getpid())
		glog.Info("subscribing to '" + a.reply_to + "'")
		sub := stompngo.Headers{
			"destination", a.reply_to,
		}
		channel, err := a.client.Subscribe(sub)
		if err != nil {
			glog.Fatal(sub)
		}
		a.channels = make(map[string]<-chan stompngo.MessageData)
		a.channels["reply"] = channel
	} else {
		glog.Fatal("server connections not yet supported")
	}
}

func (a *ActivemqConnector) Unsubscribe() {
	glog.Info("unsubscribing from channels")
	if a.app.IsClient() {
		glog.Info("Unsubscribing from '" + a.reply_to + "'")
		sub := stompngo.Headers{
			"destination", a.reply_to,
		}
		err := a.client.Unsubscribe(sub)
		if err != nil {
			glog.Fatal(err)
		}
	}
}

func (a *ActivemqConnector) Publish(msg mgollective.Message) {
	glog.Info("publishing message", msg)

}

func (a *ActivemqConnector) Loop(parsed chan mgollective.Message) {
	glog.Info("entering recieve loop")
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
