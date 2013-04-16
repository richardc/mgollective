package mgollective

import (
	"fmt"
	"os"
	"time"
)

type PingCommand struct {
}

func makePingCommand() Command {
	return &PingCommand{}
}

func init() {
	registerCommand("ping", makePingCommand)
}

func (*PingCommand) Run() {
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
