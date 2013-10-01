package mgollective

import (
	"code.google.com/p/go-commander"
	"flag"
	"github.com/golang/glog"
	"os"
)

var commands []*commander.Command
var client_config_file string = "client.cfg"
var server_config_file string = "server.cfg"

func RunApplication() {
	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", "0")

	global_flags := flag.NewFlagSet("global", flag.ExitOnError)
	cc_file := global_flags.String("client-config", "client.cfg", "usage")
	sc_file := global_flags.String("server-config", "server.cfg", "usage")
	global_flags.Parse(os.Args[1:])

	client_config_file = *cc_file
	server_config_file = *sc_file

	commander := &commander.Commander{
		Name:     "mgo",
		Commands: commands,
	}
	err := commander.Run(global_flags.Args())
	if err != nil {
		glog.Fatal(err)
	}
}

func RegisterCommand(command *commander.Command) {
	commands = append(commands, command)
}
