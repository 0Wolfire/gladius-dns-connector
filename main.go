package main

import (
	"fmt"
	"os"

	"github.com/gladiusio/gladius-dns-connector/connectors"
	"github.com/gladiusio/gladius-dns-connector/state"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app         = kingpin.New("gladns", "An application to map Gladius state to a DNS service")
	ls          = app.Command("list", "List available connectors").Action(printConnectors)
	gatewayIP   = app.Flag("gateway_ip", "The IP to connect to for the gladius network gateway").Default("127.0.0.1").IP()
	gatewayPort = app.Flag("gateway_port", "The port to connect to for the gladius network gateway").Default("3001").Uint16()
)

func printConnectors(c *kingpin.ParseContext) error {
	fmt.Println("Available connectors:")
	for _, c := range connectors.List() {
		fmt.Println("- " + c)
	}
	fmt.Println("\ntype \"help <connector>\" for more info")
	return nil
}

func main() {
	p := state.NewParser(gatewayIP.String(), *gatewayPort)

	// Regiser connector commands
	for _, name := range connectors.List() {
		command := app.Command(name, "Confgiure DNS with the "+name+" connector")
		connectors.GetConnector(name).Setup(command)
	}

	// Tell our parser to use the selected connector
	chosen := kingpin.MustParse(app.Parse(os.Args[1:]))
	if chosen != "list" {
		p.SetConnector(connectors.GetConnector(chosen))
	}
}
