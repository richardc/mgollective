package mgollective

// The application context

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/msbranco/goconfig"
	"log"
	"time"
)

type Mgollective struct {
	Connector Connector
	client    bool
	config    goconfig.ConfigFile
	logger    seelog.LoggerInterface
}

func NewFromConfigFile(file string, client bool) Mgollective {
	conf, err := goconfig.ReadConfigFile(file)
	if err != nil {
		log.Fatal("Couldn't open config ", err)
	}

	mgo := Mgollective{
		logger: seelog.Disabled,
		client: client,
		config: *conf,
	}

	connectorname := mgo.GetStringDefault("connector", "class", "redis")

	if factory, exists := connectorRegistry[connectorname]; exists {
		mgo.Connector = factory(&mgo)
	} else {
		fmt.Printf("No connector called %s", connectorname)
		panic("panic")
	}

	mgo.Connector.Connect()
	mgo.Connector.Subscribe()

	return mgo
}

func (m Mgollective) GetStringDefault(section, variable, def string) string {
	value, err := m.config.GetString(section, variable)
	if err != nil {
		return def
	}
	return value
}

func (m Mgollective) GetIntDefault(section, variable string, def int) int {
	value, err := m.config.GetInt64(section, variable)
	if err != nil {
		return def
	}
	return int(value)
}

func (m Mgollective) IsClient() bool {
	return m.client
}

// Explicit delegation over to seelog or whatever
func (m Mgollective) Error(args ...interface{}) {
	m.logger.Error(args)
}

func (m Mgollective) Errorf(fmt string, args ...interface{}) {
	m.logger.Errorf(fmt, args)
}

func (m Mgollective) Info(args ...interface{}) {
	m.logger.Info(args)
}

func (m Mgollective) Infof(fmt string, args ...interface{}) {
	m.logger.Infof(fmt, args)
}

func (m Mgollective) Debug(args ...interface{}) {
	m.logger.Debug(args)
}

func (m Mgollective) Debugf(fmt string, args ...interface{}) {
	m.logger.Debugf(fmt, args)
}

func (m Mgollective) Trace(args ...interface{}) {
	m.logger.Trace(args)
}

func (m Mgollective) Tracef(fmt string, args ...interface{}) {
	m.logger.Tracef(fmt, args)
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
