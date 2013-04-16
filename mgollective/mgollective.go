package mgollective

import (
	"time"
)

func Discover(connector Connector, config Config, callback func(Message)) {
	discovery := Message{
		Target:   "mcollective::server::agents",
		Reply_to: config.Identity(),
		Body: MessageBody{
			Agent:      "discovery",
			Body:       "ping",
			Collective: "mcollective",
			Callerid:   config.Callerid(),
			Senderid:   config.Senderid(),
			Ttl:        60,
			Msgtime:    time.Now().Unix(),
			Requestid:  "42",
		},
	}

	cb := make(chan Message)
	go connector.Loop(cb)
	connector.Publish(discovery)

	for {
		select {
		case message := <-cb:
			callback(message)
		case <-time.After(3 * time.Second):
			return
		}
	}
}
