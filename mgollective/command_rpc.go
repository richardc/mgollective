package mgollective

import (
	"code.google.com/p/go-commander"
	"encoding/json"
	"flag"
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
		if len(value) > 1 {
			params[value[0]] = value[1]
		} else {
			params[arg] = ""
		}
	}

	request := RequestMessage{
		Body: RequestBody{
			Agent:  args[0],
			Action: args[1],
			Params: params,
		},
	}

	discovered_nodes := []string{mgo.Identity()}

	var json_mode bool
	if cmd.Flag.Lookup("j").Value.String() == "true" {
		json_mode = true
	}
	json_data := make([]interface{}, 0)

	defer mgo.Shutdown()
	mgo.RpcCommand(request, discovered_nodes, func(message ResponseMessage) {
		if json_mode {
			record := map[string]interface{}{
				"agent":  args[0],
				"action": args[1],
				"sender": message.Headers["mc_identity"],
				"data":   message.Body,
			}

			json_data = append(json_data, record)
		} else {
			fmt.Printf("%-40s %v\n", message.Headers["mc_identity"], message.Body)
		}
	})

	if json_mode {
		bytes, _ := json.MarshalIndent(json_data, "  ", "  ")
		fmt.Print(string(bytes))
	}
}

func init() {
	cmd := &commander.Command{
		UsageLine: "rpc",
		Run:       runRpcCommand,
		Flag:      *flag.NewFlagSet("mgo-rpc", flag.ExitOnError),
	}
	cmd.Flag.Bool("j", false, "output as json")
	RegisterCommand(cmd)
}
