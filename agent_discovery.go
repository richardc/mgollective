package mgollective

import (
	"log"
)

type DiscoveryAgent struct {
}

func (a *DiscoveryAgent) Respond(msg *Message, connector *Connector) {
	log.Printf("Discover agent handling %+v", msg)
}

func makeDiscoveryAgent(c *Config) Agent {
	return &DiscoveryAgent{}
}

func init() {
	registerAgent("discovery", makeDiscoveryAgent)
}
