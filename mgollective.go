package mgollective

import (
	"fmt"
	"os"
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

func PingLoop() {
	start := time.Now()
	config := getconfig("mgo.conf", true)
	connectorname := config.GetStringDefault("connector", "class", "redis")
	var connector Connector
	if factory, exists := connectorRegistry[connectorname]; exists {
		connector = factory(config)
	} else {
		fmt.Printf("No connector called %s", connectorname)
		os.Exit(1)
	}

	connector.Connect()
	connector.Subscribe()

	pings := make([]time.Duration, 0)
	// Discover should be a method on *something*.  Probably want to refactor config
	Discover(connector, *config, func(message Message) {
		ping := time.Since(start)
		pings = append(pings, ping)
		fmt.Printf("%-40s time=%s\n", message.body.Senderid, ping.String())
	})

	var min, max, sum time.Duration
	min = pings[0]
	for _, ping := range pings {
		sum += ping
		if ping > max {
			max = ping
		}
		if ping < min {
			min = ping
		}
	}

	mean := time.Duration(int64(sum) / int64(len(pings)))
	fmt.Println()
	fmt.Println("--- ping statistics ---")
	fmt.Printf("%d replies max: %s min: %s avg: %s\n",
		len(pings), max.String(), min.String(), mean.String())
}

func DaemonLoop() {
	config := getconfig("mgo.conf", false)
	connectorname := config.GetStringDefault("connector", "class", "redis")
	var connector Connector
	if factory, exists := connectorRegistry[connectorname]; exists {
		connector = factory(config)
	} else {
		fmt.Printf("No connector called %s", connectorname)
		os.Exit(1)
	}

	connector.Connect()
	connector.Subscribe()

	ch := make(chan Message)
	go connector.Loop(ch)
	for {
		message := <-ch
		logger.Debugf("Recieved %+v", message)
		if agent, exists := agentRegistry[message.body.Agent]; exists {
			agent(config).Respond(message, connector)
		} else {
			logger.Debugf("No agent '%s'", message.body.Agent)
		}
	}
}
