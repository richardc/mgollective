package mgollective

import (
	"github.com/golang/glog"
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
	mgo := NewFromConfigFile("server.cfg", false)

	ch := make(chan Message)
	// go mgo.Connector.Loop(ch)
	for {
		message := <-ch
		glog.Infof("Recieved %+v", message)
		if agent, exists := agentRegistry[message.Body.Agent]; exists {
			agent(&mgo).Respond(message, mgo.Connector)
		} else {
			glog.Infof("No agent '%s'", message.Body.Agent)
		}
	}
	return 0
}
