package main

import (
	"os"

	"github.com/gladiusio/gladius-dns-connector/connectors"
	"github.com/gladiusio/gladius-dns-connector/state"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app         = kingpin.New("gladns", "An application to map Gladius state to a DNS service")
	gatewayIP   = app.Flag("gateway_ip", "The IP to connect to for the gladius network gateway").Default("127.0.0.1").IP()
	gatewayPort = app.Flag("gateway_port", "The port to connect to for the gladius network gateway").Default("3001").Uint16()
)

func main() {
	p := state.NewParser(gatewayIP.String(), *gatewayPort)

	// Regiser connector commands
	for _, name := range connectors.List() {
		command := app.Command(name, "Use the "+name+" connector")
		connectors.GetConnector(name).Setup(command)
	}

	// Tell our parser to use the selected connector
	chosen := kingpin.MustParse(app.Parse(os.Args[1:]))
	p.SetConnector(connectors.GetConnector(chosen))
}
