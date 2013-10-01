package mgollective

import (
	"code.google.com/p/go-commander"
	"fmt"
	"github.com/golang/glog"
	"strings"
)

func runRpcCommand(cmd *commander.Command, args []string) {
	mgo := NewClient()

	if len(args) < 2 {
		glog.Fatal("not enough arguments")
	}

	params := make(map[string]string)
	for _, arg := range args[2:] {
		value := strings.SplitN(arg, "=", 2)
		params[value[0]] = value[1]
	}

	request := RequestMessage{
		Body: RequestBody{
			Agent:  args[0],
			Action: args[1],
			Params: params,
		},
	}

	discovered_nodes := []string{mgo.Identity()}

	defer mgo.Shutdown()
	mgo.RpcCommand(request, discovered_nodes, func(message ResponseMessage) {
		fmt.Printf("%-40s %v\n", message.Headers["mc_identity"], message.Body)
	})
}

func init() {
	RegisterCommand(&commander.Command{
		UsageLine: "rpc",
		Run:       runRpcCommand,
	})
}
