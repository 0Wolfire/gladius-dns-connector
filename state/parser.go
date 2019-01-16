package state

import (
	"errors"

	"github.com/gladiusio/gladius-dns-connector/connectors"
)

// Parser processes the JSON state and exposes methods to view it
type Parser struct{}

// NewParser returns a parser connected to the specified Gladius Network Gateway
func NewParser(ip string, port uint16) *Parser {
	return &Parser{}
}

// SetConnector sets the connector and starts interacting with it
func (p *Parser) SetConnector(c connectors.Connector) {
	if c == nil {
		panic(errors.New("connector cannot be nil"))
	}
}
