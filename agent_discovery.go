package mgollective

import (
	"log"
)

type DiscoveryAgent struct {
	config *Config
}

func (a *DiscoveryAgent) matches(msg Message) bool {
	return true
}

func (a *DiscoveryAgent) Respond(msg Message, connector Connector) {
	log.Printf("Discover agent handling %+v", msg)
	if !a.matches(msg) {
		log.Printf("Not for us")
		return
	}

	reply := Message{target: msg.reply_to}
	if msg.Body == "ping" {
		reply.Body = "pong"
	} else {
		reply.Body = "Unknown Request: " + msg.Body
	}

	log.Printf("Going to send %+v", reply)
	connector.Publish(reply)
}

func makeDiscoveryAgent(c *Config) Agent {
	return &DiscoveryAgent{config: c}
}

func init() {
	registerAgent("discovery", makeDiscoveryAgent)
}
