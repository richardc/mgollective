package mgollective

import (
	"time"
)

type DiscoveryAgent struct {
	app *Mgollective
}

func (a *DiscoveryAgent) matches(msg Message) bool {
	return true
}

func (agent *DiscoveryAgent) Respond(msg Message, connector Connector) {
	agent.app.Infof("Discover agent handling %+v", msg)
	if !agent.matches(msg) {
		agent.app.Debugf("Not for us")
		return
	}
	var body string
	if msg.Body.Body == "ping" {
		body = "pong"
	} else {
		body = "Unknown Request: " + msg.Body.Body
	}

	reply := Message{
		Target: msg.Reply_to,
		Body: MessageBody{
			Requestid:   msg.Body.Requestid,
			Senderagent: "discovery",
			Senderid:    agent.app.Senderid(),
			Msgtime:     time.Now().Unix(),
			Body:        body,
		},
	}

	connector.Publish(reply)
}

func makeDiscoveryAgent(m *Mgollective) Agent {
	return &DiscoveryAgent{app: m}
}

func init() {
	registerAgent("discovery", makeDiscoveryAgent)
}
