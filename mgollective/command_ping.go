package mgollective

import (
	"code.google.com/p/go-commander"
	"fmt"
	"time"
)

func runPingCommand(cmd *commander.Command, args []string) {
	start := time.Now()
	mgo := NewClient()

	pings := make([]time.Duration, 0)
	mgo.Discover(func(message ResponseMessage) {
		ping := time.Since(start)
		pings = append(pings, ping)
		fmt.Printf("%-40s time=%s\n", message.Headers["sender-id"], ping.String())
	})
	defer mgo.Shutdown()

	if len(pings) == 0 {
		fmt.Println("No responses.")
		return
	}

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

func init() {
	RegisterCommand(&commander.Command{
		UsageLine: "ping",
		Run:       runPingCommand,
	})
}
