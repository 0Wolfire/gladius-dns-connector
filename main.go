package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/gladiusio/gladius-dns-connector/connectors"
	"github.com/gladiusio/gladius-dns-connector/state"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app         = kingpin.New("gladns", "An application to map Gladius state to a DNS service")
	ls          = app.Command("list", "List available connectors").Action(printConnectors)
	logPretty   = app.Flag("log_pretty", "Whether or not to use pretty output or JSON").Default("false").Bool()
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

	chosen := kingpin.MustParse(app.Parse(os.Args[1:]))

	// Configure the logger
	if *logPretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	//
	if chosen != "list" {
		err := connectors.GetConnector(chosen).Connect()
		if err != nil {
			log.Fatal().Err(err).Str("name", chosen).Msg("Error conecting to DNS connector")
		}
		p.SetConnector(connectors.GetConnector(chosen))
	}

}
