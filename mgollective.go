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
		m := <-ch
		log.Printf("%+v", m)
	}
}
