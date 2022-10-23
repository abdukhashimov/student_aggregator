package clitools

import (
	"flag"
	"os"
	"strings"
)

type command struct {
	Action string
	Port   int
}

func PasreCommand() command {
	var res = command{}
	var subcommand = flag.NewFlagSet("", flag.ExitOnError)

	subcommand.IntVar(&res.Port, "port", -1, "listen on port")
	_ = subcommand.Parse(os.Args[2:])

	res.Action = strings.ToLower(os.Args[1])
	return res
}
