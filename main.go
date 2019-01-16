package main

import (
	"os"

	"github.com/gladiusio/gladius-dns-connector/connectors"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app         = kingpin.New("gladns", "A command-line chat application.")
	gatewayIP   = app.Flag("gateway_ip", "The IP to connect to for the gladius network gateway").Default("127.0.0.1").IP()
	gatewayPort = app.Flag("gateway_port", "The port to connect to for the gladius network gateway").Default("3001").Uint16()
)

func main() {
	// Regiser connector commands
	for _, name := range connectors.List() {
		command := app.Command(name, "Use the "+name+" connector")
		connectors.GetConnector(name).SetupCommand(command)
	}

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
