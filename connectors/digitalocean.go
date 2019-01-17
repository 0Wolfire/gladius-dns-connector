package connectors

import (
	"net"

	"context"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
	"gopkg.in/alecthomas/kingpin.v2"
)

// DigitalOceanDNSConnector implements the connector interface and provides methods
// to interact with the DigitalOcean DNS API
type DigitalOceanDNSConnector struct {
	client *godo.Client
	domain string
}

// Register this connector in the list
func init() {
	RegisterConnector("digitalocean", &DigitalOceanDNSConnector{})
}

// Assert that we meet the interface at compile time
var _ Connector = (*DigitalOceanDNSConnector)(nil)

// Setup is used to setup the command line details and connect to the parser
func (do *DigitalOceanDNSConnector) Setup(app *kingpin.CmdClause) {
	token := app.Flag("api_key", "The DigitalOcean API Key [Env: DO_API_KEY]").Envar("DO_API_KEY").Required().String()
	tokenSource := &TokenSource{
		AccessToken: *token,
	}

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	do.client = godo.NewClient(oauthClient)

	do.domain = *app.Flag("domain", "The domain for on DigitalOcean DNS [Env: DO_DOMAIN]").Envar("DO_DOMAIN").Required().String()

}

// AddNode creates a new record for that node on the DO DNS API
func (do *DigitalOceanDNSConnector) AddNode(address string, ip net.IP, ttl int) error {
	recordAdd := &godo.DomainRecordEditRequest{
		Type: "A",
		Data: ip.String(),
		Name: address,
	}
	ctx := context.TODO()
	do.client.Domains.CreateRecord(ctx, do.domain, recordAdd)

	// TODO: Return real error
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

// TokenSource is a type to store tokens for Oauth
type TokenSource struct {
	AccessToken string
}

// Token returns the Oauth token
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}
