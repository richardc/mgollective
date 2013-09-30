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
	agent := "rpcutil"
	command := "ping"
	params := make(map[string]string)

	defer mgo.Shutdown()
	mgo.RpcCommand(agent, command, params, func(message ResponseMessage) {
		fmt.Printf("%-40s %s\n", message.Headers["mc_identity"], message.Body["timestamp"])
	})

	return 0
}
