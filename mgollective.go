package mgollective

import (
	"log"
	"time"
)

func Discover(connector Connector, config Config, callback func(Message)) {
	log.Println("Discovering nodes")
	discovery := Message{
		target:   "mcollective::server::agents",
		reply_to: config.identity(),
		body: MessageBody{
			Agent:      "discovery",
			Body:       "ping",
			Collective: "mcollective",
			Callerid:   config.callerid(),
			Senderid:   config.senderid(),
			Ttl:        60,
			Msgtime:    time.Now().Unix(),
			Requestid:  "42",
		},
	}

	cb := make(chan Message)
	go connector.Loop(cb)
	connector.Publish(discovery)

	for {
		select {
		case message := <-cb:
			callback(message)
		case <-time.After(3 * time.Second):
			log.Println("timed out")
			return
		}
	}
}

func PingLoop() {
	start := time.Now()
	config := getconfig("mgo.conf", true)
	connectorname := config.GetStringDefault("connector", "class", "redis")
	var connector Connector
	if factory, exists := connectorRegistry[connectorname]; exists {
		connector = factory(config)
	} else {
		log.Fatal("No connector called %s", connectorname)
	}

	log.Println(connector)
	connector.Connect()
	connector.Subscribe()

	nodes := make([]map[string]string, 0)
	// Should be a method on *something*.  Probably want to refactor config
	Discover(connector, *config, func(message Message) {
		node := map[string]string{
			"senderid": message.body.Senderid,
			"ping":     time.Since(start).String(),
		}
		nodes = append(nodes, node)
		log.Println(node)
	})
	log.Printf("Discovered %d nodes", len(nodes))
}

func DaemonLoop() {
	config := getconfig("mgo.conf", false)
	connectorname := config.GetStringDefault("connector", "class", "redis")
	var connector Connector
	if factory, exists := connectorRegistry[connectorname]; exists {
		connector = factory(config)
	} else {
		log.Fatal("No connector called %s", connectorname)
	}

	log.Println(connector)
	connector.Connect()
	connector.Subscribe()

	ch := make(chan Message)
	go connector.Loop(ch)
	for {
		message := <-ch
		log.Printf("Recieved %+v", message)
		if agent, exists := agentRegistry[message.body.Agent]; exists {
			agent(config).Respond(message, connector)
		} else {
			log.Printf("No agent '%s'", message.body.Agent)
		}
	}
}
