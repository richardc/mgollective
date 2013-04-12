package mgollective

import (
	"log"
	"time"
)

func Discover(connector Connector, config Config, timeout int) []map[string]string {
	log.Println("Discovering nodes")
	discovery := Message{
		target:   "mcollective::server::agents",
		reply_to: config.identity(),
		body: MessageBody{
			Agent:     "discovery",
			Body:      "ping",
			Callerid:  config.callerid(),
			Senderid:  config.senderid(),
			Ttl:       60,
			Msgtime:   time.Now().Unix(),
			Requestid: "42",
		},
	}

	start := time.Now()
	cb := make(chan Message)
	go connector.Loop(cb)
	connector.Publish(discovery)

	nodes := make([]map[string]string, 0)
	for {
		select {
		case message := <-cb:
			log.Printf("got response %+v", message)
			node := map[string]string{
				"senderid": message.body.Senderid,
				"ping":     time.Since(start).String(),
			}
			nodes = append(nodes, node)
			log.Println(node)
		case <-time.After(time.Duration(timeout) * time.Second):
			log.Println("timed out")
			return nodes
		}
	}
	return nodes
}

func PingLoop() {
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

	// Should be a method on *something*.  Probably want to refactor config
	timeout := 3
	nodes := Discover(connector, *config, timeout)
	log.Printf("Discovered %d nodes in %d seconds", len(nodes), timeout)

	for _, node := range nodes {
		log.Printf("ping from %+v", node)
	}
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
