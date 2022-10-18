package main

import (
	"fmt"
	"os"

	"github.com/abdukhashimov/student_aggregator/internal/pkg/clitools"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("no subcommand")
		return
	}

	comm := clitools.PasreCommand()

	switch comm.Action {
	case "http":
		fmt.Printf("HTTP Port is %d \n", comm.Port)
	case "grpc":
		fmt.Printf("GRPC Port is %d \n", comm.Port)
	default:
		fmt.Println("Unknown subcommand: " + os.Args[1])
	}
}
