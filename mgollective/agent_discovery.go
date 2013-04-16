package mgollective

import (
	"time"
)

type DiscoveryAgent struct {
	config *Config
}

func (a *DiscoveryAgent) matches(msg Message) bool {
	return true
}

func (a *DiscoveryAgent) Respond(msg Message, connector Connector) {
	logger.Infof("Discover agent handling %+v", msg)
	if !a.matches(msg) {
		logger.Debugf("Not for us")
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
			Senderid:    a.config.Senderid(),
			Msgtime:     time.Now().Unix(),
			Body:        body,
		},
	}

	connector.Publish(reply)
}

func makeDiscoveryAgent(c *Config) Agent {
	return &DiscoveryAgent{config: c}
}

func init() {
	registerAgent("discovery", makeDiscoveryAgent)
}
