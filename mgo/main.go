package main

import (
	"flag"
	"github.com/richardc/mgollective"
)

func main() {
	defer mgollective.FlushLog()
	flag.Parse()
	command := flag.Arg(0)
	mgollective.RunCommand(command)
}
