package mgollective

import (
	"fmt"
	"os"
)

type DaemonCommand struct {
}

func makeDaemonCommand() Command {
	return &DaemonCommand{}
}

func init() {
	registerCommand("daemon", makeDaemonCommand)
}

func (d *DaemonCommand) Run() {
	config := getconfig("mgo.conf", false)
	connectorname := config.GetStringDefault("connector", "class", "redis")
	var connector Connector
	if factory, exists := connectorRegistry[connectorname]; exists {
		connector = factory(config)
	} else {
		fmt.Printf("No connector called %s", connectorname)
		os.Exit(1)
	}

	connector.Connect()
	connector.Subscribe()

	ch := make(chan Message)
	go connector.Loop(ch)
	for {
		message := <-ch
		logger.Debugf("Recieved %+v", message)
		if agent, exists := agentRegistry[message.Body.Agent]; exists {
			agent(config).Respond(message, connector)
		} else {
			logger.Debugf("No agent '%s'", message.Body.Agent)
		}
	}
}
