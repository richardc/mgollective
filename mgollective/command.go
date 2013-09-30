package mgollective

import (
	"flag"
	"github.com/maruel/subcommands"
	"os"
)

var commands []*subcommands.Command

func RunApplication() {
	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", "0")
	mgollective := &subcommands.DefaultApplication{
		Name:     "mgo",
		Title:    "mgollective",
		Commands: append(commands, subcommands.CmdHelp),
	}
	subcommands.Run(mgollective, os.Args[1:])
}

func RegisterCommand(command *subcommands.Command) {
	commands = append(commands, command)
}
