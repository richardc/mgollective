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
	app           *mgollective.Mgollective
	socket        net.Conn
	client        *stompngo.Connection
	internal_chan chan mgollective.WireMessage
}

func (c *ActivemqConnector) Connect() {
	host := c.app.GetConfig("plugin.activemq.pool.1.host", "127.0.0.1")
	port := c.app.GetConfig("plugin.activemq.pool.1.port", "61613")
	glog.Info("connecting to activemq", host, port)

	connection, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		glog.Fatal(err) // Handle this ......
	}
	glog.Info("connected ...")
	c.socket = connection

	if c.app.GetConfig("plugin.activemq.pool.1.ssl", "0") == "1" {
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

	user := c.app.GetConfig("plugin.activemq.pool.1.user", "")
	connection_headers := stompngo.Headers{
		"login", user,
		"passcode", c.app.GetConfig("plugin.activemq.pool.1.password", ""),
	}

	glog.Info("logging in as ", user)
	client, err := stompngo.Connect(connection, connection_headers)
	if err != nil {
		glog.Fatal(err)
	}
	glog.Info("logged in")

	c.client = client
}

func (c *ActivemqConnector) Disconnect() {
	glog.Info("disconnecting")
	eh := stompngo.Headers{}
	e := c.client.Disconnect(eh)
	if e != nil {
		glog.Fatal(e)
	}

	glog.Info("closing socket")
	e = c.socket.Close()
	if e != nil {
		glog.Fatal(e)
	}
	glog.Info("socket closed")
}

func (c *ActivemqConnector) Subscribe() {
	glog.Info("subscribing to channels")

	var queue string
	if c.app.IsClient() {
		queue = fmt.Sprintf("/queue/%s.reply.%s_%d", c.app.Collective(), c.app.Senderid(), os.Getpid())
	} else {
		queue = fmt.Sprintf("/queue/%s.nodes", c.app.Collective())
	}

	sub := stompngo.Headers{
		"destination", queue,
	}
	glog.Info("subscribing with headers %v", sub)
	channel, err := c.client.Subscribe(sub)
	if err != nil {
		glog.Fatal(sub)
	}
	go c.recieve(channel)
}

// This recieves from a channel of stompnogo.MessageData, wraps it in a
// mgollective.WireMessage and passes it on to the internal channel
// that we recieve on
func (c *ActivemqConnector) recieve(channel <-chan stompngo.MessageData) {
	for messagedata := range channel {
		glog.Infof("STOMP recieved %+v", messagedata)
		headers := make(map[string]string)
		for i, key := range messagedata.Message.Headers {
			if i%2 == 0 {
				headers[key] = messagedata.Message.Headers.Value(key)
			}
		}
		wire := mgollective.WireMessage{
			Headers: headers,
			Body:    messagedata.Message.Body,
		}
		c.internal_chan <- wire
	}
}

func (c *ActivemqConnector) Unsubscribe() {
	glog.Info("unsubscribing from channels")
	var queue string
	if c.app.IsClient() {
		queue = fmt.Sprintf("/queue/%s.reply.%s_%d", c.app.Collective(), c.app.Senderid(), os.Getpid())
	} else {
		queue = fmt.Sprintf("/queue/%s.nodes", c.app.Collective())
	}

	sub := stompngo.Headers{
		"destination", queue,
	}
	glog.Info("Unsubscribing from %v", sub)

	err := c.client.Unsubscribe(sub)
	if err != nil {
		glog.Fatal(err)
	}
}

func (c *ActivemqConnector) Publish(queue string, destinations []string, msg mgollective.WireMessage) {
	for _, destination := range destinations {
		headers := stompngo.Headers{
			"destination", queue,
			"mc_identity", destination,
		}
		for k, v := range msg.Headers {
			headers = headers.Add(k, v)
		}
		glog.Infof("publishing message on %s with headers %v", queue, headers)
		err := c.client.Send(headers, string(msg.Body))
		if err != nil {
			glog.Fatalln(err)
		}
	}
}

func (c *ActivemqConnector) RecieveLoop(parsed chan mgollective.WireMessage) {
	glog.Info("entering recieve loop")
	for {
		message := <-c.internal_chan
		glog.Info("recieved %v", message)
		parsed <- message
	}
}

func makeActivemqConnector(app *mgollective.Mgollective) mgollective.Connector {
	return &ActivemqConnector{
		app:           app,
		internal_chan: make(chan mgollective.WireMessage),
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
