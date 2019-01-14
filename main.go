package main

import (
	"fmt"
	"os"

	"github.com/gladiusio/gladius-dns-connector/connectors"
)

func main() {
	firstArg := os.Args[1]
	if firstArg == "help" {
		fmt.Print("Gladius DNS Connector - Help\n\nUsage - connector [connector name] [options]\n")
		fmt.Print("Type \"connector list\" to see all connectors.\n")
	} else if firstArg == "list" {
		fmt.Print("Available connectors:\n\n")
		for _, c := range connectors.List() {
			fmt.Println("- " + c)
		}
	} else if connectors.Exists(firstArg) {
		c := connectors.GetConnector(firstArg)
		c.Connect(os.Args[2:])
		// TODO: Pass connector into state logic
	} else {
		fmt.Println()
	}
}
