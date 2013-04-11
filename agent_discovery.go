package mgollective

import (
	"log"
	"time"
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
	var body string
	if msg.Body == "ping" {
		body = "pong"
	} else {
		body = "Unknown Request: " + msg.Body
	}

	reply := map[string]interface{}{
		"target":       msg.reply_to,
		":requestid":   msg.Requestid,
		":senderagent": "discovery",
		":senderid":    a.config.senderid(),
		":msgtime":     time.Now().Unix(),
		"body":         body,
	}

	connector.Publish(reply)
}

func makeDiscoveryAgent(c *Config) Agent {
	return &DiscoveryAgent{config: c}
}

func init() {
	registerAgent("discovery", makeDiscoveryAgent)
}
