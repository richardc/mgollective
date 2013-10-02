package mgollective

// The application context

import (
	"github.com/golang/glog"
	"time"
)

type Mgollective struct {
	Connector        Connector
	Encoder          Encoder
	SecurityProvider SecurityProvider
	Factsource       Factsource
	client           bool
	config           map[string]string
}

func NewClient() Mgollective {
	return newFromConfigFile(client_config_file, true)
}

func NewServer() Mgollective {
	return newFromConfigFile(server_config_file, false)
}

func newFromConfigFile(file string, client bool) Mgollective {
	mgo := Mgollective{
		client: client,
		config: ParseConfig(file),
	}

	connectorname := mgo.GetConfig("connector", "redis")

	if factory, exists := connectorRegistry[connectorname]; exists {
		mgo.Connector = factory(&mgo)
	} else {
		glog.Fatalf("No connector called %s", connectorname)
	}

	if factory, exists := encoderRegistry["json"]; exists {
		mgo.Encoder = factory(&mgo)
	} else {
		glog.Fatalf("No encoder called %s", "json")
	}

	securityprovider := mgo.GetConfig("securityprovider", "null")
	if factory, exists := securityProviderRegistry[securityprovider]; exists {
		mgo.SecurityProvider = factory(&mgo)
	} else {
		glog.Fatalf("No securityprovider called %s", securityprovider)
	}

	mgo.Connector.Connect()
	mgo.Connector.Subscribe()

	return mgo
}

func (m Mgollective) Shutdown() {
	m.Connector.Unsubscribe()
	m.Connector.Disconnect()
}

func (m Mgollective) GetConfig(name, def string) string {
	if value, ok := m.config[name]; ok {
		return value
	} else {
		return def
	}
}
func (m Mgollective) IsClient() bool {
	return m.client
}

func (m Mgollective) Collectives() []string {
	return []string{"mcollective"}
}

func (m Mgollective) Collective() string {
	return "mcollective"
}

func (m Mgollective) Classes() []string {
	return []string{"mgollective"}
}

func (m Mgollective) Identity() string {
	return "foo.example.com"
}

func (m Mgollective) Callerid() string {
	return "user=meat"
}

func (m Mgollective) Senderid() string {
	return "meat.example.com"
}

func (m Mgollective) Discover(callback func(ResponseMessage)) {
	discovery := RequestMessage{
		Body: RequestBody{
			Agent:  "discovery",
			Action: "ping",
		},
	}

	cb := make(chan ResponseMessage)
	//	go m.Connector.RecieveLoop(cb)
	glog.Info(discovery)
	//m.Connector.Publish(discovery)

	for {
		select {
		case message := <-cb:
			callback(message)
		case <-time.After(3 * time.Second):
			return
		}
	}
}

func (m Mgollective) encodeResponse(message ResponseMessage) WireMessage {
	msg := WireMessage{
		Type: Response,
		Headers: WireMessageHeaders{
			"protocol_version": "2",
			"encoding":         m.Encoder.Name(),
		},
		Body: m.Encoder.EncodeResponse(message),
	}
	return msg
}

func (m Mgollective) decodeResponse(message WireMessage) ResponseMessage {
	/// XXX select encoder
	msg := m.Encoder.DecodeResponse(message.Body)
	msg.Headers = message.Headers
	glog.Infof("Decoded to %v", msg)
	return msg
}

func (m Mgollective) encodeRequest(message RequestMessage) WireMessage {
	msg := WireMessage{
		Type: Request,
		Headers: WireMessageHeaders{
			"mc_protocol_version": "2",
			"mc_encoding":         m.Encoder.Name(),
			"mc_accepts_encoding": m.Encoder.Name(), // XXX should be lookup
			"mc_requestid":        "blarb",
		},
		Body: m.Encoder.EncodeRequest(message),
	}
	return msg
}

func (m Mgollective) decodeRequest(message WireMessage) RequestMessage {
	/// XXX select encoder
	msg := m.Encoder.DecodeRequest(message.Body)
	msg.Headers = message.Headers
	glog.Infof("Decoded to %v", msg)
	return msg
}

func (m Mgollective) signMessage(message *WireMessage) {
	signature := m.SecurityProvider.Sign(message.Body)
	for k, v := range signature {
		message.Headers[k] = v
	}
}

func (m Mgollective) verifyMessage(message WireMessage) bool {
	return m.SecurityProvider.Verify(message.Body, message.Headers)
}

func (m Mgollective) RpcCommand(request RequestMessage, discovered_nodes []string, callback func(ResponseMessage)) {
	glog.Info("sending RpcCommand %v", request)
	responses := make(chan WireMessage)
	msg := m.encodeRequest(request)
	m.signMessage(&msg)

	msg.Destination = discovered_nodes

	go m.Connector.PublishRequest(msg)
	go m.Connector.RecieveLoop(responses)

	response_count := 0
	for {
		select {
		case message := <-responses:
			if m.verifyMessage(message) {
				msg := m.decodeResponse(message)
				callback(msg)
				response_count++
				if response_count >= len(discovered_nodes) {
					glog.Info("got all responses, quitting")
					return
				}
			} else {
				glog.Error("Message didn't validate")
			}
		case <-time.After(10 * time.Second):
			glog.Info("timing out")
			return
		}
	}
}

func init() {
	DeclareConfig("main_collective")
	DeclareConfig("securityprovider")
	DeclareConfig("connector")
}
