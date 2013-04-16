package mgollective

import (
	"time"
)

func Discover(connector Connector, config Config, callback func(Message)) {
	discovery := Message{
		target:   "mcollective::server::agents",
		reply_to: config.identity(),
		body: MessageBody{
			Agent:      "discovery",
			Body:       "ping",
			Collective: "mcollective",
			Callerid:   config.callerid(),
			Senderid:   config.senderid(),
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
