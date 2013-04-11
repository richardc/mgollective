package mgollective

import (
	"log"
)

func PingLoop() {
	config := getconfig("mgo.conf")
	connectorname := config.GetStringDefault("connector", "class", "redis")
	var connector Connector
	if factory, exists := connectorRegistry[connectorname]; exists {
		connector = factory(config)
	} else {
		log.Fatal("No connector called %s", connectorname)
	}

	log.Println(connector)
	connector.Connect()

}

func DaemonLoop() {
	config := getconfig("mgo.conf")
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
