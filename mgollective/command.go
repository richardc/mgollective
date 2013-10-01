package mgollective

import (
	"code.google.com/p/go-commander"
	"flag"
	"github.com/golang/glog"
	"os"
)

var commands []*commander.Command
var client_config_file string
var server_config_file string

func RunApplication() {
	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", "0")

	global_flags := flag.NewFlagSet("global", flag.ExitOnError)
	global_flags.StringVar(&client_config_file, "client-config", "client.cfg", "usage")
	global_flags.StringVar(&server_config_file, "server-config", "server.cfg", "usage")
	global_flags.Parse(os.Args[1:])

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
