package mgollective

import (
	"github.com/golang/glog"
	//"time"
)

type DiscoveryAgent struct {
	app *Mgollective
}

func (a *DiscoveryAgent) matches(msg RequestMessage) bool {
	return true
}

func (agent *DiscoveryAgent) Respond(msg RequestMessage) *ResponseMessage {
	glog.Infof("Discover agent handling %+v", msg)
	if !agent.matches(msg) {
		glog.Infof("Not for us")
		return nil
	}
	var body string
	if msg.Body.Action == "ping" {
		body = "pong"
	} else {
		body = "Unknown Request: " + msg.Body.Agent
	}

	glog.Info(body)
	return nil

	// reply := Message{
	// 	Target: msg.Reply_to,
	// 	Body: MessageBody{
	// 		Requestid:   msg.Body.Requestid,
	// 		Senderagent: "discovery",
	// 		Senderid:    agent.app.Senderid(),
	// 		Msgtime:     time.Now().Unix(),
	// 		Body:        body,
	// 	},
	// }

	// connector.Publish(reply)
}

func makeDiscoveryAgent(m *Mgollective) Agent {
	return &DiscoveryAgent{app: m}
}

func init() {
	registerAgent("discovery", makeDiscoveryAgent)
}
