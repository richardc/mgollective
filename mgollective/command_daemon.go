package mgollective

import (
	"code.google.com/p/go-commander"
	"github.com/golang/glog"
)

func runDaemonCommand(cmd *commander.Command, args []string) {
	mgo := NewServer()

	ch := make(chan WireMessage)
	go mgo.Connector.RecieveLoop(ch)
	for {
		message := <-ch
		glog.Infof("Recieved %+v", message)
		if mgo.verifyMessage(message) {
			request := mgo.decodeRequest(message)
			agentname := request.Body.Agent
			if agent, exists := agentRegistry[agentname]; exists {
				response := agent.Respond(mgo, request)
				if response == nil {
					glog.Infof("No response from agent %s", agentname)
				} else {
					glog.Infof("Sending response %v", response)
					wire_response := mgo.encodeResponse(*response)
					mgo.signMessage(&wire_response)
					wire_response.Target = request.Headers["reply-to"]
					go mgo.Connector.PublishResponse(wire_response)
				}
			} else {
				glog.Infof("No agent '%s'", agentname)
			}
		} else {
			glog.Info("Message failed verification")
		}

	}
}

func init() {
	RegisterCommand(&commander.Command{
		UsageLine: "daemon",
		Run:       runDaemonCommand,
	})
}
