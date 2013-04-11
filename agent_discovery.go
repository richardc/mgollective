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
	reply := make(map[string]interface{})

	reply["target"] = msg.reply_to
	reply[":requestid"] = msg.Requestid
	reply[":senderagent"] = "discovery"
	reply[":senderid"] = a.config.senderid()
	reply[":msgtime"] = time.Now().Unix()
	if msg.Body == "ping" {
		reply["body"] = "pong"
	} else {
		reply["body"] = "Unknown Request: " + msg.Body
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
