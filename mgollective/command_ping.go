package mgollective

import (
	"fmt"
	"github.com/maruel/subcommands"
	"time"
)

type PingCommand struct {
	subcommands.CommandRunBase
}

func init() {
	RegisterCommand(&subcommands.Command{
		UsageLine:  "ping",
		CommandRun: func() subcommands.CommandRun { return &PingCommand{} },
	})
}

func (*PingCommand) Run(a subcommands.Application, args []string) int {
	start := time.Now()
	mgo := NewClient()

	pings := make([]time.Duration, 0)
	mgo.Discover(func(message Message) {
		ping := time.Since(start)
		pings = append(pings, ping)
		fmt.Printf("%-40s time=%s\n", message.Body.Senderid, ping.String())
	})

	if len(pings) == 0 {
		fmt.Println("No responses.")
		return 1
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
	return 0
}
