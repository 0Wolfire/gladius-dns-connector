package connectors

import (
	"errors"
	"net"

	"context"

	"github.com/rs/zerolog/log"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
	"gopkg.in/alecthomas/kingpin.v2"
)

// DigitalOceanDNSConnector implements the connector interface and provides methods
// to interact with the DigitalOcean DNS API
type DigitalOceanDNSConnector struct {
	client    *godo.Client
	domain    string
	cdnDomain string
	token     string
	idMap     map[string]int
}

// Register this connector in the list
func init() {
	RegisterConnector("digitalocean", &DigitalOceanDNSConnector{})
}

// Assert that we meet the interface at compile time
var _ Connector = (*DigitalOceanDNSConnector)(nil)

// Setup is used to setup the command line details
func (do *DigitalOceanDNSConnector) Setup(app *kingpin.CmdClause) {
	app.Flag("api_key", "The DigitalOcean API Key [Env: DO_API_KEY]").Envar("DO_API_KEY").Required().StringVar(&do.token)
	app.Flag("domain", "The domain for on DigitalOcean DNS [Env: DO_DOMAIN]").PlaceHolder("yourepool.com").Envar("DO_DOMAIN").Required().StringVar(&do.domain)
	app.Flag("cdn_subdomain", "The cdn subdomain for nodes [Env: DO_CDN_SUBDOMAIN]").Default("cdn").Envar("DO_CDN_SUBDOMAIN").StringVar(&do.cdnDomain)
}

// Connect connects to the DO API
func (do *DigitalOceanDNSConnector) Connect() error {
	tokenSource := &TokenSource{
		AccessToken: do.token,
	}

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	do.client = godo.NewClient(oauthClient)

	// Check to see if we have a valid API key
	_, _, err := do.client.Account.Get(context.TODO())
	if err != nil {
		return err
	}

	// Create a list to hold our domain records
	list := []godo.DomainRecord{}

	// Create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		records, resp, err := do.client.Domains.Records(context.TODO(), do.domain, opt)
		if err != nil {
			return err
		}

		// Append the current page's droplets to our list
		for _, d := range records {
			list = append(list, d)
		}

		// If we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return err
		}

		// Set the page we want for the next request
		opt.Page = page + 1
	}

	// Populate our id map for updates
	do.idMap = make(map[string]int)
	for _, r := range list {
		do.idMap[r.Name] = r.ID
	}

	log.Debug().Interface("records_map", do.idMap).Msg("Loaded record IDs")

	return nil
}

// AddNode creates a new record for that node on the DO DNS API
func (do *DigitalOceanDNSConnector) AddNode(address string, ip net.IP, ttl int) error {
	recordAdd := &godo.DomainRecordEditRequest{
		Type: "A",
		Data: ip.String(),
		Name: do.makeName(address),
	}
	ctx := context.TODO()
	r, _, err := do.client.Domains.CreateRecord(ctx, do.domain, recordAdd)

	// Update our record map
	do.idMap[r.Name] = r.ID
	return err
}

// UpdateNode updates a record for that node on the DO DNS API
func (do *DigitalOceanDNSConnector) UpdateNode(address string, ip net.IP, ttl int) error {
	recordUpdate := &godo.DomainRecordEditRequest{
		Type: "A",
		Data: ip.String(),
		Name: do.makeName(address),
	}
	ctx := context.TODO()

	// Get the correct record ID from the node address
	id := do.idMap[do.makeName(address)]
	_, _, err := do.client.Domains.EditRecord(ctx, do.domain, id, recordUpdate)

	return err
}

// DeleteNode deletes the record for that node on the DO DNS API
func (do *DigitalOceanDNSConnector) DeleteNode(address string) error {
	if _, exists := do.idMap[address]; !exists {
		return errors.New("record does not exist: " + do.makeName(address))
	}

	_, err := do.client.Domains.DeleteRecord(context.TODO(), do.domain, do.idMap[address])
	return err
}

func (do *DigitalOceanDNSConnector) makeName(address string) string {
	return address + "." + do.cdnDomain
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
