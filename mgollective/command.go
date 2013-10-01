package mgollective

import (
	"code.google.com/p/go-commander"
	"flag"
	"github.com/golang/glog"
	"os"
)

var commands []*commander.Command
var client_config_file string = "client.cfg"

func RunApplication() {
	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", "0")

	global_flags := flag.NewFlagSet("global", flag.ExitOnError)
	config_file := global_flags.String("config", "client.cfg", "usage")
	global_flags.Parse(os.Args[1:])

	client_config_file = *config_file

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
