package mgollective

import (
	"log"
)

func Run() {
	conf := getconfig()
	connector := connectorRegistry["redis"](conf)
	log.Println(connector)
	connector.Connect()
}
