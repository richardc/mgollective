package mgollective

import (
	"fmt"
	"os"
)

type Command interface {
	Run()
}

var commandRegistry = map[string]func() Command{}

func registerCommand(name string, command func() Command) {
	commandRegistry[name] = command
}

func RunCommand(name string) {
	command, ok := commandRegistry[name]
	if !ok {
		fmt.Println("unrecognised command")
		os.Exit(1)
	}
	command().Run()
}
