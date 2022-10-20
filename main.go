package main

import (
	"fmt"

	"github.com/abdukhashimov/student_aggregator/internal/config"
)

func getGreeting() string {
	return "Hello, world"
}

func main() {
	greeting := getGreeting()
	fmt.Println(greeting)

	cfg := config.Load()

	// TODO: remove
	fmt.Println(cfg)
}
