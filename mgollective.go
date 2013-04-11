package mgollective

import (
	"log"
)

func Run() {
	conf := getconfig()
	log.Println(conf.GetString("connector", "class"))
	connector := connectorRegistry["redis"](conf)
	log.Println(connector)
	connector.Connect()
}
