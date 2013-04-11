package mgollective

import (
	"log"
	"time"
)

// top level app scope
type App struct {
	config    Config
	connector Connector
}

func Discover(connector Connector, config Config, timeout int) []map[string]string {
	log.Println("Discovering nodes")
	discovery := make(map[string]interface{}, 0)
	discovery["target"] = "mcollective::server::agents"
	discovery["reply-to"] = config.identity()
	discovery[":agent"] = "discovery"
	discovery[":body"] = "ping"
	discovery[":callerid"] = config.callerid()
	discovery[":senderid"] = config.senderid()
	discovery[":ttl"] = 60
	discovery[":msgtime"] = time.Now().Unix()
	filters := map[string][]string{
		"identity": {},
		"agent":    {},
		"fact":     {},
		"compound": {},
		"cf_class": {},
	}
	discovery[":requestid"] = "42"
	discovery[":filter"] = filters
	connector.Publish(discovery)
	cb := make(chan Message)
	go connector.Loop(cb)
	nodes := make([]map[string]string, 0)
	for {
		select {
		case message := <-cb:
			log.Printf("got response %+v", message)
			node := make(map[string]string, 0)
			nodes = append(nodes, node)
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
	timeout := 5
	nodes := Discover(connector, *config, timeout)
	log.Printf("Discovered %d nodes in %d seconds", len(nodes), timeout)

	message := make(map[string]interface{})

	for _, node := range nodes {
		log.Printf("sending ping to", node)
		message["target"] = node["target"]
		connector.Publish(message)
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
		if agent, exists := agentRegistry[message.Agent]; exists {
			agent(config).Respond(message, connector)
		} else {
			log.Printf("No agent '%s'", message.Agent)
		}
	}
}
