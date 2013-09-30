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

	ch := make(chan WireMessage)
	go mgo.Connector.RecieveLoop(ch)
	for {
		message := <-ch
		glog.Infof("Recieved %+v", message)
		if agent, exists := agentRegistry[message.Headers["agent"]]; exists {
			//agent(&mgo).Respond(message, mgo.Connector)
			//agent(&mgo).Name
			glog.Info(message, agent)
		} else {
			glog.Infof("No agent '%s'", message.Headers["agent"])
		}
	}
	return 0
}
