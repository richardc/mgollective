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
	global_flags := flag.NewFlagSet("mgo", flag.ExitOnError)
	var debug bool
	global_flags.BoolVar(
		&debug,
		"debug",
		false,
		"be very verbose",
	)
	global_flags.StringVar(
		&client_config_file,
		"client-config",
		"client.cfg",
		"specify client config file",
	)
	global_flags.StringVar(
		&server_config_file,
		"server-config",
		"server.cfg",
		"specify server config file",
	)
	global_flags.Parse(os.Args[1:])

	if debug {
		/// tell glog to turn up the noise
		flag.Set("logtostderr", "true")
	}

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
