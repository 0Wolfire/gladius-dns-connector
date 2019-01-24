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
	app            = kingpin.New("gladns", "An application to map Gladius state to a DNS service")
	ls             = app.Command("list", "List available connectors").Action(printConnectors)
	tickRate       = app.Flag("tick_rate", "How often to query the state and update DNS records").Default("5s").Duration()
	logPretty      = app.Flag("log_pretty", "Whether or not to use pretty output or JSON").Default("false").Bool()
	gatewayAddress = app.Flag("gateway_address", "The base address to connect to for the gladius network gateway").Default("http://localhost:3001").URL()
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

	var p *state.Parser

	// Get the chosen connector and start it
	if connectors.Exists(chosen) {
		c := connectors.GetConnector(chosen)

		err := c.Connect()
		if err != nil {
			log.Fatal().Err(err).Str("name", chosen).Msg("Error conecting to DNS connector")
		}
		p = state.NewParser(*gatewayAddress, *tickRate, c)
	}

	// Start the state parser
	err := p.Start()
	if err != nil {
		log.Error().Err(err).Msg("Error parsing state")
	}
}
