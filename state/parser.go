package state

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gladiusio/gladius-dns-connector/connectors"
	"github.com/rs/zerolog/log"
)

// Parser processes the JSON state and exposes methods to view it
type Parser struct {
	connector  connectors.Connector
	gatewayURL *url.URL
}

// NewParser returns a parser connected to the specified Gladius Network Gateway
func NewParser(url *url.URL, c connectors.Connector) *Parser {
	return &Parser{
		connector:  c,
		gatewayURL: url,
	}
}

// Start starts the state parsing and connector calls
func (p *Parser) Start() error {
	c := time.Tick(5 * time.Second)
	for range c {
		resp, err := http.Get(p.gatewayURL.String() + "/api/p2p/state")
		if err != nil {
			log.Error().Err(err).Msg("Error connecting to gateway")
			return err
		}

		err = p.processResponse(resp)
		if err != nil {
			log.Error().Err(err).Msg("Error processing response from gateway")
			return err
		}
	}

	return nil
}

func (p *Parser) processResponse(resp *http.Response) error {
	return nil
}
