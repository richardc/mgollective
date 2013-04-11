package main

import (
	"flag"
	"fmt"
	"github.com/richardc/mgollective"
)

func main() {
	flag.Parse()
	command := flag.Arg(0)
	if command == "ping" {
		mgollective.PingLoop()
	} else if command == "daemon" {
		mgollective.DaemonLoop()
	} else {
		fmt.Println("unrecognised command")
	}
}
