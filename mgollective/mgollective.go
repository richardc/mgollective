package mgollective

// The application context

import (
	"github.com/golang/glog"
	"time"
)

type Mgollective struct {
	Connector Connector
	client    bool
	config    map[string]string
}

func NewClient() Mgollective {
	return NewFromConfigFile("client.cfg", true)
}

func NewFromConfigFile(file string, client bool) Mgollective {
	mgo := Mgollective{
		client: client,
		config: ParseConfig(file),
	}

	connectorname := mgo.GetConfig("connector", "redis")

	if factory, exists := connectorRegistry[connectorname]; exists {
		mgo.Connector = factory(&mgo)
	} else {
		glog.Errorf("No connector called %s", connectorname)
		panic("panic")
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
	return "mcollective::agent::pies"
}

func (m Mgollective) Callerid() string {
	return "user=meat"
}

func (m Mgollective) Senderid() string {
	return "meat.example.com"
}

func (m Mgollective) Discover(callback func(Message)) {
	discovery := Message{
		Target:   m.Collective() + "::server::agents",
		Reply_to: m.Identity(),
		Body: MessageBody{
			Agent:      "discovery",
			Body:       "ping",
			Collective: "mcollective",
			Callerid:   m.Callerid(),
			Senderid:   m.Senderid(),
			Ttl:        60,
			Msgtime:    time.Now().Unix(),
			Requestid:  "42",
		},
	}

	cb := make(chan Message)
	go m.Connector.Loop(cb)
	m.Connector.Publish(discovery)

	for {
		select {
		case message := <-cb:
			callback(message)
		case <-time.After(3 * time.Second):
			return
		}
	}
}

func (m Mgollective) RpcCommand(agent, command string, params map[string]string, callback func(ResponseMessage)) {
	glog.Info("sending RpcCommand")
	responses := make(chan ResponseMessage)

	for {
		select {
		case message := <-responses:
			callback(message)
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
