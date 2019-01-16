package connectors

import (
	"net"

	"gopkg.in/alecthomas/kingpin.v2"
)

// DigitalOceanDNSConnector implements the connector interface and provides methods
// to interact with the DigitalOcean DNS API
type DigitalOceanDNSConnector struct {
}

// Register this connector in the list
func init() {
	RegisterConnector("digitalocean", &DigitalOceanDNSConnector{})
}

// Assert that we meet the interface at compile time
var _ Connector = (*DigitalOceanDNSConnector)(nil)

// Setup is used to setup the command line details and connect to the parser
func (do *DigitalOceanDNSConnector) Setup(app *kingpin.CmdClause) {
	app.Arg("", "ROO").String()
}

// AddNode creates a new record for that node on the DO DNS API
func (do *DigitalOceanDNSConnector) AddNode(address string, ip net.IP, ttl int) error {
	return nil
}

// UpdateNode updates a record for that node on the DO DNS API
func (do *DigitalOceanDNSConnector) UpdateNode(address string, ip net.IP, ttl int) error {
	return nil
}

// DeleteNode deletes the record for that node on the DO DNS API
func (do *DigitalOceanDNSConnector) DeleteNode(address string) error {
	return nil
}
