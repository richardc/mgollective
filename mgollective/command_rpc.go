package mgollective

import (
	"code.google.com/p/go-commander"
	"fmt"
)

func runRpcCommand(cmd *commander.Command, args []string) {
	mgo := NewClient()

	request := RequestMessage{
		Body: RequestBody{
			Agent:  args[0],
			Action: args[1],
			Params: make(map[string]string),
		},
	}

	discovered_nodes := []string{mgo.Identity()}

	defer mgo.Shutdown()
	mgo.RpcCommand(request, discovered_nodes, func(message ResponseMessage) {
		fmt.Printf("%-40s %s\n", message.Headers["mc_identity"], message.Body["timestamp"])
	})
}

func init() {
	RegisterCommand(&commander.Command{
		UsageLine: "rpc",
		Run:       runRpcCommand,
	})
}
