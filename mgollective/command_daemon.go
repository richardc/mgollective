package mgollective

import (
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
	mgo := NewFromConfigFile("mgo.conf", false)

	ch := make(chan Message)
	go mgo.Connector.Loop(ch)
	for {
		message := <-ch
		mgo.Debugf("Recieved %+v", message)
		if agent, exists := agentRegistry[message.Body.Agent]; exists {
			agent(&mgo).Respond(message, mgo.Connector)
		} else {
			mgo.Debugf("No agent '%s'", message.Body.Agent)
		}
	}
	return 0
}
