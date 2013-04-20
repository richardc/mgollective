package mgollective

import (
	"github.com/maruel/subcommands"
	"os"
)

var commands []*subcommands.Command

func RunApplication() {
	defer FlushLog()
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
