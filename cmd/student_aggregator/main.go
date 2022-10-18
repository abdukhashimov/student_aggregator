package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("no subcommand")
		return
	}

	var port int

	subcommand := flag.NewFlagSet("", flag.ExitOnError)

	subcommand.IntVar(&port, "port", -1, "listen on port")
	subcommand.Parse(os.Args[2:])

	comm := strings.ToLower(os.Args[1])

	switch comm {
	case "http":
		fmt.Printf("HTTP Port is %d \n", port)
	case "grpc":
		fmt.Printf("GRPC Port is %d \n", port)
	default:
		fmt.Println("Unknown subcommand: " + os.Args[1])
	}
}
