package mgollective

import (
    "log"
)

func Run() {
    conf := getconfig()
    log.Println(conf.GetString("connector", "class"))
    connector := connectorRegistry["redis"]()
    log.Println(connector)
    connector.Connect()
}
