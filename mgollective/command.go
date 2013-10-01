package mgollective

import (
	"code.google.com/p/go-commander"
	"flag"
	"github.com/golang/glog"
	"os"
)

var commands []*commander.Command

func RunApplication() {
	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", "0")
	mgollective := &commander.Commander{
		Name:     "mgo",
		Commands: commands,
	}
	err := mgollective.Run(os.Args[1:])
	if err != nil {
		glog.Fatal(err)
	}
}

func RegisterCommand(command *commander.Command) {
	commands = append(commands, command)
}
