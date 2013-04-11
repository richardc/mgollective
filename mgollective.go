package mgollective

import (
	"log"
)

func Run() {
	config := getconfig()
	connector := connectorRegistry["redis"](config)
	log.Println(connector)
	connector.Connect()
	connector.Subscribe(config)

	ch := make(chan Message)
	go connector.Loop(ch)
	for {
		message := <-ch
		log.Printf("%+v", message)
		if agent, exists := agentRegistry[message.agent]; exists {
			agent(config).Respond(&message, &connector)
		} else {
			log.Printf("No agent '%s'", message.agent)
		}

	}
}
