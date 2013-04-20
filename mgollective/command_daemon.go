package mgollective

import (
	"fmt"
	"github.com/maruel/subcommands"
)

type DaemonCommand struct {
	subcommands.CommandRunBase
}

func init() {
	RegisterCommand(&subcommands.Command{
		UsageLine:  "daemon",
		CommandRun: func() subcommands.CommandRun { return &DaemonCommand{} },
	})
}

func (d *DaemonCommand) Run(a subcommands.Application, args []string) int {
	config := getconfig("mgo.conf", false)
	connectorname := config.GetStringDefault("connector", "class", "redis")
	var connector Connector
	if factory, exists := connectorRegistry[connectorname]; exists {
		connector = factory(config)
	} else {
		fmt.Printf("No connector called %s", connectorname)
		return 1
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
	return 0
}
