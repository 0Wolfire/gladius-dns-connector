package connectors

import (
	"net"

	"gopkg.in/alecthomas/kingpin.v2"
)

// Connector is an interface for exposing methods to interact with a DNS service
type Connector interface {
	AddNode(address string, ip net.IP, ttl int) error
	UpdateNode(address string, ip net.IP, ttl int) error
	DeleteNode(address string) error
	Setup(*kingpin.CmdClause)
	Connect() error
}

var connectors map[string]Connector

// RegisterConnector registers a connector globally so it can be accessed in the CLI
func RegisterConnector(name string, c Connector) {
	if connectors == nil {
		connectors = make(map[string]Connector)
	}

	connectors[name] = c
}

// GetConnector returns the connector assosiated with the name, nil if not found
func GetConnector(name string) Connector {
	return connectors[name]
}

// List returns all of the names of the registered connectors
func List() []string {
	toReturn := make([]string, 0)
	for n := range connectors {
		toReturn = append(toReturn, n)
	}

	return toReturn
}

// Exists checks to see if there is a connector by that name registered
func Exists(name string) bool {
	_, exists := connectors[name]
	return exists
}
