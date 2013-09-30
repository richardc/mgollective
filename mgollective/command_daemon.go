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

func (c *DaemonCommand) Run(a subcommands.Application, args []string) int {
	mgo := NewFromConfigFile("server.cfg", false)

	ch := make(chan WireMessage)
	go mgo.Connector.RecieveLoop(ch)
	for {
		message := <-ch
		glog.Infof("Recieved %+v", message)
		if mgo.verifyMessage(message) {
			request := mgo.decodeRequest(message)
			agentname := request.Body.Agent
			if agent, exists := agentRegistry[agentname]; exists {
				response := agent(&mgo).Respond(request)
				if response == nil {
					glog.Infof("No response from agent %s", agentname)
				} else {
					glog.Infof("Sending response %v", response)
					/// XXX actually send a response
				}
			} else {
				glog.Infof("No agent '%s'", agentname)
			}
		} else {
			glog.Info("Message failed verification")
		}

	}
	return 0
}
