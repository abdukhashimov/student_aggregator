package main

import (
	"fmt"
)

func getGreeting() string {
	return "Hello, world"
}

func main() {
	greeting := getGreeting()
	fmt.Println(greeting)
}
