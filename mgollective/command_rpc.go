package mgollective

import (
	"fmt"
	"github.com/maruel/subcommands"
)

type RpcCommand struct {
	subcommands.CommandRunBase
}

func init() {
	RegisterCommand(&subcommands.Command{
		UsageLine:  "rpc",
		CommandRun: func() subcommands.CommandRun { return &RpcCommand{} },
	})
}

func (*RpcCommand) Run(a subcommands.Application, args []string) int {
	mgo := NewClient()

	request := RequestMessage{
		Body: RequestBody{
			Agent:  args[0],
			Action: args[1],
			Params: make(map[string]string),
		},
	}

	discovered_nodes := []string{"foo"}

	defer mgo.Shutdown()
	mgo.RpcCommand(request, discovered_nodes, func(message ResponseMessage) {
		fmt.Printf("%-40s %s\n", message.Headers["mc_identity"], message.Body["timestamp"])
	})

	return 0
}
